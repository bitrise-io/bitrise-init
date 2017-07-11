package reactnative

import (
	"errors"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Name ...
const Name = "react-native"

// Scanner ...
type Scanner struct {
	SearchDir      string
	iosScanner     *ios.Scanner
	androidScanner *android.Scanner
	hasNPMTest     bool
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (scanner *Scanner) Name() string {
	return Name
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.SearchDir = searchDir

	packageJSONPths, err := CollectPackageJSONFiles(searchDir)
	if err != nil {
		return false, err
	}

	return (len(packageJSONPths) > 0), nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}

	packageJSONPths, err := CollectPackageJSONFiles(scanner.SearchDir)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}

	// react options
	packageJSONPth := packageJSONPths[0]
	packages, err := utility.ParsePackagesJSON(packageJSONPth)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}

	if _, found := packages.Scripts["test"]; found {
		scanner.hasNPMTest = true
	}

	projectDir := filepath.Dir(packageJSONPth)

	// android options
	var androidOptions *models.OptionModel
	androidDir := filepath.Join(projectDir, "android")
	if exist, err := pathutil.IsDirExists(androidDir); err != nil {
		return models.OptionModel{}, warnings, err
	} else if exist {
		androidScanner := android.NewScanner()

		if detected, err := androidScanner.DetectPlatform(scanner.SearchDir); err != nil {
			return models.OptionModel{}, warnings, err
		} else if detected {
			options, warns, err := androidScanner.Options()
			warnings = append(warnings, warns...)
			if err != nil {
				return models.OptionModel{}, warnings, err
			}

			androidOptions = &options
			scanner.androidScanner = androidScanner
		}
	}

	// ios options
	var iosOptions *models.OptionModel
	iosDir := filepath.Join(projectDir, "ios")
	if exist, err := pathutil.IsDirExists(iosDir); err != nil {
		return models.OptionModel{}, warnings, err
	} else if exist {
		iosScanner := ios.NewScanner()

		if detected, err := iosScanner.DetectPlatform(scanner.SearchDir); err != nil {
			return models.OptionModel{}, warnings, err
		} else if detected {
			options, warns, err := iosScanner.Options()
			warnings = append(warnings, warns...)
			if err != nil {
				return models.OptionModel{}, warnings, err
			}

			iosOptions = &options
			scanner.iosScanner = iosScanner
		}
	}

	if androidOptions == nil && iosOptions == nil {
		return models.OptionModel{}, warnings, errors.New("no ios nor android config options found")
	}

	var options *models.OptionModel
	if androidOptions != nil {
		if iosOptions == nil {
			// we only found an android project
			// we need to update the config names
			lastChilds := androidOptions.LastChilds()
			for _, child := range lastChilds {
				for _, child := range child.ChildOptionMap {
					if child.Config == "" {
						return models.OptionModel{}, warnings, fmt.Errorf("no config for option: %s", child.String())
					}

					configName := configName(true, nil)
					child.Config = configName
				}
			}
		} else {
			// we have both ios and android projects
			// we need to remove the android option's config names,
			// since ios options will hold them
			androidOptions.RemoveConfigs()
		}

		options = androidOptions
	}

	if iosOptions != nil {
		lastChilds := iosOptions.LastChilds()
		for _, child := range lastChilds {
			for _, child := range child.ChildOptionMap {
				if child.Config == "" {
					return models.OptionModel{}, warnings, fmt.Errorf("no config for option: %s", child.String())
				}

				descriptor := ios.NewConfigDescriptorWithName(child.Config)
				configName := configName(scanner.androidScanner != nil, &descriptor)
				child.Config = configName
			}
		}

		if androidOptions == nil {
			// we only found an ios project
			options = iosOptions
		} else {
			// we have both ios and android projects
			// we attach ios options to the android options
			options.AttachToLastChilds(iosOptions)
		}

	}

	return *options, warnings, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{}
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	configBuilder := models.NewDefaultConfigBuilder(true)

	if scanner.androidScanner != nil {
		androidConfigBuilder := android.GenerateConfigBuilder(true)
		configBuilder = &androidConfigBuilder
	}

	// ---
	if scanner.iosScanner != nil {
		descriptors := scanner.iosScanner.ConfigDescriptors
		descriptors = ios.RemoveDuplicatedConfigDescriptors(descriptors, utility.XcodeProjectTypeIOS)

		for _, descriptor := range descriptors {
			iosConfigBuilder := ios.GenerateConfigBuilder(utility.XcodeProjectTypeIOS, descriptor.HasPodfile, descriptor.HasTest, descriptor.MissingSharedSchemes, descriptor.CarthageCommand, true)
			mergedBuilder := configBuilder.Merge(iosConfigBuilder)

			bitriseDataModel, err := mergedBuilder.Generate(Name)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			data, err := yaml.Marshal(bitriseDataModel)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			configName := configName(scanner.androidScanner != nil, &descriptor)
			configMap[configName] = string(data)
		}
	} else {
		bitriseDataModel, err := configBuilder.Generate(Name)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(bitriseDataModel)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configName := configName(scanner.androidScanner != nil, nil)
		configMap[configName] = string(data)
	}

	// ---

	return configMap, nil
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}

// ExcludedScannerNames ...
func (scanner *Scanner) ExcludedScannerNames() []string {
	return []string{
		string(utility.XcodeProjectTypeIOS),
		string(utility.XcodeProjectTypeMacOS),
		android.ScannerName,
	}
}
