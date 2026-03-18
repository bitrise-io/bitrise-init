package models

import (
	bitriseModels "github.com/bitrise-io/bitrise/v2/models"
	envmanModels "github.com/bitrise-io/envman/v2/models"
	stepmanModels "github.com/bitrise-io/stepman/models"
	"gopkg.in/yaml.v2"
)

// WorkflowID ...
type WorkflowID string

// PipelineID ...
type PipelineID string

const (
	// PrimaryWorkflowID ...
	PrimaryWorkflowID WorkflowID = "primary"
	// DeployWorkflowID ...
	DeployWorkflowID WorkflowID = "deploy"

	// FormatVersion ...
	FormatVersion = bitriseModels.FormatVersion

	defaultSteplibSource = "https://github.com/bitrise-io/bitrise-steplib.git"
)

type Container struct {
	Type    string                              `json:"type,omitempty" yaml:"type,omitempty"`
	Image   string                              `json:"image,omitempty" yaml:"image,omitempty"`
	Ports   []string                            `json:"ports,omitempty" yaml:"ports,omitempty"`
	Envs    []envmanModels.EnvironmentItemModel `json:"envs,omitempty" yaml:"envs,omitempty"`
	Options string                              `json:"options,omitempty" yaml:"options,omitempty"`
}

type BitriseConfig struct {
	bitriseModels.BitriseDataModel
	Tools      map[string]string    `json:"tools,omitempty" yaml:"tools,omitempty"`
	Containers map[string]Container `json:"containers,omitempty" yaml:"containers,omitempty"`
}

// MarshalYAML inlines BitriseDataModel's fields and appends Tools and Containers,
// overriding the embedded Containers field which has a conflicting yaml key.
func (c BitriseConfig) MarshalYAML() (interface{}, error) {
	base := c.BitriseDataModel
	base.Containers = nil

	baseBytes, err := yaml.Marshal(base)
	if err != nil {
		return nil, err
	}

	var result yaml.MapSlice
	if err := yaml.Unmarshal(baseBytes, &result); err != nil {
		return nil, err
	}

	if len(c.Tools) > 0 {
		result = append(result, yaml.MapItem{Key: "tools", Value: c.Tools})
	}
	if len(c.Containers) > 0 {
		result = append(result, yaml.MapItem{Key: "containers", Value: c.Containers})
	}

	return result, nil
}

type Step struct {
	stepmanModels.StepModel `yaml:",inline"`
	ServiceContainers       []string `json:"service_containers,omitempty" yaml:"service_containers,omitempty"`
}

// ConfigBuilderModel ...
type ConfigBuilderModel struct {
	workflowBuilderMap   map[WorkflowID]*workflowBuilderModel
	pipelineBuilderMap   map[PipelineID]*pipelineBuilderModel
	containerDefinitions map[string]Container
	tools                map[string]string
}

// NewDefaultConfigBuilder ...
func NewDefaultConfigBuilder() *ConfigBuilderModel {
	return &ConfigBuilderModel{
		workflowBuilderMap: map[WorkflowID]*workflowBuilderModel{},
		pipelineBuilderMap: map[PipelineID]*pipelineBuilderModel{},
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

// SetGraphPipelineWorkflowTo ...
func (builder *ConfigBuilderModel) SetGraphPipelineWorkflowTo(pipeline PipelineID, workflow WorkflowID, item bitriseModels.GraphPipelineWorkflowModel) {
	pipelineBuilder := builder.pipelineBuilderMap[pipeline]
	if pipelineBuilder == nil {
		pipelineBuilder = newDefaultPipelineBuilder()
		builder.pipelineBuilderMap[pipeline] = pipelineBuilder
	}
	pipelineBuilder.setGraphPipelineWorkflow(workflow, item)
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

// SetWorkflowSummaryTo ...
func (builder *ConfigBuilderModel) SetWorkflowSummaryTo(workflow WorkflowID, summary string) {
	workflowBuilder := builder.workflowBuilderMap[workflow]
	if workflowBuilder == nil {
		workflowBuilder = newDefaultWorkflowBuilder()
		builder.workflowBuilderMap[workflow] = workflowBuilder
	}
	workflowBuilder.Summary = summary
}

// SetContainerDefinitions ...
func (builder *ConfigBuilderModel) SetContainerDefinitions(containers map[string]Container) {
	builder.containerDefinitions = containers
}

// AddTool appends a tool with its version to the tools map.
func (builder *ConfigBuilderModel) AddTool(id string, version string) {
	if builder.tools == nil {
		builder.tools = map[string]string{}
	}
	builder.tools[id] = version
}

// Generate ...
func (builder *ConfigBuilderModel) Generate(projectType string, appEnvs ...envmanModels.EnvironmentItemModel) (BitriseConfig, error) {
	pipelines := map[string]bitriseModels.PipelineModel{}
	for pipelineID, pipelineBuilder := range builder.pipelineBuilderMap {
		pipelines[string(pipelineID)] = pipelineBuilder.generate()
	}

	workflows := map[string]bitriseModels.WorkflowModel{}
	for workflowID, workflowBuilder := range builder.workflowBuilderMap {
		workflows[string(workflowID)] = workflowBuilder.generate()
	}

	app := bitriseModels.AppModel{
		Environments: appEnvs,
	}

	core := bitriseModels.BitriseDataModel{
		FormatVersion:        FormatVersion,
		DefaultStepLibSource: defaultSteplibSource,
		ProjectType:          projectType,
		Pipelines:            pipelines,
		Workflows:            workflows,
		App:                  app,
	}
	return BitriseConfig{
		BitriseDataModel: core,
		Tools:            builder.tools,
		Containers:       builder.containerDefinitions,
	}, nil
}
