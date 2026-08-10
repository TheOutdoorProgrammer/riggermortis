# Schema prototype

Exploratory shapes, not settled.
The point of this document is to find where the model breaks before any generator exists.
Everything here is written to be argued with.

Field names are provisional.
The knot action vocabulary in particular is a placeholder: it needs reconciling against [Suber's actual notation](https://legacy.earlham.edu/~peters/knotting/notate.htm) before anything is authored at volume.

## The one idea holding it together

Knots and rigs share a single primitive: the **typed connection point**.

A component exposes pins with types (`closed_eye`, `open_eye`, `shank`, `split_ring`, `line_end`).
A knot declares which pin types it can terminate on.
A rig declares which pin a given knot ties to.

Because both sides speak the same vocabulary, a rig becomes a netlist and validation becomes electrical rule checking.
Design the two schemas independently and they will not compose.

## Components

The reusable atom.
Roughly 50 of these yield 100+ rigs.

```yaml
id: swivel.barrel
name: Barrel swivel
category: connector
pins:
  - {id: eye.a, type: closed_eye, wire_diameter_mm: 0.9}
  - {id: eye.b, type: closed_eye, wire_diameter_mm: 0.9}
blocks_passage: true      # nothing threaded on the line can slide past it
sizes:
  # Size labels are NOT comparable across manufacturers. rating_lb is.
  - {system: swivel_number, label: "7", manufacturer: rosco, rating_lb: 55}
  - {system: swivel_number, label: "10", manufacturer: rosco, rating_lb: 35}
```

```yaml
id: weight.egg
name: Egg sinker
category: weight
mounting: threaded        # goes ON the line, is not tied TO it
bore_mm: 2.4              # the largest line or hardware that can pass through
blocks_passage: false     # line runs through it, so it stops nothing
pins: []                  # a threaded component has no pins at all
sizes:
  - {system: ounce, label: "1/2", grams: 14.2}
```

`blocks_passage`, `bore_mm` and `mounting` exist purely so the validator can reason.
They are not display fields.

## Knots

A knot is a **procedure**, so the record is a sequence of stages rather than a shape.
Stages map one to one onto animation keyframes.

```yaml
id: knot.palomar
names:
  canonical: Palomar knot
  aliases: [Palomar]
abok_ref: null            # cite the number if one exists; never the text or figures
connects:
  from: line_end
  to: [closed_eye, split_ring]
line_types: [mono, fluoro, braid]
strength:
  # Published as attributed claims with sample size, never as bare fact.
  claims:
    - {residual_pct: 85, line: braid, n: 1, source: src.example, tier: C}
stages:
  - {id: 1, action: double,      subject: RP, note: form a bight}
  - {id: 2, action: pass_through, subject: bight, target: HookEye.1, direction: U}
  - {id: 3, action: overhand,    subject: [bight, SP]}
  - {id: 4, action: pass_over,   subject: bight, target: Hook.1}
  - {id: 5, action: dress,       subject: all, note: wet before seating}
  - {id: 6, action: seat,        pull: [SP, tag]}
```

`RP` is the running part, `SP` the standing part.
`HookEye.1` is the named external object, which is the capability no topological formalism provides and the reason this model was chosen over knot theory (see [ADR 0001](../adr/0001-encode-tying-method-not-topology.md)).

## Rigs

A rig is a **graph**, not a list.
Two edge kinds carry most of the weight:

- `tied` is a fixed connection made by a knot, through a named pin.
- `threaded` means the component rides on a line segment and can move.

### Carolina rig, a chain

```yaml
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
  - {from: main, to: bead, rel: threaded, travel: {toward_terminal: swivel}}
  - {from: main,   to: swivel, rel: tied, knot: knot.palomar, pin: eye.a}
  - {from: swivel, to: leader, rel: tied, knot: knot.palomar, pin: eye.b}
  - {from: leader, to: hook,   rel: tied, knot: knot.palomar, pin: eye}
```

`travel.toward_rod: open` is correct and intentional.
The sinker is *supposed* to slide freely up the main line.
It is the **terminal** side that must be stopped, or the weight leaves the rig.
That asymmetry is the whole reason the check has to be directional rather than "is it bounded."

### Umbrella rig, a tree

This is the specimen that kills the "ordered list of components along a line" model that both of us were assuming.

```yaml
id: rig.umbrella
names:
  canonical: Umbrella rig
  aliases: [Alabama Rig, A-Rig]
  # "Alabama Rig" and "A-Rig" are registered marks of Slick Lures, LLC.
  # Canonical name stays generic; the mark is recorded as an alias only.
nodes:
  - {id: main, type: line, role: main_line}
  - {id: head, ref: hardware.umbrella_head, arms: 5}
  - {id: arm1, ref: hardware.wire_arm}
  - {id: arm2, ref: hardware.wire_arm}
  - {id: jig1, ref: jighead.swimbait}
  - {id: jig2, ref: jighead.swimbait}
edges:
  - {from: main, to: head, rel: tied, knot: knot.palomar, pin: eye}
  - {from: head, to: arm1, rel: fixed, pin: arm_socket.1}
  - {from: head, to: arm2, rel: fixed, pin: arm_socket.2}
  - {from: arm1, to: jig1, rel: clipped, pin: snap.1}
  - {from: arm2, to: jig2, rel: clipped, pin: snap.2}
```

Same schema, no changes required.
A list model would have needed a rewrite.

## Techniques and conditions

A technique is a separate record because the relationship is many to many.
Drop shot plus shaking and drop shot plus dragging are different answers with different conditions.

**Conditions attach to the rig-and-technique pairing, never to the rig.**

```yaml
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
  tier: C                 # judgement, not fact; must carry sources
  sources: [src.example]
```

## Rules that fall out of these shapes

Each maps to an ERC-style severity.
These are the Tier A and Tier B checks the README describes, expressed against the fields above.

| # | Rule | Severity | Tier |
| --- | --- | --- | --- |
| 1 | Every `threaded` node must have a `blocks_passage` stop between it and each **free** end of its segment | error | A |
| 2 | A knot's `connects.to` must include the type of the pin it ties to | error | A |
| 3 | `bore_mm` of a threaded component must exceed the diameter of what passes through it | error | A |
| 4 | The rig graph must be acyclic and fully connected, with no orphan nodes | error | A |
| 5 | Every `ref` must resolve to an existing component | error | A |
| 6 | A knot's `line_types` should include the line type of the segment it is tied on | warning | B |
| 7 | Size labels must never be compared numerically across manufacturers | error | A |
| 8 | Any `tier: C` claim must carry at least one source | error | A |

Rule 1 is the one worth noticing.
"Sliding weight with nothing to stop it" is the most common real rigging error, and against this shape it is a decidable interval check rather than a matter of opinion.

## Open questions

- Reconcile the stage action vocabulary against Suber's real notation before authoring at volume.
- Does `rel: rigged` (bait onto hook) belong in the rig graph, or is it a presentation concern?
- Where do line segments get their diameter: on the node, or inherited from a rig-level line spec the user selects?
- Do left-handed and right-handed variants come free as a transform over stages, or do some knots need explicit handedness?
