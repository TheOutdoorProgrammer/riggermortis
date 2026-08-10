# Rigging Notation Research

Question: does an established, tried-and-true notation exist for describing how a flexible or soft body is pierced, threaded, or mounted onto a hook, needle, or similar?
If not, which existing systems are the best structural models to borrow from?

Motivating case: encoding the Texas rig "Tex-posed" method as data.
Insert the point into the nose, exit through the belly, slide the bait up the shank, rotate the hook 180 degrees, re-pierce, bury the point just under the skin.
This is the rigging equivalent of what [Peter Suber's Knot Tying Notation](https://legacy.earlham.edu/~peters/knotting/notate.htm) already does for knots in this project.

Short answer up front: no such notation exists, in fishing or anywhere else.
The best structural model to borrow is [knitout](https://textiles-lab.github.io/knitout/knitout.html), with the per-pierce parameter tuple taken from surgical needle-steering literature.
Section 8 states the case plainly.

## 1. Fishing-Specific Notation

Finding: nothing exists.
This was searched from several angles and every angle came back empty.

Searches for a formal or machine-readable system for describing bait rigging or hook placement return only patents for physical hardware: [hook baiting machines](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/4015359), [baiting methods for longline hooks](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/5934003), and a [bait rigging system](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/5377442).
These patents describe mechanisms that put bait on hooks.
None of them defines a language for describing how a bait is rigged, which is what a patent in this space would need if a notation existed to reference.

There is no fishing rig ontology, JSON schema, or interchange format in open source.
A direct search for one returned nothing fishing-related at all.

Four things in fishing come close enough to be worth naming, and all four fall short in the same way.

Named rig taxonomy is real but is a taxonomy of end states, not a procedure notation.
Texas, Carolina, wacky, Neko, drop shot, Ned, and [Tex-posed](https://www.wired2fish.com/fishing-rigs/the-texas-rig-how-to-rig-and-fish) are stable, widely understood names.
Every instructional source describing them, including [Tackle Warehouse](https://www.tacklewarehouse.com/bass-fishing/how-to/how-to-rig-a-texas-rig.html) and [Vermont Fish and Wildlife](https://www.vtfishandwildlife.com/fish/fishing-opportunities/fishing-basics/common-lures-and-rigs/bottom-lures-texas-rig), teaches them as numbered prose steps plus photographs.
The name is a label for a result, and the procedure behind the label is never encoded.

Fisheries science has a semi-standard vocabulary for hooking *location*, which is the single most notation-like thing in the sport.
Catch-and-release mortality studies consistently partition outcomes into anatomical sites: lip or jaw, inside the mouth, gills, throat or esophagus, and foul-hooked.
See [Lindsay and Wydoski, hooking mortality by anatomical location](https://myodfw.com/sites/default/files/2025-12/References_1_Lindsay-Wydoski.pdf), the [NOAA release mortality white paper](https://www.st.nmfs.noaa.gov/Assets/ecosystems/bycatch/documents/Release_Mortality_White_Paper_Final.pdf), and [Myers et al. on bait type and hooking location](https://seafwa.org/sites/default/files/journal-articles/MYERS-39-45.pdf).
This is a controlled vocabulary of anatomical sites on a soft body, which is genuinely one of the primitives a rigging notation needs.
It is convention rather than a published standard, the categories vary between studies, and it describes where a hook ended up in a caught fish rather than where an angler deliberately places one.

Fly tying recipes are structured, but they are ingredient lists rather than procedures.
A pattern specifies hook, thread, tail, body, rib, hackle, and wing as named slots, then narrates assembly in prose.
[Community databases](https://flytyingrecipes.com/) and [Global FlyFisher](https://globalflyfisher.com/patterns) both use this shape, and tyers have repeatedly asked for a real queryable schema without one emerging.
There is no standard fly pattern format.

FAO's international gear classification is sometimes offered as a counterexample.
It classifies gear types such as longline, gillnet, and trap.
It says nothing about how a bait is mounted.
UNCERTAIN: this was not verified with a primary source in this run, but the distinction between gear class and rigging method is not in doubt.

Conclusion for this section: fishing has vocabulary but no grammar.
It has names for rigs, names for hooks, and names for where a hook ends up in a fish.
It has no way to write down a rigging procedure such that two people, or a person and a machine, reconstruct the same result.

## 2. Surgical Suturing

Surgery is the closest analogue and it splits cleanly into two very different things.
One is a mature classification that is not a notation.
The other is a genuine formalism that exists only because robots needed it.

### 2.1 Suture pattern classification is prose plus diagrams

Suture patterns are classified along stable orthogonal axes: interrupted versus continuous, and appositional versus inverting versus everting, with a further simple versus tension split.
See [WikiVet](https://en.wikivet.net/Suture_Patterns) and [Veterinary Surgery Online](https://www.vetsurgeryonline.com/suture-materials-and-suture-patterns/).
Named patterns include simple interrupted, horizontal mattress, vertical mattress, cruciate, subcuticular, purse-string, Lembert, Cushing, and Connell.

This is a real, universally taught, centuries-stable classification.
It is not a notation.
Every individual pattern is defined by prose plus an illustration.
The [Saskatchewan veterinary teaching material](https://wcvm.usask.ca/vsac205/Lab3/interrupted-suture-patterns.php) is representative: horizontal mattress is described as "insert the needle 6 to 8 mm laterally on the same side as the needle exited, cross to the other side of the incision, exit directly across 6 to 8 mm lateral to the first entry".
That sentence contains exactly the primitives a rigging notation needs (entry site, exit site, offset, side, crossing) and encodes none of them symbolically.

The structural lesson is still worth taking.
The two-axis classification (topology: interrupted or continuous; tissue effect: apposing, inverting, everting) is a good model for classifying rigs along axes rather than as a flat name list.
A rigging equivalent would be something like weedless versus exposed, and fixed versus sliding.

### 2.2 Robotic suturing produced the actual formalism

Where a machine has to execute the stitch, the path becomes a parameter tuple.

The single most useful artifact found in this entire research pass is Alterovitz et al.'s needle insertion plan, from [Planning for Steerable Bevel-tip Needle Insertion Through 2D Soft Tissue](https://robotics.cs.unc.edu/publications/Alterovitz2005_ICRA.pdf).
An insertion plan is defined as:

```text
X = (y0, theta, b, d)

y0     insertion location on the tissue surface
theta  insertion angle, in [-90, +90] degrees
b      bevel rotation, in {0, 180} degrees
d      insertion distance (depth)
```

That is a four-field record fully specifying "drive a curved sharp thing into a soft body".
It maps onto bait rigging almost field for field: entry site, angle relative to the bait's long axis, hook point orientation, and how far in.

The complementary work is [Needle Path Planning for Autonomous Robotic Surgical Suturing](https://pmc.ncbi.nlm.nih.gov/articles/PMC3966119/) (Sen, Cavusoglu et al.), which parameterizes a stitch with an entry point g, an exit point f, a tip vector k, an initial insertion distance alpha, and a penetration depth beta, with needle pose carried as a homogeneous transform in SE(3).
An entry point plus an exit point plus a depth is precisely the "nose in, belly out" primitive.

[Webster et al.'s nonholonomic model of needle steering](https://limbs.lcsr.jhu.edu/wp-content/papercite-data/pdf/websternonholonomic2006.pdf) supplies the reason the bevel rotation field exists: an asymmetric tip makes the needle curve, so the roll of the point is a first-class control input.
A hook point buried under bait skin has the same property, which is why "rotate the hook 180 degrees" is a step in the Tex-posed method rather than an incidental detail.

Caveat: this is geometry, not symbols.
It is a coordinate specification for one insertion, with no vocabulary, no sequencing, and no way to say "slide the bait up the shank".
It gives the shape of a single instruction's arguments, not the shape of a program.

### 2.3 Surgical gesture vocabularies are the verb-set analogue

Separately, surgery has developed exactly the thing Suber's verb set is: a finite, agreed list of action primitives.

The [JIGSAWS dataset](https://cirl.lcsr.jhu.edu/research/hmm/datasets/jigsaws_release/) from Johns Hopkins and Intuitive defines 15 gestures, called surgemes, used to annotate suturing, knot tying, and needle passing.
The suturing subset is:

```text
G1   reaching for needle with right hand
G2   positioning needle
G3   pushing needle through tissue
G4   transferring needle from left to right
G5   moving to center with needle in grip
G6   pulling suture with left hand
G7   pulling suture with right hand
G8   orienting needle
G9   using right hand to help tighten suture
G10  loosening more suture
G11  dropping suture at end and moving to end points
```

Note how close G2, G3, and G8 are to what a rigging notation needs: position, push through, orient.
Note also that hand assignment is baked into the verb, which is a design choice worth avoiding (it multiplies the vocabulary; better to carry the hand as a slot).

The lineage matters here.
This grew out of the [Language of Surgery project](https://www.sciencedaily.com/releases/2006/12/061208101350.htm), which explicitly borrowed speech-recognition methods and treated surgical motions as words forming sentences.
That is the same intellectual move this project is making with knots.

The current state of the art is the [SAGES Delphi consensus surgical gesture taxonomy](https://link.springer.com/article/10.1007/s00464-026-12906-2) (2026), which is the most rigorously derived action vocabulary found anywhere in this research.
It started from 270 literature-derived gesture terms, used LLM-assisted semantic clustering plus multi-round expert review with an 80 percent agreement threshold, and converged on a hierarchy of 10 clusters, 24 gestures, and 46 sub-gestures.
Two findings from it are directly transferable.
First, the final structure is hierarchical (cluster, gesture, sub-gesture) rather than flat, which is how a rigging vocabulary should be organized too.
Second, the panel explicitly rejected labeling only the dominant instrument in favour of multi-instrument annotation, because the assisting action matters.
For rigging that translates to: the hand holding and compressing the bait is part of the instruction, not context.

There are also surgical process ontologies, [OntoSPM and LapOntoSPM](https://ontospm.univ-rennes1.fr/doku.php?id=ontology), which model surgery as phases, steps, and activities with instruments.
These are the OWL-ontology approach to the same problem.
Relevant as prior art for structure, heavier than this project needs.

## 3. Textile Notations

This section contains the best structural model found.

### 3.1 Knitting charts: standardized symbols, not machine-readable

The [Craft Yarn Council](https://www.craftyarncouncil.com/standards/knit-chart-symbols) publishes standardized knit and [crochet](https://www.craftyarncouncil.com/standards/crochet-chart-symbols) chart symbols, with [free downloadable symbol graphics](https://www.craftyarncouncil.com/standards/downloadable-symbols) and a broader [standards and guidelines document](https://media.craftyarncouncil.com/sites/default/files/images/standards/CYC_YarnStandards-2018-11-06.pdf).

How the symbology works: a chart is a 2D grid where each cell is one stitch and each row is one row of knitting, and each symbol depicts the stitch as it appears from the right side of the work.
Charts are read bottom to top, and alternate rows reverse direction because the work is turned.

Two limitations matter.
CYC itself notes that patterns commonly invent project-specific symbols and that a pattern should always carry its own key, so the standard is a strong convention rather than a closed set.
And a chart is fundamentally a spatial map of a finished fabric, not an ordered instruction stream, which means it does not transfer to a procedure with rotation and re-piercing.

### 3.2 knitout is the structural model to steal

[knitout](https://textiles-lab.github.io/knitout/knitout.html), the .k file format from the [Carnegie Mellon Textiles Lab](https://github.com/textiles-lab/knitout), is a machine-independent low-level instruction format for knitting machines.
It is deliberately dumb: the spec states it contains "no flow-control, abstractions, or grouping primitives", only the instructions required to fabricate the object.
It is line-oriented UTF-8, opens with a version magic line `;!knitout-2`, carries a `;;Key: value` header block (Carriers is required; Machine, Gauge, Yarn-C, Position optional), and reserves `X-` headers and `x-` opcodes for extensions.

The full opcode set:

| Opcode | Params | Meaning |
| --- | --- | --- |
| `in` | CS | bring carrier set into action from the grippers |
| `inhook` | CS | bring carrier set in using the yarn inserting hook |
| `releasehook` | CS | release yarns held in the inserting hook |
| `out` | CS | bring carriers out of action into the grippers |
| `outhook` | CS | bring carriers out via the inserting hook |
| `stitch` | L T | set pre-loop and post-loop needle pull |
| `rack` | R | set back bed offset relative to front bed |
| `knit` | D N CS | pull a loop formed in direction D through loops on needle N, dropping them |
| `tuck` | D N CS | add a loop formed in direction D to those already on needle N |
| `split` | D N N2 CS | pull a loop through loops on N, transferring the old loops to N2 |
| `drop` | N | synonym for knit with no carriers |
| `amiss` | N | synonym for tuck with no carriers |
| `xfer` | N N2 | synonym for split with no carriers |
| `miss` | D N CS | move carriers as if they had formed a loop in direction D at N |
| `pause` | none | halt until the operator continues |
| `x-*` | any | extension opcode |

The argument types are the important part.
Direction D is `+` or `-`.
Needle N is a bed prefix plus an integer: `f` front, `b` back, `fs` front slider, `bs` back slider, so `f1`, `b-2`, `fs3`.
Carrier set CS is a list of carrier names.

That gives the canonical instruction shape `VERB DIRECTION LOCATION SUBJECT`, where location is addressed as face-plus-index.
This is close to a direct answer to the rigging problem, and section 6 spells out the mapping.

Also relevant: [KnitScript](https://dl.acm.org/doi/fullHtml/10.1145/3586183.3606789) is a higher-level scripting DSL that compiles down to knitout, and [Coupling Programs and Visualization for Machine Knitting](https://dl.acm.org/doi/fullHtml/10.1145/3424630.3425410) covers the tooling layer.
The two-layer split (readable authoring language on top, dumb linear IR underneath) is a proven architecture in exactly this problem class.

Crochet has no equivalent.
Research prototypes such as CrochetPARADE encode stitches as primitives with explicit control flow, and there is [work on machine crochet](https://arxiv.org/pdf/2511.09483), but no adopted standard.
UNCERTAIN: crochet was searched once, not exhaustively.

### 3.3 Embroidery: coordinate lists, not notation

Machine embroidery formats are real and machine-readable, and they are the wrong shape.
As described in the embroidery data patents ([US 4964352](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/4964352), [US 6167823](https://image-ppubs.uspto.gov/dirsearch-public/print/downloadPdf/6167823)), a design is a numbered sequence of needle penetrations stored as relative X and Y offsets from the previous penetration, interleaved with jump, trim, and colour-change commands.

So embroidery formats encode literally "where the needle pierces", which sounds ideal, and they carry no depth, no entry surface versus exit surface, no notion of which side of the material anything is on, and no rotation of one object relative to another.
Fabric is treated as a plane with zero thickness.
That is exactly the assumption a bait rigging notation cannot make.

Hand embroidery stitches (satin, chain, French knot, couching) are named and taught with prose plus diagrams, with no notation.

### 3.4 ISO 4915: a real international standard, but for structures not procedures

[ISO 4915:1991, Textiles, Stitch types, Classification and terminology](https://standards.globalspec.com/std/473931/iso-4915) classifies 88 stitch types into six numbered classes: 100 single-thread chainstitch, 200 hand and hand-effect stitches, 300 lockstitch, 400 multi-thread chainstitch, 500 overedge, 600 covering.
Saddle stitch is type 201, plain lockstitch is 301, multi-thread chainstitch is 401.
A [technical review aligned to ISO 4915 and ASTM D6193](https://www.opastpublishers.com/open-access-articles/industrial-sewing-stitch-classification-a-technical-review-aligned-with-iso-4915-and-astm-d6193-standards.pdf) is openly available and covers the scheme.

The genuinely interesting detail: ISO 4915 assigns symbolic identifiers to thread *roles*.
Needle threads are 1, 2, 3; looper threads are a, b, c; cover threads are Z, Y, X.
That is the same design instinct as Suber's running part versus standing part, and it is the closest thing in the stitching world to role-typed nouns.

It remains a classification of resulting structures, not a procedure notation, and it is a paywalled ISO document.

## 4. Robotics and Deformable Object Manipulation

There is no action representation language for deformable object manipulation.
This was checked directly and the answer is consistent across the literature.

The [Frontiers tutorial and review on modeling deformable objects for robotic manipulation](https://www.frontiersin.org/journals/robotics-and-ai/articles/10.3389/frobt.2020.00082/full) frames the field's representations as geometric and physical models, not symbolic ones: mass-spring systems, finite element models, position-based dynamics, and particle sets.

Where symbolic structure appears at all, it is at the task-planning layer and it is natural language.
The [survey of language-conditioned robot manipulation](https://arxiv.org/pdf/2312.10807) notes that most language-conditioned work targets rigid objects and that a significant gap remains between task planning and action execution for deformables.
[CLASP](https://arxiv.org/pdf/2507.19983) proposes semantic keypoints as the interface between a vision-language model's plan and low-level execution for clothes, which is the closest thing to a shared vocabulary and is a set of labelled points on a body rather than a grammar.

The transferable idea is small but real: semantic keypoints on a deformable body are exactly the abstraction a bait needs.
A soft plastic worm's nose, collar, belly seam, and tail junction are semantic keypoints, and naming them is what makes an instruction reproducible across bait sizes.

The needle-insertion literature in section 2.2 is where robotics actually delivers, and it delivers parameters rather than symbols.

## 5. Other Domains

### 5.1 Entomology specimen pinning

Real, old, international conventions, and no notation.

Pin sizes are standardized: [insect pins](https://en.wikipedia.org/wiki/Insect_pins) are 38 mm long in sizes 000 through 8.
Specimen height on the pin is standardized via a stepping block.
Pin placement is taxon-specific and taught as a lookup: beetles through the right elytron, true bugs through the scutellum, Lepidoptera through the centre of the thorax, flies and wasps through the thorax right of the midline.
See the [NC State pinning guide](https://genent.cals.ncsu.edu/students/lab-schedule/lab-1-collecting-preserving-part-1-meadow-field-trip/insect-collection-instructions-2/a-guide-to-mounting-insects-on-pins/), [Oregon State extension](https://extension.oregonstate.edu/sites/extd8/files/documents/1/pinninginsects.pdf), [Kansas State](https://entomology.k-state.edu/doc/4-h-collections/pinning.pdf), and the [UC Riverside Entomology Research Museum](https://entmuseum.ucr.edu/guides-faqs/specimen-preparation).

This is the closest real-world analogue to bait rigging in one specific respect: it is a mapping from body type to insertion site on that body, with a rule for depth and angle, applied to a soft body with a rigid pin.
It is delivered entirely as prose plus annotated photographs, published independently by extension services and museums, with no symbolic encoding and no central standards body.

### 5.2 Leatherwork saddle stitch

Prose plus pictures.
Every serious source ([Tandy](https://tandyleather.com/blogs/tandy-blog/tandy-skills-saddle-stitching), [Walsall Leather Skills Centre](https://walsall-leather.org/a-beginners-guide-to-saddle-stitching-in-leatherwork/), [Squushed's illustrated guide](https://squushed.com/blogs/guides/how-to-saddle-stitch-an-illustrated-guide)) teaches it as numbered steps with photographs, including the conventions that matter (which needle goes through first, keeping the same needle always on the same side, angling the awl consistently).
Those conventions are real and consistently taught, and they are never written symbolically.
The only formal handle on saddle stitch anywhere is its ISO 4915 number, 201.

### 5.3 Bookbinding

Named structures with sewing diagrams, no notation.
Coptic, pamphlet stitch, long stitch, French link, Japanese stab binding, and Smyth sewing are stable names covered by guides such as [Talas](https://blog.talasonline.com/post/a-comprehensive-guide-to-bookbinding-from-sewing-signatures-to-casing-in) and the [Dartmouth book arts research guide](https://researchguides.dartmouth.edu/c.php?g=1022238&p=7404661).
Sewing station diagrams are a genuine convention within the craft.
UNCERTAIN: no formal published notation was found, and this domain got one search rather than a deep sweep.

### 5.4 Taxidermy and upholstery

UNCERTAIN, effectively unsearched.
Neither surfaced anything in the one combined query run, and neither is likely to hold a formal notation given that adjacent, better-documented crafts do not.
If someone wants to close this gap, upholstery is the more plausible of the two because it has an industrial side.

### 5.5 Fly tying

Covered in section 1.
Structured ingredient lists, prose procedures, no schema.

## 6. Best Structural Model to Borrow

Nothing transfers whole.
The recommendation is a composite, and the composite has a clear primary parent.

Primary parent: knitout.
Take its architecture wholesale, because it already solved the same meta-problem, which is representing a physical fabrication procedure as a flat, machine-independent, human-readable instruction stream.

What to take from each source:

| Source | What to take |
| --- | --- |
| knitout | File architecture: version magic line, `;;Key: value` header declaring hook and bait, flat instruction list, no flow control, `x-` extension namespace. Instruction shape `VERB DIRECTION LOCATION SUBJECT`. Face-plus-index addressing. |
| Alterovitz needle insertion plan | The per-pierce argument tuple: entry site, angle, point roll, depth. |
| Sen and Cavusoglu suture path planning | Entry point plus exit point as a pair, rather than depth alone, for through-piercings. |
| SAGES gesture taxonomy | Hierarchical vocabulary (cluster, verb, sub-verb) rather than a flat verb list. Multi-agent annotation: the supporting hand is part of the instruction. |
| Suber knot notation | Stage markers for simultaneity, role-typed nouns, terse two-letter mnemonics. |
| ISO 4915 | Symbolic identifiers for component roles, and a numbered library of named canonical patterns. |
| Suture pattern classification | Classify rigs on orthogonal axes (weedless or exposed, fixed or sliding) rather than as a flat name list. |
| CLASP semantic keypoints | Name the anatomy of the soft body so instructions survive changes in bait size. |
| Entomology pinning | The body-type to insertion-site lookup table as a first-class artifact. |
| KnitScript | Two-layer design: readable authoring DSL on top, dumb IR underneath. Build the IR first. |

The single highest-value steal is knitout's needle addressing.
`f1` and `b-2` say "front bed, position 1" and "back bed, position minus 2".
A bait needs the identical idea: a surface prefix plus a position along the body.
Something like `n` nose, `d` dorsal, `v` ventral, `l` and `r` lateral, plus a normalized position from 0 at the nose to 1 at the tail, giving `v0.35` for "belly, 35 percent of the way down".
Normalized rather than millimetres is what makes one encoded rig work across a 4 inch and a 7 inch worm.

The hook needs the same treatment, addressed by feature: `eye`, `shank@0.5`, `bend`, `barb`, `point`.
Without hook addressing, "slide the bait up the shank" is inexpressible, and that is a mandatory step in the motivating example.

Verb set sketch, drawn from the surgical and knitting vocabularies:

```text
PIERCE   enter a body at a site, at an angle, to a depth
EXIT     emerge from a body at a site
SKIN     pass tangentially just under the surface (the subcuticular case)
SLIDE    translate a body along the shank to a hook feature
ROTATE   rotate one object relative to another, in degrees, about a named axis
BURY     terminate the point subsurface without exiting
SECURE   peg, screw-lock, or bait-keeper engagement
```

`SKIN` and `BURY` are the two verbs no existing system has, and they are exactly the two the Tex-posed method needs.
The nearest prior art for `SKIN` is the subcuticular suture, which is a real named pattern for running a needle just under a surface without breaching it, so the concept is at least borrowed rather than invented.

Keep Suber's stage marker for simultaneity.
Rigging is bimanual, and the SAGES panel's rejection of dominant-instrument-only labeling is direct evidence that discarding the supporting hand loses information that practitioners consider essential.

Finally, treat named rigs the way surgery treats named suture patterns and the way ISO 4915 treats numbered stitch types: as a canonical library of macros expanding to IR sequences.
`texas`, `texposed`, `wacky`, `neko`, `bridle`, `lip`, `dorsal` become library entries, not special cases in the language.

## 7. Licensing

| Thing | Status | Can this project use it |
| --- | --- | --- |
| Suber Knot Tying Notation | No copyright or licence statement found on the page. | The notation itself (verb names, abbreviations, method) is a system, not protected expression, so adopting the vocabulary is fine. Do not reproduce his tables or prose verbatim without asking. Suber is a leading open access advocate, so an email is very likely to get a clear permissive answer. UNCERTAIN until asked. |
| knitout | Verified this run: the [GitHub repo](https://github.com/textiles-lab/knitout) contains only `README.md`, `extensions.html`, and `knitout.html`, the GitHub licence API returns null, and the spec text contains no licence or copyright string. | No licence means all rights reserved on the spec *text*. The architecture and instruction shape are ideas and are freely borrowable. Do not copy spec prose or the opcode table verbatim into project docs. CMU published it explicitly as an open interchange format and lists `knitout-feedback@cs.cmu.edu`, so asking is cheap and likely to be welcomed. |
| Craft Yarn Council symbols | Free download, published for industry-wide use. CYC retains copyright on the graphics. | Concepts and stitch names are fine. Do not redistribute their symbol artwork without checking terms. UNCERTAIN on the exact grant. |
| ISO 4915:1991 | Paid ISO standard, ISO copyright. | Referencing stitch type numbers as facts ("ISO 4915 type 201") is fine and normal. Reproducing the standard's text or drawings is not. The [open-access technical review](https://www.opastpublishers.com/open-access-articles/industrial-sewing-stitch-classification-a-technical-review-aligned-with-iso-4915-and-astm-d6193-standards.pdf) is a citable substitute for the scheme's shape. |
| JIGSAWS | Requires a data use agreement with JHU and Intuitive for the dataset. The site was Cloudflare-blocked during this run. | The 15-gesture vocabulary is published in the peer-reviewed literature and is citable and conceptually reusable. The video and kinematic data is not needed here. UNCERTAIN on exact DUA terms. |
| SAGES Delphi gesture taxonomy | [Springer, Surgical Endoscopy, 2026](https://link.springer.com/article/10.1007/s00464-026-12906-2). | Taxonomy terms are facts and can be cited and adapted. UNCERTAIN whether the article is open access; check before quoting at length. |
| OntoSPM and LapOntoSPM | Academic ontology, described as ready for wider public use. | Licence not verified. UNCERTAIN. Only relevant as prior art, not as something to import. |
| Alterovitz, Sen, Cavusoglu, Webster papers | Standard academic publication. | Mathematical formulations are not copyrightable. Cite and reuse the parameterization freely. |

Nothing in this list creates a licensing obstacle, because in every case the useful thing is a structure or a set of names rather than a copyrightable artifact.
The single action item is to email CMU's knitout contact and Peter Suber if the project wants to reproduce their material verbatim rather than merely follow their design.

## 8. Honest Assessment

No tried-and-true notation exists for this.
Not in fishing, and not anywhere close enough to adopt.

Fishing is a genuine void.
There is no formal or semi-formal notation, no structured format, no standards body attempt, and no academic attempt.
What fishing has is a stable set of names for finished rigs, a semi-standard vocabulary for where a hook ends up in a caught fish, and an enormous volume of numbered prose steps with photographs.
This is not a case of the research failing to find an obscure system.
The absence is consistent across patents, fisheries literature, instructional material, and open source, and it is what you would expect from a domain where the knowledge transfers fine by demonstration and nobody has ever needed to write it down precisely.

The adjacent domains are more interesting than expected but none is a drop-in.
Surgery has the right primitives and split them across two incompatible layers: a mature named-pattern classification that is entirely prose, and a rigorous geometric parameterization that exists only because robots needed coordinates and which has no vocabulary or sequencing.
Textiles has the right architecture in knitout and applies it to a domain with no depth axis, no rotation of one object relative to another, and no concept of piercing a body and coming back out.
Embroidery formats encode needle penetrations and explicitly treat material as a zero-thickness plane, which is the one assumption this problem cannot tolerate.
Entomology has the right conceptual object, a body-type to insertion-site lookup, and publishes it as photographs.
Robotics has no action language for deformables at all, which the surveys say outright.

So: build one.
The recommendation is knitout's architecture as the skeleton, with Alterovitz's `(y0, theta, b, d)` as the argument shape of the core `PIERCE` instruction, hierarchical verbs in the SAGES style, Suber's stage markers retained for bimanual simultaneity so the knot and rig notations stay coherent with each other, and named rigs implemented as a macro library rather than as language features.

The reason knitout is the right parent, specifically, is that it is the only system found that solved the full meta-problem rather than a piece of it.
It is machine-independent, deliberately low-level with no abstractions, human-readable, versioned, extensible through a reserved namespace, and it separates a declarative header describing the equipment from a flat imperative body describing the procedure.
A bait rig has exactly that split: the hook and bait are the equipment header, and the rigging steps are the body.
Every other candidate gives one good idea; knitout gives a working shape for the whole file.

The two verbs that have no prior art anywhere, `SKIN` and `BURY`, are unavoidable inventions.
That is the honest measure of the gap: even the closest analogue in the world, surgical suturing, only has a name for one of them, and only because subcuticular closure happens to need it.
