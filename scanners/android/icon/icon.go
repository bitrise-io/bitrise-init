package icon

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/beevik/etree"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/sliceutil"
)

type icon struct {
	prefix       string
	fileNameBase string
}

func lookupIconName(manifestPth string) (icon, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(manifestPth); err != nil {
		return icon{}, err
	}

	log.Debugf("Looking for app icons. Manifest path: %s", manifestPth)
	parsedIcon, err := parseIconName(doc)
	if err != nil {
		return icon{}, err
	}

	return parsedIcon, nil
}

// parseIconName fetches icon name from AndroidManifest.xml
func parseIconName(doc *etree.Document) (icon, error) {
	man := doc.SelectElement("manifest")
	if man == nil {
		log.Debugf("Key manifest not found in manifest file")
		return icon{}, nil
	}
	app := man.SelectElement("application")
	if app == nil {
		log.Debugf("Key application not found in manifest file")
		return icon{}, nil
	}
	ic := app.SelectAttr("android:icon")
	if ic == nil {
		log.Debugf("Attribute not found in manifest file")
		return icon{}, nil
	}

	iconPathParts := strings.Split(strings.TrimPrefix(ic.Value, "@"), "/")
	if len(iconPathParts) != 2 {
		return icon{}, fmt.Errorf("unsupported icon key")
	}
	return icon{
		prefix:       iconPathParts[0],
		fileNameBase: iconPathParts[1],
	}, nil
}

func lookupIconPaths(resPth string, icon icon) ([]string, error) {
	var resourceSuffixes = [...]string{"xxxhdpi", "xxhdpi", "xhdpi", "hdpi", "mdpi", "ldpi"}
	resourceDirs := make([]string, len(resourceSuffixes))
	for _, mipmapSuffix := range resourceSuffixes {
		resourceDirs = append(resourceDirs, icon.prefix+"-"+mipmapSuffix)
	}

	for _, dir := range resourceDirs {
		iconPaths, err := filepath.Glob(filepath.Join(regexp.QuoteMeta(resPth), dir, icon.fileNameBase+".png"))
		if err != nil {
			return nil, err
		}
		if len(iconPaths) != 0 {
			return iconPaths, nil
		}
	}
	return nil, nil
}

func lookupPossibleMatches(projectDir string, basepath string) ([]string, error) {
	variantPaths := filepath.Join(regexp.QuoteMeta(projectDir), "*", "src", "*")
	manifestPaths, err := filepath.Glob(filepath.Join(variantPaths, "AndroidManifest.xml"))
	if err != nil {
		return nil, err
	}
	resourcesPaths, err := filepath.Glob(filepath.Join(variantPaths, "res"))
	if err != nil {
		return nil, err
	}

	iconNames := []icon{
		{
			prefix:       "mipmap",
			fileNameBase: "ic_launcher",
		},
		{
			prefix:       "mipmap",
			fileNameBase: "ic_launcher_round",
		},
	}
	for _, manifestPath := range manifestPaths {
		icon, err := lookupIconName(manifestPath)
		if err != nil {
			return nil, err
		}
		iconNames = append(iconNames, icon)
	}

	var iconPaths []string
	for _, resourcesPath := range resourcesPaths {
		for _, icon := range iconNames {
			foundIconPaths, err := lookupIconPaths(resourcesPath, icon)
			if err != nil {
				return nil, err
			}
			iconPaths = append(iconPaths, foundIconPaths...)
		}
	}
	return sliceutil.UniqueStringSlice(iconPaths), nil
}

// LookupPossibleMatches returns the largest resolution for all potential android icons
// It does look up all possible files project_dir/*/src/*/AndroidManifest.xml,
// then looks up the icon referenced in the res directory
func LookupPossibleMatches(projectDir string, basepath string) (models.Icons, error) {
	iconPaths, err := lookupPossibleMatches(projectDir, basepath)
	if err != nil {
		return nil, err
	}

	icons, err := utility.ConvertPathsToUniqueFileNames(iconPaths, basepath)
	if err != nil {
		return nil, err
	}
	return icons, nil
}
