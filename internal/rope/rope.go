// SPDX-License-Identifier: Apache-2.0

// Package rope models cord in three dimensions and projects it to two.
//
// Earlier attempts authored the 2D projection directly, which meant guessing
// where strands go and which passes over. Both are consequences of the 3D
// shape, not choices: project the curve and whichever strand is nearer the
// camera at a crossing is the one on top. Over and under stop being authored.
package rope

import (
	"math"
	"sort"

	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

type V3 struct{ X, Y, Z float64 }

// Cord is a space curve sampled densely enough that straight segments between
// samples are visually indistinguishable from the curve.
type Cord struct {
	ID string
	P  []V3
}

// Twist wraps two cords around a shared axis. Turns counts half turns and
// carries handedness in its sign: reef is +1 then -1, granny +1 then +1.
type Twist struct {
	Turns int
}

// Layout describes a two-cord knot as a sequence of twists along one axis.
type Layout struct {
	Twists []Twist
	// Radius of the helix. Seating a knot shrinks it.
	Radius float64
	// Length of each half turn along the axis.
	Pitch float64
	// How far the free ends run past the last twist.
	Tail float64
	// Flare lifts the ends away from the axis so they read as loose.
	Flare float64
	// Stub is how far a working end runs past the knot. A bend joins two
	// ropes: each cord has a long standing part heading off one way and a
	// short working end the other, and they head off opposite ways. Equal
	// tails on both sides is what made this read as a twisted pair.
	Stub float64
}

const samplesPerHalfTurn = 34

// Build lays two cords out as counter-phased helices. Accumulating signed
// half turns means a reversal genuinely unwinds rather than continuing.
func Build(l Layout) []Cord {
	a := Cord{ID: "a"}
	b := Cord{ID: "b"}

	total := 0
	for _, t := range l.Twists {
		total += abs(t.Turns)
	}
	axisLen := float64(total) * l.Pitch
	x := -axisLen / 2

	// Lead-in tails, flared away from the axis so the standing parts read as
	// separate cords rather than as the start of the weave.
	a.P = append(a.P, tail(x, -l.Tail, l.Radius, 0, l.Radius+l.Flare, 0)...)
	b.P = append(b.P, tail(x, -l.Tail, -l.Radius, 0, -(l.Radius+l.Flare), 0)...)

	theta := 0.0
	for _, t := range l.Twists {
		dir := 1.0
		if t.Turns < 0 {
			dir = -1
		}
		steps := abs(t.Turns) * samplesPerHalfTurn
		for i := 1; i <= steps; i++ {
			f := float64(i) / float64(steps)
			th := theta + dir*math.Pi*float64(abs(t.Turns))*f
			xx := x + l.Pitch*float64(abs(t.Turns))*f
			a.P = append(a.P, V3{xx, l.Radius * math.Cos(th), l.Radius * math.Sin(th)})
			b.P = append(b.P, V3{xx, -l.Radius * math.Cos(th), -l.Radius * math.Sin(th)})
		}
		theta += dir * math.Pi * float64(abs(t.Turns))
		x += l.Pitch * float64(abs(t.Turns))
	}

	ay, az := l.Radius*math.Cos(theta), l.Radius*math.Sin(theta)
	outA, outB := l.Stub, l.Tail
	if outA == 0 {
		outA = l.Tail
	}
	a.P = append(a.P, tail(x, outA, ay, az, ay+sign(ay)*l.Flare*0.4, az)...)
	b.P = append(b.P, tail(x, outB, -ay, -az, -ay-sign(ay)*l.Flare, -az)...)

	return []Cord{a, b}
}

// sampleCR3 walks a Catmull-Rom spline through waypoints, densely enough that
// crossing detection reads the curve rather than the chords between points.
func sampleCR3(w []V3, n int) []V3 {
	at := func(i int) V3 {
		if i < 0 {
			return w[0]
		}
		if i >= len(w) {
			return w[len(w)-1]
		}
		return w[i]
	}
	out := make([]V3, 0, n)
	segs := len(w) - 1
	for s := 0; s < segs; s++ {
		p0, p1, p2, p3 := at(s-1), at(s), at(s+1), at(s+2)
		steps := n / segs
		for i := 0; i < steps; i++ {
			t := float64(i) / float64(steps)
			t2, t3 := t*t, t*t*t
			f := func(a, b, c, d float64) float64 {
				return 0.5 * ((2 * b) + (-a+c)*t +
					(2*a-5*b+4*c-d)*t2 + (-a+3*b-3*c+d)*t3)
			}
			out = append(out, V3{
				f(p0.X, p1.X, p2.X, p3.X),
				f(p0.Y, p1.Y, p2.Y, p3.Y),
				f(p0.Z, p1.Z, p2.Z, p3.Z),
			})
		}
	}
	return append(out, w[len(w)-1])
}

// tail runs a straight-ish lead from the axis out to a flared end.
func tail(x, dx, y0, z0, y1, z1 float64) []V3 {
	const n = 16
	out := make([]V3, 0, n)
	for i := 0; i <= n; i++ {
		f := float64(i) / n
		// Ease so the tail leaves the helix tangentially rather than kinking.
		e := f * f * (3 - 2*f)
		out = append(out, V3{x + dx*f, y0 + (y1-y0)*e, z0 + (z1-z0)*e})
	}
	if dx < 0 {
		for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
			out[i], out[j] = out[j], out[i]
		}
	}
	return out
}

type crossing struct {
	i    int     // segment index within the cord
	ti   float64 // parameter along cord i's segment
	over bool    // whether this cord is in front
}

// Project flattens the cords, resolving paint order by depth. Each cord is cut
// at every crossing so it can be in front at one and behind at the next.
// Box is a projection frame shared by every stage.
type Box struct{ MinX, MaxX, MinY, MaxY float64 }

// Frame returns bounds covering every stage, so one transform serves them all.
func Frame(stages [][]Cord) Box {
	b := Box{math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)}
	for _, cords := range stages {
		for _, c := range cords {
			for _, p := range c.P {
				b.MinX, b.MaxX = math.Min(b.MinX, p.X), math.Max(b.MaxX, p.X)
				b.MinY, b.MaxY = math.Min(b.MinY, p.Y), math.Max(b.MaxY, p.Y)
			}
		}
	}
	return b
}

func Project(cords []Cord, w, h float64) []spec.Segment {
	return ProjectIn(cords, Frame([][]Cord{cords}), w, h)
}

// ProjectIn projects within a fixed frame.
func ProjectIn(cords []Cord, box Box, w, h float64) []spec.Segment {
	pts := make([][]spec.Point, len(cords))
	depth := make([][]float64, len(cords))

	minX, maxX, minY, maxY := box.MinX, box.MaxX, box.MinY, box.MaxY
	pad := 34.0
	sx := (w - 2*pad) / (maxX - minX)
	sy := (h - 2*pad) / (maxY - minY)
	s := math.Min(sx, sy)
	ox := (w - (maxX-minX)*s) / 2
	oy := (h - (maxY-minY)*s) / 2

	for ci, c := range cords {
		for _, p := range c.P {
			pts[ci] = append(pts[ci], spec.Point{ox + (p.X-minX)*s, oy + (p.Y-minY)*s})
			depth[ci] = append(depth[ci], p.Z)
		}
	}

	// Cut points per cord, as indices where a crossing happens.
	cuts := make([][]crossing, len(cords))
	for ci := 0; ci < len(cords); ci++ {
		for cj := ci + 1; cj < len(cords); cj++ {
			for i := 0; i+1 < len(pts[ci]); i++ {
				for j := 0; j+1 < len(pts[cj]); j++ {
					ti, tj, ok := intersect(pts[ci][i], pts[ci][i+1], pts[cj][j], pts[cj][j+1])
					if !ok {
						continue
					}
					zi := depth[ci][i] + (depth[ci][i+1]-depth[ci][i])*ti
					zj := depth[cj][j] + (depth[cj][j+1]-depth[cj][j])*tj
					cuts[ci] = append(cuts[ci], crossing{i: i, ti: ti, over: zi > zj})
					cuts[cj] = append(cuts[cj], crossing{i: j, ti: tj, over: zj > zi})
				}
			}
		}
	}

	var out []spec.Segment
	const halo = 9 // samples either side of a crossing that share its paint order

	for ci, c := range cords {
		sort.Slice(cuts[ci], func(a, b int) bool { return cuts[ci][a].i < cuts[ci][b].i })
		cuts[ci] = cluster(cuts[ci])

		// Paint order per sample: 0 free, 1 behind, 2 in front.
		order := make([]int, len(pts[ci]))
		for _, x := range cuts[ci] {
			z := 1
			if x.over {
				z = 2
			}
			for k := x.i - halo; k <= x.i+halo+1; k++ {
				if k >= 0 && k < len(order) && z > order[k] {
					order[k] = z
				}
			}
		}

		start := 0
		for k := 1; k <= len(order); k++ {
			if k == len(order) || order[k] != order[start] {
				// Overlap generously: adjacent pieces must share enough length
				// that no seam is visible where paint order changes.
				const lap = 6
				lo, hi := start-lap, k+lap
				if lo < 0 {
					lo = 0
				}
				if hi > len(pts[ci]) {
					hi = len(pts[ci])
				}
				for hi-lo < 3 {
					if lo > 0 {
						lo--
					} else if hi < len(pts[ci]) {
						hi++
					} else {
						break
					}
				}
				if hi-lo >= 2 {
					out = append(out, spec.Segment{
						Cord: c.ID, Z: order[start], Points: pts[ci][lo:hi],
					})
				}
				start = k
			}
		}
	}
	return out
}

// cluster collapses the run of intersections a polyline registers around a
// single true crossing. Two counter-phased helices are coincident in
// projection where they cross, so nearby segment pairs all report a hit and
// the cord would otherwise be cut into dashes.
func cluster(xs []crossing) []crossing {
	if len(xs) == 0 {
		return xs
	}
	const gap = 6
	out := []crossing{xs[0]}
	for _, x := range xs[1:] {
		last := &out[len(out)-1]
		if x.i-last.i <= gap {
			// Whichever reading is in front wins: a true crossing has one
			// answer, and the coincident samples around it are noise.
			last.over = last.over || x.over
			continue
		}
		out = append(out, x)
	}
	return out
}

// intersect returns the parameters where two segments cross, if they do.
func intersect(p1, p2, p3, p4 spec.Point) (float64, float64, bool) {
	d := (p2[0]-p1[0])*(p4[1]-p3[1]) - (p2[1]-p1[1])*(p4[0]-p3[0])
	if math.Abs(d) < 1e-9 {
		return 0, 0, false
	}
	t := ((p3[0]-p1[0])*(p4[1]-p3[1]) - (p3[1]-p1[1])*(p4[0]-p3[0])) / d
	u := ((p3[0]-p1[0])*(p2[1]-p1[1]) - (p3[1]-p1[1])*(p2[0]-p1[0])) / d
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return 0, 0, false
	}
	return t, u, true
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(f float64) float64 {
	if f < 0 {
		return -1
	}
	return 1
}
