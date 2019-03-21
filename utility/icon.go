package utility

import (
	"crypto/sha256"
	"fmt"

	"github.com/bitrise-core/bitrise-init/models"
)

// ConvertPathsToUniqueFileNames returns a map whose values are the input array elements
// keys are a sha256 hash of input strings
func ConvertPathsToUniqueFileNames(appIconPaths []string) models.Icons {
	iconIDToPath := models.Icons{}
	for _, appIconPath := range appIconPaths {
		hash := sha256.Sum256([]byte(appIconPath))
		hashStr := fmt.Sprintf("%x", hash)
		iconIDToPath[hashStr] = appIconPath
	}
	return iconIDToPath
}
