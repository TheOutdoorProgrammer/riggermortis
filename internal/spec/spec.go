// SPDX-License-Identifier: Apache-2.0

// Package spec loads riggermortis records and indexes them for cross-record
// checks. Shape and enum validation belong to CUE; this package exists for the
// rules CUE cannot express, which are the ones that need to resolve references
// or walk a graph.
package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Names struct {
	Canonical     string   `yaml:"canonical"`
	Aliases       []string `yaml:"aliases"`
	TrademarkNote string   `yaml:"trademark_note"`
}

type Pin struct {
	ID             string  `yaml:"id"`
	Type           string  `yaml:"type"`
	WireDiameterMM float64 `yaml:"wire_diameter_mm"`
}

type Variants struct {
	Axis   string    `yaml:"axis"`
	Values []float64 `yaml:"values"`
}

type Connects struct {
	From string   `yaml:"from"`
	To   []string `yaml:"to"`
}

type Caution struct {
	Severity string   `yaml:"severity"`
	Note     string   `yaml:"note"`
	Sources  []string `yaml:"sources"`
}

type Claim struct {
	ResidualPct *float64 `yaml:"residual_pct"`
	N           *int     `yaml:"n"`
	Source      string   `yaml:"source"`
	Tier        string   `yaml:"tier"`
	Rank        string   `yaml:"rank"`
}

type Strength struct {
	Claims []Claim `yaml:"claims"`
}

type Diameter struct {
	BreakingLoadKG float64 `yaml:"breaking_load_kg"`
	DiameterMM     float64 `yaml:"diameter_mm"`
	Source         string  `yaml:"source"`
}

// Node carries every variant-axis field a rig may select. Only the one
// matching the referenced component's axis is meaningful, which is rule 32.
type Node struct {
	ID             string  `yaml:"id"`
	Ref            string  `yaml:"ref"`
	Type           string  `yaml:"type"`
	Role           string  `yaml:"role"`
	Material       string  `yaml:"material"`
	BreakingLoadKG float64 `yaml:"breaking_load_kg"`
	MassG          float64 `yaml:"mass_g"`
	GapMM          float64 `yaml:"gap_mm"`
	DiameterMM     float64 `yaml:"diameter_mm"`
	RatingKG       float64 `yaml:"rating_kg"`

	// Scalar on a component node (a variant selection), a two-element range on
	// a line node (a cut length). Same name, two meanings, disambiguated here.
	LengthMM any `yaml:"length_mm"`
}

type Travel struct {
	TowardRod      string `yaml:"toward_rod"`
	TowardTerminal string `yaml:"toward_terminal"`
}

type Edge struct {
	From    string  `yaml:"from"`
	To      string  `yaml:"to"`
	Rel     string  `yaml:"rel"`
	Knot    string  `yaml:"knot"`
	Rigging string  `yaml:"rigging"`
	Pin     string  `yaml:"pin"`
	Travel  *Travel `yaml:"travel"`
}

type Point [2]float64

// Segment is a run of cord between crossings. Paint order lives here rather
// than on the cord because a square knot alternates over and under.
type Segment struct {
	Cord   string  `yaml:"cord"`
	Z      int     `yaml:"z"`
	Points []Point `yaml:"points"`
	// Stage that introduces this piece, so the renderer can call out what the
	// current step does while still drawing the whole knot. Zero means the
	// piece is structural rather than belonging to any one step.
	Stage int `yaml:"stage,omitempty"`
}

type StageGeometry struct {
	Stage    int       `yaml:"stage"`
	Segments []Segment `yaml:"segments"`
	// Where the working end travels next. A knot diagram teaches by showing
	// motion, not by showing a sequence of finished positions.
	Arrow *Arrow `yaml:"arrow,omitempty"`
	// Ends marks each working end so a reader can track which is which.
	Ends []Point `yaml:"ends,omitempty"`
}

type Arrow struct {
	From Point `yaml:"from"`
	Ctrl Point `yaml:"ctrl"`
	To   Point `yaml:"to"`
}

type Geometry struct {
	Width  float64         `yaml:"width"`
	Height float64         `yaml:"height"`
	Cords  []string        `yaml:"cords"`
	Stages []StageGeometry `yaml:"stages"`
}

type Action struct {
	Verb string `yaml:"verb"`
	// Subject is what the verb acts on: one entity or several. Records write
	// both forms, so it stays untyped here and Subjects normalises it.
	Subject   any    `yaml:"subject"`
	Names     string `yaml:"names"`
	Through   string `yaml:"through"`
	Chirality string `yaml:"chirality"`
	// Direction is which way a reeve passes through its loop, "U-D" or "D-U".
	// It is the sign of the normal the end travels along, and it is the only
	// thing distinguishing a tuck from the same tuck made backwards.
	Direction string `yaml:"direction"`
	Rotation  string `yaml:"rotation"`
	Force     string `yaml:"force"`
	Repeat    any    `yaml:"repeat"`
}

// Subjects normalises the subject to a list. A scalar is a list of one.
func (a Action) Subjects() []string {
	switch t := a.Subject.(type) {
	case string:
		return []string{t}
	case []any:
		out := make([]string, 0, len(t))
		for _, v := range t {
			if s, ok := v.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

type Stage struct {
	ID       int      `yaml:"id"`
	Prose    string   `yaml:"prose"`
	Notation string   `yaml:"notation"`
	Actions  []Action `yaml:"actions"`
}

type Validation struct {
	Status string `yaml:"status"`
}

// Record is a superset of every kind. CUE has already rejected any record
// whose fields do not belong to its kind, so a flat struct here trades a
// little precision for a great deal less machinery.
type Record struct {
	Kind          string     `yaml:"kind"`
	ID            string     `yaml:"id"`
	SchemaVersion int        `yaml:"schema_version"`
	Name          string     `yaml:"name"`
	Names         *Names     `yaml:"names"`
	Validation    Validation `yaml:"validation"`

	Mounting      string    `yaml:"mounting"`
	BlocksPassage bool      `yaml:"blocks_passage"`
	BoreMM        float64   `yaml:"bore_mm"`
	Soft          bool      `yaml:"soft"`
	Pins          []Pin     `yaml:"pins"`
	Variants      *Variants `yaml:"variants"`

	Role      string     `yaml:"role"`
	Connects  *Connects  `yaml:"connects"`
	Cautions  []Caution  `yaml:"cautions"`
	Strength  *Strength  `yaml:"strength"`
	Diameters []Diameter `yaml:"diameters"`

	Stages   []Stage   `yaml:"stages"`
	Geometry *Geometry `yaml:"geometry"`

	Nodes []Node `yaml:"nodes"`
	Edges []Edge `yaml:"edges"`

	LineTypes    []string `yaml:"line_types"`
	FailureModes []string `yaml:"failure_modes"`

	Path string `yaml:"-"`
}

// Set is every record, indexed by id.
type Set struct {
	ByID map[string]*Record
	All  []*Record
}

// dirKind enforces rule 7: directory determines kind.
var dirKind = map[string]string{
	"components": "component",
	"lines":      "line",
	"knots":      "knot",
	"riggings":   "rigging",
	"rigs":       "rig",
	"sources":    "source",
	"patterns":   "pattern",
	"techniques": "technique",
	"species":    "species",
}

func Load(root string) (*Set, error) {
	set := &Set{ByID: map[string]*Record{}}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".yaml" {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var r Record
		if err := yaml.Unmarshal(raw, &r); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		r.Path = path
		if prev, dup := set.ByID[r.ID]; dup {
			return fmt.Errorf("%s: id %q already used by %s", path, r.ID, prev.Path)
		}
		set.ByID[r.ID] = &r
		set.All = append(set.All, &r)
		return nil
	})
	return set, err
}

// ExpectedKind returns the kind a path's directory implies, and whether the
// directory is one the spec knows about.
func ExpectedKind(path string) (string, bool) {
	k, ok := dirKind[filepath.Base(filepath.Dir(path))]
	return k, ok
}

// Category is the part of an id before the dot.
func Category(id string) string {
	if i := strings.IndexByte(id, '.'); i > 0 {
		return id[:i]
	}
	return ""
}

// Pin finds a pin by id on a component.
func (r *Record) Pin(id string) *Pin {
	for i := range r.Pins {
		if r.Pins[i].ID == id {
			return &r.Pins[i]
		}
	}
	return nil
}

// SelectedVariant returns the value a rig node picked along the given axis.
func (n *Node) SelectedVariant(axis string) (float64, bool) {
	switch axis {
	case "mass_g":
		return n.MassG, n.MassG != 0
	case "gap_mm":
		return n.GapMM, n.GapMM != 0
	case "diameter_mm":
		return n.DiameterMM, n.DiameterMM != 0
	case "rating_kg":
		return n.RatingKG, n.RatingKG != 0
	case "breaking_load_kg":
		return n.BreakingLoadKG, n.BreakingLoadKG != 0
	case "length_mm":
		v, ok := n.LengthMM.(float64)
		if !ok {
			if i, isInt := n.LengthMM.(int); isInt {
				v, ok = float64(i), true
			}
		}
		return v, ok
	}
	return 0, false
}
