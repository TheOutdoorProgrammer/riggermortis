// SPDX-License-Identifier: Apache-2.0

package solve

import (
	"math"
	"sort"

	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

// Flip swaps which cord is on top at the given crossing numbers, labelled by
// column then top to bottom. Path depths cannot do this: each moves several
// crossings and changes how many exist, renumbering the rest.
func Flip(g *spec.Geometry, stage int, which ...int) {
	if g == nil || stage < 0 || stage >= len(g.Stages) {
		return
	}
	segs := g.Stages[stage].Segments

	type ref struct {
		i    int
		x, y float64
		cord string
		over bool
	}
	var ps []ref
	for i, s := range segs {
		if s.Z == 0 || len(s.Points) == 0 {
			continue
		}
		m := s.Points[len(s.Points)/2]
		ps = append(ps, ref{i, m[0], m[1], s.Cord, s.Z == 2})
	}

	type pair struct {
		a, b int
		x, y float64
	}
	var pairs []pair
	for _, a := range ps {
		if !a.over {
			continue
		}
		best, bd := -1, math.Inf(1)
		for j, b := range ps {
			if b.over || b.cord == a.cord {
				continue
			}
			if d := math.Hypot(a.x-b.x, a.y-b.y); d < bd {
				best, bd = j, d
			}
		}
		if best >= 0 && bd < 90 {
			pairs = append(pairs, pair{a.i, ps[best].i, (a.x + ps[best].x) / 2, (a.y + ps[best].y) / 2})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if math.Abs(pairs[i].x-pairs[j].x) > 40 {
			return pairs[i].x < pairs[j].x
		}
		return pairs[i].y < pairs[j].y
	})

	for _, n := range which {
		if n < 1 || n > len(pairs) {
			continue
		}
		p := pairs[n-1]
		segs[p.a].Z, segs[p.b].Z = segs[p.b].Z, segs[p.a].Z
	}
}
