package kmp

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android"
	"github.com/bitrise-io/bitrise-init/scanners/gradle"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/steps"
)

/*
Relevant Gradle dependencies:
	plugins:
		org.jetbrains.kotlin.multiplatform -> kotlin("multiplatform")
			This plugin is used to enable Kotlin Multiplatform projects, allowing you to share code between different platforms (e.g., JVM, JS, Native).
		org.jetbrains.kotlin.plugin.compose -> kotlin("plugin.compose")
			This plugin is used to add support for Jetpack Compose in Kotlin Multiplatform projects. It allows you to use Compose UI components across multiple platforms.
*/

const (
	scannerName       = "kmp"
	configName        = "kotlin-multiplatform-config"
	defaultConfigName = "default-kotlin-multiplatform-config"
	testWorkflowID    = "run_tests"
	gradleTestTask    = "test"

	gradlewPathInputEnvKey  = "GRADLEW_PATH"
	gradlewPathInputTitle   = "The project's Gradle Wrapper script (gradlew) path."
	gradlewPathInputSummary = "The project's Gradle Wrapper script (gradlew) path."
)

type Scanner struct {
	gradleProject gradle.Project
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return scannerName
}

func (s *Scanner) DetectPlatform(searchDir string) (bool, error) {
	gradleProject, err := gradle.ScanProject(searchDir)
	if err != nil {
		return false, err
	}
	if gradleProject == nil {
		return false, nil
	}

	kotlinMultiplatformDetected, err := gradleProject.DetectAnyDependencies([]string{
		"org.jetbrains.kotlin.multiplatform",
		"org.jetbrains.kotlin.plugin.compose",
		`kotlin("multiplatform")`,
		`kotlin("plugin.compose")`,
	})
	if err != nil {
		return false, err
	}

	s.gradleProject = *gradleProject

	fmt.Println(gradleProject.ToJSON())

	return kotlinMultiplatformDetected, nil
}

func (s *Scanner) ExcludedScannerNames() []string {
	return []string{android.ScannerName, string(ios.XcodeProjectTypeIOS)}
}

func (s *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	gradlewPathOption := models.NewOption(gradlewPathInputTitle, gradlewPathInputSummary, gradlewPathInputEnvKey, models.TypeSelector)
	configOption := models.NewConfigOption(configName, nil)
	gradlewPathOption.AddConfig(s.gradleProject.GradlewPath, configOption)
	return *gradlewPathOption, nil, nil, nil
}

func (s *Scanner) DefaultOptions() models.OptionNode {
	gradlewPathOption := models.NewOption(gradlewPathInputTitle, gradlewPathInputSummary, gradlewPathInputEnvKey, models.TypeUserInput)
	configOption := models.NewConfigOption(defaultConfigName, nil)
	gradlewPathOption.AddConfig(models.UserInputOptionDefaultValue, configOption)
	return *gradlewPathOption
}

func (s *Scanner) Configs(sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	gradlewPath := "$" + gradlewPathInputEnvKey

	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: sshKeyActivation})...,
	)
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.GradleRunnerStepListItem(gradlewPath, gradleTestTask),
	)

	config, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	bitriseDataMap := models.BitriseConfigMap{}
	bitriseDataMap[configName] = string(data)

	return bitriseDataMap, nil
}

func (s *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	gradlewPath := "$" + gradlewPathInputEnvKey

	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: models.SSHKeyActivationConditional})...,
	)
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.GradleRunnerStepListItem(gradlewPath, gradleTestTask),
	)

	config, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	bitriseDataMap := models.BitriseConfigMap{}
	bitriseDataMap[defaultConfigName] = string(data)

	return bitriseDataMap, nil
}
