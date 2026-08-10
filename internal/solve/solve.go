// SPDX-License-Identifier: Apache-2.0

// Package solve derives knot geometry from the tying actions.
//
// The notation already contains the diagram. A knot diagram is a sequence of
// crossings between adjacent strands, each with a handedness, which is a braid
// word; and the stages already say where crossings happen and which way they
// go. `ML` with chirality "/" is a positive crossing and "\" is negative, so
// the square knot's stages read as σ σ σ⁻¹ σ⁻¹ while a granny reads σ σ σ σ.
// One field differs and the picture differs with it.
//
// Nothing here is authored. Coordinates come out of the crossing sequence, so
// adding a knot record adds a diagram without anyone drawing one.
package solve

import (
	"github.com/theoutdoorprogrammer/riggermortis/internal/rope"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

// Canvas and track geometry. Two strands run left to right and swap places at
// each crossing, which is the clearest way to show an interweave.
const (
	width      = 620.0
	height     = 300.0
	margin     = 78.0
	trackTop   = 68.0
	trackLow   = 132.0
	maxSwapPad = 46.0
	tailFade   = 6.0
)

// Crossing is one exchange between the two strands.
type Crossing struct {
	// Sign +1 means the strand travelling downward passes over.
	Sign int
	// Stage that introduced it, so a diagram can be built up step by step.
	Stage int
}

// Knot is the combinatorial reading of a knot's stages.
type Knot struct {
	Crossings []Crossing
	// Tighten is the stage at which the knot is drawn seated rather than open.
	Tighten int
}

// Read maps tying actions onto crossings. Reeving through a named loop
// inherits that loop's handedness; gripping and moving change nothing.
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
				k.Crossings = append(k.Crossings, Crossing{Sign: s, Stage: stage})
			case "RV":
				if s, ok := loopSign[a.Through]; ok {
					k.Crossings = append(k.Crossings, Crossing{Sign: s, Stage: stage})
				}
			case "MT":
				s := 1
				if a.Rotation == "CCW" {
					s = -1
				}
				for n := 0; n < repeatCount(a.Repeat); n++ {
					k.Crossings = append(k.Crossings, Crossing{Sign: s, Stage: stage})
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

// Geometry lays the crossing sequence out, one drawing per stage. Returns nil
// when there are no crossings: no diagram beats a misleading one.
// Geometry draws the knot at each stage, growing the weave as it is tied.
//
// Stage k shows every crossing introduced up to k, so stepping forward adds a
// turn to the same rope rather than swapping in a new picture. Returns nil
// when the actions yield no crossings: no diagram beats a misleading one.
func Geometry(r *spec.Record) *spec.Geometry {
	k := Read(r)
	if len(k.Crossings) == 0 {
		return nil
	}

	g := &spec.Geometry{Width: width, Height: height, Cords: []string{"a", "b"}}

	for i, st := range r.Stages {
		stage := st.ID
		if stage == 0 {
			stage = i + 1
		}
		var upto []Crossing
		for _, c := range k.Crossings {
			if c.Stage <= stage {
				upto = append(upto, c)
			}
		}
		if len(upto) == 0 {
			upto = k.Crossings[:1]
		}

		l := rope.Layout{Twists: twists(upto), Radius: 26, Pitch: 74,
			Tail: 150, Flare: 40, Stub: 60}
		cords := rope.Build(l)

		// Every stage settles a little so the cord hangs rather than tracing a
		// perfect helix. The stage that pulls settles hard, and that is where a
		// reef lying flat and a granny cocking over become visibly different.
		settle := rope.Settle{Iterations: 40, Diameter: 30, Tension: 0.12, Stiffness: 0.16}
		if k.Tighten > 0 && stage >= k.Tighten {
			settle = rope.Settle{Iterations: 160, Diameter: 30, Tension: 0.55, Stiffness: 0.18}
		}
		rope.Relax(cords, settle)

		g.Stages = append(g.Stages, spec.StageGeometry{
			Stage:    stage,
			Segments: rope.Project(cords, width, height),
		})
	}
	return g
}

// twists collapses a run of same-handed crossings into one twist, which is
// what a half knot physically is: the cords wrapping around each other.
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
