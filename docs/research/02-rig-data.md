# Rig Data Research

Research into source data for FISHING RIGS (terminal tackle assemblies) for the `riggermortis` open-source project.

**Bottom line up front:** no structured, machine-readable dataset of fishing rigs exists.
Not on GitHub, not on Kaggle, not in any government open-data portal, not in any formal ontology.
The rig catalog must be hand-authored.

The saving grace is that this is legally easier than it looks.
The *facts* of a rig (a Carolina rig is: line, bullet weight, bead, swivel, leader, hook) are uncopyrightable facts and methods under 17 USC 102(b).
What is copyrighted is the *expression*: the specific diagram artwork, the specific prose wording, the selection and arrangement of a particular publication's rig list.
So the correct build is: read many sources, learn the facts, author original component sequences and original SVG artwork, cite sources as references rather than as content.
Do not scrape, do not trace diagrams, do not copy prose.

---

## 1. Existing Structured Datasets

### VERDICT: NONE EXIST (this is a real finding, not a failed search)

Searched across GitHub code search, GitHub topic search, Kaggle, the FAO data catalog, ontology registries, and adjacent framings ("tackle ontology", "fishing gear taxonomy", "terminal tackle database", "fishing rig schema").

**Hard negative evidence:**

- GitHub code search over the grep.app index (>1M public repos) for the literal string `carolina rig` (case-insensitive) returned **zero results**.
  There is not a single public repository anywhere with Carolina rig data as code, JSON, YAML, or CSV.
  If any rig dataset existed in the open-source world, that string would hit.
- Kaggle: no tackle, rig, or terminal-tackle dataset surfaced. Kaggle's fishing datasets are catch statistics and fish-image classification sets.
- FAO's own data catalog exposes ISSCFG (gear classification) and nothing about terminal tackle. See section 4 for why ISSCFG is useless to us.

**What exists instead, and why none of it is a dataset:**

| Thing | What it actually is | Why it does not help |
| --- | --- | --- |
| [alaskafishcounts/adfg-sport-dataset](https://github.com/alaskafishcounts/adfg-sport-dataset) | Sport fish *counts*, 1882 to 2025, JSON | Fish, not tackle |
| [opendata/Hunting-and-Fishing](https://github.com/opendata/Hunting-and-Fishing) | Hunting/fishing *regulations* as JSON (RETIRED) | Regulations, not gear. Also dead. |
| [alzayats/fish-datasets](https://github.com/alzayats/fish-datasets) | Computer-vision training sets (COCO, Pascal VOC) | Fish images |
| [Global Fishing Watch](https://globalfishingwatch.org/datasets-and-code/) | Commercial vessel AIS activity | Industrial fleet tracking |
| [Fish Ontology (FO)](https://peerj.com/articles/3811/) | 1,830 classes: fish anatomy, morphology, ecology, developmental stages | About the animal, zero tackle classes |
| Fishery Ontology Service (FOS) | ~35,000 classes reengineered from FAO thesauri ([overview](https://www.researchgate.net/publication/228569311_A_Core_Ontology_of_Fishery_and_its_use_in_the_Fishery_Ontology_Service_Project)) | Document indexing for fishery literature, commercial framing |
| Retailer catalogs (Tackle Warehouse, TackleDirect, FishUSA, Omnia) | Structured product data behind a UI | Proprietary, ToS-protected, product SKUs not rig assemblies |

**Conclusion:** riggermortis would be the first structured rig dataset in existence.
That is a genuine reason for the project to exist, and it is also the whole cost of the project.
Budget for hand-authoring every rig.

---

## 2. Licensing Verdicts

**The single most important thing on this page:** federal *funding* does not create public domain.
Nearly every state "learn to fish" curriculum in the US is paid for by the USFWS Sport Fish Restoration program, and every one of them is still copyrighted by the state agency that wrote it.
The Minnesota MinnAqua guide is the cleanest proof (see the table).

### 2.1 The federal rule

17 USC 105 removes copyright protection only from a work "prepared by an officer or employee of the United States Government as part of that person's official duties."
Two exceptions matter to us:

1. **Contractors and grantees are not federal employees.** A work produced by a grantee is not a US Government work and can be copyrighted. NOAA's own guidance says scientific and technical reports prepared under NOAA grant programs may be copyrighted by grantees, with the Government retaining a blanket reprint license. See [NAO 205-17A](https://www.noaa.gov/organization/administration/nao-205-17a-information-access-dissemination) and the [NOAA Library copyright guide](https://library.noaa.gov/openaccesspublishing/copyright).
2. **This is exactly what Sea Grant is.** Sea Grant programs are university-based grantees (NC State, VIMS/William and Mary, UC San Diego, etc.), not NOAA staff. A Sea Grant publication is presumptively copyrighted by its host university unless it says otherwise. Treat "Sea Grant = federal = public domain" as **false**, and check each publication's own colophon.

### 2.2 Verdict table

| Publisher | License status | Can we use it? | Verbatim caveat / evidence |
| --- | --- | --- | --- |
| **US federal employee works** (NOAA staff, USFWS staff) | Public domain, 17 USC 105 | **Yes**, text and diagrams, verbatim | Must confirm the specific work was authored by federal staff, not a grantee or contractor |
| **NOAA Sea Grant program pubs** | Presumptively **copyrighted by the host university** (grantee) | **No, not by default.** Read facts only | NOAA: works of grantees "may be copyrighted by grantees" ([NAO 205-17A](https://www.noaa.gov/organization/administration/nao-205-17a-information-access-dissemination)) |
| **Minnesota DNR** | All rights reserved, permission required | **NO** | "The Minnesota Department of Natural Resources (DNR) claims copyright on all intellectual property created by the department." / "Individuals are granted permission to view or print information from this site for personal use and under fair use rights." / "For commercial uses, **including print publications, publication on other websites** or in other formats or media, permission is required." ([source](https://www.dnr.state.mn.us/aboutdnr/disclaimers_and_policies.html)) |
| **MN DNR MinnAqua Leader's Guide** (USFWS Sport Fish Restoration funded) | © Minnesota DNR, on every single page | **NO** | Page footer reads verbatim: `© 2010 Minnesota DNR • MinnAqua • USFWS Sport Fish Restoration`. Federal funding, state copyright. This is **the trap.** ([PDF](https://files.dnr.state.mn.us/education_safety/education/minnaqua/leadersguide/chapter_5/5_1_freshwater_rods_and_reels.pdf)) |
| **Wisconsin DNR** | Noncommercial general-public use, fair use only | **NO** | Material is "for noncommercial use of the general public"; outside fair use, permission must be sought, and if granted the credit line "Reproduced with permission from the Department of Natural Resources" is mandatory ([source](https://dnr.wisconsin.gov/legal/acceptableuse)) |
| **Texas Parks & Wildlife (TPWD)** | Copyrighted; noncommercial/educational reuse **by request only** | **NO** for a public repo | TPWD holds copyright; no part may be copied, reproduced, or translated without prior written consent except where noted. Even granted educational use requires the credit line "provided with permission of the Texas Parks and Wildlife Department" plus a date, **and TPWD must be sent a copy of the final product** ([policy](https://tpwd.texas.gov/site/policies/copyright-policy)) |
| **Indiana DNR** | "Copyright © 2026 State of Indiana - All rights reserved." | **NO** | Site footer ([source](https://www.in.gov/dnr/fish-and-wildlife/fishing/)) |
| **Cornell Cooperative Extension** (eCommons) | **No rights statement at all** | **NO**, and this is the worst kind of no | Verified via the DSpace REST API: the item's metadata has **no `dc.rights` field whatsoever** (keys present: contributor.author, date.*, description.abstract, format.*, identifier.uri, language.iso, publisher, subject, title, title.alternative, type). Silence defaults to all rights reserved, held by Cornell, with no contact path baked into the record. |
| **MU Extension (Missouri)** | Case-by-case, **fee may apply** | **NO** | "Permission requests to print, adapt or post MU Extension publications, webpages or graphics will be considered on a case-by-case basis, and a processing fee may apply." ([source](https://extension.missouri.edu/publications/copy)) |
| **UF/IFAS EDIS (Florida)** | All rights retained, but broad educational grant | **MAYBE, the most permissive found** | UF/IFAS retains all rights but grants permission to others to use the materials in part or in full **for educational purposes** with full credit citing the publication, source, and date ([source](https://edis.ifas.ufl.edu/es/copyright)). "Educational purposes" for a public GitHub repo is arguable but not clearly granted. Get written permission if you want their diagrams. |
| **Take Me Fishing / RBFF** | Hard no | **ABSOLUTELY NOT** | "view and download material from this site only for personal, non-commercial home use. All rights reserved." / "No Site Content may be modified, copied, distributed, framed, reproduced, republished, downloaded, **scraped**, displayed, posted, transmitted..." / users may not "republish Site Content on any Internet, Intranet or Extranet site **or incorporate the information in any other database or compilation**" ([source](https://www.takemefishing.org/site-terms-of-use/)). That last clause names our exact use case and forbids it. |

### 2.3 K-State

**UNCERTAIN.** Could not locate a K-State Research and Extension copyright/permissions page in search.
The prior run's note that "K-State has a notable copyright policy" is **unverified** here.
What is verifiable and adjacent: the Ned rig's originator, Ned Kehde, is from Lawrence, Kansas, and Kansas extension/parks material is a plausible primary source for it, which may be where that thread started.
Someone should check `ksre.k-state.edu` and `bookstore.ksre.ksu.edu` directly before relying on any K-State document.

### 2.4 The practical rule for this repo

Since the repo is public and MIT/permissive-licensed:

1. **Never copy diagrams.** Every rig SVG must be original artwork drawn from the component list, not traced.
2. **Never copy prose.** Descriptions must be written fresh.
3. **Do copy facts freely.** Component order, weight ranges, hook sizes, target species, and technique are facts and methods, not expression. 17 USC 102(b) excludes "any idea, procedure, process, system, method of operation" from copyright, and a rig is a method of operation.
4. **Cite sources per rig** in a `sources: []` field. It is honest, it is good practice, and it makes the "we learned this, we did not copy it" argument concrete if anyone ever asks.
5. **Prefer multiple corroborating sources per rig.** A rig described identically by five independent authors is a fact, not one author's expression.

---

## 3. Government & Extension Sources That Actually Contain Rig Content

| Publisher | Document | Rig content | URL |
| --- | --- | --- | --- |
| Texas Parks & Wildlife | *Fish Texas: A Basic Guide for the Beginning Angler* (PWD BK K0700-0639) | Verified by extracting the PDF: sections on Hooks, Line, Sinkers, Bobbers. Covers split shot, slip-cork rigging, sinker placement 4 to 8 inches above the hook. Beginner-level, **shallow**: 16 hits on "rig", only 1 on "swivel", 1 on "split shot". No named rigs like Carolina or drop shot. | <https://tpwd.texas.gov/publications/pwdpubs/media/pwd_bk_k0700_0639.pdf> |
| Minnesota DNR (MinnAqua) | *MinnAqua Leader's Guide*, Ch. 5 "Freshwater Rods and Reels" (and sibling chapters) | Full angler-education curriculum with tackle vocabulary and illustrations. Chapterized PDFs. **Copyrighted, see section 2.** | <https://files.dnr.state.mn.us/education_safety/education/minnaqua/leadersguide/chapter_5/5_1_freshwater_rods_and_reels.pdf> |
| Wisconsin DNR | Angler Education, *Fishing for Dinner* online guide | R3 angler-education material, tackle and rigging basics | <https://dnr.wisconsin.gov/sites/default/files/topic/Fishing/AnglerEd_FishingForDinnerOnlineGuide.pdf> |
| Wisconsin DNR | Angler and aquatic education program hub | Index to further guides | <https://dnr.wisconsin.gov/topic/Fishing/anglereducation/index> |
| Cornell Cooperative Extension | *Let's Go Fishing (A Fish and Fishing Project)*, Decker, Howard & Kelley, April 1991, pamphlet, 2.3 MB PDF | Subject tags: 4-H, **knots**, bait, lures, casting, pike, bass, bullhead, bluegill, trout. "Introduces fishing tackle and basic techniques." This is the strongest 4-H-lineage rig/knot source found, and it has **no license**. | <https://hdl.handle.net/1813/3306> (handle) / <https://ecommons.cornell.edu/items/85c2d331-11ff-4c30-b6e2-bdfb072c9bc1>. A second copy sits at handle `1813/9492` in the 4-H Youth Development Resources collection. |
| Illinois 4-H Extension | *Illinois 4-H Bass Fishing SPIN Club Curriculum* | Bass-specific club curriculum, likely covers Texas/Carolina rigs | <https://4h.extension.illinois.edu/sites/default/files/2022-07/Illinois%204-H%20Bass%20Fishing%20SPIN%20Club%20Curriculum.pdf> |
| National 4-H Council | *Fishing Curriculum Level 1: Take the Bait* (and higher levels) | The canonical 4-H Sport Fishing curriculum. **Sold, not free.** Commercial product, fully copyrighted. | <https://shop4-h.org/products/fishing-curriculum-1-take-the-bait> |
| NOAA Office of National Marine Sanctuaries | *Fishery Basics: Fishing Gear, Hook & Line* (Voices of the Bay education program) | Federal-employee-authored education material on hook-and-line gear, J-hooks vs circle hooks, trolling/longlining/jigging/pole-and-line. **Likely genuine public domain.** Commercial framing though, not recreational rigs. | <http://sanctuaries.noaa.gov/education/voicesofthebay/pdfs/hookandline.pdf> |
| New Hampshire Fish & Game (via eRegulations) | *Saltwater Rigging Basics* | Actual named saltwater rigs in a state regulations digest | <https://www.eregulations.com/newhampshire/fishing/saltwater/saltwater-rigging-basics> |
| Oregon Dept. of Fish & Wildlife | *Gearing up for fishing* | Tackle basics | <https://myodfw.com/articles/gearing-fishing> |
| South Carolina DNR | Aquatic Education / Family Fishing Clinics | Teaches knot tying and rod/reel rigging. Program pages, curriculum PDFs not directly located. | <https://www.dnr.sc.gov/aquaticed> |
| USFWS | Sport Fish Restoration program | The *funding* behind most of the above. Program page only, no rig content. | <https://www.fws.gov/program/sport-fish-restoration> |

**Honest note on this table:** none of these are deep enough on their own.
State agency material is written for a 10-year-old holding a Zebco and stops at "hook, line, sinker, bobber."
The genuinely detailed rig knowledge (drop shot vs Neko vs Tokyo, when a Mojo beats a Carolina) lives in enthusiast media (Wired2Fish, BassResource, In-Fisherman, Salt Strong), YouTube, and manufacturer marketing.
All of it is copyrighted, all of it is facts we may learn from, none of it is content we may copy.

Sea Grant specifically: repeated searches did not surface a Sea Grant publication that is primarily a rig diagram catalog.
The [Sea Grant Collection in the NOAA repository](https://repository.library.noaa.gov/cbrowse?pid=noaa%3A11&parentId=noaa%3A11) is the place to dig further if someone wants to keep looking. **UNCERTAIN** whether a good one exists there.

---

## 4. Formal Gear Taxonomies

### VERDICT: They stop several orders of magnitude above us, so do not adopt, invent

### 4.1 FAO ISSCFG

The International Standard Statistical Classification of Fishing Gear, adopted at the 10th CWP Session (1980), current revision adopted at the 25th CWP Session (FAO 2016).
Landing page: <https://www.fao.org/cwp-on-fishery-statistics/handbook/capture-fisheries-statistics/fishing-gear-classification/en/>
Full standard PDF: <https://openknowledge.fao.org/server/api/core/bitstreams/fe19f36f-7e88-4304-b7e0-aad1e3bdd7e6/content>
Data catalog entry: <https://data.apps.fao.org/catalog/dataset/cwp-isscfg>

I extracted the actual code table from the PDF. Here is verbatim what the standard says about our entire domain:

```text
09.0.0  HOOKS AND LINES
09.1.0  Handlines and pole-lines (hand-operated)     LHP
09.2.0  Handlines and pole-lines (mechanized)        LHM
09.3.0  Set longlines                                LLS
09.4.0  Drifting longlines                           LLD
09.5.0  Longlines (not specified)                    LL
09.6.0  Trolling lines                               LTL
09.9.0  Hooks and lines (not specified)              LX
...
25.0.0  RECREATIONAL FISHING GEAR                    RG
```

**That is it.**
All of recreational angling, every rig riggermortis will ever model, is a single leaf code: `25.0.0 / RG`.
The claim on the FAO landing page that the classification "applies to commercial, subsistence and recreational fisheries" is technically true and practically meaningless: recreational gear gets exactly one bucket with no subdivision.

The standard does say it "provides for national or regional variations to be included at sub-levels of the classification," which is FAO's way of saying *you're on your own below this line*.
That is an invitation to invent, not a taxonomy to adopt.

**Verdict: ISSCFG is not adoptable.** At most, tag rigs with `isscfg: "25.0.0"` as a courtesy cross-reference so the data is nominally interoperable with fisheries systems. It carries zero information for us.

FAO gear fact sheets: <https://www.fao.org/fishery/en/geartype/search> (same granularity, illustrated commercial gear).
Companion reference: Nédélec & Prado, *Definition and classification of fishing gear categories*, FAO Fisheries Technical Paper 222 Rev.1 (1990), <http://www.fao.org/3/a-t0367t.pdf>.

### 4.2 AGROVOC

Queried the live AGROVOC SKOS API directly rather than trusting search snippets.

`fishing gear` (`c_2946`) has exactly these 13 narrower concepts:

```text
fishing rods (c_f6510256)          fishing lines (c_2949)
fishing dredges (c_4079990c)       electrified gear (c_746e8fe9)
fishing nets (c_2951)              otter boards (c_4a6d0936)
fishing pots (c_2953)              deck equipment (c_2148)
bycatch excluder devices (c_fe6a0492)   wounding gear (c_9997d1df)
fish aggregating devices (c_265f64c3)   grappling gear (c_01d4db4f)
abandoned, lost or discarded fishing gear (c_9d5917cc)
```

`fishing lines` (`c_2949`) has exactly **one** narrower concept:

```text
hooks (c_3bdbab0f)
```

And `hooks` is a leaf.
The prior run's finding that "hooks is under fishing lines" is **CONFIRMED**, and it is the entire depth of the tree.

Direct concept searches for the rest of our component vocabulary returned **nothing**:

| Term searched | AGROVOC result |
| --- | --- |
| sinkers | no concept |
| swivels | no concept |
| floats | no concept |
| fishing weights | no concept |
| bobbers | no concept |
| leaders | `c_15550` exists, but almost certainly the *person who leads* sense, not fishing leader. **UNCERTAIN**, do not rely on it. |

Browse pages: <https://agrovoc.fao.org/browse/agrovoc/en/page/c_2946> (fishing gear), <https://agrovoc.fao.org/browse/agrovoc/en/page/c_2949> (fishing lines).

**Verdict: AGROVOC is not adoptable either.** It has two relevant URIs total (`fishing lines`, `hooks`) and no concept for a sinker. We can cite those two as `sameAs`/`broader` links for provenance, and that is the entire extractable value.

### 4.3 Fishing gear ontology (ODP)

The prior run reported an ontology-design-pattern URL for a fishing gear ontology that is **dead**.
Not re-verified in this run to save budget. Treat as **UNCERTAIN but probably dead**; ODP-portal link rot is endemic and nothing else in this search surfaced a live tackle ontology.

### 4.4 Consequence for the schema

Nobody has modeled this. That means:

- There is no vocabulary to conform to, so design the schema for *rendering an animated SVG* and for *querying by species/water/depth/cover*, which are the actual use cases.
- Emit stable IDs of our own (`rig:carolina`, `component:bullet-weight`) and treat any external URI as an optional annotation.
- Publishing this as SKOS/JSON-LD later would make riggermortis the de facto vocabulary, since there is no incumbent. That is worth designing toward even if v1 is plain JSON.

---

## 5. Component Vocabularies & Sizing Standardization

### VERDICT: Sizing is genuinely, structurally unstandardized, so the schema must treat every size as an opaque, manufacturer-scoped label

### 5.1 Hook sizes: NOT standardized (confirmed)

The aught system (size 32 up to size 1, then 1/0 up to 20/0) is a *shared convention*, not a *standard*.
There is no ISO, ANSI, or trade-association specification defining what "2/0" measures.
Multiple independent sources state this plainly:

- Hook sizing is not standardized across manufacturers; a Mustad 7/0 differs in size from an Eagle Claw 7/0, though each manufacturer is internally consistent across its own range ([go-saltwater-fishing.com](https://www.go-saltwater-fishing.com/fishing-hook-sizes.html)).
- A Gamakatsu 3/0 wide-gap may be physically a different size from an Owner 3/0 wide-gap; compare physical dimensions, not size numbers, when switching brands ([knots.fish](https://knots.fish/guides/fishing-hook-size-guide/)).
- Size interacts with *pattern*, not just brand: a 2/0 octopus hook is much smaller than a 2/0 worm hook ([Tackle Warehouse gear guide](https://www.tacklewarehouse.com/bass-fishing/gear-guides/hook-sizing-gear-guide-for-bass-fishing.html)).

So there are **two** independent axes of ambiguity: manufacturer *and* hook pattern.
`2/0` alone is not a size; `Gamakatsu 2/0 EWG` is.

Note also that "wire size" (1X strong, 2X strong, 3X fine, etc.) is a *relative* modifier meaning "one wire gauge heavier/lighter than standard for this pattern," and "standard for this pattern" is itself manufacturer-defined. It compounds, it does not resolve, the ambiguity.

### 5.2 Split rings: confirmed mess, and worse than reported

The prior run's finding is **confirmed and extends further**: the modifier suffix system means a single manufacturer ships multiple products at the same nominal size number with wildly different ratings.
Rosco alone sells size 3 and size 3H:

| Product | Approx. test rating |
| --- | --- |
| Rosco Stainless Split Ring **size 3** | ~25 lb |
| Rosco Stainless Split Ring **size 3H** | ~65 lb |

Rosco's line runs to 4XXH and 6XH. ([Rosco size listings on njtackle.com](https://www.njtackle.com/terminal-tackle/rosco-stainless-steel-split-rings/), [Rosco product page](https://www.roscotackle.com/product/rosco-split-ring/))

Industry practice is that each manufacturer publishes its own chart and prints size plus pound-test on the packaging ([Sport Fishing Mag](https://www.sportfishingmag.com/attention-to-detail/)).
That is an admission there is no cross-brand meaning to a size number.
The single reliable, comparable field is **pound test**, not size.

### 5.3 Swivels, sinkers, floats, leader

**UNCERTAIN / lower confidence, not directly verified in this run.** Flagged for follow-up:

- **Swivels** use an inverted numbering convention in many catalogs (higher number equals smaller swivel), overlaid with per-manufacturer ratings. Same structural problem as split rings. The comparable field is again pound test.
- **Sinker styles** are named descriptively and the names are broadly consistent in the wild (egg, bullet/worm, bank, pyramid, split shot, no-roll, walking, drop-shot/cylinder, pencil), but no body defines them. Weight in ounces/grams is the only unambiguous field. TPWD's *Fish Texas* confirms the practical range: "from BB split shot to five pounds."
- **Line and leader** are sold by nominal pound test, which is not standardized either: most line breaks *above* its rating. The exception worth encoding is **IGFA class line**, which is specified to break *at or below* its stated rating for record eligibility. **UNCERTAIN**, verify against IGFA rules before relying on it.
- **Beads, floats, snaps, clevises**: no vocabulary found anywhere. Purely descriptive.

### 5.4 Direct schema consequence

Do **not** model size as a scalar.
Model it as a structured, scoped value:

```jsonc
{
  "type": "hook",
  "pattern": "ewg_worm",          // our own controlled vocab
  "size": {
    "system": "aught",            // aught | metric | mm | oz | g | free
    "label": "3/0",               // opaque string, never compared numerically
    "manufacturer": "Gamakatsu"   // optional; a size means nothing without it
  },
  "measured": { "gap_mm": null, "shank_mm": null },  // optional ground truth
  "rating_lb": null               // the ONLY cross-brand-comparable field
}
```

Rules that fall out of the research:

1. Never sort, compare, or arithmetic on a size label across manufacturers.
2. `rating_lb` is the canonical comparable field for rings, swivels, snaps, line, and leader.
3. Weight in grams/ounces is the canonical field for sinkers and jigheads.
4. For rendering, the SVG needs *relative visual proportion*, not a real size. Store an explicit `render_scale` hint rather than deriving geometry from a size label, because you cannot derive it.

---

## 6. Regional Naming & Trademarks

### 6.1 Documented same-rig-different-name collisions

| Canonical | Also called | Evidence |
| --- | --- | --- |
| **Carolina rig** | Fish finder rig | Widely used interchangeably. The genuine distinction: a Carolina rig uses an in-line sinker (egg, barrel, bullet) threaded on the main line; a fish finder rig uses a heavy sinker on a **sliding clip/slider**. Functionally the same presentation. ([castandspear.com](https://castandspear.com/carolina-rig/)) |
| **High-low rig** | Dropper loop rig, flapper rig | "The High-Low Rig is also known as a Dropper Loop Rig or Flapper Rig." It is the same setup as a three-way/dropper-loop rig with a second hook added above the first. ([greatdaysoutdoors.com](https://greatdaysoutdoors.com/surf-fishing-rigs/)) |
| **High-low rig** | "Chicken rig" | Commonly asserted in surf-fishing communities. **UNVERIFIED** in this run; searches did not corroborate it in a citable source. Mark UNCERTAIN until a source is found. |

The freshwater/saltwater split is the deeper naming fault line: a bass angler's "Carolina rig" (bullet weight, glass bead, swivel, 18 to 36 inch fluoro leader, soft plastic) and a surf angler's "Carolina rig" (egg sinker, bead, swivel, mono leader, circle hook, cut bait) share a topology and share nothing else about size, material, or use.

**Schema consequence:** a rig needs `aliases: []` as a first-class field, and probably a `variant_of` relation so a bass Carolina and a surf Carolina are two records sharing one topology rather than one record with contradictory metadata.
Search must resolve aliases or users will not find anything.

### 6.2 Trademarks

| Name | Status | Evidence | Practical risk |
| --- | --- | --- | --- |
| **"THE ALABAMA RIG"**, **"A-RIG"** | **Registered trademarks, owned by Slick Lures, LLC** | [Justia trademark listing for Slick Lures, LLC](https://trademark.justia.com/owners/slick-lures-llc-2231169), both registered for fishing rigs. Inventor Andy Poss filed with USPTO 2011-02-12 and signed an OEM licensee and trademark agreement with Mann's Bait Company 2011-12-04. His *patent* attempt was abandoned after a similar 2009 filing surfaced. ([Wired2Fish profile](https://www.wired2fish.com/people-profiles-and-passings/guest-blog-the-man-behind-the-alabama-rig), [KentuckyAngling](https://www.kentuckyangling.com/magazine/the-alabama-rig-exposed/)) | **Highest of the three.** Prefer the generic **"umbrella rig"** as the canonical record name and list "Alabama rig"/"A-rig" as aliases. |
| **"Tokyo Rig"** | VMC/Rapala markets "Tokyo Rig®" with the registered symbol on product lines (Heavy Duty Wide Gap, Finesse Neko, Heavy Duty Worm, Heavy Duty Flippin') | [rapala.com VMC Tokyo Rigs](https://www.rapala.com/us_en/vmc/terminal-tackle/tokyo-rigs). Direct USPTO/Justia verification **blocked** (Justia returned HTTP 403; the USPTO search API endpoint no longer resolves). | **UNCERTAIN but assume registered.** The ® is being asserted publicly by a company that owns [dozens of registered marks](https://trademarks.justia.com/owners/rapala-vmc-oyj-42127). Consider a generic canonical name (e.g. "drop-arm rig" / "wire-drop rig") with "Tokyo rig" as alias. |
| **"Ned Rig"** | **No trademark registration found.** Named after Ned Kehde of Lawrence, Kansas, who popularized it. Z-Man markets heavily around it but its registered marks are for products (Z-MAN, TRD, ShroomZ), not the rig name. | [Z-Man Ned Rig page](https://zmanfishing.com/pages/ned-rig), [Z-Man's Justia trademark list](https://trademarks.justia.com/owners/z-man-fishing-products-inc-1191698) | **Low.** Appears to have become generic. Still **UNCERTAIN**, not exhaustively searched. |
| Carolina, Texas, drop shot, wacky, Neko, three-way, slip bobber, Sabiki, Mojo, high-low | No evidence of registration; all long-established generic technique names | | Low |

**The legal reality, stated plainly:** trademark protects *source identification in commerce*, not descriptive reference.
Nominative fair use lets us say "this is the rig commonly called the Alabama rig."
What we must not do is name a product, brand, or the project itself after someone's mark, or imply endorsement.
Since riggermortis is a reference dataset and not a tackle product, the exposure is low.
But **choosing generic canonical names and demoting trademarked names to aliases costs nothing and removes the risk entirely.**
Do that.

Also worth knowing: **rigs get patented.** Example: [US 8,407,929 "Surf fishing rig"](https://patents.justia.com/patent/8407929).
Patents do not restrict *describing* a rig, only making/selling it, so this is irrelevant to the dataset. Noted so nobody panics later.

---

## 7. Honest Assessment

**There is no dataset. Build it by hand. That is the finding.**

Six independent search strategies (GitHub code search over a million-repo index, GitHub topic search, Kaggle, the FAO data catalog, ontology registries, and four adjacent framings) produced zero structured rig data.
A literal code search for `carolina rig` across public GitHub returned **nothing**.
That is about as strong a negative as this kind of research produces.

**The formal standards are a dead end, definitively.**
FAO ISSCFG collapses all of recreational angling into one code, `25.0.0 / RG`, with no subdivision.
AGROVOC's entire tackle vocabulary is two URIs: `fishing lines` and `hooks`, and it has no concept for a sinker.
Neither is worth adopting. Cite them as courtesy annotations and move on.
The upside: there is no incumbent vocabulary to conflict with, so riggermortis can define one.

**The licensing situation is better than the prior run's "critical trap" framing suggests, but only if the project takes the right path.**
The trap is real and it is this: US federal *funding* is not the same as federal *authorship*.
Sport Fish Restoration pays for nearly every state angler-education curriculum in the country, and every one is still © the state agency.
MinnAqua stamps `© 2010 Minnesota DNR • MinnAqua • USFWS Sport Fish Restoration` on every page.
The same trap catches Sea Grant: those are university grantees, not NOAA employees, so 17 USC 105 does not reach them.
Anyone who assumes "government fishing PDF equals public domain" will put infringing artwork in a public repo.

But the trap only bites if you copy *content*.
Rig topology is a method of operation and is excluded from copyright by 17 USC 102(b).
Draw the SVGs from scratch, write the prose from scratch, cite sources as references, and the licensing problem evaporates.

**What that costs.** Every rig is: read three to five sources, reconcile their disagreements, decide the canonical component order, decide the canonical name and its aliases, author the SVG.
That is real work and it is the bulk of the project.
It does not parallelize into a scraper. Accept that up front.

**Best citable sources, ranked:**

1. **NOAA federal-employee education material** (e.g. the Voices of the Bay *Hook & Line* PDF). Genuinely public domain, but commercial-gear framed and thin on recreational rigs.
2. **Cornell CCE *Let's Go Fishing*** (1991). The best 4-H-lineage tackle/knot pamphlet found, freely downloadable, and it has **no rights statement at all**, which is the most dangerous kind of source. Learn from it, cite it, copy nothing.
3. **TPWD *Fish Texas*.** Verified to contain hooks/line/sinkers/bobbers sections. Beginner-level and shallow: one mention of "swivel" in the whole document. Copyrighted with an onerous permission process.
4. **State angler-education curricula** (MinnAqua, WI *Fishing for Dinner*, IL 4-H Bass SPIN). Structured, pedagogically ordered, all copyrighted.
5. **Enthusiast media and manufacturer technique pages** (Wired2Fish, BassResource, Z-Man, VMC/Rapala). By far the deepest and most current on modern finesse rigs. All copyrighted. Facts only.

**Three schema decisions this research forces:**

1. **Size is never a scalar.** Store `{system, label, manufacturer}` plus an optional `rating_lb`. Hook sizes vary by brand *and* by pattern; split ring "size 3" and "size 3H" from one maker differ by 40 lb of rating. Never compare size labels numerically.
2. **`aliases[]` is mandatory, and variants need their own records.** Carolina/fish-finder and high-low/dropper-loop/flapper are documented collisions, and a bass Carolina and a surf Carolina share only a topology.
3. **Canonical names should be generic; trademarked names go in aliases.** "Umbrella rig" canonical, "Alabama rig" alias, because Slick Lures, LLC owns that mark. Same treatment for Tokyo Rig®. Costs nothing, removes the risk.

**Open items, honestly flagged:**

- K-State's copyright policy: **not verified.** The prior run's note is unconfirmed here.
- "Tokyo Rig" registration details: **not verified.** Justia 403'd and USPTO's search API endpoint is gone. Needs a manual TSDR lookup.
- Whether a genuinely good Sea Grant rig publication exists: **unknown.** The [NOAA Sea Grant Collection](https://repository.library.noaa.gov/cbrowse?pid=noaa%3A11&parentId=noaa%3A11) was not exhaustively browsed.
- Swivel, sinker, float, and leader vocabularies: characterized from general knowledge and one PDF, **not systematically verified** the way hooks and split rings were.
- The dead fishing-gear ontology URL from the prior run: not re-checked.
