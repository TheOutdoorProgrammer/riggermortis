# 5. Dual licence: ODbL for the data, Apache-2.0 for everything else

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

The repository was published deliberately without a `LICENSE`, which meant all rights reserved, because a code-and-data project needs two licences and choosing wrongly is expensive to reverse once a dataset has contributors.

The requirement, stated by the owner:

> Attribution is required. The work may be used as-is, commercially, without the user making their own work public. If they **modify the data**, it has to be contributed back or made public in some fashion.

## Decision Drivers

- The dataset does not exist anywhere else. Research established that no structured fishing rig dataset exists on GitHub, Kaggle, or in any public catalogue. It is the project's only real asset.
- Every record is hand-authored, so the cost of the data vastly exceeds the cost of the code.
- The **format** should be freely implementable by anyone. Adoption of the spec is a goal.
- Well-funded incumbents exist in adjacent products and could ingest an unprotected dataset without returning anything.

## Considered Options

1. **MIT for everything.**
2. **CC BY 4.0 for the data.**
3. **CC BY-SA 4.0 for the data.**
4. **ODbL 1.0 for the data, Apache-2.0 for everything else.**
5. **CDLA-Sharing-1.0 for the data.**

## Decision Outcome

Chosen: **option 4.**

- **`data/` is ODbL-1.0.** The directory is the licence boundary; a file inside it is part of the database and a file outside it is not.
- **Everything else is Apache-2.0**, including the spec, the documentation, the ADRs, and all code.

ODbL is the only candidate that expresses the stated requirement without ambiguity, because it distinguishes two things the Creative Commons licences do not:

- A **Produced Work**, meaning something built *from* the database such as an app, a site, or a rendered diagram, requires **attribution only**. It may be commercial and closed.
- A **Derivative Database**, meaning the data itself modified, must be **offered under ODbL** if a Produced Work from it is publicly used.

That maps one to one onto "use as-is freely, but modifications become public."
It also closes the hosted-service gap, since ODbL triggers on public *use* rather than only on distribution.

### The spec is Apache-2.0 on purpose

Putting share-alike on the *format* would discourage independent implementations, which is the opposite of the intent.
Anyone should be able to read `docs/spec.md`, write a renderer or a validator, and owe nothing.
Reciprocity belongs on the data, which is the expensive part, not on the description of its shape.

### What no licence can do

ODbL cannot compel contributions back upstream, and neither can anything else.
The mechanism is publish-alike: a modifier must make their version available, and we may then choose to merge it.
Anyone hoping for mandatory upstreaming should know that is not a property any OSI or Open Data Commons licence has.

### Consequences

Good:

- The stated requirement is satisfied precisely, including the hosted-service case.
- The format stays freely implementable, so third-party tooling is encouraged.
- The dataset cannot be silently absorbed into a closed product and improved in private.
- ODbL has a decade of real-world precedent through OpenStreetMap, so the edge cases have been litigated by someone else.

Bad:

- ODbL is unfamiliar to most developers and will cause hesitation.
- The Produced Work boundary still needs judgement in odd cases, and a generated SVG shipped as part of a page is a Produced Work while a bulk export of records is a Derivative Database.
- Two licences mean a boundary to police. Mitigated by making the boundary a directory rather than a file-by-file judgement.
- Share-alike data cannot be merged into permissively licensed datasets, which is friction for anyone trying to combine sources.

## Pros and Cons of the Options

### MIT for everything

- Good: universally understood, zero friction.
- Bad: MIT is a **software** licence. Applied to a database its effect is ambiguous, and it does not address database rights at all.
- Bad: no reciprocity whatsoever, which fails the requirement outright.

### CC BY 4.0 for the data

- Good: simple, maximally adoptable, covers *sui generis* database rights since 4.0.
- Bad: attribution only. A competitor may take the dataset, extend it substantially, and publish nothing. Fails the requirement.

### CC BY-SA 4.0 for the data

- Good: familiar, covers database rights since 4.0, and does impose share-alike.
- Bad: **leaves it genuinely unclear whether a product built from the data is itself infected.** That ambiguity is precisely why OpenStreetMap migrated off CC BY-SA to ODbL in 2012, and inheriting a known-broken arrangement with the case study already written would be careless.
- Bad: weaker on hosted services than ODbL's public-use trigger.

### CDLA-Sharing-1.0

- Good: modern, purpose-built for data, from the Linux Foundation, and shaped very similarly to ODbL.
- Bad: little real-world precedent, so its edge cases are untested where ODbL's have been argued through a large project.
- **Revisit if** ODbL's Produced Work boundary proves to be genuine friction rather than a useful distinction.

## More Information

- [LICENSE](../LICENSE), [LICENSE-DATA](../LICENSE-DATA), [NOTICE](../NOTICE), and [`data/README.md`](../data/README.md) for the boundary in practice.
- Every record carries `# SPDX-License-Identifier: ODbL-1.0`, and every spec document carries the Apache-2.0 equivalent, so a file separated from its repository keeps its terms. CI enforces this.
