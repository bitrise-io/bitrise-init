package models

import (
	"encoding/json"
	"fmt"

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
		Components:     []string{},
	}
}

// NewConfigOption ...
func NewConfigOption(name string) *OptionModel {
	return &OptionModel{
		ChildOptionMap: map[string]*OptionModel{},
		Config:         name,
		Components:     []string{},
	}
}

func (option *OptionModel) String() string {
	if option.Config != "" {
		return fmt.Sprintf(`Config Option:
  config: %s
`, option.Config)
	}

	values := option.GetValues()
	return fmt.Sprintf(`Option:
  title: %s
  env_key: %s
  values: %v
`, option.Title, option.EnvKey, values)
}

// AddOption ...
func (option *OptionModel) AddOption(forValue string, newOption *OptionModel) {
	option.ChildOptionMap[forValue] = newOption

	if newOption != nil {
		newOption.Components = append(option.Components, forValue)

		if option.Head == nil {
			// first option's head is nil
			newOption.Head = option
		} else {
			newOption.Head = option.Head
		}
	}
}

// AddConfig ...
func (option *OptionModel) AddConfig(forValue string, newConfigOption *OptionModel) {
	option.ChildOptionMap[forValue] = newConfigOption
	if newConfigOption != nil {
		newConfigOption.Components = append(option.Components, forValue)

		if option.Head == nil {
			// first option's head is nil
			newConfigOption.Head = option
		} else {
			newConfigOption.Head = option.Head
		}
	}
}

// Parent ...
func (option *OptionModel) Parent() (*OptionModel, string, bool) {
	if option.Head == nil {
		return nil, "", false
	}

	parentComponents := option.Components[:len(option.Components)-1]
	parentOption, ok := option.Head.Child(parentComponents...)
	if !ok {
		return nil, "", false
	}
	underKey := option.Components[len(option.Components)-1:][0]
	return parentOption, underKey, true
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

// LastChilds ...
func (option *OptionModel) LastChilds() []*OptionModel {
	lastOptions := []*OptionModel{}

	var walk func(option *OptionModel)
	walk = func(option *OptionModel) {
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

			walk(childOption)
		}
	}

	walk(option)

	return lastOptions
}

// Copy ...
func (option *OptionModel) Copy() *OptionModel {
	bytes, err := json.Marshal(*option)
	if err != nil {
		return nil
	}

	var optionCopy OptionModel
	if err := json.Unmarshal(bytes, &optionCopy); err != nil {
		return nil
	}

	return &optionCopy
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
