# Coverage

<!-- SPDX-License-Identifier: Apache-2.0 -->
> **Licence.** This document is **Apache-2.0**.
> The dataset it tracks, under [`data/`](../data/), is **ODbL-1.0**.
> See [LICENSE](../LICENSE), [LICENSE-DATA](../LICENSE-DATA), and [ADR 0005](../adr/0005-dual-licence-odbl-and-apache.md).

What exists, what does not, and what a human has actually confirmed.

This file is hand-maintained today.
Once the schema lands it should be **generated** from the records themselves, since every column below is already a field in the data.
Until then, treat a mismatch between this file and the data as a bug in this file.

## Legend

| Column | Means |
| --- | --- |
| **Data** | A record exists and passes every Tier A rule in CI. |
| **Draw** | The generated diagram or animation has been checked and is correct. |
| **Field** | Someone who fishes it confirmed the record matches reality. This is `validation.status: field-tested` and it is the only thing that closes the Tier C gap. |

☐ not done · ☑ done · ◐ partial · ✕ blocked

**Bold** entries are the v1 core.
Everything else is catalogued so the scope is honest, not because it ships first.

A record can pass CI and still be wrong. **Data ☑ with Field ☐ means "structurally valid, unverified."** That distinction is the entire point of this file.

---

## Knots

### Terminal: line to hardware

| Knot | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Palomar** | | ☑ | ☐ | ☐ | Doubled line through an eye. Core specimen. |
| **Improved clinch** | fisherman's knot | ☐ | ☐ | ☐ | Wraps plus tuck. Tests `repeat` range 5–7. |
| Clinch | | ☐ | ☐ | ☐ | Variant of the above; candidate `variant_of`. |
| Trilene | double-loop clinch | ☐ | ☐ | ☐ | Two passes through the eye. |
| **Uni** | Duncan loop, grinner | ☐ | ☐ | ☐ | Basis for the double uni bend. |
| **Snell** | | ☐ | ☐ | ☐ | Wraps the **shank**, not the eye. Tests `pin_type: shank`. |
| Egg loop snell | bait loop | ☐ | ☐ | ☐ | Snell that traps bait. |
| San Diego jam | reverse clinch | ☐ | ☐ | ☐ | |
| **Non-slip loop** | Kreh loop, Lefty's loop | ☐ | ☐ | ☐ | Loop knot at the lure; tests free-swinging attachment. |
| Rapala knot | | ☐ | ☐ | ☐ | |
| Homer Rhode loop | | ☐ | ☐ | ☐ | |
| Orvis knot | | ☐ | ☐ | ☐ | |
| Berkley braid knot | | ☐ | ☐ | ☐ | Braid-specific; tests `knot_sensitivity`. |
| Pitzen | Eugene bend | ☐ | ☐ | ☐ | |
| Turle knot | | ☐ | ☐ | ☐ | |
| Double Palomar | | ☐ | ☐ | ☐ | |

### Bends: line to line

| Knot | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **FG knot** | | ☐ | ☐ | ☐ | Braid to fluoro, **no hardware**. Repeated weave; the specimen that motivated `patterns`. |
| **Double uni** | | ☑ | ☐ | ☐ | Two unis facing each other. Exercises `patterns` in a knot. |
| **Blood knot** | barrel knot | ☐ | ☐ | ☐ | Similar diameters only. |
| **Surgeon's knot** | double/triple surgeon's | ☐ | ☐ | ☐ | Tests `repeat` on turns. |
| Alberto | crazy Alberto | ☐ | ☐ | ☐ | |
| Albright special | | ☐ | ☐ | ☐ | Joins very unequal diameters. |
| Slim Beauty | | ☐ | ☐ | ☐ | |
| Yucatan knot | | ☐ | ☐ | ☐ | |
| Nail knot | | ☐ | ☐ | ☐ | Needs a tool or tube as an external object. |
| PR bobbin knot | | ☐ | ☐ | ☐ | Requires a bobbin; external-object stress case. |

### Loops and structure: knots as nodes

| Knot | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Dropper loop** | | ☐ | ☐ | ☐ | `role: loop`. Three pins. Powers high-low, Sabiki, crappie, spider. |
| **Surgeon's loop** | double surgeon's loop | ☐ | ☐ | ☐ | |
| Perfection loop | | ☐ | ☐ | ☐ | |
| **Bimini twist** | | ☐ | ☐ | ☐ | 20+ twists. Heaviest `patterns` stress case. |
| Spider hitch | | ☐ | ☐ | ☐ | Faster Bimini substitute. |
| Figure-eight loop | | ☐ | ☐ | ☐ | |

### Stoppers and spool

| Knot | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Bobber stop knot** | | ☐ | ☐ | ☐ | `role: stopper`, `blocks_passage: true`. Makes slip-float rigs validate. |
| Arbor knot | | ☐ | ☐ | ☐ | `role: arbor`. Line to spool. |

*Fly-specific knots (needle knot, nail knot to fly line, etc.) are out of scope for v1.*

---

## Rigs

### Bass and freshwater artificial

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Texas rig** | T-rig | ☐ | ☐ | ☐ | Pegged vs free-sliding bullet weight. |
| **Carolina rig** | C-rig | ☑ | ☐ | ☐ | Sliding weight + stopper. Core specimen. |
| **Drop shot** | | ☐ | ☐ | ☐ | Hook mid-line, tag end below. |
| **Ned rig** | | ☐ | ☐ | ☐ | Simplest possible: jighead + plastic. |
| **Wacky rig** | | ☐ | ☐ | ☐ | Hook through the middle; O-ring variant. |
| Neko rig | | ☐ | ☐ | ☐ | Wacky plus a nail weight. |
| Shaky head | | ☐ | ☐ | ☐ | |
| Jika rig | | ☐ | ☐ | ☐ | Split ring junction. |
| Tokyo rig | | ☐ | ☐ | ☐ | ⚠ Trademark asserted (VMC/Rapala), **unverified**. Needs `trademark_note`. |
| Free rig | | ☐ | ☐ | ☐ | |
| Punch rig | | ☐ | ☐ | ☐ | |
| Split shot rig | | ☐ | ☐ | ☐ | Crimped weight; unusual mounting. |
| Mojo rig | | ☐ | ☐ | ☐ | |
| **Umbrella rig** | Alabama Rig®, A-Rig® | ☐ | ☐ | ☐ | ⚠ Registered marks (Slick Lures LLC). **Requires `legality` block** for hook limits. Tree topology. |
| Damiki rig | | ☐ | ☐ | ☐ | |
| Flick shake | | ☐ | ☐ | ☐ | |
| Weightless / Texposed | | ☐ | ☐ | ☐ | |
| Scrounger | | ☐ | ☐ | ☐ | |

### Live and cut bait, freshwater

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Slip bobber rig** | slip float | ☐ | ☐ | ☐ | Bobber stop as a stopper node. Strict ordering. |
| Fixed bobber rig | | ☐ | ☐ | ☐ | |
| **Three-way rig** | Wolf River rig | ☐ | ☐ | ☐ | Three-pin junction. |
| Lindy rig | walleye slip sinker | ☐ | ☐ | ☐ | |
| Bottom bouncer | | ☐ | ☐ | ☐ | |
| Spinner harness | crawler harness | ☐ | ☐ | ☐ | Beads + clevis + blade in series. |
| Santee Cooper rig | | ☐ | ☐ | ☐ | Catfish; peg float on the leader. |
| Quick-strike rig | | ☐ | ☐ | ☐ | Pike; two trebles. Check hook limits. |
| Paternoster | | ☐ | ☐ | ☐ | Naming overlaps high-low regionally. |

### Saltwater inshore

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Fish finder rig** | slider rig | ☐ | ☐ | ☐ | Sliding sinker on the main line. |
| **Popping cork rig** | | ☐ | ☐ | ☐ | Cork with beads on a wire. |
| Knocker rig | | ☐ | ☐ | ☐ | Weight rides directly on the hook. |
| **High-low rig** | chicken rig, double drop | ☐ | ☐ | ☐ | Two dropper loops. Killed "knot is always an edge." |
| Free line rig | | ☐ | ☐ | ☐ | No weight at all; degenerate case worth testing. |

### Surf

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| Pompano rig | | ☐ | ☐ | ☐ | Floats + droppers. |
| Double drop bottom rig | | ☐ | ☐ | ☐ | Regional twin of high-low. |
| Sputnik / storm sinker rig | | ☐ | ☐ | ☐ | Anchoring sinker with wire arms. |
| Fireball rig | | ☐ | ☐ | ☐ | |

### Offshore and trolling

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| Ballyhoo rig | naked / skirted | ☐ | ☐ | ☐ | Bait rigging, not just tackle. Boundary question. |
| **Spreader bar** | | ☐ | ☐ | ☐ | Heavy `patterns` case; N teasers. |
| Daisy chain | | ☐ | ☐ | ☐ | Baits in series; `patterns`. |
| Deep drop rig | | ☐ | ☐ | ☐ | Many droppers. Check hook limits. |
| Kite fishing rig | | ☐ | ☐ | ☐ | |
| Downrigger release rig | | ☐ | ☐ | ☐ | External object not on the line. |
| Planer rig | | ☐ | ☐ | ☐ | |
| Wire leader rig | | ☐ | ☐ | ☐ | Tests `line_material: wire` and crimps over knots. |

### Panfish and crappie

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Crappie double-minnow rig** | | ☐ | ☐ | ☐ | Two droppers. |
| Kentucky rig | | ☐ | ☐ | ☐ | |
| Spider rig | | ☐ | ☐ | ☐ | Multiple rods; may not be a rig at all. Boundary question. |
| Roadrunner / underspin | | ☐ | ☐ | ☐ | Assembled or bought? Boundary rule applies. |

### Catfish

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Slip sinker catfish rig** | | ☐ | ☐ | ☐ | |
| Slip float catfish rig | | ☐ | ☐ | ☐ | |
| Three-way catfish rig | | ☐ | ☐ | ☐ | |

### Ice

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| Tip-up rig | | ☐ | ☐ | ☐ | |
| Jigging spoon with dropper | | ☐ | ☐ | ☐ | |
| Deadstick rig | | ☐ | ☐ | ☐ | |

### Bait catching

| Rig | Also known as | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- | --- |
| **Sabiki rig** | bait rig, piscator | ☐ | ☐ | ☐ | Six identical droppers. The `patterns` motivating case. Check hook limits. |

*Fly rigs (dry-dropper, nymph, streamer) are out of scope for v1.*

---

## Techniques

| Technique | Data | Draw | Field | Notes |
| --- | --- | --- | --- | --- |
| **Drag** | ☐ | ☐ | ☐ | Continuous bottom contact. |
| **Hop** | ☐ | ☐ | ☐ | |
| **Shake** | ☐ | ☐ | ☐ | In-place; hardest to animate legibly. |
| **Deadstick** | ☐ | ☐ | ☐ | Degenerate case: no motion at all. |
| Yo-yo | ☐ | ☐ | ☐ | |
| **Swim** | ☐ | ☐ | ☐ | |
| Rip / snap | ☐ | ☐ | ☐ | |
| Walk the dog | ☐ | ☐ | ☐ | Surface only. |
| Twitch and pause | ☐ | ☐ | ☐ | |
| Pitch / flip | ☐ | ☐ | ☐ | Delivery, not retrieve. Boundary question. |
| Punch | ☐ | ☐ | ☐ | |
| Slow roll | ☐ | ☐ | ☐ | |
| Burn | ☐ | ☐ | ☐ | |
| Stroll / power-shot | ☐ | ☐ | ☐ | |
| **Vertical jig** | ☐ | ☐ | ☐ | |
| Troll | ☐ | ☐ | ☐ | Boat motion, not rod motion. |
| Drift | ☐ | ☐ | ☐ | |
| Dead drift | ☐ | ☐ | ☐ | Current-driven. |
| Anchor and wait | ☐ | ☐ | ☐ | |

---

## Components

Roughly fifty of these unlock everything above.
This is the highest-leverage block in the entire project and should be authored first.

| Group | Members | Data | Draw | Notes |
| --- | --- | --- | --- | --- |
| **Hooks** | EWG, offset worm, straight shank, octopus, circle, treble, drop shot, wacky, weedless wacky, aberdeen, kahle, siwash, live bait, baitholder, sabiki | ☐ | ☐ | Sizing is per-manufacturer **and** per-pattern. |
| **Weights** | bullet, egg, split shot, drop shot cylinder, drop shot teardrop, bank, pyramid, sputnik, no-roll, walking, tungsten nail, mojo, keel | ☐ | ☐ | `mounting` and `bore_mm` decide sliding behaviour. |
| **Connectors** | barrel swivel, ball bearing swivel, three-way swivel, snap swivel, snap, split ring, coastlock | ☐ | ☐ | `blocks_passage: true` for most. Split ring ratings are chaotic. |
| **Floats** | fixed round, slip float, popping cork, pencil, cigar, foam peg | ☐ | ☐ | |
| **Beads and stops** | glass bead, plastic bead, faceted bead, bobber stop, peg | ☐ | ☐ | |
| **Blades and skirts** | willow, Colorado, Indiana, clevis, skirt, rattle | ☐ | ☐ | |
| **Jigheads** | ball, football, swimbait, shaky, ned, tube | ☐ | ☐ | Hook plus weight in one; single component. |
| **Hardware** | umbrella head, wire arm, leader sleeve, crimp | ☐ | ☐ | |

---

## Lines

| Line | Data | Field | Notes |
| --- | --- | --- | --- |
| **Monofilament** | ☐ | ☐ | |
| **Fluorocarbon** | ☐ | ☐ | `knot_sensitivity: high`. |
| **Braid** | ☐ | ☐ | `knot_sensitivity: high` for a different reason. |
| Wire | ☐ | ☐ | Crimped, not tied. |
| Leader material | ☐ | ☐ | May be a variant of mono or fluoro rather than its own record. |

Breaking load per stated diameter is **not** standardised, so every `breaking_load_n` to `diameter_mm` pairing is a sourced measurement rather than a property of the material.

---

## Species

Only species actually referenced by a technique need records.
Listing them here prevents bare strings creeping back into `applicability`.

| Species | Data | Notes |
| --- | --- | --- |
| **Largemouth bass** | ☐ | ⚠ Black bass taxonomy was revised recently. Scientific name **must** be sourced, not recalled. |
| **Smallmouth bass** | ☐ | Same caveat. |
| Spotted bass | ☐ | Same caveat. |
| **Bluegill** | ☐ | Ambiguous alias: "sunfish". |
| **Crappie** (black, white) | ☐ | Regionally called "sunfish" and "speckled perch". |
| **Walleye** | ☐ | |
| **Channel / blue / flathead catfish** | ☐ | |
| Northern pike | ☐ | |
| **Redfish** | ☐ | Aliases: red drum, channel bass, spottail. |
| **Spotted seatrout** | ☐ | Aliases: speckled trout, speck. Not a trout. |
| Snook | ☐ | |
| Flounder | ☐ | Several unrelated species share the name. |
| Striped bass | ☐ | Aliases: striper, rockfish. ⚠ "Rockfish" means an entirely different family on the Pacific coast. |
| Sheepshead | ☐ | ⚠ Means two unrelated species depending on region. The canonical `ambiguous_aliases` case. |

---

## Summary

| Category | Catalogued | Data | Draw | Field |
| --- | --- | --- | --- | --- |
| Knots | 34 | 0 | 0 | 0 |
| Rigs | 47 | 0 | 0 | 0 |
| Techniques | 19 | 0 | 0 | 0 |
| Component groups | 8 | 0 | 0 | n/a |
| Lines | 5 | 0 | n/a | 0 |
| Species | 14 | 0 | n/a | n/a |

Nothing is built.
That is the accurate state and this table should never flatter it.
