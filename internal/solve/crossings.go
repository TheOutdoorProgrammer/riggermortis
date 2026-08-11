// SPDX-License-Identifier: Apache-2.0

package solve

import (
	"github.com/theoutdoorprogrammer/riggermortis/internal/rope"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

// CordReading is the over/under sequence along one cord, left to right.
type CordReading struct {
	Cord        string
	Sequence    []string
	Alternating bool
}

// ReadCrossings reports what a stage actually draws. Eyeballing a render cannot
// tell a reef from a granny; this can.
func ReadCrossings(r *spec.Record, stage int) []CordReading {
	all := Cords(r)
	if stage < 0 || stage >= len(all) {
		return nil
	}
	var out []CordReading
	for ci, seq := range rope.Reading(all[stage]) {
		cr := CordReading{Cord: all[stage][ci].ID, Alternating: true}
		for i, over := range seq {
			z := "U"
			if over {
				z = "O"
			}
			cr.Sequence = append(cr.Sequence, z)
			if i > 0 && over == seq[i-1] {
				cr.Alternating = false
			}
		}
		out = append(out, cr)
	}
	return out
}
