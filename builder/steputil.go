package builder

import "github.com/bitrise-io/bitrise-init/steps"

func DefaultPrepareStepsTemplate(isIncludeCache bool) *Steps {
	stepsList := ActivateSSHKeyStepTemplate()
	stepsList.Append(Step{ID: steps.GitCloneID})

	if isIncludeCache {
		stepsList.Append(Step{ID: steps.CachePullID})
	}

	return stepsList.Append(Step{
		ID:    steps.ScriptID,
		Title: steps.ScriptDefaultTitle,
	})
}

func ActivateSSHKeyStepTemplate() *Steps {
	stepList := Steps{}
	return stepList.Append(Step{
		ID:    steps.ActivateSSHKeyID,
		RunIf: `{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}`,
	})
}
