package utility

import (
	"sort"
	"strings"
)

const (
	buildGradleBasePath = "build.gradle"
	gradlewBasePath     = "gradlew"
)

// FixedGradlewPath ...
func FixedGradlewPath(gradlewPth string) string {
	split := strings.Split(gradlewPth, "/")
	if len(split) != 1 {
		return gradlewPth
	}

	if !strings.HasPrefix(gradlewPth, "./") {
		return "./" + gradlewPth
	}
	return gradlewPth
}

// FilterRootBuildGradleFiles ...
func FilterRootBuildGradleFiles(fileList []string) ([]string, error) {
	allowBuildGradleBaseFilter := BaseFilter(buildGradleBasePath, true)
	gradleFiles, err := FilterPaths(fileList, allowBuildGradleBaseFilter)
	if err != nil {
		return []string{}, err
	}

	if len(gradleFiles) == 0 {
		return []string{}, nil
	}

	sortableFiles := []SortablePath{}
	for _, pth := range gradleFiles {
		sortable, err := NewSortablePath(pth)
		if err != nil {
			return []string{}, err
		}
		sortableFiles = append(sortableFiles, sortable)
	}

	sort.Sort(BySortablePathComponents(sortableFiles))
	mindDepth := len(sortableFiles[0].Components)

	rootGradleFiles := []string{}
	for _, sortable := range sortableFiles {
		depth := len(sortable.Components)
		if depth == mindDepth {
			rootGradleFiles = append(rootGradleFiles, sortable.Pth)
		}
	}

	return rootGradleFiles, nil
}
