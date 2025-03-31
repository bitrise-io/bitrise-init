package kmp

import (
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners"
	"github.com/bitrise-io/go-utils/log"
)

/*
	Relevant Gradle dependencies:
		plugins:
			org.jetbrains.kotlin.multiplatform -> kotlin("multiplatform")
				This plugin is used to enable Kotlin Multiplatform projects, allowing you to share code between different platforms (e.g., JVM, JS, Native).
			org.jetbrains.kotlin.plugin.compose -> kotlin("plugin.compose")
				This plugin is used to add support for Jetpack Compose in Kotlin Multiplatform projects. It allows you to use Compose UI components across multiple platforms.
*/

type ProjectStructure struct {
	GradleConfigurationDirPath string
	UsesVersionCatalogFile     bool
	Projects                   []string
	ProjectDirPaths            []string
}

const scannerName = "kmp"

type Scanner struct {
}

func NewScanner() scanners.ScannerInterface {
	return &Scanner{}
}

func (s Scanner) Name() string {
	return scannerName
}

func (s Scanner) DetectPlatform(searchDir string) (bool, error) {
	projectStructure, err := s.detectProjectStructure(searchDir)
	if err != nil {
		return false, err
	}
	if projectStructure == nil {
		return false, nil
	}

	printProjectStructure(*projectStructure)

	return false, nil
}

func (s Scanner) ExcludedScannerNames() []string {
	//TODO implement me
	return nil
}

func (s Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	//TODO implement me
	return models.OptionNode{}, nil, nil, nil
}

func (s Scanner) DefaultOptions() models.OptionNode {
	//TODO implement me
	return models.OptionNode{}
}

func (s Scanner) Configs(sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	//TODO implement me
	return nil, nil
}

func (s Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	//TODO implement me
	return nil, nil
}

func (s Scanner) detectProjectStructure(searchDir string) (*ProjectStructure, error) {
	repoEntries, err := listDirEntries(searchDir, 4)
	if err != nil {
		return nil, err
	}

	// KMP project detection
	gradleConfigurationDirectories, err := detectGradleConfigurationDirectories(repoEntries)
	if err != nil {
		return nil, err
	}
	if len(gradleConfigurationDirectories) == 0 {
		return nil, nil
	}

	gradleConfigurationDirectory := gradleConfigurationDirectories[0]
	if len(gradleConfigurationDirectories) > 1 {
		log.Warnf("Multiple gradle configuration directories found: %v, using the first one: %s", gradleConfigurationDirectories, gradleConfigurationDirectory)
	}

	projectTypeDetected := false
	versionCatalogFile := detectVersionCatalogFile(gradleConfigurationDirectory.Path, repoEntries)
	usesVersionCatalogFile := versionCatalogFile != nil
	if usesVersionCatalogFile {
		detected, err := detectAnyDependencies(*versionCatalogFile, []string{
			"org.jetbrains.kotlin.multiplatform",
			"org.jetbrains.kotlin.plugin.compose",
		})
		if err != nil {
			return nil, err
		}
		projectTypeDetected = detected
	}

	gradleProjectRootDirPth := filepath.Dir(gradleConfigurationDirectory.Path)
	if !projectTypeDetected {
		projectGradleBuildScriptFiles := detectGradleBuildScriptFiles(gradleProjectRootDirPth, repoEntries)
		if !usesVersionCatalogFile {
			if len(projectGradleBuildScriptFiles) == 0 {
				return nil, nil
			}

			for _, projectGradleBuildScriptFile := range projectGradleBuildScriptFiles {
				detected, err := detectAnyDependencies(projectGradleBuildScriptFile, []string{
					"org.jetbrains.kotlin.multiplatform",
					"org.jetbrains.kotlin.plugin.compose",
					`kotlin("multiplatform")`,
					`kotlin("plugin.compose")`,
				})
				if err != nil {
					return nil, err
				}
				if !detected {
					return nil, nil
				}
			}
		}
	}
	// ---

	// Included projects detection
	settingsGradleFile := detectSettingsGradleFile(gradleProjectRootDirPth, repoEntries)
	if settingsGradleFile == nil {
		return nil, nil
	}

	projectIncludes, err := detectProjectIncludes(*settingsGradleFile)
	if err != nil {
		return nil, err
	}

	var projects []string
	var projectDirPaths []string
	for _, projectInclude := range projectIncludes {
		projectDirPath := detectProjectDirPath(gradleProjectRootDirPth, projectInclude, repoEntries)
		if projectDirPath != "" {
			projects = append(projects, projectInclude)
			projectDirPaths = append(projectDirPaths, projectDirPath)
		}
	}

	projectStructure := ProjectStructure{
		GradleConfigurationDirPath: gradleConfigurationDirectory.Path,
		UsesVersionCatalogFile:     usesVersionCatalogFile,
		Projects:                   projects,
		ProjectDirPaths:            projectDirPaths,
	}

	return &projectStructure, nil
}
