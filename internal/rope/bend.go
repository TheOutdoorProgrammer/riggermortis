// SPDX-License-Identifier: Apache-2.0

package rope

// wp is an authored waypoint: position plus depth, splined before use.
type wp struct{ x, y, z float64 }

// BuildBend is the SPECIAL CASE for the two-half-knot bends (square, granny):
// authored waypoints, not solved, because counter-phased helices read as the
// twisted pair they are. Crossing signs still pick every over and under.
func BuildBend(s1, s2, n int, tight bool) []Cord {
	m := 1.0
	if s1 < 0 {
		m = -1
	}
	reef := (s1 < 0) != (s2 < 0)

	var a, b []wp
	switch {
	case n <= 1:
		a = bendCross(m)
	case n == 2:
		a = bendHalfKnot(m)
	case n == 3:
		a = bendCrossedEnds(m, reef)
	case reef:
		// The clasp is point-symmetric, not mirror-symmetric: mirroring
		// would land paired crossings on the axis at the same x, where
		// the left-to-right reading cannot order them.
		a = bendReefClasp(m, tight)
		b = bendRotate(a)
	default:
		a = bendStacked(m)
	}
	if b == nil {
		b = bendMirror(a)
	}
	return []Cord{bendCord("a", a), bendCord("b", b)}
}

// Cord a enters lower left and crosses the other end.
func bendCross(m float64) []wp {
	return []wp{
		{-300, 85, 0}, {-120, 45, 0}, {0, 5, 14 * m}, {130, -55, 0}, {230, -90, 0},
	}
}

// First half knot: one full wrap, the ends swap sides.
func bendHalfKnot(m float64) []wp {
	return []wp{
		{-290, 90, 0}, {-120, 58, 0}, {0, 42, 14 * m}, {-68, -4, 0},
		{0, -50, -14 * m}, {110, -80, 0}, {230, -98, 0},
	}
}

// Half knot seated, with the ends crossed again above it.
func bendCrossedEnds(m float64, reef bool) []wp {
	z3 := 14 * m
	if !reef {
		z3 = -z3
	}
	return []wp{
		{-290, 110, 0}, {-120, 78, 0}, {0, 52, 14 * m}, {-65, 10, 0},
		{0, -32, -14 * m}, {95, -62, 0}, {45, -92, 0}, {0, -102, z3},
		{-135, -125, 0},
	}
}

// Both half knots tied the same hand, not yet drawn up. Granny only: the
// reef's stages 4 and 5 are both bendReefClasp, because a vertical stack
// puts every crossing at the same x and the reading cannot order them.
func bendStacked(m float64) []wp {
	return []wp{
		{-280, 135, 0}, {-110, 100, 0}, {0, 74, 14 * m}, {-70, 34, 0},
		{0, -2, -14 * m},
		{-42, -36, 0}, {0, -68, -14 * m}, {-70, -104, 0},
		{0, -140, 14 * m}, {110, -165, 0}, {230, -180, 0},
	}
}

// bendReefClasp seats the square knot: two interlocked bights, each bend
// collaring the other cord's paired legs, ends home beside their standing
// parts. Cord a weaves O,U,O,U,O,U left to right; loose opens it up.
func bendReefClasp(m float64, tight bool) []wp {
	z := 13 * m
	w := []wp{
		{-275, 41, 0}, {-182, 34, 0},
		{-151, 32, z}, {-105, 31, z}, // standing part over b's bend
		{-72, 30, 0},
		{-6, 31, -z}, {44, 55, -z}, // under b's returning leg
		{77, 58, 0}, {96, 58, 0},
		{113, 47, z}, {129, 19, z}, // own bend over that leg
		{144, -6, 0},
		{150, -19, -z}, {138, -44, -z}, // own bend under b's standing part
		{116, -56, 0}, {77, -58, 0},
		{19, -54, z}, {-30, -33, z}, // working end back over it
		{-72, -30, 0},
		{-105, -31, -z}, {-151, -32, -z}, // and under b's bend, home
		{-182, -34, 0}, {-275, -41, 0},
	}
	if tight {
		return w
	}
	for i := range w {
		w[i].x *= 1.08
		w[i].y *= 1.55
	}
	w[0].y, w[len(w)-1].y = 78, -78
	return w
}

// bendMirror reflects a path for the second cord: x for side, z so every
// shared crossing resolves the other way.
func bendMirror(w []wp) []wp {
	out := make([]wp, len(w))
	for i, p := range w {
		out[i] = wp{-p.x, p.y, -p.z}
	}
	return out
}

// bendRotate spins a path a half turn in the plane for the second cord,
// keeping z so each crossing keeps its winner.
func bendRotate(w []wp) []wp {
	out := make([]wp, len(w))
	for i, p := range w {
		out[i] = wp{-p.x, -p.y, p.z}
	}
	return out
}

func bendCord(id string, w []wp) Cord {
	pts := make([]V3, len(w))
	for i, p := range w {
		pts[i] = V3{p.x, p.y, p.z}
	}
	return Cord{ID: id, P: sampleCR3(pts, 26*(len(pts)-1))}
}
