# Knot Data Research

Research into data sources and formal notation systems for fishing knots, for `riggermortis` (public, open-source, renders knot-tying animations from structured data).

Every claim carries a URL.
Anything not directly verified is marked **UNCERTAIN**.
Where a prior research run's conclusion was wrong, it is called out explicitly.

## Two headline findings, read these first

**1. The Ashley Book of Knots is NOT in the public domain.**
Its US copyright was renewed on 1971-08-02 (registration R510678).
Primary evidence in [§ Licensing Verdicts](#5-licensing-verdicts).
This reverses a prior run's conclusion and a widely-repeated claim online.

**2. The notation this project actually wants already exists, and it is not from knot theory.**
[Peter Suber's Knot Tying Notation](https://legacy.earlham.edu/~peters/knotting/notate.htm) (2002-2004, v0.9) encodes *methods of tying*, distinguishes running part from standing part, handles tying **around a named external object**, and is organised into **stages** that map one-to-one onto animation keyframes.
Detail in [§ Notation Systems](#2-notation-systems).
Knot theory is the wrong layer for this project; Suber is the right one.

## 1. The Open Knot Problem

### The actual problem

Classical knot theory studies **closed loops**: embeddings of S¹ in S³ or R³.
Every open arc in R³ is ambient-isotopic to a straight line, so classically *every fishing knot is the unknot*.
That is not a technicality to wave away.
It is the reason no off-the-shelf knot invariant will ever tell an Improved Clinch from a Palomar.

A fishing knot breaks three different assumptions at once:

1. **Open ends.** There is a working (tag) end and a standing part. Not a loop.
2. **Tied around an object.** The hook eye, swivel, or a second line is a distinct topological object the line passes through and around. The interesting structure lives in the relationship between line and object, not in the line alone.
3. **Held by friction, not topology.** A clinch knot's holding power is mechanical and tribological. A topologically identical arrangement in slippery braid fails. Topology cannot express this at all.

### Closure schemes (the mainstream workaround)

The protein-knot community hit problem (1) first, because a protein backbone is an open chain, and solved it by *closing* the chain and then applying classical invariants.

| Scheme | What it does | Source |
| --- | --- | --- |
| **Direct closure** | Join the two endpoints by the shortest segment (a straight chord). Deterministic, one answer, but sensitive to endpoint geometry and the chord can cut back through the structure. | [Knoto-ID, Bioinformatics 34(19):3402](https://academic.oup.com/bioinformatics/article/34/19/3402/4990827) |
| **Single-ended / outward closure** | Extend each end radially outward to a large enclosing sphere, then close along the sphere surface. Reduces the chord artefact. | [KnotProt knot detection](https://knotprot.cent.uw.edu.pl/knot_detection) |
| **Stochastic / probabilistic (random) closure** | Close the chain hundreds of times to random points on a large enclosing sphere, compute the invariant each time, and report a **probability distribution over knot types**. KnotProt connects endpoints "several hundred times to two points randomly chosen from a set of vertices of the truncated icosahedron positioned on a large sphere enclosing the analyzed chain", then joins them by an arc on the sphere surface. | [KnotProt knot detection](https://knotprot.cent.uw.edu.pl/knot_detection) |
| **Uniform closure** | The deterministic-average variant Knoto-ID exposes alongside direct closure. | [Knoto-ID paper](https://academic.oup.com/bioinformatics/article/34/19/3402/4990827) |

The honest read: closure turns "what knot is this open chain?" into "what knot does this open chain *most probably* become?"
That is a statistical answer to a geometric question.
Fine for classifying tens of thousands of PDB structures.
**Useless for us**, because we already know which knot we mean.
We need to *describe* it, not *identify* it.

### Knotoids (the formalism that does not require closure)

**Knotoids** were introduced by **Vladimir Turaev in 2012** as a generalisation of knots to open curves.

- A knotoid diagram is *"an immersed oriented interval in a surface, with classical crossing information at double points, considered up to Reidemeister moves performed away from the endpoints"*, per [Algebraic and Geometric Aspects of Non-Classical Knots, arXiv:2607.04445](https://arxiv.org/pdf/2607.04445).
- The critical rule: **endpoints cannot be pulled across strands**. These are the *forbidden endpoint moves*. That restriction is exactly what preserves the information a closure destroys (same source).
- Knotoids generalise the **1-1 tangle** by allowing the two endpoints to sit in *different regions* of the diagram (same source).
- Applied to proteins: *"The knotoid approach doesn't require the closure of chains into loops which implies that the geometry of analysed chains does not need to be changed by closure"*, and it *"detects topologically nontrivial protein folds that are not detected using the knotting approach"*, per [Goundaroulis et al., arXiv:1705.07849](https://arxiv.org/pdf/1705.07849) / [PubMed 28740166](https://pubmed.ncbi.nlm.nih.gov/28740166/).
- Invariants used: the **Jones polynomial** for projections to a sphere, or the **Turaev loop bracket polynomial** for projections to a plane ([Knoto-ID paper](https://academic.oup.com/bioinformatics/article/34/19/3402/4990827)).
- A systematic classification of planar and spherical knotoids exists: [arXiv:1902.07277](https://arxiv.org/pdf/1902.07277).

### Implementations

| Tool | What it does | License | URL |
| --- | --- | --- | --- |
| **Knoto-ID** | Reference implementation of the knotoid formalism. Classifies open curves without closure; also implements direct and uniform closure. Global topology, exhaustive subchain analysis, knotted-core identification. Input: text files of Cα coordinates from PDB. C++ with R for viz. | Paper states **GPL v2 or later**. GitHub license auto-detection returns **null** (no `LICENSE` file detected by the API), so the in-repo grant is **UNCERTAIN**. Last push 2019-03-11, 9 stars. | [github.com/sib-swiss/Knoto-ID](https://github.com/sib-swiss/Knoto-ID) |
| **pyknotid `OpenKnot`** | Python. Stochastic-closure engine plus virtual-knot detection. Verified in detail below. | **MIT** (confirmed via GitHub API on `SPOCKnots/pyknotid`, 47 stars, last push 2023-02-19) | [github.com/SPOCKnots/pyknotid](https://github.com/SPOCKnots/pyknotid) / [docs](https://pyknotid.readthedocs.io/en/latest/sources/spacecurves/openknot.html) |
| **KnotProt** | Database and web service using stochastic closure over the PDB. | Web service | [knotprot.cent.uw.edu.pl](https://knotprot.cent.uw.edu.pl/knot_detection) |

**pyknotid `OpenKnot` is real, and the prior run was right that it exists.**
It subclasses `SpaceCurve`, holds the vertices of a single open line, and replaces the closed-curve defaults.
Confirmed methods, per the [class docs](https://pyknotid.readthedocs.io/en/latest/sources/spacecurves/openknot.html):

- `alexander_polynomials()`, closes on a sphere of given radius at N approximately-evenly-distributed sample points. This is spherical stochastic closure.
- `closure_alexander_polynomial(theta, phi)`, invariant for one specific closure direction.
- `alexander_fractions()`, the **knot spectrum**: each polynomial with the *fraction* of closure directions producing it.
- `virtual_check()` / `virtual_checks()` / `virtual_fractions()`, checks whether the projection's **Gauss code** corresponds to a virtual knot; `projection_invariant()` computes a virtual invariant if so.
- `self_linking()` / `self_linkings()` / `self_linking_fractions()`, the self-linking number J(K), a knotoid-adjacent invariant.
- `plot_alexander_shell()`, plots the curve inside a translucent sphere **coloured by the knot type obtained by closing at each point**. The clearest possible picture of what closure costs you.
- `raw_crossings()`, crossings computed **without** the closing segment. This is the honest open-curve crossing set and the one genuinely reusable primitive here.
- `arclength()` (excludes closure), `closing_distance()`, `vassiliev_degree_2_average()`.

Note what it *is*: a **classifier**, not a representation.
It answers "which knot is this curve probably equivalent to". It does not store how to tie anything.

### Does the around-an-object part break it? Yes, and this is where I want to be blunt

- Knotoids are defined as an **immersed interval in a surface** ([arXiv:2607.04445](https://arxiv.org/pdf/2607.04445)). One strand. There is no second object in the definition.
- The natural extension exists but moves the goalposts. Knotoids and braidoids have been developed **on the torus** ([arXiv:2103.16433](https://arxiv.org/pdf/2103.16433)). A hook eye is topologically a solid torus, so "knotoid in the complement of a solid torus" is the formally correct object for a fishing knot. **UNCERTAIN:** I found no worked-out theory, and definitely no software, for knotoids in a solid-torus complement applied to practical knots. The torus paper is about knotoids *on* the torus surface, a related but different construction.
- Multi-component fishing knots are not knotoids at all. A bend joining two lines (blood knot, Albright, FG) is a **tangle** or a **spatial graph**, not a single interval.
- Nothing in any of this expresses **friction, line diameter, dressing, seating order, or tightening sequence**. For a fishing knot those are not decoration, they are the entire point. Whether a Palomar or a Trilene Clinch slips at 80% of breaking strain is topologically invisible.

### Architectural conclusion

**Do not build `riggermortis` on knot theory.** It is the wrong abstraction layer:

- It answers "are these two knots the same?" We need "what does the line do next?"
- Its canonical forms deliberately discard exactly what an animation needs: which end moves, in what order, around what.
- Its ambient spaces do not contain a hook.

Use knot theory **only** as an optional validation and dedup layer.
A knotoid or closure-spectrum fingerprint computed from generated 3D geometry would let you detect that two authored knots are secretly the same, and `pyknotid` (MIT) is a reasonable tool for that.
That is a nice-to-have, not a foundation.

The foundation should be a **procedural notation**: an ordered sequence of operations on a strand relative to named objects. See the next section, because that already exists.

## 2. Notation Systems

### Mathematical notations, and whether any can encode a working knot

| Notation | Encodes | Usable for a fishing knot? |
| --- | --- | --- |
| **Gauss code** | Sequence of crossings encountered walking the curve, with over/under and sign. | **Partially.** Defined for closed curves, but `pyknotid` computes it for open-curve projections, and an open Gauss code is exactly the combinatorial core of a knotoid diagram. Captures crossing order and over/under, which is a real part of what a tying sequence needs. Does **not** capture the hook eye or which end is the working end. |
| **PD (planar diagram) code** | Set of crossings, each listing its four incident edge labels in cyclic order. Used by the Knot Atlas and `KnotTheory\``. | Closed-diagram oriented. Would need extension for free ends and for a non-line object. **UNCERTAIN** whether a standard open-PD extension exists. |
| **DT (Dowker-Thistlethwaite) code** | A sequence of even integers, sign giving over/under. Extremely compact. | **No.** Only valid for *prime* knots and relies on a closed traversal returning to its start. [Wikipedia: DT notation](https://en.wikipedia.org/wiki/Dowker%E2%80%93Thistlethwaite_notation) |
| **Conway notation** | Rational-tangle arithmetic on a closed diagram. | **No** for our purposes, though the *tangle* substrate is conceptually the closest classical object to "tied around something". [Knot Atlas: Conway Notation](https://katlas.org/wiki/Conway_Notation) |
| **Braid word** | Word in braid generators σᵢ; the closure of the braid gives a link. | **Partially interesting.** Braid words are genuinely sequential and read a bit like instructions ("strand i crosses over strand i+1"). For the wrap phase of a clinch or a Bimini twist, a braid-word-ish encoding of N wraps is plausible. But braids assume all strands run monotonically in one direction and are closed at the end. |

**Bottom line:** none of the classical notations is a drop-in.
The closest usable primitive is an **open/extended Gauss code**, which is also the combinatorial data the knotoid literature uses.

### Peter Suber's Knot Tying Notation (the real find)

[legacy.earlham.edu/~peters/knotting/notate.htm](https://legacy.earlham.edu/~peters/knotting/notate.htm), Peter Suber (Earlham College, later Harvard Open Access Project), v0.9, first online 2002-07-26, last revised 2004-07-28.

Suber states the distinction outright:

> "Mathematical knot theory has several systems of notation for describing the structure of knots already tied. By contrast, the present notation describes *methods of tying knots*."

That single sentence is the architectural thesis of this project.

**Structure.** Predicate-term syntax, `Predicate(term, term, ...)`.
Sentences are grouped into **stages** delimited by `*`.
All actions within a stage happen simultaneously; stages sequence the process in time.
**Stages are animation keyframes.** That mapping is free.

**Nouns (cord parts).** `C` cord, `RP` running part, `SP` standing part, `E` end, `LP` loop, `BT` bight, `CR` crossing, `CS` cord segment, `T` turn, `K` knot, `EY` eye, `LE` link eye, `ST` site, `SD` side. Hands: `LH`/`RH`, `LF`/`RF` with fingers 1-5. Unit `u` = one cord diameter (case-sensitive against direction `U`).

**Directions.** `L`/`R`, `U`/`D`, `F`/`A` (fore/aft), `CW`/`CCW`. Planes `HP`/`VP`/`EP`. Lines `LR`/`UD`/`FA`.

**Verbs.** `GP` grip, `ML` make loop (with `/` and `\` modifiers distinguishing crossing sense), `MB` make bight, `MT` make turn (wrap around anchor), `MV` move, `RV` reeve (thread through an opening), `TW` twist, `PT` point, `Do` universal verb.

**Descriptors (state predicates, no action implied).** `L()`/`R()`/`F()`/`A()`/`U()`/`D()` relative position, `BN` between, `CO` crossover, `LG` length, `PN` plane, `NX` next, `OL` overlay, `PL` parallel, `NM` name, `RS` results.

**Modifiers.** `.` which-one, `:` how, `=` naming, `^` push/pull, `-` direction change, `~` negation/release, `[...]` clusters.

**It handles tying around an object.** The `MT` verb binds around a named anchor, and external objects enter by user extension:

```text
MT(RP, Tree.1:CW:U):2
```

= "make two turns with the running part around tree #1, clockwise, viewed from above".

Substitute `HookEye.1` and you have the exact primitive fishing knots need. The document explicitly demonstrates naming an object not predefined in the notation (a marlinspike) and notating actions performed by tools.

**Working end vs standing part is fundamental**, not an afterthought: `RP` is the active section, `SP` the load-bearing one, `E.RP` the running end specifically.

**Worked example, the overhand knot:**

```text
* ML(/RP=LP.1)
* RV(E.RP, LP.1:A-F)
* MV(E.RP^:R, E.SP^:L)
```

Stage 1: make a forward-slash loop in the running part, name it LP.1.
Stage 2: reeve the running end through LP.1, aft to fore.
Stage 3: pull the running end right and the standing end left (tighten).

**Limits, stated by the author.** No formal BNF; syntax is described narratively.
Suber explicitly chose "manageability over adequacy": *"I'd rather have a fairly simple and learnable notation adequate for about 80% of useful tying methods than an extremely complex notation that is adequate for 98%."*
Scope covers practical knots, hitches, bends, lashings and rigging; it explicitly excludes decorative knotwork, splicing, and webbing techniques.
There is **no 3D geometry**: the notation says what to do, not where the curve is in space. A renderer still has to solve the geometry.

**License: UNCERTAIN, needs resolution.** The page carries a bare `Copyright © 2002-2006, Peter Suber.` with **no license grant**. See [§ Licensing Verdicts](#5-licensing-verdicts) for what that does and does not block.

### Fink and Mao necktie grammar (the same idea, formalised harder)

Fink and Mao's necktie-knot work established a formal alphabet where a knot is a **sequence of moves**, each saying which way the active blade drapes over the partially-built knot, with `U` meaning "tuck the blade under". Extended in "More ties than we thought" ([arXiv:1401.8242](https://arxiv.org/pdf/1401.8242), published open access as [PeerJ Computer Science, peerj.com/articles/cs-2/](https://peerj.com/articles/cs-2/); the PeerJ page returned HTTP 403 to automated fetch, so its exact license string is **UNCERTAIN**, though PeerJ CS is CC BY by default).
Plain-language write-up by one of the authors: [Vejdemo-Johansson, "Necktie knots, formal languages and network security"](https://medium.com/message/necktie-knots-formal-languages-and-network-security-2f703632a527).
Related: [The Mathematics of Tie Knots, arXiv:2005.13000](https://arxiv.org/pdf/2005.13000).

Relevance: it is proof that a **regular-language encoding of tying moves** works and is enumerable. It is domain-specific to neckties (fixed anchor = the neck, fixed move alphabet) and is not directly reusable, but the design pattern (alphabet of moves + grammar + generator) transfers cleanly.

### The IGKT does not have a notation standard

Searched for an International Guild of Knot Tyers notation standard.
**Negative result.** The IGKT ([igkt.net](https://igkt.net/), founded 1982) publishes knot charts and educational material, not a formal encoding.
There is a forum thread literally titled ["Is Knot Notation still being pursued?"](https://forum.igkt.net/index.php?topic=6930.0), which is itself a signal that the community tried and it stalled. That host would not resolve from this machine (`getaddrinfo ENOTFOUND forum.igkt.net`) so the thread contents are **UNVERIFIED**.
The de-facto "notation" for practical knots is the **ABoK number**, which is an arbitrary catalogue index, not a structural encoding.

### Mathesis `.knt` / `.dgr` (chased, confirmed, and more useful than expected)

[github.com/Mathesis-Software/Knots](https://github.com/Mathesis-Software/Knots) ("KnotEditor"), C++/Qt6, **Apache-2.0** ([LICENSE](https://raw.githubusercontent.com/Mathesis-Software/Knots/master/LICENSE)), 19 stars.
Draws 2D diagrams, converts them to 3D knots, computes Seifert surfaces and invariants.

Both formats are plain **JSON**. Verified by direct fetch of the sample data.

`.dgr`, the 2D diagram ([3_1.dgr](https://raw.githubusercontent.com/Mathesis-Software/Knots/master/data/3_1.dgr)):

```json
{"type":"diagram","name":"Trefoil","components":[{
  "vertices":[[0,56,224],[1,125,85],[2,183,85],[3,318,215],[4,356,380]],
  "crossings":[{"down":2,"up":8},{"down":7,"up":13},{"down":12,"up":3}],
  "isClosed":true}]}
```

Vertices are `[index, x, y]` in screen coordinates.
Crossings reference **edge indices** and name which strand goes `up` and which goes `down`.
Note the explicit **`isClosed`** boolean: the format already anticipates open curves.

`.knt`, the 3D embedding ([3_1.knt](https://raw.githubusercontent.com/Mathesis-Software/Knots/master/data/3_1.knt)):

```json
{"type":"link","name":"Trefoil","components":[{"points":[[-0.9463,0.14736,0.11234]]}]}
```

A dense polyline of `[x,y,z]` floats. Nothing else. No crossings, no semantics, no end metadata.

**Assessment.** This is a *drawing* format, not a *tying* format.
The `.dgr` crossing model (`{down: edgeIdx, up: edgeIdx}`) is a clean idea worth stealing for our schema.
`.knt` is just a sampled curve; any spline library gives you the same.
Neither has any notion of a hook, an object, an operation order, or a working end.
There is **no schema documentation** in the repo; the above is reverse-engineered from the shipped data.
The repo ships **~250 sample `.knt`/`.dgr` files** (3_1 through 17_1, plus connected sums like `3_1#-3_1`) under Apache-2.0. Usable, but they are mathematical knots, not fishing knots.

## 3. Available Databases

| Name | Contents | Format | License | URL |
| --- | --- | --- | --- | --- |
| **KnotInfo** | Table of knot invariants for prime knots. Mathematical only. Created by Chuck Livingston (2004), redesigned with Allison Moore (2019). Partly NSF/Indiana University funded. | Web tables, **CSV download** of query results; larger knot/link databases as **XLS/XLSX**. Also wrapped by SageMath (see [sagemath#30352](https://github.com/sagemath/sage/issues/30352)). | No explicit open license found. Users are *asked to cite* it. Effectively academic-courtesy terms, **UNCERTAIN** whether redistribution is permitted. Canonical host `knotinfo.math.indiana.edu` did **not resolve** from this machine; mirror at [knotinfo.org](https://knotinfo.org/homelinks/about.html). | [linkinfo.math.indiana.edu/homelinks/about.html](https://linkinfo.math.indiana.edu/homelinks/about.html) |
| **The Knot Atlas** | Wiki encyclopedia by Dror Bar-Natan and Scott Morrison. Rolfsen table (≤10 crossings), Hoste-Thistlethwaite (11 crossings), Thistlethwaite link table, 36 torus knots to 36 crossings, per-knot images, and the `KnotTheory\`` Mathematica package. Invariants include Alexander and Jones polynomials, braid words and indices, bridge index, genus, Conway notation, determinant. Mathematical only. | Bulk **"Take Home Database"** as gzipped **RDF**: `katlas.rdf.gz` (~50 MB gz, ~400 MB raw), plus split files `Rolfsen.rdf.gz`, `Knots11..15.rdf.gz`, `Links.rdf.gz`, `TorusKnots.rdf.gz`. | **UNVERIFIED.** No license statement is present on the [main page](https://katlas.org/wiki/Main_Page), the [Take Home Database page](https://katlas.org/wiki/The_Take_Home_Database), or the [Wikipedia article](https://en.wikipedia.org/wiki/The_Knot_Atlas). The disclaimer page 404s. Commonly assumed GNU FDL because it is MediaWiki, but **I could not confirm this and you should not assume it**. | [katlas.org](https://katlas.org/) |
| **SnapPy** | Studies geometry and topology of 3-manifolds; ships the SnapPea cusped-manifold census. Mathematical only, and one step further from us than the others (manifolds, not diagrams). | Python package + census databases | Open source, **GPL-family; exact version UNCERTAIN**. Listed at [CompuTop.org](https://nmd.web.illinois.edu/computop/) | [CompuTop archive](https://nmd.web.illinois.edu/computop/) |
| **Regina** | 3- and 4-manifold topology: triangulations, knots and links, normal surfaces, angle structures. Includes knot censuses, notably an independently tabulated **virtual knot census** (reflections, reversals and flips identified) with classical knots included and all entries certified distinct, plus cross-references to other databases. Includes code from SnapPea/SnapPy and Normaliz. | C++/Python library plus [supporting data downloads](https://regina-normal.github.io/data.html) | **GNU GPL** ([handbook](https://regina-normal.github.io/docs/)) | [regina-normal.github.io](https://regina-normal.github.io/) |
| **Mathesis/Knots data** | ~250 sample knots as JSON 2D diagrams and 3D polylines. Mathematical only. | JSON (`.knt`, `.dgr`) | **Apache-2.0** | [github.com/Mathesis-Software/Knots](https://github.com/Mathesis-Software/Knots) |
| **`dariusk/corpora`** | Verified: **`data/technology/knots.json` exists** and contains a flat list of **212 knot names** under key `knots`, with `description: "A list of knot names."`. Genuinely includes fishing knots (Albright special, Arbor knot, Bimini twist, blood knot, angler's loop, anchor bend). **Names only. No structure, no steps, no geometry, no ABoK numbers.** | JSON | **CC0.** The prior run's observation is confirmed and resolved: there is **no root `LICENSE` file** (GitHub's license API returns "Not Found"), but the [README](https://raw.githubusercontent.com/dariusk/corpora/master/README.md) states *"Since Corpora is more data than code, I have chosen to CC0 license this"* and *"To the extent possible under law, Darius Kazemi has waived all copyright and related or neighboring rights to Corpora."* Contributors agree to CC0 on submission. That is a valid CC0 dedication. | [data/technology/knots.json](https://raw.githubusercontent.com/dariusk/corpora/master/data/technology/knots.json) |

**Verdict across all of them:** every mathematical database contains *closed prime knots indexed by crossing number*.
None contains a single fishing knot.
None contains a tying procedure.
Their value to `riggermortis` is approximately zero for content and mildly useful for schema inspiration.

## 4. Fishing-Specific Sources

### Wikidata: verified, and it is thin

Queried the live [Wikidata SPARQL endpoint](https://query.wikidata.org/) directly.

- The property is **[P1806, "ABoK number"](https://www.wikidata.org/wiki/Property:P1806)**, described as *"identifier for a knot in the Ashley Book of Knots"*.
- `SELECT (COUNT(DISTINCT ?item)) WHERE { ?item wdt:P1806 ?v }` returns **37**. **The prior run's figure is confirmed exactly.**
- There are considerably more *statements* than items, because one practical knot maps to many ABoK entries: Constrictor knot carries #176, #1188, #1249, #1250, #1251; cow hitch carries #1184, #1673, #1694, #1698, #1700. Useful signal about how Ashley catalogued variants and use-contexts separately.
- What is actually there: bowline, clove hitch, constrictor knot, taut-line hitch, figure-eight loop, sheepshank, carrick bend, hangman's knot, Klemheist, distel hitch, cow hitch, grief knot, thief knot, shoelace knot, Hunter's bend, farmer's loop, artillery loop, adjustable bend, bowline on a bight, overhand loop, Ashley's bend.
- **Essentially none are fishing knots.** No Palomar, Improved Clinch, FG, Uni, Albright, Blood, Bimini, Trilene, Arbor, or Snell in the ABoK-tagged set.
- `?item wdt:P31/wdt:P279* wd:Q1093` returned **0**, so Wikidata's knot class hierarchy is not modelled the way you would guess. **UNCERTAIN** what the correct class QID is. Either way, knot modelling in Wikidata is not systematic.

**Verdict: Wikidata is useless as a knot data source.** No tying steps, no geometry, no fishing coverage. At best a source of stable identifiers to link out to.

### Wikimedia Commons: genuinely useful for reference imagery

Queried the Commons API directly.

**[Category:Fishing knots](https://commons.wikimedia.org/wiki/Category:Fishing_knots)** has **40 members**, including subcategories: Albright special, Angler's loops, Arbor knots, Bimini twist, Blood knots, Cat's paw, Domhof knot, Fisherman's knots, Fishing loop, Palomarknots, Surgeon's knot, Trilene knot, World's fair knots.

Verified licenses on specific files (via `extmetadata`):

| File | License |
| --- | --- |
| [Albright knot diagram retouched.png](https://commons.wikimedia.org/wiki/File:Albright_knot_diagram_retouched.png) | **Public domain** |
| [Bimini Twist knot.svg](https://commons.wikimedia.org/wiki/File:Bimini_Twist_knot.svg) | **Public domain** |
| [Bumper knot diagram.png](https://commons.wikimedia.org/wiki/File:Bumper_knot_diagram.png) | **Public domain** |
| [FMIB 46700 Turle knot.jpeg](https://commons.wikimedia.org/wiki/File:FMIB_46700_Turle_knot.jpeg) | **Public domain** |
| [PalomarKnotSequence.jpg](https://commons.wikimedia.org/wiki/File:PalomarKnotSequence.jpg) | **CC BY-SA 3.0** |
| [Uni knot.jpg](https://commons.wikimedia.org/wiki/File:Uni_knot.jpg) | **CC BY-SA 3.0** |
| [BloodKnot HowTo.jpg](https://commons.wikimedia.org/wiki/File:BloodKnot_HowTo.jpg) | **CC BY-SA 3.0** |

The **FMIB** prefix is the Freshwater and Marine Image Bank, a large set of scanned pre-1929 angling illustrations, uniformly public domain. `Category:Fishing knots` contains a run of them (Turle knot, loop knot for drop flies, Francis knot, knot for turned-down metal-eyed hook, whipping on a treble hook, etc.), plus BHL scans from *American Game Fishes*.

Commons-wide policy, per the [category page](https://commons.wikimedia.org/wiki/Category:Knot_animations): *"All structured data from the file namespace is available under the Creative Commons CC0 License; all unstructured text is available under the Creative Commons Attribution-ShareAlike License."* Individual media files carry their own license, so **check per file**.

[Category:Knot animations](https://commons.wikimedia.org/wiki/Category:Knot_animations) has **28 files**, and [Category:Blue knots rendered; animated](https://commons.wikimedia.org/wiki/Category:Blue_knots_rendered;_animated) has 27. **None are fishing knots**: they are bowlines, carrick bends, trefoils, Turk's heads and mathematical unknottings, mostly CC BY-SA 4.0.

**Verdict: Commons is the best free source of fishing-knot *reference imagery*, and a fair chunk of it is public domain.** It is a source for humans authoring data, and legally safe as visual reference. It contains **zero** structured or animatable data.

### GitHub: nothing exists

Searched the GitHub API and code search.

- `topic:knots` returns mathematical tools and rope-physics simulators. Top hits: `lennart-finke/knottingham` (MIT, JS, 67 stars, actively maintained, "A Tool for Drawing Pretty Knots"), `Mathesis-Software/Knots` (Apache-2.0), `QuantuMope/imc-der` (MIT, C++, elastic-rod contact model "framework for knot tying"), `StructuresComp/rod-contact-sim` (GPL-3.0), `dmackinnon1/celtic` (Celtic-pattern generator, no license).
- `fishing knot` repo search returns **658 results that are entirely noise** (emoji cheat sheets, OSINT lists, AIS libraries). There is no fishing-knot project.
- Code search for the literal string `"palomar"` in JSON files across public GitHub: **no results**.

**Verdict: there is no open-source structured fishing-knot dataset. Not a partial one, not a stale one. None.**
If `riggermortis` ships one, it will be the first.

## 5. Licensing Verdicts

| Source | License | Can we use it? | Caveat |
| --- | --- | --- | --- |
| **Ashley Book of Knots** text/illustrations | **In copyright until 2040-01-01.** Renewed 1971. | ❌ **No** | See below. Number *citation* is fine; content is not. |
| **ABoK reference numbers** (e.g. `#1249`) | Short factual identifiers | ✅ **Yes, cite freely** | Do not ship a wholesale ABoK concordance. See caveat below. |
| **animatedknots.com / Grog LLC** | All rights reserved, actively enforced | ❌ **Absolutely not** | Not even as an animation reference that would produce a substantially similar sequence. |
| **Wikimedia Commons, PD files** | Public domain | ✅ **Yes, unrestricted** | Verify per file. FMIB/BHL angling scans are the richest PD seam. |
| **Wikimedia Commons, CC BY-SA files** | CC BY-SA 3.0 / 4.0 | ⚠️ **Yes with care** | ShareAlike is viral for derivatives of the *image*. Safe as reference for a human; do not trace and ship without complying. |
| **Commons structured data** | CC0 | ✅ **Yes** | Metadata only, not the media. |
| **`dariusk/corpora` knots.json** | **CC0** (README dedication; no `LICENSE` file) | ✅ **Yes** | Only 212 knot *names*. Fine as a seed vocabulary. |
| **Mathesis/Knots `.knt`/`.dgr` + data** | **Apache-2.0** | ✅ **Yes** | Mathematical knots only. Schema is undocumented; reverse-engineer it. |
| **pyknotid** | **MIT** | ✅ **Yes** | Classifier, not a representation. Last release 2023. |
| **Knoto-ID** | Paper says **GPL v2+**; repo license detection empty | ⚠️ **Probably, verify first** | GPL is copyleft. Fine as a CLI you shell out to; do not link into a permissively-licensed binary without checking. Unmaintained since 2019. |
| **Regina** | **GPL** | ⚠️ **Copyleft** | Same reasoning. Also irrelevant content-wise. |
| **KnotInfo** | No open license; "please cite" | ⚠️ **UNCERTAIN** | Do not redistribute the tables without asking. Mathematical only anyway. |
| **Knot Atlas** | **No license statement found anywhere** | ⚠️ **UNCERTAIN, treat as all-rights-reserved** | Commonly assumed FDL. I could not confirm it. Do not redistribute the RDF dumps on that assumption. |
| **Peter Suber's Knot Tying Notation** | `Copyright © 2002-2006, Peter Suber`, **no license grant** | ⚠️ **Implement the system, do not copy the page** | See below. |

### 🔴 The Ashley Book of Knots is under copyright until 2040. The "public domain" claim is false

This is the most important finding in this document, and it reverses a prior conclusion.

**Primary evidence.** The NYPL `cce-renewals` dataset, which is tab-delimited transcriptions of the US *Catalog of Copyright Entries* renewal volumes released **CC0 1.0** ([LICENSE](https://raw.githubusercontent.com/NYPL/cce-renewals/master/LICENSE), [repo](https://github.com/NYPL/cce-renewals)), contains this record in `data/1971-1.tsv`:

```text
ASHLEY, CLIFFORD W.  The Ashley book of knots.
© 21Jul44; A181853.  Sarah R. Delano (W); 2Aug71; R510678.
```

Decoded:

| Field | Value |
| --- | --- |
| Author | Clifford W. Ashley |
| Original copyright registration | **A181853**, dated **1944-07-21** |
| Renewal registration | **R510678**, dated **1971-08-02** |
| Renewal claimant | Sarah R. Delano, `(W)` = widow of the author |

A 1944 work required renewal in its 28th year (1971-72) to retain copyright.
**It was renewed, on time, by the correct statutory claimant.**
Under the Copyright Act as extended (28 + 67 years for renewed pre-1978 works), the term runs **95 years from publication**, so the book enters the US public domain on **2040-01-01**.

**The prior run's claim that this was "independently confirmed by NYPL's CC0 dataset" is exactly backwards.** The NYPL dataset is the strongest available evidence that it is *not* public domain.
The widely-cited [smallboatsmonthly.com 2017 article "Ashley Book of Knots Now in the Public Domain"](https://smallboatsmonthly.com/2017/08/ashley-book-knots-now-public-domain/) appears to be **wrong**.
Consistent with that, the [Internet Archive copy](https://archive.org/details/ashleybookofknot0000clif_z0z2) is a **lending-library** item (borrow, not free download), which is how IA handles in-copyright books.

Consequences for `riggermortis`:

- ❌ Do not copy, trace, redraw, or derive from ABoK **illustrations**. Ashley drew all ~7,000 of them; they are protected pictorial works.
- ❌ Do not copy ABoK descriptive **text**.
- ✅ **ABoK reference numbers can be cited.** A number like `#1249` is a short factual identifier. Individual numbers are not copyrightable subject matter, and citing them is universal practice in the field: [Wikipedia does it](https://en.wikipedia.org/wiki/The_Ashley_Book_of_Knots) ("The numbers Ashley assigned to each knot can be used to unambiguously identify them"), the IGKT does it, every knot site does it. The book has 3,857 numbered entries; three of them are non-integers (794.5, 1034.5, 2585.5).
- ⚠️ **Caveat that matters:** reproducing the *numbering scheme as a whole*, i.e. shipping a complete ABoK-number-to-knot-name concordance, is a different question, because a comprehensive catalogue can attract thin compilation protection. Cite numbers on knots you independently document; do not ship the index as a dataset. **UNCERTAIN**, this is a judgement call and not settled law I can cite. If the project ever wants a large ABoK concordance, involve an actual lawyer.

### 🔴 animatedknots.com / Grog LLC: hard no

[Copyright & Privacy Policy](https://www.animatedknots.com/copyright-and-privacy-policy). Copyright holder: **Grog LLC**. Verbatim:

> "The images and the text contained in this website are protected under United States Copyright law."
>
> Grog LLC "vigorously and regularly identifies examples of Copyright infringement and protects its interests by insisting that such material identified on the web or elsewhere be either used under license or destroyed."

The **only** reuse permission found anywhere on the site: if you run a website and link to them, you may use **a single image from an animation sequence as the link**.
Nothing else.
No fair-use carve-out, no educational exception, no attribution-based permission.

**Verdict: cannot use animatedknots.com content, frames, animations, or text in any form.**
Do not scrape it, do not trace it, do not use it as a frame-by-frame animation reference.
Their entire business model is licensing and they say plainly that they enforce.
Treat the site as a competitor, not a source.
The same reasoning applies to [Knots 3D](https://knots3d.com/), a commercial 200-knot animated app.

### ⚠️ Peter Suber's notation: use the system, not the document

The page carries `Copyright © 2002-2006, Peter Suber.` and no license grant.

What that means in practice:

- ✅ **The notation itself is a system, and systems are not copyrightable.** Under the idea/expression dichotomy (and *Baker v. Selden* reasoning), you can implement Suber's symbols and semantics in a parser, a schema, and a renderer. Symbol names like `RP`, `MT`, `RV` are short functional identifiers.
- ❌ **Do not copy his prose, his tables, or his worked examples verbatim** into project docs. Paraphrase and attribute.
- ✅ **Attribute him anyway.** It is the right thing to do and costs nothing.
- 💡 **Just ask him.** Peter Suber is the single most prominent open-access advocate alive (Harvard Open Access Project). The odds he refuses a CC license or a clarifying statement for an open-source knot project are low. This is the highest-value, lowest-effort action available on this whole list. **UNCERTAIN** whether the page is already covered by a blanket dedication elsewhere on his site; I did not find one.

## 6. Prior Art in Code

| Project | What it does | Representation | License | URL |
| --- | --- | --- | --- | --- |
| **Animated Knots in Three.js** (Alessandro Stefanini) | Animates rope tying a knot in the browser. **Closest thing to what riggermortis is doing.** | **Hand-authored curves.** Rope built on `CatmullRomCurve3` interpolating control positions in space, animated with **Theatre.js** (a keyframe/timeline tool). Nothing is generated from notation; a human places control points and keyframes them. | Blog post, license of any code **UNCERTAIN** | [alessandrostefanini.it/2021/12/24/knots-three-js](https://alessandrostefanini.it/2021/12/24/knots-three-js/) |
| **knottingham** | Interactive knot drawing tool, JS. Actively maintained (last push 2026-07). | Diagram-based, mathematical knots | **MIT**, 67 stars | [github.com/lennart-finke/knottingham](https://github.com/lennart-finke/knottingham) |
| **Mathesis/Knots (KnotEditor)** | Draw 2D diagram, convert to 3D, smooth, Seifert surfaces, invariants. **The 2D-to-3D pipeline is the single most relevant piece of prior art**: it turns a crossing diagram into a 3D embedding, which is a problem we will have to solve too. | JSON `.dgr` (2D + crossings) to `.knt` (3D polyline) | **Apache-2.0** | [github.com/Mathesis-Software/Knots](https://github.com/Mathesis-Software/Knots) |
| **imc-der** | Implicit contact model for 3D elastic rod simulation, explicitly framed as a "framework for knot tying". Physics, not notation. | Discrete elastic rods + contact | **MIT** | [github.com/QuantuMope/imc-der](https://github.com/QuantuMope/imc-der) |
| **rod-contact-sim** | Sibling elastic-rod contact simulator. | Discrete elastic rods | **GPL-3.0** | [github.com/StructuresComp/rod-contact-sim](https://github.com/StructuresComp/rod-contact-sim) |
| **blender-rope-sim** | Blender rope simulation for robotics knot-tying/untangling research; models self-collision and knots with realistic rope appearance. | Blender physics | **UNCERTAIN** | [github.com/priyasundaresan/blender-rope-sim](https://github.com/priyasundaresan/blender-rope-sim) |
| **Knot Madness** | Hackathon project: 3D knot practice with hand tracking. Three.js + Cannon.js + MediaPipe. | Physics sandbox, no knot data | No license file | [github.com/Pho86/mountain-madness-2025](https://github.com/Pho86/mountain-madness-2025) |
| **gridmaker** | Generates tying instructions for Turk's head knots. Rare example of *generated instructions*, but for one narrow knot family. | Grid/algorithmic | No license file | [github.com/AllwineDesigns/gridmaker](https://github.com/AllwineDesigns/gridmaker) |
| **Knots 3D** (commercial) | 200+ animated 3D knots. The market leader. | Proprietary | ❌ Commercial | [knots3d.com](https://knots3d.com/) |

**The pattern is unambiguous.** Every project that renders a *recognisable practical knot* uses **hand-authored curves** (control points, splines, keyframes).
Every project that *generates* geometry works on **mathematical knots** (from PD/DT/Gauss codes) or on **physics** (elastic rods, which produce a rope but do not know what knot it is).
**Nobody generates a practical knot animation from a structured tying description.** That gap is the project.

## 7. Honest Assessment

### What genuinely exists and is usable

1. **A procedural notation designed for exactly this problem.** Suber's Knot Tying Notation: stages, running part vs standing part, `MT` for wrapping a named external object, `RV` for reeving through an opening. Adopt its *model*; the licensing needs a five-minute email to resolve. This is the single most valuable find here.
2. **A proven design pattern for a move grammar.** Fink/Mao/Vejdemo-Johansson necktie work proves a regular language over tying moves is tractable and enumerable ([arXiv:1401.8242](https://arxiv.org/pdf/1401.8242)).
3. **A concrete crossing-encoding idea, Apache-2.0.** Mathesis `.dgr`'s `{down: edgeIdx, up: edgeIdx}` with an `isClosed` flag, plus its 2D-to-3D conversion as reference for the hard geometry step.
4. **Public-domain reference imagery.** Commons `Category:Fishing knots`, especially the FMIB and BHL angling scans. Legally clean, useful to whoever authors the data.
5. **A CC0 seed vocabulary.** 212 knot names in `dariusk/corpora`, including real fishing knots.
6. **An MIT-licensed topology classifier if you want dedup/validation.** `pyknotid`, applied to generated 3D geometry.

### What does not exist and must be authored from scratch

1. **All fishing-knot data.** There is no dataset, anywhere, in any license. Not a partial one, not an abandoned one. GitHub code search for `"palomar"` in JSON returns nothing. Wikidata has 37 ABoK-tagged items and essentially no fishing knots. This is not a gap you can fill by aggregating; somebody has to sit down and encode Palomar, Improved Clinch, Uni, Blood, Albright, FG, Bimini, Trilene, Arbor, Snell, Loop knots, one at a time, by hand.
2. **The object model.** No formalism covers "line tied around a hook eye" in a computationally usable way. Knotoids handle open ends but assume a single strand in a surface. Torus knotoids exist but are not the right construction and have no software. You will define `Eye`, `Shank`, `Swivel`, `OtherLine` as first-class objects yourself, and no paper will tell you the right shape for that.
3. **Multi-strand knots.** Bends (blood knot, Albright, FG) are tangles, not knotoids. Whatever you design has to handle two independent strands from day one, or you will retrofit it painfully later.
4. **Notation-to-geometry.** Suber tells you what to do, not where the rope is. Turning stages into 3D control curves is unsolved for practical knots and is the actual engineering risk in this project. Mathesis's diagram-to-3D code and the elastic-rod simulators (`imc-der`, MIT) are the closest starting points, in that order.
5. **Everything mechanical.** Friction, line diameter, dressing, seating, tightening order, breaking strength, which knot suits braid vs fluorocarbon. Zero formal support anywhere. Model it as plain domain metadata and stop looking for theory.

### Things that look promising and are not

- **Knot theory generally.** Real, deep, well-tooled, and answers a question we do not have. Do not let its rigour seduce the schema design. Optional fingerprint layer at most.
- **KnotInfo / Knot Atlas / SnapPy / Regina.** Every one of them is closed prime knots indexed by crossing number. Zero fishing knots, zero tying steps. The Knot Atlas additionally has **no discoverable license**, so redistributing its dumps is a risk with no upside.
- **Wikidata.** 37 items. The ABoK property exists and is basically unpopulated. Link out to it; do not source from it.
- **ABoK itself.** The obvious canonical source is under copyright for another 14 years. Cite the numbers, use nothing else.
- **animatedknots.com.** The best free fishing-knot animations on the internet are also the most aggressively defended. Treat as radioactive.
- **Rope physics simulators.** They produce beautiful rope and have no idea what knot they are tying. Useful as a *renderer* backend, useless as a data model. Do not mistake one for the other.

### Recommended architecture

1. **Author a Suber-derived procedural schema.** Stages as animation keyframes; explicit `working_end` / `standing_part`; first-class named objects (`eye`, `shank`, `swivel`, `line_b`); verbs approximately `wrap`, `reeve`, `loop`, `bight`, `cross`, `pull`, `dress`, `seat`. Attribute Suber, and email him.
2. **Keep an optional crossing-level layer** (Mathesis-style `{over, under}` on strand segments) so a diagram renderer and a topology fingerprint are both derivable later. Do not make it the primary representation.
3. **Author the fishing-knot corpus by hand**, using Commons PD imagery as reference. Ship it CC0 or CC BY so it becomes the thing that did not exist.
4. **Solve notation-to-geometry with control curves first** (`CatmullRomCurve3`-style, as in the Three.js prior art), and consider an elastic-rod relaxation pass (`imc-der`, MIT) only if hand-tuned curves prove unmaintainable.
5. **Cite ABoK numbers, ship no ABoK content.**

### Open items

- Peter Suber's licensing (email him).
- Knot Atlas license (email Bar-Natan or Morrison) if you ever want their dumps, which you probably do not.
- Knoto-ID's in-repo license, if you end up shelling out to it.
- The IGKT forum thread on knot notation (`forum.igkt.net` did not resolve from this machine); worth a manual look for prior community attempts.
- Whether a "knotoid in a solid-torus complement" formalism has been developed. **UNCERTAIN**, and probably academic curiosity rather than a blocker.
