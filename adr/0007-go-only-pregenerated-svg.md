# 7. Go only, with SVGs pre-generated at build time

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10
- Supersedes: [ADR 0004](0004-polyglot-seam-go-for-spec-typescript-for-geometry.md)

## Context and Problem Statement

[ADR 0004](0004-polyglot-seam-go-for-spec-typescript-for-geometry.md) split the project across two languages: Go for spec tooling, TypeScript for geometry and rendering, joined by a seam of validated JSON.

Two constraints have since been stated that were not known when it was written, and either one alone would have changed the outcome.

1. **SVGs are pre-generated at build time and served as static assets.** The user interface does not generate anything on the fly.
2. **There is no rig configurator.** Anyone wanting a rig that is not in the catalogue contributes a record for it.

ADR 0004 rested almost entirely on "the browser is the target anyway." Under these constraints it is not, so that decision has to be revisited rather than inherited.

## What changed

**The renderer is a build tool, not a runtime.** It runs in CI, reads validated records, and writes files. It never ships to a browser, so nothing about the browser's ecosystem is relevant to choosing its language.

**The three.js argument was wrong on its own terms.** ADR 0004 cited `TubeGeometry` as making "rope along a path" a library call. That produces a **WebGL mesh**, and the output here is **2D SVG paths**. Generating a mesh and flattening it is strictly worse than computing the 2D outline directly, which is what SVG rope actually needs: take the centreline curve, offset it perpendicular by the radius on both sides, emit one filled path. That is a few hundred lines of curve mathematics owned by this project in any language.

**The iteration-speed argument dissolves.** Browser hot reload was cited as the better loop for adjusting how a knot looks. A Go file watcher that regenerates on save, with the browser reloading the `.svg`, is the same loop.

**Without a configurator, no geometry ever needs to run client-side.** That was the only requirement that would have forced the browser's hand.

## Decision Outcome

**Go for everything: spec definitions in CUE, and one Go binary that validates, expands patterns, runs the conformance corpus, and generates SVG.**

There is no seam, because there are no longer two sides.

### The geometry package stays pure

Curves in, paths out. No file access, no globals, no configuration.

This is not stylistic. A pure geometry package is a **compile target away from running anywhere**, so if a client-side need ever appears, WebAssembly is a build flag rather than a rewrite. Go to WASM costs a few megabytes, which is real but survivable, and TinyGo trims it. That escape hatch costs nothing to preserve and everything to retrofit.

### Nothing runs in production

**Go is a build-time compiler, not a server.** The binary runs in CI, reads records, and writes JSON, SVG and HTML. The deployed artifact is a directory of static files. There is no process, no runtime, no database and no operational surface.

That is worth stating plainly because "Go for everything" could otherwise be read as a backend. It is not one, and the project does not have one.

Deployment is therefore any static host, and Cloudflare Pages is the obvious fit given it is already the pattern in use elsewhere.

### The site is a separate, later, low-stakes decision

It consumes JSON and pre-rendered SVG and contains **no geometry whatsoever**. Go templates executed at build time, Astro, Hugo, anything that emits static files. Deferring it costs nothing because nothing else depends on it.

One consequence to note early: with no server, search and any conditions-to-rig lookup run **client-side over a shipped JSON index**. At the catalogue's expected size that index is small enough to ship whole and filter in a few kilobytes of script, so this is a constraint rather than a problem. It does mean the query layer is bounded by what fits in a browser, and if that ever stops being true it is a decision to revisit, not a thing to quietly work around.

### Consequences

Good:

- One language, one toolchain, one binary, one CI path.
- No interface to version between halves, and no possibility of two implementations disagreeing.
- Matches the house preference without having to argue for it, because the technical case now points the same way.
- Pre-generated output means the site can be entirely static, which is cheap, fast and hard to break.

Bad:

- Go has no 2D vector geometry library comparable to paper.js for path offsetting and boolean operations. That code gets written here. It is bounded and well understood, but it is real work that TypeScript would have partly avoided.
- Should a client-side need ever appear despite this decision, WASM is the answer, and its binary size is a genuine cost.
- Animated SVG must be expressed declaratively inside the file, through CSS or SMIL, or through a small playback script that toggles pre-generated states. Neither computes geometry, but it does constrain how the animation is authored.

## Why no configurator

Recorded here because it is the load-bearing constraint, and a future reader will otherwise propose it as an obvious feature.

- **It is a machine for producing unvalidated rigs.** A user assembles a sliding weight with nothing to stop it, the site renders it faithfully, and the weight leaves on the first cast. The catalogue path runs rule 1 before anything is drawn. A surface that renders whatever it is handed contradicts the premise of the project.
- **Nothing accrues.** A configurator session ends and the dataset is unchanged. A contributed record is validated, reviewed, eventually field-tested, and joins the dataset under ODbL where everyone benefits.
- **Nobody asked for it.** Community research found the unmet need is findability, not composition. No thread examined contained a request for a rig builder.

A contribution flow with good authoring ergonomics, a form that emits a record and opens a pull request, would deliver the same user value while producing something durable. It also runs server-side, so it does not reopen this decision. Worth considering later; not now.

## More Information

- Superseded: [ADR 0004](0004-polyglot-seam-go-for-spec-typescript-for-geometry.md), retained per the project's rule that reversed decisions are recorded rather than rewritten.
- Enforcement stack and the case for CUE: [ADR 0002](0002-schema-enforcement-stack.md).
- Testing the spec itself: [ADR 0003](0003-test-the-spec-as-an-artifact.md).
