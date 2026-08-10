# 6. Invent the rigging notation, on borrowed architecture

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

The `rigging` kind describes mounting a soft body on a hook, such as the Texas rig's Texposed method.
Its structure was settled by analogy with knots, a staged procedure where each stage is a keyframe.
Its **vocabulary** was not.

Knots got their vocabulary by adopting an existing system, [Suber's Knot Tying Notation](https://legacy.earlham.edu/~peters/knotting/notate.htm), on the principle that reinventing a considered system is strictly worse than adopting one.
The question was whether an equivalent exists for a rigid object piercing a soft body, or whether this one has to be invented.

Inventing a notation is the expensive answer, so it needed ruling out properly rather than assumed.

## What the research found

Full detail with citations in [`docs/research/05-rigging-notation.md`](../docs/research/05-rigging-notation.md).

**Nothing exists in fishing.** Searches across patents, fisheries literature, instructional material, competitive angling and open source found no formal or semi-formal notation, no structured format, and no attempt at a standard. Patent results are all physical baiting machinery rather than a language. Three near-misses fail identically: named rig types are labels for *end states* taught as numbered prose plus photographs; fisheries-science hooking-location categories describe where a hook ended up in a *caught fish* rather than how bait was mounted; fly-tying recipes are ingredient lists with prose assembly.

**No adjacent domain is a drop-in**, and they fail in an instructive way: each has half the problem and never joins the halves.

- **Surgery** has a mature named-pattern classification that is entirely prose and diagrams, *and* a rigorous geometric parameterisation that exists only because robots needed coordinates. The two never meet.
- **Embroidery** formats encode needle penetrations precisely but assume zero-thickness material, which is fatal when the whole question is what happens *inside* the body.
- **Robotics** has no action language for deformable objects at all.
- **Entomology** pinning has exactly the right conceptual object, a body-type to insertion-site lookup, published as photographs.

## Decision Outcome

**Invent the notation, assembled from three proven pieces rather than from nothing.**

### Architecture from knitout

[knitout](https://textiles-lab.github.io/knitout/knitout.html), from the CMU Textiles Lab, is the only system found that solved the whole meta-problem rather than a piece of it. Adopted:

- A **versioned header**, which `schema_version` already provides.
- An **equipment declaration** stating what the instructions assume, which is what `applies_to` does.
- A **flat instruction list with no flow control.** knitout is a compiled target; loops are expanded before emission. That matches this project's `patterns`, which expand before validation, so the emitted notation is always the flat form.
- An **`x-` extension namespace**, so a third party can extend without forking. Newly adopted here.
- **`VERB DIRECTION LOCATION SUBJECT` ordering** with face-plus-index addressing such as `f1` and `b-2`.

That addressing is the most directly transferable idea found, and it converges with a decision already made independently: normalising positions so one rigging record works across every bait length.

### Geometry from needle-insertion planning

Alterovitz's insertion plan, `X = (y₀, θ, b, d)` for location, angle, bevel roll and depth, maps onto bait rigging field for field.

This exposed a real gap. The provisional draft carried location, surface and depth but had **no angle and no roll**, and both matter: entry angle determines where the point exits, and the point's rotational orientation determines whether it tracks straight or wanders. Both are now fields.

### Verb-set scale from surgical gesture taxonomies

JIGSAWS defines 15 surgemes, and the 2026 SAGES Delphi taxonomy defines 24 gestures across 10 clusters. Together with Suber's 7 verbs, they establish that a closed set of roughly this size is the right granularity, and the final set of 7 verbs plus 2 descriptors sits comfortably inside it.

### The two verbs that are genuinely ours

`SK` skin, passing tangentially just beneath a surface, and `BU` bury, terminating subsurface without exiting, **have no prior art in any domain examined.** Surgery names only one of them, and only because subcuticular closure happens to require it.

They are recorded as inventions rather than borrowings, because that is the honest measure of the gap this notation fills.

### Consequences

Good:

- Every structural decision has a precedent that has already survived contact with real users, so the novelty is confined to the two verbs that genuinely had no answer.
- Adopting knitout's `x-` namespace means third parties extend rather than fork.
- The angle and roll fields were a real hole and are closed before any data was authored.
- Position addressing is normalised, so a rigging record is bait-length independent.

Bad:

- It is still an invented notation, so it carries the full burden of one: nobody knows it, no tooling exists, and it will be wrong in ways only authoring will reveal.
- Borrowing across four domains risks incoherence. Mitigated by the verb set staying small and the grammar staying knitout's.
- The two invented verbs have no external validation whatsoever.

## Licensing

Nothing blocks adoption, because in every case the useful thing is a structure or a set of names rather than copyrightable expression.

Verified: knitout's repository carries **no LICENSE file and no copyright statement in the spec**. The architecture may be borrowed freely; its prose may not be copied. The lab lists `knitout-feedback@cs.cmu.edu` for contact, and reaching out is courteous and worth doing.

Suber's page likewise carries no licence statement, which is the same position already recorded for the knot vocabulary.

## Considered Alternatives

**Extend Suber's notation to cover rigging.**
Tempting for consistency, and rejected because the nouns are wrong. Suber's vocabulary assumes cord: running part, standing part, bight, loop. A soft plastic is a solid body with anatomy, and you cannot make a bight in a worm. Sharing the *stage architecture* while keeping separate vocabularies is the correct amount of reuse.

**Adopt a surgical suture notation.**
Rejected because the half that is standardised is prose, and the half that is formal exists only inside robotic surgery research as coordinates for a specific manipulator. Neither half is a notation an author could write.

**Use free-text prose and skip notation entirely.**
Rejected because it forfeits the entire premise. Prose cannot be animated, validated, diffed, or checked for the missing alignment step that makes a Texas rig spin.
