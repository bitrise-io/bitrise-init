package jekyll

import (
	"github.com/bitrise-core/bitrise-init/utility"
)

const (
	// ScannerName ...
	ScannerName = "jekyll"
	configYmlFile = "_config.yml"
	gemfile = "Gemfile"
)

// filterProjectFile ...
func filterProjectFile(fileName string, fileList []string) (string, error) {
	allowGivenFileBaseFilter := utility.BaseFilter(fileName, true)
	filePaths, err := utility.FilterPaths(fileList, allowGivenFileBaseFilter)
	if err != nil {
		return "", err
	}

	if len(filePaths) == 0 {
		return "", nil
	}

	return filePaths[0], nil
}
