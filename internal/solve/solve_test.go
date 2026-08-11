// SPDX-License-Identifier: Apache-2.0

package solve_test

import (
	"strings"
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

// A correctly tied square knot, confirmed by someone who ties them. Note what
// it is not: alternating. Each eye tucks under the whole of the other cord's
// pair, and enforcing alternation instead drew a knot that was merely close.
func TestSquareKnotReading(t *testing.T) {
	want := map[string]string{"a": "U O U U O U", "b": "O U O O U O"}

	rec := load(t).ByID["knot.square"]
	g := solve.Geometry(rec)
	if g == nil {
		t.Fatal("square knot solves to no geometry")
	}
	for _, r := range solve.ReadCrossings(rec, len(g.Stages)-1) {
		if got := strings.Join(r.Sequence, " "); got != want[r.Cord] {
			t.Errorf("cord %s reads [%s], want [%s]", r.Cord, got, want[r.Cord])
		}
	}
}

// Every knot must still draw something rather than nothing, and the two cords
// must agree: one reads over wherever the other reads under.
func TestEveryStageDraws(t *testing.T) {
	for _, r := range load(t).All {
		if r.Kind != "knot" {
			continue
		}
		g := solve.Geometry(r)
		if g == nil {
			continue
		}
		for i := range g.Stages {
			rd := solve.ReadCrossings(r, i)
			if len(rd) != 2 || len(rd[0].Sequence) == 0 {
				t.Errorf("%s stage %d: no crossings drawn", r.ID, g.Stages[i].Stage)
				continue
			}
			if len(rd[0].Sequence) != len(rd[1].Sequence) {
				t.Errorf("%s stage %d: cords disagree on crossing count, %v vs %v",
					r.ID, g.Stages[i].Stage, rd[0].Sequence, rd[1].Sequence)
			}
		}
	}
}
