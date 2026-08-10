# 8. Go templates and a few kilobytes of vanilla JavaScript for the site

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

[ADR 0007](0007-go-only-pregenerated-svg.md) settled that SVGs are pre-generated, nothing runs in production, and the deployed artifact is a directory of static files.
It deliberately left the site itself undecided.

The question is what builds the HTML and what provides the interactivity, with React as the obvious default worth testing against.

## Decision Drivers

- **Findability is the product.** Community research found the unmet need is not better diagrams, it is that existing information cannot be sorted or searched. People will arrive from a search engine on a single knot.
- The interactive state on this site is a search string, a set of active filters, and a current animation stage. That is the whole of it. No authentication, no forms, no persistence, no real-time anything.
- [ADR 0007](0007-go-only-pregenerated-svg.md) just removed the second toolchain. Reintroducing one needs to earn its place.

## Decision Outcome

**Go `html/template` at build time, plus a few kilobytes of hand-written JavaScript.**

The same binary that validates records and generates SVG also emits the HTML. One toolchain, standard library only, no Node in the build.

Interactivity is two small things:

- A **stage player** for knot and rigging animations, toggling visibility across pre-generated stages. Roughly forty lines.
- A **filter** over a shipped JSON index. At the catalogue's size that is an array filter.

### Why not React

Not dogma, and not a claim that React is bad. It is that this site has three pieces of state and React exists to manage many.

The stronger objection is architectural. Findability is the entire product thesis, so pages must be static HTML with the SVG inline, indexable without executing anything. Shipping a JavaScript shell that hydrates works against the one thing the site is for.

### Progressive enhancement falls out for free

With stages pre-generated and toggled by CSS, **the knot steps work with JavaScript disabled.** For a teaching reference that is a real property, not a box tick: the instructions are never behind a runtime.

### Consequences

Good:

- One toolchain, one binary, no `node_modules`, nothing to keep current.
- Pages are static HTML that a crawler reads without executing anything.
- Content works with scripting disabled.
- The frontend stays genuinely swappable, because the JSON and SVG bundle is a clean boundary rather than an internal call.

Bad:

- No component model. Layout reuse is template composition, which is workable and less pleasant than components.
- No hot reload out of the box. A file watcher covers it, but that is something to build.
- Hand-written JavaScript has no framework guard rails, so it stays small by discipline rather than by structure.
- Anything genuinely interactive later, a comparison view or a conditions query with real interaction, will strain this and is the trigger to revisit.

## Considered Alternatives

**Astro.**
The framework built precisely for mostly-static content with islands of interactivity, shipping zero JavaScript by default and hydrating only what is marked. Better authoring ergonomics, content collections, real layouts. Rejected **for now** because it reintroduces Node to a build that has just been reduced to one Go binary, for state that does not yet exist. **Held in reserve**, and the first choice if templates start to hurt. Its islands can be React where a component genuinely warrants it, so this is not a decision against React so much as a decision to defer needing it.

**Alpine.js on static HTML.**
Around fifteen kilobytes, no build step, sits between hand-written script and a framework. A reasonable middle and **held in reserve** alongside Astro. Rejected now because forty lines of script does not need a library.

**htmx.**
Rejected outright rather than deferred. It is built around server round-trips and there is no server.

**React with a static export.**
Heaviest option for the least benefit here. It would ship a runtime and a hydration step to manage three pieces of state, on a site whose thesis is that pages must be readable without executing anything.

## More Information

- The frontend is swappable by construction. Everything it consumes is the published JSON and SVG bundle, and it contains no geometry and no validation.
- Revisit when a view needs genuine interaction beyond filtering and stage playback. Astro first, then Alpine, and React only inside an island that has earned it.
