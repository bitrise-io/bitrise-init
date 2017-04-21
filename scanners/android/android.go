package android

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
)

const scannerName = "android"

const (
	// GradleFileInputKey ...
	GradleFileInputKey = "gradle_file"
	// GradleFileInputEnvKey ...
	GradleFileInputEnvKey = "GRADLE_BUILD_FILE_PATH"
	// GradleFileInputTitle ...
	GradleFileInputTitle = "Path to the gradle file to use"
)

const (
	pathInputKey              = "path"
	projectRootDirInputEnvKey = "PROJECT_ROOT_DIR"
	projectRootDirInputTitle  = "Path to the Android project root directory"
)

const (
	// GradleTaskInputKey ...
	GradleTaskInputKey = "gradle_task"
	// GradleTaskInputEnvKey ...
	GradleTaskInputEnvKey = "GRADLE_TASK"
	// GradleTaskInputTitle ...
	GradleTaskInputTitle = "Gradle task to run"
)

const (
	// GradlewPathInputKey ...
	GradlewPathInputKey = "gradlew_path"
	// GradlewPathInputEnvKey ...
	GradlewPathInputEnvKey = "GRADLEW_PATH"
	// GradlewPathInputTitle ...
	GradlewPathInputTitle = "Gradlew file path"
)

var defaultGradleTasks = []string{
	"assemble",
	"assembleDebug",
	"assembleRelease",
}

const (
	// ConfigName ...
	ConfigName = "android-config"
	// DefaultConfigName ...
	DefaultConfigName = "default-android-config"
)

// ConfigDescriptor ...
type ConfigDescriptor struct {
	MissingGradlew bool
}

// NewConfigDescriptor ...
func NewConfigDescriptor(missingGradlew bool) ConfigDescriptor {
	return ConfigDescriptor{
		MissingGradlew: missingGradlew,
	}
}

// ConfigName ...
func (descriptor ConfigDescriptor) ConfigName() string {
	name := "android"
	if descriptor.MissingGradlew {
		name += "-missing-gradlew"
	}
	name += "-config"
	return name
}

// Scanner ...
type Scanner struct {
	FileList    []string
	GradleFiles []string

	configDescriptors []ConfigDescriptor
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (scanner Scanner) Name() string {
	return scannerName
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return false, fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}
	scanner.FileList = fileList

	// Search for gradle file
	log.Infoft("Searching for build.gradle files")

	gradleFiles, err := utility.FilterRootBuildGradleFiles(fileList)
	if err != nil {
		return false, fmt.Errorf("failed to search for build.gradle files, error: %s", err)
	}
	scanner.GradleFiles = gradleFiles

	log.Printft("%d build.gradle files detected", len(gradleFiles))
	for _, file := range gradleFiles {
		log.Printft("- %s", file)
	}

	if len(gradleFiles) == 0 {
		log.Printft("platform not detected")
		return false, nil
	}

	log.Doneft("Platform detected")

	return true, nil
}

// GenerateOption ...
func (scanner *Scanner) GenerateOption(allowMissingGradlew, addProjectRootDirOption bool) (models.OptionModel, []ConfigDescriptor, models.Warnings, error) {
	// Search for gradlew_path input
	log.Infoft("Searching for gradlew files")

	warnings := models.Warnings{}
	gradlewFiles, err := utility.FilterGradlewFiles(scanner.FileList)
	if err != nil {
		return models.OptionModel{}, []ConfigDescriptor{}, warnings, fmt.Errorf("Failed to list gradlew files, error: %s", err)
	}

	log.Printft("%d gradlew files detected", len(gradlewFiles))
	for _, file := range gradlewFiles {
		log.Printft("- %s", file)
	}

	rootGradlewPath := ""
	gradlewFilesCount := len(gradlewFiles)
	switch {
	case gradlewFilesCount == 0:
		if allowMissingGradlew {
			log.Warnft("No gradle wrapper (gradlew) found")
			log.Warnft(`<b>No Gradle Wrapper (gradlew) found.</b> 
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`)
		} else {
			log.Errorft("No gradle wrapper (gradlew) found")
			return models.OptionModel{}, []ConfigDescriptor{}, warnings, fmt.Errorf(`<b>No Gradle Wrapper (gradlew) found.</b> 
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`)
		}
	case gradlewFilesCount == 1:
		rootGradlewPath = gradlewFiles[0]
	case gradlewFilesCount > 1:
		rootGradlewPath = gradlewFiles[0]
		log.Warnft("Multiple gradlew file, detected:")
		for _, gradlewPth := range gradlewFiles {
			log.Warnft("- %s", gradlewPth)
		}
		log.Warnft("Using: %s", rootGradlewPath)
	}

	// Inspect Gradle files

	configDescriptors := []ConfigDescriptor{}
	gradleFileOption := models.NewOption(GradleFileInputTitle, GradleFileInputEnvKey)

	for _, gradleFile := range scanner.GradleFiles {
		log.Infoft("Inspecting gradle file: %s", gradleFile)

		// generate-gradle-wrapper step will generate the wrapper
		if rootGradlewPath == "" {
			gradleFileDir := filepath.Dir(gradleFile)
			rootGradlewPath = filepath.Join(gradleFileDir, "gradlew")

			configDescriptors = append(configDescriptors, NewConfigDescriptor(true))
		} else {
			configDescriptors = append(configDescriptors, NewConfigDescriptor(false))
		}
		// ---

		gradlewPthOption := models.NewOption(GradlewPathInputTitle, GradlewPathInputEnvKey)

		if addProjectRootDirOption {
			projectRootDirOption := models.NewOption(projectRootDirInputTitle, projectRootDirInputEnvKey)
			gradleFileOption.AddOption(gradleFile, projectRootDirOption)

			projectRootDirOption.AddOption(filepath.Join("$BITRISE_SOURCE_DIR", filepath.Dir(gradleFile)), gradlewPthOption)
		} else {
			gradleFileOption.AddOption(gradleFile, gradlewPthOption)
		}

		gradleTaskOption := models.NewOption(GradleTaskInputTitle, GradleTaskInputEnvKey)
		gradlewPthOption.AddOption(rootGradlewPath, gradleTaskOption)

		log.Printft("%d gradle tasks", len(defaultGradleTasks))

		for _, gradleTask := range defaultGradleTasks {
			log.Printft("- %s", gradleTask)

			configOption := models.NewConfigOption(ConfigName)
			gradleTaskOption.AddConfig(gradleTask, configOption)
		}
	}

	configDescriptors = plain(configDescriptors)

	return *gradleFileOption, configDescriptors, warnings, nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	option, configDescriptors, warnings, err := scanner.GenerateOption(false, true)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}
	scanner.configDescriptors = configDescriptors
	return option, warnings, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	gradleFileOption := models.NewOption(GradleFileInputTitle, GradleFileInputEnvKey)

	projectRootOption := models.NewOption(projectRootDirInputTitle, projectRootDirInputEnvKey)
	gradleFileOption.AddOption("_", projectRootOption)

	gradlewPthOption := models.NewOption(GradlewPathInputTitle, GradlewPathInputEnvKey)
	projectRootOption.AddOption("_", gradlewPthOption)

	gradleTaskOption := models.NewOption(GradleTaskInputTitle, GradleTaskInputEnvKey)
	gradlewPthOption.AddOption("_", gradleTaskOption)

	configOption := models.NewConfigOption(DefaultConfigName)
	gradleTaskOption.AddConfig("_", configOption)

	return *gradleFileOption
}

// GenerateConfigBuilder ...
func GenerateConfigBuilder(missingGradlew bool) models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	configBuilder.AppendPreparStepList(steps.ChangeWorkDirStepListItem(envmanModels.EnvironmentItemModel{pathInputKey: "$" + projectRootDirInputEnvKey}))
	if missingGradlew {
		configBuilder.AppendPreparStepList(steps.GenerateGradleWrapperStepListItem())
	}
	configBuilder.AppendPreparStepList(steps.InstallMissingAndroidToolsStepListItem())

	configBuilder.AppendMainStepList(steps.GradleRunnerStepListItem(
		envmanModels.EnvironmentItemModel{GradleFileInputKey: "$" + GradleFileInputEnvKey},
		envmanModels.EnvironmentItemModel{GradleTaskInputKey: "$" + GradleTaskInputEnvKey},
		envmanModels.EnvironmentItemModel{GradlewPathInputKey: "$" + GradlewPathInputEnvKey},
	))

	return *configBuilder
}

func plain(configDescriptors []ConfigDescriptor) []ConfigDescriptor {
	descriptors := []ConfigDescriptor{}
	descritorNameMap := map[string]bool{}
	for _, descriptor := range configDescriptors {
		_, exist := descritorNameMap[descriptor.ConfigName()]
		if !exist {
			descriptors = append(descriptors, descriptor)
		}
	}
	return descriptors
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	bitriseDataMap := models.BitriseConfigMap{}
	for _, descriptor := range scanner.configDescriptors {
		configBuilder := GenerateConfigBuilder(descriptor.MissingGradlew)

		config, err := configBuilder.Generate()
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		bitriseDataMap[descriptor.ConfigName()] = string(data)
	}

	return bitriseDataMap, nil
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	configBuilder.AppendPreparStepList(steps.ChangeWorkDirStepListItem(envmanModels.EnvironmentItemModel{pathInputKey: "$" + projectRootDirInputEnvKey}))
	configBuilder.AppendPreparStepList(steps.InstallMissingAndroidToolsStepListItem())

	configBuilder.AppendMainStepList(steps.GradleRunnerStepListItem(
		envmanModels.EnvironmentItemModel{GradleFileInputKey: "$" + GradleFileInputEnvKey},
		envmanModels.EnvironmentItemModel{GradleTaskInputKey: "$" + GradleTaskInputEnvKey},
		envmanModels.EnvironmentItemModel{GradlewPathInputKey: "$" + GradlewPathInputEnvKey},
	))

	config, err := configBuilder.Generate()
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		DefaultConfigName: string(data),
	}, nil
}
