package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/godrei/xcodeproj/pbxproj"
)

func main() {
	pbxprojPth := "./_tmp/project.pbxproj"
	content, err := fileutil.ReadStringFromFile(pbxprojPth)
	if err != nil {
		fmt.Printf("failed to read file (%s), error: %s\n", pbxprojPth, err)
		os.Exit(1)
	}

	isaSectionLinesMap, err := pbxproj.SplitObjectsSections(content)
	if err != nil {
		fmt.Printf("failed to split project (%s), error: %s\n", pbxprojPth, err)
		os.Exit(1)
	}

	for isa, lines := range isaSectionLinesMap {
		fmt.Printf("%s\n", isa)
		fmt.Println("")
		for _, line := range lines {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("-----")
		fmt.Println("")

		if isa == "PBXNativeTarget" {
			pbxNativeTargets := pbxproj.ParsePBXNativeTargetSection(lines)
			for _, pbxNativeTarget := range pbxNativeTargets {
				fmt.Println("-----")
				fmt.Printf("%#v\n", pbxNativeTarget)
				fmt.Println("-----")
			}
		}
	}
}
