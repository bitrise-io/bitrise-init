package ruby

import (
	bitriseModels "github.com/bitrise-io/bitrise/v2/models"
	envmanModels "github.com/bitrise-io/envman/v2/models"
	stepmanModels "github.com/bitrise-io/stepman/models"
)

// serviceContainerDefinition is a mock for the upstream Container model with a Type field.
// TODO: Remove when upstream bitrise models include Type in Container.
type serviceContainerDefinition struct {
	Type    string                              `yaml:"type"`
	Image   string                              `yaml:"image"`
	Ports   []string                            `yaml:"ports,omitempty"`
	Envs    []envmanModels.EnvironmentItemModel `yaml:"envs,omitempty"`
	Options string                              `yaml:"options,omitempty"`
}

// bitriseConfigWithContainers wraps BitriseDataModel fields to produce the desired YAML
// with a `containers:` section that includes `type: service`.
// Uses explicit fields (not yaml:",inline") to avoid serializing BitriseDataModel.Services.
// TODO: Remove once upstream models support Type field on Container.
type bitriseConfigWithContainers struct {
	FormatVersion        string                                `yaml:"format_version"`
	DefaultStepLibSource string                                `yaml:"default_step_lib_source,omitempty"`
	ProjectType          string                                `yaml:"project_type"`
	Containers           map[string]serviceContainerDefinition `yaml:"containers,omitempty"`
	App                  bitriseModels.AppModel                `yaml:"app,omitempty"`
	Pipelines            map[string]bitriseModels.PipelineModel `yaml:"pipelines,omitempty"`
	Workflows            map[string]bitriseModels.WorkflowModel `yaml:"workflows,omitempty"`
}

// stepModelWithServiceContainers extends StepModel with service_containers field.
// TODO: Remove when upstream stepman models include ServiceContainers.
type stepModelWithServiceContainers struct {
	stepmanModels.StepModel `yaml:",inline"`
	ServiceContainers       []string `yaml:"service_containers,omitempty"`
}
