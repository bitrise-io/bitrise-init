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

func lookupResourceBasedOnManifest(manifestPth, resPth string) (string, error) {
	// Fetch icon name from AndroidManifest.xml
	var filenameBase string
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
		filenameBase = strings.TrimPrefix(ic.Value, `@mipmap/`)
	}
	{
		mipmapDirs := []string{"mipmap-xxxhdpi", "mipmap-xxhdpi", "mipmap-xhdpi", "mipmap-hdpi", "mipmap-mdpi", "mipmap-ldpi"}

		for _, dir := range mipmapDirs {
			filePath := path.Join(dir, filenameBase+".png")
			if exists, err := pathutil.IsPathExists(filePath); err != nil {
				return "", err
			} else if exists {
				return filePath, nil
			}
		}
	}
	return "", nil
}

// LookupPossibleMatches returns all potential android icons
func LookupPossibleMatches(projectDir string, basepath string) (models.Icons, error) {
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
				iconPath, err := lookupResourceBasedOnManifest(manifestPth, resourcesPth)
				if err != nil {
					return nil, err
				}
				if iconPath != "" {
					iconPaths = append(iconPaths, iconPath)
				}
			}
		}
	}
	icons, err := utility.ConvertPathsToUniqueFileNames(iconPaths, basepath)
	if err != nil {
		return nil, err
	}
	return icons, nil
}
