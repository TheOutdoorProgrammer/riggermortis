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

	var a []wp
	switch {
	case n <= 1:
		a = bendCross(m)
	case n == 2:
		a = bendHalfKnot(m)
	case n == 3:
		a = bendCrossedEnds(m, reef)
	case tight && reef:
		a = bendReefFlat(m)
	default:
		a = bendStacked(m, reef)
	}
	return []Cord{bendCord("a", a), bendCord("b", bendMirror(a))}
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

// Both half knots tied, not yet drawn up: a reef's second wrap mirrors the
// first, a granny's repeats it.
func bendStacked(m float64, reef bool) []wp {
	w := []wp{
		{-280, 135, 0}, {-110, 100, 0}, {0, 74, 14 * m}, {-70, 34, 0},
		{0, -2, -14 * m},
	}
	if reef {
		return append(w,
			wp{42, -36, 0}, wp{0, -68, 14 * m}, wp{70, -104, 0},
			wp{0, -140, -14 * m}, wp{-110, -165, 0}, wp{-230, -180, 0})
	}
	return append(w,
		wp{-42, -36, 0}, wp{0, -68, -14 * m}, wp{-70, -104, 0},
		wp{0, -140, 14 * m}, wp{110, -165, 0}, wp{230, -180, 0})
}

// The seated square knot: each cord a bight, legs through the other bight's
// eye and home alongside their own standing part. Six alternating crossings.
func bendReefFlat(m float64) []wp {
	return []wp{
		{-240, 22, 0}, {-185, 22, 0}, {-138, 22, 13 * m}, {-85, 22, 0},
		{-32, 25, 0}, {0, 36, -13 * m}, {32, 47, 0}, {85, 52, 0},
		{122, 42, 6 * m}, {145, 20, 13 * m}, {152, 0, 0},
		{145, -20, -13 * m}, {122, -42, -6 * m}, {85, -52, 0},
		{32, -47, 0}, {0, -36, 13 * m}, {-32, -25, 0}, {-85, -22, 0},
		{-138, -22, -13 * m}, {-185, -22, 0}, {-240, -22, 0},
	}
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

func bendCord(id string, w []wp) Cord {
	pts := make([]V3, len(w))
	for i, p := range w {
		pts[i] = V3{p.x, p.y, p.z}
	}
	return Cord{ID: id, P: sampleCR3(pts, 26*(len(pts)-1))}
}
