// SPDX-License-Identifier: Apache-2.0

// Command look renders rope test cases to a page for visual inspection.
//
// Three attempts at this renderer shipped without anyone looking at the
// output. `make look` screenshots the result so that cannot happen again.
package main

import (
	"fmt"
	"os"

	"github.com/theoutdoorprogrammer/riggermortis/internal/render"
	"github.com/theoutdoorprogrammer/riggermortis/internal/rope"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

func outPath() string {
	if p := os.Getenv("LOOK_OUT"); p != "" {
		return p
	}
	return "look.html"
}

func main() {
	cases := []struct {
		name  string
		turns []int
	}{
		{"reef", []int{1, -1}},
		{"reef2", []int{2, -2}},
		{"granny", []int{1, 1}},
	}
	body := ""
	for _, c := range cases {
		var tw []rope.Twist
		for _, t := range c.turns {
			tw = append(tw, rope.Twist{Turns: t})
		}
		cords := rope.Build(rope.Layout{Twists: tw, Radius: 26, Pitch: 74, Tail: 120, Flare: 34})
		g := &spec.Geometry{Width: 620, Height: 300, Cords: []string{"a", "b"},
			Stages: []spec.StageGeometry{{Stage: 1, Segments: rope.Project(cords, 620, 300)}}}
		body += fmt.Sprintf(`<h2 style="font:14px monospace;color:#bd93f9">%s %v</h2><div>%s</div>`,
			c.name, c.turns, render.Stage(g, 0))
	}
	os.WriteFile(outPath(), []byte(
		`<body style="background:#282a36;margin:0;padding:16px;max-width:680px">`+body+`</body>`), 0o644)
	fmt.Println("wrote /tmp/look.html")
}
