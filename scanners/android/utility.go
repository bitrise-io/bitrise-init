package android

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/pathutil"
)

type fileGroups [][]string

var pathUtilIsPathExists = pathutil.IsPathExists
var filePathWalk = filepath.Walk

// Project is an Android project on the filesystem
type Project struct {
	RelPath  string
	Icons    models.Icons
	Warnings models.Warnings
}

func walk(src string, fn func(path string, info os.FileInfo) error) error {
	return filePathWalk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == src {
			return nil
		}
		return fn(path, info)
	})
}

func checkFileGroups(path string, fileGroups fileGroups) (bool, error) {
	for _, fileGroup := range fileGroups {
		found := false
		for _, file := range fileGroup {
			exists, err := pathUtilIsPathExists(filepath.Join(path, file))
			if err != nil {
				return found, err
			}
			if exists {
				found = true
			}
		}
		if !found {
			return false, nil
		}
	}
	return true, nil
}

func walkMultipleFileGroups(searchDir string, fileGroups fileGroups, skipDirs []string) (matches []string, err error) {
	match, err := checkFileGroups(searchDir, fileGroups)
	if err != nil {
		return nil, err
	}
	if match {
		matches = append(matches, searchDir)
	}
	return matches, walk(searchDir, func(path string, info os.FileInfo) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if nameMatchSkipDirs(info.Name(), skipDirs) {
				return filepath.SkipDir
			}
			match, err := checkFileGroups(path, fileGroups)
			if err != nil {
				return err
			}
			if match {
				matches = append(matches, path)
			}
		}
		return nil
	})
}

func nameMatchSkipDirs(name string, skipDirs []string) bool {
	for _, skipDir := range skipDirs {
		if skipDir == "" {
			continue
		}
		if name == skipDir {
			return true
		}
	}
	return false
}

func containsLocalProperties(projectDir string) (bool, error) {
	return pathutil.IsPathExists(filepath.Join(projectDir, "local.properties"))
}

func checkGradlew(projectDir string) error {
	gradlewPth := filepath.Join(projectDir, "gradlew")
	exist, err := pathutil.IsPathExists(gradlewPth)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(`<b>No Gradle Wrapper (gradlew) found.</b>
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`)
	}
	return nil
}
