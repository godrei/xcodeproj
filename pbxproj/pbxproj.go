package pbxproj

import (
	"bufio"
	"strings"
)

// Model ...
type Model struct {
	archiveVersion string
	classes        []string
	objectVersion  string
	rootObject     string
	objects        ObjectsModel
}

// ObjectsModel ...
type ObjectsModel struct {
	PBXBuildFiles            []PBXBuildFileModel
	PBXContainerItemProxys   []PBXContainerItemProxy
	PBXFileReferences        []PBXFileReference
	PBXFrameworksBuildPhases []PBXFrameworksBuildPhase
	PBXGroups                []PBXGroup
	PBXNativeTargets         []PBXNativeTarget
	PBXProject               PBXProject
	PBXResourcesBuildPhases  []PBXResourcesBuildPhase
	PBXSourcesBuildPhases    []PBXSourcesBuildPhase
	PBXTargetDependencies    []PBXTargetDependency
	PBXVariantGroups         []PBXVariantGroup
	XCBuildConfigurations    []XCBuildConfiguration
	XCConfigurationLists     []XCConfigurationList
	XCVersionGroups          []XCVersionGroup
}

// SplitObjectsSections ...
func SplitObjectsSections(pbxprojContent string) (map[string][]string, error) {
	isaObjectSectionMap := map[string][]string{}

	isPBXNativeTargetSection := false
	pbxNativeTargetSectionLines := []string{}

	isPBXTargetDependencySection := false
	pbxTargetDependencySectionLines := []string{}

	scanner := bufio.NewScanner(strings.NewReader(pbxprojContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// PBXNativeTarget
		if line == pbxNativeTargetSectionEndPattern {
			isaObjectSectionMap[pbxNativeTargetISA] = pbxNativeTargetSectionLines

			isPBXNativeTargetSection = false
			continue
		}

		if line == pbxNativeTargetSectionBeginPattern {
			isPBXNativeTargetSection = true
			continue
		}

		if isPBXNativeTargetSection {
			pbxNativeTargetSectionLines = append(pbxNativeTargetSectionLines, line)
		}
		// -----

		// PBXTargetDependency
		if line == pbxTargetDependencySectionEndPattern {
			isaObjectSectionMap[pbxTargetDependencyISA] = pbxTargetDependencySectionLines

			isPBXTargetDependencySection = false
			continue
		}

		if line == pbxTargetDependencySectionBeginPattern {
			isPBXTargetDependencySection = true
			continue
		}

		if isPBXTargetDependencySection {
			pbxTargetDependencySectionLines = append(pbxTargetDependencySectionLines, line)
		}
		// -----

	}
	if err := scanner.Err(); err != nil {
		return map[string][]string{}, err
	}

	return isaObjectSectionMap, nil
}
