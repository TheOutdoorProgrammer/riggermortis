// SPDX-License-Identifier: Apache-2.0
//
// The whole dataset as one value, keyed by id.
//
// Loaded with `-l '"records"' -l 'id'`, which places each data file at a path
// computed from its own contents rather than at the root, where they would
// collide. Binding the label to the id also proves that a record's id matches
// its index key, for free.
// Kept out of the schema package on purpose: a published schema must not
// carry the dataset. With records/ in cue/, `cue def --out jsonschema` emits
// the dataset shape instead of the record definitions.
package validate

import rm "github.com/theoutdoorprogrammer/riggermortis/cue:riggermortis"

records: [ID=rm.#Id]: rm.#Dispatch & {id: ID}

// Every id present in the loaded dataset. Cross-record rules narrow against
// this rather than against #Id, which only checks grammar.
_ids: {for id, _ in records {(id): true}}
