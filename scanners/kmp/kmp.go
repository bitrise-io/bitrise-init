package kmp

import (
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners"
	"github.com/bitrise-io/go-utils/log"
)

type ProjectStructure struct {
	GradleConfigurationDirPath string
	UsesVersionCatalogFile     bool
	ProjectDirPaths            []string
	ComposeAppDirPath          string
	ComposeAppProjectsDirPaths []string
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
	repoEntries, err := listDirEntries(searchDir, 4)
	if err != nil {
		return false, err
	}

	gradleConfigurationDirectories, err := detectGradleConfigurationDirectories(repoEntries)
	if err != nil {
		return false, err
	}
	if len(gradleConfigurationDirectories) == 0 {
		return false, nil
	}

	gradleConfigurationDirectory := gradleConfigurationDirectories[0]
	if len(gradleConfigurationDirectories) > 1 {
		log.Warnf("Multiple gradle configuration directories found: %v, using the first one: %s", gradleConfigurationDirectories, gradleConfigurationDirectory)
	}

	versionCatalogFile := detectVersionCatalogFile(gradleConfigurationDirectory.Path, repoEntries)
	usesVersionCatalogFile := versionCatalogFile != nil
	if usesVersionCatalogFile {
		detected, err := detectAnyDependencies(*versionCatalogFile, []string{"org.jetbrains.kotlin.multiplatform"})
		if err != nil {
			return false, err
		}
		if !detected {
			return false, nil
		}
	}

	gradleProjectRootDirPth := filepath.Dir(gradleConfigurationDirectory.Path)
	projectGradleBuildScriptFiles := detectProjectGradleBuildScriptFiles(gradleProjectRootDirPth, repoEntries)
	if !usesVersionCatalogFile {
		if len(projectGradleBuildScriptFiles) == 0 {
			return false, nil
		}

		for _, projectGradleBuildScriptFile := range projectGradleBuildScriptFiles {
			detected, err := detectAnyDependencies(projectGradleBuildScriptFile, []string{"org.jetbrains.kotlin.multiplatform"})
			if err != nil {
				return false, err
			}
			if !detected {
				return false, nil
			}
		}
	}

	var projectDirPaths []string
	var composeAppDirPath string
	for _, projectGradleBuildScriptFile := range projectGradleBuildScriptFiles {
		projectDirPath := filepath.Dir(projectGradleBuildScriptFile.Path)
		projectDirPaths = append(projectDirPaths, projectDirPath)
		if filepath.Base(projectDirPath) == "composeApp" {
			composeAppDirPath = projectDirPath
		}
	}

	var composeAppProjectDirPaths []string
	composeAppProjectDirectories := detectComposeAppProjectDirectories(composeAppDirPath, repoEntries)
	for _, composeAppProjectDirectory := range composeAppProjectDirectories {
		composeAppProjectDirPaths = append(composeAppProjectDirPaths, composeAppProjectDirectory.Path)
	}

	projectStructure := ProjectStructure{
		UsesVersionCatalogFile:     usesVersionCatalogFile,
		GradleConfigurationDirPath: gradleConfigurationDirectory.Path,
		ProjectDirPaths:            projectDirPaths,
		ComposeAppDirPath:          composeAppDirPath,
		ComposeAppProjectsDirPaths: composeAppProjectDirPaths,
	}

	printProjectStructure(projectStructure)

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
