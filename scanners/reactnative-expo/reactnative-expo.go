package expo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/scanners/reactnative"
	"github.com/bitrise-core/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	yaml "gopkg.in/yaml.v2"
)

// Name ...
const Name = "react-native-expo"

const (
	workDirInputKey = "workdir"
)

// Scanner ...
type Scanner struct {
	searchDir          string
	iosScanner         *ios.Scanner
	androidScanner     *android.Scanner
	reactnativeScanner *reactnative.Scanner
	hasNPMTest         bool
	packageJSONPth     string
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return Name
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.searchDir = searchDir

	log.TInfof("Collect package.json files")

	packageJSONPths, err := reactnative.CollectPackageJSONFiles(searchDir)
	if err != nil {
		return false, err
	}

	if len(packageJSONPths) == 0 {
		return false, nil
	}

	dependencyFound := false
	for _, packageJSONPth := range packageJSONPths {
		dependency := "expo"
		dependencyFound, err := FindDependency(packageJSONPth, dependency)
		if err != nil {
			fmt.Printf("Error during finding dependency: %s", err.Error())
			return false, err
		}

		fmt.Printf("%s found: %t\n", dependency, dependencyFound)
		if dependencyFound {
			break
		}
	}

	if !dependencyFound {
		return false, nil
	}

	log.TPrintf("%d package.json file detected", len(packageJSONPths))

	log.TInfof("Filter relevant package.json files")

	iosScanner := ios.NewScanner()
	androidScanner := android.NewScanner()
	for _, packageJSONPth := range packageJSONPths {
		log.TPrintf("checking: %s", packageJSONPth)

		projectDir := filepath.Dir(packageJSONPth)

		iosDir := filepath.Join(projectDir, "ios")
		if exist, err := pathutil.IsDirExists(iosDir); err != nil {
			return false, err
		} else if exist {
			if detected, err := iosScanner.DetectPlatform(scanner.searchDir); err != nil {
				return false, err
			} else if detected {
				return false, nil
			}
		}

		androidDir := filepath.Join(projectDir, "android")
		if exist, err := pathutil.IsDirExists(androidDir); err != nil {
			return false, err
		} else if exist {
			if detected, err := androidScanner.DetectPlatform(scanner.searchDir); err != nil {
				return false, err
			} else if detected {
				return false, nil
			}
		}
	}
	if err := ejectProject(searchDir); err != nil {
		log.Errorf("ERROR DURING EJECTING THE PROJECT: %s", err)
		return false, nil
	}

	reactnativeScanner := reactnative.NewScanner()
	scanner.reactnativeScanner = reactnativeScanner

	return reactnativeScanner.DetectPlatform(scanner.searchDir)
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	return scanner.reactnativeScanner.Options()
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionModel {
	gradleFileOption := models.NewOption(android.GradleFileInputTitle, android.GradleFileInputEnvKey)

	gradlewPthOption := models.NewOption(android.GradlewPathInputTitle, android.GradlewPathInputEnvKey)
	gradleFileOption.AddOption("_", gradlewPthOption)

	projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
	gradlewPthOption.AddOption("_", projectPathOption)

	schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)
	projectPathOption.AddOption("_", schemeOption)

	exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
	schemeOption.AddOption("_", exportMethodOption)

	for _, exportMethod := range ios.IosExportMethods {
		configOption := models.NewConfigOption(defaultConfigName())
		exportMethodOption.AddConfig(exportMethod, configOption)
	}

	return *gradleFileOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return scanner.reactnativeScanner.Configs()
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	// ci
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{"command": "install"}))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{"command": "test"}))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

	// cd
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{"command": "install"}))

	// android
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.GradleRunnerStepListItem(
		envmanModels.EnvironmentItemModel{android.GradleFileInputKey: "$" + android.GradleFileInputEnvKey},
		envmanModels.EnvironmentItemModel{android.GradleTaskInputKey: "assembleRelease"},
		envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.GradlewPathInputEnvKey},
	))

	// ios
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
		envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

	bitriseDataModel, err := configBuilder.Generate(Name)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(bitriseDataModel)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	configName := defaultConfigName()
	configMap := models.BitriseConfigMap{
		configName: string(data),
	}
	return configMap, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return nil
}

func ejectProject(pth string) error {
	log.Infof("Eject project")

	cmd := command.New("npm", "run", "eject")
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	return cmd.Run()
}
