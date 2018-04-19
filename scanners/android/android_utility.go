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

	ModuleInputKey    = "module"
	ModuleInputEnvKey = "MODULE"
	ModuleInputTitle  = "Module in an Android project"

	VariantInputKey    = "variant"
	VariantInputEnvKey = "VARIANT"
	VariantInputTitle  = "The variant in the selected module"

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

func detect(searchDir string, files ...string) (matches []string, err error) {
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

		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey)
		modules, err := detect(projectRoot, "build.gradle", "src")
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
		for _, module := range modules {
			module = filepath.Base(module)

			proj, err := gradle.NewProject(projectRoot)
			if err != nil {
				return models.OptionModel{}, warnings, err
			}
			variants, err := proj.GetModule(module).GetTask("test").GetVariants()
			if err != nil {
				return models.OptionModel{}, warnings, err
			}

			variantOption := models.NewOption(VariantInputTitle, VariantInputEnvKey)

			for _, variant := range variants {
				configOption := models.NewConfigOption(ConfigName)

				variant = strings.Split(strings.Split(variant, "Debug")[0], "Release")[0]
				variantOption.AddOption(variant, configOption)
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

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.InstallMissingAndroidToolsStepListItem())
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidUnitTestStepListItem(
		envmanModels.EnvironmentItemModel{ProjectLocationInputKey: "$" + ProjectLocationInputEnvKey},
		envmanModels.EnvironmentItemModel{ModuleInputKey: "$" + ModuleInputEnvKey},
		envmanModels.EnvironmentItemModel{VariantInputKey: "$" + VariantInputEnvKey},
	))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.GradleRunnerStepListItem(
	// envmanModels.EnvironmentItemModel{GradleFileInputKey: "$" + GradleFileInputEnvKey},
	// envmanModels.EnvironmentItemModel{GradleTaskInputKey: "assembleRelease"},
	// envmanModels.EnvironmentItemModel{GradlewPathInputKey: "$" + GradlewPathInputEnvKey},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)

	return *configBuilder
}
