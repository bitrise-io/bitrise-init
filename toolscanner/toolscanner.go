package toolscanner

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
)

// ProjectTypeTemplateKey is the name of the enviroment variable used to substitute the project type for
// automation tool scanner's config
const (
	ProjectTypeUserTitle   = "Project type"
	ProjectTypeTemplateKey = "PROJECT_TYPE"
)

// AddProjectTypeToConfig get cartesian product: 'Existing tool scanner generated config' X 'Detected project type'
func AddProjectTypeToConfig(scannerConfigMap models.BitriseConfigMap, detectedProjectTypes []string) (models.BitriseConfigMap, error) {
	configMapWithProjecTypes := map[string]string{}
	for _, projectType := range detectedProjectTypes {
		for configName, config := range scannerConfigMap {
			configWithProjectType, err := evaluateConfigTemplate(config,
				map[string]string{ProjectTypeTemplateKey: projectType})
			if err != nil {
				return nil,
					fmt.Errorf("failed to add project type to tool scanner bitrise.yml, error: %s", err)
			}
			configMapWithProjecTypes[appendProjectTypeToConfigName(configName, projectType)] = configWithProjectType
		}
	}
	return configMapWithProjecTypes, nil
}

// AddProjectTypeToOptions adds a project type question to automation tool scanners's option tree
func AddProjectTypeToOptions(scannerOptionTree models.OptionNode, detectedProjectTypes []string) models.OptionNode {
	optionsTreeWithProjectTypeRoot := models.NewOption(ProjectTypeUserTitle, ProjectTypeTemplateKey)
	for _, projectType := range detectedProjectTypes {
		optionsTreeWithProjectTypeRoot.AddOption(projectType, appendProjectTypeToConfig(scannerOptionTree, projectType))
	}
	return *optionsTreeWithProjectTypeRoot
}

func evaluateConfigTemplate(configStr string, substitutions map[string]string) (string, error) {
	// Parse bitrise.yml as a templated text, and substitute options
	tmpl, err := template.New("bitrise.yml with scanner defined options").
		Delims(utility.TemplateDelimiterLeft, utility.TemplateDelimiterRight).
		Parse(configStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse bitrise.yml template, error: %s", err)
	}
	var byteBuffer bytes.Buffer
	err = tmpl.Execute(&byteBuffer, substitutions)
	if err != nil {
		return "", fmt.Errorf("failed to execute bitrise.yml tempalte, error: %s", err)
	}
	return byteBuffer.String(), nil
}

func appendProjectTypeToConfigName(configName string, projectType string) string {
	return configName + "_" + projectType
}

func appendProjectTypeToConfig(options models.OptionNode, projectType string) *models.OptionNode {
	var appendToConfigNames func(*models.OptionNode)
	appendToConfigNames = func(node *models.OptionNode) {
		if (*node).IsConfigOption() || (*node).ChildOptionMap == nil {
			(*node).Config = appendProjectTypeToConfigName((*node).Config, projectType)
			return
		}
		for _, child := range (*node).ChildOptionMap {
			appendToConfigNames(child)
		}
	}
	appendToConfigNames(&options)
	return &options
}
