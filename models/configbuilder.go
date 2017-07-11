package models

import (
	"errors"

	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
)

// WorkflowID ...
type WorkflowID string

const (
	// PrimaryWorkflowID ...
	PrimaryWorkflowID WorkflowID = "primary"
	// DeployWorkflowID ...
	DeployWorkflowID WorkflowID = "deploy"

	// FormatVersion ...
	FormatVersion = bitriseModels.Version

	defaultSteplibSource = "https://github.com/bitrise-io/bitrise-steplib.git"
)

// ConfigBuilderModel ...
type ConfigBuilderModel struct {
	workflowBuilderMap map[WorkflowID]*workflowBuilderModel
}

// NewDefaultConfigBuilder ...
func NewDefaultConfigBuilder(isIncludeCache bool) *ConfigBuilderModel {
	return &ConfigBuilderModel{
		workflowBuilderMap: map[WorkflowID]*workflowBuilderModel{
			PrimaryWorkflowID: newDefaultWorkflowBuilder(isIncludeCache),
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

// RemoveStepListItem ...
func (builder *ConfigBuilderModel) RemoveStepListItem(workflowID WorkflowID, stepListItemID string) {
	_, found := builder.workflowBuilderMap[workflowID]
	if !found {
		return
	}

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
