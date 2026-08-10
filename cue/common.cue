// SPDX-License-Identifier: Apache-2.0
//
// Fields every record carries, plus the shapes shared across kinds.
package riggermortis

// Exactly one dot. Category from #IdCategory, name in lowercase kebab-case.
#Id: string & =~#"^[a-z][a-z0-9]*\.[a-z0-9]+(-[a-z0-9]+)*$"#

#Range: [number, number]

#ValidationEvent: {
	method:  #ValidationMethod
	result:  #ValidationResult
	by:      string
	date:    string & =~#"^\d{4}-\d{2}-\d{2}$"#
	detail?: string
	sources?: [...#Id]
}

#Validation: {
	status: #ValidationStatus
	events: [...#ValidationEvent]
}

#ExternalId: {
	catalog: #Catalog
	// Always a string. Catalogue identifiers carry prefixes, leading zeros
	// and whole URIs, so a number would silently corrupt them.
	id:      string
	source?: #Id
	rank?:   #Rank
}

#Common: {
	schema_version: int & >=0
	id:             #Id
	validation:     #Validation
	former_ids?: [...#Id]
	external_ids?: [...#ExternalId]
}

#Names: {
	canonical: string
	aliases?: [...string]
	trademark_note?: string
}

// A claim that cannot be machine-proven. Sample size is mandatory because
// most published knot tests are n=1 and a reader deserves to know.
#StrengthClaim: {
	residual_pct: number | null
	line?:        #LineMaterial
	n:            int & >=0
	note?:        string
	source:       #Id
	tier:         #Tier
	rank:         #Rank
}

#Caution: {
	severity: #CautionSeverity
	note:     string
	sources: [...#Id]
}

#Variants: {
	axis: #VariantAxis
	values: [...number]
}

#Pin: {
	id:                string
	type:              #PinType
	wire_diameter_mm?: number
}
