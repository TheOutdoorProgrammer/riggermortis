// SPDX-License-Identifier: Apache-2.0

package solve

import (
	"fmt"

	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

// CordReading is the over/under sequence along one cord, left to right.
type CordReading struct {
	Cord        string
	Sequence    []string
	Alternating bool
}

// ReadCrossings reports what a rendered stage actually draws. Eyeballing a
// render cannot tell a reef from a granny; the alternation can.
func ReadCrossings(g *spec.Geometry, stage int) []CordReading {
	if g == nil || stage < 0 || stage >= len(g.Stages) {
		return nil
	}
	var out []CordReading
	for _, cord := range g.Cords {
		var hits []string
		for _, s := range g.Stages[stage].Segments {
			if s.Cord != cord || s.Z == 0 || len(s.Points) == 0 {
				continue
			}
			z := "U"
			if s.Z == 2 {
				z = "O"
			}
			hits = append(hits, z)
		}

		r := CordReading{Cord: cord, Alternating: true}
		for i, h := range hits {
			r.Sequence = append(r.Sequence, h)
			if i > 0 && h == hits[i-1] {
				r.Alternating = false
			}
		}
		out = append(out, r)
	}
	return out
}

// CheckAlternating returns an error naming the cord that does not alternate.
func CheckAlternating(g *spec.Geometry, stage int) error {
	for _, r := range ReadCrossings(g, stage) {
		if !r.Alternating {
			return fmt.Errorf("cord %s does not alternate: %v", r.Cord, r.Sequence)
		}
	}
	return nil
}
