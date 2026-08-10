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
| `rigging` | A procedure for mounting a soft body on a hook | `data/riggings/` | `rigging` |
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

A rig calls for a circle hook with an eleven millimetre gap, braid rated near four and a half kilograms, and a fourteen gram egg sinker.
It never calls for a named product.

The reason is not squeamishness about commerce, it is that a trade label carries no information.
Research on this domain established that size labels are unstandardised across manufacturers *and* across patterns within one manufacturer: a 2/0 octopus hook is not a 2/0 worm hook, and one maker's size 3 split ring is rated near 25 lb where its own 3H is near 65 lb.

An earlier draft responded to that by making `manufacturer` a **required** field, so a label could at least be disambiguated.
That was working around the problem instead of removing it.
The brand was only ever a proxy for a physical property, so the spec records the **property itself** and the proxy becomes unnecessary:

| Instead of | Record |
| --- | --- |
| size 3 split ring | `rating_kg: 11` |
| 2/0 octopus hook | `gap_mm: 11`, `point_style: turned-in`, `shank: short` |
| 10 lb fluorocarbon | `material: fluoro`, `breaking_load_kg: 4.5`, `diameter_mm: 0.285` |
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

### Common fields

These fields exist on every kind and are not repeated in each kind's table.

| Field | Required | Meaning |
| --- | --- | --- |
| `kind` | yes | The record type. |
| `id` | yes | `{category}.{name}`, per the grammar above. |
| `schema_version` | yes | Which version of this spec the record targets. |
| `validation` | yes | How it was checked. |
| `former_ids` | no | Previous IDs, for reference migration. |
| `refs` | no | External identifiers. See below. |

### `refs`, external identifiers

An earlier draft carried `abok_ref` on knots.
That is a field per catalogue waiting to happen, and `wikidata_ref`, `gbif_ref` and `agrovoc_ref` would have followed it.

One generic list instead, on every kind:

```yaml
refs:
  - {system: abok, id: "1204", source: src.abok-general}
  - {system: wikidata, id: "Q207754"}
```

**An identifier is not a citation**, and the two do not collapse into each other.
`sources` says *this claim came from there*.
`refs` says *this same thing is called that elsewhere*.
Wikidata keeps them apart for the same reason and so do we.

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `system` | `ref_system` | yes | Which catalogue. Registered below. |
| `id` | **string** | yes | Always a string, never a number. Catalogue identifiers carry leading zeros, prefixes such as `Q207754`, and full URIs. |
| `source` | ref | conditional | Required when the system is not resolvable, since nothing can check it automatically. |
| `rank` | `rank` | no | A ref later found wrong is demoted, not deleted, exactly like any other claim. |

An unknown identifier is **omitted**, never guessed.
`null` was doing that job on `abok_ref` and an absent entry says the same thing without occupying a field.

### The `ref_system` registry

A **resolvable** system can be fetched and checked, so its identifiers are Tier B and machine-verifiable: resolve it, confirm the record exists, confirm its label matches. A print-only system is Tier C and needs a human with the book.

| System | Identifies | Resolvable | Notes |
| --- | --- | --- | --- |
| `abok` | Knots, by Ashley number | No | Print only. The **number** may be cited; text and illustrations may never be reproduced, since the book is in copyright until 2040. |
| `wikidata` | Anything, by Q-id | Yes | |
| `gbif` | Species, by usage key | Yes | |
| `itis` | Species, by TSN | Yes | |
| `worms` | Marine species, by AphiaID | Yes | |
| `fao-isscfg` | Fishing gear types | Yes | Records the link and nothing more. All recreational angling is the single leaf `25.0.0`, so it cannot describe a rig. |
| `agrovoc` | Vocabulary concepts, by URI | Yes | Tackle coverage is two concepts total, `fishing lines` and `hooks`. Almost nothing here will have one. |

The URL template belongs to the system, not the record, so a record stores the bare identifier and the renderer builds the link.

### Units

**No physical quantity is ever a bare number.**
The unit is part of the field name, always, so a value can never be misread and never needs parsing.

**SI where possible, with genuinely universal fishing standards as the only exception.**

| Quantity | Suffix | Unit |
| --- | --- | --- |
| Small dimensions | `_mm` | millimetre |
| Distance and depth | `_m` | metre |
| Mass | `_g` | gram |
| Breaking load and rating | `_kg` | kilogram, see below |
| Temperature | `_c` | degree Celsius |
| Time | `_s` | second |
| Angle | `_deg` | degree |

"SI" here means SI base and derived units, their prefixed forms, and the units the SI Brochure accepts for use with SI, which is what makes degree Celsius and the plane degree legal.

An earlier draft stored pounds, ounces, inches and Fahrenheit on the grounds that this is how the domain writes things.
That was the same mistake as recording a manufacturer: it stored the **label** instead of the **property**.
How anglers write a number is a display concern, and the display layer converts freely.
A reader in Ohio sees "10 lb braid" and a reader in Queensland sees "4.5 kg", from one stored value.

### The exception test

A non-SI unit is permitted only when the fishing world has **one** convention for that quantity worldwide.
Where usage is split by region, SI wins and the display layer converts.

| Quantity | Convention | Verdict |
| --- | --- | --- |
| Line diameter | Millimetres everywhere, including on United States packaging | Already SI |
| Breaking load | Pounds in the United States, kilograms elsewhere; IGFA publishes both | Split, so SI |
| Sinker and lure mass | Ounces in the United States, grams elsewhere | Split, so SI |
| Depth | Feet in the United States, metres elsewhere | Split, so SI |
| Water temperature | Fahrenheit in the United States, Celsius elsewhere | Split, so SI |
| Leader length | Inches and feet in the United States, centimetres elsewhere | Split, so SI |

**The exception register is currently empty**, and that is worth stating rather than inventing an entry to justify the rule.
Every genuinely universal non-SI standard in angling turns out to belong to fly fishing: AFTMA line weight, which is a grain-based scale used identically worldwide, and tippet X sizing, likewise universal.
Fly fishing is out of scope for v1, so nothing qualifies yet.

The clause stays because that scope will change, and because it records the reasoning so a future contributor does not "correct" a legitimate exception into SI and break it.

### Breaking load is in kilograms

`breaking_load_kg` and `rating_kg` follow the international line-class convention, which IGFA publishes in kilograms and which is what metric-market packaging prints.

**Precisely, this is a load, not a mass.** What is being described is the force at which something fails, so the strictly correct SI unit is the newton. Kilogram-force is used here because it is the settled convention in the one place the world agrees, and storing newtons would produce a spec no angler could sanity-check while gaining nothing a conversion cannot supply. The field name carries `_kg` and this paragraph carries the caveat.

Authoring in any unit is a tooling concern.
The authoring CLI accepts pounds, ounces or newtons and normalises on write, because hand-conversion is where errors come from and no author should be doing it.

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
| `kind` | `component`, `line`, `knot`, `rigging`, `rig`, `pattern`, `technique`, `species`, `source` |
| `id_category` | **components:** `hook`, `weight`, `swivel`, `snap`, `split-ring`, `bead`, `float`, `stop`, `blade`, `skirt`, `jighead`, `bait`, `lure`, `hardware`, `sleeve` · **others:** `line`, `knot`, `rigging`, `rig`, `pattern`, `technique`, `species`, `src` |
| `pattern_target` | `rig`, `knot` |
| `pin_type` | `closed-eye`, `open-eye`, `split-ring`, `shank`, `snap`, `arm-socket`, `line-end`, `in`, `out`, `loop` |
| `severity` | `error`, `warning`, `info` |
| `ref_system` | `abok`, `wikidata`, `gbif`, `itis`, `worms`, `fao-isscfg`, `agrovoc` |

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
| `mounting` | `tied`, `threaded`, `rigged` |
| `caution_severity` | `advisory`, `strong`, `do-not-use` |
| `unit` | `mm`, `m`, `g`, `kg`, `c`, `s`, `deg` |
| `variant_axis` | `mass_g`, `gap_mm`, `diameter_mm`, `rating_kg`, `breaking_load_kg`, `length_mm`, `supports_g` |
| `point_style` | `straight`, `turned-in`, `turned-out`, `knife-edge`, `needle` |
| `shank` | `short`, `standard`, `long`, `extra-long` |
| `bend_style` | `round`, `octopus`, `circle`, `worm`, `ewg`, `offset-worm`, `straight-shank`, `kahle`, `aberdeen`, `siwash`, `treble` |
| `body_profile` | `straight`, `curly-tail`, `paddle-tail`, `creature`, `craw`, `fluke`, `tube`, `grub`, `worm-live`, `baitfish`, `cut` |
| `landmark_surface` | `n` nose, `d` dorsal (back), `v` ventral (belly), `l` lateral (side), `t` tail |
| `rig_verb` | `IN` insert, `OUT` exit, `TH` thread, `SL` slide, `RO` rotate, `SK` skin, `BU` bury |
| `rig_descriptor` | `AL` aligned straight, `CE` centred |
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
  axis: rating_kg
  values: [11, 16, 25, 36, 59]
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
| `mounting` | `mounting` | yes | `tied` connects through a pin, `threaded` rides on a line segment, `rigged` is a soft body mounted on a hook by a procedure. This single field separates a swivel from a sinker from a worm, and determines which edge types are legal. |
| `soft` | bool | no | Whether the body can be pierced. Gates `landmarks` and makes `rel: rigged` legal against it. |
| `landmarks` | map | soft only | Named positions along the body, normalised from 0 at the nose to 1 at the tail. A `rigging` addresses these. |
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
  - {breaking_load_kg: 4.5, diameter_mm: 0.285, source: src.example}
  - {breaking_load_kg: 9.1, diameter_mm: 0.405, source: src.example}
validation: {status: unvalidated, events: []}
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `material` | `line_material` | yes | The material family. |
| `properties.knot_sensitivity` | `knot_sensitivity` | yes | How much the material punishes a poorly chosen knot. Fluorocarbon and braid are both `high` for different reasons, which is what makes rule 16 worth having. |
| `diameters[].breaking_load_kg` | number | yes | Breaking load in kilograms, the international line-class convention. The familiar "10 lb test" is a display conversion of this. |
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

### Arity is what makes a knot placeable

An earlier draft made this a property of `role`: terminal and bend knots were edges, loop and stopper knots were nodes.
Authoring a drop shot broke that.

A drop shot ties a Palomar to the hook and leaves a long tag that carries the weight, so the same Palomar that is an edge in a Texas rig has **three** connection points here: the main line, the hook eye, and the tag continuing down. Role has not changed; the number of things attached to it has.

**The rule is arity, and role only predicts it:**

| Connection points | Placement |
| --- | --- |
| 2 | an **edge** carrying `knot:` |
| 3 or more | a **node** whose `ref` is the knot |

| Role | Joins | Usual arity |
| --- | --- | --- |
| `terminal` | line end to a hardware pin | 2, but 3 when a tag is load-bearing |
| `bend` | two line ends | 2 |
| `loop` | tied *within* one line, producing a branch | 3: `in`, `out`, `loop` |
| `stopper` | tied *on* one line, blocking passage | 2, but a node because it blocks passage |
| `arbor` | line to a reel spool | 2 |

A knot acting as a node is structurally a component made of line, so `blocks_passage` applies to it exactly as it does to hardware.

An edge leaving a knot node carries **no `knot:` reference**, because the node is already the knot. The pin names which part of it: `in`, `out`, `loop`, `tag`, or `eye`.

```yaml
kind: knot
schema_version: 0
id: knot.palomar
role: terminal
names:
  canonical: Palomar knot
  aliases: [Palomar]
refs: []
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
| `strength.claims[].n` | int | yes | Sample size. **Required.** Most published knot tests are n=1, and `n: 0` with a null `residual_pct` is the legitimate way to say "no number is published here, and here is why." |
| `strength.claims[].note` | string | no | Why a figure is absent, or what qualifies it. |
| `strength.claims[].tier` | `tier` | yes | Always `C`. |
| `strength.claims[].rank` | `rank` | yes | |
| `stages[].descriptors` | list | no | Suber descriptors such as `PL` parallel or `CO` crosses-over. Assertions about state rather than motions. A stage may hold only descriptors. |
| `cautions` | list | no | Structured warnings, each with a `caution_severity`, a note, and sources. A `do-not-use` caution is rendered prominently and is how the reference says "this knot exists and you should not use it for this." |
| `failure_modes` | list | no | How the knot goes wrong in practice. |

---

## `rigging`

**What it is.** A procedure for mounting a soft body on a hook.

**What belongs here.** Soft plastic presentations such as Texposed, exposed Texas, wacky, Neko, nose-hooked, threaded on a jighead, and screw-lock. Live and cut bait hooking is the same class of problem: lip-hooked, dorsal-hooked, and bridle-rigged are procedures with a spatial outcome, not attributes.

**What does not.** Hard lures, which arrive with hooks attached and are single components. Anything about how the finished rig is *fished*, which is a `technique`.

### Why this is a kind and not a field

An earlier draft carried `rel: rigged, style: texposed` on a rig edge and listed "does bait onto hook belong in the graph at all" as an open question.

It does, and the reason is that **rigging a soft plastic is structurally the same problem as tying a knot**:

- an ordered sequence of actions,
- performed on a flexible body,
- producing a specific spatial relationship,
- which fails in a specific way when done wrong (a bait rigged crooked spins and twists line, exactly as a badly dressed knot slips),
- and which a static image cannot teach.

So `knot` and `rigging` are two instances of one shared shape, a **staged procedure**: ordered stages, each holding simultaneous actions, each stage a keyframe.
The two kinds differ only in vocabulary, so the animation player is written once and driven by both.

In CUE this is a shared definition both kinds embed, not duplicated structure.

### Soft bodies need anatomy

A hook attaches at a pin. A bait is *pierced somewhere*, so a soft-body component declares **landmarks**: named positions along the body, normalised from nose to tail.

Normalised rather than absolute, so one rigging record works across a 100 mm worm and a 180 mm worm without change.

```yaml
kind: component
schema_version: 0
id: bait.stick-worm
name: Stick worm
body_profile: straight
soft: true
landmarks:
  nose:     0.00
  collar:   0.08
  egg-sack: 0.45
  tail:     1.00
variants:
  axis: length_mm
  values: [100, 125, 150, 180]
```

An action names both a landmark and a **surface**, because "insert at the nose" and "exit through the belly" are different faces of the same body.

### Texposed, worked

### The notation

No notation for this exists, in fishing or in any adjacent domain.
That was established rather than assumed, and the search is recorded in [`docs/research/05-rigging-notation.md`](research/05-rigging-notation.md) with the reasoning in [ADR 0006](../adr/0006-rigging-notation.md).

So it is invented, but assembled from proven pieces rather than from nothing:

- **Grammar and architecture from [knitout](https://textiles-lab.github.io/knitout/knitout.html)**, the CMU Textiles Lab's machine-knitting format. `VERB DIRECTION LOCATION SUBJECT` ordering, a flat instruction list with no flow control, and surface-plus-index addressing.
- **Geometry from needle-insertion planning.** Alterovitz's insertion plan of location, angle, bevel roll and depth maps onto bait rigging field for field, and closed a real hole: an earlier draft had position and depth but no **angle** and no **roll**, both of which decide where a point actually exits and whether it tracks straight.
- **Verb-set scale from surgical gesture taxonomies.** JIGSAWS' 15 surgemes and the SAGES Delphi taxonomy's 24 gestures put the right granularity at roughly this size.

**Addressing.** A location is a surface letter and a normalised position from 0 at the nose to 1 at the tail, so `v0.30` is the ventral surface thirty percent down the body. Normalised rather than absolute, so one record works on every bait length. Named landmarks resolve to positions, so `collar` and `n0.08` are the same place.

**Two verbs are genuinely ours.** `SK` skin, passing tangentially just beneath a surface, and `BU` bury, terminating subsurface without exiting, have **no prior art in any domain examined**. Surgery names only one of them, and only because subcuticular closure needs it. They are recorded as inventions rather than borrowings, because that is the honest measure of the gap.

**Extensions** use knitout's `x-` namespace, so a third party can add a verb without forking the spec.

```yaml
kind: rigging
schema_version: 0
id: rigging.texposed
names:
  canonical: Texposed
  aliases: [Texas rigged weedless, tex-posed, skin hooked]
weedless: true
applies_to:
  body_profile: [straight, curly-tail, creature, craw]
  hook_bend: [ewg, offset-worm, straight-shank]
stages:
  - id: 1
    actions:
      - {verb: IN,  at: n0.00, angle_deg: 90, roll_deg: 0, depth: 0.06}
      - {verb: OUT, at: v0.06}
    prose: Insert the point into the centre of the nose and bring it out through the belly about a quarter inch back.
    notation: "IN n0.00 a90 r0 z0.06, OUT v0.06"
  - id: 2
    actions:
      - {verb: SL, subject: bait, along: shank, until: eye}
    prose: Slide the bait up the shank until the nose seats against the hook eye.
    notation: "SL bait shank>eye"
  - id: 3
    actions:
      - {verb: RO, subject: hook, relative_to: bait, angle_deg: 180}
    prose: Rotate the hook 180 degrees so the point faces the body.
    notation: "RO hook/bait a180"
  - id: 4
    descriptors:
      - {desc: AL, subject: bait, against: shank}
    prose: Lay the hook alongside the bait and check it hangs perfectly straight before piercing.
    notation: "AL bait/shank"
  - id: 5
    actions:
      - {verb: IN, at: v0.30, angle_deg: 90, roll_deg: 0}
      - {verb: BU, at: d0.30, depth: 0.01}
    prose: Bring the point back through the body and bury it just under the skin on the back.
    notation: "IN v0.30 a90 r0, BU d0.30 z0.01"
failure_modes:
  - Bait bunched or curved on the shank, which makes it spin and twist the line. Stage 4 exists to prevent this and is the step most often skipped.
  - Point buried too deep, which costs hooksets.
  - Entry not centred in the nose, which makes the bait ride crooked.
```

Stage 4 is not padding.
A Texas rig that is not straight on the shank spins, and that is the single most common way this presentation is got wrong.
It is an action with no piercing in it, which is exactly the kind of step photo sequences omit and a staged procedure keeps.

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `weedless` | bool | yes | Whether the finished presentation hides the point. Drives cover suitability in a technique. |
| `applies_to.body_profile` | list of `body_profile` | yes | Which bait shapes this works on. Mismatch is a hard error in a rig. |
| `applies_to.hook_bend` | list of `bend_style` | yes | Which hook geometries this works on. |
| `stages[].actions[].verb` | `rig_verb` | yes | `IN`, `OUT`, `TH`, `SL`, `RO`, `SK`, `BU`. |
| `actions[].at` | address | piercing verbs | Surface letter plus normalised position, such as `v0.30`. A declared landmark name resolves to one, so `collar` and `n0.08` are the same place. |
| `actions[].angle_deg` | number | `IN`, `RO` | On `IN`, the entry angle against the surface, which decides where the point exits. On `RO`, the rotation applied. |
| `actions[].roll_deg` | number | `IN` | Rotational orientation of the point about its own axis. A bevelled point tracks in the direction it faces, so this decides whether it runs true or wanders. |
| `actions[].depth` | number | no | Normalised depth of travel through the body. |
| `stages[].descriptors` | list | no | Assertions rather than motions, from `rig_descriptor`. A stage may hold only descriptors, which is how the alignment check exists as a step. |
| `stages[].notation` | string | derived | Emitted from the structured actions. Never hand-edited; CI regenerates and diffs it. |
| `failure_modes` | list | no | How this presentation goes wrong. |
| `x-*` | any | no | Extension namespace, after knitout. A third party may add fields here without forking, and the validator ignores them. |

Scalar prefixes in the emitted notation are `a` angle, `r` roll, `z` depth.
Depth is `z` rather than `d` because `d` already means the dorsal surface, and a notation that collides with itself is worse than a verbose one.

### In a rig

The parallel with knots is exact, which is the point:

```yaml
- {from: leader, to: hook, rel: tied,   knot: knot.palomar,      pin: eye}
- {from: hook,   to: bait, rel: rigged, rigging: rigging.texposed}
```

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
  - {id: main,   type: line, role: main-line, material: braid, breaking_load_kg: 13.6}
  - {id: sinker, ref: weight.egg,    mass_g: 14.2}
  - {id: bead,   ref: bead.glass,    diameter_mm: 8}
  - {id: swivel, ref: swivel.barrel, rating_kg: 25}
  - {id: leader, type: line, role: leader, material: fluoro, breaking_load_kg: 6.8, length_mm: [305, 1220]}
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
| `failure_modes` | list | no | How the assembled rig goes wrong. Distinct from a technique's failure modes, which are about fishing it rather than building it. |
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

**What belongs here.** Any multi-action fragment that appears more than once, either within one record or across several.

An earlier draft said "more than twice." Authoring the double uni disproved it: that knot is one uni tied twice facing opposite ways, and writing it out longhand duplicates four stages for no reason. Twice is already duplication when the fragment has more than one action in it.

**What does not.** A fragment used once, or a single action repeated, which is what `repeat` is for. Indirection has a readability cost and neither of those pays it.

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
| 33 | Every quantity field carries an SI suffix from the `unit` enum. Non-SI suffixes such as `_lb`, `_oz`, `_in`, `_ft` are rejected unless registered in the exception table, which is currently empty | error | A |
| 35 | A knot with three or more connection points in a given rig must appear as a node, never an edge | error | A |
| 36 | An edge leaving a knot node must not carry a `knot:` reference | error | A |
| 37 | `rel: rigged` requires the target component to have `soft: true` | error | A |
| 38 | A `rigging` address must resolve to a declared landmark or a position in `[0, 1]` | error | A |
| 39 | A `do-not-use` caution requires at least one source | error | A |
| 40 | Records referencing `src.needs-citation` are counted and reported as outstanding citation debt | warning | B |
| 41 | A `refs` entry on a non-resolvable system requires a `source` | error | A |
| 42 | A `refs` `id` is always a string, never a bare number | error | A |
| 43 | A `refs` entry on a resolvable system is fetched and its label checked against the record's canonical name | warning | B |
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
- **What notation should `rigging` use?** The structure is settled and the verb vocabulary is not. Whether an established system exists for describing a rigid object piercing a soft body, and what to take from surgical suture paths, textile stitch charts, or needle-insertion planning, is under research in [`docs/research/05-rigging-notation.md`](research/05-rigging-notation.md).
- **Suber's notation carries a bare copyright notice with no licence grant.** Implementable under the idea and expression distinction, but adopting his vocabulary wholesale warrants asking him directly.
