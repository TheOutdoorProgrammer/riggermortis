# riggermortis: agent context

An open, structured reference for every way a recreational angler can rig terminal tackle.
Rigs and knots are **data**; diagrams are **generated** from that data.

## Read before touching anything

| File | Why |
| --- | --- |
| [`docs/spec.md`](docs/spec.md) | The data spec. Kinds, enums, rules. Everything else is downstream of it. |
| [`adr/README.md`](adr/README.md) | Seven decisions and, more usefully, what was rejected and why. |
| [`docs/coverage.md`](docs/coverage.md) | What exists, what does not, what a human has actually confirmed. |

Do not propose an architectural change without checking whether an ADR already settled it.
If one did and you disagree, that is a **new** ADR superseding it, never an edit to the old one.

## The licence boundary is a directory

**`data/` is ODbL-1.0. Everything else is Apache-2.0.**

Putting a file in the wrong directory changes its licence.
Records go in `data/`. The spec, docs, ADRs and code go outside it.
Every file under `data/` carries `# SPDX-License-Identifier: ODbL-1.0` as its first line, because records get copied out of the repository and a file that loses its directory loses its terms.

## Non-negotiables

These are settled. Violating one is a bug, not a preference.

1. **Vendor neutral.** No brand, manufacturer, product line or model name in any normative field. Record the physical property, not the trade label: `rating_kg: 11`, never "size 3 split ring". The carve-outs are listed in the spec and are narrow.
2. **SI units.** `_mm`, `_m`, `_g`, `_kg`, `_c`. No bare numbers for physical quantities, ever, and the unit lives in the field name. Non-SI is permitted only for a fishing standard that is universal worldwide, and that register is currently empty.
3. **Enums by default.** If a field can be an enum it is one, registered in the Enumerations section. Free text is for names, prose and notes.
4. **Nothing claims to be true without saying how it was checked.** Every record carries a `validation` block recording the *method*. Never raise a record's `validation.status` for work you did yourself; `field-tested` means a human who fishes confirmed it.
5. **Copy nothing.** Rig topology is a method of operation and is not copyrightable, so learning facts from a copyrighted guide is fine. Reproducing its diagrams or prose is not. `animatedknots.com` is actively enforced; treat it as radioactive.

## Traps that have already bitten

**Storing the label instead of the property.** This has happened three separate times: a required `manufacturer` field, then imperial units, then `abok_ref` as a field-per-catalogue. Each time the fix was to record the underlying reality and let display convert. Assume the next one is out there.

**Asserting an identifier from memory.** An Ashley number, a scientific name, a breaking strength. These look like facts and are not recallable. Omit it, or cite it. `null` is not a placeholder; an absent entry says the same thing without occupying a field.

**Over-constraining the spec.** The two real schema bugs so far were both the spec forbidding something legitimate: "an ordered list of components along a line" forbade the umbrella rig, and "a knot is always an edge" forbade the dropper loop. When adding a constraint, look for the specimen that breaks it.

## Deliberately not built

- **No configurator.** Users do not compose rigs in the UI. It would be a machine for producing unvalidated rigs, and a contributed record is validated, reviewed and durable instead. See [ADR 0007](adr/0007-go-only-pregenerated-svg.md).
- **No server.** Go is a build-time compiler: records in, static files out. Nothing runs in production.
- **No client-side geometry.** SVGs are pre-generated in CI.

## Layout

```text
riggermortis/
├── adr/          # Decision records, indexed in adr/README.md
├── cue/          # Spec definitions, the source of truth for enums and shapes
├── data/         # ODbL. One record per file, SPDX header on each
├── docs/         # spec.md, coverage.md, research/
└── LICENSE       # Apache-2.0. data/LICENSE is ODbL-1.0
```

## Conventions

- **One record per file**, named after its id without the category prefix. `hook.offset-worm` lives in `data/components/hook-offset-worm.yaml`.
- **IDs** are `{category}.{name}`, exactly one dot, lowercase kebab-case, category from the `id_category` enum.
- **Commits** are one-line conventional subjects, no bodies.
- **Never push without being asked.** Committing as you go is expected; pushing is not.
- Run `yamllint data/` before committing records.

## Status

The spec is written and the data is not enforced by it yet.
There are records on disk and rules in a markdown table with nothing connecting them, so **every rule is currently aspirational**.
Building the CUE definitions and the validator is the active work, and the first useful output is finding out how many hand-authored records are wrong.
