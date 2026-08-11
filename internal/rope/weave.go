// SPDX-License-Identifier: Apache-2.0

package rope

import "math"

// Weave builds a bend with no per-knot geometry in it. Each half knot swaps
// which side the working ends are on, so half knots that cancel bring every end
// back beside its own standing part: the cord has folded, and that is a bight.
type Weave struct {
	// Signs is the handedness of each half knot, in tying order.
	Signs []int
	// Crossings is how many times one cord passes the other.
	Crossings int
	// Diameter of the cord. Every other length is a multiple of it: a knot is
	// as big as the rope it is tied in and has no scale of its own.
	Diameter float64
	// Lead is how far the standing parts and ends run outside the knot.
	Lead float64
	// Tight closes the bights up, which is what pulling the standing parts does.
	Tight bool
}

// Build returns the two cords, depths already resolved at every crossing.
func (w Weave) Build() []Cord {
	d := w.Diameter
	if d == 0 {
		d = 30
	}

	// Legs lie against each other at the mouth, spread at the belly to collar
	// the other cord's pair, and the turn reaches past it so they interlock.
	mouth := d * 0.58
	belly := d * 2.1
	reach := d * 0.8 * float64(w.Crossings)
	lead := w.Lead
	if lead == 0 {
		lead = d * 10
	}
	if w.Tight {
		mouth *= 0.86
		belly *= 0.82
		reach *= 0.9
	}

	a := Cord{ID: "a", P: sampleCR3(bight(mouth, belly, reach, lead), 26*8)}

	// The second bight is the first turned half a revolution about the knot's
	// centre. Mirroring it instead would land the paired crossings at the same
	// x, where a left-to-right reading cannot put them in order.
	b := Cord{ID: "b", P: make([]V3, len(a.P))}
	for i, p := range a.P {
		b.P[i] = V3{-p.X, -p.Y, p.Z}
	}

	cords := []Cord{a, b}
	Alternate(cords, d*0.55, w.Signs)
	return cords
}

// bight folds one cord back on itself: a long teardrop, mouth left, turn right,
// every control point a multiple of diameter. clasp drags the eye back past the
// origin, or the second bight (this one turned about it) would only meet it.
func bight(mouth, belly, reach, lead float64) []V3 {
	const clasp = 0.34
	x := func(f float64) float64 { return reach*f - reach*clasp }
	return []V3{
		{-lead, mouth, 0},
		{x(-1.4), mouth, 0},
		{x(0), mouth * 1.4, 0},
		{x(0.3), belly * 0.8, 0},
		{x(0.62), belly, 0},
		{x(0.9), belly * 0.52, 0},
		{x(1), 0, 0},
		{x(0.9), -belly * 0.52, 0},
		{x(0.62), -belly, 0},
		{x(0.3), -belly * 0.8, 0},
		{x(0), -mouth * 1.4, 0},
		{x(-1.4), -mouth, 0},
		{-lead, -mouth, 0},
	}
}

// Alternate sets which cord is in front at each crossing, overriding the depth
// the layout arrived with: a helix retraces its depths when it winds back and
// leaves two crossings the same way round. Handedness lays out, never sinks.
func Alternate(cords []Cord, lift float64, signs []int) {
	if len(cords) < 2 {
		return
	}
	a, b := cords[0], cords[1]
	over := len(signs) == 0 || signs[0] > 0
	for _, h := range meetings(a.P, b.P) {
		z := lift
		if !over {
			z = -lift
		}
		bump(a.P, h.i, z)
		bump(b.P, h.j, -z)
		over = !over
	}
}

type meeting struct{ i, j int }

// meetings returns crossings in the xy plane, ordered along the first path. A
// polyline registers one true crossing several times over, so runs collapse.
func meetings(p, q []V3) []meeting {
	var out []meeting
	last := -100
	for i := 0; i+1 < len(p); i++ {
		for j := 0; j+1 < len(q); j++ {
			if !segmentsCross(p[i], p[i+1], q[j], q[j+1]) {
				continue
			}
			if i-last <= 8 {
				continue
			}
			last = i
			out = append(out, meeting{i, j})
			break
		}
	}
	return out
}

func segmentsCross(a, b, c, d V3) bool {
	den := (b.X-a.X)*(d.Y-c.Y) - (b.Y-a.Y)*(d.X-c.X)
	if math.Abs(den) < 1e-9 {
		return false
	}
	t := ((c.X-a.X)*(d.Y-c.Y) - (c.Y-a.Y)*(d.X-c.X)) / den
	u := ((c.X-a.X)*(b.Y-a.Y) - (c.Y-a.Y)*(b.X-a.X)) / den
	return t >= 0 && t <= 1 && u >= 0 && u <= 1
}

// bump eases a stretch of cord toward a depth, so it lifts over its neighbour
// rather than stepping over it. Easing toward rather than adding leaves the
// cord's own shape intact everywhere except where a crossing needs an answer.
func bump(p []V3, at int, z float64) {
	const span = 22
	for k := at - span; k <= at+span; k++ {
		if k < 0 || k >= len(p) {
			continue
		}
		f := (1 + math.Cos(float64(k-at)/span*math.Pi)) / 2
		p[k].Z += (z - p[k].Z) * f
	}
}
