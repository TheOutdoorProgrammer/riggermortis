// SPDX-License-Identifier: Apache-2.0

// Command site builds the static site.
//
// Build time only. Records in, HTML and SVG out, nothing running in production.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/theoutdoorprogrammer/riggermortis/internal/site"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

func main() {
	data := flag.String("data", "data", "dataset root")
	out := flag.String("out", "site", "output directory")
	flag.Parse()

	set, err := spec.Load(*data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load:", err)
		os.Exit(2)
	}
	if err := site.Build(set, *out); err != nil {
		fmt.Fprintln(os.Stderr, "build:", err)
		os.Exit(1)
	}

	knots := 0
	for _, r := range set.All {
		if r.Kind == "knot" {
			knots++
		}
	}
	fmt.Printf("built %s: %d knot pages\n", *out, knots)
}
