---
dusk: v1alpha1
namespace: stout
kind: repository
name: riggermortis
title: Riggermortis
attributes:
  language: go
  visibility: public
  licence: Apache-2.0 and ODbL-1.0
---

An open reference for every way a recreational angler can rig terminal tackle, in which rigs and knots are **records** and the diagrams are generated from them.
The premise is that everything which renders a recognisable practical knot uses hand-authored splines, everything which generates geometry works on mathematical knots or pure physics, and nobody generates a practical knot from a structured description of how it is tied.
The dataset is the expensive part; the renderer, the search index and the site are all projections of it.

## The model

One primitive holds it together: the **typed pin**.
A component exposes pins with types, a knot declares which pin types it may terminate on, and a rig says which pin a given knot ties to.
Because knots and rigs share that vocabulary a rig is a netlist and validation is electrical rule checking, which is where the whole validation design is borrowed from.
Design the two schemas independently and they will not compose.

There are six kinds (`component`, `line`, `knot`, `rigging`, `rig`, `source`), one record per file under `data/`, ids of the form `{category}.{name}` with exactly one dot, and a file named after its id minus the category.
The shapes live in `cue/`, the full prose spec in `docs/spec.md`, and `docs/coverage.md` tracks what exists against what does not.

## Validation is three layers doing three different jobs

`make validate` runs `cue vet ./validate/ data/*/*.yaml -l '"records"' -l 'id'`.
Binding the label to each record's own id places records side by side instead of collapsing onto each other at the root, and proves for free that a record's id matches the key it lands under.
`cue/dispatch.cue` then switches on `kind` so a record picks its own schema, written as explicit conditionals rather than a disjunction because one failing branch of a disjunction reports every other branch's kind conflict and buries the real error.
Definitions are closed, so an unregistered field is an error and a vendor field has nowhere to live.

`make conformance` checks the inverse.
Every fixture under `conformance/invalid/` **must be rejected**, and one that validates is reported as a hole in the spec rather than as a passing test.

`make rules` runs `cmd/rig`, the Go pass for everything CUE cannot express: reference resolution, graph walks, path checks.
Each rule carries the number it has in the spec's rules table so the two can be diffed, and severity is split between error and warning deliberately.
The warnings are the project's honest debt, unvalidated records and outstanding citations, and CI publishes the count on every merge so it cannot rot quietly.

`make check` is all three and is what CI runs.
`make schema` regenerates the published JSON Schema, one per kind, for consumers who do not use CUE.
`make site` compiles records into static HTML and SVG; Go is a build-time compiler here, not a server, and no geometry ever runs client-side.

## What will bite

**The licence boundary is a directory.**
`data/` is ODbL-1.0 and everything else is Apache-2.0, and every file under `data/` carries an SPDX line as its first line because records get copied out of the repository and a file that loses its directory loses its terms.
Moving a file across that line changes its licence.

**Vendor neutrality and SI units are structural, not editorial.**
Record the physical property with the unit in the field name, `rating_kg` rather than a trade name or an imperial suffix, because no field exists for either and the conformance corpus has a fixture for both.

**Never raise a record's `validation.status` for your own work.**
`field-tested` means a human who fishes confirmed it.
The project's own stated largest gap is that software can prove the encoded knot is the knot claimed, and cannot prove that a person following the written steps arrives there.

**The recurring schema failure is over-constraining, not under-constraining.**
"An ordered list of components along a line" forbade the umbrella rig and "a knot is always an edge" forbade the dropper loop.
When adding a constraint, go and find the specimen that breaks it first.

Decisions live in `adr/` in MADR form, indexed by `adr/README.md`, and a reversal is a new record rather than an edit.
Note that `README.md` still opens with a research-phase status and "no application code yet", which `cmd/`, `internal/`, `validate/` and `site/` have long since overtaken.
