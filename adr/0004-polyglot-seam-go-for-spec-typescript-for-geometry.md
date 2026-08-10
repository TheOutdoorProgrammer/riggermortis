# 4. Go for the spec, TypeScript for geometry, split at validated JSON

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

[ADR 0002](0002-schema-enforcement-stack.md) chose Go for the validator, on the basis that Go is the house language and produces a single static binary with no runtime dependencies.

That reasoning holds for spec tooling.
It does not obviously hold for the other half of the project, which turns a tying description into curves and emits SVG.
Choosing Go there because Go was already chosen would be deciding by momentum.

## Decision Drivers

- Spec tooling runs in CI on every pull request. Startup time and dependency footprint matter.
- The traceability graph and conformance runner from [ADR 0003](0003-test-the-spec-as-an-artifact.md) are graph algorithms and CLI ergonomics.
- The rendering half is not really an SVG problem. Emitting SVG is text formatting and any language does it. The hard part is **geometry**: fitting splines through control points, sweeping a tube along a 3D curve, resolving crossings so the right strand passes over, and projecting to 2D.
- The end product is a web page, so something must run in a browser regardless.

## Decision Outcome

**A polyglot split with one seam: validated JSON.**

| Layer | Language | Why |
| --- | --- | --- |
| Spec definition | CUE | Per ADR 0002. Generates JSON Schema. |
| Validation, conformance corpus, traceability, codegen | **Go** | One static binary, fast CI, graph work, no runtime deps. |
| Geometry, rendering, the site | **TypeScript** | The ecosystem for this problem, and the browser is the target anyway. |

**The seam is strict.** Go validates and emits a JSON bundle plus the generated JSON Schema. TypeScript consumes that bundle and never re-implements validation. Data crossing the seam has already been proven correct, so there is exactly one validator and no possibility of two implementations disagreeing.

### Why TypeScript for geometry specifically

Go's numerical and geometric ecosystem is thin. `gonum` is solid for linear algebra and there is no strong story above it for curves and tube meshes.

TypeScript has the pieces already built: a Catmull-Rom curve through control points with a tube swept along it is precisely how rope is modelled, and it is a library call rather than a project. The same code renders interactively in the browser and headlessly to static SVG, so the animation and the still diagram come from one implementation instead of two that drift.

Rendering is also where iteration speed dominates. Getting a knot to look right is a loop of adjusting and looking, and a browser with hot reload is a better instrument for that than a compile-and-write-a-file cycle.

### Consequences

Good:

- Each half uses the tool suited to it rather than the tool already present.
- The seam is a published artifact, so a third party could write their own renderer against the same validated JSON. That is a real benefit for an openly licensed dataset.
- The conformance corpus is already language-neutral by [ADR 0003](0003-test-the-spec-as-an-artifact.md), so the TypeScript side can run the identical fixtures.

Bad:

- Two toolchains, two dependency stories, two CI paths.
- A contributor may need both to work end to end, which raises the barrier.
- The JSON bundle is a real interface with real versioning obligations, not an internal call.
- It contradicts a general preference for single-binary Go tooling, and does so knowingly.

## Considered Alternatives

**All Go.**
One toolchain, one binary, matches the house preference. Rejected because the geometry ecosystem is not there and a browser renderer would still be needed, so the practical outcome is Go plus TypeScript anyway, with the geometry written twice.

**All TypeScript.**
Genuinely tempting. One language end to end, and Zod or TypeBox can generate JSON Schema from types. Rejected because it gives up the single static binary for the part that runs on every pull request, puts Node in the CI path for spec validation, and is a weaker fit for the traceability graph work. The spec tooling is where the durability matters most.

**Python for geometry.**
Strong numerics with numpy and scipy. Rejected because the output target is a browser, so it adds a third language without removing the need for the second.

## More Information

- If the seam proves to be friction rather than a boundary, the fallback is all-TypeScript, not all-Go. The geometry constraint is real and the Go preference is a preference.
