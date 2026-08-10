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
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

func main() {
	root := flag.String("data", "data", "dataset root")
	strict := flag.Bool("strict", false, "treat warnings as failures")
	list := flag.Bool("rules", false, "list registered rules and exit")
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
