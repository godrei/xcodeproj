package pbxproj

import (
	"regexp"
	"strings"
)

const (
	pbxTargetDependencyISA                 = "PBXTargetDependency"
	pbxTargetDependencySectionBeginPattern = `/* Begin PBXTargetDependency section */`
	pbxTargetDependencySectionEndPattern   = `/* End PBXTargetDependency section */`
)

// PBXTargetDependency ...
type PBXTargetDependency struct {
	id            string
	isa           string
	targetID      string
	targetProxyID string
}

// ParsePBXTargetDependencySection ...
func ParsePBXTargetDependencySection(pbxTargetDependencySectionLines []string) []PBXTargetDependency {
	pbxTargetDependencies := []PBXTargetDependency{}
	pbxTargetDependency := PBXTargetDependency{}
	isPBXTargetDependency := false

	// BAC384251BA9F569005CFE20 /* PBXTargetDependency */ = {
	pbxTargetDependencyBeginRegexp := regexp.MustCompile(`(?P<id>[A-Z0-9]+) /\* (?P<isa>.*) \*/ = {`)
	pbxTargetDependencyEndPattern := `};`

	// isa = PBXNativeTarget;
	isaRegexp := regexp.MustCompile(`isa = (?P<isa>.*);`)

	// target = BAAFFED019EE788800F3AC91 /* SampleAppWithCocoapods */;
	targetRegexp := regexp.MustCompile(`target = (?P<id>[A-Z0-9]+) /\* (?P<name>.*) \*/;`)

	// targetProxy = BAC3842F1BA9F569005CFE20 /* PBXContainerItemProxy */;
	targetProxyRegexp := regexp.MustCompile(`targetProxy = (?P<id>[A-Z0-9]+) /\* (?P<isa>.*) \*/;`)

	for _, line := range pbxTargetDependencySectionLines {
		if line == pbxTargetDependencyEndPattern {
			pbxTargetDependencies = append(pbxTargetDependencies, pbxTargetDependency)

			pbxTargetDependency = PBXTargetDependency{}

			isPBXTargetDependency = false
			continue
		}

		if matches := pbxTargetDependencyBeginRegexp.FindStringSubmatch(line); len(matches) == 3 {
			pbxTargetDependency.id = strings.Trim(matches[1], `"`)

			isPBXTargetDependency = true
			continue
		}

		if !isPBXTargetDependency {
			continue
		}

		// PBXTargetDependency
		if matches := isaRegexp.FindStringSubmatch(line); len(matches) == 2 {
			pbxTargetDependency.isa = strings.Trim(matches[1], `"`)
		}

		if matches := targetRegexp.FindStringSubmatch(line); len(matches) == 3 {
			pbxTargetDependency.targetID = strings.Trim(matches[1], `"`)
		}

		if matches := targetProxyRegexp.FindStringSubmatch(line); len(matches) == 3 {
			pbxTargetDependency.targetProxyID = strings.Trim(matches[1], `"`)
		}
		// -----
	}

	return pbxTargetDependencies
}
