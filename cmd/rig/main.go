// SPDX-License-Identifier: Apache-2.0

// Command rig validates the riggermortis dataset.
//
// CUE checks shape, enums and grammar. This checks everything that needs to
// resolve a reference or walk a graph, which is where the rules with teeth are.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/theoutdoorprogrammer/riggermortis/internal/rules"
	"github.com/theoutdoorprogrammer/riggermortis/internal/solve"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

func main() {
	root := flag.String("data", "data", "dataset root")
	strict := flag.Bool("strict", false, "treat warnings as failures")
	list := flag.Bool("rules", false, "list registered rules and exit")
	cross := flag.String("crossings", "", "report the crossing sequence of a knot id and exit")
	flag.Parse()

	if *list {
		for _, r := range rules.All() {
			fmt.Printf("%s  %-8s %s\n", r.ID, r.Severity, r.Summary)
		}
		return
	}

	set, err := spec.Load(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load:", err)
		os.Exit(2)
	}

	if *cross != "" {
		r := set.ByID[*cross]
		if r == nil {
			fmt.Fprintln(os.Stderr, "no such record:", *cross)
			os.Exit(2)
		}
		g := solve.Geometry(r)
		if g == nil {
			fmt.Fprintln(os.Stderr, *cross, "solves to no crossings")
			os.Exit(2)
		}
		for i := range g.Stages {
			for _, rd := range solve.ReadCrossings(r, i) {
				fmt.Printf("stage %d cord %s  n=%d  %v\n",
					g.Stages[i].Stage, rd.Cord, len(rd.Sequence), rd.Sequence)
			}
		}
		return
	}

	var errs, warns int
	for _, f := range rules.Run(set) {
		fmt.Printf("%-9s %-6s %s: %s\n", f.Rule, f.Severity, f.Path, f.Message)
		if f.Severity == rules.Error {
			errs++
		} else {
			warns++
		}
	}

	fmt.Printf("\n%d records, %d rules, %d errors, %d warnings\n",
		len(set.All), len(rules.All()), errs, warns)

	if errs > 0 || (*strict && warns > 0) {
		os.Exit(1)
	}
}
