package reactnative

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"gopkg.in/yaml.v2"
)

const (
	expoConfigName        = "react-native-expo-config"
	expoDefaultConfigName = "default-" + expoConfigName
)

const (
	projectRootDirInputTitle   = "Project root directory"
	projectRootDirInputSummary = "The directory of the 'app.json' or 'package.json' file of your React Native project."
)

const wordirEnv = "WORKDIR"

// expoOptions implements ScannerInterface.Options function for Expo based React Native projects.
func (scanner *Scanner) expoOptions() (models.OptionNode, models.Warnings, error) {
	return models.OptionNode{}, models.Warnings{}, nil
}

// expoConfigs implements ScannerInterface.Configs function for Expo based React Native projects.
func (scanner *Scanner) expoConfigs(isPrivateRepo bool) (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	// determine workdir
	packageJSONDir := filepath.Dir(scanner.packageJSONPth)
	relPackageJSONDir, err := utility.RelPath(scanner.searchDir, packageJSONDir)
	if err != nil {
		return models.BitriseConfigMap{}, fmt.Errorf("Failed to get relative package.json dir path, error: %s", err)
	}
	if relPackageJSONDir == "." {
		// package.json placed in the search dir, no need to change-dir in the workflows
		relPackageJSONDir = ""
	}
	log.TPrintf("Working directory: %v", relPackageJSONDir)

	// primary workflow
	configBuilder := models.NewDefaultConfigBuilder()
	primaryDescription := primaryExpoWorkflowNoTestsDescription
	if scanner.hasTest {
		primaryDescription = primaryExpoWorkflowDescription
	}

	configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, primaryDescription)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
		ShouldIncludeCache:       false,
		ShouldIncludeActivateSSH: isPrivateRepo,
	})...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, scanner.getTestSteps(relPackageJSONDir)...)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepListV2(false)...)

	// deploy workflow
	// TODO: deploy wf description update
	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployExpoWorkflowDescription)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
		ShouldIncludeCache:       false,
		ShouldIncludeActivateSSH: isPrivateRepo,
	})...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, scanner.getTestSteps(relPackageJSONDir)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.RunEASBuildStepListItem(envmanModels.EnvironmentItemModel{"work_dir": relPackageJSONDir}))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

	// generate bitrise.yml
	bitriseDataModel, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(bitriseDataModel)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	configMap[expoConfigName] = string(data)

	return configMap, nil
}

// expoDefaultOptions implements ScannerInterface.DefaultOptions function for Expo based React Native projects.
func (Scanner) expoDefaultOptions() models.OptionNode {
	// TODO: update options with Expo wording
	workDirOption := models.NewOption(projectRootDirInputTitle, projectRootDirInputSummary, wordirEnv, models.TypeUserInput)
	return *workDirOption
}

// expoDefaultConfigs implements ScannerInterface.DefaultConfigs function for Expo based React Native projects.
func (scanner Scanner) expoDefaultConfigs() (models.BitriseConfigMap, error) {
	// TODO: should we ask if test, if yarn, which platform to deploy?
	configMap := models.BitriseConfigMap{}

	// primary workflow
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, primaryExpoWorkflowDescription)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
		ShouldIncludeCache:       false,
		ShouldIncludeActivateSSH: true,
	})...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, getTestSteps("$WORKDIR", true, true)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepListV2(false)...)

	// deploy workflow
	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployExpoWorkflowDescription)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
		ShouldIncludeCache:       false,
		ShouldIncludeActivateSSH: true,
	})...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, getTestSteps("$WORKDIR", true, true)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.RunEASBuildStepListItem(envmanModels.EnvironmentItemModel{"work_dir": "$WORKDIR"}))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

	// generate bitrise.yml
	bitriseDataModel, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(bitriseDataModel)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	configMap[expoDefaultConfigName] = string(data)

	return configMap, nil
}
