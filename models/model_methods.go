package models

import (
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
)

const (
	// FormatVersion ...
	FormatVersion        = "1.3.1"
	defaultSteplibSource = "https://github.com/bitrise-io/bitrise-steplib.git"
	primaryWorkflowID    = "primary"
	deployWorkflowID     = "deploy"
)

// ---
// OptionModel

// NewOption ...
func NewOption(title, envKey string) *OptionModel {
	return &OptionModel{
		Title:          title,
		EnvKey:         envKey,
		ChildOptionMap: map[string]*OptionModel{},
	}
}

// NewConfigOption ...
func NewConfigOption(name string) *OptionModel {
	return &OptionModel{
		ChildOptionMap: map[string]*OptionModel{},
		Config:         name,
	}
}

// AddOption ...
func (option *OptionModel) AddOption(forValue string, newOption *OptionModel) {
	option.ChildOptionMap[forValue] = newOption
}

// AddConfig ...
func (option *OptionModel) AddConfig(forValue string, newConfigOption *OptionModel) {
	option.ChildOptionMap[forValue] = newConfigOption
}

// Child ...
func (option *OptionModel) Child(components ...string) (*OptionModel, bool) {
	currentOption := option
	for _, component := range components {
		childOption := currentOption.ChildOptionMap[component]
		if childOption == nil {
			return nil, false
		}
		currentOption = childOption
	}
	return currentOption, true
}

// LastOptions ...
func (option *OptionModel) LastOptions() []*OptionModel {
	lastOptions := []*OptionModel{}

	var walkDepth func(option *OptionModel)

	walkDepth = func(option *OptionModel) {
		if len(option.ChildOptionMap) == 0 {
			// no more child, this is the last option in this branch
			lastOptions = append(lastOptions, option)
			return
		}
		for _, childOption := range option.ChildOptionMap {
			if childOption == nil {
				// values are set to this option, but has value without child
				lastOptions = append(lastOptions, option)
				return
			}

			walkDepth(childOption)
		}
	}

	walkDepth(option)

	return lastOptions
}

// GetValues ...
func (option *OptionModel) GetValues() []string {
	if option.Config != "" {
		return []string{option.Config}
	}

	values := []string{}
	for value := range option.ChildOptionMap {
		values = append(values, value)
	}
	return values
}

// ---

// AddError ...
func (result *ScanResultModel) AddError(platform string, errorMessage string) {
	if result.PlatformErrorsMap == nil {
		result.PlatformErrorsMap = map[string]Errors{}
	}
	if result.PlatformErrorsMap[platform] == nil {
		result.PlatformErrorsMap[platform] = []string{}
	}
	result.PlatformErrorsMap[platform] = append(result.PlatformErrorsMap[platform], errorMessage)
}

// BitriseDataWithCIWorkflow ...
func BitriseDataWithCIWorkflow(appEnvs []envmanModels.EnvironmentItemModel, steps []bitriseModels.StepListItemModel) bitriseModels.BitriseDataModel {
	workflows := map[string]bitriseModels.WorkflowModel{
		primaryWorkflowID: bitriseModels.WorkflowModel{
			Steps: steps,
		},
	}

	triggerMap := []bitriseModels.TriggerMapItemModel{
		bitriseModels.TriggerMapItemModel{
			PushBranch: "*",
			WorkflowID: primaryWorkflowID,
		},
		bitriseModels.TriggerMapItemModel{
			PullRequestSourceBranch: "*",
			WorkflowID:              primaryWorkflowID,
		},
	}

	app := bitriseModels.AppModel{
		Environments: appEnvs,
	}

	return bitriseModels.BitriseDataModel{
		FormatVersion:        FormatVersion,
		DefaultStepLibSource: defaultSteplibSource,
		TriggerMap:           triggerMap,
		Workflows:            workflows,
		App:                  app,
	}
}

// BitriseDataWithCIAndCDWorkflow ...
func BitriseDataWithCIAndCDWorkflow(appEnvs []envmanModels.EnvironmentItemModel, ciSteps, deploySteps []bitriseModels.StepListItemModel) bitriseModels.BitriseDataModel {
	workflows := map[string]bitriseModels.WorkflowModel{
		primaryWorkflowID: bitriseModels.WorkflowModel{
			Steps: ciSteps,
		},
		deployWorkflowID: bitriseModels.WorkflowModel{
			Steps: deploySteps,
		},
	}

	triggerMap := []bitriseModels.TriggerMapItemModel{
		bitriseModels.TriggerMapItemModel{
			PushBranch: "*",
			WorkflowID: primaryWorkflowID,
		},
		bitriseModels.TriggerMapItemModel{
			PullRequestSourceBranch: "*",
			WorkflowID:              primaryWorkflowID,
		},
	}

	app := bitriseModels.AppModel{
		Environments: appEnvs,
	}

	return bitriseModels.BitriseDataModel{
		FormatVersion:        FormatVersion,
		DefaultStepLibSource: defaultSteplibSource,
		TriggerMap:           triggerMap,
		Workflows:            workflows,
		App:                  app,
	}
}
