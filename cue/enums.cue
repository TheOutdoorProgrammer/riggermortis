// SPDX-License-Identifier: Apache-2.0
//
// Every closed vocabulary in the spec. This file is the source of truth;
// docs/spec.md's Enumerations section is generated from it.
package riggermortis

#Kind: "component" | "line" | "knot" | "rigging" | "rig" | "pattern" |
	"technique" | "species" | "source"

#IdCategory: #ComponentCategory | "line" | "knot" | "rigging" | "rig" |
	"pattern" | "technique" | "species" | "src"

#ComponentCategory: "hook" | "weight" | "swivel" | "snap" | "split-ring" |
	"bead" | "float" | "stop" | "blade" | "skirt" | "jighead" | "bait" |
	"lure" | "hardware" | "sleeve"

#PinType: "closed-eye" | "open-eye" | "split-ring" | "shank" | "snap" |
	"arm-socket" | "line-end" | "in" | "out" | "loop" | "tag"

#Severity: "error" | "warning" | "info"

#Catalog: "abok" | "wikidata" | "gbif" | "itis" | "worms" | "fao-isscfg" | "agrovoc"

// Catalogues that can be fetched and machine-checked. The rest need a human.
#ResolvableCatalog: "wikidata" | "gbif" | "itis" | "worms" | "fao-isscfg" | "agrovoc"

#Tier: "A" | "B" | "C"
#Rank: "preferred" | "normal" | "deprecated"

#ValidationStatus: "unvalidated" | "machine-only" | "corroborated" |
			"field-tested" | "disputed"
#ValidationMethod: "ci" | "source-corroborated" | "manufacturer-confirmed" |
			"photo-compared" | "tied" | "fished" | "expert-review"
#ValidationResult: "pass" | "fail" | "partial" | "disputed"

#CautionSeverity: "advisory" | "strong" | "do-not-use"

#Mounting:    "tied" | "threaded" | "rigged"
#Unit:        "mm" | "m" | "g" | "kg" | "c" | "s" | "deg"
#VariantAxis: "mass_g" | "gap_mm" | "diameter_mm" | "rating_kg" |
	"breaking_load_kg" | "length_mm" | "supports_g"

#PointStyle: "straight" | "turned-in" | "turned-out" | "knife-edge" | "needle"
#Shank:      "short" | "standard" | "long" | "extra-long"
#BendStyle:  "round" | "octopus" | "circle" | "worm" | "ewg" | "offset-worm" |
		"straight-shank" | "kahle" | "aberdeen" | "siwash" | "treble"
#BodyProfile: "straight" | "curly-tail" | "paddle-tail" | "creature" | "craw" |
			"fluke" | "tube" | "grub" | "worm-live" | "baitfish" | "cut"
#LandmarkSurface: "n" | "d" | "v" | "l" | "t"

#LineMaterial: "mono" | "fluoro" | "braid" | "wire" | "backing" |
		"leader-mono" | "leader-fluoro"
#Low3:    "low" | "medium" | "high"
#Density: "sinking" | "neutral" | "floating"

#KnotRole:       "terminal" | "bend" | "loop" | "stopper" | "arbor"
#Verb:           "GP" | "MB" | "ML" | "MT" | "MV" | "RV" | "TW"
#Direction:      "F" | "A" | "L" | "R" | "U" | "D"
#ReeveDirection: "F-A" | "A-F" | "L-R" | "R-L" | "U-D" | "D-U"
#Rotation:       "CW" | "CCW"
#Chirality:      "/" | "\\"
#Force:          "push" | "pull"
#Plane:          "HP" | "VP" | "EP"

#RigVerb:       "IN" | "OUT" | "TH" | "SL" | "RO" | "SK" | "BU"
#RigDescriptor: "AL" | "CE"

#EdgeRel: "tied" | "threaded" | "fixed" | "clipped" | "crimped" |
		"continuous" | "looped" | "rigged"
#NodeRole: "main-line" | "leader" | "dropper" | "tag"

#MotionPrimary: "drag" | "hop" | "shake" | "deadstick" | "yoyo" | "swim" |
		"rip" | "walk" | "troll" | "vertical" | "drift" | "twitch"
#Contact:     "bottom" | "suspended" | "surface"
#WaterColumn: "bottom" | "lower" | "middle" | "upper" | "surface"
#Clarity:     "muddy" | "stained" | "clear" | "gin-clear"
#Cover:       "none" | "sand" | "gravel" | "rock" | "rubble" | "sparse-grass" |
	"heavy-grass" | "wood" | "laydown" | "dock" | "bridge" | "reef" |
		"oyster" | "mangrove" | "current-seam"
#Season: "prespawn" | "spawn" | "postspawn" | "summer" | "fall" | "winter"

#Water:      "freshwater" | "saltwater" | "both"
#SourceType: "book" | "article" | "agency" | "extension" | "manufacturer" |
		"video" | "expert" | "forum"
#CopyPolicy:  "cite-only" | "quotable" | "adaptable"
#Reliability: "low" | "medium" | "high"

#PatternTarget: "rig" | "knot"
