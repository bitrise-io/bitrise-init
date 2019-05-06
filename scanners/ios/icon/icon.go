package icon

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/xcode-project/xcodeproj"
)

// LookupByScheme returns possible ios app icons for a scheme,
// Icons key: unique id for relative paths under basepath(sha256 hash converted to string) as a filename,
// with the original (png) file extension appended
// Icons value: absolute icon path
func LookupByScheme(projectPath string, schemeName string, basepath string) (models.Icons, error) {
	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open project file: %s, error: %s", projectPath, err)
	}

	scheme, found := project.Scheme(schemeName)
	if !found {
		return nil, fmt.Errorf("scheme (%s) not found in project", schemeName)
	}

	mainTarget, err := mainTargetOfScheme(project, scheme.Name)

	return lookupByTarget(projectPath, mainTarget, basepath)
}

// LookupByTarget returns possible ios app icons for a scheme,
// Icons key: unique id for relative paths under basepath(sha256 hash converted to string) as a filename,
// with the original (png) file extension appended
// Icons value: absolute icon path
func LookupByTarget(projectPath string, targetName string, basepath string) (models.Icons, error) {
	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open project file: %s, error: %s", projectPath, err)
	}

	target, found, err := targetByName(project, targetName)
	if err != nil {
		return models.Icons{}, err
	} else if !found {
		return models.Icons{}, fmt.Errorf("not found target: %s, in project: %s", targetName, projectPath)
	}

	return lookupByTarget(projectPath, target, basepath)
}

func lookupByTarget(projectPath string, target xcodeproj.Target, basepath string) (models.Icons, error) {
	targetToAppIconSetPaths, err := xcodeproj.AppIconSetPaths(projectPath)
	appIconSetPaths, ok := targetToAppIconSetPaths[target.ID]
	if !ok {
		return nil, fmt.Errorf("target not found in project")
	}

	iconPaths := []string{}
	for _, appIconSetPath := range appIconSetPaths {
		icon, found, err := parseResourceSet(appIconSetPath)
		if err != nil {
			return nil, fmt.Errorf("could not get icon, error: %s", err)
		} else if !found {
			return nil, nil
		}
		log.Debugf("App icons: %s", icon)

		iconPath := filepath.Join(appIconSetPath, icon.Filename)

		if _, err := os.Stat(iconPath); err != nil && os.IsNotExist(err) {
			return nil, fmt.Errorf("icon file does not exist: %s, error: %err", iconPath, err)
		}
		iconPaths = append(iconPaths, iconPath)
	}

	iconIDToPath, err := utility.ConvertPathsToUniqueFileNames(iconPaths, basepath)
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
