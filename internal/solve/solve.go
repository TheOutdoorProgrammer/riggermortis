// SPDX-License-Identifier: Apache-2.0

// Package solve derives knot geometry from the tying actions.
//
// A reeve is two crossings, not one. Passing an end through a loop enters the
// region the loop bounds and then leaves it, and each of those is a crossing in
// the diagram. Counting it once is what made the square knot read as four
// crossings when a dressed one has six, and the missing pair is why earlier
// attempts had to correct the picture by hand.
//
// Nothing here is authored. Coordinates come out of the crossing sequence, so
// adding a knot record adds a diagram without anyone drawing one.
package solve

import (
	"github.com/theoutdoorprogrammer/riggermortis/internal/rope"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

const (
	width  = 620.0
	height = 300.0
	// diameter is the cord's own width. Every knot length is a multiple of it.
	diameter = 30.0
	lead     = 300.0
)

// Kind is what a crossing is doing, which decides nothing about the picture but
// everything about how many crossings an action is worth.
type Kind int

const (
	// Lay is one part put across another.
	Lay Kind = iota
	// ReeveIn enters the region a loop bounds, ReeveOut leaves it.
	ReeveIn
	ReeveOut
)

// Crossing is one passage of a working end past another strand.
type Crossing struct {
	// Sign +1 means the strand travelling downward passes over.
	Sign int
	// Stage that introduced it, so a diagram can be built up step by step.
	Stage int
	Kind  Kind
}

// Knot is the combinatorial reading of a knot's stages.
type Knot struct {
	Crossings []Crossing
	// Halves is the handedness of each half knot in tying order. Reef is +1
	// then -1, granny +1 then +1, and that difference is the whole knot.
	Halves []int
	// Tighten is the stage at which the knot is drawn seated rather than open.
	Tighten int
}

// Bights reports whether the half knots cancel, bringing every working end back
// beside its own standing part. Half knots of one hand stand the knot on edge
// instead, and that shape has no model here yet, so it falls through unfolded.
func (k Knot) Bights() bool {
	return len(k.Halves) == 2 && k.Halves[0]+k.Halves[1] == 0
}

// Read maps tying actions onto crossings. Reeving through a named loop inherits
// that loop's handedness; reeving through a hook eye or other named object is
// not a crossing with the rope at all, so it contributes none.
func Read(r *spec.Record) Knot {
	var k Knot
	k.Tighten = -1
	loopSign := map[string]int{}

	for i, st := range r.Stages {
		stage := st.ID
		if stage == 0 {
			stage = i + 1
		}
		for _, a := range st.Actions {
			switch a.Verb {
			case "ML":
				s := signOf(a.Chirality)
				if s == 0 {
					continue
				}
				if a.Names != "" {
					loopSign[a.Names] = s
				}
				k.Crossings = append(k.Crossings, Crossing{Sign: s, Stage: stage, Kind: Lay})
				k.Halves = append(k.Halves, s)
			case "RV":
				s, ok := loopSign[a.Through]
				if !ok {
					continue
				}
				k.Crossings = append(k.Crossings,
					Crossing{Sign: s, Stage: stage, Kind: ReeveIn},
					Crossing{Sign: s, Stage: stage, Kind: ReeveOut})
			case "MT":
				s := 1
				if a.Rotation == "CCW" {
					s = -1
				}
				for n := 0; n < repeatCount(a.Repeat); n++ {
					k.Crossings = append(k.Crossings, Crossing{Sign: s, Stage: stage, Kind: Lay})
				}
			case "MV":
				if a.Force == "pull" && len(k.Crossings) > 0 {
					k.Tighten = stage
				}
			}
		}
	}
	return k
}

func signOf(chirality string) int {
	switch chirality {
	case "/":
		return 1
	case `\`:
		return -1
	}
	return 0
}

// repeatCount collapses a range to its lower bound: a diagram showing the
// fewest turns that still reads correctly is clearer than one showing seven.
func repeatCount(v any) int {
	switch t := v.(type) {
	case nil:
		return 1
	case int:
		return t
	case float64:
		return int(t)
	case []any:
		if len(t) > 0 {
			if n, ok := t[0].(int); ok {
				return n
			}
			if f, ok := t[0].(float64); ok {
				return int(f)
			}
		}
	}
	return 1
}

// Cords returns the solved 3D cords for each stage, before projection.
func Cords(r *spec.Record) [][]rope.Cord {
	k := Read(r)
	if len(k.Crossings) == 0 {
		return nil
	}
	var out [][]rope.Cord
	for i, st := range r.Stages {
		stage := st.ID
		if stage == 0 {
			stage = i + 1
		}
		out = append(out, buildStage(k, stage))
	}
	return out
}

// Geometry draws the knot at each stage. Stage k shows every crossing made up
// to k, so stepping forward adds a turn to the same rope rather than swapping
// in a new picture. Nil when there are no crossings: no diagram beats a lie.
func Geometry(r *spec.Record) *spec.Geometry {
	all := Cords(r)
	if all == nil {
		return nil
	}

	g := &spec.Geometry{Width: width, Height: height, Cords: []string{"a", "b"}}
	frame := rope.Frame(all)
	for i, st := range r.Stages {
		stage := st.ID
		if stage == 0 {
			stage = i + 1
		}
		g.Stages = append(g.Stages, spec.StageGeometry{
			Stage:    stage,
			Segments: rope.ProjectIn(all[i], frame, width, height),
		})
	}
	return g
}

// buildStage produces the 3D cords for one stage, shared by the SVG path and
// any other renderer that wants the same curves.
func buildStage(k Knot, stage int) []rope.Cord {
	var done []Crossing
	for _, c := range k.Crossings {
		if c.Stage <= stage {
			done = append(done, c)
		}
	}
	if len(done) == 0 {
		done = k.Crossings[:1]
	}
	tight := k.Tighten > 0 && stage >= k.Tighten

	// A cord only folds once the last tuck is made. Until then the ends have
	// not come back to their own standing parts and the pair is still just
	// twisted, so a part-tied knot is a shorter weave and not a cropped one.
	var cords []rope.Cord
	if k.Bights() && len(done) == len(k.Crossings) {
		// Bights come out dressed, so there is nothing for physics to do. Tension
		// does not seat this knot, it flattens the eyes the clasp is made of, and
		// contact pushes the paired legs apart into a shape rope does not take.
		cords = rope.Weave{
			Signs:     k.Halves,
			Crossings: len(k.Crossings),
			Diameter:  diameter,
			Lead:      lead,
			Tight:     tight,
		}.Build()
	} else {
		cords = rope.Build(rope.Layout{Twists: twists(done), Radius: 26,
			Pitch: 74, Tail: 150, Flare: 40, Stub: 60})
		settle := rope.Settle{Iterations: 40, Diameter: diameter, Tension: 0.12, Stiffness: 0.16}
		if tight {
			settle = rope.Settle{Iterations: 160, Diameter: diameter, Tension: 0.55, Stiffness: 0.18}
		}
		rope.Relax(cords, settle)
		rope.Alternate(cords, diameter*0.55, halves(done))
	}
	return cords
}

// halves is the handedness of each half knot among the crossings given, which
// is what a partly tied knot has to alternate from.
func halves(cs []Crossing) []int {
	var out []int
	for _, c := range cs {
		if c.Kind == Lay {
			out = append(out, c.Sign)
		}
	}
	return out
}

// twists collapses a run of same-handed crossings into one twist, which is what
// a half knot physically is: the cords wrapping around each other.
func twists(cs []Crossing) []rope.Twist {
	var out []rope.Twist
	for _, c := range cs {
		if n := len(out); n > 0 && sameSign(out[n-1].Turns, c.Sign) {
			out[n-1].Turns += c.Sign
			continue
		}
		out = append(out, rope.Twist{Turns: c.Sign})
	}
	return out
}

func sameSign(a, b int) bool { return (a < 0) == (b < 0) }
