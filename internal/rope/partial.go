// SPDX-License-Identifier: Apache-2.0

package rope

import "math"

// Partial traces a cord only as far as the working end has reached. A
// part-tied knot is the finished path with the end not yet all the way
// through, which is what makes stepping show a knot actually being tied.
func Partial(c Cord, frac float64) Cord {
	if frac >= 1 || len(c.P) < 3 {
		return c
	}
	if frac <= 0 {
		frac = 0.02
	}

	total, cum := arclength(c.P)
	target := total * frac

	cut := 0
	for cut < len(cum)-1 && cum[cut+1] < target {
		cut++
	}
	if cut < 2 {
		cut = 2
	}

	head := append([]V3(nil), c.P[:cut+1]...)

	// The untied remainder leaves along the tangent, holding the cord's total
	// length so it reads as the same rope rather than a shorter one.
	tangent := norm(sub(head[len(head)-1], head[len(head)-4]))
	// Rope lies down; it does not fly off along its tangent. Leave on the
	// tangent, then settle toward horizontal over the loose length.
	rest := V3{sign(tangent.X), 0, 0}
	if tangent.X == 0 {
		rest = V3{1, 0, 0}
	}
	remain := (total - cum[cut]) * 0.42
	const steps = 30
	last := head[len(head)-1]
	pos := last
	for i := 1; i <= steps; i++ {
		f := float64(i) / steps
		ease := math.Sin(f * math.Pi / 2)
		d := norm(lerp(tangent, rest, math.Min(1, ease*1.4)))
		step := remain / steps
		pos = V3{pos.X + d.X*step, pos.Y + d.Y*step*0.8, pos.Z + d.Z*step*0.4}
		head = append(head, pos)
	}
	return Cord{ID: c.ID, P: head}
}

// Advance returns the stretch of path a working end covers between two stages,
// which is the arrow a tying diagram needs.
func Advance(c Cord, from, to float64) []V3 {
	if to <= from || len(c.P) < 3 {
		return nil
	}
	total, cum := arclength(c.P)
	a, b := total*from, total*to
	var out []V3
	for i, d := range cum {
		if d >= a && d <= b {
			out = append(out, c.P[i])
		}
	}
	if len(out) < 2 {
		return nil
	}
	return out
}

func arclength(p []V3) (float64, []float64) {
	cum := make([]float64, len(p))
	total := 0.0
	for i := 1; i < len(p); i++ {
		total += dist(p[i-1], p[i])
		cum[i] = total
	}
	return total, cum
}
