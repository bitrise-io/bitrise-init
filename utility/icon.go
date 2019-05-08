package utility

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/sliceutil"
)

// ConvertPathsToUniqueFileNames returns a sorted array of unique icons.
func ConvertPathsToUniqueFileNames(absoluteIconPaths []string, basepath string) (models.Icons, error) {
	paths := sliceutil.UniqueStringSlice(absoluteIconPaths)
	sort.Strings(paths)

	icons := models.Icons{}
	for _, iconPath := range absoluteIconPaths {
		relativePath, err := filepath.Rel(basepath, iconPath)
		if err != nil {
			return nil, err
		}
		hash := sha256.Sum256([]byte(relativePath))
		hashStr := fmt.Sprintf("%x", hash) + filepath.Ext(iconPath)

		icons = append(icons, models.Icon{
			Filename: hashStr,
			Path:     iconPath,
		})
	}
	return icons, nil
}
