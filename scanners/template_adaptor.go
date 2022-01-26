package scanners

import (
	"github.com/bitrise-io/bitrise-init/builder"
	"github.com/bitrise-io/bitrise-init/models"
	"gopkg.in/yaml.v2"
)

type TemplateAdapter struct {
	template Template

	results map[string]builder.Result
}

func NewTemplateAdapter(template Template) *TemplateAdapter {
	return &TemplateAdapter{
		template: template,
	}
}

func (t *TemplateAdapter) Name() string {
	return t.template.Name()
}

func (t *TemplateAdapter) DetectPlatform(searchDir string) (bool, error) {
	return t.template.DetectPlatform(searchDir)
}

func (t *TemplateAdapter) ExcludedScannerNames() []string {
	return t.template.ExcludedScannerNames()
}

func (t *TemplateAdapter) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	templatepNode, err := t.template.Get()
	if err != nil {
		return models.OptionNode{}, nil, nil, err
	}

	answerTree, err := templatepNode.GetAnswers(map[string]builder.Question{}, []interface{}{})
	if err != nil {
		return models.OptionNode{}, nil, nil, err
	}

	options, results, err := builder.Export(templatepNode, answerTree, map[string]string{}, t.Name())
	if err != nil {
		return models.OptionNode{}, nil, nil, err
	}
	t.results = results

	return *options, nil, nil, nil
}

func (t *TemplateAdapter) Configs(isPrivateRepository bool) (models.BitriseConfigMap, error) {
	configMap := make(models.BitriseConfigMap)

	for configKey, result := range t.results {
		configContents, err := yaml.Marshal(result.Config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configMap[configKey] = string(configContents)
	}

	return configMap, nil
}

func (t *TemplateAdapter) DefaultOptions() models.OptionNode {
	panic("not implemented")
}

func (t *TemplateAdapter) DefaultConfigs() (models.BitriseConfigMap, error) {
	panic("not impemented")
}
