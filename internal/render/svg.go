// SPDX-License-Identifier: Apache-2.0

// Package render turns solved geometry into SVG.
//
// Cord is drawn twice per piece: a wide casing in the page ground, then the
// cord colour on top. Painting in ascending z means a later piece's casing
// erases the one beneath it at a crossing, which is how a strand reads as
// passing over. No masks, clip paths or z-buffer.
package render

import (
	"fmt"
	"strings"

	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

// Casings are painted in the page ground, so this must match it.
const (
	Background = "#282a36"
	Muted      = "#6272a4"
	Accent     = "#bd93f9"
)

var palette = []string{"#ff79c6", "#8be9fd", "#50fa7b", "#ffb86c"}

const (
	cordWidth   = 17.0
	casingWidth = 29.0
	ringWidth   = 25.0
)

// Stage renders the whole knot, calling out the crossing this step introduces.
// The reader always sees both ends and the finished shape, never a fragment.
func Stage(g *spec.Geometry, index int) string {
	if g == nil || index < 0 || index >= len(g.Stages) {
		return ""
	}
	sg := g.Stages[index]

	colour := map[string]string{}
	for i, c := range g.Cords {
		colour[c] = palette[i%len(palette)]
	}

	segs := append([]spec.Segment(nil), sg.Segments...)
	for i := 1; i < len(segs); i++ {
		for j := i; j > 0 && segs[j].Z < segs[j-1].Z; j-- {
			segs[j], segs[j-1] = segs[j-1], segs[j]
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg viewBox="0 0 %g %g" class="knot" role="img" aria-label="stage %d">`,
		g.Width, g.Height, sg.Stage)

	// Highlight rings go under everything, so a called-out crossing glows
	// behind the cord rather than obscuring it.
	for _, s := range segs {
		if s.Stage == sg.Stage && s.Stage != 0 {
			fmt.Fprintf(&b,
				`<path d="%s" fill="none" stroke="%s" stroke-width="%g" stroke-linecap="round" opacity=".9"/>`,
				path(s.Points), Accent, ringWidth+8)
		}
	}
	for _, s := range segs {
		d := path(s.Points)
		fmt.Fprintf(&b,
			`<path d="%s" fill="none" stroke="%s" stroke-width="%g" stroke-linecap="round"/>`,
			d, Background, casingWidth)
		fmt.Fprintf(&b,
			`<path d="%s" fill="none" stroke="%s" stroke-width="%g" stroke-linecap="round"/>`,
			d, colour[s.Cord], cordWidth)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// path emits a Catmull-Rom spline through every point, converted to cubic
// Béziers. Points sit on the curve rather than pulling it, which is what makes
// solved coordinates land where the solver put them.
func path(pts []spec.Point) string {
	if len(pts) < 2 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "M %g %g", pts[0][0], pts[0][1])

	at := func(i int) spec.Point {
		switch {
		case i < 0:
			return pts[0]
		case i >= len(pts):
			return pts[len(pts)-1]
		}
		return pts[i]
	}
	for i := 0; i < len(pts)-1; i++ {
		p0, p1, p2, p3 := at(i-1), at(i), at(i+1), at(i+2)
		c1x, c1y := p1[0]+(p2[0]-p0[0])/6, p1[1]+(p2[1]-p0[1])/6
		c2x, c2y := p2[0]-(p3[0]-p1[0])/6, p2[1]-(p3[1]-p1[1])/6
		fmt.Fprintf(&b, " C %.2f %.2f, %.2f %.2f, %g %g", c1x, c1y, c2x, c2y, p2[0], p2[1])
	}
	return b.String()
}

// Colour returns the palette entry a cord is drawn in, for the legend.
func Colour(g *spec.Geometry, cord string) string {
	for i, c := range g.Cords {
		if c == cord {
			return palette[i%len(palette)]
		}
	}
	return Muted
}
