// SPDX-License-Identifier: Apache-2.0

package rope

import "math"

// Weave builds a bend with no per-knot geometry in it. Each half knot swaps
// which side the working ends are on, so half knots that cancel bring every end
// back beside its own standing part: the cord has folded, and that is a bight.
type Weave struct {
	// Signs is the handedness of each half knot, in tying order.
	Signs []int
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

	// Legs lie against each other at the mouth; the eye has to swallow the other
	// cord's pair, so it is that pair's width plus a cord. Reach lands the eye
	// just past the other one, where that pair has closed up again.
	mouth := d * 0.58
	eyeR := mouth + d*1.6
	reach := eyeR * 2.3
	lead := w.Lead
	if lead == 0 {
		lead = d * 10
	}
	if !w.Tight {
		mouth *= 1.25
		eyeR *= 1.1
		reach *= 1.1
	}

	a := Cord{ID: "a", P: bight(mouth, eyeR, reach, lead)}

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

// wrap is how far the eye closes around itself. Past half a turn the eye is a
// ring the other cord's pair cannot leave sideways, which is the grip.
const wrap = 2.36

// bight folds a cord into a lasso: two legs opening into a round eye. That eye
// collars the other cord's paired legs, well outside the other eye and never
// through it. Threading it instead ties a different knot that looks close.
func bight(mouth, eyeR, reach, lead float64) []V3 {
	cx := reach - eyeR
	tip := V3{cx + eyeR*math.Cos(wrap), eyeR * math.Sin(wrap), 0}
	// Travelling the eye with the angle decreasing, so the tangent arriving at
	// the leg join points back along the leg.
	tan := V3{math.Sin(wrap), -math.Cos(wrap), 0}

	// The legs stay paired until clear of the other cord's eye, which sits
	// opposite this one. Spreading any earlier lays them across that ring wide
	// instead of running through it together, and the pairs then never cross.
	hold := V3{eyeR*0.7 - cx, mouth, 0}
	in := append(run(V3{-lead, mouth, 0}, hold), leg(hold, tip, tan)...)

	out := append([]V3(nil), in...)
	const steps = 72
	for i := 1; i < steps; i++ {
		a := wrap - 2*wrap*float64(i)/steps
		out = append(out, V3{cx + eyeR*math.Cos(a), eyeR * math.Sin(a), 0})
	}
	// Only the leg comes back: the eye already carries itself from one join to
	// the other, and mirroring that too would retrace it and pinch it shut.
	for i := len(in) - 1; i >= 0; i-- {
		out = append(out, V3{in[i].X, -in[i].Y, in[i].Z})
	}
	return out
}

// run is the straight stretch where a cord's legs lie against each other.
func run(from, to V3) []V3 {
	const steps = 24
	out := make([]V3, 0, steps)
	for i := range steps {
		out = append(out, lerp(from, to, float64(i)/steps))
	}
	return out
}

// leg spreads from the paired run into the eye, leaving flat and arriving on
// the eye's tangent so the cord does not kink. The arrival tangent stays short:
// at full chord the curve dips back under the run.
func leg(from, to, tan V3) []V3 {
	k := dist(from, to)
	t0 := V3{k * 0.75, 0, 0}
	t1 := scale(tan, k*0.35)
	const steps = 30
	out := make([]V3, 0, steps)
	for i := range steps {
		t := float64(i) / steps
		t2, t3 := t*t, t*t*t
		h00, h10 := 2*t3-3*t2+1, t3-2*t2+t
		h01, h11 := -2*t3+3*t2, t3-t2
		out = append(out, V3{
			h00*from.X + h10*t0.X + h01*to.X + h11*t1.X,
			h00*from.Y + h10*t0.Y + h01*to.Y + h11*t1.Y,
			0,
		})
	}
	return out
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
