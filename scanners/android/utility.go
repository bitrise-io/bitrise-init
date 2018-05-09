package android

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-android/gradle"
)

// Constants ...
const (
	ScannerName       = "android"
	ConfigName        = "android-config"
	DefaultConfigName = "default-android-config"

	ProjectLocationInputKey    = "project_location"
	ProjectLocationInputEnvKey = "PROJECT_LOCATION"
	ProjectLocationInputTitle  = "The root directory of an Android project"

	ModuleBuildGradlePathInputKey = "build_gradle_path"

	ModuleInputKey    = "module"
	ModuleInputEnvKey = "MODULE"
	ModuleInputTitle  = "Module in an Android project"

	VariantInputKey         = "variant"
	TestVariantInputEnvKey  = "TEST_VARIANT"
	BuildVariantInputEnvKey = "BUILD_VARIANT"
	TestVariantInputTitle   = "The variant for testing"
	BuildVariantInputTitle  = "The variant for building"

	gradlewPathInputKey    = "gradlew_path"
	gradlewPathInputEnvKey = "GRADLEW_PATH"
	gradlewPathInputTitle  = "Gradlew file path"
)

func walk(src string, fn func(path string, info os.FileInfo) error) error {
	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			if err := walk(filepath.Join(src, info.Name()), fn); err != nil {
				return err
			}
		}
		if err := fn(filepath.Join(src, info.Name()), info); err != nil {
			return err
		}
	}
	return nil
}

func checkFiles(path string, files ...string) (bool, error) {
	for _, file := range files {
		exists, err := pathutil.IsPathExists(filepath.Join(path, file))
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}
	return true, nil
}

func walkMultipleFiles(searchDir string, files ...string) (matches []string, err error) {
	match, err := checkFiles(searchDir, files...)
	if err != nil {
		return nil, err
	}
	if match {
		matches = append(matches, searchDir)
	}
	return matches, walk(searchDir, func(path string, info os.FileInfo) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			match, err := checkFiles(path, files...)
			if err != nil {
				return err
			}
			if match {
				matches = append(matches, path)
			}
		}
		return nil
	})
}

func checkLocalProperties(projectDir string) string {
	localPropertiesPth := filepath.Join(projectDir, "local.properties")
	exist, err := pathutil.IsPathExists(localPropertiesPth)
	if err == nil && exist {
		return fmt.Sprintf(`The local.properties file must NOT be checked into Version Control Systems, as it contains information specific to your local configuration.
The location of the file is: %s`, localPropertiesPth)
	}
	return ""
}

func checkGradlew(projectDir string) error {
	gradlewPth := filepath.Join(projectDir, "gradlew")
	exist, err := pathutil.IsPathExists(gradlewPth)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(`<b>No Gradle Wrapper (gradlew) found.</b> 
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`)
	}
	return nil
}

func (scanner *Scanner) generateOptions(searchDir string) (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}

	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputEnvKey)

	for _, projectRoot := range scanner.ProjectRoots {
		if warning := checkLocalProperties(projectRoot); warning != "" {
			warnings = append(warnings, warning)
		}

		if err := checkGradlew(projectRoot); err != nil {
			return models.OptionModel{}, warnings, err
		}

		proj, err := gradle.NewProject(projectRoot)
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
		testVariantsMap, err := proj.GetTask("test").GetVariants()
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
		buildVariantsMap, err := proj.GetTask("assemble").GetVariants()
		if err != nil {
			return models.OptionModel{}, warnings, err
		}

		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey)

		for module, variants := range testVariantsMap {
			variantOption := models.NewOption(TestVariantInputTitle, TestVariantInputEnvKey)
			buildVariantOption := models.NewOption(BuildVariantInputTitle, BuildVariantInputEnvKey)

			for _, variant := range variants {
				variant = strings.TrimSuffix(variant, "UnitTest")
				variantOption.AddOption(variant, buildVariantOption)
			}

			for _, variant := range buildVariantsMap[module] {
				configOption := models.NewConfigOption(ConfigName)
				buildVariantOption.AddOption(variant, configOption)
			}

			moduleOption.AddOption(module, variantOption)
		}

		gradlewPthOption := models.NewOption(gradlewPathInputTitle, gradlewPathInputEnvKey)
		gradlewPthOption.AddOption(filepath.Join(projectRoot, "gradlew"), moduleOption)

		projectLocationOption.AddOption(projectRoot, gradlewPthOption)
	}
	return *projectLocationOption, warnings, nil
}

// GenerateConfigBuilder ...
func GenerateConfigBuilder(isIncludeCache bool) models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	projectLocationEnv, moduleEnv, testVariantEnv, buildVariantEnv := "$"+ProjectLocationInputEnvKey, "$"+ModuleInputEnvKey, "$"+TestVariantInputEnvKey, "$"+BuildVariantInputEnvKey

	//-- primary
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.InstallMissingAndroidToolsStepListItem())
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidLintStepListItem(
		envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
		envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
		envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidUnitTestStepListItem(
		envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
		envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
		envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)

	//-- deploy
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem())

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ChangeAndroidVersionCodeAndVersionNameStepListItem(
		envmanModels.EnvironmentItemModel{ModuleBuildGradlePathInputKey: filepath.Join(projectLocationEnv, moduleEnv, "build.gradle")},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidLintStepListItem(
		envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
		envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
		envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidUnitTestStepListItem(
		envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
		envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
		envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
		envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
		envmanModels.EnvironmentItemModel{VariantInputKey: buildVariantEnv},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.SignAPKStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)

	return *configBuilder
}
