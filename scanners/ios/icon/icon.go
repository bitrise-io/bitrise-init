package icon

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xcode-project/xcodeproj"
)

// GetAllIcons returns a map with key of a id (sha256 hash converted to string), value of icon path
func GetAllIcons(path string) (models.Icons, error) {
	appIconSets, err := getAppIconSetDirs(path)
	if err != nil {
		return nil, err
	}

	var appIconPaths []string
	for _, appIconSetPath := range appIconSets {
		log.Printf("%s", appIconSetPath)
		icon, err := parseAppIconSet(appIconSetPath)
		if err != nil {
			return nil, fmt.Errorf("could not get icon, error: %s", err)
		} else if icon == nil {
			continue
		}

		iconPath := filepath.Join(appIconSetPath, icon.Filename)

		if _, err := os.Stat(iconPath); os.IsNotExist(err) {
			log.Errorf("Can not open icon file: %s, error: %err", iconPath, err)
			continue
		}
		appIconPaths = append(appIconPaths, iconPath)
	}

	iconIDToPath := models.Icons{}
	for _, appIconPath := range appIconPaths {
		hash := sha256.Sum256([]byte(appIconPath))
		hashStr := fmt.Sprintf("%x", hash)
		iconIDToPath[hashStr] = appIconPath
	}
	return iconIDToPath, nil
}

func getAppIconSetDirs(path string) ([]string, error) {
	appIconSets := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && strings.HasSuffix(path, ".appiconset") {
			appIconSets = append(appIconSets, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk path %s, error: %s", path, err)
	}
	log.Debugf("%s", appIconSets)

	return appIconSets, nil
}

func getIcon(projectPath string, schemeName string) (string, error) {
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

	appIconSetName, err := getAppIconSetName(project, mainTarget)
	if err != nil {
		return "", fmt.Errorf("app icon set name not found in project, error: %s", err)
	}

	assetCatalogPaths, err := getAssetCatalogPaths(project, mainTarget)
	if err != nil {
		return "", fmt.Errorf("failed to get asset catalog paths, error: %s", err)
	}

	appIconPath, found, err := lookupAppIconPath(projectPath, assetCatalogPaths, appIconSetName)
	if err != nil {
		return "", err
	} else if !found {
		return "", err
	}

	log.Printf("%s", appIconPath)
	icon, err := parseAppIconSet(appIconPath)
	if err != nil {
		return "", fmt.Errorf("could not get icon: ")
	}

	iconPath := filepath.Join(appIconPath, icon.Filename)
	_, err = os.Open(iconPath)
	if err != nil {
		return "", fmt.Errorf("Can not open icon file: %s, error: %err", iconPath, err)
	}

	return icon.Filename, nil
}

type assetIcon struct {
	Size     string
	Filename string
}

type assetInfo struct {
	Version int
	Author  string
}
type assetMetadata struct {
	Images []assetIcon
	Info   assetInfo
}

type appIcon struct {
	Size     int
	Filename string
}

func parseAppIconSet(appIconPath string) (*appIcon, error) {
	const metaDataFileName = "Contents.json"
	file, err := os.Open(filepath.Join(appIconPath, metaDataFileName))
	if err != nil {
		return nil, fmt.Errorf("failed to open file, error: %s", err)
	}

	appIcons, err := parseAppIconMetadata(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse asset metadata, error: %s", err)
	}

	if len(appIcons) == 0 {
		return nil, nil
	}
	largestIcon := appIcons[0]
	for _, icon := range appIcons {
		if icon.Size > largestIcon.Size {
			largestIcon = icon
		}
	}
	return &largestIcon, nil
}

func parseAppIconMetadata(input io.Reader) ([]appIcon, error) {
	decoder := json.NewDecoder(input)
	var decoded assetMetadata
	err := decoder.Decode(&decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode asset metadata file, error: %s", err)
	}

	if decoded.Info.Version != 1 {
		return nil, fmt.Errorf("unsupported asset metadata version")
	}
	var icons []appIcon
	for _, icon := range decoded.Images {
		fmt.Printf("%s", icon)
		fmt.Println()
		sizeParts := strings.Split(icon.Size, "x")
		if len(sizeParts) != 2 {
			return nil, fmt.Errorf("invalid image size format")
		}
		iconSize, err := strconv.ParseFloat(sizeParts[0], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid image size, error: %s", err)
		}
		// If icon is not set for a given usage, filaname key is missing
		if icon.Filename != "" {
			icons = append(icons, appIcon{
				Size:     int(iconSize),
				Filename: icon.Filename,
			})
		}
	}
	return icons, nil
}

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
