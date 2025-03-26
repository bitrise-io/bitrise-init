package kmp

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type DirEntry struct {
	os.DirEntry
	Path string
}

var ignoreDirs = []string{".git", ".gradle", ".idea", "build", ".kotlin", ".fleet"}

// TODO: ignore well known directories like .git, .gradle, .idea, build, etc.
func listDirEntries(root string, depth uint) ([]DirEntry, error) {
	var entries []DirEntry
	dirsToRead := []string{root}
	for i := 0; i < int(depth); i++ {
		var nextDirsToRead []string
		for _, dir := range dirsToRead {
			dirEntries, err := os.ReadDir(dir)
			if err != nil {
				// log.Warnf("Failed to read dir: %s", dir)
				continue
			}

			for _, entry := range dirEntries {
				if entry.IsDir() {
					if slices.Contains(ignoreDirs, entry.Name()) {
						continue
					}

					nextDirsToRead = append(nextDirsToRead, filepath.Join(dir, entry.Name()))
				}
				entries = append(entries, DirEntry{entry, path.Join(dir, entry.Name())})
			}
		}
		if len(nextDirsToRead) == 0 {
			break
		}
		dirsToRead = nextDirsToRead
	}

	slices.SortFunc(entries, func(a, b DirEntry) int {
		componentsA := strings.Split(a.Path, string(os.PathSeparator))
		componentsB := strings.Split(b.Path, string(os.PathSeparator))
		if len(componentsA) < len(componentsB) {
			return -1
		} else if len(componentsA) > len(componentsB) {
			return 1
		}
		for i := 0; i < len(componentsA); i++ {
			if componentsA[i] < componentsB[i] {
				return -1
			} else if componentsA[i] > componentsB[i] {
				return 1
			}
		}
		return 0
	})
	return entries, nil
}

func detectGradleConfigurationDirectories(repoEntries []DirEntry) ([]DirEntry, error) {
	var gradleConfigurationDirectories []DirEntry
	for _, repoEntry := range repoEntries {
		if !repoEntry.IsDir() {
			continue
		}
		if repoEntry.Name() == "gradle" {
			gradleConfigurationDirectories = append(gradleConfigurationDirectories, repoEntry)
		}
	}
	return gradleConfigurationDirectories, nil

}

func detectVersionCatalogFile(gradleConfigurationDirectoryPth string, repoEntries []DirEntry) *DirEntry {
	var versionCatalogFile *DirEntry
	// gradle/libs.versions.toml
	expectedPathComponentNum := strings.Count(gradleConfigurationDirectoryPth, string(os.PathSeparator)) + 1

	for _, repoEntry := range repoEntries {
		if !strings.HasPrefix(repoEntry.Path, gradleConfigurationDirectoryPth) {
			continue
		}

		entryPathComponentsNum := strings.Count(repoEntry.Path, string(os.PathSeparator))
		if entryPathComponentsNum > expectedPathComponentNum {
			break
		}
		if entryPathComponentsNum != expectedPathComponentNum {
			continue
		}
		if repoEntry.IsDir() {
			continue
		}
		if repoEntry.Name() == "libs.versions.toml" {
			versionCatalogFile = &repoEntry
			break
		}
	}

	return versionCatalogFile
}

func detectProjectGradleBuildScriptFiles(gradleProjectRootDirPth string, repoEntries []DirEntry) []DirEntry {
	var projectGradleBuildScriptFiles []DirEntry
	// composeApp/build.gradle.kts
	// server/build.gradle.kts
	expectedPathComponentNum := strings.Count(gradleProjectRootDirPth, string(os.PathSeparator)) + 2

	for _, repoEntry := range repoEntries {
		if !strings.HasPrefix(repoEntry.Path, gradleProjectRootDirPth) {
			continue
		}

		entryPathComponentsNum := strings.Count(repoEntry.Path, string(os.PathSeparator))
		if entryPathComponentsNum > expectedPathComponentNum {
			break
		}
		if entryPathComponentsNum != expectedPathComponentNum {
			continue
		}
		if repoEntry.IsDir() {
			continue
		}
		if repoEntry.Name() == "build.gradle.kts" || repoEntry.Name() == "build.gradle" {
			projectGradleBuildScriptFiles = append(projectGradleBuildScriptFiles, repoEntry)
		}
	}

	return projectGradleBuildScriptFiles
}

func detectComposeAppProjectDirectories(composeAppDirPth string, repoEntries []DirEntry) []DirEntry {
	var projectDirectories []DirEntry
	// composeApp/src/androidMain
	composeAppSrcDirPth := filepath.Join(composeAppDirPth, "src")
	expectedPathComponentNum := strings.Count(composeAppSrcDirPth, string(os.PathSeparator)) + 1

	for _, repoEntry := range repoEntries {
		if !strings.HasPrefix(repoEntry.Path, composeAppSrcDirPth) {
			continue
		}

		entryPathComponentsNum := strings.Count(repoEntry.Path, string(os.PathSeparator))
		if entryPathComponentsNum > expectedPathComponentNum {
			break
		}
		if entryPathComponentsNum != expectedPathComponentNum {
			continue
		}
		if !repoEntry.IsDir() {
			continue
		}
		if strings.HasSuffix(repoEntry.Name(), "Main") {
			projectDirectories = append(projectDirectories, repoEntry)
		}
	}

	return projectDirectories
}

func detectAnyDependencies(entry DirEntry, dependencies []string) (bool, error) {
	file, err := os.Open(entry.Path)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// log.Warnf("Failed to close file: %s", versionCatalogEntry.Path)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}

	for _, dependency := range dependencies {
		if strings.Contains(string(content), dependency) {
			return true, nil
		}
	}

	return false, nil
}

func printProjectStructure(projectStructure ProjectStructure) {
	content, err := json.MarshalIndent(projectStructure, "", "  ")
	if err == nil {
		log.Println(string(content))
	}
}
