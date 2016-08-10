package pbxproj

import (
	"regexp"
	"strings"
)

const (
	pbxNativeTargetISA                 = "PBXNativeTarget"
	pbxNativeTargetSectionBeginPattern = `/* Begin PBXNativeTarget section */`
	pbxNativeTargetSectionEndPattern   = `/* End PBXNativeTarget section */`
)

// PBXNativeTarget ...
type PBXNativeTarget struct {
	id          string
	isa         string
	name        string
	productName string
	productType string

	productReferenceID       string
	buildConfigurationListID string

	buildPhasesIDs  []string
	buildRulesIDs   []string
	dependenciesIDs []string
}

// ParsePBXNativeTargetSection ...
func ParsePBXNativeTargetSection(pbxNativeTargetSectionLines []string) []PBXNativeTarget {
	pbxNativeTargets := []PBXNativeTarget{}

	pbxNativeTarget := PBXNativeTarget{}

	isPBXNativeTarget := false

	// // BAAFFED019EE788800F3AC91 /* SampleAppWithCocoapods */ = {
	pbxNativeTargetBeginRegexp := regexp.MustCompile(`(?P<id>[A-Z0-9]+) /\* (?P<name>.*) \*/ = {`)
	pbxNativeTargetEndPattern := `};`

	// isa = PBXNativeTarget;
	isaRegexp := regexp.MustCompile(`isa = (?P<isa>.*);`)

	// buildConfigurationList = BAC3843A1BA9F569005CFE20 /* Build configuration list for PBXNativeTarget "BitriseXcode7SampleTests" */;
	buildConfigurationListRegexp := regexp.MustCompile(`buildConfigurationList = (?P<id>[A-Z0-9]+) /\* .* \*/;`)

	// name = SampleAppWithCocoapods;
	nameRegexp := regexp.MustCompile(`name = (?P<name>.*);`)

	// productName = BitriseXcode7SampleTests;
	productNameRegexp := regexp.MustCompile(`productName = (?P<name>.*);`)

	// productReference = BAAFFEED19EE788800F3AC91 /* SampleAppWithCocoapodsTests.xctest */;
	productReferenceRegexp := regexp.MustCompile(`productReference = (?P<id>[A-Z0-9]+) /\* (?P<path>.*) \*/;`)

	// productType = "com.apple.product-type.bundle.unit-test";
	productTypeRegexp := regexp.MustCompile(`productType = (?P<productType>.*);`)

	buildPhasesBeginPattern := `buildPhases = (`
	buildPhasesEndPattern := `);`
	buildPhaseRegexp := regexp.MustCompile(`\s*(?P<id>[A-Z0-9]+) /\* (?P<ype>.*) \*/,`)
	isBuildPhases := false

	buildRulesBeginPattern := `buildRules = (`
	buildRulesEndPattern := `);`
	buildRuleRegexp := regexp.MustCompile(`\s*(?P<id>[A-Z0-9]+) /\* (?P<ype>.*) \*/,`)
	isBuildRules := false

	dependenciesBeginPattern := `dependencies = (`
	dependenciesEndPattern := `);`
	dependencieRegexp := regexp.MustCompile(`\s*(?P<id>[A-Z0-9]+) /\* (?P<isa>.*) \*/,`)
	isDependencies := false

	for _, line := range pbxNativeTargetSectionLines {
		if line == pbxNativeTargetEndPattern {
			pbxNativeTargets = append(pbxNativeTargets, pbxNativeTarget)

			pbxNativeTarget = PBXNativeTarget{}

			isPBXNativeTarget = false
			continue
		}

		if matches := pbxNativeTargetBeginRegexp.FindStringSubmatch(line); len(matches) == 3 {
			pbxNativeTarget.id = strings.Trim(matches[1], `"`)
			pbxNativeTarget.name = strings.Trim(matches[2], `"`)

			isPBXNativeTarget = true
			continue
		}

		if !isPBXNativeTarget {
			continue
		}

		// PBXNativeTarget

		if matches := isaRegexp.FindStringSubmatch(line); len(matches) == 2 {
			pbxNativeTarget.isa = strings.Trim(matches[1], `"`)
		}

		if matches := buildConfigurationListRegexp.FindStringSubmatch(line); len(matches) == 2 {
			pbxNativeTarget.buildConfigurationListID = strings.Trim(matches[1], `"`)
		}

		if matches := nameRegexp.FindStringSubmatch(line); len(matches) == 2 {
			pbxNativeTarget.name = strings.Trim(matches[1], `"`)
		}

		if matches := productNameRegexp.FindStringSubmatch(line); len(matches) == 2 {
			pbxNativeTarget.productName = strings.Trim(matches[1], `"`)
		}

		if matches := productReferenceRegexp.FindStringSubmatch(line); len(matches) == 3 {
			pbxNativeTarget.productReferenceID = strings.Trim(matches[1], `"`)
		}

		if matches := productTypeRegexp.FindStringSubmatch(line); len(matches) == 2 {
			pbxNativeTarget.productType = strings.Trim(matches[1], `"`)
		}

		// buildPhases
		if isBuildPhases && line == buildPhasesEndPattern {
			isBuildPhases = false
			continue
		}

		if line == buildPhasesBeginPattern {
			isBuildPhases = true
			continue
		}

		if isBuildPhases {
			if matches := buildPhaseRegexp.FindStringSubmatch(line); len(matches) == 3 {
				buildPhaseID := strings.Trim(matches[1], `"`)

				pbxNativeTarget.buildPhasesIDs = append(pbxNativeTarget.buildPhasesIDs, buildPhaseID)
			}
		}
		// -----

		// buildRules
		if isBuildRules && line == buildRulesEndPattern {
			isBuildRules = false
			continue
		}

		if line == buildRulesBeginPattern {
			isBuildRules = true
			continue
		}

		if isBuildRules {
			if matches := buildRuleRegexp.FindStringSubmatch(line); len(matches) == 3 {
				buildRuleID := strings.Trim(matches[1], `"`)

				pbxNativeTarget.buildRulesIDs = append(pbxNativeTarget.buildRulesIDs, buildRuleID)
			}
		}
		// -----

		// dependencies
		if isDependencies && line == dependenciesEndPattern {
			isDependencies = false
			continue
		}

		if line == dependenciesBeginPattern {
			isDependencies = true
			continue
		}

		if isDependencies {
			if matches := dependencieRegexp.FindStringSubmatch(line); len(matches) == 3 {
				dependencieID := strings.Trim(matches[1], `"`)

				pbxNativeTarget.dependenciesIDs = append(pbxNativeTarget.dependenciesIDs, dependencieID)
			}
		}
		// -----

	}

	return pbxNativeTargets
}
