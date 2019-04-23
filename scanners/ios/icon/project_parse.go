package icon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/xcode-project/xcodeproj"
)

// ToDo: use file paths based on xcode project
func lookupAppIconPath(projectPath string, assetCatalogPaths []string, appIconSetName string) (string, bool, error) {
	projectDir := strings.TrimSuffix(projectPath, ".xcodeproj")
	for _, assetCatalogPath := range assetCatalogPaths {
		var matches []string
		err := filepath.Walk(projectDir, func(path string, f os.FileInfo, err error) error {
			if _, name := filepath.Split(path); name == assetCatalogPath {
				matches = append(matches, path)
			}
			return nil
		})
		if err != nil {
			return "", false, err
		}

		log.Printf("%s %s", assetCatalogPath, matches)
		if len(matches) > 0 {
			iconSetMatches, err := filepath.Glob(filepath.Join(matches[0], appIconSetName+".appiconset"))
			if err != nil {
				return "", false, err
			}
			if len(iconSetMatches) > 0 {
				return iconSetMatches[0], true, nil
			}
		}
	}
	return "", false, nil
}

// mainTargetOfScheme return the main target
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

func getAppIconSetName(project xcodeproj.XcodeProj, target xcodeproj.Target) (string, error) {
	const appIconSetNameKey = "ASSETCATALOG_COMPILER_APPICON_NAME"

	found, defaultConfiguration := defaultConfiguration(target)
	if !found {
		return "", fmt.Errorf("default configuraion not founf for target: %s", target)
	}

	log.Printf("%s", defaultConfiguration)

	appIconSetNameRaw, ok := defaultConfiguration.BuildSettings[appIconSetNameKey]
	if !ok {
		return "", nil
	}

	appIconSetName, ok := appIconSetNameRaw.(string)
	if !ok {
		return "", fmt.Errorf("type assertion failed for value of key %s", appIconSetNameKey)
	}
	log.Printf("asstets: %s", appIconSetName)
	return appIconSetName, nil
}

func defaultConfiguration(target xcodeproj.Target) (bool, xcodeproj.BuildConfiguration) {
	defaultConfigurationName := target.BuildConfigurationList.DefaultConfigurationName
	for _, configuration := range target.BuildConfigurationList.BuildConfigurations {
		if configuration.Name == defaultConfigurationName {
			return true, configuration
		}
	}
	return false, xcodeproj.BuildConfiguration{}
}
