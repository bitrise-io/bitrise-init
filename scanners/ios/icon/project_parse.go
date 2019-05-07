package icon

import (
	"fmt"

	"github.com/bitrise-io/xcode-project/xcodeproj"
	"github.com/bitrise-io/xcode-project/xcworkspace"
)

func mainTargetOfScheme(proj xcodeproj.XcodeProj, scheme string) (xcodeproj.Target, error) {
	projTargets := proj.Proj.Targets
	sch, ok := proj.Scheme(scheme)
	if !ok {
		return xcodeproj.Target{}, fmt.Errorf("Failed to find scheme (%s) in project", scheme)
	}

	var blueIdent string
	for _, entry := range sch.BuildAction.BuildActionEntries {
		if entry.BuildableReference.IsAppReference() {
			blueIdent = entry.BuildableReference.BlueprintIdentifier
			break
		}
	}

	// Search for the main target
	for _, t := range projTargets {
		if t.ID == blueIdent {
			return t, nil

		}
	}
	return xcodeproj.Target{}, fmt.Errorf("failed to find the project's main target for scheme (%s)", scheme)
}

func targetByName(proj xcodeproj.XcodeProj, target string) (xcodeproj.Target, bool, error) {
	projTargets := proj.Proj.Targets
	for _, t := range projTargets {
		if t.Name == target {
			return t, true, nil
		}
	}
	return xcodeproj.Target{}, false, nil
}

// ProjectPathByScheme ...
func ProjectPathByScheme(workspacePath string, scheme string) (string, bool, error) {
	workspace, err := xcworkspace.Open(workspacePath)
	if err != nil {
		return "", false, err
	}

	_, projectPath, err := workspace.Scheme(scheme)
	if err != nil {
		return "", false, err
	}

	return projectPath, true, nil
}
