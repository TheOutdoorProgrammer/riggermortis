// SPDX-License-Identifier: Apache-2.0
//
// Dispatch on the kind field, so a record picks its own schema and the caller
// does not have to name one. This replaces a per-directory table in the
// Makefile that had to be kept in step with the spec by hand.
//
// An explicit conditional rather than a bare disjunction: with `#Component |
// #Line | ...` a single failing branch collapses the whole thing and reports
// every other branch's kind conflict, which buries the actual error.
package riggermortis

#Dispatch: {
	kind!: #Kind

	if kind == "component" {#Component}
	if kind == "line" {#Line}
	if kind == "knot" {#Knot}
	if kind == "rigging" {#Rigging}
	if kind == "rig" {#Rig}
	if kind == "source" {#Source}
}
