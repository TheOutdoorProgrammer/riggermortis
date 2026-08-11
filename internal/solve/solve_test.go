// SPDX-License-Identifier: Apache-2.0

package solve_test

import (
	"testing"

	"github.com/theoutdoorprogrammer/riggermortis/internal/solve"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

func load(t *testing.T) *spec.Set {
	t.Helper()
	set, err := spec.Load("../../data")
	if err != nil {
		t.Fatal(err)
	}
	return set
}

// A reeve is two crossings. Counting it once gave the square knot four, and no
// four-crossing diagram can be a square knot, which is what forced the geometry
// to be corrected by hand.
func TestReeveIsTwoCrossings(t *testing.T) {
	k := solve.Read(load(t).ByID["knot.square"])
	if got := len(k.Crossings); got != 6 {
		t.Errorf("square knot has %d crossings, want 6", got)
	}
	if want := []int{1, -1}; len(k.Halves) != 2 || k.Halves[0] != want[0] || k.Halves[1] != want[1] {
		t.Errorf("half knots %v, want %v: opposite hands are the whole knot", k.Halves, want)
	}
	if !k.Bights() {
		t.Error("half knots cancel, so every end folds back and the bend is two bights")
	}
}

// Every knot in the dataset must draw an alternating diagram at every stage.
// Eyeballing a render cannot tell a reef from a granny; this can.
func TestStagesAlternate(t *testing.T) {
	for _, r := range load(t).All {
		if r.Kind != "knot" {
			continue
		}
		g := solve.Geometry(r)
		if g == nil {
			continue
		}
		for i := range g.Stages {
			if err := solve.CheckAlternating(g, i); err != nil {
				t.Errorf("%s stage %d: %v", r.ID, g.Stages[i].Stage, err)
			}
		}
	}
}

// The finished square knot must read over, under, over along both cords. Six
// crossings in the record mean nothing if the picture only draws four.
func TestSquareKnotDiagram(t *testing.T) {
	g := solve.Geometry(load(t).ByID["knot.square"])
	if g == nil {
		t.Fatal("square knot solves to no geometry")
	}
	last := len(g.Stages) - 1
	for _, r := range solve.ReadCrossings(g, last) {
		if len(r.Sequence) != 6 {
			t.Errorf("cord %s reads %v, want six crossings", r.Cord, r.Sequence)
		}
	}
}
