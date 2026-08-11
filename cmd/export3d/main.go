// SPDX-License-Identifier: Apache-2.0

// Command export3d dumps a knot's cords as 3D polylines.
//
// The topology is already solved in Go; this hands the same curves to a
// renderer that can show depth, so flat SVG and lit 3D can be compared on
// identical geometry.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/theoutdoorprogrammer/riggermortis/internal/solve"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

type outCord struct {
	ID     string       `json:"id"`
	Points [][3]float64 `json:"points"`
}

func main() {
	id := flag.String("knot", "knot.square", "record id")
	out := flag.String("out", ".look/cords.json", "output file")
	flag.Parse()

	set, err := spec.Load("data")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	r := set.ByID[*id]
	if r == nil {
		fmt.Fprintln(os.Stderr, "no such record:", *id)
		os.Exit(1)
	}

	stages := solve.Cords(r)
	payload := make([][]outCord, 0, len(stages))
	for _, cords := range stages {
		var st []outCord
		for _, c := range cords {
			pts := make([][3]float64, len(c.P))
			for i, p := range c.P {
				pts[i] = [3]float64{p.X, p.Y, p.Z}
			}
			st = append(st, outCord{ID: c.ID, Points: pts})
		}
		payload = append(payload, st)
	}
	b, _ := json.Marshal(payload)
	if err := os.WriteFile(*out, b, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s: %d stages\n", *out, len(payload))
}
