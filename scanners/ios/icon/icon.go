package icon

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xcode-project/xcodeproj"
)

func getIcon(projectPath string, scheme string) (string, error) {
	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to open project file: %s, error: %s", projectPath, err)
	}

	log.Printf("name: %s", project.Name)

	scheme, found := project.Scheme(schemeName)
	if !found {
		return "", fmt.Errorf("scheme (%s) not found in project", schemeName)
	}

	mainTarget, err := mainTargetOfScheme(project, scheme.Name)
	log.Printf("main target: %s", mainTarget.Name)

	assetCatalogPath, err := getAssetCatalogPath(projectPath, scheme)
	if err != nil {
		return "", fmt.Errorf("failed to get asset catalog path, error: %s", err)
	}

	openAssetCatalog(assetCatalogPath, projectPath)
	return "", nil
}

func openAssetCatalog(assetCatalogName string, projectPath string) (string, error) {
	return "", nil
}

// mainTargetOfScheme return the main target
func mainTargetOfScheme(proj xcodeproj.XcodeProj, scheme string) (xcodeproj.Target, error) {
	projTargets := proj.Proj.Targets
	sch, ok := proj.Scheme(scheme)
	if !ok {
		return xcodeproj.Target{}, fmt.Errorf("Failed to found scheme (%s) in project", scheme)
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

const appIconSetNameKey = "ASSETCATALOG_COMPILER_APPICON_NAME"

func getAppIconSetName(project xcodeproj.XcodeProj, target xcodeproj.Target) (string, error) {
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

func getAssetCatalogPaths(project xcodeproj.XcodeProj, target xcodeproj.Target) ([]string, error) {
	log.Printf("assets in project: %v+", project.Proj.TargetToAssetCatalogs)
	log.Printf("target ID: %s", target.ID)
	assetCatalogs, ok := project.Proj.TargetToAssetCatalogs[target.ID]
	if !ok {
		return nil, fmt.Errorf("asset catalog path not found in project")
	}

	log.Printf("asset catalog path: %s", assetCatalogs)
	return assetCatalogs, nil
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
