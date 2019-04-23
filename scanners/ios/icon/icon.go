package icon

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/xcode-project/xcodeproj"
)

// LookupPossibleMatches returns possible ios app icons,
// in a map with key of a id (sha256 hash converted to string), value of icon path
func LookupPossibleMatches(projectPath string, schemeName string, basepath string) (models.Icons, error) {
	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open project file: %s, error: %s", projectPath, err)
	}

	log.Printf("name: %s", project.Name)

	scheme, found := project.Scheme(schemeName)
	if !found {
		return nil, fmt.Errorf("scheme (%s) not found in project", schemeName)
	}

	mainTarget, err := mainTargetOfScheme(project, scheme.Name)
	log.Printf("main target: %s", mainTarget.Name)

	appIconSetName, err := getAppIconSetName(project, mainTarget)
	if err != nil {
		return nil, fmt.Errorf("app icon set name not found in project, error: %s", err)
	}

	targetsToAssetCatalogs, err := xcodeproj.AssetCatalogs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset catalogs for project: %s, error: %s", projectPath, err)
	}
	assetCatalogPaths, ok := targetsToAssetCatalogs[mainTarget.ID]
	if !ok {
		return nil, fmt.Errorf("target not found in project")
	}

	appIconPath, found, err := lookupAppIconPath(projectPath, assetCatalogPaths, appIconSetName)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, err
	}

	log.Printf("%s", appIconPath)
	icon, found, err := parseResourceSet(appIconPath)
	if err != nil {
		return nil, fmt.Errorf("could not get iconm path: %s, error: %s", appIconPath, err)
	} else if !found {
		return nil, fmt.Errorf("icon not found, path: %s", appIconPath)
	}

	iconPath := filepath.Join(appIconPath, icon.Filename)

	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("icon file does not exist: %s, error: %err", iconPath, err)
	}

	iconIDToPath, err := utility.ConvertPathsToUniqueFileNames([]string{iconPath}, basepath)
	if err != nil {
		return nil, err
	}
	return iconIDToPath, nil
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

func parseResourceSet(resourceSetPath string) (appIcon, bool, error) {
	const resourceMetadataFileName = "Contents.json"
	file, err := os.Open(filepath.Join(resourceSetPath, resourceMetadataFileName))
	if err != nil {
		return appIcon{}, false, fmt.Errorf("failed to open file, error: %s", err)
	}

	appIcons, err := parseResourceSetMetadata(file)
	if err != nil {
		return appIcon{}, false, fmt.Errorf("failed to parse asset metadata, error: %s", err)
	}

	if len(appIcons) == 0 {
		return appIcon{}, false, nil
	}
	largestIcon := appIcons[0]
	for _, icon := range appIcons {
		if icon.Size > largestIcon.Size {
			largestIcon = icon
		}
	}
	return largestIcon, true, nil
}

func parseResourceSetMetadata(input io.Reader) ([]appIcon, error) {
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
