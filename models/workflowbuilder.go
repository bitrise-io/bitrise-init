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

func (builder *workflowBuilderModel) removePreparStepList(item bitriseModels.StepListItemModel) {
	updatedPrepareSteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.PrepareSteps {
		for stepID := range stepListItem {
			shouldRemove := false

			for stepIDToRemove := range item {
				if stepIDToRemove == stepID {
					shouldRemove = true
				}
				break
			}

			if !shouldRemove {
				updatedPrepareSteps = append(updatedPrepareSteps, stepListItem)
			}
			break
		}
	}
	builder.PrepareSteps = updatedPrepareSteps
}

func (builder *workflowBuilderModel) appendDependencyStepList(items ...bitriseModels.StepListItemModel) {
	builder.DependencySteps = append(builder.DependencySteps, items...)
}

func (builder *workflowBuilderModel) removeDependencyStepList(item bitriseModels.StepListItemModel) {
	updatedDependencySteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.DependencySteps {
		for stepID := range stepListItem {
			shouldRemove := false

			for stepIDToRemove := range item {
				if stepIDToRemove == stepID {
					shouldRemove = true
				}
				break
			}

			if !shouldRemove {
				updatedDependencySteps = append(updatedDependencySteps, stepListItem)
			}
			break
		}
	}
	builder.DependencySteps = updatedDependencySteps
}

func (builder *workflowBuilderModel) appendMainStepList(items ...bitriseModels.StepListItemModel) {
	builder.MainSteps = append(builder.MainSteps, items...)
}

func (builder *workflowBuilderModel) removeMainStepList(item bitriseModels.StepListItemModel) {
	updatedMainSteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.MainSteps {
		for stepID := range stepListItem {
			shouldRemove := false

			for stepIDToRemove := range item {
				if stepIDToRemove == stepID {
					shouldRemove = true
				}
				break
			}

			if !shouldRemove {
				updatedMainSteps = append(updatedMainSteps, stepListItem)
			}
			break
		}
	}
	builder.MainSteps = updatedMainSteps
}

func (builder *workflowBuilderModel) appendDeployStepList(items ...bitriseModels.StepListItemModel) {
	builder.DeploySteps = append(builder.DeploySteps, items...)
}

func (builder *workflowBuilderModel) removeDeployStepList(item bitriseModels.StepListItemModel) {
	updatedDeploySteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.DeploySteps {
		for stepID := range stepListItem {
			shouldRemove := false

			for stepIDToRemove := range item {
				if stepIDToRemove == stepID {
					shouldRemove = true
				}
				break
			}

			if !shouldRemove {
				updatedDeploySteps = append(updatedDeploySteps, stepListItem)
			}
			break
		}
	}
	builder.DeploySteps = updatedDeploySteps
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

// RemoveStepListItem ...
func (builder *workflowBuilderModel) RemoveStepListItem(stepListItemID string) {
	updatedPrepareSteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.PrepareSteps {
		for stepID := range stepListItem {
			if stepID != stepListItemID {
				updatedPrepareSteps = append(updatedPrepareSteps, stepListItem)
			}
			break
		}
	}
	builder.PrepareSteps = updatedPrepareSteps

	updatedDependencySteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.DependencySteps {
		for stepID := range stepListItem {
			if stepID != stepListItemID {
				updatedDependencySteps = append(updatedDependencySteps, stepListItem)
			}
			break
		}
	}
	builder.DependencySteps = updatedDependencySteps

	updatedMainSteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.MainSteps {
		for stepID := range stepListItem {
			if stepID != stepListItemID {
				updatedMainSteps = append(updatedMainSteps, stepListItem)
			}
			break
		}
	}
	builder.MainSteps = updatedMainSteps

	updatedDeploySteps := []bitriseModels.StepListItemModel{}
	for _, stepListItem := range builder.DeploySteps {
		for stepID := range stepListItem {
			if stepID != stepListItemID {
				updatedDeploySteps = append(updatedDeploySteps, stepListItem)
			}
			break
		}
	}
	builder.DeploySteps = updatedDeploySteps
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
