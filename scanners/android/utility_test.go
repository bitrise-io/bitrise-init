package android

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/assert"
)

type TestFileInfo struct {
	name string
}

func (t TestFileInfo) Name() string {
	return t.name
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
	pathSeparator := string(os.PathSeparator)
	rootPath := "1"
	paths := []string{rootPath, "2", "3", "4", "5", "5" + pathSeparator + "6"}
	filePathWalk = func(root string, walkFn filepath.WalkFunc) error {
		var skipPaths []string
		for _, path := range paths {
			skip := false
			for _, skipPath := range skipPaths {
				if strings.HasPrefix(path, skipPath) {
					skip = true
				}
			}
			if skip {
				continue
			}
			err := walkFn(path, TestFileInfo{name: path}, nil)
			if err != nil {
				if err == filepath.SkipDir {
					skipPaths = append(skipPaths, path)
				} else {
					return err
				}
			}
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
		skip       []string
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
		{
			name:       "Some child folders contains build.gradle and settings.gradle and one is on skip list",
			pathExists: buildMatcher(map[string][]string{paths[1]: {"build.gradle", "settings.gradle"}, paths[2]: {"build.gradle.kts", "settings.gradle.kts"}, paths[3]: {"build.gradle", "settings.gradle"}}),
			skip:       []string{paths[2]},
			expect:     []string{paths[1], paths[3]},
		},
		{
			name:       "Skipped directory's child contains build.gradle and settings.gradle",
			pathExists: buildMatcher(map[string][]string{paths[5]: {"build.gradle", "settings.gradle"}}),
			skip:       []string{paths[4]},
			expect:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			pathUtilIsPathExists = tc.pathExists
			// Act
			groups, err := walkMultipleFileGroups(rootPath, projectFiles, tc.skip)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, groups, tc.expect)
		})
	}
}

func Test_ContainsLocalProperties(t *testing.T) {
	testCases := []struct {
		name                    string
		containsLocalProperties bool
		want                    bool
	}{
		{
			name:                    "1. Given the project directory contains the local.properties file when the existance of the local.properties file is checked then expect TRUE to be returned",
			containsLocalProperties: true,
			want:                    true,
		},
		{
			name:                    "2. Given the project directory not contains the local.properties file when the existance of the local.properties file is checked then expect FLASE to be returned",
			containsLocalProperties: false,
			want:                    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			projectPath := createProjectDirectory(t, testCase.containsLocalProperties)

			// When
			exists, err := containsLocalProperties(projectPath)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, testCase.want, exists)
		})
	}
}

func buildMatcher(rootsAndPaths map[string][]string) func(path string) (bool, error) {
	return func(path string) (bool, error) {
		for key, paths := range rootsAndPaths {
			for _, p := range paths {
				if path == filepath.Join(key, p) {
					return true, nil
				}
			}
		}
		return false, nil
	}
}

func createProjectDirectory(t *testing.T, containsLocalProperties bool) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	assert.NoError(t, err)

	createFile(t, tmpDir, "build.gradle")

	if containsLocalProperties {
		createFile(t, tmpDir, "local.properties")
	}

	return tmpDir
}

func createFile(t *testing.T, path, fileName string) {
	_, err := os.Create(filepath.Join(path, fileName))
	assert.NoError(t, err)
}
