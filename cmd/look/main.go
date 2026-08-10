// SPDX-License-Identifier: Apache-2.0

// Command look renders every stage of a knot to a page for visual inspection.
//
// Three attempts at this renderer shipped without anyone looking at the
// output. `make look` screenshots the result so that cannot happen again.
package main

import (
	"fmt"
	"os"

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
	if r == nil {
		fmt.Fprintln(os.Stderr, "no such record:", target)
		os.Exit(1)
	}
	g := solve.Geometry(r)
	if g == nil {
		fmt.Fprintln(os.Stderr, target, "solves to no crossings")
		os.Exit(1)
	}

	body := ""
	for i, st := range g.Stages {
		prose := ""
		if i < len(r.Stages) {
			prose = r.Stages[i].Prose
		}
		body += fmt.Sprintf(
			`<p style="font:13px/1.4 monospace;color:#bd93f9;margin:14px 0 2px">%d &middot; `+
				`<span style="color:#6272a4">%s</span></p><div>%s</div>`,
			st.Stage, prose, render.Stage(g, i))
	}
	os.WriteFile(outPath(), []byte(
		`<body style="background:#282a36;margin:0;padding:14px;max-width:680px">`+body+`</body>`), 0o644)
	fmt.Println("wrote", outPath())
}
