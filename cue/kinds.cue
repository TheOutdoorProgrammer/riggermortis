// SPDX-License-Identifier: Apache-2.0
//
// The record kinds. Definitions are closed, so an unregistered field is an
// error rather than something silently ignored. That is what stops a typo
// quietly dropping a value, and what makes vendor fields structurally
// impossible: there is nowhere to put one.
package riggermortis

#Record: #Component | #Line | #Knot | #Rigging | #Rig | #Source

#Component: {
	#Common
	kind:           "component"
	name:           string
	mounting:       #Mounting
	blocks_passage: bool
	pins: [...#Pin]

	// Validator-only. Never rendered.
	bore_mm?: number

	soft?: bool
	landmarks?: [string]: number & >=0 & <=1
	body_profile?: #BodyProfile
	bend_style?:   #BendStyle
	point_style?:  #PointStyle
	shank?:        #Shank
	variants?:     #Variants
}

#Line: {
	#Common
	kind:     "line"
	name:     string
	material: #LineMaterial
	properties: {
		stretch:               #Low3
		density:               #Density
		abrasion_resistance:   #Low3
		underwater_visibility: #Low3
		knot_sensitivity:      #Low3
	}
	// Breaking load per diameter is not standardised, so each pairing is a
	// sourced measurement rather than a property of the material.
	diameters: [...{
		breaking_load_kg: number & >0
		diameter_mm:      number & >0
		source:           #Id
		rank?:            #Rank
	}]
}

#Action: {
	verb: #Verb
	subject: string | [...string]
	names?:   string
	through?: string
	around?: string | [...string]
	over?:      string
	direction?: #Direction | #ReeveDirection
	rotation?:  #Rotation
	chirality?: #Chirality
	force?:     #Force
	repeat?:    int | #Range | string
	wet?:       bool
	length_mm?: number
}

#Descriptor: {
	desc:     string
	subject?: string
	to?:      string
	against?: string
}

#Stage: {
	id?: int
	actions?: [...#Action]
	descriptors?: [...#Descriptor]
	prose?:    string
	notation?: string
	expand?:   #Expand
}

#Expand: {
	pattern: string
	count:   int & >0 | #Range
	from?:   string
	to?:     string
	with?: [string]: _
}

#PatternDef: {
	target: #PatternTarget
	params?: [...{
		name:     string
		type:     "ref" | "number" | "string" | "range"
		required: bool
	}]
	stages?: [...#Stage]
	nodes?: [...]
	edges?: [...]
	exposes?: {in: string, out: string}
}

#Knot: {
	#Common
	kind:  "knot"
	names: #Names
	role:  #KnotRole
	connects: {
		from: #PinType
		to: [...#PinType]
	}
	line_types: [...#LineMaterial]
	objects?: [...{ref: string, pin_type: #PinType}]
	patterns?: [string]: #PatternDef
	stages: [...#Stage]
	strength?: claims: [...#StrengthClaim]
	cautions?: [...#Caution]
	failure_modes?: [...string]
}

#Rigging: {
	#Common
	kind:     "rigging"
	names:    #Names
	weedless: bool
	applies_to: {
		body_profile: [...#BodyProfile]
		hook_bend: [...#BendStyle]
	}
	stages: [...{
		id?: int
		actions?: [...{
			verb: #RigVerb
			// A landmark name, or a surface letter with a normalised position.
			at?:          string
			subject?:     string
			along?:       string
			until?:       string
			relative_to?: string
			angle_deg?:   number
			roll_deg?:    number
			depth?:       number & >=0 & <=1
		}]
		descriptors?: [...{
			desc:     #RigDescriptor
			subject?: string
			against?: string
		}]
		prose?:    string
		notation?: string
	}]
	failure_modes?: [...string]
}

#Node: {
	id:                string
	ref?:              #Id
	type?:             "line"
	role?:             #NodeRole
	material?:         #LineMaterial
	breaking_load_kg?: number
	length_mm?:        number | #Range
	mass_g?:           number
	gap_mm?:           number
	diameter_mm?:      number
	rating_kg?:        number
	arms?:             int
}

#Edge: {
	from: string
	to:   string
	rel:  #EdgeRel
	// Absent when the source node is itself the knot.
	knot?:    #Id
	rigging?: #Id
	pin?:     string
	travel?: {
		toward_rod?:      string
		toward_terminal?: string
	}
}

#Rig: {
	#Common
	kind:        "rig"
	names:       #Names
	variant_of?: #Id
	legality?: {
		general_warning: bool
		notes?:          string
		restrictions?: [...{
			jurisdiction: string
			rule:         string
			source:       #Id
			tier:         #Tier
		}]
	}
	patterns?: [string]: #PatternDef
	nodes: [...#Node]
	edges: [...#Edge]
	expand?: [...#Expand]
	failure_modes?: [...string]
}

#Source: {
	#Common
	kind:      "source"
	type:      #SourceType
	title:     string
	author?:   string
	year:      int | null
	url?:      string
	accessed?: string
	// SPDX identifier, or the literals below. Not an enum: SPDX is the registry.
	license:     string
	copy_policy: #CopyPolicy
	reliability: #Reliability
	notes?:      string
}
