package models

import (
	"testing"

	bitriseModels "github.com/bitrise-io/bitrise/v2/models"
	stepmanModels "github.com/bitrise-io/stepman/models"
	"github.com/stretchr/testify/require"
)

func TestConfigGenerateHaveProjectType(t *testing.T) {
	config := NewDefaultConfigBuilder()
	config.AppendStepListItemsTo("primary", []bitriseModels.StepListItemModel{
		{"step-id": stepmanModels.StepModel{}},
	}...)

	model, err := config.Generate("iOS")

	require.Nil(t, err)
	require.Equal(t, "iOS", model.ProjectType)
}

func TestConfigDoesNotGenerateTriggerMap(t *testing.T) {
	config := NewDefaultConfigBuilder()
	config.AppendStepListItemsTo("primary", []bitriseModels.StepListItemModel{
		{"step-id": stepmanModels.StepModel{}},
	}...)

	model, err := config.Generate("iOS")

	require.Nil(t, err)
	require.Nil(t, model.TriggerMap)
}
