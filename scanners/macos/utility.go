package macos

import (
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
)

const packageBase = "Package.swift"

// Detect ...
func Detect(searchDir string) (bool, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return false, err
	}

	log.TInfof("Filter relevant Package.swift files")

	filters := []utility.FilterFunc{
		utility.BaseFilter(packageBase, true),
	}
	packageFileList, err := utility.FilterPaths(fileList, filters...)
	if err != nil {
		return false, err
	}

	if len(packageFileList) == 0 {
		log.TPrintf("platform not detected")
		return false, nil
	}

	log.TSuccessf("Platform detected")

	return true, nil
}
