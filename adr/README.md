# Decision records

<!-- SPDX-License-Identifier: Apache-2.0 -->

[MADR](https://adr.github.io/madr/) records for riggermortis.

Every entry here exists because a decision had a real trade-off, and the value is in the **rejected** alternatives, not the choice.
A reversed decision gets a **new** record that supersedes the old one.
The old one stays, wrong, with its reasoning intact, because the reason it was wrong is the useful part.

| # | Decision | Status |
| --- | --- | --- |
| [0001](0001-encode-tying-method-not-topology.md) | Encode the tying method, not the knot's topology | accepted |
| [0002](0002-schema-enforcement-stack.md) | Enforce the spec with CUE, JSON Schema, and a Go validator | accepted |
| [0003](0003-test-the-spec-as-an-artifact.md) | Test the spec itself, not only its implementations | accepted |
| [0004](0004-polyglot-seam-go-for-spec-typescript-for-geometry.md) | Go for the spec, TypeScript for geometry, split at validated JSON | **superseded by [0007](0007-go-only-pregenerated-svg.md)** |
| [0005](0005-dual-licence-odbl-and-apache.md) | Dual licence: ODbL for the data, Apache-2.0 for everything else | accepted |
| [0006](0006-rigging-notation.md) | Invent the rigging notation, on borrowed architecture | accepted |
| [0007](0007-go-only-pregenerated-svg.md) | Go only, with SVGs pre-generated at build time | accepted |

## What each one settled

**[0001](0001-encode-tying-method-not-topology.md): tying method over topology.**
Knot theory is defined on closed loops, and a fishing knot is an open curve tied around an object. Knotoids are the honest treatment of an open curve and still break, because a knotoid has no second object. So the encoding is the *procedure*, after Suber, and topology is kept only as an optional validation fingerprint.

**[0002](0002-schema-enforcement-stack.md): CUE, JSON Schema, Go.**
Sorting the rules by what could possibly enforce each makes the answer obvious: roughly a third is declarative and the rest is graph traversal and cross-file joins. Rejected OpenAPI (describes APIs, this is a dataset) and Protocol Buffers (a wire format for machines, and its type system is weaker exactly where strength is needed). LinkML was genuinely close.

**[0003](0003-test-the-spec-as-an-artifact.md): test the spec.**
The failure that matters is **over-constraining**, and this project has hit it twice: "an ordered list along a line" forbade the umbrella rig, and "a knot is always an edge" forbade the dropper loop. Both were caught by hand-authoring an awkward specimen, which is not repeatable. A language-neutral conformance corpus and a traceability graph make it repeatable.

**[0004](0004-polyglot-seam-go-for-spec-typescript-for-geometry.md): superseded.**
Rested on "the browser is the target anyway." Once SVGs became pre-generated and the configurator was ruled out, that premise was false. Retained because the three.js reasoning inside it is a specific, instructive mistake.

**[0005](0005-dual-licence-odbl-and-apache.md): ODbL and Apache-2.0.**
ODbL is the only licence that cleanly separates a **Produced Work**, which needs attribution only, from a **Derivative Database**, which must be published. That distinction is why OpenStreetMap left CC BY-SA, and inheriting a known-broken arrangement would have been careless. The spec is Apache-2.0 on purpose, so the format stays freely implementable.

**[0006](0006-rigging-notation.md): the rigging notation.**
No notation exists for a rigid object piercing a soft body, in fishing or anywhere. Established, not assumed. So it is invented on borrowed architecture: grammar from knitout, geometry from needle-insertion planning, verb-set scale from surgical gesture taxonomies. Two verbs, `SK` and `BU`, have no prior art anywhere and are recorded as inventions.

**[0007](0007-go-only-pregenerated-svg.md): Go only.**
Pre-generated SVG plus no configurator means no geometry ever runs client-side, so the renderer is a build tool and the browser's ecosystem is irrelevant to choosing its language. Go is a build-time compiler here, not a server: records in, static files out, nothing running in production.

**[0010](0010-knot-shape-from-symmetry-not-alternation.md): shape from symmetry, not alternation.**
A reeve is two crossings, not one, and undercounting it is what forced the square knot to be hand-drawn. A bend whose half knots cancel is two interlocked bights, each collaring the other cord's *pair* rather than threading its eye, and the two cords are related by an inversion because the square knot is amphichiral. The rejected rule matters most: a dressed square knot is **not** an alternating diagram, and asserting that it was locked in the very bug the assertion was meant to catch.

## Writing one

Keep the MADR shape: context and problem statement, considered options, decision outcome with consequences both good and bad.

Number sequentially and never renumber.
State the bad consequences honestly, because an ADR that only lists upsides is marketing and will not be trusted the next time it matters.
