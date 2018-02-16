package jekyll

import (
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/fileutil"
)

const (
	// ScannerName ...
	ScannerName   = "jekyll"
	configYmlFile = "_config.yml"
	gemfileFile   = "Gemfile"
	jekyllGemName = "jekyll"
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

// ParseConfigXML ...
func readGemfileToString() (string, error) {
	content, err := fileutil.ReadStringFromFile(gemfileFile)
	if err != nil {
		return "", err
	}
	return content, nil
}