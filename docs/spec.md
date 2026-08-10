# riggermortis spec

<!-- SPDX-License-Identifier: Apache-2.0 -->
> **Licence.** This document is **Apache-2.0**, so anyone may implement the spec freely.
> The dataset it describes, under [`data/`](../data/), is **ODbL-1.0**.
> See [LICENSE](../LICENSE), [LICENSE-DATA](../LICENSE-DATA), and [ADR 0005](../adr/0005-dual-licence-odbl-and-apache.md).

Exploratory shapes, not settled.
The purpose of this document is to find where the model breaks before any generator exists.
Everything here is written to be argued with.

Every field is documented, every kind carries a definition of what it *is* and what does **not** belong in it, and every closed vocabulary is an enumeration registered in one place.

Three principles run through the whole document:

- **The spec is vendor neutral.** Line is line, hooks are hooks. What a thing is made of and how it behaves is normative; who sold it is not.
- **Free text is the exception, not the default.** If a field can be an enum, it is one. Free strings are reserved for names, prose, and notes.
- **Nothing claims to be true without saying how it was checked.** Every record carries a `validation` block recording the method, not merely a boolean.

Fields that exist only so a validator can reason are marked **validator-only** and must never be rendered to a reader.

## The kinds

| Kind | One line | Lives in | ID category |
| --- | --- | --- | --- |
| `component` | A physical thing you can hold, buy, and attach | `data/components/` | many, see below |
| `line` | A cord material, and its diameter per stated strength | `data/lines/` | `line` |
| `knot` | A procedure for tying, expressed as stages | `data/knots/` | `knot` |
| `rig` | A graph of components, lines, and knots | `data/rigs/` | `rig` |
| `pattern` | A named, parameterised fragment instantiated many times | `data/patterns/` | `pattern` |
| `technique` | How a rig is fished, and the conditions it suits | `data/techniques/` | `technique` |
| `species` | A target fish, with its regional names reconciled | `data/species/` | `species` |
| `source` | A citation, with its copyright position recorded | `data/sources/` | `src` |

## Conventions

### `kind` and `id`

Every record carries both, and they are not redundant.

- **`kind`** is the record type. It is the schema discriminator and tells a validator which rules apply.
- **`id`** is a globally unique, stable, human-legible handle of the form `{category}.{name}`.

**The grammar is strict and machine-checked:**

```text
id      ::= category "." name
category::= a member of the id_category enum   (see Enumerations)
name    ::= [a-z0-9]+ ("-" [a-z0-9]+)*         lowercase kebab-case
```

Regex: `^[a-z][a-z0-9]*\.[a-z0-9]+(-[a-z0-9]+)*$`

**Exactly one dot.** The category is always a registered enum member; the name never contains a dot.
An earlier draft used `technique.carolina.drag`, which has two and is now `technique.carolina-drag`.
Underscores are not permitted anywhere in an ID.

The namespace is deliberate.
References appear far from the definitions they point at, and `ref: weight.egg` is self-documenting at the call site in a way that `ref: egg` is not.
In hand-authored data that readability is worth more than the purity of an opaque identifier.

There is no separate `category` field.
An earlier draft carried both `id: weight.egg` and `category: weight`, which is exactly the duplication that drifts.
The ID prefix **is** the category, and it must be a member of the `id_category` enum.

The cost is real: recategorizing something means renaming it, which breaks every reference.
Two mitigations, both enforced:

- `former_ids` records any previous ID, so references migrate mechanically and old links redirect.
- An ID, once published, is never reused for a different thing.

`kind` duplicating the directory is intentional.
Records get bundled, exported as one JSON blob, and passed around individually, all of which lose directory context.
A self-describing record survives that.

### Vendor neutrality

**No brand, manufacturer, product line, or model name appears in any normative field.**

A rig calls for a circle hook with an eleven millimetre gap, braid breaking near forty-five newtons, and a fourteen gram egg sinker.
It never calls for a named product.

The reason is not squeamishness about commerce, it is that a trade label carries no information.
Research on this domain established that size labels are unstandardised across manufacturers *and* across patterns within one manufacturer: a 2/0 octopus hook is not a 2/0 worm hook, and one maker's size 3 split ring is rated near 25 lb where its own 3H is near 65 lb.

An earlier draft responded to that by making `manufacturer` a **required** field, so a label could at least be disambiguated.
That was working around the problem instead of removing it.
The brand was only ever a proxy for a physical property, so the spec records the **property itself** and the proxy becomes unnecessary:

| Instead of | Record |
| --- | --- |
| size 3 split ring | `rating_n: 111` |
| 2/0 octopus hook | `gap_mm: 11`, `point_style: turned-in`, `shank: short` |
| 10 lb fluorocarbon | `material: fluoro`, `breaking_load_n: 44.5`, `diameter_mm: 0.285` |
| 1/2 oz egg sinker | `mass_g: 14.2`, `mounting: threaded` |

Two consequences fall out, both good.
Comparability stops being a hazard, because nothing compares labels any more.
And the dataset works anywhere, since an angler holding entirely different brands still has a circle hook with a measurable gap.

**Enforcement is structural rather than a lint rule.** There is no field to put a brand in. With `additionalProperties: false` set everywhere, a record that tries is rejected by the validator rather than caught in review.

Three carve-outs, stated explicitly so they do not read as leaks:

1. **Sources name their authors.** A citation that concealed its publisher would not be a citation. `reliability` may even turn on it, since a maker documenting their own product is a strong source.
2. **Aliases may record a trademarked common name.** Recording that some people say "Alabama Rig" is how a search for that term reaches the generic umbrella rig entry. It is never the canonical name, and it always carries a `trademark_note`. Naming a mark in order to avoid using it is the opposite of vendor specificity.
3. **`validation_method: manufacturer-confirmed`** describes where evidence came from, not what to buy.
4. **A knot's established name is its name, even when a brand is inside it.** The Berkley braid knot and the Rapala knot have no generic equivalents in use, exactly as the Albright special carries a person's name. Renaming them would make the spec less usable and no more neutral. The line is that a *name* identifies a procedure, while a *specification* prescribes a purchase, and only the second is forbidden.

Free-text fields such as `name`, `notes`, `prose` and `cadence` can still smuggle a brand in.
Nothing mechanical catches that, so it is a review responsibility and rule 22 states it.

### Units

**No physical quantity is ever a bare number.**
The unit is part of the field name, always, so a value can never be misread and never needs parsing.

**The spec is SI only.**

| Quantity | Suffix | Unit |
| --- | --- | --- |
| Small dimensions | `_mm` | millimetre |
| Distance and depth | `_m` | metre |
| Mass | `_g` | gram |
| Force and breaking load | `_n` | newton |
| Temperature | `_c` | degree Celsius |
| Time | `_s` | second |
| Angle | `_deg` | degree |

"SI only" here means SI base and derived units, their prefixed forms, and the units the SI Brochure accepts for use with SI, which is what makes degree Celsius and the plane degree legal.

An earlier draft stored pounds, ounces, inches and Fahrenheit on the grounds that this is how the domain writes things.
That was the same mistake as recording a manufacturer: it stored the **label** instead of the **property**.
How anglers write a number is a display concern, and the display layer converts freely.
A reader in Ohio sees "10 lb braid" and a reader in Queensland sees "4.5 kg", from one stored value that is neither.

### Breaking load is stored in newtons

This is the one judgement worth arguing with.

Line "test" and hardware "rating" describe a **force**, and the SI unit of force is the newton.
The international fishing convention states line classes in kilograms, which is metric but is a mass standing in for a force, and adopting it would reintroduce exactly the proxy the vendor-neutrality rule removed.

So `breaking_load_n: 44.5` rather than `test_lb: 10` or `test_kg: 4.5`.
Both familiar labels are then display conversions, and neither is privileged.

The honest cost is authoring ergonomics: nobody has a feel for a knot holding 44.5 newtons.
That is a tooling problem rather than a spec problem, so the authoring CLI accepts a value in any unit and normalises on write, and CI range-checks the result for plausibility.
Hand-converting is where the errors would come from, and no author should be doing it.

Anything a person would state as a range is stored as a range, `[min, max]`, never flattened to a midpoint.

Anything a person would state as a range is stored as a range, `[min, max]`, never flattened to a midpoint.

### `schema_version`

Every record carries one.
The dataset is published under a reciprocal licence, so downstream consumers need to know what shape they are consuming and when it changed.

### Tier, rank, and validation

Three orthogonal things, routinely confused.

- **`tier`** is *how knowable* a claim is at all.
- **`rank`** is *how much we currently believe it*, after Wikidata. A claim later found wrong is **demoted with a reason, never deleted**, so the record shows its own history and nobody re-adds it.
- **`validation`** is *how it was actually checked*. This is the only mechanism that closes the Tier C gap, and a boolean is not good enough: reading three websites and physically tying the knot are not the same evidence.

```yaml
validation:
  status: unvalidated
  events:
    - method: source-corroborated
      result: pass
      by: joeystout
      date: 2026-08-10
      detail: Three independent sources agree on component order.
      sources: [src.example-a, src.example-b, src.example-c]
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `status` | `validation_status` | yes | The highest rung reached. Derived from events, but stored so it is greppable and diffable. |
| `events` | list | yes | Append-only history. May be empty; may never be rewritten. |
| `events[].method` | `validation_method` | yes | **How** it was checked. This is the field that carries the actual weight. |
| `events[].result` | `validation_result` | yes | What happened. A `fail` event is kept, not deleted. |
| `events[].by` | string | yes | Handle of the person or `ci` for automation. |
| `events[].date` | date | yes | ISO date. |
| `events[].detail` | string | no | What specifically was checked. |
| `events[].sources` | list of refs | conditional | Required when the method is `source-corroborated`. |

The ladder runs `unvalidated` → `machine-only` → `corroborated` → `field-tested`, with `disputed` as a side state that any rung can fall into.
`field-tested` is the top because someone tying the rig and fishing it is the only evidence that closes the gap between our written instructions and reality.

---

## Enumerations

Every closed vocabulary in the schema, in one place.
This section is the source of truth from which the JSON Schema is generated, so adding a value here is the only way to make it legal anywhere.

### Structural

| Enum | Values |
| --- | --- |
| `kind` | `component`, `line`, `knot`, `rig`, `pattern`, `technique`, `species`, `source` |
| `id_category` | **components:** `hook`, `weight`, `swivel`, `snap`, `split-ring`, `bead`, `float`, `stop`, `blade`, `skirt`, `jighead`, `bait`, `lure`, `hardware`, `sleeve` · **others:** `line`, `knot`, `rig`, `pattern`, `technique`, `species`, `src` |
| `pattern_target` | `rig`, `knot` |
| `pin_type` | `closed-eye`, `open-eye`, `split-ring`, `shank`, `snap`, `arm-socket`, `line-end`, `in`, `out`, `loop` |
| `severity` | `error`, `warning`, `info` |

### Evidence

| Enum | Values |
| --- | --- |
| `tier` | `A` machine-provable · `B` heuristically checkable · `C` citation or human only |
| `rank` | `preferred`, `normal`, `deprecated` |
| `validation_status` | `unvalidated`, `machine-only`, `corroborated`, `field-tested`, `disputed` |
| `validation_method` | `ci`, `source-corroborated`, `manufacturer-confirmed`, `photo-compared`, `tied`, `fished`, `expert-review` |
| `validation_result` | `pass`, `fail`, `partial`, `disputed` |

`tied` means someone physically made it and compared it to the record.
`fished` means someone used it and it behaved as described.
`photo-compared` means the record was checked against a reference image rather than a physical object, which is weaker than `tied` and stronger than `source-corroborated`.

### Components and line

| Enum | Values |
| --- | --- |
| `mounting` | `tied`, `threaded` |
| `unit` | `mm`, `m`, `g`, `n`, `c`, `s`, `deg` |
| `variant_axis` | `mass_g`, `gap_mm`, `diameter_mm`, `rating_n`, `breaking_load_n`, `length_mm`, `supports_g` |
| `point_style` | `straight`, `turned-in`, `turned-out`, `knife-edge`, `needle` |
| `shank` | `short`, `standard`, `long`, `extra-long` |
| `bend_style` | `round`, `octopus`, `circle`, `worm`, `ewg`, `kahle`, `aberdeen`, `siwash`, `treble` |
| `line_material` | `mono`, `fluoro`, `braid`, `wire`, `backing`, `leader-mono`, `leader-fluoro` |
| `stretch` | `low`, `medium`, `high` |
| `density` | `sinking`, `neutral`, `floating` |
| `abrasion_resistance` | `low`, `medium`, `high` |
| `underwater_visibility` | `low`, `medium`, `high` |
| `knot_sensitivity` | `low`, `medium`, `high` |

### Knots

| Enum | Values |
| --- | --- |
| `knot_role` | `terminal`, `bend`, `loop`, `stopper`, `arbor` |
| `verb` | `GP`, `MB`, `ML`, `MT`, `MV`, `RV`, `TW` |
| `direction` | `F`, `A`, `L`, `R`, `U`, `D` |
| `reeve_direction` | `F-A`, `A-F`, `L-R`, `R-L`, `U-D`, `D-U` |
| `rotation` | `CW`, `CCW` |
| `chirality` | `/`, `\` |
| `force` | `push`, `pull` |
| `plane` | `HP`, `VP`, `EP` |

### Rigs

| Enum | Values |
| --- | --- |
| `edge_rel` | `tied`, `threaded`, `fixed`, `clipped`, `crimped`, `continuous`, `looped` |
| `node_role` | `main-line`, `leader`, `dropper`, `tag` |
| `travel_stop` | a node id, or the literal `open` |

### Techniques and conditions

| Enum | Values |
| --- | --- |
| `motion_primary` | `drag`, `hop`, `shake`, `deadstick`, `yoyo`, `swim`, `rip`, `walk`, `troll`, `vertical`, `drift`, `twitch` |
| `contact` | `bottom`, `suspended`, `surface` |
| `water_column` | `bottom`, `lower`, `middle`, `upper`, `surface` |
| `clarity` | `muddy`, `stained`, `clear`, `gin-clear` |
| `cover` | `none`, `sand`, `gravel`, `rock`, `rubble`, `sparse-grass`, `heavy-grass`, `wood`, `laydown`, `dock`, `bridge`, `reef`, `oyster`, `mangrove`, `current-seam` |
| `season` | `prespawn`, `spawn`, `postspawn`, `summer`, `fall`, `winter` |

### Species and sources

| Enum | Values |
| --- | --- |
| `water` | `freshwater`, `saltwater`, `both` |
| `source_type` | `book`, `article`, `agency`, `extension`, `manufacturer`, `video`, `expert`, `forum` |
| `copy_policy` | `cite-only`, `quotable`, `adaptable` |
| `reliability` | `low`, `medium`, `high` |

`license` is deliberately **not** an enum. It is an SPDX identifier, or the literals `all-rights-reserved` or `unknown`.

---

## The primitive that holds it together

Knots and rigs share one vocabulary: the **typed connection point**, or pin.

A component exposes pins with types.
A knot declares which pin types it may terminate on.
A rig declares which pin a given knot ties to.

Because both speak the same vocabulary, a rig is a netlist and validation is electrical rule checking.
Design the knot schema and the rig schema independently and they will not compose.

---

## `component`

**What it is.** A discrete physical object an angler can buy, hold, and attach to a rig. It has no behaviour of its own; it is a part.

**What belongs here.** Hooks, weights, swivels, snaps, split rings, beads, floats, manufactured bobber stops, blades, skirts, umbrella heads, wire arms, jigheads, soft plastics, and hard lures.

**What does not.** Line, which has its own kind because its properties vary continuously with stated strength. Knots, which are procedures. Anything that exists only as an assembly.

**Boundary rule for lures.** If you buy it assembled, it is a **component**. If you assemble it, it is a **rig**. A spinnerbait is internally a wire, two blades, a skirt and a hook, but it arrives finished, so it is one component.

```yaml
kind: component
schema_version: 0
id: swivel.barrel
name: Barrel swivel
mounting: tied
pins:
  - {id: eye-a, type: closed-eye, wire_diameter_mm: 0.9}
  - {id: eye-b, type: closed-eye, wire_diameter_mm: 0.9}
blocks_passage: true
variants:
  axis: rating_n
  values: [111, 156, 245, 356, 578]
validation:
  status: unvalidated
  events: []
```

```yaml
kind: component
schema_version: 0
id: hook.circle
name: Circle hook
mounting: tied
bend_style: circle
point_style: turned-in
shank: short
pins:
  - {id: eye, type: closed-eye, wire_diameter_mm: 1.1}
blocks_passage: true
variants:
  axis: gap_mm
  values: [6, 8, 11, 14, 18, 23, 30]
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `mounting` | `mounting` | yes | `tied` if it connects through a pin, `threaded` if it rides on a line segment instead. This single field separates a swivel from a sinker and determines which edge types are legal. |
| `pins[].id` | string | yes | Local identifier, referenced by rig edges. |
| `pins[].type` | `pin_type` | yes | A knot must declare this type in `connects.to` to be legal here. |
| `pins[].wire_diameter_mm` | number | no | Physical gauge of the eye. Feeds knot-fit warnings. |
| `blocks_passage` | bool | yes | Whether something threaded on the line is stopped here. **Validator-only.** Makes the sliding-weight check decidable. |
| `bore_mm` | number | threaded only | Largest thing that can pass through. **Validator-only.** |
| `bend_style` | `bend_style` | hooks only | The hook's geometry family. This, not a trade name, is what distinguishes a circle from an EWG. |
| `point_style` | `point_style` | hooks only | Where the point aims. Determines hooking behaviour. |
| `shank` | `shank` | hooks only | Relative shank length. |
| `variants.axis` | `variant_axis` | yes | The **physical quantity** along which this component varies. Never a trade label. |
| `variants.values` | list of numbers | yes | The values the component is commonly available in, in the axis's units. |
| `former_ids` | list | no | Previous IDs, for reference migration. |

A component record describes a **type**, and a rig instantiates it with a value from that type's axis: `{ref: hook.circle, gap_mm: 11}`.
This is the same idea as pattern parameters, applied one level down.

Mapping a shop label such as "size 4/0" onto a gap in millimetres is a genuinely useful thing, and it is deliberately **not** in scope here.
It is a shopping aid rather than a fact about a rig, it changes per manufacturer, and the charts that carry it are copyrighted catalogue content.
If it is ever built it is a separate, clearly non-normative dataset with its own licensing position, and nothing in the spec depends on it.

---

## `line`

**What it is.** A cord material together with the table mapping stated strength to actual diameter.

**Why it is not a component.** A component has one identity. Line is a continuum, and rule 3 cannot run without a real diameter to compare against a bore.

```yaml
kind: line
schema_version: 0
id: line.fluorocarbon
name: Fluorocarbon
material: fluoro
properties:
  stretch: low
  density: sinking
  abrasion_resistance: high
  underwater_visibility: low
  knot_sensitivity: high
diameters:
  - {breaking_load_n: 44.5, diameter_mm: 0.285, source: src.example}
  - {breaking_load_n: 89.0, diameter_mm: 0.405, source: src.example}
validation: {status: unvalidated, events: []}
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `material` | `line_material` | yes | The material family. |
| `properties.knot_sensitivity` | `knot_sensitivity` | yes | How much the material punishes a poorly chosen knot. Fluorocarbon and braid are both `high` for different reasons, which is what makes rule 16 worth having. |
| `diameters[].breaking_load_n` | number | yes | Breaking load in newtons. The familiar "10 lb test" is a display conversion of this. |
| `diameters[].diameter_mm` | number | yes | Actual measured diameter. Feeds rule 3. |
| `diameters[].source` | ref | yes | Required. Diameter per stated test is **not** standardised, so every pairing is a sourced measurement rather than a fact about the material. Where sources disagree, both rows exist and `rank` separates them. |

---

## `knot`

**What it is.** A **procedure**, not a shape. A sequence of stages, each of which is one animation keyframe.

**What belongs here.** Anything tied in cord: terminal knots, bends, loops, stoppers, arbor knots.

**What does not.** Manufactured connectors that do a knot's job, such as a crimp sleeve or a snap. Those are components.

### On Suber's notation

The vocabulary is taken from [Peter Suber's Knot Tying Notation](https://legacy.earlham.edu/~peters/knotting/notate.htm) (2002 to 2004).
It is not a widely adopted standard, so this is adoption on merit rather than for interoperability.

**Adopted:** the verb set, the noun vocabulary, **stages** in which all actions happen simultaneously, the direction system, **loop chirality** `/` versus `\` which fixes which strand passes over and is invariant under rotation, the **push-pull operator**, and the **external object extension**, which Suber explicitly intended for fishhooks among other hardware.

**Not adopted:** the string syntax itself, since structured actions diff readably in a pull request and validate without a parser; the `®` until-operator, since a runtime condition cannot render deterministically; and `Do(*n)` stage iteration, since repetition attaches to the action instead.

**Extended for fishing:** `repeat` accepts a **range**, because real instructions say "five to seven turns" and collapsing that to one number is false precision.

Suber's own caveats are recorded honestly: the notation is not adequate for decorative knots, webbing, splicing, or whippings, and he deliberately made it incomplete to keep it learnable.
None of those gaps touch fishing knots.

Note which layer each verb serves.
`GP` grips with a hand, generating **prose** but no geometry.
`ML`, `MT`, `RV` and `MB` manipulate cord, generating **geometry**.

### `role` is what makes a knot placeable

This is the fix for the dropper loop.

| Role | Joins | Appears in a rig as |
| --- | --- | --- |
| `terminal` | line end to a hardware pin | an **edge** |
| `bend` | two line ends | an **edge** |
| `loop` | tied *within* one line, producing a branch | a **node** with pins `in`, `out`, `loop` |
| `stopper` | tied *on* one line, blocking passage | a **node** with pins `in`, `out` |
| `arbor` | line to a reel spool | an **edge**, terminal at the reel |

A `loop` knot is structurally a three-pin component made of line.
That unification is what lets a high-low rig and a Sabiki exist at all, and it means `blocks_passage` applies to knots-as-nodes exactly as it does to components.

```yaml
kind: knot
schema_version: 0
id: knot.palomar
role: terminal
names:
  canonical: Palomar knot
  aliases: [Palomar]
abok_ref: null
connects:
  from: line-end
  to: [closed-eye, split-ring]
line_types: [mono, fluoro, braid]
objects:
  - {ref: HookEye.1, pin_type: closed-eye}
stages:
  - id: 1
    actions:
      - {verb: MB, subject: RP, names: BT.1, length_mm: 150}
    prose: Double roughly six inches of line back on itself to form a bight.
    notation: "* MB(RP=BT.1), LG(BT.1:150mm)"
  - id: 2
    actions:
      - {verb: RV, subject: BT.1, through: HookEye.1, direction: F-A, force: push}
    prose: Pass the bight through the hook eye.
    notation: "* RV(^BT.1, HookEye.1:F-A)"
  - id: 3
    actions:
      - {verb: ML, subject: [BT.1, SP], chirality: "/", names: LP.1}
      - {verb: RV, subject: E.BT.1, through: LP.1, direction: U-D}
    prose: Tie a loose overhand knot with the doubled line, leaving the hook hanging.
    notation: "* ML(/[BT.1, SP]=LP.1), RV(E.BT.1, LP.1:U-D)"
  - id: 4
    actions:
      - {verb: MV, subject: BT.1, over: Hook.1, direction: D, force: pull}
    prose: Pass the bight completely over the hook.
    notation: "* MV(BT.1^, Hook.1:D)"
  - id: 5
    actions:
      - {verb: MV, subject: [SP, E.RP], force: pull, wet: true}
    prose: Moisten the line and draw both the standing line and the tag end to seat the knot.
    notation: "* MV([SP, E.RP]^)"
strength:
  claims:
    - {residual_pct: 85, line: braid, n: 1, source: src.example, tier: C, rank: normal}
validation: {status: unvalidated, events: []}
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `role` | `knot_role` | yes | Determines whether the knot is an edge or a node. |
| `abok_ref` | int | no | Ashley Book of Knots reference number. The **number** may be cited; text and figures may not be reproduced, as ABoK is in copyright until 2040. |
| `connects.to` | list of `pin_type` | conditional | Absent for `loop` and `stopper` roles, which terminate on nothing. |
| `line_types` | list of `line_material` | yes | Mismatch is a warning, not an error, because it is a judgement. |
| `objects` | list | no | Suber user-extension nouns bound to the pin type they stand for. `HookEye.1` is what lets one knot record be reused against any component exposing a `closed-eye`. |
| `stages[].actions` | list | yes | Performed **simultaneously** within the stage, per Suber's semantics. |
| `actions[].verb` | `verb` | yes | |
| `actions[].subject` | string or list | yes | A list is a cluster acted on as one object. |
| `actions[].names` | string | no | Binds a name to the result, Suber's `=`. |
| `actions[].direction` | `direction` or `reeve_direction` | no | |
| `actions[].rotation` | `rotation` | no | |
| `actions[].chirality` | `chirality` | ML only | Fixes which strand passes over. |
| `actions[].force` | `force` | no | Supplies the animation its motion. |
| `actions[].repeat` | int or range | no | `[5, 7]` is legal and preferred. |
| `actions[].wet` | bool | no | Structured rather than buried in prose, because it materially affects the result. |
| `stages[].prose` | string | yes | Human instruction. Authored, see open questions. |
| `stages[].notation` | string | derived | Emitted from the structured actions. Never hand-edited; CI regenerates and diffs it. |
| `strength.claims[].n` | int | yes | Sample size. **Required.** Most published knot tests are n=1. |
| `strength.claims[].tier` | `tier` | yes | Always `C`. |
| `strength.claims[].rank` | `rank` | yes | |

---

## `rig`

**What it is.** A **graph** of components, line segments and knots forming a complete terminal presentation.

**What belongs here.** Anything an angler assembles.

**What does not.** How you fish it, which is a `technique`. Conditions, which attach to the rig-and-technique pairing.

**The rod boundary.** A node with `role: main-line` has a rod-ward end that is **not free**, because it continues to the reel. Only the terminal-ward end is free. That asymmetry is what makes rule 1 correct rather than merely plausible.

### Carolina rig, a chain

```yaml
kind: rig
schema_version: 0
id: rig.carolina
names: {canonical: Carolina rig, aliases: [C-rig]}
nodes:
  - {id: main,   type: line, role: main-line, material: braid, breaking_load_n: 133}
  - {id: sinker, ref: weight.egg,    mass_g: 14.2}
  - {id: bead,   ref: bead.glass,    diameter_mm: 8}
  - {id: swivel, ref: swivel.barrel, rating_n: 245}
  - {id: leader, type: line, role: leader, material: fluoro, breaking_load_n: 67, length_mm: [305, 1220]}
  - {id: hook,   ref: hook.ewg,      gap_mm: 14}
edges:
  - {from: main,   to: sinker, rel: threaded, travel: {toward_rod: open, toward_terminal: swivel}}
  - {from: main,   to: bead,   rel: threaded, travel: {toward_rod: open, toward_terminal: swivel}}
  - {from: main,   to: swivel, rel: tied, knot: knot.palomar, pin: eye-a}
  - {from: swivel, to: leader, rel: tied, knot: knot.palomar, pin: eye-b}
  - {from: leader, to: hook,   rel: tied, knot: knot.palomar, pin: eye}
validation: {status: unvalidated, events: []}
```

`travel.toward_rod: open` is correct and intentional.
The sinker is supposed to slide freely up the main line; that is the entire point.
It is the **terminal** side that must be stopped, or the weight leaves.

### Umbrella rig, a tree

The specimen that killed "an ordered list of components along a line."

```yaml
kind: rig
schema_version: 0
id: rig.umbrella
names:
  canonical: Umbrella rig
  aliases: [Alabama Rig, A-Rig]
  trademark_note: >
    "Alabama Rig" and "A-Rig" are registered marks of Slick Lures, LLC.
    The canonical name stays generic; the marks are aliases only.
legality:
  general_warning: true
  notes: >
    Hook-count limits apply in many jurisdictions and some prohibit
    multi-lure rigs outright. Check local regulations before use.
nodes:
  - {id: main, type: line, role: main-line}
  - {id: head, ref: hardware.umbrella-head, arms: 5}
  - {id: arm1, ref: hardware.wire-arm}
  - {id: jig1, ref: jighead.swimbait}
edges:
  - {from: main, to: head, rel: tied,    knot: knot.palomar, pin: eye}
  - {from: head, to: arm1, rel: fixed,   pin: arm-socket-1}
  - {from: arm1, to: jig1, rel: clipped, pin: snap-1}
```

### High-low rig, a knot as a node

The specimen that killed "a knot is always an edge."

```yaml
kind: rig
schema_version: 0
id: rig.high-low
names:
  canonical: High-low rig
  aliases: [chicken rig, double drop, two-hook bottom rig]
nodes:
  - {id: main,   type: line, role: main-line}
  - {id: drop1,  ref: knot.dropper-loop}
  - {id: seg1,   type: line, role: leader}
  - {id: drop2,  ref: knot.dropper-loop}
  - {id: seg2,   type: line, role: leader}
  - {id: hook1,  ref: hook.baitholder}
  - {id: hook2,  ref: hook.baitholder}
  - {id: sinker, ref: weight.bank}
edges:
  - {from: main,  to: drop1,  rel: continuous, pin: in}
  - {from: drop1, to: hook1,  rel: looped,     pin: loop}
  - {from: drop1, to: seg1,   rel: continuous, pin: out}
  - {from: seg1,  to: drop2,  rel: continuous, pin: in}
  - {from: drop2, to: hook2,  rel: looped,     pin: loop}
  - {from: drop2, to: seg2,   rel: continuous, pin: out}
  - {from: seg2,  to: sinker, rel: tied, knot: knot.improved-clinch, pin: eye}
```

Same graph, same pin vocabulary, no schema change.

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `names.trademark_note` | string | conditional | Required whenever any alias is a registered mark. |
| `variant_of` | ref | no | The rig this is a regional or minor variant of, after the accepted-name and synonym model. Prevents near-duplicates drifting apart. |
| `legality` | object | conditional | Required for any rig carrying more than three hooks, or with a known jurisdictional restriction. |
| `legality.general_warning` | bool | yes | Whether the interface surfaces a check-your-regulations notice. |
| `legality.restrictions[]` | list | no | Per-jurisdiction rules. Always `tier: C`, always sourced. Never presented as legal advice. |
| `nodes[].ref` | ref | no | The component **or knot** this node instantiates. Absent for plain line segments. |
| `nodes[].role` | `node_role` | line only | Determines which end is free. |
| `edges[].rel` | `edge_rel` | yes | |
| `edges[].knot` | ref | tied and bend only | Absent when the knot is a node instead. |
| `edges[].pin` | pin id | yes except threaded | Its type must appear in the knot's `connects.to`. |
| `travel.toward_rod` / `toward_terminal` | `travel_stop` | threaded only | What stops the component on each side. `open` against a free end is a hard error. |

---

## `pattern`

**What it is.** A named, parameterised fragment that is instantiated many times. A macro, not a loop.

**Why not call it a loop.** The word is already taken three times over: Suber's `LP`, the `loop` knot role, and the `loop` pin. A macro called "loop" inside a knot-tying schema would be permanently ambiguous.

**Why it exists.** Rigs repeat themselves constantly, and spelling the repetition out longhand is copy-paste inside data, which is the thing that rots first. A Sabiki is six identical dropper-and-hook units. An umbrella rig is five identical arm-and-jighead units. A spreader bar, a daisy chain, a crappie spider rig and a high-low are all the same shape at different counts. Written out, a Sabiki is roughly twenty near-identical nodes; as a pattern it is four lines.

**What belongs here.** Any fragment that appears more than twice, either within one record or across several.

**What does not.** A fragment used once. Indirection has a readability cost and a single use never pays it.

### Two homes, one shape

- **Inline**, under a `patterns:` key on the record, for repetition local to that rig or knot.
- **Shared**, as its own `kind: pattern` record, for a fragment several records use. `dropper-unit` is used by the high-low, Sabiki, crappie and spider rigs, so defining it once means a fix propagates instead of drifting across four files.

The shape is identical either way. A record may reference a shared pattern by id or define its own.

### Parameters and the instance index

This is the part a naive repeat cannot do.
Each instantiation needs a unique node id, and usually its own values.

- `{i}` interpolates the one-based instance index.
- `params` declares what the pattern accepts; `with` supplies it at the call site.
- `exposes` declares the pattern's own `in` and `out` connection points, which is what allows instances to chain onto each other and onto the surrounding graph.

```yaml
kind: pattern
schema_version: 0
id: pattern.dropper-unit
target: rig
params:
  - {name: hook, type: ref, required: true}
  - {name: spacing_mm, type: number, required: true}
nodes:
  - {id: "drop-{i}", ref: knot.dropper-loop}
  - {id: "hook-{i}", ref: "{hook}"}
  - {id: "seg-{i}",  type: line, role: leader, length_mm: "{spacing_mm}"}
edges:
  - {from: in,        to: "drop-{i}", rel: continuous, pin: in}
  - {from: "drop-{i}", to: "hook-{i}", rel: looped,    pin: loop}
  - {from: "drop-{i}", to: "seg-{i}",  rel: continuous, pin: out}
exposes: {in: "drop-{i}", out: "seg-{i}"}
validation: {status: unvalidated, events: []}
```

### Sabiki, using it

```yaml
kind: rig
schema_version: 0
id: rig.sabiki
names: {canonical: Sabiki rig, aliases: [bait rig, piscator rig]}
nodes:
  - {id: main,   type: line, role: main-line}
  - {id: sinker, ref: weight.bank}
expand:
  - pattern: pattern.dropper-unit
    count: 6
    from: main
    to: sinker
    with: {hook: hook.sabiki, spacing_mm: 150}
edges:
  - {from: expand.last, to: sinker, rel: tied, knot: knot.improved-clinch, pin: eye}
```

Six units, four lines.
Changing the spacing is one edit rather than six, and the high-low rig above becomes the same declaration with `count: 2`.

### Knots use the same mechanism

Per-action `repeat` covers a single repeated action, such as five turns in an improved clinch.
It does **not** cover a repeated *sequence*, which is what an FG knot's alternating weave actually is.

```yaml
patterns:
  fg-weave:
    target: knot
    stages:
      - actions: [{verb: MT, subject: RP, around: SP, rotation: CW,  direction: U}]
      - actions: [{verb: MT, subject: RP, around: SP, rotation: CCW, direction: D}]
stages:
  - {id: 1, actions: [...]}
  - {id: 2, expand: {pattern: fg-weave, count: [15, 20]}}
  - {id: 3, actions: [...]}
```

`count` accepts a range for the same reason `repeat` does: real instructions say "fifteen to twenty weaves."

**The rule dividing them:** one action uses `repeat`, more than one uses a pattern.
A pattern containing exactly one action and taking no parameters should have been a `repeat`, and CI warns about it.

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `target` | `pattern_target` | yes | Whether the fragment expands into a rig graph or a knot stage list. A pattern is never valid in both. |
| `params[]` | list | no | Declared inputs. A pattern with no params is a literal fragment. |
| `params[].type` | enum | yes | `ref`, `number`, `string`, `range`. |
| `params[].required` | bool | yes | An unsupplied required param is a hard error, never a silent default. |
| `nodes` / `edges` | list | rig target | Fragment graph. May use `{i}` and `{param}` interpolation. |
| `stages` | list | knot target | Fragment stages. Same interpolation. |
| `exposes.in` / `.out` | node id | rig target | The fragment's own connection points, without which instances cannot chain. |
| `expand[].pattern` | ref or local name | yes | Shared pattern id, or a key from the inline `patterns:` block. |
| `expand[].count` | int or range | yes | How many instances. |
| `expand[].from` / `.to` | node id | rig target | What the first instance chains onto and the last chains into. |
| `expand[].with` | object | conditional | Values for the pattern's params. |
| `expand.last` | reserved id | n/a | Refers to the final instance's `out`, for edges written after the expansion. |

### What expansion means for validation

Patterns are expanded **before** any graph rule runs.
Rules 1 through 5 operate on the fully expanded graph, never on the compressed form, so a pattern cannot hide a structural error.

Expansion is deterministic: the same record always yields the same node ids, so diffs stay meaningful.

The honest cost: a pattern-heavy rig is less literal to read than one written longhand, and a reader has to follow an indirection to see the whole graph. That is paid for by the diff, where changing spacing across six droppers is one line instead of six, and by the fix that propagates to four rigs at once instead of drifting.

---

## `technique`

**What it is.** How a rig is worked in the water, and the conditions under which that pairing is a good answer.

**Why it is separate from the rig.** The relationship is many to many. Drop shot with a shaking retrieve and drop shot with a dragging retrieve are different answers suiting different conditions.

**Conditions attach to the pairing, never to the rig.**
Hang them on the rig and that distinction is destroyed before anyone notices.

```yaml
kind: technique
schema_version: 0
id: technique.carolina-drag
rig: rig.carolina
motion:
  primary: drag
  contact: bottom
  cadence: slow sweep, pause two to five seconds
water_column: bottom
applicability:
  species: [species.largemouth-bass, species.smallmouth-bass]
  cover: [gravel, sand, sparse-grass]
  clarity: [stained, clear]
  depth_m: [2.4, 7.6]
  water_temp_c: [13, 27]
  season: [prespawn, postspawn, summer]
  tier: C
  rank: normal
  sources: [src.example]
failure_modes:
  - Dragged too fast, the weight outruns the bait and the presentation collapses.
  - In heavy grass the leader fouls and the rig fishes dead.
validation: {status: unvalidated, events: []}
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `rig` | ref | yes | Many techniques may reference one rig. |
| `motion.primary` | `motion_primary` | yes | Drives the animated lure path. |
| `motion.contact` | `contact` | no | |
| `motion.cadence` | string | no | Free text by necessity; rhythm resists enumeration. |
| `water_column` | `water_column` | yes | Used by the side-view animation. |
| `applicability.species` | list of refs | yes | References, not bare strings. |
| `applicability.cover` | list of `cover` | no | |
| `applicability.clarity` | list of `clarity` | no | |
| `applicability.season` | list of `season` | no | |
| `applicability.depth_m` | range | no | A range, never a point value. |
| `applicability.water_temp_c` | range | no | Degrees Celsius, as with everything else. |
| `applicability.tier` | `tier` | yes | Effectively always `C`. |
| `applicability.sources` | list | yes | At least one, enforced. |
| `failure_modes` | list | no | How this presentation goes wrong. Rare in existing references and disproportionately useful. |

---

## `species`

**What it is.** A target fish, with regional common names reconciled to one accepted name.

**Why it exists.** Regional naming chaos is worse for fish than for rigs. "Sheepshead" means two unrelated species depending on where you are standing.

```yaml
kind: species
schema_version: 0
id: species.bluegill
common_name: Bluegill
scientific_name: Lepomis macrochirus
scientific_name_source: src.example
aliases: [bream, brim, copperbelly]
ambiguous_aliases:
  - {name: sunfish, note: Also used for other Lepomis species and regionally for crappie.}
water: freshwater
validation: {status: unvalidated, events: []}
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `scientific_name_source` | ref | yes | **Required.** Black bass taxonomy was revised recently, so this is not safe to assert from memory. |
| `ambiguous_aliases` | list | no | Names meaning different fish in different places. Never used for automatic matching. |
| `water` | `water` | yes | |

---

## `source`

**What it is.** A citation, together with an explicit record of what we are legally permitted to do with it.

**Why `copy_policy` exists.** The legal position is read everything, cite everything, draw everything ourselves, copy nothing. That rule only holds if the constraint travels with the citation, so a contributor reading a source record is told at the point of use whether they may quote it.

```yaml
kind: source
schema_version: 0
id: src.example
type: agency
title: Example Fishing Guide
author: Example State Department of Natural Resources
year: 2019
url: https://example.gov/guide.pdf
accessed: 2026-08-10
license: all-rights-reserved
copy_policy: cite-only
reliability: medium
notes: >
  State agency publication. State works are not covered by 17 USC 105,
  so this is presumptively copyrighted despite being government-produced.
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `type` | `source_type` | yes | |
| `url` / `accessed` | string / date | conditional | Required for anything online. Link rot is the default outcome. |
| `license` | SPDX or literal | yes | An SPDX identifier, `all-rights-reserved`, or `unknown`. Not an enum, because SPDX is the registry. |
| `copy_policy` | `copy_policy` | yes | **The guardrail.** Anything not verified as openly licensed defaults to `cite-only`. |
| `reliability` | `reliability` | yes | A manufacturer's diagram of their own product is `high`; a forum post is `low`. |

---

## Rules that fall out

| # | Rule | Severity | Tier |
| --- | --- | --- | --- |
| 1 | Every `threaded` node must have a `blocks_passage` stop between it and each **free** end of its segment | error | A |
| 2 | A knot's `connects.to` must include the type of the pin it ties to | error | A |
| 3 | `bore_mm` must exceed the `diameter_mm` of whatever passes through it | error | A |
| 4 | The rig graph must be connected and acyclic, with no orphan nodes | error | A |
| 5 | Every `ref` must resolve, and every `pin` must exist on its target | error | A |
| 6 | `id` must match the grammar, and its category must be an `id_category` member | error | A |
| 7 | `kind` must match the directory the record lives in | error | A |
| 8 | Every enum-typed field must hold a registered value | error | A |
| 9 | An emitted `notation` string must match its structured actions | error | A |
| 10 | Any `tier: C` claim must carry at least one source | error | A |
| 11 | A strength claim without `n` is rejected | error | A |
| 12 | Any alias that is a registered mark requires `trademark_note` | error | A |
| 13 | A rig with more than three hooks requires a `legality` block | error | A |
| 14 | A knot with `role: loop` or `stopper` must appear as a node, never an edge | error | A |
| 15 | A `deprecated` rank requires a reason | error | A |
| 16 | Any source with `license: unknown` is forced to `copy_policy: cite-only` | error | A |
| 17 | `ambiguous_aliases` are never used for automatic species matching | error | A |
| 18 | A `validation.status` above `unvalidated` requires at least one matching event | error | A |
| 19 | `source-corroborated` events require at least two distinct sources | error | A |
| 20 | `validation.events` is append-only; rewriting history fails the diff check | error | A |
| 21 | A knot's `line_types` should include the line type of the segment it ties on | warning | B |
| 22 | No brand, manufacturer, product line, or model name appears in a normative field. Structural for typed fields, since none exists to hold one; review-enforced for free text | error | A / B |
| 31 | A component's `variants.axis` must be a physical quantity from `variant_axis`, never a trade label | error | A |
| 32 | A value a rig selects must exist in the referenced component's `variants.values` | error | A |
| 33 | Every quantity field carries an SI suffix from the `unit` enum. Non-SI suffixes such as `_lb`, `_oz`, `_in`, `_ft` are rejected outright | error | A |
| 34 | Quantity values must fall in a plausible range for their unit, catching unconverted imperial figures entered by mistake | warning | B |
| 23 | Nothing publishes with `validation.status: unvalidated` | warning | B |
| 24 | Patterns expand before rules 1 to 5 run, so a pattern can never hide a structural error | error | A |
| 25 | Expanded node ids must not collide with declared ids, or with each other | error | A |
| 26 | Every required `param` must be supplied by `with`; there are no silent defaults | error | A |
| 27 | A rig-target pattern must declare `exposes.in` and `exposes.out`, or it cannot chain | error | A |
| 28 | A pattern's `target` must match where it is expanded; a knot pattern in a rig is rejected | error | A |
| 29 | A pattern with one action and no params should be a `repeat` instead | warning | B |
| 30 | `expand.count` must be a positive integer or range, bounded by a sanity cap | error | A |

Rule 1 is the one worth noticing.
"Sliding weight with nothing to stop it" is the most common real rigging error, and against this shape it is a decidable interval check rather than a matter of opinion.

## Open questions

- **Is `prose` authored or generated?** Generating it from the Suber verbs would guarantee the words match the animation, which is precisely the Tier C gap nothing else closes. But generated prose reads like a machine, and this is a teaching site. Currently authored, with CI checking only that stage counts agree.
- **Where does geometry come from?** Suber's notation says what to do, not where anything is. Either a solver derives positions from the actions, or stages need an optional `geometry` block of seed control points. That field cannot honestly be designed before a renderer exists, so it is absent rather than guessed at.
- **Do handedness variants come free** as a transform over `chirality` and the direction fields, or do some knots need explicit encodings?
- **Where do generated assets live**, and are they committed or built?
- **Does `rel: rigged`** (bait onto a hook) belong in the rig graph, or is it presentation?
- **Suber's notation carries a bare copyright notice with no licence grant.** Implementable under the idea and expression distinction, but adopting his vocabulary wholesale warrants asking him directly.
