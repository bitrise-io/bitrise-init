package builder

import "github.com/bitrise-io/bitrise-init/steps"

type PrepareListParams struct {
	ShouldIncludeCache       bool
	ShouldIncludeActivateSSH bool
}

func DefaultPrepareStepsTemplate(params PrepareListParams) *Steps {
	var stepList Steps
	if params.ShouldIncludeActivateSSH {
		stepList.Append(Step{ID: steps.ActivateSSHKeyID})
	}

	stepList.Append(Step{ID: steps.GitCloneID})

	if params.ShouldIncludeCache {
		stepList.Append(Step{ID: steps.CachePullID})
	}

	return stepList.Append(Step{
		ID:    steps.ScriptID,
		Title: steps.ScriptDefaultTitle,
	})
}
