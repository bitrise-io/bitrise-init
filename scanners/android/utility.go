package android

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	// ScannerName ...
	ScannerName = "android"
	// ConfigName ...
	ConfigName = "android-config"
	// DefaultConfigName ...
	DefaultConfigName = "default-android-config"

	// Step Inputs
	gradlewPathInputKey    = "gradlew_path"
	gradlewPathInputEnvKey = "GRADLEW_PATH"
	gradlewPathInputTitle  = "Gradlew file path"

	gradleFileInputKey    = "gradle_file"
	gradleFileInputEnvKey = "GRADLE_BUILD_FILE_PATH"
	gradleFileInputTitle  = "Path to the gradle file to use"

	gradleTaskInputKey = "gradle_task"
)

// CollectRootBuildGradleFiles - Collects the most root (mint path depth) build.gradle files
// May the searchDir contains multiple android projects, this case it return multiple builde.gradle path
// searchDir/android-project1/build.gradle, searchDir/android-project2/build.gradle, ...
func CollectRootBuildGradleFiles(searchDir string) ([]string, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, false)
	if err != nil {
		return nil, fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}

	return utility.FilterRootBuildGradleFiles(fileList)
}

// CheckLocalProperties - Returns warning if local.properties exists
// Local properties may contains absolute paths (sdk.dir=/Users/xyz/Library/Android/sdk),
// it should be gitignored
func CheckLocalProperties(buildGradleFile string) string {
	projectDir := filepath.Dir(buildGradleFile)
	localPropertiesPth := filepath.Join(projectDir, "local.properties")
	exist, err := pathutil.IsPathExists(localPropertiesPth)
	if err == nil && exist {
		return fmt.Sprintf(`The local.properties file must NOT be checked into Version Control Systems, as it contains information specific to your local configuration.
The location of the file is: %s`, localPropertiesPth)
	}
	return ""
}

// EnsureGradlew - Retuns the gradle wrapper path, or error if not exists
func EnsureGradlew(buildGradleFile string) (string, error) {
	projectDir := filepath.Dir(buildGradleFile)
	gradlewPth := filepath.Join(projectDir, "gradlew")
	if exist, err := pathutil.IsPathExists(gradlewPth); err != nil {
		return "", err
	} else if !exist {
		return "", errors.New(`<b>No Gradle Wrapper (gradlew) found.</b> 
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`)
	}
	return gradlewPth, nil
}

// GenerateOptions ...
func GenerateOptions(searchDir string) (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}

	buildGradlePths, err := CollectRootBuildGradleFiles(searchDir)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}

	gradleFileOption := models.NewOption(gradleFileInputTitle, gradleFileInputEnvKey)

	for _, buildGradlePth := range buildGradlePths {
		if warning := CheckLocalProperties(buildGradlePth); warning != "" {
			warnings = append(warnings, warning)
		}

		gradlewPth, err := EnsureGradlew(buildGradlePth)
		if err != nil {
			return models.OptionModel{}, warnings, err
		}

		gradlewPthOption := models.NewOption(gradlewPathInputTitle, gradlewPathInputEnvKey)
		gradleFileOption.AddOption(buildGradlePth, gradlewPthOption)

		configOption := models.NewConfigOption(ConfigName)
		gradlewPthOption.AddConfig(gradlewPth, configOption)
	}

	return *gradleFileOption, warnings, nil
}

// GenerateConfigBuilder ...
func GenerateConfigBuilder() models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder(true)

	configBuilder.AppendPreparStepList(steps.InstallMissingAndroidToolsStepListItem())

	configBuilder.AppendMainStepList(steps.GradleRunnerStepListItem(
		envmanModels.EnvironmentItemModel{gradleFileInputKey: "$" + gradleFileInputEnvKey},
		envmanModels.EnvironmentItemModel{gradleTaskInputKey: "assembleDebug"},
		envmanModels.EnvironmentItemModel{gradlewPathInputKey: "$" + gradlewPathInputEnvKey},
	))

	configBuilder.AddDefaultWorkflowBuilder(models.DeployWorkflowID, true)
	configBuilder.AppendPreparStepListTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem())

	configBuilder.AppendMainStepListTo(models.DeployWorkflowID, steps.GradleRunnerStepListItem(
		envmanModels.EnvironmentItemModel{gradleFileInputKey: "$" + gradleFileInputEnvKey},
		envmanModels.EnvironmentItemModel{gradleTaskInputKey: "assembleRelease"},
		envmanModels.EnvironmentItemModel{gradlewPathInputKey: "$" + gradlewPathInputEnvKey},
	))

	return *configBuilder
}
