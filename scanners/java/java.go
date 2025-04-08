package java

import (
	"github.com/bitrise-io/bitrise-init/detectors/gradle"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/log"
)

const (
	projectType       = "java"
	configName        = "java-config"
	defaultConfigName = "default-java-config"
	testWorkflowID    = "run_tests"
)

type Scanner struct {
	gradleProject gradle.Project
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return projectType
}

func printGradleProject(gradleProject gradle.Project) {
	log.TPrintf("Project root dir: %s", gradleProject.RootDirEntry.RelPath)
	log.TPrintf("Gradle wrapper script: %s", gradleProject.GradlewFileEntry.RelPath)
	if gradleProject.ConfigDirEntry != nil {
		log.TPrintf("Gradle config dir: %s", gradleProject.ConfigDirEntry.RelPath)
	}
	if gradleProject.VersionCatalogFileEntry != nil {
		log.TPrintf("Version catalog file: %s", gradleProject.VersionCatalogFileEntry.RelPath)
	}
	if gradleProject.SettingsGradleFileEntry != nil {
		log.TPrintf("Gradle settings file: %s", gradleProject.SettingsGradleFileEntry.RelPath)
	}
	if len(gradleProject.IncludedProjects) > 0 {
		log.TPrintf("Included projects:")
		for _, includedProject := range gradleProject.IncludedProjects {
			log.TPrintf("- %s: %s", includedProject.Name, includedProject.BuildScriptFileEntry.RelPath)
		}
	}
}

func (s *Scanner) DetectPlatform(searchDir string) (bool, error) {
	log.TInfof("Searching for Gradle project files...")

	gradleProject, err := gradle.ScanProject(searchDir)
	if err != nil {
		return false, err
	}

	log.TDonef("Gradle project found: %v", gradleProject != nil)
	if gradleProject != nil {
		printGradleProject(*gradleProject)
	} else {

	}

	// TODO: implement
	return false, nil
}

func (s *Scanner) ExcludedScannerNames() []string {
	// TODO: implement
	return []string{}
}

func (s *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	// TODO: implement
	return models.OptionNode{}, nil, nil, nil
}

func (s *Scanner) DefaultOptions() models.OptionNode {
	// TODO: implement
	return models.OptionNode{}
}

func (s *Scanner) Configs(sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	// TODO: implement
	return models.BitriseConfigMap{}, nil
}

func (s *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	// TODO: implement
	return models.BitriseConfigMap{}, nil
}
