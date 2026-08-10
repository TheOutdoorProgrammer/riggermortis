// SPDX-License-Identifier: Apache-2.0

// Package rules implements the checks CUE cannot express: anything that
// resolves a reference or walks a graph.
//
// Every rule carries the ID it has in docs/spec.md so the two can be checked
// against each other, per ADR 0003.
package rules

import (
	"fmt"
	"sort"

	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
)

type Finding struct {
	Rule     string
	Severity Severity
	Path     string
	Message  string
}

type Rule struct {
	ID       string
	Severity Severity
	Summary  string
	Check    func(*spec.Set) []Finding
}

var registry []Rule

func register(r Rule) { registry = append(registry, r) }

// All returns every registered rule, ordered by ID.
func All() []Rule {
	out := append([]Rule(nil), registry...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func Run(s *spec.Set) []Finding {
	var all []Finding
	for _, r := range All() {
		for _, f := range r.Check(s) {
			f.Rule, f.Severity = r.ID, r.Severity
			all = append(all, f)
		}
	}
	return all
}

func find(r *spec.Record, format string, a ...any) Finding {
	return Finding{Path: r.Path, Message: fmt.Sprintf(format, a...)}
}

func init() {
	register(Rule{
		ID: "R001", Severity: Error,
		Summary: "a threaded component must be stopped before every free end",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				if r.Kind != "rig" {
					continue
				}
				nodes := map[string]*spec.Node{}
				for i := range r.Nodes {
					nodes[r.Nodes[i].ID] = &r.Nodes[i]
				}
				for _, e := range r.Edges {
					if e.Rel != "threaded" {
						continue
					}
					if e.Travel == nil {
						out = append(out, find(r, "%s: threaded %q declares no travel bounds", r.ID, e.To))
						continue
					}
					// The rod-ward end of a main line is not free; it continues
					// to the reel. Every other open end lets the component off.
					seg := nodes[e.From]
					if e.Travel.TowardRod == "open" && (seg == nil || seg.Role != "main-line") {
						out = append(out, find(r, "%s: threaded %q is open toward a free rod-side end", r.ID, e.To))
					}
					out = append(out, checkStop(s, r, nodes, e.To, "terminal", e.Travel.TowardTerminal)...)
					if e.Travel.TowardRod != "open" {
						out = append(out, checkStop(s, r, nodes, e.To, "rod", e.Travel.TowardRod)...)
					}
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R002", Severity: Error,
		Summary: "a knot must accept the pin type it ties to",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				if r.Kind != "rig" {
					continue
				}
				nodes := map[string]*spec.Node{}
				for i := range r.Nodes {
					nodes[r.Nodes[i].ID] = &r.Nodes[i]
				}
				for _, e := range r.Edges {
					if e.Rel != "tied" || e.Knot == "" || e.Pin == "" {
						continue
					}
					knot, target := s.ByID[e.Knot], nodes[e.To]
					if knot == nil || knot.Connects == nil || target == nil {
						continue
					}
					comp := s.ByID[target.Ref]
					if comp == nil {
						continue
					}
					pin := comp.Pin(e.Pin)
					if pin == nil {
						out = append(out, find(r, "%s: %q has no pin %q", r.ID, target.Ref, e.Pin))
						continue
					}
					if !contains(knot.Connects.To, pin.Type) {
						out = append(out, find(r, "%s: %s does not terminate on %s (pin %s.%s)",
							r.ID, e.Knot, pin.Type, target.Ref, e.Pin))
					}
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R005", Severity: Error,
		Summary: "every reference resolves",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				check := func(ref, what string) {
					if ref != "" && s.ByID[ref] == nil {
						out = append(out, find(r, "%s: %s %q does not resolve", r.ID, what, ref))
					}
				}
				for _, n := range r.Nodes {
					check(n.Ref, "node ref")
				}
				for _, e := range r.Edges {
					check(e.Knot, "knot")
					check(e.Rigging, "rigging")
				}
				for _, c := range r.Cautions {
					for _, src := range c.Sources {
						check(src, "caution source")
					}
				}
				if r.Strength != nil {
					for _, c := range r.Strength.Claims {
						check(c.Source, "strength source")
					}
				}
				for _, d := range r.Diameters {
					check(d.Source, "diameter source")
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R007", Severity: Error,
		Summary: "kind matches the directory",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				want, known := spec.ExpectedKind(r.Path)
				if !known {
					out = append(out, find(r, "%s: not in a recognised data directory", r.ID))
				} else if r.Kind != want {
					out = append(out, find(r, "%s: kind %q in a %q directory", r.ID, r.Kind, want))
				}
				if got := spec.Category(r.ID); got == "" {
					out = append(out, find(r, "id %q has no category", r.ID))
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R010", Severity: Error,
		Summary: "a tier C claim carries a source",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				if r.Strength == nil {
					continue
				}
				for i, c := range r.Strength.Claims {
					if c.Tier == "C" && c.Source == "" {
						out = append(out, find(r, "%s: strength claim %d is tier C with no source", r.ID, i))
					}
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R011", Severity: Error,
		Summary: "a strength claim states its sample size",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				if r.Strength == nil {
					continue
				}
				for i, c := range r.Strength.Claims {
					if c.N == nil {
						out = append(out, find(r, "%s: strength claim %d omits n", r.ID, i))
					}
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R032", Severity: Error,
		Summary: "a selected variant exists on the component",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				if r.Kind != "rig" {
					continue
				}
				for _, n := range r.Nodes {
					comp := s.ByID[n.Ref]
					if comp == nil || comp.Variants == nil {
						continue
					}
					v, set := n.SelectedVariant(comp.Variants.Axis)
					if !set {
						out = append(out, find(r, "%s: node %q selects no %s from %s",
							r.ID, n.ID, comp.Variants.Axis, n.Ref))
						continue
					}
					if !hasValue(comp.Variants.Values, v) {
						out = append(out, find(r, "%s: node %q has %s %g, not offered by %s",
							r.ID, n.ID, comp.Variants.Axis, v, n.Ref))
					}
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R037", Severity: Error,
		Summary: "only a soft body can be rigged",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				nodes := map[string]*spec.Node{}
				for i := range r.Nodes {
					nodes[r.Nodes[i].ID] = &r.Nodes[i]
				}
				for _, e := range r.Edges {
					if e.Rel != "rigged" {
						continue
					}
					target := nodes[e.To]
					if target == nil {
						continue
					}
					if comp := s.ByID[target.Ref]; comp != nil && !comp.Soft {
						out = append(out, find(r, "%s: %q is rigged but is not soft", r.ID, target.Ref))
					}
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R039", Severity: Error,
		Summary: "a do-not-use caution carries a source",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				for i, c := range r.Cautions {
					if c.Severity == "do-not-use" && len(c.Sources) == 0 {
						out = append(out, find(r, "%s: caution %d is do-not-use with no source", r.ID, i))
					}
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R040", Severity: Warning,
		Summary: "citation debt is reported, not hidden",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				n := 0
				if r.Strength != nil {
					for _, c := range r.Strength.Claims {
						if c.Source == "src.needs-citation" {
							n++
						}
					}
				}
				for _, d := range r.Diameters {
					if d.Source == "src.needs-citation" {
						n++
					}
				}
				if n > 0 {
					out = append(out, find(r, "%s: %d outstanding citation(s)", r.ID, n))
				}
			}
			return out
		},
	})

	register(Rule{
		ID: "R023", Severity: Warning,
		Summary: "nothing publishes unvalidated",
		Check: func(s *spec.Set) []Finding {
			var out []Finding
			for _, r := range s.All {
				if r.Validation.Status == "unvalidated" {
					out = append(out, find(r, "%s: no human has confirmed this", r.ID))
				}
			}
			return out
		},
	})
}

// checkStop verifies a named travel bound exists and actually blocks passage.
func checkStop(s *spec.Set, r *spec.Record, nodes map[string]*spec.Node, threaded, side, stop string) []Finding {
	if stop == "" {
		return []Finding{find(r, "%s: threaded %q has no %s-side bound", r.ID, threaded, side)}
	}
	if stop == "open" {
		return []Finding{find(r, "%s: threaded %q is open toward the %s end and will come off", r.ID, threaded, side)}
	}
	n := nodes[stop]
	if n == nil {
		return []Finding{find(r, "%s: threaded %q names unknown %s-side stop %q", r.ID, threaded, side, stop)}
	}
	comp := s.ByID[n.Ref]
	if comp != nil && !comp.BlocksPassage {
		return []Finding{find(r, "%s: threaded %q is stopped by %q, which does not block passage", r.ID, threaded, stop)}
	}
	return nil
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}

func hasValue(xs []float64, want float64) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
