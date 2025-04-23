package java

import (
	"github.com/bitrise-io/bitrise-init/detectors/gradle"
	"github.com/bitrise-io/bitrise-init/detectors/maven"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/v2/models"
	"github.com/bitrise-io/go-utils/log"
	"gopkg.in/yaml.v2"
)

const (
	ProjectType = "java"

	buildToolInputTitle   = "Build tool"
	buildToolInputSummary = "The build tool used in the project. Supported options: Gradle, Maven."
	buildToolGradle       = "Gradle"
	buildToolMaven        = "Maven"

	testWorkflowID = "run_tests"

	gradleConfigName        = "java-gradle-config"
	defaultGradleConfigName = "default-java-gradle-config"
	gradlewPathInputEnvKey  = "GRADLEW_PATH"
	gradlewPathInputTitle   = "The project's Gradle Wrapper script (gradlew) path."
	gradlewPathInputSummary = "The project's Gradle Wrapper script (gradlew) path."

	mavenConfigName                 = "java-maven-config"
	defaultMavenConfigName          = "default-java-maven-config"
	mavenProjectRootDirInputEnvKey  = "MAVEN_PROJECT_ROOT_DIR"
	mavenProjectRootDirInputTitle   = "The root directory of the Maven project (where the pom.xml file is located)."
	mavenProjectRootDirInputSummary = "The root directory of the Maven project (where the pom.xml file is located)."
	mavenTestScriptTitle            = `Run Maven tests`
	mavenTestScriptContent          = `#!/usr/bin/env bash
# fail if any commands fails
set -e
# make pipelines' return status equal the last command to exit with a non-zero status, or zero if all commands exit successfully
set -o pipefail
# debug log
set -x

mvm test
`
	installMavenScriptTitle   = `Install Maven`
	installMavenScriptContent = `#!/usr/bin/env bash
# fail if any commands fails
set -e
# make pipelines' return status equal the last command to exit with a non-zero status, or zero if all commands exit successfully
set -o pipefail
# debug log
set -x

sudo apt install maven
`
)

type Scanner struct {
	gradleProject *gradle.Project
	mavenProject  *maven.Project
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return ProjectType
}

func (s *Scanner) DetectPlatform(searchDir string) (bool, error) {
	log.TInfof("Searching for Gradle project files...")

	gradleProject, err := gradle.ScanProject(searchDir)
	if err != nil {
		return false, err
	}

	log.TDonef("Gradle project found: %v", gradleProject != nil)
	if gradleProject != nil {
		s.gradleProject = gradleProject
		printGradleProject(*gradleProject)
		return true, nil
	}

	log.TInfof("Searching for Maven project files...")
	mavenProject, err := maven.ScanProject(searchDir)
	if err != nil {
		return false, err
	}

	log.TDonef("Maven project found: %v", mavenProject != nil)
	if mavenProject != nil {
		s.mavenProject = mavenProject
		printMavenProject(*mavenProject)
		return true, nil
	}

	return false, nil
}

func (s *Scanner) ExcludedScannerNames() []string {
	return []string{}
}

func (s *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	if s.gradleProject != nil {
		gradlewPathOption := models.NewOption(gradlewPathInputTitle, gradlewPathInputSummary, gradlewPathInputEnvKey, models.TypeSelector)
		configOption := models.NewConfigOption(gradleConfigName, nil)
		gradlewPathOption.AddConfig(s.gradleProject.GradlewFileEntry.RelPath, configOption)
		return *gradlewPathOption, nil, nil, nil
	}

	if s.mavenProject != nil {
		mavenProjectRootDirOption := models.NewOption(mavenProjectRootDirInputTitle, mavenProjectRootDirInputSummary, mavenProjectRootDirInputEnvKey, models.TypeSelector)
		configOption := models.NewConfigOption(mavenConfigName, nil)
		mavenProjectRootDirOption.AddConfig(s.mavenProject.RootDirEntry.RelPath, configOption)
		return *mavenProjectRootDirOption, nil, nil, nil
	}

	return models.OptionNode{}, nil, nil, nil
}

func (s *Scanner) DefaultOptions() models.OptionNode {
	buildToolOption := models.NewOption(buildToolInputTitle, buildToolInputSummary, "", models.TypeSelector)

	gradlewPathOption := models.NewOption(gradlewPathInputTitle, gradlewPathInputSummary, gradlewPathInputEnvKey, models.TypeUserInput)
	buildToolOption.AddOption(buildToolGradle, gradlewPathOption)

	gradleConfigOption := models.NewConfigOption(defaultGradleConfigName, nil)
	gradlewPathOption.AddConfig("", gradleConfigOption)

	mavenProjectRootDirOption := models.NewOption(mavenProjectRootDirInputTitle, mavenProjectRootDirInputSummary, mavenProjectRootDirInputEnvKey, models.TypeUserInput)
	buildToolOption.AddOption(buildToolMaven, mavenProjectRootDirOption)

	mavenConfigOption := models.NewConfigOption(defaultMavenConfigName, nil)
	mavenProjectRootDirOption.AddConfig("", mavenConfigOption)

	return *buildToolOption
}

func (s *Scanner) Configs(sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	bitriseDataMap := models.BitriseConfigMap{}

	if s.gradleProject != nil {
		gradlewPath := "$" + gradlewPathInputEnvKey
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: sshKeyActivation})...,
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.GradleUnitTestStepListItem(gradlewPath),
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultDeployStepList()...,
		)

		config, err := configBuilder.Generate(ProjectType)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		bitriseDataMap[gradleConfigName] = string(data)
	}

	if s.mavenProject != nil {
		mavenProjectRootDir := "$" + mavenProjectRootDirInputEnvKey
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: sshKeyActivation})...,
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.ScriptStepListItem(installMavenScriptTitle, installMavenScriptContent),
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.ScriptStepListItem(mavenTestScriptTitle, mavenTestScriptContent, envmanModels.EnvironmentItemModel{
				"working_dir": mavenProjectRootDir,
			}),
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultDeployStepList()...,
		)
		config, err := configBuilder.Generate(ProjectType)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}
		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}
		bitriseDataMap[mavenConfigName] = string(data)
	}

	return bitriseDataMap, nil
}

func (s *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	bitriseDataMap := models.BitriseConfigMap{}

	if s.gradleProject != nil {
		gradlewPath := "$" + gradlewPathInputEnvKey
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: models.SSHKeyActivationConditional})...,
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.GradleUnitTestStepListItem(gradlewPath),
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultDeployStepList()...,
		)

		config, err := configBuilder.Generate(ProjectType)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		bitriseDataMap[defaultGradleConfigName] = string(data)
	}

	if s.mavenProject != nil {
		mavenProjectRootDir := "$" + mavenProjectRootDirInputEnvKey
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: models.SSHKeyActivationConditional})...,
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.ScriptStepListItem(installMavenScriptTitle, installMavenScriptContent),
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.ScriptStepListItem(mavenTestScriptTitle, mavenTestScriptContent, envmanModels.EnvironmentItemModel{
				"working_dir": mavenProjectRootDir,
			}),
		)
		configBuilder.AppendStepListItemsTo(testWorkflowID,
			steps.DefaultDeployStepList()...,
		)
		config, err := configBuilder.Generate(ProjectType)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}
		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}
		bitriseDataMap[defaultMavenConfigName] = string(data)
	}

	return bitriseDataMap, nil
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

func printMavenProject(mavenProject maven.Project) {
	log.TPrintf("Project root dir: %s", mavenProject.RootDirEntry.RelPath)
	log.TPrintf("Maven POM file: %s", mavenProject.ProjectObjectModelFileEntry.RelPath)
}
