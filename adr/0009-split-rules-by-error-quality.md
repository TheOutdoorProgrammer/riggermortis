# 9. Split validation between CUE and Go by error quality

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

Research into CUE's capabilities ([`docs/research/06-cue-capabilities.md`](../docs/research/06-cue-capabilities.md)) established that CUE can do far more than [ADR 0002](0002-schema-enforcement-stack.md) assumed. The belief that CUE validates one file at a time was a property of the Makefile, not of CUE: the `-l` flag loads an entire dataset into one value keyed by id, and 616 records vet in half a second.

Nine of the eleven Go rules were shown to be expressible in CUE, including referential integrity, which was the largest single chunk of the Go program.

The question is how many of them should actually move.

## Decision Drivers

- The dataset is hand-authored, and it wants contributions from people who fish rather than people who write CUE. For them, **the error message is the validator's entire product.**
- [ADR 0003](0003-test-the-spec-as-an-artifact.md) states that a validator which silently passes is worse than none, because it manufactures confidence.
- Duplication is worse than either language. A rule expressed twice will drift.

## Decision Outcome

**A rule lives in CUE when its failure is structural. It lives in Go when its failure needs explaining.**

Moved to CUE, or deleted:

| Change | Effect |
| --- | --- |
| Whole-dataset loading via `-l` | The Makefile's per-directory `PAIRS` table is gone |
| `#Dispatch` on the `kind` field | A record picks its own schema; no definition named per file |
| `!` required markers | "field is required but not present" replaces "incomplete value string" |
| JSON Schema generation | Six published schemas, which [ADR 0002](0002-schema-enforcement-stack.md) promised and nothing had delivered |
| **R011 deleted** | Dead code. The schema already makes `n` non-optional |

Kept in Go: **R001, R002, R005, R007, R010, R023, R032, R037, R039, R040.**

### Why not move the rest

The report is clear that they work in CUE. It is equally clear about what it costs:

> Go produces `rig.carolina: threaded "bead" is stopped by "sinker", which does not block passage`.
> CUE produces `conflicting values false and true`.

That Go message is not decoration. It is how a real defect in the Carolina rig was found and understood: an egg sinker has `blocks_passage: false` because line runs through it, so it stops nothing. "Conflicting values false and true" would have been stared at and shrugged off.

Two further costs, both from the report:

- **Nil-guarding gets worse.** A missing field in CUE is *incomplete*, not an error, so a forgotten presence guard does not fail loudly. It degrades to `some instances are incomplete`, which points at nothing. Go's `if comp == nil { continue }` is uglier and far easier to get right.
- **`x & true` silently passes when `x` has a default.** Unification picks the matching branch instead of comparing. This cost the researcher two rounds to find, and it is a validator that silently passes, which is the precise failure ADR 0003 exists to prevent.

Referential integrity (R005) is the closest call, since it is the biggest chunk and CUE does it in about ten lines. It stays in Go because its message names the offending reference, and because the type-based CUE form emits one error line **per existing id**: nine lines at sixteen records, six hundred at six hundred.

### Consequences

Good:

- Contributors get messages that say what to fix.
- CUE is used for what it is unambiguously best at, and the wins there were substantial and are already taken.
- No rule is expressed twice, so nothing can drift.

Bad:

- The Go program stays larger than it needs to be, and roughly two hundred lines of it could be ten lines of CUE.
- Two languages hold validation logic, so a contributor must know where a given rule lives.
- If CUE's error messages improve materially, this is worth revisiting, and nothing here is hard to move later.

## Confirmed impossible, permanently

Graph traversal. `#Chain & {n: 3}` produces `structural cycle`, and the detector fires on the definition's shape rather than actual depth, so no decreasing counter rescues it. Any future rule that must walk a graph stays in Go by necessity rather than preference.

## More Information

- Full research, with every claim verified by running commands: [`docs/research/06-cue-capabilities.md`](../docs/research/06-cue-capabilities.md).
- The schema package had to split from the dataset overlay for JSON Schema generation to produce anything usable. With `records:` inside the schema package, `cue def --out jsonschema` emitted the dataset shape with zero definitions.
