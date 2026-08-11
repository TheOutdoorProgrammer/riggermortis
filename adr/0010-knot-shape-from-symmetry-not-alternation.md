# 10. Derive a bend's shape from its symmetry, not from alternation

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-11

## Context and Problem Statement

The square knot was drawn from nineteen hand-authored waypoints in `internal/rope/bend.go`, plus a per-record table of crossings to invert (`squareFlips{"knot.square": {1,2,4}}`), plus a table of magic percentages deciding how much of the finished knot each stage showed.
The package doc directly above that table read "Nothing here is authored."

That is the thing the project exists not to do.
The README's whole claim is that every other project renders practical knots from hand-authored splines, and that generating them from a structured tying description is the point.

So: can the square knot be generated from `data/knots/square.yaml`?

The blocker was that the record's stages read as the braid word σ σ σ⁻¹ σ⁻¹, and the two-strand braid group is ℤ, so that word is the identity.
The solver was being asked to draw nothing, and somebody drew it by hand instead.

## Considered Options

1. **Keep the authored waypoints.** Honest about what it is, but it is the incumbents' approach with better typography.
2. **Relax a topologically correct start into shape.** Let physics do the dressing.
3. **Derive the dressed shape from the knot's own symmetry.**

## Decision Outcome

Chosen: **option 3**, with the crossing count fixed first.

**A reeve is two crossings, not one.**
Passing an end through a loop enters the region the loop bounds and then leaves it.
Counting it once gave the square knot four crossings; a dressed reef has six, and no four-crossing diagram can be a square knot at all.
That missing pair is the entire reason the picture had to be corrected by hand.

**A bend whose half knots cancel is two interlocked bights.**
Each half knot swaps which side the working ends are on, so an even number brings every end back beside its own standing part.
A cord that comes back to where it started has folded, and a folded cord is a bight.

**Each bight collars the other cord's paired legs. It does not thread the other bight's eye.**
This was got wrong once, with a constant named `clasp` deliberately dragging the eyes together to force the threading, and a confident comment explaining why that was necessary.
The result looked close and was a different knot.

**The two cords are related by an inversion, not a rotation.**
The square knot is amphichiral: it is its own mirror image, which is exactly what makes it a reef and not a granny.
So cord b is `(-x, -y, -z)` of cord a.
Using `(-x, -y, +z)`, a rotation, is the granny's symmetry, and it produces a knot that is structurally right and chirally wrong.

**Depth belongs to the section of cord, not to the crossings.**
The eye tucks under the whole of the other cord's pair; the legs ride over.
Rejected along the way: the rule that a practical knot diagram alternates over and under along each cord.

## Consequences

Good:

- No per-knot geometry anywhere. `bend.go`, `flip.go`, `squareFlips` and the stage percentage table are deleted.
- Stages are real tying states: stage *k* is a weave of the crossings made by stage *k*, not a crop of the finished knot.
- Every length is a multiple of the cord's own diameter. A knot has no scale of its own; it is as big as the rope it is tied in.
- `TestSquareKnotReading` locks the measured reading of a knot confirmed correct by someone who ties them.

Bad, and worth being plain about:

- **The bight is a parameterised primitive, not a general solver.** The shape is derived from the record, but *that a cancelling bend is two interlocked bights* is knowledge in the code, not knowledge in the notation. This generates the bend family. It does not generate arbitrary knots.
- **Relaxation is switched off for bights.** It flattened the eyes the clasp is made of. The dressing comes from the primitive, which is the opposite of the direction ADR-less commit `c19b2ca` was heading.
- **Granny knots fall through undrawn.** Half knots of one hand stand the knot on edge. There is no verified model for that here and no granny record to check one against, so such a record gets the helix rather than a fake reef.
- **Palomar is untouched and still wrong.** It renders as a twisted pair with no hook, no doubled line and no loop passed over the hook. The reeve fix moved it from two crossings to three. It was broken before this change and it is broken after.

## The rule that was wrong

**A dressed square knot is not an alternating diagram.**

Measured off the reference: cord a reads `U O U U O U`, cord b reads `O U O O U O`.
The doubled pass is where a cord's own eye goes under the whole of the other cord's pair. Collaring a pair means passing to one side of both strands, not weaving between them.

Alternation had been asserted twice in this repo and was never true:
`CheckAlternating` was added as a check and never called, and a test was later written that enforced it, which locked in the bug it was meant to catch.
Both are gone.

A related trap: `ReadCrossings` counted rendered *segments*, and two crossings a cord passes the same way merge into one segment, so it could not see a doubled `U U` and silently reported five crossings out of six.
It now counts real crossings via `rope.Reading`.
