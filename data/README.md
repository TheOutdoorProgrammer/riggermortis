# data

<!-- SPDX-License-Identifier: ODbL-1.0 -->
> **Licence.** Everything in this directory, at any depth, is licensed under the
> **Open Database License (ODbL) v1.0**. See [LICENSE-DATA](../LICENSE-DATA).
> The rest of the repository, including the spec and all code, is Apache-2.0.

**This directory is the licence boundary.** A file inside it is part of the database. A file outside it is not.

The records here are the dataset: components, lines, knots, rigs, patterns, techniques, species, and sources.
Their shape is defined by [`docs/spec.md`](../docs/spec.md), which is Apache-2.0 so that anyone may implement the format without inheriting the dataset's terms.

## What ODbL means in practice

- **Build whatever you like on it, commercially, closed source.** A site, an app, a printed chart. Those are Produced Works and they need only attribution.
- **Attribution is required** on any Produced Work. The exact notice is in [NOTICE](../NOTICE).
- **If you modify the database and publicly use something built on your modified version, you must offer that modified database** in machine-readable form under ODbL. Improvements stay open.

You are not required to send changes back here.
No licence can compel that.
The obligation is to publish, and a pull request is simply the easiest way to satisfy it.

## Every record carries its own header

Records are copied out of repositories, bundled, and pasted into issues.
A file that loses its directory loses its terms, so each one repeats them:

```yaml
# SPDX-License-Identifier: ODbL-1.0
kind: component
schema_version: 0
id: swivel.barrel
```

CI enforces the header's presence, exactly as it enforces every other rule.
