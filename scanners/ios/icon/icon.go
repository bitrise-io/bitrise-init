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
)

// GetAllIcons returns possible ios app icons,
// in a map with key of a id (sha256 hash converted to string), value of icon path
func GetAllIcons(searchPath string, basepath string) (models.Icons, error) {
	appIconSets, err := getAppIconSetDirs(searchPath)
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

	iconIDToPath, err := utility.ConvertPathsToUniqueFileNames(appIconPaths, basepath)
	if err != nil {
		return nil, err
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
