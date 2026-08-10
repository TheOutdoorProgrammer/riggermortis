# 2. Enforce the schema with CUE, JSON Schema, and a Go validator

- Status: accepted
- Deciders: Joey Stout
- Date: 2026-08-10

## Context and Problem Statement

[`docs/schema-prototype.md`](../docs/schema-prototype.md) defines eight record kinds, roughly thirty enumerations, and thirty validation rules.
All of it is currently prose.
Prose does not fail a build, so without machine enforcement the specification and the data will diverge, and the project's entire claim is that its data is verifiably correct.

The question is what technology holds us to our own standard, and whether an existing specification framework covers it.

## Decision Drivers

- The data is **hand-authored by humans and reviewed in pull requests**. Authoring ergonomics and diff readability outrank serialization efficiency.
- The rules span very different classes, from "this string matches a regex" to "walk the graph and prove no threaded component is unbounded."
- The dataset is published under a reciprocal licence, so **downstream consumers need to validate it too**, using tooling we do not control.
- Documentation drift is the specific failure mode to design against. A thirty-rule prose list and a separate implementation will disagree within a month.
- Go is the house language for tooling: one static binary, no runtime dependencies.

## The rules are not one kind of problem

Sorting them by what could possibly enforce each is what makes the answer obvious.

| Class | Example rules | Enforceable by |
| --- | --- | --- |
| Shape, type, enum, regex, conditional requirement | 6, 7, 8, 11, 15, 16, 26, 30 | A schema language |
| Referential integrity across files | 5 | Code |
| Cross-record semantics (a knot's `connects.to` versus a pin's type) | 2 | Code |
| Graph algorithms (connectivity, acyclicity, the interval check) | 1, 4 | Code |
| Derived-value consistency (emitted notation versus structured actions) | 9 | Code |
| Expansion ordering (patterns before graph rules) | 24 | Code |
| Process and history (append-only validation events) | 20 | Code plus CI |

Roughly a third is declarative and the rest is not.
Any answer of the form "use tool X" is therefore wrong on its face.

## Considered Options

1. **OpenAPI**
2. **Protocol Buffers**
3. **Hand-written JSON Schema plus a Go validator**
4. **LinkML**
5. **CUE as the source of truth, generating JSON Schema, plus a Go validator**

## Decision Outcome

Chosen: **option 5**.

A three-layer stack, where each layer does only what it is good at:

1. **CUE** is the single source of truth. Kinds, fields, enumerations, regex constraints, and conditional requirements are defined once. CUE is written in Go and embeddable in Go, so the toolchain stays one binary.
2. **JSON Schema** is generated from CUE and published alongside the dataset. It gives editors inline validation and autocomplete while authoring YAML, and it gives downstream consumers a standard way to validate without adopting our tooling.
3. **A Go validator** implements everything declarative schemas cannot express: reference resolution, cross-record semantics, pattern expansion, and the graph rules.

`additionalProperties: false` is set everywhere. Typos in hand-authored data are the most common error by a wide margin, and silently ignoring an unknown key is how a field goes missing without anyone noticing.

### Two things that matter more than the tool choice

**A rule registry with traceability.**
Every rule gets a stable ID.
The registry records its ID, statement, severity, and tier, and it is the source from which the rules table in the specification is generated.
CI asserts that every documented rule has an implementation and every implementation has a documented rule.
Without this, the prose and the code drift and the specification quietly becomes fiction.

**A failing fixture for every rule.**
Each rule ships with a record that must fail it and a record that must pass.
A validator that silently passes everything is worse than no validator, because it manufactures confidence.
This project already documents that exact failure mode elsewhere: a check whose output nobody inspected was measuring how fast it errored.

### Consequences

Good:

- One definition of every enumeration and constraint, so the specification cannot drift from what is enforced.
- Downstream consumers validate with standard JSON Schema tooling in any language.
- Authors get inline editor errors while typing, not after pushing.
- The Go validator, the pattern expander, and the eventual generator are one binary with no runtime dependencies.

Bad:

- CUE has a genuinely steep learning curve. Its unification model is unlike anything else and its error messages can be cryptic.
- It is a smaller ecosystem than JSON Schema alone, so some problems will have no Stack Overflow answer.
- Two artifacts to keep aligned, CUE and generated JSON Schema, though generation makes that mechanical rather than manual.
- Roughly two thirds of the rules still require hand-written Go, which no tool choice avoids.

## Pros and Cons of the Options

### OpenAPI

- Good: the standard for describing HTTP interfaces, and OAS 3.1 is a proper superset of JSON Schema 2020-12.
- Bad: it describes **APIs**, and this project is a dataset. There are no endpoints, methods, or responses to describe.
- Bad: adopting it to reach the JSON Schema inside means taking an entire interface-description layer for the one part that applies.
- **Revisit if** an HTTP API is ever published over the dataset. It would then describe the API and reference these same JSON Schemas, replacing nothing.

### Protocol Buffers

- Good: excellent wire format, strong cross-language codegen, compact binary encoding.
- Bad: optimised for **machine-to-machine serialization**, which is the opposite of hand-authored files reviewed in pull requests. A human cannot reasonably hand-write a message, so a text layer would be needed anyway.
- Bad: the type system is weaker where we need strength. No regex constraints, no conditional requirements such as "required only when `mounting` is `threaded`", no numeric ranges. Reaching parity means bolting on a CEL-based validation layer.
- Bad: proto3 enums require a zero value, which becomes an implicit default. A schema whose stated principle is "no silent defaults" should not begin with one.
- **Revisit if** a compact binary distribution or gRPC services are ever needed.

### Hand-written JSON Schema plus a Go validator

- Good: boring, standard, universally supported, no new language to learn.
- Good: identical runtime outcome to the chosen option.
- Bad: the enumeration registry would live in the specification prose **and** in eight hand-maintained schema files. That is the duplication this project keeps identifying and removing.
- Bad: JSON Schema is verbose to hand-write and unpleasant to review at this size.

### LinkML

- Good: purpose-built for **datasets** rather than APIs, which is exactly our shape, and proven at scale in biomedical data commons.
- Good: generates JSON Schema, SHACL, documentation, and code from one source. Its documentation generation directly attacks our drift problem.
- Good: enumerations are first-class and a permissible value can carry its own definition and source URI, which fits a project that cites everything.
- Bad: Python-based, against a house preference for single static Go binaries with no runtime dependencies.
- Bad: oriented toward linked data and ontologies, bringing modelling concepts this project does not need.
- **Genuinely close.** If documentation generation later proves more valuable than Go-nativity, this is the option to reconsider, and the CUE definitions would port to it without a redesign.

## More Information

- Rules and enumerations: [`docs/schema-prototype.md`](../docs/schema-prototype.md).
- Rejected here but relevant elsewhere: Rego and Open Policy Agent express cross-record relational rules declaratively and would suit rules 2, 5, 8, 10, 13, 16 and 17 well. They are awkward for the graph traversals in rules 1 and 4, and adding a second policy runtime for a subset of the rules is not justified at this size. Worth revisiting if the relational rule set outgrows readable Go.
