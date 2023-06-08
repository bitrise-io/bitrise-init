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
	defaultBranch      string
	workflowBuilderMap map[WorkflowID]*workflowBuilderModel
}

// NewDefaultConfigBuilder ...
func NewDefaultConfigBuilder(defaultBranch string) *ConfigBuilderModel {
	return &ConfigBuilderModel{
		defaultBranch: defaultBranch,
		workflowBuilderMap: map[WorkflowID]*workflowBuilderModel{
			PrimaryWorkflowID: newDefaultWorkflowBuilder(),
		},
	}
}

// AppendStepListItemsTo ...
func (builder *ConfigBuilderModel) AppendStepListItemsTo(workflow WorkflowID, items ...bitriseModels.StepListItemModel) {
	workflowBuilder := builder.workflowBuilderMap[workflow]
	if workflowBuilder == nil {
		workflowBuilder = newDefaultWorkflowBuilder()
		builder.workflowBuilderMap[workflow] = workflowBuilder
	}
	workflowBuilder.appendStepListItems(items...)
}

// SetWorkflowDescriptionTo ...
func (builder *ConfigBuilderModel) SetWorkflowDescriptionTo(workflow WorkflowID, description string) {
	workflowBuilder := builder.workflowBuilderMap[workflow]
	if workflowBuilder == nil {
		workflowBuilder = newDefaultWorkflowBuilder()
		builder.workflowBuilderMap[workflow] = workflowBuilder
	}
	workflowBuilder.Description = description
}

// Generate ...
func (builder *ConfigBuilderModel) Generate(projectType string, appEnvs ...envmanModels.EnvironmentItemModel) (bitriseModels.BitriseDataModel, error) {
	primaryWorkflowBuilder, ok := builder.workflowBuilderMap[PrimaryWorkflowID]
	if !ok || primaryWorkflowBuilder == nil || len(primaryWorkflowBuilder.Steps) == 0 {
		return bitriseModels.BitriseDataModel{}, errors.New("primary workflow not defined")
	}

	workflows := map[string]bitriseModels.WorkflowModel{}
	for workflowID, workflowBuilder := range builder.workflowBuilderMap {
		workflows[string(workflowID)] = workflowBuilder.generate()
	}

	pushBranchPattern := "*"
	if builder.defaultBranch != "" {
		pushBranchPattern = builder.defaultBranch
	}
	triggerMap := []bitriseModels.TriggerMapItemModel{
		bitriseModels.TriggerMapItemModel{
			PushBranch: pushBranchPattern,
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
