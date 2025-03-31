package kmp

import (
	"github.com/bitrise-io/bitrise-init/scanners/gradle"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners"
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
	gradleProject, err := gradle.ScanProject(searchDir)
	if err != nil {
		return false, err
	}
	if gradleProject == nil {
		return false, nil
	}

	kotlinMultiplatformDetected := false
	if gradleProject.VersionCatalogFilePath != "" {
		detected, err := gradle.DetectAnyDependencies(gradleProject.VersionCatalogFilePath, []string{
			"org.jetbrains.kotlin.multiplatform",
			"org.jetbrains.kotlin.plugin.compose",
		})
		if err != nil {
			return false, err
		}
		kotlinMultiplatformDetected = detected
	}

	if !kotlinMultiplatformDetected {
		for _, buildScriptPath := range gradleProject.BuildScriptPaths {
			detected, err := gradle.DetectAnyDependencies(buildScriptPath, []string{
				"org.jetbrains.kotlin.multiplatform",
				"org.jetbrains.kotlin.plugin.compose",
				`kotlin("multiplatform")`,
				`kotlin("plugin.compose")`,
			})
			if err != nil {
				return false, err
			}
			if detected {
				kotlinMultiplatformDetected = true
				break
			}
		}
	}

	gradle.PrintProject(*gradleProject)

	return kotlinMultiplatformDetected, nil
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
