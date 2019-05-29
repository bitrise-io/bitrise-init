package scanner

import (
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/models"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/goinp/goinp"
)

func askForOptionValue(option models.OptionNode) (string, string, error) {
	// this options is a last element in a tree, contains only config name
	if option.Config != "" {
		return "", option.Config, nil
	}

	var (
		answer string
		err    error
	)

	switch option.Type {
	case models.TypeSelector, models.TypeOptionalSelector:
		answer, err = goinp.AskOptions(option.Title, "", option.Type == models.TypeOptionalSelector, option.GetValues()...)
	case models.TypeUserInput, models.TypeOptionalUserInput:
		var defaultValue string
		for key := range option.ChildOptionMap {
			defaultValue = key
			break
		}
		answer, err = goinp.AskOptions(option.Title, defaultValue, option.Type == models.TypeOptionalUserInput)
	}

	return option.EnvKey, strings.TrimSpace(answer), err
}

// AskForOptions ...
func AskForOptions(options models.OptionNode) (string, []envmanModels.EnvironmentItemModel, error) {
	configPth := ""
	appEnvs := []envmanModels.EnvironmentItemModel{}

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
			appEnvs = append(appEnvs, envmanModels.EnvironmentItemModel{
				optionEnvKey: selectedValue,
			})
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
		return "", []envmanModels.EnvironmentItemModel{}, err
	}

	if configPth == "" {
		return "", nil, errors.New("no config selected")
	}

	return configPth, appEnvs, nil
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
		platform, err = goinp.AskOptions("platform", "", false, platforms...)
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

	configPth, appEnvs, err := AskForOptions(options)
	if err != nil {
		return bitriseModels.BitriseDataModel{}, err
	}
	// --

	//
	// Build config
	configMap := scanResult.ScannerToBitriseConfigMap[platform]
	configStr := configMap[configPth]

	var config bitriseModels.BitriseDataModel
	if err := yaml.Unmarshal([]byte(configStr), &config); err != nil {
		return bitriseModels.BitriseDataModel{}, fmt.Errorf("failed to unmarshal config, error: %s", err)
	}

	config.App.Environments = append(config.App.Environments, appEnvs...)
	// ---

	return config, nil
}
