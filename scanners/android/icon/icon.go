package icon

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/bitrise-core/bitrise-init/utility"

	"github.com/beevik/etree"
	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-io/go-utils/pathutil"
)

// // UsesPermission ...
// type UsesPermission struct {
// 	Name string `xml:"name,attr"`
// }

// // Application ...
// type Application struct {
// 	Icon string `xml:"android:icon,attr"`
// }

// // Manifest ...
// type Manifest struct {
// 	UsesPermission UsesPermission
// 	Application    Application `xml:"application"`
// }

// FindIcon ...
func findIcon(manifestPth, resPth string) (string, error) {
	//
	// Fetch icon name from AndroidManifest.xml
	var iconName string
	{
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(manifestPth); err != nil {
			return "", err
		}

		log.Printf("XML: %+v", doc)

		man := doc.SelectElement("manifest")
		if man == nil {
			return "", fmt.Errorf("key manifest not found in manifest file")
		}
		app := man.SelectElement("application")
		if app == nil {
			return "", fmt.Errorf("key application not found in manifest file")
		}
		ic := app.SelectAttr("android:icon")
		if ic == nil {
			return "", fmt.Errorf("attribute not found in manifest file")
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
	return "", fmt.Errorf("could not found any .png icon")
}

func pathsByPattern(paths ...string) ([]string, error) {
	pattern := filepath.Join(paths...)
	return filepath.Glob(pattern)
}

func fetchIcon(appPath string) (string, error) {
	xmlPth := path.Join(appPath, "AndroidManifest.xml")
	resPth := path.Join(appPath, "res")
	return findIcon(xmlPth, resPth)
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
				iconPaths = append(iconPaths, iconPath)
			}
		}
	}
	return utility.ConvertPathsToUniqueFileNames(iconPaths), nil
}
