# CUE Capabilities

What CUE v0.17.1 can and cannot do for this dataset, tested against the real records in this repository rather than read off a doc page.
Every section is marked **verified by running** (with the command and its real output) or **UNCERTAIN**.

Experiments ran in a throwaway copy of the repo so the working tree stayed clean.
The extra CUE files written for them are reproduced inline and are the deliverable, not the prose.

## Headline

CUE can do referential integrity across the whole dataset.
It can also do every cross-record rule currently written in Go except one, plus the kind dispatch the Makefile does by hand.
It cannot do a bounded graph walk, and it produces worse error messages than the Go rules unless the assertions are shaped carefully.

The suspicion in the brief is correct: CUE was under-used, and by a lot.

## 1. Referential integrity across a whole dataset

**Verified by running. It works.**

The assumed blocker, "CUE validates one file at a time", is a property of how the Makefile invokes CUE, not of CUE.
Passing several data files to `cue vet` unifies them into one value, and the reason that looks useless is that they collide at the root:

```console
$ cue eval ./cue/ data/components/weight-egg.yaml data/components/weight-bullet.yaml -e id
id: conflicting values "weight.egg" and "weight.bullet":
    ./data/components/weight-bullet.yaml:4:5
    ./data/components/weight-egg.yaml:4:5
```

The fix is the `-l` (path) flag, which places each orphan data file at a computed path instead of at the root.
The expression is evaluated against the file's own contents, so `-l 'id'` keys each file by its own `id`, and `-l` repeats for multiple path components:

```console
$ cue vet ./cue/ data/*/*.yaml -l '"records"' -l 'id'
$ echo $?
0
$ cue eval ./cue/ data/*/*.yaml -l '"records"' -l 'id' -e records --out json | jq 'keys|length'
16
```

One command loads all 16 records into a single value at `records["<id>"]`.
Everything after that is ordinary CUE.

### The schema side

```cue
// cue/index.cue
package riggermortis

// Binding the label to the id also proves file id equals index key, for free.
records: [ID=#Id]: #Dispatch & {id: ID}

// The set of ids that actually exist in the loaded dataset.
_ids: {for id, _ in records {(id): true}}

// A ref must be one of the literal ids present.
#ExistingId: or([for id, _ in records {id}])
```

```cue
// cue/refs.cue, which is R005 expressed as a type.
// Definitions unify across files in a package, so this narrows the existing
// #Id-typed ref fields without editing kinds.cue.
// #Knot.objects[].ref is a local object name, not a record id, so it is untouched.
package riggermortis

#Node: ref?:     #ExistingId
#Edge: knot?:    #ExistingId
#Edge: rigging?: #ExistingId
#StrengthClaim: source: #ExistingId
#Caution: sources: [...#ExistingId]
#ValidationEvent: sources?: [...#ExistingId]
#Line: diameters: [...{source: #ExistingId}]
```

There is **no reference cycle**, which was the obvious worry: `records` is constrained by a type derived from the keys of `records`.
CUE terminates because the keys come from the file injection and the constraint applies to values.

### It catches dangling refs at the right path

Pointing `data/rigs/texas.yaml` at a component that does not exist:

```console
$ cue vet ./cue/ data/*/*.yaml -l '"records"' -l 'id'
records."rig.texas".nodes.1.ref: 9 errors in empty disjunction:
records."rig.texas".nodes.1.ref: conflicting values "bait.stick-worm" and "weight.doesnotexist":
```

Nested refs behave identically, for example `records."knot.palomar".strength.claims.0.source` for a bad source inside a strength claim.
That is the whole of **R005**, plus the id-matches-key half of **R007**, in about ten lines.

### Two caveats worth knowing before you commit

The `or([...ids])` type is quadratic in error output.
A dangling ref prints one "conflicting values" line **per existing id**, so at 16 records it is 9 lines and at 600 records it is 600.
For a dataset that is going to grow, either accept the noise, pipe through `head`, or express R005 as a `rules` entry like the ones in section 4 instead of as a type.

Performance is a non-issue.
616 records (the real 16 plus 600 generated) vet in half a second:

```console
$ time cue vet ./cue/ data/*/*.yaml -l '"records"' -l 'id'
cue vet ...  0.59s user 0.03s system 119% cpu 0.511 total
```

## 2. The required field marker `!`

**Verified by running.**

```cue
#Plain: {a: string}
#Req:   {a!: string}
#Opt:   {a?: string}
```

Against `{}`:

```console
$ cue vet ./ empty.json -d '#Plain'
a: incomplete value string:
    ./s.cue:2:13

$ cue vet ./ empty.json -d '#Req'
a: field is required but not present:
    ./s.cue:3:10

$ cue vet ./ empty.json -d '#Opt'
(no output, passes)
```

`a: string` and `a!: string` both reject the missing field, but only `!` says *why* in language a contributor understands.
"incomplete value string" is CUE jargon for "you did not fill this in".
"field is required but not present" is an error message.

Two things `!` does **not** change, both tested:

1. Generated JSON Schema is identical.
   `#Plain` and `#Req` both emit `"required": ["a"]`; `#Opt` emits none.
2. Concreteness checking is already on.
   `cue vet` without `-c` still reported the incomplete value, so `-c` is not what is making the plain form fail.

**Recommendation: switch every non-optional field in `cue/` to `!`.**
It is a pure error-message upgrade with no downstream cost, and it makes the schema self-documenting about intent instead of relying on the reader knowing that a bare `field: type` in CUE means required-but-unfilled.

## 3. Validating many files at once

**Verified by running. The Makefile loop and the `PAIRS` table can both be deleted.**

Two separate discoveries here.

### `-d` validates each data file separately

Passing many data files with `-d` does **not** unify them.
Each is checked against the schema on its own, which is why the id collision from section 1 does not happen:

```console
$ cue vet ./cue/ data/*/*.yaml -d '#Dispatch'
$ echo $?
0
```

That is not a silent no-op.
Introducing one bad enum in one file is caught, and the failing data file is named in the position list:

```console
$ cue vet ./cue/ data/*/*.yaml -d '#Dispatch'
mounting: 3 errors in empty disjunction:
mounting: conflicting values "rigged" and "bolted":
    ./cue/dispatch.cue:7:26
    ./cue/enums.cue:38:15
    ./cue/kinds.cue:15:18
    ./data/components/weight-egg.yaml:6:11
```

### The `kind` field can pick the schema, so `-d '#Component'` per directory is unnecessary

A bare disjunction (`#Record: #Component | #Line | ...`) technically works but is unusable, because a failure in one branch collapses the whole disjunction and CUE reports every *other* branch's `kind` conflict instead of the real cause.
A dangling ref in a rig produced five "conflicting values \"component\" and \"rig\"" style errors and never mentioned the ref.

Explicit dispatch fixes it:

```cue
// cue/dispatch.cue
#Dispatch: {
	kind: #Kind
	if kind == "component" {#Component}
	if kind == "line" {#Line}
	if kind == "knot" {#Knot}
	if kind == "rigging" {#Rigging}
	if kind == "rig" {#Rig}
	if kind == "source" {#Source}
}
```

With that, the same dangling ref reports at `records."rig.texas".nodes.1.ref`, and every conformance fixture and every data file validates against a single `-d '#Dispatch'` with no per-file kind:

```console
$ for f in conformance/invalid/*.yaml; do
    cue vet ./cue/ "$f" -d '#Dispatch' >/dev/null 2>&1 && echo "HOLE $f" || echo "ok $f rejected"
  done
ok bad-enum.yaml rejected
ok bad-id-two-dots.yaml rejected
ok imperial-unit.yaml rejected
ok missing-validation.yaml rejected
ok numeric-external-id.yaml rejected
ok unknown-field.yaml rejected
```

The whole `validate` target collapses to `cue vet ./cue/ data/*/*.yaml -d '#Dispatch'`, and `conformance/manifest.txt` loses its kind column.

One thing `-d '#Dispatch'` does **not** give you is R007's directory check.
`data/knots/weight-egg.yaml` would pass, because nothing ties the directory to the kind.
That check stays in Go, or becomes a two-line shell loop, or moves into `_ids` keyed by directory.

## 4. Cross-record and cross-field constraints

**Verified by running. R001, R002, R032 and R037 all work in CUE.**

Once the dataset is loaded as one value, these are comprehensions.
The key design decision is the pattern constraint on line 3, which is what turns a computed boolean into an assertion:

```cue
// cue/rules.cue
package riggermortis

import "list"

// Visible, so `cue vet` checks it. The pattern constraint is what turns a
// rule into an assertion: any entry that evaluates to false conflicts here.
rules: [string]: [string]: true
rules: {
	// R032: a node that references a component with variants must select a
	// value the component actually offers, on the component's own axis.
	R032: {
		for id, r in records if r.kind == "rig" {
			for i, n in r.nodes {
				for f, _ in n if f == "ref" {
					let c = records[n.ref]
					for cf, _ in c if cf == "variants" {
						"\(id).nodes[\(i)] offers \(c.variants.axis)":
							list.Contains(c.variants.values, n[c.variants.axis])
					}
				}
			}
		}
	}

	// R037: only a soft body can be rigged.
	R037: {
		for id, r in records if r.kind == "rig" {
			let byID = {for n in r.nodes {(n.id): n}}
			for i, e in r.edges if e.rel == "rigged" {
				for f, _ in byID[e.to] if f == "ref" {
					"\(id).edges[\(i)] rigs a soft body": (records[byID[e.to].ref].soft) == true
				}
			}
		}
	}

	// R002: a knot must accept the pin type it ties to.
	R002: {
		for id, r in records if r.kind == "rig" {
			let byID = {for n in r.nodes {(n.id): n}}
			for i, e in r.edges if e.rel == "tied" {
				// Both guards matter: a tied edge may target a line node with no
				// ref (rig.carolina ties to its leader), and CUE calls a missing
				// field "incomplete", not an error.
				for f, _ in e if f == "pin" {
					for g, _ in byID[e.to] if g == "ref" {
						let comp = records[byID[e.to].ref]
						let pins = [for p in comp.pins if p.id == e.pin {p}]
						"\(id).edges[\(i)] pin \(e.pin) exists": (len(pins) == 1)
						if len(pins) == 1 {
							"\(id).edges[\(i)] \(e.knot) accepts \(pins[0].type)":
								list.Contains(records[e.knot].connects.to, pins[0].type)
						}
					}
				}
			}
		}
	}
}

// R001, one hop: a threaded component's declared travel bounds must name a
// node whose component actually blocks passage.
rules: R001: {
	for id, r in records if r.kind == "rig" {
		let byID = {for n in r.nodes {(n.id): n}}
		for i, e in r.edges if e.rel == "threaded" {
			for side in ["toward_rod", "toward_terminal"] {
				for f, bound in e.travel if f == side if bound != "open" {
					for g, _ in byID[bound] if g == "ref" {
						"\(id).edges[\(i)] \(side) stop \(bound) blocks passage":
							(records[byID[bound].ref].blocks_passage) == true
					}
				}
			}
		}
	}
}
```

On the real data this evaluates to a readable audit trail rather than a silent pass:

```console
$ cue eval ./cue/ data/*/*.yaml -l '"records"' -l 'id' -e rules.R002
"rig.carolina.edges[2] pin eye-a exists":                true
"rig.carolina.edges[2] knot.palomar accepts closed-eye": true
"rig.carolina.edges[4] pin eye exists":                  true
"rig.texas.edges[1] pin eye exists":                     true
"rig.texas.edges[1] knot.palomar accepts closed-eye":    true
```

Each rule was then broken deliberately in `data/rigs/texas.yaml` and each was caught.
The rule key *is* the error message, which is why the naming matters:

```console
--- R037: rig the sinker instead of the bait ---
rules.R037."rig.texas.edges[2] rigs a soft body": conflicting values false and true

--- R032: mass_g 7.0 -> 6.9, not on weight.bullet's axis ---
rules.R032."rig.texas.nodes[1] offers mass_g": conflicting values false and true

--- R002: tie the palomar to the shank pin instead of the eye ---
rules.R002."rig.texas.edges[1] knot.palomar accepts shank": conflicting values false and true

--- R001: point the sinker's terminal stop at the bait ---
rules.R001."rig.texas.edges[0] toward_terminal stop bait blocks passage": conflicting values false and true
```

That is four of the eleven Go rules, including all three that the brief singled out as "cross-record".
Note that R001 is not a graph walk in the Go version either.
`checkStop` is a single hop: name a node, resolve its component, read `blocks_passage`.

### Three traps that will cost an afternoon if nobody writes them down

**A missing field is "incomplete", not an error.**
`rig.carolina` ties a knot to its leader, which is a line node with no `ref`.
Reading `byID[e.to].ref` unguarded does not fail loudly, it makes the whole instance incomplete, and `cue vet` degrades to the useless `some instances are incomplete; use the -c flag to show errors`.
Every optional-field access inside a comprehension needs a `for f, _ in x if f == "field"` presence guard.
This is CUE's equivalent of the `if comp == nil { continue }` lines all over `rules.go`, and it is strictly more annoying.

**`x & true` is not an assertion when `x` has a default.**
Adding `#Component: soft: *false | bool` so that absent means false is correct and necessary, but it silently defeats `& true`:

```console
$ cue eval -e a -e b   # where x: *false | bool, a: x & true, b: x == true
// a
true
// b
false
```

Unification picks the branch that matches instead of comparing.
Use `== true` plus a `[string]: true` pattern constraint, as above.

**Optional booleans need defaults.**
`soft?: bool` means a cross-record rule reading `comp.soft` on a swivel gets incomplete rather than false.
Give every "absent means false" boolean an explicit `*false | bool`.

## 5. Graph traversal

**Verified by running. Confident no.**

CUE rejects self-referential definitions outright, regardless of a decreasing counter that would obviously terminate:

```cue
#Chain: {
	n: int
	let m = n - 1
	if n > 0 {next: #Chain & {n: m}}
	if n <= 0 {done: true}
}
walk: #Chain & {n: 3}
```

```console
$ cue eval ./ -e walk
walk.next.next.next: structural cycle
```

The detector fires on the *shape* of the definition, not on actual depth, so it triggers at depth 3 and would trigger at depth 1.
There is no bounded-recursion escape hatch.

The only way to walk a graph of unknown depth in CUE is to unroll it by hand to a fixed maximum, which means writing `hop1`, `hop2`, `hop3` and accepting silent under-checking past the limit.
Do not do that.

This is not a real loss for this project, because nothing in `rules.go` actually walks.
It becomes a loss the day the spec grows a rule like "no cycle in a rig's edge graph" or "every node is reachable from the main line".
Those stay in Go permanently.

## 6. Generating artefacts

**Verified by running. JSON Schema and OpenAPI both work today, against the unmodified `cue/`.**

```console
$ cue def ./cue/ --out openapi | head
{
    "openapi": "3.0.0",
    "info": {
        "title": "SPDX-License-Identifier: Apache-2.0\n\nFields every record carries, ...",
        "version": "no version"
    },
    "paths": {},
    "components": {
        "schemas": {
            "Action": {
                "type": "object",
                "required": ["verb", "subject"],
```

```console
$ cue def ./cue/ -e '#Component' --out jsonschema | head
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$defs": {
        "BendStyle": {"enum": ["round", "octopus", "circle", ...]},
```

Enums come through as JSON Schema `enum`, required fields as `required`, closed structs as `additionalProperties: false`.
This is a genuinely free artefact.

Three constraints, all found by running it:

1. **The `if kind == ...` dispatch breaks OpenAPI export.**
   `cue def ./cue/ --out openapi` on the package *with* `#Dispatch` fails with `unresolved disjunction "component" | "line" | ...`.
   JSON Schema survives it, OpenAPI does not.
2. **`--out openapi` needs a package of definitions, not a selected value.**
   `cue def ./cue/ -e '#Component' --out openapi` fails with `unsupported top-level field "schema_version"`.
   Export the package, not one definition.
3. **The doc comments become the schema title and description verbatim**, SPDX header and all.
   Whoever wires this into CI should expect `"title": "SPDX-License-Identifier: Apache-2.0\n\n..."` and either live with it or move the licence header to a non-doc position.

Points 1 and 2 together force a **two-package layout**, which is the single most important structural conclusion in this document:

- `cue/` keeps the publishable, self-contained schema (enums, common, kinds, plus `!` markers).
  It generates JSON Schema and OpenAPI, and it is what a third party consumes.
- a second package, say `validate/`, holds `index.cue`, `dispatch.cue`, `refs.cue` and `rules.cue`.
  These are dataset-scoped and meaningless without the data.

Keeping them together also poisons the publishable schema in a subtler way: `#ExistingId` is `or([])` when no data is loaded, and `cue eval ./cue/ -e '#ExistingId'` on an empty dataset returns an internal export error.
A downstream consumer who imports the schema and evaluates a `#Node` would hit that.

## 7. CUE modules and publishing

**Partly verified by running, partly UNCERTAIN.**

Verified locally:

```console
$ cue mod resolve github.com/foo/bar
registry.cue.works/github.com/foo/bar
```

The default registry is the CUE Central Registry, `registry.cue.works`, and `cue mod publish`, `cue login`, `cue mod get`, `cue mod mirror` and `cue mod tidy` all exist in v0.17.1.
Modules are published as OCI artifacts, so any OCI registry works via `$CUE_REGISTRY`, and `cue mod publish --out` can write an OCI Image Layout directory for a registry-less workflow.

This module is not publishable as it stands, and CUE says exactly why:

```console
$ cue mod publish --dry-run v0.0.1
publishing a module requires a source field in cue.mod/module.cue;
choose a source with 'cue mod edit --source'
```

That is a one-line fix (`cue mod edit --source git`), after which the module name is already correct: `github.com/theoutdoorprogrammer/riggermortis`.
A third party would then run `cue mod get github.com/theoutdoorprogrammer/riggermortis@v0` and `cue vet` their own YAML against `#Dispatch`.

**UNCERTAIN:** I could not complete an end-to-end publish-and-consume round trip.
`cue mod registry` (the local in-memory registry) blocks as a server, and the network calls to `registry.cue.works` did not return inside the tool timeout, so I verified the commands and the resolution path but not a real fetch.
Before promising this to anyone, do one real `cue mod publish --dry-run --json` with a source field set, and one `cue mod get` from a clean module.

The other thing worth flagging: publishing the schema is only useful if the published package is the *self-contained* one from section 6.
Publishing a package whose `#Node.ref` is constrained to ids that exist in this repository's `data/` would be actively broken for a consumer.

## 8. Testing the schema itself

**Verified by running for the mechanism. UNCERTAIN that any established tooling exists.**

There is no built-in "this must fail" test runner in CUE, and I found no standard pattern for one.
`cue vet` exit codes in a loop, which is what `make conformance` already does, is the state of the art.
The one improvement available today is dropping the kind column from `conformance/manifest.txt`, since `-d '#Dispatch'` handles every fixture (section 3).

Two things do change with the whole-dataset model, and they matter:

1. **Invalid fixtures must never be loaded into the index.**
   A fixture with a dangling ref, added to `data/`, would break the entire load rather than fail in isolation.
   Conformance fixtures stay a separate per-file `-d` run.
   This is a good reason to keep `-d` in the Makefile even after the loop is deleted.
2. **The cross-record rules need their own negative fixtures**, and those are whole-dataset fixtures, not single files.
   A directory per scenario (`conformance/rulesets/r002-wrong-pin/`) vetted with the same `-l` invocation is the natural shape.
   That is exactly what I did by hand for section 4, and it worked.

The tooling layer (`cue cmd`, `*_tool.cue`) does work and can drive this:

```cue
package micro

import (
	"tool/exec"
	"tool/cli"
)

command: check: {
	vet: exec.Run & {cmd: ["cue", "version"]}
	say: cli.Print & {$after: vet, text: "tool layer ran"}
}
```

```console
$ cue cmd check ./
cue version v0.17.1
```

It is a real workflow engine with dependency ordering via `$after`.
Whether it earns its place over the existing Makefile is a judgement call, and my read is no: the Makefile is legible to contributors who do not know CUE, and `cue cmd` is not.

## 9. Things a project of this shape would normally use

**Verified by running unless noted.**

**Closed structs already work and are doing real work.**
`manufacturer: acme` appended to a component is rejected with `field not allowed`, at the right line.
This is the mechanism behind the CLAUDE.md promise that vendor fields are structurally impossible, and it is holding.

**Defaults (`*x | y`)** are under-used and are needed for the cross-record rules to work at all (section 4).

**`cue fmt` is already clean**: `cue fmt --check ./cue/` exits 0.
Worth adding to CI so it stays that way.

**Attributes (`@rule(R005)`)** parse fine and are the obvious way to carry the spec's rule IDs on the fields they constrain, which is the ADR 0003 "check the two against each other" idea made mechanical.
They are stripped from generated JSON Schema, so they cost nothing downstream, and they are readable from Go via `cuelang.org/go/cue`.
This is the cleanest available answer to keeping `docs/spec.md` and the schema in sync.

**`@embed`** exists (`@extern(embed)` at file level, then `@embed(file=...)`) and can pull a YAML file into evaluation directly.
Not obviously useful here, since `-l` already does the loading, but it is an alternative to the `-l` invocation that puts the file list under version control instead of in the Makefile.
**UNCERTAIN**, not tested end to end.

**Not used and probably should be:** `cue vet -E` / `--all-errors`.
The default truncates, and for a validator whose job is "find out how many hand-authored records are wrong", truncation is the wrong default.

## What should move into CUE

### Should move, with confidence

| Rule | Why |
| --- | --- |
| **R005** every reference resolves | Section 1. Ten lines of type narrowing, no cycle, catches nested refs. This is the largest single chunk of `rules.go`. |
| **R002** a knot must accept the pin type | Section 4. Works, with a self-describing error. |
| **R032** a selected variant exists | Section 4. Works, including the dynamic axis lookup `n[c.variants.axis]`. |
| **R037** only a soft body can be rigged | Section 4, once `soft` gets a `*false \| bool` default. |
| **R001** threaded travel bounds block passage | Section 4. It is a one-hop lookup, not a walk, despite living in the "graph" bucket. |
| **R010** tier C claim carries a source | Trivial in CUE: `if tier == "C" {source: #Id}`. Not separately tested, but it is the same conditional shape as `#Dispatch`, which is. |
| **R011** a strength claim states n | Already enforced: `#StrengthClaim.n` is `int & >=0`, non-optional. This rule is dead code in Go today. |
| **R039** do-not-use caution carries a source | Same conditional shape as R010: `if severity == "do-not-use" {sources: [_, ...#Id]}`. |
| id matches filename / index key | Free side effect of `records: [ID=#Id]: {id: ID}`. |
| kind dispatch (the Makefile `PAIRS` table) | Section 3. `-d '#Dispatch'` replaces the loop and the table. |

That is nine of eleven registered rules plus the Makefile's kind table.

### Should stay in Go

| Rule | Why |
| --- | --- |
| **R007** kind matches the directory | CUE sees values, not file paths. `-l` can key by `id` but not by directory. A shell loop or Go keeps this. |
| **R040** citation debt is reported | It is a **warning**, and CUE has one outcome: valid or not. Anything that must report without failing needs a host program. |
| **R023** nothing publishes unvalidated | Same reason. Warning, not error. |
| Any future graph traversal | Section 5. Structural cycle detection makes this impossible, permanently. |

### The honest cost

Three things get worse, and they should be weighed rather than waved away:

1. **Error messages.**
   Go produces `rig.texas: hook.offset-worm does not terminate on shank (pin hook.eye)`.
   CUE produces `rules.R002."rig.texas.edges[1] knot.palomar accepts shank": conflicting values false and true`.
   The second is decipherable, but "conflicting values false and true" is boilerplate on every single rule failure, and the `or()` form of R005 emits one line per existing id.
2. **Nil-guarding is worse, not better.**
   Every optional-field access needs a `for f, _ in x if f == "field"` presence comprehension, and forgetting one does not fail, it produces `some instances are incomplete`, which points at nothing.
   Go's `if comp == nil { continue }` is uglier to read and far easier to get right.
3. **The rules are harder for a contributor to extend.**
   `rules.go` can be read by anyone who knows Go.
   The R002 comprehension above cannot be read by anyone who does not know CUE well, and the `& true` versus `== true` trap is not discoverable.

### Recommendation

Move R005, R001, R002, R032, R037, R010, R039 and the kind dispatch into CUE.
Keep R007, R023 and R040 in Go, and keep the Go program itself, because it is still the only thing that can emit warnings, check paths, and one day walk a graph.
Delete R011, which the schema already enforces.

Split `cue/` into a publishable schema package and a dataset-scoped `validate/` package, or JSON Schema generation and module publishing both break.
Switch required fields to `!` while doing it.

The net effect is that `rules.go` shrinks to three rules and a reporter, the Makefile's `validate` target becomes one command, and the project gains generated JSON Schema and OpenAPI for free.
