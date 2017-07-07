package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bitrise-core/bitrise-init/steps"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
)

const (
	// FormatVersion ...
	FormatVersion = bitriseModels.Version

	defaultSteplibSource = "https://github.com/bitrise-io/bitrise-steplib.git"
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
	bytes, err := json.MarshalIndent(option, "", "\t")
	if err != nil {
		return fmt.Sprintf("failed to marshal, error: %s", err)
	}

	return string(bytes)
}

// IsConfigOption ...
func (option *OptionModel) IsConfigOption() bool {
	return option.Config != ""
}

// IsValueOption ...
func (option *OptionModel) IsValueOption() bool {
	return option.Title != ""
}

// IsEmpty ...
func (option *OptionModel) IsEmpty() bool {
	return !option.IsValueOption() && !option.IsConfigOption()
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
			lastOptions = append(lastOptions, option)
			return
		}

		for _, childOption := range option.ChildOptionMap {
			if childOption == nil {
				lastOptions = append(lastOptions, option)
				return
			}

			if childOption.IsConfigOption() {
				lastOptions = append(lastOptions, option)
				return
			}

			if childOption.IsEmpty() {
				lastOptions = append(lastOptions, option)
				return
			}

			walk(childOption)
		}
	}

	walk(option)

	return lastOptions
}

// RemoveConfigs ...
func (option *OptionModel) RemoveConfigs() {
	lastChilds := option.LastChilds()
	for _, child := range lastChilds {
		for _, child := range child.ChildOptionMap {
			child.Config = ""
		}
	}
}

// AttachToLastChilds ...
func (option *OptionModel) AttachToLastChilds(opt *OptionModel) {
	childs := option.LastChilds()
	for _, child := range childs {
		values := child.GetValues()
		for _, value := range values {
			child.AddOption(value, opt)
		}
	}
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

// ---
// Config Builder

func newDefaultWorkflowBuilder(isIncludeCache bool) *workflowBuilderModel {
	return &workflowBuilderModel{
		PrepareSteps:    steps.DefaultPrepareStepList(isIncludeCache),
		DependencySteps: []bitriseModels.StepListItemModel{},
		MainSteps:       []bitriseModels.StepListItemModel{},
		DeploySteps:     steps.DefaultDeployStepList(isIncludeCache),
	}
}

func newWorkflowBuilder(items ...bitriseModels.StepListItemModel) *workflowBuilderModel {
	return &workflowBuilderModel{
		steps: items,
	}
}

func stepListItemEquals(stepListItem1, stepListItem2 bitriseModels.StepListItemModel) bool {
	stepID1 := ""
	for key := range stepListItem1 {
		stepID1 = key
		break
	}

	stepID2 := ""
	for key := range stepListItem2 {
		stepID2 = key
		break
	}

	return stepID1 == stepID2
}

func (builder *workflowBuilderModel) merge(workflowBuilder *workflowBuilderModel) {
	for _, stepListItemToCheck := range workflowBuilder.PrepareSteps {
		contains := false
		for _, stepListItem := range builder.PrepareSteps {
			if stepListItemEquals(stepListItem, stepListItemToCheck) {
				contains = true
				break
			}
		}
		if !contains {
			builder.appendPreparStepList(stepListItemToCheck)
		}
	}

	for _, stepListItemToCheck := range workflowBuilder.DependencySteps {
		contains := false
		for _, stepListItem := range builder.DependencySteps {
			if stepListItemEquals(stepListItem, stepListItemToCheck) {
				contains = true
				break
			}
		}
		if !contains {
			builder.appendPreparStepList(stepListItemToCheck)
		}
	}

	for _, stepListItemToCheck := range workflowBuilder.MainSteps {
		contains := false
		for _, stepListItem := range builder.MainSteps {
			if stepListItemEquals(stepListItem, stepListItemToCheck) {
				contains = true
				break
			}
		}
		if !contains {
			builder.appendPreparStepList(stepListItemToCheck)
		}
	}

	for _, stepListItemToCheck := range workflowBuilder.DeploySteps {
		contains := false
		for _, stepListItem := range builder.DeploySteps {
			if stepListItemEquals(stepListItem, stepListItemToCheck) {
				contains = true
				break
			}
		}
		if !contains {
			builder.appendPreparStepList(stepListItemToCheck)
		}
	}
}

func (builder *workflowBuilderModel) appendPreparStepList(items ...bitriseModels.StepListItemModel) {
	builder.PrepareSteps = append(builder.PrepareSteps, items...)
}

func (builder *workflowBuilderModel) appendDependencyStepList(items ...bitriseModels.StepListItemModel) {
	builder.DependencySteps = append(builder.DependencySteps, items...)
}

func (builder *workflowBuilderModel) appendMainStepList(items ...bitriseModels.StepListItemModel) {
	builder.MainSteps = append(builder.MainSteps, items...)
}

func (builder *workflowBuilderModel) appendDeployStepList(items ...bitriseModels.StepListItemModel) {
	builder.DeploySteps = append(builder.DeploySteps, items...)
}

func (builder *workflowBuilderModel) stepList() []bitriseModels.StepListItemModel {
	if len(builder.steps) > 0 {
		return builder.steps
	}

	stepList := []bitriseModels.StepListItemModel{}
	stepList = append(stepList, builder.PrepareSteps...)
	stepList = append(stepList, builder.DependencySteps...)
	stepList = append(stepList, builder.MainSteps...)
	stepList = append(stepList, builder.DeploySteps...)
	return stepList
}

func (builder *workflowBuilderModel) generate() bitriseModels.WorkflowModel {
	return bitriseModels.WorkflowModel{
		Steps: builder.stepList(),
	}
}

// NewDefaultConfigBuilder ...
func NewDefaultConfigBuilder(isIncludeCache bool) *ConfigBuilderModel {
	return &ConfigBuilderModel{
		workflowBuilderMap: map[WorkflowID]*workflowBuilderModel{
			PrimaryWorkflowID: newDefaultWorkflowBuilder(isIncludeCache),
		},
	}
}

// NewConfigBuilder ...
func NewConfigBuilder(primarySteps []bitriseModels.StepListItemModel) *ConfigBuilderModel {
	return &ConfigBuilderModel{
		workflowBuilderMap: map[WorkflowID]*workflowBuilderModel{
			PrimaryWorkflowID: newWorkflowBuilder(primarySteps...),
		},
	}
}

// AddDefaultWorkflowBuilder ...
func (builder *ConfigBuilderModel) AddDefaultWorkflowBuilder(workflow WorkflowID, isIncludeCache bool) {
	builder.workflowBuilderMap[workflow] = newDefaultWorkflowBuilder(isIncludeCache)
}

// AppendPreparStepListTo ...
func (builder *ConfigBuilderModel) AppendPreparStepListTo(workflow WorkflowID, items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[workflow]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[workflow] = workflowBuilder
	}
	workflowBuilder.appendPreparStepList(items...)
}

// AppendDependencyStepListTo ...
func (builder *ConfigBuilderModel) AppendDependencyStepListTo(workflow WorkflowID, items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[workflow]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[workflow] = workflowBuilder
	}
	workflowBuilder.appendDependencyStepList(items...)
}

// AppendMainStepListTo ...
func (builder *ConfigBuilderModel) AppendMainStepListTo(workflow WorkflowID, items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[workflow]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[workflow] = workflowBuilder
	}
	workflowBuilder.appendMainStepList(items...)
}

// AppendDeployStepListTo ...
func (builder *ConfigBuilderModel) AppendDeployStepListTo(workflow WorkflowID, items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[workflow]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[workflow] = workflowBuilder
	}
	workflowBuilder.appendDeployStepList(items...)
}

// AppendPreparStepList ...
func (builder *ConfigBuilderModel) AppendPreparStepList(items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[PrimaryWorkflowID]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[PrimaryWorkflowID] = workflowBuilder
	}
	workflowBuilder.appendPreparStepList(items...)
}

// AppendDependencyStepList ...
func (builder *ConfigBuilderModel) AppendDependencyStepList(items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[PrimaryWorkflowID]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[PrimaryWorkflowID] = workflowBuilder
	}
	workflowBuilder.appendDependencyStepList(items...)
}

// AppendMainStepList ...
func (builder *ConfigBuilderModel) AppendMainStepList(items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[PrimaryWorkflowID]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[PrimaryWorkflowID] = workflowBuilder
	}
	workflowBuilder.appendMainStepList(items...)
}

// AppendDeployStepList ...
func (builder *ConfigBuilderModel) AppendDeployStepList(items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[PrimaryWorkflowID]
	if workflowBuilder == nil {
		workflowBuilder = &workflowBuilderModel{}
		builder.workflowBuilderMap[PrimaryWorkflowID] = workflowBuilder
	}
	workflowBuilder.appendDeployStepList(items...)
}

// Merge ...
func (builder ConfigBuilderModel) Merge(configBuilder ConfigBuilderModel) ConfigBuilderModel {
	workflowIDs := []WorkflowID{}
	for workflowID := range builder.workflowBuilderMap {
		workflowIDs = append(workflowIDs, workflowID)
	}

	for _, workfloID := range workflowIDs {
		toMergeWorkflowBuilder, ok := configBuilder.workflowBuilderMap[workfloID]
		if ok {
			originalWorkflowBuilder := builder.workflowBuilderMap[workfloID]
			originalWorkflowBuilder.merge(toMergeWorkflowBuilder)
		}
	}

	return builder
}

// Generate ...
func (builder *ConfigBuilderModel) Generate(projectType string, appEnvs ...envmanModels.EnvironmentItemModel) (bitriseModels.BitriseDataModel, error) {
	primaryWorkflowBuilder, ok := builder.workflowBuilderMap[PrimaryWorkflowID]
	if !ok || primaryWorkflowBuilder == nil || len(primaryWorkflowBuilder.stepList()) == 0 {
		return bitriseModels.BitriseDataModel{}, errors.New("primary workflow not defined")
	}

	workflows := map[string]bitriseModels.WorkflowModel{}
	for workflowID, workflowBuilder := range builder.workflowBuilderMap {
		workflows[string(workflowID)] = workflowBuilder.generate()
	}

	triggerMap := []bitriseModels.TriggerMapItemModel{
		bitriseModels.TriggerMapItemModel{
			PushBranch: "*",
			WorkflowID: string(PrimaryWorkflowID),
		},
		bitriseModels.TriggerMapItemModel{
			PullRequestSourceBranch: "*",
			WorkflowID:              string(PrimaryWorkflowID),
		},
	}

	app := bitriseModels.AppModel{
		Environments: appEnvs,
	}

	return bitriseModels.BitriseDataModel{
		FormatVersion:        FormatVersion,
		DefaultStepLibSource: defaultSteplibSource,
		ProjectType:          projectType,
		TriggerMap:           triggerMap,
		Workflows:            workflows,
		App:                  app,
	}, nil
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
