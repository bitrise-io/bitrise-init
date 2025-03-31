package kmp

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func detectGradleConfigurationDirectories(repoEntries []DirEntry) ([]DirEntry, error) {
	var gradleConfigurationDirectories []DirEntry
	for _, repoEntry := range repoEntries {
		if !repoEntry.IsDir {
			continue
		}
		if repoEntry.Name == "gradle" {
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
		if repoEntry.IsDir {
			continue
		}
		if repoEntry.Name == "libs.versions.toml" {
			versionCatalogFile = &repoEntry
			break
		}
	}

	return versionCatalogFile
}

func detectSettingsGradleFile(gradleProjectRootDirPth string, repoEntries []DirEntry) *DirEntry {
	var settingsGradleFile *DirEntry
	expectedPathComponentNum := strings.Count(gradleProjectRootDirPth, string(os.PathSeparator)) + 1

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
		if repoEntry.IsDir {
			continue
		}
		if repoEntry.Name == "settings.gradle.kts" || repoEntry.Name == "settings.gradle" {
			settingsGradleFile = &repoEntry
			break
		}
	}

	return settingsGradleFile
}

func detectGradleBuildScriptFiles(gradleProjectRootDirPth string, repoEntries []DirEntry) []DirEntry {
	var projectGradleBuildScriptFiles []DirEntry
	// composeApp/build.gradle.kts
	// server/build.gradle.kts
	minPathComponentNum := strings.Count(gradleProjectRootDirPth, string(os.PathSeparator)) + 1

	for _, repoEntry := range repoEntries {
		if !strings.HasPrefix(repoEntry.Path, gradleProjectRootDirPth) {
			continue
		}

		entryPathComponentsNum := strings.Count(repoEntry.Path, string(os.PathSeparator))
		if entryPathComponentsNum < minPathComponentNum {
			break
		}
		if repoEntry.IsDir {
			continue
		}
		if repoEntry.Name == "build.gradle.kts" || repoEntry.Name == "build.gradle" {
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
		if !repoEntry.IsDir {
			continue
		}
		if strings.HasSuffix(repoEntry.Name, "Main") {
			projectDirectories = append(projectDirectories, repoEntry)
		}
	}

	return projectDirectories
}

func detectProjectIncludes(settingGradleFile DirEntry) ([]string, error) {
	file, err := os.Open(settingGradleFile.Path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// log.Warnf("Failed to close file: %s", settingGradleFile.Path)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return detectProjectIncludesInContent(string(content)), nil
}

func detectProjectIncludesInContent(settingGradleFileContent string) []string {
	var includedProjects []string
	lines := strings.Split(settingGradleFileContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "include") {
			continue
		}

		includedModules := strings.TrimPrefix(line, "include")
		includedModules = strings.Trim(includedModules, "()")
		includedModulesSplit := strings.Split(includedModules, ",")

		for _, includedModule := range includedModulesSplit {
			includedModule = strings.TrimSpace(includedModule)
			includedModule = strings.Trim(includedModule, `"'`)
			if !strings.HasPrefix(includedModule, ":") {
				includedModule = ":" + includedModule
			}
			includedProjects = append(includedProjects, includedModule)
		}
	}
	sort.Strings(includedProjects)

	return includedProjects
}

func detectProjectDirPath(gradleProjectRootDirPth, projectInclude string, repoEntries []DirEntry) string {
	projectInclude = strings.TrimPrefix(projectInclude, ":")
	projectIncludeComponents := strings.Split(projectInclude, ":")

	projectDirPath := gradleProjectRootDirPth
	for _, projectIncludeComponent := range projectIncludeComponents {
		if projectIncludeComponent == "" {
			continue
		}
		projectIncludeComponent = strings.TrimSpace(projectIncludeComponent)
		projectDirPath += "/" + projectIncludeComponent
	}

	expectedPathComponentNum := strings.Count(projectDirPath, string(os.PathSeparator)) + 1

	for _, repoEntry := range repoEntries {
		if !strings.HasPrefix(repoEntry.Path, projectDirPath) {
			continue
		}

		entryPathComponentsNum := strings.Count(repoEntry.Path, string(os.PathSeparator))
		if entryPathComponentsNum > expectedPathComponentNum {
			break
		}
		if entryPathComponentsNum != expectedPathComponentNum {
			continue
		}
		if repoEntry.IsDir {
			continue
		}
		if repoEntry.Name == "build.gradle.kts" || repoEntry.Name == "build.gradle" {
			return filepath.Dir(repoEntry.Path)
		}
	}

	return ""
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
