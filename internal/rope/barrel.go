// SPDX-License-Identifier: Apache-2.0

package rope

import "math"

// Barrels is a knot whose turns are taken about a core rather than about the
// other cord. Each end folds back over both standing parts and coils along
// them, so what grips is a collar of wraps and not an interlock.
type Barrels struct {
	// Turns is how many wraps each end takes about the pair.
	Turns int
	// Diameter of the cord. Every other length is a multiple of it.
	Diameter float64
	// Lead is how far the standing parts run outside the knot.
	Lead float64
	// Tight slides the two coils together, which is what seating them does.
	Tight bool
}

// Build returns the two cords. The second is the first turned half a
// revolution: the lines face opposite ways and each coils over the other.
func (b Barrels) Build() []Cord {
	d := b.Diameter
	if d == 0 {
		d = 30
	}
	turns := max(b.Turns, 1)

	// The core is the two standing parts lying against each other. A wrap has
	// to clear both of them plus its own thickness, which fixes its radius.
	// Wraps lie against each other, so the pitch is one cord thickness. Any
	// less and neighbouring turns merge into a solid mass when drawn.
	core := d * 0.5
	radius := core + d*0.8
	pitch := d
	// Seated still leaves a coil's width between them: the two knots jam on
	// each other's buried tags, not by grinding their wraps together.
	gap := d * 2.4
	if b.Tight {
		gap = d * 1.2
	}

	span := float64(turns) * pitch
	fold := gap/2 + span

	a := Cord{ID: "a", P: barrel(core, radius, pitch, gap, fold, b.Lead, turns)}
	c := Cord{ID: "b", P: make([]V3, len(a.P))}
	for i, p := range a.P {
		c.P[i] = V3{-p.X, -p.Y, p.Z}
	}
	return []Cord{a, c}
}

// barrel is one line: a standing part running the length of the knot, a fold at
// the far end, then wraps coiling back over the pair to the tag.
func barrel(core, radius, pitch, gap, fold, lead float64, turns int) []V3 {
	var out []V3

	turnR := (radius - core) / 2
	stop := fold + turnR*1.6

	straight := paces(stop + lead)
	for i := range straight {
		out = append(out, V3{lead*-1 + (stop+lead)*float64(i)/float64(straight), core, 0})
	}

	// The fold, swung in the plane so the wraps start out at coil radius rather
	// than kinking straight up off the standing part.
	bendSteps := paces(math.Pi * turnR)
	cy := (core + radius) / 2
	for i := range bendSteps {
		t := math.Pi * float64(i) / float64(bendSteps)
		out = append(out, V3{stop + turnR*1.6*math.Sin(t), cy - turnR*math.Cos(t), 0})
	}

	// Wraps, running back along the core toward where the other coil will sit.
	steps := paces(float64(turns) * 2 * math.Pi * radius)
	for i := range steps + 1 {
		f := float64(i) / float64(steps)
		ang := math.Pi/2 + 2*math.Pi*float64(turns)*f
		x := stop - (stop-(fold-float64(turns)*pitch))*f
		out = append(out, V3{x, radius * math.Sin(ang), radius * math.Cos(ang)})
	}

	// The end drops to the core within half a wrap, then runs level into the
	// other line's coil and finishes under its wraps. That is the lock: each
	// coil sits on the other's end, so neither can back off its own turns.
	// It lies against the inside of the coil, on the far side from its own
	// standing part. On the axis the two ends would be in the same place, and
	// beside a standing part each one hides behind it.
	rest := V3{0, -(radius - core*0.5), 0}
	last := out[len(out)-1]

	drop := paces(pitch)
	for i := 1; i <= drop; i++ {
		e := float64(i) / float64(drop)
		e = e * e * (3 - 2*e)
		out = append(out, V3{
			last.X - pitch*0.5*e,
			last.Y + (rest.Y-last.Y)*e,
			last.Z + (rest.Z-last.Z)*e,
		})
	}

	from := out[len(out)-1]
	reach := gap + pitch*float64(buried)
	level := paces(reach)
	for i := 1; i <= level; i++ {
		out = append(out, V3{from.X - reach*float64(i)/float64(level), rest.Y, rest.Z})
	}
	return out
}

// paces samples a stretch at the same density as every other. Crossings resolve
// over a fixed window of samples, so a coarser stretch merges neighbours into
// one answer and a strand threading a coil paints straight over it instead.
func paces(length float64) int {
	return max(int(length/1.4), 8)
}

// buried is how many of the other coil's wraps end up sitting on this tag.
// Fewer and the tag is not held; the knot relies on it being pinned.
const buried = 3
