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
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

// Canvas and track geometry. Two strands run left to right and swap places at
// each crossing, which is the clearest way to show an interweave.
const (
	width      = 560.0
	height     = 200.0
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
		seated := k.Tighten > 0 && stage >= k.Tighten
		g.Stages = append(g.Stages, spec.StageGeometry{
			Stage:    stage,
			Segments: weave(k.Crossings, seated),
		})
	}
	return g
}

// weave builds both cords as continuous polylines, then cuts each into pieces
// at the crossings so paint order can differ per crossing. Pieces overlap
// slightly at their joins, so a casing never bites into its own cord.
func weave(cs []Crossing, seated bool) []spec.Segment {
	n := len(cs)

	left, right := margin, width-margin
	top, low := trackTop, trackLow
	if seated {
		// Seating pulls the crossings together and flattens the strands.
		mid := (left + right) / 2
		left, right = mid-(mid-left)*0.55, mid+(right-mid)*0.55
		c := (top + low) / 2
		top, low = c-(c-top)*0.62, c+(low-c)*0.62
	}

	span := (right - left) / float64(n)
	pad := span * 0.42
	if pad > maxSwapPad {
		pad = maxSwapPad
	}
	mid := (top + low) / 2

	xs := make([]float64, n)
	for i := range xs {
		xs[i] = left + span*(float64(i)+0.5)
	}

	// Track occupancy: cord "a" starts on top, "b" below, and they exchange at
	// every crossing.
	yFor := func(cord string, after int) float64 {
		swapped := after%2 == 1
		aTop := (cord == "a") != swapped
		if aTop {
			return top
		}
		return low
	}

	var out []spec.Segment

	for _, cord := range []string{"a", "b"} {
		// Lead-in.
		out = append(out, spec.Segment{
			Cord: cord, Z: 0,
			Points: []spec.Point{
				{0, yFor(cord, 0)},
				{(xs[0] - pad) * 0.5, yFor(cord, 0)},
				{xs[0] - pad + tailFade, yFor(cord, 0)},
			},
		})

		for i := 0; i < n; i++ {
			yIn, yOut := yFor(cord, i), yFor(cord, i+1)
			descending := yOut > yIn
			over := (cs[i].Sign > 0) == descending

			z := 1
			if over {
				z = 2
			}
			out = append(out, spec.Segment{
				Cord: cord, Z: z, Stage: cs[i].Stage,
				Points: []spec.Point{
					{xs[i] - pad - tailFade, yIn},
					{xs[i] - pad*0.45, yIn + (mid-yIn)*0.30},
					{xs[i], mid},
					{xs[i] + pad*0.45, yOut - (yOut-mid)*0.30},
					{xs[i] + pad + tailFade, yOut},
				},
			})

			// Run to the next crossing, or out to the edge.
			nextX := width
			if i+1 < n {
				nextX = xs[i+1] - pad
			}
			out = append(out, spec.Segment{
				Cord: cord, Z: 0,
				Points: []spec.Point{
					{xs[i] + pad - tailFade, yOut},
					{(xs[i] + pad + nextX) / 2, yOut},
					{nextX + tailFade, yOut},
				},
			})
		}
	}
	return out
}
