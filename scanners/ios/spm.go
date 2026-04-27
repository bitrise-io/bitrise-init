package ios

import (
	"fmt"

	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/pathutil"
)

func HasSPMDependencies(fileList []string) (bool, error) {
	// Pure SPM projects: find project-level `Package.swift` file
	pureSwiftFilters := append(utility.CommonExcludeFilters(), pathutil.BaseFilter("Package.swift", true))
	matches, err := pathutil.FilterPaths(fileList, pureSwiftFilters...)
	if err != nil {
		return false, fmt.Errorf("couldn't detect SPM dependencies: %w", err)
	}

	if len(matches) > 0 {
		return true, nil
	}

	// Xcode projects: find lockfile inside `xcodeproj`
	xcodeFilters := append(utility.CommonExcludeFilters(), pathutil.BaseFilter("Package.resolved", true))
	matches, err = pathutil.FilterPaths(fileList, xcodeFilters...)
	if err != nil {
		return false, fmt.Errorf("couldn't detect SPM dependencies: %w", err)
	}
	return len(matches) > 0, nil
}
