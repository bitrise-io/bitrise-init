package gradle

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bitrise-io/bitrise-init/direntry"
)

/*
Settings File (settings.gradle[.kts]): The settings file is the entry point of every Gradle project.
	The primary purpose of the settings file is to add subprojects to your build.
	Gradle supports single and multi-project builds.
	- For single-project builds, the settings file is optional.
	- For multi-project builds, the settings file is mandatory and declares all subprojects.

Gradle wrapper scripts (gradlew and gradlew.bat)
	The presence of the gradlew and gradlew.bat files in the root directory of a project is a clear indicator that Gradle is used.
*/

type Project struct {
	RootDirPath            string
	GradlewPath            string
	ConfigDirPath          string
	VersionCatalogFilePath string
	SettingsGradleFilePath string

	IncludedProjects                []string
	IncludedProjectBuildScriptPaths []string
	BuildScriptPaths                []string
}

func ScanProject(searchDir string) (*Project, error) {
	rootEntry, err := direntry.ListEntries(searchDir, 4)
	if err != nil {
		return nil, err
	}

	projectRoot, err := detectGradleProjectRoot(*rootEntry)
	if err != nil {
		return nil, err
	}
	if projectRoot == nil {
		return nil, nil
	}
	projects, err := detectIncludedProjects(*projectRoot)
	if err != nil {
		return nil, err
	}

	var includedProjectBuildScriptPaths []string
	for _, projectBuildScript := range projects.projectBuildScripts {
		includedProjectBuildScriptPaths = append(includedProjectBuildScriptPaths, projectBuildScript.Path)
	}

	var buildScriptPaths []string
	for _, buildScript := range projects.buildScripts {
		buildScriptPaths = append(buildScriptPaths, buildScript.Path)
	}

	project := Project{
		RootDirPath:            projectRoot.rootDirEntry.Path,
		GradlewPath:            projectRoot.gradlewFileEntry.Path,
		ConfigDirPath:          projectRoot.configDirEntry.Path,
		VersionCatalogFilePath: projectRoot.versionCatalogFileEntry.Path,
		SettingsGradleFilePath: projectRoot.settingsGradleFileEntry.Path,

		IncludedProjects:                projects.projects,
		IncludedProjectBuildScriptPaths: includedProjectBuildScriptPaths,
		BuildScriptPaths:                buildScriptPaths,
	}

	return &project, nil
}

func (proj Project) DetectAnyDependencies(dependencies []string) (bool, error) {
	detected, err := proj.detectAnyDependenciesInVersionCatalogFile(dependencies)
	if err != nil {
		return false, err
	}
	if detected {
		return true, nil
	}

	detected, err = proj.detectAnyDependenciesInIncludedProjectBuildScripts(dependencies)
	if err != nil {
		return false, err
	}
	if detected {
		return true, nil
	}

	return proj.detectAnyDependenciesInBuildScripts(dependencies)
}

func (proj Project) detectAnyDependenciesInVersionCatalogFile(dependencies []string) (bool, error) {
	if proj.VersionCatalogFilePath == "" {
		return false, nil
	}
	return proj.detectAnyDependencies(proj.VersionCatalogFilePath, dependencies)
}

func (proj Project) detectAnyDependenciesInIncludedProjectBuildScripts(dependencies []string) (bool, error) {
	for _, includedProjectBuildScriptPth := range proj.IncludedProjectBuildScriptPaths {
		detected, err := proj.detectAnyDependencies(includedProjectBuildScriptPth, dependencies)
		if err != nil {
			return false, err
		}
		if detected {
			return true, nil
		}
	}
	return false, nil
}

func (proj Project) detectAnyDependenciesInBuildScripts(dependencies []string) (bool, error) {
	for _, buildScriptPth := range proj.BuildScriptPaths {
		detected, err := proj.detectAnyDependencies(buildScriptPth, dependencies)
		if err != nil {
			return false, err
		}
		if detected {
			return true, nil
		}
	}
	return false, nil
}

func (proj Project) detectAnyDependencies(pth string, dependencies []string) (bool, error) {
	file, err := os.Open(pth)
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

type gradleProjectRootEntry struct {
	rootDirEntry            direntry.DirEntry
	gradlewFileEntry        direntry.DirEntry
	configDirEntry          *direntry.DirEntry
	versionCatalogFileEntry *direntry.DirEntry
	settingsGradleFileEntry *direntry.DirEntry
}

func detectGradleProjectRoot(rootEntry direntry.DirEntry) (*gradleProjectRootEntry, error) {
	gradlewFileEntry := rootEntry.FindEntryByName("gradlew", false)
	if gradlewFileEntry == nil {
		if len(rootEntry.Entries) == 0 {
			return nil, nil
		}

		for _, entry := range rootEntry.Entries {
			if entry.IsDir {
				return detectGradleProjectRoot(entry)
			}
		}
		return nil, nil
	}

	projectRoot := gradleProjectRootEntry{
		rootDirEntry:     rootEntry,
		gradlewFileEntry: *gradlewFileEntry,
	}

	configDirEntry := rootEntry.FindEntryByName("gradle", true)
	if configDirEntry != nil {
		projectRoot.configDirEntry = configDirEntry

		versionCatalogFileEntry := configDirEntry.FindEntryByName("libs.versions.toml", false)
		if versionCatalogFileEntry != nil {
			projectRoot.versionCatalogFileEntry = versionCatalogFileEntry
		}
	}

	settingsFileEntry := rootEntry.FindEntryByName("settings.gradle", false)
	if settingsFileEntry == nil {
		settingsFileEntry = rootEntry.FindEntryByName("settings.gradle.kts", false)
	}
	if settingsFileEntry != nil {
		projectRoot.settingsGradleFileEntry = settingsFileEntry
	}

	return &projectRoot, nil
}

type includedProjects struct {
	buildScripts        []direntry.DirEntry
	projects            []string
	projectBuildScripts []direntry.DirEntry
}

func detectIncludedProjects(projectRootEntry gradleProjectRootEntry) (*includedProjects, error) {
	projects := includedProjects{}
	projects.buildScripts = projectRootEntry.rootDirEntry.FindAllEntriesByName("build.gradle", false)
	projects.buildScripts = append(projects.buildScripts, projectRootEntry.rootDirEntry.FindAllEntriesByName("build.gradle.kts", false)...)

	if projectRootEntry.settingsGradleFileEntry != nil {
		includes, err := detectProjectIncludes(*projectRootEntry.settingsGradleFileEntry)
		if err != nil {
			return nil, err
		}
		projects.projects = includes

		for _, include := range includes {
			var components []string

			include = strings.TrimPrefix(include, ":")
			includeComponents := strings.Split(include, ":")
			for _, includeComponent := range includeComponents {
				if includeComponent == "" {
					continue
				}
				includeComponent = strings.TrimSpace(includeComponent)
				components = append(components, includeComponent)
			}

			projectBuildScript := projectRootEntry.rootDirEntry.FindEntryByPath(false, append(components, "build.gradle")...)
			if projectBuildScript != nil {
				projects.projectBuildScripts = append(projects.projectBuildScripts, *projectBuildScript)
			}

			projectBuildScript = projectRootEntry.rootDirEntry.FindEntryByPath(false, append(components, "build.gradle.kts")...)
			if projectBuildScript != nil {
				projects.projectBuildScripts = append(projects.projectBuildScripts, *projectBuildScript)
			}
		}
	}

	return &projects, nil
}

func detectProjectIncludes(settingGradleFile direntry.DirEntry) ([]string, error) {
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
		if !strings.HasPrefix(line, "include(") && !strings.HasPrefix(line, "include ") {
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

func PrintProject(proj Project) {
	content, err := json.MarshalIndent(proj, "", "  ")
	if err == nil {
		log.Println(string(content))
	}
}
