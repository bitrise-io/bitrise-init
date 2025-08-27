package kmp

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/detectors/direntry"
	"github.com/bitrise-io/bitrise-init/detectors/gradle"
	"github.com/bitrise-io/bitrise-init/detectors/kmp"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/scanners/java"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/log"
)

/*
Relevant Gradle dependencies:
	plugins:
		org.jetbrains.kotlin.multiplatform -> kotlin("multiplatform")
			This plugin is used to enable Kotlin Multiplatform projects, allowing you to share code between different platforms (e.g., JVM, JS, Native).
*/

const (
	projectType       = "kotlin-multiplatform"
	configName        = "kotlin-multiplatform-config"
	defaultConfigName = "default-kotlin-multiplatform-config"
	testWorkflowID    = "run_tests"

	gradleProjectRootDirInputEnvKey  = "PROJECT_ROOT_DIR"
	gradleProjectRootDirInputTitle   = "The root directory of the Gradle project."
	gradleProjectRootDirInputSummary = "The root directory of the Gradle project, which contains all source files from your project, as well as Gradle files, including the Gradle Wrapper (`gradlew`) file."
)

type Scanner struct {
	kmpProject *kmp.Project
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return projectType
}

func (s *Scanner) DetectPlatform(searchDir string) (bool, error) {
	log.TInfof("Searching for Gradle project files...")

	rootEntry, err := direntry.WalkDir(searchDir, 6)
	if err != nil {
		return false, err
	}

	gradleWrapperScripts := rootEntry.FindAllEntriesByName("gradlew", false)

	log.TDonef("%d Gradle wrapper script(s) found", len(gradleWrapperScripts))
	if len(gradleWrapperScripts) == 0 {
		return false, nil
	}
	gradleWrapperScript := gradleWrapperScripts[0]

	log.TInfof("Scanning project with Gradle wrapper script: %s", gradleWrapperScript.AbsPath)

	projectRootDir := gradleWrapperScript.Parent()
	if projectRootDir == nil {
		return false, fmt.Errorf("failed to get parent directory of %s", gradleWrapperScript.AbsPath)
	}
	gradleProject, err := gradle.ScanProject(*projectRootDir)
	if err != nil {
		return false, err
	}
	if gradleProject == nil {
		log.TWarnf("No Gradle project found in %s", projectRootDir.AbsPath)
		return false, nil
	}

	kmpProject, err := kmp.ScanProject(*gradleProject)
	if err != nil {
		return false, fmt.Errorf("failed to scan Kotlin Multiplatform project: %w", err)
	}

	printKMPProject(*kmpProject)

	s.kmpProject = kmpProject

	return kmpProject != nil, nil
}

func (s *Scanner) ExcludedScannerNames() []string {
	return []string{
		android.ScannerName,
		string(ios.XcodeProjectTypeIOS),
		java.ProjectType,
	}
}

func (s *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	gradleProjectRootDirOption := models.NewOption(gradleProjectRootDirInputTitle, gradleProjectRootDirInputSummary, gradleProjectRootDirInputEnvKey, models.TypeSelector)
	configOption := models.NewConfigOption(configName, nil)
	gradleProjectRootDirOption.AddConfig(s.kmpProject.GradleProject.RootDirEntry.RelPath, configOption)

	return *gradleProjectRootDirOption, nil, nil, nil
}

func (s *Scanner) DefaultOptions() models.OptionNode {
	gradleProjectRootDirOption := models.NewOption(gradleProjectRootDirInputTitle, gradleProjectRootDirInputSummary, gradleProjectRootDirInputEnvKey, models.TypeUserInput)
	configOption := models.NewConfigOption(defaultConfigName, nil)
	gradleProjectRootDirOption.AddConfig(models.UserInputOptionDefaultValue, configOption)
	return *gradleProjectRootDirOption
}

func (s *Scanner) Configs(sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	gradleProjectRootDir := "$" + gradleProjectRootDirInputEnvKey
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: sshKeyActivation})...,
	)
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.GradleUnitTestStepListItem(gradleProjectRootDir),
	)
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.DefaultDeployStepList()...,
	)

	config, err := configBuilder.Generate(projectType)
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

	gradleProjectRootDir := "$" + gradleProjectRootDirInputEnvKey
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: models.SSHKeyActivationConditional})...,
	)
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.GradleUnitTestStepListItem(gradleProjectRootDir),
	)
	configBuilder.AppendStepListItemsTo(testWorkflowID,
		steps.DefaultDeployStepList()...,
	)

	config, err := configBuilder.Generate(projectType)
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

func printKMPProject(kmpProject kmp.Project) {
	log.TPrintf("Project root dir: %s", kmpProject.GradleProject.RootDirEntry.RelPath)
	log.TPrintf("Gradle wrapper script: %s", kmpProject.GradleProject.GradlewFileEntry.RelPath)
	if kmpProject.GradleProject.ConfigDirEntry != nil {
		log.TPrintf("Gradle config dir: %s", kmpProject.GradleProject.ConfigDirEntry.RelPath)
	}
	if kmpProject.GradleProject.VersionCatalogFileEntry != nil {
		log.TPrintf("Version catalog file: %s", kmpProject.GradleProject.VersionCatalogFileEntry.RelPath)
	}
	if kmpProject.GradleProject.SettingsGradleFileEntry != nil {
		log.TPrintf("Gradle settings file: %s", kmpProject.GradleProject.SettingsGradleFileEntry.RelPath)
	}
	if len(kmpProject.GradleProject.IncludedProjects) > 0 {
		log.TPrintf("Included projects:")
		for _, includedProject := range kmpProject.GradleProject.IncludedProjects {
			log.TPrintf("- %s: %s", includedProject.Name, includedProject.BuildScriptFileEntry.RelPath)
		}
	}

	if kmpProject.IOSAppDetectResult != nil {
		log.TPrintf("iOS App target: %s", kmpProject.IOSAppDetectResult.Projects[0].RelPath)
	}
	if kmpProject.AndroidAppDetectResult != nil {
		log.TPrintf("Android App target: %s", kmpProject.AndroidAppDetectResult.Modules[0].BuildScriptPth)
	}
}
