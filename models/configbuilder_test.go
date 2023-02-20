package models

import (
	"testing"

	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/stretchr/testify/require"
)

func TestConfigGenerateHaveProjectType(t *testing.T) {
	config := NewDefaultConfigBuilder()
	config.AppendStepListItemsTo("primary", steps.DefaultPrepareStepList(steps.PrepareListParams{})...)

	model, err := config.Generate("iOS")

	require.Nil(t, err)
	require.Equal(t, "iOS", model.ProjectType)
}

func TestConfigGenerateDoesNotHaveTriggerMap(t *testing.T) {
	config := NewDefaultConfigBuilder()
	config.AppendStepListItemsTo("primary", steps.DefaultPrepareStepList(steps.PrepareListParams{})...)

	model, err := config.Generate("iOS")

	require.Nil(t, err)
	require.Equal(t, 0, len(model.TriggerMap))
}
