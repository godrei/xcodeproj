package pbxproj

const (
	pbxTargetDependencyISA                 = "PBXTargetDependency"
	pbxTargetDependencySectionBeginPattern = `/* Begin PBXTargetDependency section */`
	pbxTargetDependencySectionEndPattern   = `/* End PBXTargetDependency section */`
)

// PBXTargetDependency ...
type PBXTargetDependency struct {
}
