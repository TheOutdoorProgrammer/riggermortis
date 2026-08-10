// SPDX-License-Identifier: Apache-2.0

package rope

import "math"

// Settle is how hard a knot is pulled up.
type Settle struct {
	Iterations int
	// Diameter of the cord. Non-adjacent parts never come closer than this,
	// which is what stops a knot collapsing through itself.
	Diameter float64
	// Tension draws the standing parts apart each step, cinching the knot.
	Tension float64
	// Stiffness resists bending, so cord curves rather than kinking.
	Stiffness float64
}

// Relax settles rope into the shape tension and contact give it. Dressing is
// what makes a reef lie flat and a granny cock over, so an unsettled knot
// cannot show the difference. Constraints move points directly rather than
// integrating forces, so this cannot explode and needs no timestep tuning.
func Relax(cords []Cord, s Settle) {
	if s.Iterations == 0 {
		return
	}
	rest := make([][]float64, len(cords))
	for i, c := range cords {
		rest[i] = make([]float64, len(c.P)-1)
		for j := range rest[i] {
			rest[i][j] = dist(c.P[j], c.P[j+1])
		}
	}

	for it := 0; it < s.Iterations; it++ {
		pull(cords, s.Tension)
		for pass := 0; pass < 2; pass++ {
			keepLength(cords, rest)
		}
		unkink(cords, s.Stiffness)
		separate(cords, s.Diameter)
		keepLength(cords, rest)
	}
}

// pull draws each cord's ends apart along its own direction, which is how a
// knot is actually tightened.
func pull(cords []Cord, t float64) {
	if t == 0 {
		return
	}
	for _, c := range cords {
		n := len(c.P)
		if n < 4 {
			continue
		}
		d0 := norm(sub(c.P[0], c.P[3]))
		d1 := norm(sub(c.P[n-1], c.P[n-4]))
		c.P[0] = add(c.P[0], scale(d0, t))
		c.P[n-1] = add(c.P[n-1], scale(d1, t))
	}
}

// keepLength holds the cord inextensible. Ends are pinned so tension has
// something to act against.
func keepLength(cords []Cord, rest [][]float64) {
	for ci, c := range cords {
		for j := 0; j+1 < len(c.P); j++ {
			a, b := c.P[j], c.P[j+1]
			d := sub(b, a)
			l := length(d)
			if l == 0 {
				continue
			}
			corr := scale(d, (l-rest[ci][j])/l*0.5)
			pinA := j == 0
			pinB := j+2 == len(c.P)
			switch {
			case pinA && !pinB:
				c.P[j+1] = sub(b, scale(corr, 2))
			case pinB && !pinA:
				c.P[j] = add(a, scale(corr, 2))
			case !pinA && !pinB:
				c.P[j] = add(a, corr)
				c.P[j+1] = sub(b, corr)
			}
		}
	}
}

// unkink pulls each point toward the midpoint of its neighbours, which is a
// cheap stand-in for bending stiffness and keeps curves smooth under load.
func unkink(cords []Cord, k float64) {
	if k == 0 {
		return
	}
	for _, c := range cords {
		next := make([]V3, len(c.P))
		copy(next, c.P)
		for j := 1; j+1 < len(c.P); j++ {
			mid := scale(add(c.P[j-1], c.P[j+1]), 0.5)
			next[j] = lerp(c.P[j], mid, k)
		}
		copy(c.P[1:len(c.P)-1], next[1:len(next)-1])
	}
}

// separate pushes apart any two parts of the rope closer than its diameter,
// including a cord against itself. Without this a tightening knot passes
// straight through itself and the topology is lost.
func separate(cords []Cord, dia float64) {
	if dia == 0 {
		return
	}
	// Adjacent points are meant to touch; only parts this far apart along the
	// cord are treated as separate rope.
	const skip = 8

	type ref struct{ c, i int }
	var all []ref
	for ci, c := range cords {
		for i := range c.P {
			all = append(all, ref{ci, i})
		}
	}

	for x := 0; x < len(all); x++ {
		for y := x + 1; y < len(all); y++ {
			a, b := all[x], all[y]
			if a.c == b.c && b.i-a.i < skip {
				continue
			}
			pa, pb := cords[a.c].P[a.i], cords[b.c].P[b.i]
			d := sub(pb, pa)
			l := length(d)
			if l >= dia || l == 0 {
				continue
			}
			push := scale(d, (dia-l)/l*0.5)
			if a.i > 0 && a.i < len(cords[a.c].P)-1 {
				cords[a.c].P[a.i] = sub(pa, push)
			}
			if b.i > 0 && b.i < len(cords[b.c].P)-1 {
				cords[b.c].P[b.i] = add(pb, push)
			}
		}
	}
}

func sub(a, b V3) V3           { return V3{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func add(a, b V3) V3           { return V3{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func scale(a V3, f float64) V3 { return V3{a.X * f, a.Y * f, a.Z * f} }
func length(a V3) float64      { return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z) }
func dist(a, b V3) float64     { return length(sub(a, b)) }

func norm(a V3) V3 {
	l := length(a)
	if l == 0 {
		return V3{}
	}
	return scale(a, 1/l)
}

func lerp(a, b V3, f float64) V3 {
	return V3{a.X + (b.X-a.X)*f, a.Y + (b.Y-a.Y)*f, a.Z + (b.Z-a.Z)*f}
}
