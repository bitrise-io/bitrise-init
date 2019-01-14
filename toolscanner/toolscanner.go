package toolscanner

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
)

// ProjectTypeTemplateKey is the name of the enviroment variable used to substitute the project type for
// automation tool scanner's config
const (
	ProjectTypeUserTitle   = "Project type"
	ProjectTypeTemplateKey = "PROJECT_TYPE"
)

// AddProjectTypeToToolScanner is used to add a project type for automation tool scanners's option tree and config map
func AddProjectTypeToToolScanner(scannerOptionTree models.OptionNode, scannerConfigMap models.BitriseConfigMap, detectedProjectTypes []string) (models.OptionNode, models.BitriseConfigMap, error) {
	log.Printf("AddProjectTypeToToolScanner old toolScannerOptionTree: %s, detectedProjectTypes: %s", scannerOptionTree, detectedProjectTypes)

	// For each tool scanner generated config, multiply it to get all combinations of 'Existing config' X 'Detected project type'
	configMapWithProjecTypes := map[string]string{}
	for _, projectType := range detectedProjectTypes {
		for configName, config := range scannerConfigMap {
			configWithProjectType, err := evaluateConfigTemplate(config, map[string]string{ProjectTypeTemplateKey: projectType})
			if err != nil {
				return scannerOptionTree, scannerConfigMap, fmt.Errorf("failed to add project type to tool scanner bitrise.yml, error: %s", err)
			}
			configMapWithProjecTypes[configName+"_"+projectType] = configWithProjectType
		}
	}

	// add the possible project types as a question to the option map
	optionsTreeWithProjectTypeRoot := models.NewOption(ProjectTypeUserTitle, ProjectTypeTemplateKey)
	for _, projectType := range detectedProjectTypes {
		optionsTreeWithProjectTypeRoot.AddOption(projectType, appendProjectTypeToConfig(scannerOptionTree, projectType))
	}

	log.Printf("AddProjectTypeToToolScanner new options root: %s", optionsTreeWithProjectTypeRoot)
	return *optionsTreeWithProjectTypeRoot, configMapWithProjecTypes, nil
}

func evaluateConfigTemplate(configStr string, substitutions map[string]string) (string, error) {
	log.Printf("substituteChosenOptionsInConfig configStr: %s", configStr)
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
	for _, leafNode := range options.LastChilds() {
		if leafNode.Config != "" {
			leafNode.Config = appendProjectTypeToConfigName(options.Config, projectType)
		}
	}
	return &options
}
