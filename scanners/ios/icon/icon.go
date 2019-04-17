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
	"github.com/bitrise-tools/xcode-project/xcodeproj"
)

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
	icon, err := parseResourceSet(appIconPath)
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

// LookupPossibleMatches returns possible ios app icons,
// in a map with key of a id (sha256 hash converted to string), value of icon path
func LookupPossibleMatches(searchPath string, basepath string) (models.Icons, error) {
	appIconSets, err := getResourceSetDirs(searchPath)
	if err != nil {
		return nil, err
	}

	var appIconPaths []string
	for _, appIconSetPath := range appIconSets {
		log.Printf("%s", appIconSetPath)
		icon, err := parseResourceSet(appIconSetPath)
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

	iconIDToPath, err := utility.ConvertPathsToUniqueFileNames(appIconPaths, basepath)
	if err != nil {
		return nil, err
	}
	return iconIDToPath, nil
}

func getResourceSetDirs(path string) ([]string, error) {
	appIconSets := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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

func parseResourceSet(appIconPath string) (*appIcon, error) {
	const metaDataFileName = "Contents.json"
	file, err := os.Open(filepath.Join(appIconPath, metaDataFileName))
	if err != nil {
		return nil, fmt.Errorf("failed to open file, error: %s", err)
	}

	appIcons, err := parseResourceSetMetadata(file)
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
