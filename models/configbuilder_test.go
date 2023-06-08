package models

import (
	"testing"

	bitriseModels "github.com/bitrise-io/bitrise/models"
	stepmanModels "github.com/bitrise-io/stepman/models"
	"github.com/stretchr/testify/require"
)

func TestConfigGenerateHaveProjectType(t *testing.T) {
	config := NewDefaultConfigBuilder("")
	config.AppendStepListItemsTo("primary", []bitriseModels.StepListItemModel{
		{"step-id": stepmanModels.StepModel{}},
	}...)

	model, err := config.Generate("iOS")

	require.Nil(t, err)
	require.Equal(t, "iOS", model.ProjectType)
}

func TestConfigGenerateHaveTriggerMap(t *testing.T) {
	config := NewDefaultConfigBuilder("main")
	config.AppendStepListItemsTo("primary", []bitriseModels.StepListItemModel{
		{"step-id": stepmanModels.StepModel{}},
	}...)

	model, err := config.Generate("iOS")

	require.Nil(t, err)
	require.Equal(t, bitriseModels.TriggerMapModel{
		bitriseModels.TriggerMapItemModel{
			PushBranch: "main",
			WorkflowID: "primary",
		},
		bitriseModels.TriggerMapItemModel{
			PullRequestSourceBranch: "*",
			WorkflowID:              "primary",
		},
	}, model.TriggerMap)
}
