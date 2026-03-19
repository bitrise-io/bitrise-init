package models

import (
	bitriseModels "github.com/bitrise-io/bitrise/v2/models"
	envmanModels "github.com/bitrise-io/envman/v2/models"
	stepmanModels "github.com/bitrise-io/stepman/models"
	"gopkg.in/yaml.v2"
)

// TODO: The structs and functions in this file exist as a workaround because the Bitrise CLI requires
// 	a higher Go version than what is pre-installed on stacks. Since Steps run as Go source code, this
// 	tool must target the stack's pre-installed Go version. Once the Go version constraint is resolved,
// 	these types should be removed in favor of importing them directly from the bitrise/v2 dependency.

type BitriseConfig struct {
	bitriseModels.BitriseDataModel
	Tools      map[string]string    `json:"tools,omitempty" yaml:"tools,omitempty"`
	Containers map[string]Container `json:"containers,omitempty" yaml:"containers,omitempty"`
}

type Container struct {
	Type    string                              `json:"type,omitempty" yaml:"type,omitempty"`
	Image   string                              `json:"image,omitempty" yaml:"image,omitempty"`
	Ports   []string                            `json:"ports,omitempty" yaml:"ports,omitempty"`
	Envs    []envmanModels.EnvironmentItemModel `json:"envs,omitempty" yaml:"envs,omitempty"`
	Options string                              `json:"options,omitempty" yaml:"options,omitempty"`
}

type Step struct {
	stepmanModels.StepModel `yaml:",inline"`
	ServiceContainers       []string `json:"service_containers,omitempty" yaml:"service_containers,omitempty"`
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
