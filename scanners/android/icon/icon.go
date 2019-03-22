package icon

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

func findIcon(manifestPth, resPth string) (string, error) {
	// Fetch icon name from AndroidManifest.xml
	var iconName string
	{
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(manifestPth); err != nil {
			return "", err
		}

		man := doc.SelectElement("manifest")
		if man == nil {
			log.TPrintf("Key manifest not found in manifest file")
			return "", nil
		}
		app := man.SelectElement("application")
		if app == nil {
			log.TPrintf("Key application not found in manifest file")
			return "", nil
		}
		ic := app.SelectAttr("android:icon")
		if ic == nil {
			log.TPrintf("Attribute not found in manifest file")
			return "", nil
		}
		iconName = strings.TrimPrefix(ic.Value, `@mipmap/`)
	}
	{
		mipmapDirs := []string{"mipmap-anydpi*", "mipmap-xxxhdpi", "mipmap-xxhdpi", "mipmap-xhdpi", "mipmap-hdpi", "mipmap-mdpi", "mipmap-ldpi"}

		for _, dir := range mipmapDirs {
			pths, err := pathsByPattern(resPth, dir)
			if err != nil {
				return "", err
			}

			for _, pth := range pths {
				if exists, err := pathutil.IsPathExists(path.Join(pth, iconName+".png")); err != nil {
					continue
				} else if exists {
					return path.Join(pth, iconName+".png"), nil
				}
			}
		}
	}
	return "", nil
}

func pathsByPattern(paths ...string) ([]string, error) {
	pattern := filepath.Join(paths...)
	return filepath.Glob(pattern)
}

// GetAllIcons returns all potential android icons
func GetAllIcons(projectDir string) (models.Icons, error) {
	children, err := ioutil.ReadDir(projectDir)
	if err != nil {
		return nil, err
	}

	var iconPaths []string
	for _, object := range children {
		if object.IsDir() {
			manifestPth := filepath.Join(projectDir, object.Name(), "src", "main", "AndroidManifest.xml")
			resourcesPth := filepath.Join(projectDir, object.Name(), "src", "main", "res")
			if exist, err := pathutil.IsPathExists(manifestPth); err != nil {
				return nil, err
			} else if exist {
				iconPath, err := findIcon(manifestPth, resourcesPth)
				if err != nil {
					return nil, err
				}
				if iconPath != "" {
					iconPaths = append(iconPaths, iconPath)
				}
			}
		}
	}
	return utility.ConvertPathsToUniqueFileNames(iconPaths), nil
}
