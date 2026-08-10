# Schema prototype

Exploratory shapes, not settled.
The point of this document is to find where the model breaks before any generator exists.
Everything here is written to be argued with.

Every field is documented.
Where a field exists only so a validator can reason about it, that is stated, because those fields must never be shown to a reader.

## Conventions

### `kind` and `id`

Every record carries both, and they are not redundant.

- **`kind`** is the record type: `component`, `knot`, `rig`, `technique`, `source`. It is the schema discriminator, and it is what tells a validator which rules apply.
- **`id`** is a globally unique, stable, human-legible handle, namespaced by category: `weight.egg`, `swivel.barrel`, `knot.palomar`.

The namespace is deliberate.
References appear far from the definitions they point at, and `ref: weight.egg` is self-documenting at the call site in a way that `ref: egg` is not.
In hand-authored data that readability is worth more than the purity of an opaque identifier.

There is no separate `category` field.
An earlier draft carried both `id: weight.egg` and `category: weight`, which is exactly the duplication that drifts.
The ID prefix is the category, and CI enforces that the prefix is a member of the allowed set.

The cost of this choice is real: recategorizing a component means renaming it, which breaks every reference to it.
Two mitigations, both enforced in CI:

- `former_ids` records any previous ID, so references can be migrated mechanically and old links redirect.
- An ID, once published, is never reused for a different thing.

The `kind` field duplicates what the directory already implies, and that is intentional.
Records get concatenated into bundles, passed around individually, and exported as one JSON blob, all of which lose directory context.
A self-describing record survives that.

### Tiers

Any field carrying a factual claim may be annotated with the tier it earned, per the README's correctness model.
`A` is machine-provable, `B` is heuristically checkable, `C` is citation or human expert only.
Every `C` claim must carry at least one source.

## The primitive that holds it together

Knots and rigs share one vocabulary: the **typed connection point**, or pin.

A component exposes pins with types.
A knot declares which pin types it can terminate on.
A rig declares which pin a given knot ties to.

Because both speak the same vocabulary, a rig is a netlist and validation is electrical rule checking.
Design the knot schema and the rig schema independently and they will not compose.

Current pin types: `closed_eye`, `open_eye`, `split_ring`, `shank`, `snap`, `line_end`, `arm_socket`.

## Components

The reusable atom.
Roughly 50 of these yield 100+ rigs.

```yaml
kind: component
id: swivel.barrel
name: Barrel swivel
pins:
  - {id: eye.a, type: closed_eye, wire_diameter_mm: 0.9}
  - {id: eye.b, type: closed_eye, wire_diameter_mm: 0.9}
blocks_passage: true
sizes:
  - {system: swivel_number, label: "7",  manufacturer: rosco, rating_lb: 55}
  - {system: swivel_number, label: "10", manufacturer: rosco, rating_lb: 35}
```

```yaml
kind: component
id: weight.egg
name: Egg sinker
mounting: threaded
bore_mm: 2.4
blocks_passage: false
pins: []
sizes:
  - {system: ounce, label: "1/2", grams: 14.2}
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `kind` | enum | yes | Always `component`. |
| `id` | string | yes | Namespaced identifier. Prefix is the category and must be in the allowed set. |
| `name` | string | yes | Display name. |
| `mounting` | enum | yes | `tied` if it connects through a pin, `threaded` if it rides on a line segment instead. This is the distinction that separates a swivel from a sinker, and it changes which edge types are legal. |
| `pins` | list | yes | Connection points. A threaded component has none, and an empty list is meaningful rather than missing data. |
| `pins[].id` | string | yes | Local identifier, referenced by rig edges. |
| `pins[].type` | enum | yes | Pin type. A knot must declare this type in `connects.to` to be legal here. |
| `pins[].wire_diameter_mm` | number | no | Physical wire gauge of the eye. Feeds knot-fit warnings. |
| `blocks_passage` | bool | yes | Whether a threaded component sliding down the line is stopped by this one. **Validator-only field, never displayed.** It is what makes the sliding-weight check decidable. |
| `bore_mm` | number | threaded only | Largest thing that can pass through. **Validator-only.** |
| `sizes` | list | no | Available sizes. Never a bare number, see below. |
| `sizes[].system` | enum | yes | Which sizing system the label belongs to (`ounce`, `aught`, `swivel_number`, `wire_size`). |
| `sizes[].label` | string | yes | The printed label, as a **string**. `"1/2"` and `"3H"` are not numbers. |
| `sizes[].manufacturer` | string | no | Required whenever the system is manufacturer-specific. |
| `sizes[].rating_lb` | number | no | The **only** field comparable across brands. |
| `former_ids` | list | no | Previous IDs, for reference migration. |

Size labels are not comparable across manufacturers, and often not within one.
A 2/0 octopus hook is not a 2/0 worm hook, and Rosco's own size 3 split ring is rated around 25 lb where its 3H is around 65 lb.
Comparing size labels numerically is a hard validation error.

## Knots

A knot is a **procedure**, not a shape, so the record is a sequence of stages.

### On Suber's notation

The vocabulary here is taken from [Peter Suber's Knot Tying Notation](https://legacy.earlham.edu/~peters/knotting/notate.htm) (2002 to 2004).

Adopted, because reinventing them would be strictly worse:

- The **verb set**: `GP` grip, `MB` make bight, `ML` make loop, `MT` make turn, `MV` move, `RV` reeve, `TW` twist.
- The **noun vocabulary**: `RP` running part, `SP` standing part, `BT` bight, `LP` loop, `CS` cord segment, `E` end.
- **Stages**, where all actions in a stage happen simultaneously. This maps one to one onto animation keyframes and is the single most useful thing in the system.
- **Directions** `F A L R U D`, rotations `CW` and `CCW`, and reeve direction pairs like `F-A`.
- **Loop chirality**, `/` versus `\`, which is the crossing-handedness primitive a renderer needs.
- The **push-pull operator** `^`, which gives an animation its motion direction.
- The **external object extension**, which Suber explicitly intended for fishhooks among others.

Deliberately not adopted:

- **The string syntax itself.** Actions are stored as structured objects and the Suber string is emitted as a derived `notation` field. Structure diffs readably in a pull request, validates under JSON Schema without a parser, and leaves somewhere to hang rendering data that Suber's notation has no room for. The string remains the citable, human-facing form.
- **The `®` until-operator.** A runtime condition cannot be rendered deterministically. Fishing knots do not need it.
- **`Do(*n)` stage iteration.** Repetition attaches to the action instead, see `repeat` below.

Extended for fishing:

- `repeat` accepts a **range**. Real instructions say "five to seven turns," and flattening that to a single number is a false precision.

Suber's own caveats are worth recording: he states the notation is not adequate for decorative knots, webbing, splicing, or whippings, and that he made it deliberately incomplete to keep it manageable. None of those gaps affect fishing knots.

Note also which layer each verb serves.
`GP` grips with a hand, which generates **prose** but no geometry.
`ML`, `MT`, `RV` and `MB` manipulate cord, which generates **geometry**.
Both are kept, because they feed different consumers.

```yaml
kind: knot
id: knot.palomar
names:
  canonical: Palomar knot
  aliases: [Palomar]
abok_ref: null
connects:
  from: line_end
  to: [closed_eye, split_ring]
line_types: [mono, fluoro, braid]
objects:
  - {ref: HookEye.1, pin_type: closed_eye}
stages:
  - id: 1
    actions:
      - {verb: MB, subject: RP, names: BT.1, length: {value: 6, unit: in}}
    prose: Double roughly six inches of line back on itself to form a bight.
    notation: "* MB(RP=BT.1), LG(BT.1:6in)"
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
    prose: Moisten the line and draw both the standing line and tag end to seat the knot.
    notation: "* MV([SP, E.RP]^)"
strength:
  claims:
    - {residual_pct: 85, line: braid, n: 1, source: src.example, tier: C}
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `kind` | enum | yes | Always `knot`. |
| `id` | string | yes | Namespaced, always prefixed `knot.`. |
| `names.canonical` | string | yes | Preferred display name. Generic where a trademark exists. |
| `names.aliases` | list | no | Regional and colloquial names. Drives search and disambiguation. |
| `abok_ref` | int | no | Ashley Book of Knots reference number. The number may be cited; the text and figures may not be reproduced, since ABoK is in copyright until 2040. |
| `connects.from` | enum | yes | What the knot starts on, almost always `line_end`. |
| `connects.to` | list | yes | Pin types this knot may terminate on. A rig edge tying this knot to a pin of another type is a hard error. |
| `line_types` | list | yes | Line materials the knot is appropriate for. Mismatch is a warning, not an error, since it is a judgement call. |
| `objects` | list | no | Suber user-extension nouns used in the stages, bound to the pin type they stand for. `HookEye.1` here is what lets the same knot record be reused against any component exposing a `closed_eye`. |
| `stages` | list | yes | Ordered. Each stage is one animation keyframe. |
| `stages[].id` | int | yes | Sequential stage number. |
| `stages[].actions` | list | yes | Actions performed **simultaneously** within the stage, per Suber's semantics. |
| `actions[].verb` | enum | yes | Suber verb: `GP`, `MB`, `ML`, `MT`, `MV`, `RV`, `TW`. |
| `actions[].subject` | string or list | yes | What is acted on. A list means a cluster acted on as one object, Suber's `[...]`. |
| `actions[].names` | string | no | Binds a name to what results, Suber's `=` operator. Later stages refer to it. |
| `actions[].through` | string | RV only | What is reeved through. |
| `actions[].around` | string | MT only | What the turn is made around. |
| `actions[].over` | string | MV only | What is passed over. |
| `actions[].direction` | string | no | Six-way direction, or a reeve pair like `F-A`. |
| `actions[].rotation` | enum | no | `CW` or `CCW`. |
| `actions[].chirality` | enum | ML only | `/` or `\`. Determines which strand passes over at the crossing, and is invariant under rotation, which is why it is the right primitive for a renderer. |
| `actions[].force` | enum | no | `push` or `pull`, Suber's `^`. Supplies the animation its motion direction. |
| `actions[].repeat` | int or range | no | Number of times. A range such as `[5, 7]` is legal and preferred where real instructions give one. |
| `actions[].wet` | bool | no | Whether the line should be moistened. Matters enough in practice to be structured rather than buried in prose. |
| `stages[].prose` | string | yes | Human-readable instruction. **Authored, not generated.** See the open question below. |
| `stages[].notation` | string | derived | Suber string emitted from the structured actions. Never hand-edited; CI regenerates and diffs it. |
| `strength.claims` | list | no | Attributed strength claims. Never a bare number. |
| `claims[].residual_pct` | number | yes | Percentage of line strength retained. |
| `claims[].n` | int | yes | Sample size. **Required.** Most published knot tests are n=1 and readers deserve to know. |
| `claims[].source` | ref | yes | Citation. |
| `claims[].tier` | enum | yes | Always `C` for strength. No strength claim is machine-provable. |

## Rigs

A rig is a **graph**, not a list.
Two edge kinds carry most of the weight: `tied` is a fixed connection made by a knot through a named pin, and `threaded` means a component rides on a line segment and can move along it.

### Carolina rig, a chain

```yaml
kind: rig
id: rig.carolina
names:
  canonical: Carolina rig
  aliases: [C-rig]
nodes:
  - {id: main,   type: line, role: main_line}
  - {id: sinker, ref: weight.egg}
  - {id: bead,   ref: bead.glass}
  - {id: swivel, ref: swivel.barrel}
  - {id: leader, type: line, role: leader, length_in: [12, 48]}
  - {id: hook,   ref: hook.ewg}
edges:
  - from: main
    to: sinker
    rel: threaded
    travel: {toward_rod: open, toward_terminal: swivel}
  - {from: main,   to: bead,   rel: threaded, travel: {toward_terminal: swivel}}
  - {from: main,   to: swivel, rel: tied, knot: knot.palomar, pin: eye.a}
  - {from: swivel, to: leader, rel: tied, knot: knot.palomar, pin: eye.b}
  - {from: leader, to: hook,   rel: tied, knot: knot.palomar, pin: eye}
```

`travel.toward_rod: open` is correct and intentional.
The sinker is supposed to slide freely up the main line, which is the entire point of the rig.
It is the **terminal** side that must be stopped, or the weight leaves.
That asymmetry is why the validation rule has to be directional rather than simply asking whether a component is bounded.

### Umbrella rig, a tree

This is the specimen that kills the "ordered list of components along a line" model.

```yaml
kind: rig
id: rig.umbrella
names:
  canonical: Umbrella rig
  aliases: [Alabama Rig, A-Rig]
  trademark_note: >
    "Alabama Rig" and "A-Rig" are registered marks of Slick Lures, LLC.
    Canonical name stays generic; the marks are recorded as aliases only.
nodes:
  - {id: main, ref: null, type: line, role: main_line}
  - {id: head, ref: hardware.umbrella_head, arms: 5}
  - {id: arm1, ref: hardware.wire_arm}
  - {id: arm2, ref: hardware.wire_arm}
  - {id: jig1, ref: jighead.swimbait}
  - {id: jig2, ref: jighead.swimbait}
edges:
  - {from: main, to: head, rel: tied,    knot: knot.palomar, pin: eye}
  - {from: head, to: arm1, rel: fixed,   pin: arm_socket.1}
  - {from: head, to: arm2, rel: fixed,   pin: arm_socket.2}
  - {from: arm1, to: jig1, rel: clipped, pin: snap.1}
  - {from: arm2, to: jig2, rel: clipped, pin: snap.2}
```

Same schema, no changes required.
A list model would have needed a rewrite.

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `kind` | enum | yes | Always `rig`. |
| `names.trademark_note` | string | no | Present whenever any alias is a registered mark. Its presence is what CI checks to confirm the canonical name is generic. |
| `nodes[].id` | string | yes | Local identifier, unique within the rig. Referenced by edges. |
| `nodes[].ref` | ref | no | The component this node instantiates. Absent for line segments. |
| `nodes[].type` | enum | no | `line` for a cord segment. Everything else gets its type from the referenced component. |
| `nodes[].role` | enum | line only | `main_line` or `leader`. Determines which end is free. |
| `nodes[].length_in` | range | no | Length range in inches. A range because real instructions give one. |
| `edges[].from` / `to` | node id | yes | Direction runs rod-ward to terminal-ward, which is what makes `travel` interpretable. |
| `edges[].rel` | enum | yes | `tied`, `threaded`, `fixed`, `clipped`, `crimped`. |
| `edges[].knot` | ref | tied only | Which knot makes the connection. |
| `edges[].pin` | pin id | yes except threaded | Which pin on the target component. Its type must appear in the knot's `connects.to`. |
| `edges[].travel` | object | threaded only | Movement bounds along the segment. |
| `travel.toward_rod` | node id or `open` | yes | What stops the component on the rod side, or `open` if nothing does. |
| `travel.toward_terminal` | node id or `open` | yes | Same, terminal side. `open` on a free end is a hard error. |

## Techniques and conditions

A technique is a separate record because the relationship is many to many.
Drop shot with a shaking retrieve and drop shot with a dragging retrieve are different answers with different conditions.

**Conditions attach to the rig-and-technique pairing, never to the rig.**
Hang them on the rig and the distinction is lost before it is ever noticed.

```yaml
kind: technique
id: technique.carolina.drag
rig: rig.carolina
motion:
  primary: drag
  cadence: slow sweep, pause two to five seconds
  contact: continuous bottom
water_column: bottom
applicability:
  species: [largemouth, smallmouth]
  cover: [gravel, sand, sparse_grass]
  clarity: [stained, clear]
  depth_ft: [8, 25]
  season: [prespawn, postspawn, summer]
  tier: C
  sources: [src.example]
```

| Field | Type | Required | Meaning |
| --- | --- | --- | --- |
| `kind` | enum | yes | Always `technique`. |
| `rig` | ref | yes | The rig this technique applies to. Many techniques may reference one rig. |
| `motion.primary` | enum | yes | `drag`, `hop`, `shake`, `deadstick`, `yoyo`, `swim`, `rip`. Drives the animated lure path. |
| `motion.cadence` | string | no | Rhythm, in prose. |
| `motion.contact` | enum | no | Where the rig rides: `bottom`, `suspended`, `surface`. |
| `water_column` | enum | yes | Depth band, used for the side-view animation. |
| `applicability` | object | yes | The conditions edge. Everything under here is judgement and carries `tier: C`. |
| `applicability.depth_ft` | range | no | A range, never a point value. |
| `applicability.tier` | enum | yes | Effectively always `C`. |
| `applicability.sources` | list | yes | At least one, enforced. |

## Rules that fall out

Each maps to an ERC-style severity.

| # | Rule | Severity | Tier |
| --- | --- | --- | --- |
| 1 | Every `threaded` node must have a `blocks_passage` stop between it and each **free** end of its segment | error | A |
| 2 | A knot's `connects.to` must include the type of the pin it ties to | error | A |
| 3 | `bore_mm` of a threaded component must exceed the diameter of what passes through it | error | A |
| 4 | The rig graph must be connected and acyclic, with no orphan nodes | error | A |
| 5 | Every `ref` must resolve, and every `pin` must exist on its component | error | A |
| 6 | `id` prefix must be an allowed category, and `kind` must match the directory | error | A |
| 7 | An emitted `notation` string must match its structured actions | error | A |
| 8 | Any `tier: C` claim must carry at least one source | error | A |
| 9 | A strength claim without `n` is rejected | error | A |
| 10 | Any alias that is a registered mark requires `trademark_note` | error | A |
| 11 | A knot's `line_types` should include the line type of the segment it ties on | warning | B |
| 12 | Size labels are never compared numerically across manufacturers | error | A |

Rule 1 is the one worth noticing.
"Sliding weight with nothing to stop it" is the most common real rigging error, and against this shape it is a decidable interval check rather than a matter of opinion.

## Open questions

- **Is `prose` authored or generated?** Generating it from Suber verbs would guarantee it matches the animation, which is attractive, since divergence between the words and the picture is exactly the Tier C gap the README describes. But generated prose reads like a robot, and this is a teaching site. Currently authored, with CI checking only that stage counts agree.
- **Where does line diameter live?** On the node, or inherited from a rig-level line spec the reader selects? Inheriting makes the rig reusable across line weights, which is probably right, but rule 3 needs a concrete number to check against.
- **Does `rel: rigged` (bait onto hook) belong in the rig graph at all**, or is it presentation?
- **Do handedness variants come free** as a transform over `chirality` and direction fields, or do some knots need explicit left and right encodings?
- **Suber's notation carries a bare copyright notice with no license grant.** The system is implementable under the idea and expression distinction, but adopting his vocabulary wholesale warrants asking him directly.
