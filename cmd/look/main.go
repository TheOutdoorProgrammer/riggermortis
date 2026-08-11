// SPDX-License-Identifier: Apache-2.0

// Command look renders a knot's stages for visual inspection. Visual work gets
// looked at, so this exists to put every stage on one page at size.
package main

import (
	"fmt"
	"os"
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

	var b strings.Builder
	b.WriteString(`<body style="background:#282a36;margin:0;padding:14px;max-width:700px">`)
	for i, st := range g.Stages {
		prose := ""
		if i < len(r.Stages) {
			prose = r.Stages[i].Prose
		}
		fmt.Fprintf(&b,
			`<p style="font:13px/1.4 monospace;color:#bd93f9;margin:16px 0 2px">%d &middot; `+
				`<span style="color:#6272a4">%s</span></p><div>%s</div>`,
			st.Stage, prose, render.Stage(g, i))
	}
	b.WriteString(`</body>`)
	if err := os.WriteFile(outPath(), []byte(b.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("wrote", outPath())
}
