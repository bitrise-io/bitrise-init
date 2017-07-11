package models

import (
	"github.com/bitrise-core/bitrise-init/steps"
	bitriseModels "github.com/bitrise-io/bitrise/models"
)

type workflowBuilderModel struct {
	PrepareSteps    []bitriseModels.StepListItemModel
	DependencySteps []bitriseModels.StepListItemModel
	MainSteps       []bitriseModels.StepListItemModel
	DeploySteps     []bitriseModels.StepListItemModel
}

func newDefaultWorkflowBuilder(isIncludeCache bool) *workflowBuilderModel {
	return &workflowBuilderModel{
		PrepareSteps:    steps.DefaultPrepareStepList(isIncludeCache),
		DependencySteps: []bitriseModels.StepListItemModel{},
		MainSteps:       []bitriseModels.StepListItemModel{},
		DeploySteps:     steps.DefaultDeployStepList(isIncludeCache),
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
