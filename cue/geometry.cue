// SPDX-License-Identifier: Apache-2.0
//
// Authored geometry for a knot's stages.
//
// This is a staging post, not the destination. The intent is that a solver
// derives positions from the Suber verbs in `stages`; until one exists the
// coordinates are hand-authored so the renderer and the site can be built and
// proven end to end. A record carrying `geometry` is asserting how it looks,
// which is a Tier C claim like any other and needs a human to confirm it.
package riggermortis

#Point: [number, number]

// A run of cord between crossings, not a whole cord.
//
// Over and under cannot be a property of a cord: a square knot alternates,
// with one cord over at the first crossing and under at the second. Paint
// order therefore belongs to the segment. Higher z paints later, and each
// segment is stroked with a background casing first, so whatever paints last
// visibly passes over what came before.
#Segment: {
	cord!: string
	z:     *0 | int
	points!: [#Point, #Point, ...#Point]
}

#StageGeometry: {
	stage!: int
	segments!: [...#Segment]
}

#Geometry: {
	width!:  number
	height!: number
	// Draw order for colour assignment, and the legend's order.
	cords!: [string, ...string]
	stages!: [...#StageGeometry]
}
