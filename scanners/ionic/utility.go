package ionic

import (
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/pathutil"
)

// FilterRootFile ...
func FilterRootFile(fileList []string, fileName string) (string, error) {
	filters := append(utility.CommonExcludeFilters(), pathutil.BaseFilter(fileName, true))
	files, err := pathutil.FilterPaths(fileList, filters...)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	return files[0], nil
}
