package scanner

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/goinp/goinp"
	yaml "gopkg.in/yaml.v2"
)

func askForOptionValue(option models.OptionNode) (string, string, error) {
	optionValues := option.GetValues()

	selectedValue := ""
	if len(optionValues) == 1 {
		if optionValues[0] == "_" {
			// provide option value
			question := fmt.Sprintf("Provide: %s", option.Title)
			answer, err := goinp.AskForString(question)
			if err != nil {
				return "", "", err
			}

			selectedValue = answer
		} else {
			// auto select the only one value
			selectedValue = optionValues[0]
		}
	} else {
		// select from values
		question := fmt.Sprintf("Select: %s", option.Title)
		answer, err := goinp.SelectFromStrings(question, optionValues)
		if err != nil {
			return "", "", err
		}

		selectedValue = answer
	}

	return option.EnvKey, selectedValue, nil
}

// AskForOptions ...
func AskForOptions(options models.OptionNode) (string, map[string]string, error) {
	log.Printf("AskForOptions options: %s", options)
	configPth := ""
	substitutions := map[string]string{}

	var walkDepth func(models.OptionNode) error
	walkDepth = func(opt models.OptionNode) error {
		optionEnvKey, selectedValue, err := askForOptionValue(opt)
		if err != nil {
			return fmt.Errorf("Failed to ask for value, error: %s", err)
		}

		if opt.Title == "" {
			// last option selected, config got
			configPth = selectedValue
			return nil
		} else if optionEnvKey != "" {
			// env's value selected
			substitutions[optionEnvKey] = selectedValue
		}

		var nestedOptions *models.OptionNode
		if len(opt.ChildOptionMap) == 1 {
			// auto select the next option
			for _, childOption := range opt.ChildOptionMap {
				nestedOptions = childOption
				break
			}
		} else {
			// go to the next option, based on the selected value
			childOptions, found := opt.ChildOptionMap[selectedValue]
			if !found {
				return nil
			}
			nestedOptions = childOptions
		}

		return walkDepth(*nestedOptions)
	}

	if err := walkDepth(options); err != nil {
		return "", map[string]string{}, err
	}

	if configPth == "" {
		return "", nil, errors.New("no config selected")
	}

	log.Printf("AskForOptions configPth: %s, appEnvs: %s", configPth, substitutions)
	return configPth, substitutions, nil
}

// AskForConfig ...
func AskForConfig(scanResult models.ScanResultModel) (bitriseModels.BitriseDataModel, error) {

	//
	// Select platform
	platforms := []string{}
	for platform := range scanResult.ScannerToOptionRoot {
		platforms = append(platforms, platform)
	}

	platform := ""
	if len(platforms) == 0 {
		return bitriseModels.BitriseDataModel{}, errors.New("no platform detected")
	} else if len(platforms) == 1 {
		platform = platforms[0]
	} else {
		var err error
		platform, err = goinp.SelectFromStrings("Select platform", platforms)
		if err != nil {
			return bitriseModels.BitriseDataModel{}, err
		}
	}
	// ---

	//
	// Select config
	options, ok := scanResult.ScannerToOptionRoot[platform]
	if !ok {
		return bitriseModels.BitriseDataModel{}, fmt.Errorf("invalid platform selected: %s", platform)
	}

	configPth, substitutions, err := AskForOptions(options)
	if err != nil {
		return bitriseModels.BitriseDataModel{}, err
	}
	// --

	//
	// Build config
	configMap := scanResult.ScannerToBitriseConfigMap[platform]
	configStr := configMap[configPth]
	return substituteChosenOptionsInConfig(configStr, substitutions)
}

func substituteChosenOptionsInConfig(configStr string, substitutions map[string]string) (bitriseModels.BitriseDataModel, error) {
	log.Printf("substituteChosenOptionsInConfig configStr: %s", configStr)

	executeTemplate := func(text string) (string, error) {
		tmpl, err := template.New("bitrise.yml with scanner defined options").
			Delims(utility.TemplateDelimiterLeft, utility.TemplateDelimiterRight).
			Parse(text)
		if err != nil {
			return text, fmt.Errorf("failed to parse bitrise.yml template, error: %s", err)
		}
		var byteBuffer bytes.Buffer
		err = tmpl.Execute(&byteBuffer, substitutions)
		if err != nil {
			return text, fmt.Errorf("failed to execute bitrise.yml tempalte, error: %s", err)
		}
		return byteBuffer.String(), nil
	}

	// Parse bitrise.yml as a templated text, and substitute options
	var err error
	configStr, err = executeTemplate(configStr)
	if err != nil {
		return bitriseModels.BitriseDataModel{}, err
	}

	var config bitriseModels.BitriseDataModel
	if err := yaml.Unmarshal([]byte(configStr), &config); err != nil {
		return bitriseModels.BitriseDataModel{}, fmt.Errorf("failed to unmarshal config, error: %s", err)
	}
	return config, nil
}
