package android

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type TestFileInfo struct {
}

func (t TestFileInfo) Name() string {
	panic("implement me")
}

func (t TestFileInfo) Size() int64 {
	panic("implement me")
}

func (t TestFileInfo) Mode() os.FileMode {
	panic("implement me")
}

func (t TestFileInfo) ModTime() time.Time {
	panic("implement me")
}

func (t TestFileInfo) IsDir() bool {
	return true
}

func (t TestFileInfo) Sys() interface{} {
	panic("implement me")
}

func TestWalkMultipleFileGroups(t *testing.T) {
	rootPath := "1"
	paths := []string{rootPath, "2", "3", "4", "5"}
	filePathWalk = func(root string, walkFn filepath.WalkFunc) error {
		for _, path := range paths {
			_ = walkFn(path, TestFileInfo{}, nil)
		}
		return nil
	}

	projectFiles := fileGroups{
		{"build.gradle", "build.gradle.kts"},
		{"settings.gradle", "settings.gradle.kts"},
	}

	testCases := []struct {
		name       string
		pathExists func(string) (bool, error)
		expect     []string
	}{
		{
			name:       "Root folder contains build.gradle and settings.gradle",
			pathExists: buildMatcher(map[string][]string{rootPath: {"build.gradle", "settings.gradle"}}),
			expect:     []string{rootPath},
		},
		{
			name:       "Root folder contains build.gradle.kts and settings.gradle.kts",
			pathExists: buildMatcher(map[string][]string{rootPath: {"build.gradle.kts", "settings.gradle.kts"}}),
			expect:     []string{rootPath},
		},
		{
			name:       "Root folder contains build.gradle and settings.gradle.kts",
			pathExists: buildMatcher(map[string][]string{rootPath: {"build.gradle", "settings.gradle.kts"}}),
			expect:     []string{rootPath},
		},
		{
			name:       "Non-root folder contains build.gradle and settings.gradle",
			pathExists: buildMatcher(map[string][]string{paths[1]: {"build.gradle.kts", "settings.gradle.kts"}}),
			expect:     []string{paths[1]},
		},
		{
			name:       "Non-root folder contains build.gradle.kts and settings.gradle.kts",
			pathExists: buildMatcher(map[string][]string{paths[2]: {"build.gradle.kts", "settings.gradle.kts"}}),
			expect:     []string{paths[2]},
		},
		{
			name:       "Non-root folder contains build.gradle.kts and settings.gradle",
			pathExists: buildMatcher(map[string][]string{paths[2]: {"build.gradle.kts", "settings.gradle"}}),
			expect:     []string{paths[2]},
		},
		{
			name:       "Root folder and child folder contains build.gradle and settings.gradle",
			pathExists: buildMatcher(map[string][]string{rootPath: {"build.gradle", "settings.gradle"}, paths[2]: {"build.gradle.kts", "settings.gradle.kts"}}),
			expect:     []string{rootPath, paths[2]},
		},
		{
			name:       "Two child folders contains build.gradle and settings.gradle",
			pathExists: buildMatcher(map[string][]string{paths[1]: {"build.gradle", "settings.gradle"}, paths[2]: {"build.gradle.kts", "settings.gradle.kts"}}),
			expect:     []string{paths[1], paths[2]},
		},
		{
			name:       "No folder contains any gradle files",
			pathExists: buildMatcher(map[string][]string{}),
			expect:     nil,
		},
		{
			name:       "Root folder only contains settings.gradle",
			pathExists: buildMatcher(map[string][]string{rootPath: {"settings.gradle"}}),
			expect:     nil,
		},
		{
			name:       "Root folder only contains build.gradle",
			pathExists: buildMatcher(map[string][]string{rootPath: {"build.gradle"}}),
			expect:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			pathUtilIsPathExists = tc.pathExists
			// Act
			groups, _ := walkMultipleFileGroups(rootPath, projectFiles)
			// Assert
			assert.Equal(t, groups, tc.expect)
		})
	}
}

func buildMatcher(rootsAndPaths map[string][]string) func(path string) (bool, error) {
	return func(path string) (bool, error) {
		for key, paths := range rootsAndPaths {
			for _, p := range paths {
				if path == fmt.Sprintf("%s%c%s", key, filepath.Separator, p) {
					return true, nil
				}
			}
		}
		return false, nil
	}
}
