# Prior Art & Competitive Landscape

Research pass: 2026-08-10.
Scope: prior art for a mobile-first, open-source, programmatically-rendered reference for recreational terminal-tackle rigs.

Claims tagged `UNCERTAIN` were not verifiable within this pass.
Everything else carries a citation.
Where a claim came from inspecting page source rather than marketing copy, the method is stated so it can be re-run.

## The Gap

Four questions, four explicit answers.

| # | Thing | Exists? | Closest thing that does exist |
| --- | --- | --- | --- |
| a | Structured, queryable **rig database** | **No** | MyRigs (a static 140-rig picture book, abandoned 2020); Wired2Fish `/fishing-rigs` (a WordPress blog archive) |
| b | Conditions-to-rig **recommendation engine** | **Partially, and not in this shape** | Deep Dive, Fishbox, Rigline all recommend *lure / depth / location*, not *rig*, and all are closed and commercial |
| c | **Configurator** rendering arbitrary user-composed rigs | **No** | Nothing found. Every product is a fixed catalog |
| d | **Openly licensed rig dataset** | **No** | Nothing. The entire practical-knot corpus on GitHub is 1 repo with 0 stars |

### (a) Structured, queryable rig database: NO

The closest artifact is **MyRigs (Deep Sea Fishing Rigs)**, iOS, $3.99, "over 140 fishing rigs across 7 categories," shown as static diagrams.
It is a fixed catalog, not a database: no query, no API, no schema, and it was last updated **May 26, 2020** with **3 total ratings**.
<https://apps.apple.com/us/app/id1378115973>

The web equivalent is **Wired2Fish's `/fishing-rigs` section**, which is a WordPress category archive.
Fetching <https://www.wired2fish.com/fishing-rigs> and extracting links yields **16 rig articles on page 1** plus a `page/2`, so the total corpus is roughly 30 to 40 long-form prose posts.
Each is a hand-written article ("The Texas Rig: How to Rig and Fish", "Neko Rig: Complete Rigging and Fishing Guide").
There is no shared schema across them, no way to ask "show me every rig that uses a sliding weight," and no machine-readable representation of any rig.

Catch-logging apps (Fishbrain, FishAngler, ANGLR) are **social/logging products**, not rig references.
Their content model is catches, waypoints, and species, not tackle topology.
`UNCERTAIN`: none of the three was inspected directly in this pass, so a buried rig-reference section cannot be fully excluded, but none surfaced in any "best fishing apps" roundup as a rig reference.

**Verdict: no structured, queryable rig database exists. This is a genuine hole.**

### (b) Conditions-to-rig recommendation engine: PARTIALLY, but aimed elsewhere

Recommendation engines do exist in fishing, and they are better funded than anything this project will ship.
But they answer a different question.

- **Deep Dive (bass)** ships a "Bait Tool & Forecast" that recommends **best lure, depth, and presentation** from user inputs plus live weather, and a satellite-derived **Water Clarity Map covering 170+ lakes** used to pick lure color and depth. <https://apps.apple.com/us/app/deep-dive-bass-fishing-app/id1549715715>
- **Fishbox** ships an AI catch assistant that answers questions about "lure, rig, timing, and location," over 100,000+ waterbodies, with a PRO tier adding depth, submerged vegetation, and real-time water clarity. <https://fishbox.com/blog/best-fishing-apps>
- **Rigline** does estuary-specific hotspot scoring from tides, wind, salinity, clarity, habitat and structure, per target species. <https://riglineoffshore.com/>

Three things separate these from what riggermortis proposes:

1. **They recommend the lure and the spot, not the rig.** The output is "throw a green pumpkin craw at 12 feet," not "use a Carolina rig with a 3/4oz bullet weight and an 18 inch fluorocarbon leader, because you need bottom contact while keeping the bait above the grass."
2. **They are unexplainable by design.** They are prediction products. None exposes the reasoning chain, which is precisely the thing an educational reference is for.
3. **They are all closed and commercial**, several subscription-gated.

`UNCERTAIN`: Fishbox's marketing copy uses the word "rig," and this pass did not verify whether its AI assistant actually returns rig topology or just repeats lure advice. Given it is an LLM wrapper over general fishing text, assume the latter until tested.

**Verdict: a conditions-to-rig engine that explains *why* does not exist. Conditions-to-lure engines do exist and are well capitalized. Do not claim the "recommendation engine" as the differentiator, because it is the crowded part of the map.**

### (c) Configurator rendering arbitrary user-composed rigs: NO

Nothing found, in any product, at any price, on any platform.
Every rig reference located in this pass (MyRigs, Wired2Fish, FishUSA, Cast and Spear, qwikfishing, manufacturer content) is a **fixed list of preset rigs presented as pre-made artwork**.

This follows directly from how the incumbents build assets: if your rig diagram is a photograph or a hand-drawn illustration, you physically cannot render a rig that nobody drew.
The absence of a configurator is not an oversight in the market, it is a **consequence of the asset pipeline everyone chose**.

**Verdict: no configurator exists. This is the strongest and most defensible differentiator in the whole project.**

### (d) Openly licensed rig dataset: NO, emphatically

Searched the GitHub API directly:

| Query | Total results | What is actually there |
| --- | --- | --- |
| `topic:knots` | 7 repos | All **mathematical knot theory** (Seifert surfaces, knot invariants, lattice knots). Zero practical knot tying. |
| `fishing rig in:name,description` | 21 repos | All noise: blockchain "fishing rights" contracts, a Discord fishing bot, two empty repos. |
| `knot tying animation in:name,description` | **1 repo** | `dingonzo/animated-knots-website`, 1 star, no license. |

The single most substantial practical-knot repository located is **`antontsv/knots`**: 0 stars, MIT, described as "Mostly notes about climbing knots."

Meanwhile the two authoritative sources are explicitly closed:

- Animated Knots: "Copyright Protected... Please respect this. Many years of work underlie this website," permitting only linking and the use of *a single* image as a link. <https://www.animatedknots.com/faq>
- Knots 3D: proprietary commercial app, no data export.

**Verdict: there is no open rig or knot dataset. Not a small one, not a stale one, not a partial one. The commons is empty.**

## Knot References

| Product | Coverage | Rendering technique (verified) | Business model | Mobile | Notably bad at |
| --- | --- | --- | --- | --- | --- |
| **Animated Knots by Grog** | ~200 knots, large fishing section | **JPEG photo sequence**, ~13 to 19 frames, swapped into one `<img>` by hand-rolled JS. No video, no canvas, no 3D. | Display ads (Deployads/Snigel) + paid iOS/Android/macOS apps | Responsive site, works fine, ad-heavy | Cannot render anything not photographed; knots only, no rigs; no conditions context; legally unforkable |
| **Knots 3D** | 220+ knots incl. fishing | **Genuine interactive 3D models.** 360 degree rotation, scrubbable animation, photorealistic rope textures | **$5.99 one-time**, no ads, no IAP, works offline | Native app, strong reviews | Knots only, not rigs; app-only, no web; closed source; no "which knot for what" reasoning |
| **netknots.com** | Large fishing knot library | **Adobe Animate (ex-Flash) HTML5 Canvas export, driven by CreateJS.** Vector animation, but from a frozen authoring pipeline. See detail below. | Ads / affiliate | Dated, fixed-size iframe | Frozen 2015-era pipeline; 440x440 fixed canvas in a 500x500 iframe is hostile on mobile; no source of truth beyond the `.fla` |
| **Pro Knot** | Fishing/boating/climbing card sets | **Not a reference site.** proknot.com is a **Shopify storefront** (`cdn/shop/files/` assets, zero video/canvas/iframe) selling waterproof plastic knot cards | E-commerce, physical product + companion app | n/a, it is a store | Not a competitor to a web reference at all; sells laminated cards |

### Animated Knots by Grog: the definitive read

This is the market leader and the most important competitor to understand, because the project's core bet is programmatic rendering.

Verified by fetching <https://www.animatedknots.com/palomar-knot> and inspecting the raw DOM:

```html
<img ID='Knot' class='rotation0' style='width: 100%;'
     src="https://www.animatedknots.com/photos/palomar/palomarR13.jpg"
     alt="Palomar Knot, Step-by-Step Animation">
```

Findings, each independently checkable:

- **Zero `<video>` elements. Zero `<canvas>` elements.** The page has 36 `<img>` tags and no rendering surface of any kind.
- Playback is bespoke JavaScript mutating `img.src`. The controls call `nextFrame()`, `prevFrame()`, `stopGo()`, `setSpeed()`, `setLoop()`, and `setReflection(1|2)`.
- Frame assets follow `<knot><handedness><n>.jpg`. Probing the CDN directly: `palomarR1.jpg` returns 200, `palomarR5` 200, `palomarR13` 200, `palomarR20` 404. So one knot is roughly **13 to 19 hand-shot photographs**.
- The left-handed/right-handed toggle (`setReflection`) applies a **CSS mirror transform to a bitmap** rather than loading a second photo set. That is the tell: there is no model underneath, only a flat image.
- Provenance confirms it. The originals were shot by Alan "Grog" Grogono himself "on his kitchen table with mostly natural light," and re-shot in 2010 for the iPhone app. <https://www.animatedknots.com/acknowledgements>
- Monetisation is display ads. The page loads `//tags-cdn.deployads.com/a/animatedknots.com.js` **eight times**.

**Why this matters more than anything else in this document:** the category leader in "animated knots" is a photo flipbook.
Every single asset is a manual photo shoot.
That means it structurally **cannot** render a knot it has not photographed, cannot change line diameter or colour, cannot change hook style, cannot compose two knots into one rig, cannot adapt to a phone-sized viewport beyond scaling a JPEG, and cannot be forked or extended by anyone.

A programmatic renderer beats it on exactly those axes.
It does **not** beat it on visual realism, trust, SEO, or twenty years of accumulated authority.
Be honest about which fight you are picking.

### Knots 3D: the real technical competitor

Knots 3D is the one incumbent that is **not** a flipbook.
It uses real interactive 3D geometry: rotate any knot 360 degrees, scrub the tying animation, photorealistic rope textures.
<https://knots3d.com/> and <https://knots3d.com/faq>

It is $5.99 one-time with no ads, no in-app purchases, and full offline operation, which is an unusually pro-user model and makes it hard to out-compete on goodwill.

The openings against it are narrow but real:

1. It is **app-only**. There is no shareable web URL for a knot, which kills the "someone posts a link in a forum thread" distribution loop entirely.
2. It covers **knots, not rigs**. A knot is a single connection. A rig is a composed assembly of line, terminal tackle, weights and hooks. Knots 3D has nothing to say about a Carolina rig.
3. It is **closed**, so nothing can build on it.

`UNCERTAIN`: the specific engine and pipeline behind Knots 3D (whether the models are hand-authored in a DCC tool and baked, or procedurally generated) is not published. The observable behaviour is consistent with pre-authored, baked animations rather than a live physical solver.

### netknots: the one vector-animated incumbent, and why it is a cautionary tale

netknots is the only incumbent found that animates knots as **vector graphics on a canvas** rather than as photographs.
That makes it the closest thing to prior art for this project's technical approach, so it is worth being precise about.

Fetching the animation iframe target directly, <https://www.netknots.com/application/themes/netknots/swf/PALOMAR_KNOT_net.html>, shows:

- One `<canvas>` element, zero `<video>` elements.
- `<script src="https://code.createjs.com/createjs-2015.11.26.min.js">`
- Adobe Animate export markers throughout (`AdobeAn`, `lib.Stage`, `createjs.Stage`)
- `lib.properties = { id: '4D0928AFB811E646B5E67CD5A6A6FE49', width: 440, height: 440, fps: 12, ... }`

So: a **Flash `.fla` timeline animation, exported by Adobe Animate to an HTML5 Canvas build running on CreateJS/EaselJS, pinned to a 2015 CDN, hard-coded to 440x440 at 12fps.**

The lesson is sharper than "netknots is old."
It proves vector knot animation genuinely works and looks fine.
It also shows the failure mode: the artwork is authored **by hand in a GUI timeline tool**, so the real source of truth is a `.fla` file, not data.
The result is the same structural dead end as the photo pipeline.
It cannot generate a knot nobody drew, cannot reflow to a phone viewport, cannot be composed into a rig, and is welded to a dead authoring tool and a decade-old runtime.

**The differentiator is not "vector instead of photos." netknots already did vector. The differentiator is "generated from data instead of authored by hand."** That distinction should be stated precisely in the project's own positioning, because "animated diagrams" alone describes three existing competitors.

## Rig References

| Source | What it actually is | Structured? | Verdict |
| --- | --- | --- | --- |
| **Wired2Fish** `/fishing-rigs` | WordPress category archive, ~16 posts on page 1 plus page 2 | No | The best free rig content on the open web, and it is still just blog posts |
| **BassResource** | Forum plus article library | No | Deep, unstructured, dated presentation |
| **Salt Strong** | Membership/coaching business | No | Paywalled, video-first, inshore-focused |
| **Fishbrain / FishAngler / ANGLR** | Catch logging, social, waypoints | Rig content is not the product | Not rig references |
| **Rapala / Berkley / Bass Pro / Tackle Warehouse** | Manufacturer and retailer marketing content | No | Content marketing designed to sell the SKU shown in the diagram, inherently non-neutral |
| **MyRigs (app)** | 140 static rig diagrams | No | Abandoned 2020, 3 ratings, but the closest thing to a rig catalogue in existence |
| **Wikipedia** | Individual articles: Texas rig, Carolina rig, wacky rig, chod rig | No | Encyclopedic stubs, no diagrams worth anything, no cross-rig structure |

**The honest state of the art: there is no single comprehensive rig reference anywhere.**
It is genuinely scattered blog posts, YouTube videos, forum folklore, and manufacturer marketing.
Every rig article is written independently, so terminology drifts, the same rig appears under three names, and nothing is comparable across sources.

That is real, and it is the strongest supporting evidence for gap (a) and (c).

The caveat worth stating plainly: **content being scattered is not proof that anyone wants it consolidated.**
Plenty of scattered domains stay scattered because the scattered form is good enough. See Community Signal below.

## Design References Worth Stealing

### Agrawala et al., "Designing Effective Step-By-Step Assembly Instructions" (SIGGRAPH 2003)

<https://graphics.stanford.edu/papers/assembly_instructions/assembly.pdf>

The single most directly applicable piece of prior art in this document, and it is not about fishing at all.
The principles were derived empirically: subjects assembled a TV stand and drew their own instructions, a second group ranked those instructions, a third group was timed and error-counted using the ranked instructions.
Highly rated instructions measurably produced **faster assembly and fewer errors**.

The paper's principles, with the direct translation to rigs:

**1. Hierarchy and grouping of parts.**
"People think of assemblies as a hierarchy of parts... parts that are disjoint are more likely to be segmented... people prefer that parts within a group are added to the assembly at the same time, or in sequence one after another."
*Steal this:* model a rig as a hierarchy (mainline, leader, terminal group, hook group), not a flat list of steps. Group operations on the same component.

**2. Hierarchy of operations, and two levels is usually enough.**
"a two-level hierarchy (significant parts and less important parts + fasteners) is common for many build-at-home objects."
*Steal this:* split rig steps into **significant components** (hook, weight, swivel, leader) and **fasteners** (the knots and crimps that attach them). Resist deeper nesting. Two levels is the empirically supported default.

**3. Step-by-step beats one diagram, and one significant part per step.**
"people prefer instructions that present the assembly operations across a sequence of diagrams rather than a single diagram... people generally prefer that each diagram show how to attach only one significant part at a time. However, each diagram will usually show multiple non-significant part attachments."
*Steal this:* one terminal component per step, with its knot shown in the same step rather than as a separate step. This is a concrete, testable layout rule.

**4. But do not over-fragment.**
"each diagram should also present as much information as possible. If instructions are split across too many diagrams, they become tedious to use... A better approach is to skip repetitive operations after they have been presented in detail a few times."
*Steal this:* directly relevant, because rigs repeat knots constantly. Show the improved clinch in full the first time; collapse it to a reference thereafter. Do not make the user watch the same knot animation five times in one rig.

**5. Action diagrams beat structural diagrams.**
"Structural diagrams present all the parts of the assembly in their final assembled positions; users must compare two consecutive diagrams to infer which parts are to be attached. Action diagrams spatially separate the parts to be attached from the parts that are already attached and use guidelines to indicate where the new parts attach... We found that action diagrams are superior to structural diagrams."
*Steal this:* this is the highest-value single takeaway. Do not render the finished rig and highlight a part. **Separate the incoming component from the assembled portion and draw a guideline to where it goes.** This is trivially expressible in SVG and is exactly what a programmatic renderer can do that a photograph cannot.

**6. Orientation and natural views.**
"Most objects have a set of natural orientations or preferred views... These orientations maximize the number of important features that are visible."
*Steal this:* rigs have a canonical orientation (line entering from the top, terminal end down, matching how it hangs in water). Pick it and hold it. Reorient only when a step genuinely needs it, as the paper's bookcase does.

**7. Visibility is the strongest principle.**
"Perhaps the strongest design principle is that all the new parts added in each step of the assembly must be visible... While the new parts have to be visible, the parts attached in earlier steps should also be visible to provide context."
*Steal this:* every step must show both the new component and enough of the already-built rig for context. On a phone, this is a hard constraint on zoom and framing, and it should drive the camera/viewBox logic rather than being an afterthought.

**8. Symmetry exemption.**
"maintaining visibility for all parts in a symmetric group is less important... it is usually enough that at least one part in the group is visible."
*Steal this:* treble hooks, multi-dropper rigs, tandem hooks. Show one clearly, imply the rest.

The paper also validates the entire architecture: it separates a **planner** (what to show in each step) from a **presenter** (how to draw it), taking as input part geometry, orientations, groupings, and ordering constraints.
That is close to a blueprint for the rig-rendering pipeline, and it is 23 years old and freely readable.

### Other domains worth mining

`UNCERTAIN` on specifics for the items below, which were not individually verified in this pass. They are listed as leads, not findings.

- **Origami (Yoshizawa-Randlett system).** The relevant lesson is that a domain with genuinely 3D, sequence-dependent manipulations converged on a **standardised symbolic vocabulary** (dashed vs dot-dash lines for valley vs mountain, distinct arrow forms for fold, fold-and-unfold, turn over, repeat). Rigging has no such vocabulary. Inventing one, publishing it, and rendering it consistently is a credible contribution in itself, arguably more durable than the site.
- **Guitar chord diagram libraries. This is the single best proof that the configurator model works, and it is verified, not speculative.** A small fixed grammar (strings, frets, finger numbers, barre, open, muted) renders **any** chord from a compact data representation, including voicings nobody ever drew by hand. GitHub returns 131 repos for `chord diagram guitar`, with several mature MIT-licensed implementations: `acspike/ChordJS` (148 stars, MIT, "Draw guitar chord diagrams on HTML5 canvas"), `BeauNouvelle/SwiftyGuitarChords` (144 stars, MIT), and most relevantly `andygock/chordy-svg` (28 stars, MIT, "generating guitar chord diagrams in SVG format"). *Steal this:* read `chordy-svg` before designing the rig schema. The interesting design work is not the renderer, it is how tersely the **data model** captures a chord. The equivalent question for this project is what the minimum sufficient data model for a rig is, and that decision constrains everything downstream.
- **Knitting and crochet stitch libraries.** Standardised chart symbols plus written pattern notation. The lesson is the **dual representation**: a symbolic chart for the visual learner and a terse written form for the experienced user. A rig probably wants both a diagram and a copy-pasteable text notation.
- **iFixit.** The strongest mobile-instruction UX reference: step-by-step with one action per step, explicit tool call-outs, and community-contributed guides under an open licence. Worth studying its **contribution and moderation model** as much as its layout, since an open rig dataset has exactly the same trust problem.
- **Tie-a-tie sites.** Direct analogue to Animated Knots, same limitation, same photo-sequence approach.
- **Surgical suture training.** Notable mainly for how much it relies on video and physical simulators, which is evidence that some manipulations resist static diagramming.

## Rendering Libraries

| Library | License | Maturity (Aug 2026) | Fit |
| --- | --- | --- | --- |
| **GSAP** | **Not OSI open source.** GitHub API reports `license: null`, no `LICENSE.md` in the repo. Ships a custom "Standard License." Made **free for all commercial use** after the Webflow acquisition, but free is not the same as open source. | 27.6k stars, actively maintained | Best-in-class SVG animation, incl. MorphSVG, now free. **But an explicitly open-source project cannot ship a dependency with no OSI licence without qualifying the claim.** Avoid, or accept and document the caveat. |
| **anime.js** | MIT | 72.0k stars, pushed 2026-08-09 | Strongest recommendation for SVG/DOM timeline animation. Genuinely open source, huge, current, small. |
| **Motion** (ex Framer Motion) | MIT | 33.2k stars, pushed 2026-08-10 | Excellent if the site is React. Gesture handling is a real advantage for scrub-through-steps on touch. |
| **three.js** | MIT | 114.4k stars, pushed 2026-08-10 | The route for real 3D rope. `TubeGeometry` extruded along a `CatmullRomCurve3` is the standard technique for rope/cable, and it is the credible path to beating Knots 3D at its own game. Heavier payload, harder on low-end phones. |
| **react-three-fiber** | MIT | 31.7k stars, pushed 2026-08-10 | Declarative three.js. Worth it only if already committed to React. |
| **D3** | ISC | 113.4k stars | Not an animation library, but its data-join and path/curve generators (`d3-shape`, `d3-path`) are excellent for **generating** rig geometry as SVG paths. Pair with anime.js. |
| **Pixi.js** | MIT | 48.0k stars | 2D WebGL. Overkill unless doing many animated elements at once. |
| **SVG.js** | `NOASSERTION` (GitHub cannot auto-detect; believed MIT, **verify before adopting**) | 11.8k stars, pushed 2026-08-04 | Lightweight programmatic SVG construction. Licence needs manual confirmation. |
| **rough-notation** | MIT | 9.7k stars, last pushed **2024-03-18** | Stale. Interesting only for hand-drawn annotation aesthetics. |

**On GSAP specifically:** the prior run's finding holds and should be stated precisely.
GSAP became **free for commercial use** in late 2024 after Webflow acquired GreenSock, including previously paid plugins like MorphSVG and SplitText (<https://webflow.com/blog/gsap-becomes-free>).
It did **not** become open source.
The GitHub API returns `license: null` for `greensock/GSAP` and there is no `LICENSE.md`.
For a project whose identity is "open source," shipping a core rendering dependency under a bespoke proprietary licence is a real inconsistency, not a technicality. anime.js plus D3 path generation covers the same ground under MIT/ISC.

**Recommended stack for this project:** SVG generated from a rig data model using D3 path/curve generators, animated with anime.js, all MIT/ISC.
Reach for three.js `TubeGeometry` only if the 2D diagram proves insufficient for showing line-over-line topology, which is the one place 2D genuinely struggles and where Knots 3D earns its money.

## Community Signal

This is the weakest section in the report, and that weakness is itself a finding.

**Methodological caveat, stated up front:** Reddit is now largely closed to automated access. `reddit.com` is blocked to the search tool's user agent, and both `reddit.com/*.json` and `old.reddit.com/*.json` return HTTP 403. Only `old.reddit.com` HTML search remained accessible, and its relevance ranking is poor. The sample below is therefore thin and should not be treated as representative.

### What was actually found

A real, upvoted r/Fishing thread, "Leaders: Which, How, and When to Use Them?"
<https://old.reddit.com/r/Fishing/comments/36y1yt/leaders_which_how_and_when_to_use_them/>

> "Ever since I was a kid I've been fishing, but I was never taught about the different rigs and how to tie them. The one that still doesn't make that much sense to me is the use of leaders. [...] I don't know what kind of leader I need, what strength leader to get, or even know when and where to use it."

That is a clean articulation of the target user, and notably it is about **when and why**, not **how to tie**.
The how-to is solved. The selection logic is not.

The replies are the more interesting data:

> "Leader lengths are a personal preference as is the use of them. Most prefer 12 to 18 inches"
>
> "That's for wire leaders generally. If you're going for the other end of the spectrum [...] you'll want a lot more than 18 inches. Most fly leaders are 6-12 feet [...] Guys trolling leadcore line use 10-30 feet"
>
> "I guess it all depends. If your going the thinner approach maybe 4 to 6 pound. Depends on what your fishing for."

**Read this carefully, because it cuts both ways.**
It confirms the information is genuinely scattered and inconsistent, which supports the gap.
It also shows that the community's own answer is *"it depends" and *"personal preference."*
That is a serious warning for feature (b): a conditions-to-rig recommendation engine is asserting authority over a question that experienced anglers openly treat as contested and situational.
Build it wrong and the response will not be gratitude, it will be argument.

### Someone literally asks for this product

r/Fishing, "Best place to learn how to set up rods?"
<https://old.reddit.com/r/Fishing/comments/14tf6ol/best_place_to_learn_how_to_set_up_rods/>

> "I'm not really new to fishing, but I have always fished with family who have rods already set up for what we are after. I really want to fish independently but want to know if there is a **website, YouTube channel, app, etc. that gives step by step guides on how to actually set up rods for different styles of fishing.** (Bobbers, jigs, etc.)"

That is the product brief, written by a user, unprompted.
The thread's sole substantive reply is effectively "go google it":

> "You can absolutely find this info online, I would start by looking up which specific rig you want first. Aka 'bass rod setup'."

Nobody could name a resource. That is the gap, observed in the wild.

### The real complaint is navigability, not rendering

r/Fishing, "Best online resources"
<https://old.reddit.com/r/Fishing/comments/1782gqo/best_online_resources/>

The original poster, a complete novice who lost the uncle who would have taught him:

> "He knew all the knots, gear, types of fish and proper baits [...] I, of course, never paid any attention. Now, I have a 9 month old son and really want to be able to provide those same sort of experiences for him, but I know absolutely [nothing]"

The top advice, and the OP's reaction, are the most useful two lines found in this entire research pass:

> "**YouTube has basically all your answers but it'll be tedious to sort through the information.**"
>
> "Thanks! I definitely went down a **YouTube rabbit hole and got overwhelmed rather quickly** 😆"

**This reframes the opportunity.**
The complaint is not "the diagrams are bad."
It is "the information exists, in enormous quantity, and I cannot navigate it."
That is a **structure, indexing, and progressive-disclosure problem**, which is precisely what a queryable structured database solves and what a wall of YouTube videos cannot.
It also happens to be the part of the project that needs no rendering breakthrough at all.

### A caution from the same thread

The same commenter dispensing advice also said:

> "how-to tie [knot type - **clinch, palomar, and uni are basically all you ever need**]"

Take that seriously.
A stated goal of documenting *every* way to rig terminal tackle runs directly into experienced anglers' view that a handful of knots covers almost everything.
Exhaustiveness may be a collector's goal rather than a user's goal.
The long tail is worth having for credibility and for search traffic, but the product should be designed around the short head.

### Honest read on demand

- Adjacent products exist and are commercially alive (Fishbox, Deep Dive, Rigline, Knots 3D), which proves anglers **will** install reference and advice tools.
- Anglers **do** explicitly ask for "a website, YouTube channel, app, etc. that gives step by step guides on how to actually set up rods for different styles of fishing," and the community cannot name one. Verified above.
- The recurring, concrete, repeated complaint is **information overload and poor navigability**, not diagram quality. "Tedious to sort through," "overwhelmed rather quickly," "went down a YouTube rabbit hole."
- **No evidence was found, in any thread, of anyone complaining that existing rig or knot diagrams are visually inadequate.** Nobody said Animated Knots is hard to follow. Nobody asked for better animation.
- Animated Knots' longevity and heavy ad load suggest large, sustained traffic for knot instruction, so demand for *tying* content is well established. `UNCERTAIN`: no traffic figures were obtained.

**Summary: the demand signal is real, and it is stronger than the first pass suggested, but it points at a different feature than the project leads with.**
Users are asking for **a single navigable structured reference**.
They are not asking for **better rendering**.
Those are separable, and only the first is evidenced.

## Honest Assessment

**Is there a real gap, or is this re-treading a solved problem?**

There is a real gap, but it is **narrower than the project framing implies**, and the framing currently bundles one genuinely novel idea with two crowded ones.

### What is genuinely unoccupied, with high confidence

**The configurator, and the programmatic rendering that makes it possible.** This is the real find.
Every incumbent, without exception, is built on hand-made assets: Animated Knots shoots 13 to 19 photographs per knot, netknots is still serving a converted Flash file, MyRigs is 140 static drawings, Wired2Fish is prose plus stock photography.
Knots 3D is the sole exception and it stops at knots, not rigs, and is app-only.
Because assets are hand-made, none of them can render a rig that does not already exist in their library.
A rig-as-data model with a renderer on top is a **categorically different product**, not a better version of an existing one.

**The open dataset.** The commons is empty in a way that is almost startling: one relevant GitHub repo with one star.
A well-modelled, openly licensed rig schema with real data would be the only one, and it is the kind of artifact that outlives the site that produced it.
This is also the piece with the least competitive risk, because nobody is trying.

### What is crowded, and where the pitch is currently weak

**The recommendation engine is the weakest pillar and should not lead the pitch.**
Deep Dive, Fishbox, and Rigline are funded, shipping, and already fusing satellite water clarity, weather, tides, and species models.
They are ahead on data and will stay ahead.
Worse, the community evidence suggests the underlying question ("which rig for these conditions") is treated by experienced anglers as a matter of preference rather than fact, so a confident recommendation engine invites dispute and is hard to validate.
There is a defensible version of this feature, which is an **explainable** mapping ("this rig, for this reason, with this tradeoff") rather than a prediction. That framing is a teaching tool, not a competitor to a forecast product. Frame it that way or drop it.

**"Comprehensive rig reference" is a content problem wearing an engineering costume.**
The reason no comprehensive reference exists is not that nobody thought of it.
It is that writing accurate, neutral, well-reviewed coverage of hundreds of rigs across freshwater, saltwater, surf, ice, and fly is a large sustained editorial effort with no obvious funding model.
Wired2Fish has a staff and got to roughly 30 to 40 rig articles.
The rendering technology does not solve this. An empty configurator is worth nothing.
**The binding constraint on this project is content, not code, and the plan should say so.**

### The uncomfortable questions

1. **Do anglers actually experience current rig diagrams as broken?** No. Across every thread read in this pass, **not one person complained about diagram or animation quality.** What they complained about, repeatedly and specifically, is that the information is unnavigable: "tedious to sort through," "overwhelmed rather quickly." They also explicitly asked for a single step-by-step resource and could not be pointed to one. So the demand is real, but it is demand for **structure and navigation**, and better rendering does not deliver that. This is the riskiest assumption in the project: the evidenced need and the headline feature are not the same thing.
2. **Does animation genuinely help for rigs, as opposed to knots?** A knot is a sequence of manipulations, so animation is obviously useful. A rig is largely a static topology: line, swivel, leader, weight, hook, in order. Agrawala's finding, that **action diagrams** beat structural ones, suggests the win comes from *spatially separating the incoming component and drawing a guideline*, which is a static technique. Animation may be doing less work than assumed. Worth prototyping before committing to it as the headline feature.
3. **Is "not video" a user preference or an engineering preference?** The stated constraint is explicitly anti-video. Meanwhile the community overwhelmingly consumes rigging content on YouTube. The advantages of programmatic rendering are real (linkable, composable, tiny, forkable, accessible, offline), but they are **producer-side and ecosystem-side advantages**. None of them is a thing a user asks for. That is fine, but it should be an explicit strategic bet rather than an assumed user need.

### Bottom line

**There is a real gap. The owner is not re-treading a solved problem. But the gap is not exactly where the project is currently pointing.**

What the evidence supports, in descending order of confidence:

1. **The structured, queryable rig database is the product.** It does not exist (verified), the commons is empty (verified), and users explicitly ask for it and get told to google (verified). This is the highest-confidence finding in the report and it should be the headline.
2. **The open dataset is the most durable artifact.** One relevant repo with one star is a startling vacuum. This outlives the site.
3. **The configurator is the strongest technical differentiator**, and it is genuinely unoccupied, because every incumbent's asset pipeline structurally forbids it. Animated Knots shoots photographs. netknots hand-authors Flash timelines. MyRigs draws static art. None can render a rig nobody made.
4. **Programmatic rendering is the enabler for 3, not a user-facing benefit in its own right.** No angler asked for it. Justify it as what makes the configurator and the open dataset possible, not as a feature.
5. **The recommendation engine should be demoted** to an explainable teaching layer. As a prediction product it loses to funded incumbents fusing satellite and weather data, and the community treats the underlying question as preference rather than fact.

The two things most likely to kill this:

- **Content, not code, is the binding constraint.** Wired2Fish has staff and reached roughly 30 to 40 rig articles. An empty configurator is worth nothing, and no rendering technology solves the editorial problem.
- **The evidenced need is navigability; the planned differentiator is rendering.** Those can both be true and still not be the same product. Ship the structured, searchable, well-organised reference first, using whatever rendering is cheapest, and confirm people use it. Then earn the renderer.

Put bluntly: the research says build it, but it also says the animated diagrams are the third most valuable thing here, not the first.
