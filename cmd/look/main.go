// SPDX-License-Identifier: Apache-2.0

// Command look renders a knot's stages for visual inspection, optionally
// numbering each crossing so a reader can name which one is wrong.
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/theoutdoorprogrammer/riggermortis/internal/render"
	"github.com/theoutdoorprogrammer/riggermortis/internal/solve"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

func outPath() string {
	if p := os.Getenv("LOOK_OUT"); p != "" {
		return p
	}
	return "look.html"
}

type mark struct {
	x, y float64
	over string
}

// crossings pairs each over piece with the under piece it sits on, so one
// label lands per crossing rather than one per strand.
func crossings(g *spec.Geometry, stage int) []mark {
	type piece struct {
		x, y float64
		cord string
		over bool
	}
	var ps []piece
	for _, s := range g.Stages[stage].Segments {
		if s.Z == 0 || len(s.Points) == 0 {
			continue
		}
		m := s.Points[len(s.Points)/2]
		ps = append(ps, piece{m[0], m[1], s.Cord, s.Z == 2})
	}
	var out []mark
	for _, a := range ps {
		if !a.over {
			continue
		}
		best, bd := -1, math.Inf(1)
		for i, b := range ps {
			if b.over || b.cord == a.cord {
				continue
			}
			d := math.Hypot(a.x-b.x, a.y-b.y)
			if d < bd {
				best, bd = i, d
			}
		}
		if best >= 0 && bd < 90 {
			out = append(out, mark{(a.x + ps[best].x) / 2, (a.y + ps[best].y) / 2, a.cord})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if math.Abs(out[i].x-out[j].x) > 40 {
			return out[i].x < out[j].x
		}
		return out[i].y < out[j].y
	})
	return out
}

func main() {
	target := "knot.square"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}
	set, err := spec.Load("data")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	r := set.ByID[target]
	g := solve.Geometry(r)
	if g == nil {
		fmt.Fprintln(os.Stderr, target, "solves to no crossings")
		os.Exit(1)
	}

	var b strings.Builder
	b.WriteString(`<body style="background:#282a36;margin:0;padding:14px;max-width:700px">`)
	for i, st := range g.Stages {
		prose := ""
		if i < len(r.Stages) {
			prose = r.Stages[i].Prose
		}
		svg := render.Stage(g, i)

		// Label beside the crossing, never on it: the whole point is to see
		// which strand is on top.
		var lbl strings.Builder
		for n, m := range crossings(g, i) {
			// Park labels on a rail clear of the knot so nothing overlaps rope.
			ly := 22.0
			if m.y > g.Height/2 {
				ly = g.Height - 22
			}
			col := "#ff79c6"
			if m.over == "b" {
				col = "#8be9fd"
			}
			lbl.WriteString(fmt.Sprintf(
				`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1.5" opacity=".5" stroke-dasharray="3 4"/>`+
					`<circle cx="%.1f" cy="%.1f" r="14" fill="%s" stroke="#282a36" stroke-width="2"/>`+
					`<text x="%.1f" y="%.1f" font-family="monospace" font-size="15" font-weight="700"`+
					` fill="#282a36" text-anchor="middle" dominant-baseline="central">%d</text>`,
				m.x, m.y, m.x, ly, col, m.x, ly, col, m.x, ly, n+1))
		}
		svg = strings.Replace(svg, `</svg>`, lbl.String()+`</svg>`, 1)

		fmt.Fprintf(&b,
			`<p style="font:13px/1.4 monospace;color:#bd93f9;margin:16px 0 2px">%d &middot; `+
				`<span style="color:#6272a4">%s</span></p><div>%s</div>`,
			st.Stage, prose, svg)
	}
	b.WriteString(`</body>`)
	os.WriteFile(outPath(), []byte(b.String()), 0o644)
	fmt.Println("wrote", outPath())
}
