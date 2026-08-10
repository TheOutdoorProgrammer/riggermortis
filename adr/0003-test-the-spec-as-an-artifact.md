# 3. Test the spec itself, not only its implementations

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

[ADR 0002](0002-schema-enforcement-stack.md) settled how records are validated against the spec.
It did not answer a different and harder question: **how do we know a change to the spec itself is safe before it goes live?**

These are not the same thing.

- An **implementation test** asks whether the validator correctly rejects a bad record.
- A **spec test** asks whether the spec is internally coherent, whether it still admits everything it should, whether it still forbids everything it should, and what a proposed edit breaks downstream.

A validator can be flawless and the spec still wrong.
The spec is the artifact that other people build against, published under a reciprocal licence, so a silent narrowing invalidates data we do not control.

The obstacle is that **prose cannot be tested.**
Thirty rules and thirty enumerations written as markdown are unrunnable, so the first requirement is that the spec exist as a machine-readable artifact.
That is a stronger argument for the CUE decision in ADR 0002 than validation ever was, and it is worth stating plainly: the spec is source code that happens to render as documentation.

## What actually goes wrong with a spec

Naming the failure modes first, because each mechanism below targets specific ones.

| # | Failure | Example |
| --- | --- | --- |
| F1 | Internal inconsistency | A rule references an enum value the registry does not define. |
| F2 | Orphan definitions | An enum is defined and no field ever uses it. Dead vocabulary misleads authors. |
| F3 | Under-constrained | Nothing forbids a `threaded` component from declaring pins, so eventually one does. |
| F4 | **Over-constrained** | The spec forbids something legitimate, and nobody finds out until a real rig cannot be expressed. |
| F5 | Unknown blast radius | A field is renamed and no one can say which rules, kinds, examples, and records are affected. |
| F6 | Silent breaking change | An edit invalidates previously valid data without a version bump. |
| F7 | Doc and code drift | The rules table and the enforced behaviour disagree. |

F4 is the one that has already bitten this project twice: "an ordered list of components along a line" forbade the umbrella rig, and "a knot is always an edge" forbade the dropper loop.
Both were caught by hand-authoring an awkward specimen.
That is not a repeatable process, which is exactly what this ADR is for.

## Decision Outcome

Six mechanisms, each targeting named failures.

### 1. A language-neutral conformance corpus

The centrepiece.
A directory of records, each labelled with its expected outcome:

```text
conformance/
├── manifest.yaml
├── valid/
│   └── rig-carolina-minimal.yaml
└── invalid/
    ├── R001-threaded-weight-unbounded.yaml
    └── R014-loop-knot-used-as-edge.yaml
```

The manifest declares, for every fixture, whether it must pass or **which rule it must fail by**.

Running the corpus after a spec change classifies every difference:

- A must-pass fixture that now fails is **over-constraining** (F4).
- A must-fail fixture that now passes is a **hole** (F3).
- A must-fail fixture that fails by a *different* rule means rule attribution moved, which is a silent behaviour change even though the count looks unchanged.

**The corpus is plain YAML with a manifest, never Go test files.**
That is deliberate. Downstream consumers of an openly licensed dataset can run it against their own tooling, a future TypeScript renderer can run the same fixtures, and the corpus outlives whichever implementation currently exists.
This is the model the JSON Schema Test Suite uses, and the reason it works is precisely that it belongs to no implementation.

### 2. Every example in the spec is executable

Each YAML block in `docs/spec.md` is extracted and run through the validator in CI.

Documentation examples rot silently and are the first thing a new contributor copies.
This is the cheapest mechanism here and it catches a disproportionate share of F1 and F7.

### 3. Spec coverage

The corpus is measured against the spec, not the code.

- Every enum value must appear in at least one fixture.
- Every rule must have at least one must-fail fixture.
- Every field must appear in at least one valid fixture.

Anything unexercised is untested spec surface and is reported, not silently tolerated.
An enum value no fixture ever uses is indistinguishable from one that does not work.

### 4. Property-based generation and mutation

A hand-written corpus only covers problems someone already imagined, which makes it structurally blind to F3.

Two generators close that:

- **Generate valid records from the spec.** The validator must accept every one, and parse to emit to parse must be a fixed point.
- **Mutate valid records into invalid ones.** Delete a required field, swap an enum value for an unregistered one, point a `ref` at nothing, introduce a cycle. The validator must reject, and should cite the rule that ought to have caught it.

A mutation that survives is a hole in the spec, and it is promoted into the permanent corpus so it can never reopen.

### 5. Breaking-change detection

On every spec change, the **previous** corpus runs against the **new** spec, and the delta is classified:

- **Widening**, such as a new optional field or a new enum value, is backward compatible.
- **Narrowing**, such as a new required field, a removed enum value, or a tightened pattern, is breaking and requires a `schema_version` bump.

CI fails a narrowing change that does not bump the version.
This directly targets F6, which matters more here than in a private schema because the dataset is published and consumed by tooling we do not control.

Protobuf was rejected in ADR 0002, but its tooling culture around `buf breaking` is the right idea and is worth taking even though the format is not.

### 6. A traceability graph, and an impact query

The spec is loaded as a graph: kinds own fields, fields have types and enums, rules read fields, examples and fixtures use fields and exercise rules.

That makes blast radius a query rather than a guess (F5):

```console
$ rig spec impact blocks_passage

field blocks_passage
  kinds        component, knot (role: stopper)
  rules        R001 threaded bounding, R014 loop-knot placement
  enums        none
  examples     4 in docs/spec.md
  fixtures     7 (2 must-pass, 5 must-fail)
  downstream   generated JSON Schema, pattern expander
```

The same graph finds orphans (F2) by inverting the question: which enums have no field, which rules read no field, which fields no rule constrains.

It is also what generates the rules table in the spec, so documentation drift (F7) becomes structurally impossible rather than a review responsibility.

## Consequences

Good:

- A spec change produces a diff of *behaviour*, not merely a diff of text.
- Over-constraining, the failure this project has actually hit twice, becomes a failing test rather than a discovery months later.
- Downstream consumers can verify their own implementations against the same corpus, which is a real obligation of publishing an open dataset.
- The rules table, the enum registry, and the enforced behaviour cannot disagree, because two of the three are generated from the first.

Bad:

- Real work before there is any data. The corpus, the generators, and the traceability tool are infrastructure for records that do not exist yet.
- The corpus must be maintained. A fixture set that stops tracking the spec is worse than none, because it manufactures confidence.
- Property-based generation over a spec this expressive is fiddly, particularly generating records that are *interestingly* valid rather than trivially so.
- Some failures remain out of reach. Nothing here catches a rule that is coherent, tested, and simply wrong about fishing. That is what `validation.status: field-tested` is for, and no amount of tooling substitutes for it.

## Considered Alternatives

**Test only the implementation.**
Standard practice, and insufficient. It answers whether the code matches the spec, never whether the spec is right. Both of this project's real schema bugs would have passed a complete implementation test suite.

**Rely on review.**
A thirty-rule spec exceeds what a human reliably holds in their head, and the whole premise of the project is not trusting that.

**Version the spec and never change it.**
Not credible for a spec with no data behind it yet. It will change substantially, which is the argument for making change *safe* rather than rare.

## More Information

- The spec: [`docs/spec.md`](../docs/spec.md).
- Enforcement stack: [ADR 0002](0002-schema-enforcement-stack.md).
- Rules need stable IDs (`R001` and so on) for the corpus to reference them. That numbering, currently informal in the spec's rules table, becomes normative as part of this decision.
