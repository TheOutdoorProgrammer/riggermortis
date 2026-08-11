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

// casingHalo is half the width the renderer paints a cord's casing at. A
// crossing has to hold its paint order over at least that much rope.
const casingHalo = 15.0

// pace is a cord's mean sample spacing, so a halo in real units converts to
// samples for cords built at different densities.
func pace(pts []spec.Point) float64 {
	if len(pts) < 2 {
		return 1
	}
	total := 0.0
	for i := 1; i < len(pts); i++ {
		total += math.Hypot(pts[i][0]-pts[i-1][0], pts[i][1]-pts[i-1][1])
	}
	return math.Max(total/float64(len(pts)-1), 0.01)
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

// findCuts records every crossing and which cord is in front. Cords are
// compared with themselves too: a knot of wraps crosses its own standing part
// more often than the other line, and skipping that left half of it unresolved.
func findCuts(pts [][]spec.Point, depth [][]float64) [][]crossing {
	// Samples this close together along one cord are the same piece of rope
	// bending, not two parts of it meeting.
	const selfGap = 12

	cuts := make([][]crossing, len(pts))
	for ci := range pts {
		for cj := ci; cj < len(pts); cj++ {
			for i := 0; i+1 < len(pts[ci]); i++ {
				start := 0
				if ci == cj {
					start = i + selfGap
				}
				for j := start; j+1 < len(pts[cj]); j++ {
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
	for ci := range cuts {
		sort.Slice(cuts[ci], func(a, b int) bool { return cuts[ci][a].i < cuts[ci][b].i })
		cuts[ci] = cluster(cuts[ci])
	}
	return cuts
}

// Reading is the over/under sequence along each cord, in path order. Segments
// cannot answer this: two crossings a cord passes the same way merge into one
// piece of rope, and a square knot has exactly that pair.
func Reading(cords []Cord) [][]bool {
	pts := make([][]spec.Point, len(cords))
	depth := make([][]float64, len(cords))
	for ci, c := range cords {
		for _, p := range c.P {
			pts[ci] = append(pts[ci], spec.Point{p.X, p.Y})
			depth[ci] = append(depth[ci], p.Z)
		}
	}
	out := make([][]bool, len(cords))
	for ci, cs := range findCuts(pts, depth) {
		for _, c := range cs {
			out[ci] = append(out[ci], c.over)
		}
	}
	return out
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

	cuts := findCuts(pts, depth)

	var out []spec.Segment

	for ci, c := range cords {
		// Wide enough to cover the drawn cord where it crosses, measured in this
		// cord's own samples so a densely sampled part is not swallowed whole.
		halo := max(int(casingHalo/pace(pts[ci])), 3)

		// Paint order per sample: 0 free, 1 behind, 2 in front. Nearest crossing
		// wins, not frontmost within reach: frontmost lets one pass in front
		// decide a whole run, painting a threaded strand over the coil it threads.
		order := make([]int, len(pts[ci]))
		near := 0
		for k := range order {
			for near+1 < len(cuts[ci]) &&
				abs(cuts[ci][near+1].i-k) <= abs(cuts[ci][near].i-k) {
				near++
			}
			if len(cuts[ci]) == 0 || abs(cuts[ci][near].i-k) > halo {
				continue
			}
			order[k] = 1
			if cuts[ci][near].over {
				order[k] = 2
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
