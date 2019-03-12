package icon

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
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

// // Version ...
// type Version struct {
// 	ID string
// }

// // Versions ...
// type Versions []Version

// func (v Versions) Len() int {
// 	return len(v)
// }
// func (v Versions) Swap(i, j int) {
// 	v[i], v[j] = v[j], v[i]
// }
// func (v Versions) Less(i, j int) bool {
// 	return len(v[i].ID) < len(v[j].ID)
// }

// FindIcon ...
func findIcon(manifestPth, resPth string) (string, error) {
	//
	// Fetch icon name from AndroidManifest.xml
	var iconName string
	{
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(manifestPth); err != nil {
			panic(err)
		}

		log.Printf("XML: %+v", doc)

		man := doc.SelectElement("manifest")
		app := man.SelectElement("application")
		ic := app.SelectAttr("android:icon")
		iconName = strings.TrimPrefix(ic.Value, `@mipmap/`)
	}

	//
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

// FetchIcon ...
func FetchIcon(appPath string) (string, error) {
	xmlPth := path.Join(appPath, "AndroidManifest.xml")
	resPth := path.Join(appPath, "res")
	return findIcon(xmlPth, resPth)
}
