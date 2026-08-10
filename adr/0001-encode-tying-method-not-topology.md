# 1. Encode the tying method, not the knot's topology

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

riggermortis renders knot-tying animations from structured data rather than from photographs or video.
That requires a machine-readable representation of a knot.
Two fundamentally different things could be represented: the **finished knot as an object** (its topology), or the **procedure that produces it** (its tying method).

The choice is architectural.
It determines the schema, the renderer, what can be validated automatically, and whether decades of existing mathematical tooling is an asset or a distraction.

The initial instinct was topology, because knot theory is rigorous, has libraries, and offers the appealing promise of machine-provable correctness.
Research found that instinct to be wrong in a way that is not obvious from the outside.

## Decision Drivers

- A fishing knot is tied **around an object** (a hook eye, a swivel, a second line).
- The user-facing artifact is an **animation of a procedure**, not a static depiction of a finished object.
- Correctness must be checkable in CI wherever that is possible.
- Friction, dressing, and tightening order are load-bearing properties of a working knot.

## Considered Options

1. **Classical knot theory** (Gauss codes, PD codes, invariants), the standard mathematical encoding.
2. **Knotoids** (Turaev 2012), the formalism built specifically for open curves.
3. **A tying-action notation**, encoding the procedure, after Peter Suber's Knot Tying Notation.
4. **Hand-authored splines per knot**, what every existing renderer does.

## Decision Outcome

Chosen: **option 3, a tying-action notation**, with topology retained only as an optional validation fingerprint.

Suber's notation separates the running part from the standing part, groups actions into **stages that map one-to-one onto animation keyframes**, and can wrap a **named external object**.
That last capability is the deciding factor: `MT(RP, Tree.1:CW:U):2` becomes `HookEye.1` and expresses the exact primitive that every fishing knot needs and that no topological formalism provides.

### Consequences

Good:

- The representation matches the deliverable. Stages are keyframes, so the animation falls out of the data rather than being authored alongside it.
- Tying around an object is expressible at all, which is the requirement that eliminated the alternatives.
- Left-handed and right-handed variants become a transform on the description rather than a mirrored bitmap.
- The data stays diffable and reviewable, so a change to a knot is legible in a pull request.

Bad:

- The mathematical ecosystem does not apply directly. Every knot library assumes a closed loop, so none of them can consume our primary representation.
- We take on the hard geometry step ourselves: turning an action sequence into renderable curves. Mathesis (Apache-2.0) is the best available prior art for this and is not a drop-in.
- Suber's notation carries a bare copyright notice with no license grant. The system is implementable regardless under the idea/expression distinction, but adopting the notation wholesale is a separate open question.
- Topological validation becomes a second, derived pipeline rather than a property of the primary encoding.

## Pros and Cons of the Options

### Classical knot theory

- Good: rigorous, mature libraries, exact unknot recognition via Regina.
- Bad: defined only on closed loops. A fishing knot is an open curve, which has no knot type at all. That is precisely *why* a badly chosen knot slips.
- Bad: expresses the finished object, not the procedure, so an animation cannot be derived from it.
- Bad: silent on friction, dressing, and tightening order.

### Knotoids

- Good: the mathematically honest treatment of an open curve, and `pyknotid`'s `OpenKnot` implements the closure machinery (MIT, pip-installable).
- Bad: a knotoid is a **single immersed interval with no second object**, so tying around a hook eye breaks the formalism outright.
- Bad: a bend joining two lines is a tangle, not a knotoid.
- Bad: `pyknotid` has been unmaintained since 2018, and random-closure methods return a distribution, which flakes in CI.

### Hand-authored splines

- Good: proven. It is what every recognizable knot renderer actually ships.
- Bad: it is the thing this project exists not to do. Authoring cost scales linearly with the catalogue, nothing is queryable, and a correction means re-authoring art.

## More Information

- Full research with citations: [`docs/research/01-knot-data.md`](../docs/research/01-knot-data.md) and [`docs/research/03-validation.md`](../docs/research/03-validation.md).
- [Peter Suber, Knot Tying Notation](https://legacy.earlham.edu/~peters/knotting/notate.htm).
- Topology's retained role, and the three-tier correctness model it feeds, is described in the README under Correctness posture.
