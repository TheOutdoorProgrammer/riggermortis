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
	radius := core + d*0.6
	pitch := d
	gap := d * 0.9
	if b.Tight {
		gap = d * 0.15
	}

	span := float64(turns) * pitch
	fold := gap/2 + span

	a := Cord{ID: "a", P: barrel(core, radius, pitch, fold, b.Lead, turns)}
	c := Cord{ID: "b", P: make([]V3, len(a.P))}
	for i, p := range a.P {
		c.P[i] = V3{-p.X, -p.Y, p.Z}
	}
	return []Cord{a, c}
}

// barrel is one line: a standing part running the length of the knot, a fold at
// the far end, then wraps coiling back over the pair to the tag.
func barrel(core, radius, pitch, fold, lead float64, turns int) []V3 {
	var out []V3

	turnR := (radius - core) / 2
	stop := fold + turnR*1.6

	const straight = 40
	for i := range straight {
		out = append(out, V3{lead*-1 + (stop+lead)*float64(i)/straight, core, 0})
	}

	// The fold, swung in the plane so the wraps start out at coil radius rather
	// than kinking straight up off the standing part.
	const bendSteps = 22
	cy := (core + radius) / 2
	for i := range bendSteps {
		t := math.Pi * float64(i) / bendSteps
		out = append(out, V3{stop + turnR*1.6*math.Sin(t), cy - turnR*math.Cos(t), 0})
	}

	// Wraps, running back along the core toward where the other coil will sit.
	steps := turns * 40
	for i := range steps + 1 {
		f := float64(i) / float64(steps)
		ang := math.Pi/2 + 2*math.Pi*float64(turns)*f
		x := stop - (stop-(fold-float64(turns)*pitch))*f
		out = append(out, V3{x, radius * math.Sin(ang), radius * math.Cos(ang)})
	}

	// The tag leaves away from the axis rather than along it. Running it on
	// would carry it straight into the other line's coil, which is where the
	// two knots jam against each other and where there is no room.
	const tag = 30
	last := out[len(out)-1]
	for i := 1; i <= tag; i++ {
		f := float64(i) / tag
		out = append(out, V3{
			last.X - pitch*0.5*f,
			last.Y - (radius*1.5+core)*f,
			last.Z * (1 - f),
		})
	}
	return out
}
