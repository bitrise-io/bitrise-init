package reactnative

import (
	"errors"
	"fmt"
	"path/filepath"

	yaml "gopkg.in/yaml.v1"

	"strings"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/scanners/xcode"
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

	packageJSONPth := packageJSONPths[0]
	projectDir := filepath.Dir(packageJSONPth)

	hasAndroidProjectDir := false
	androidDir := filepath.Join(projectDir, "android")
	{
		var err error
		hasAndroidProjectDir, err = pathutil.IsDirExists(androidDir)
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
	}

	var androidOptions *models.OptionModel
	if hasAndroidProjectDir {
		androidScanner := android.NewScanner()

		if detected, err := androidScanner.DetectPlatform(scanner.SearchDir); err != nil {
			return models.OptionModel{}, warnings, err
		} else if detected {
			options, warns, err := androidScanner.Options()
			warnings = append(warnings, warns...)
			if err != nil {
				return models.OptionModel{}, warnings, err
			}

			options.RemoveConfigs()
			androidOptions = &options

			scanner.androidScanner = androidScanner
		}
	}

	hasIosProjectDir := false
	iosDir := filepath.Join(projectDir, "ios")
	{
		var err error
		hasIosProjectDir, err = pathutil.IsDirExists(iosDir)
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
	}

	var iosOptions *models.OptionModel
	if hasIosProjectDir {
		iosScanner := ios.NewScanner()

		if detected, err := iosScanner.DetectPlatform(scanner.SearchDir); err != nil {
			return models.OptionModel{}, warnings, err
		} else if detected {
			options, warns, err := iosScanner.Options()
			warnings = append(warnings, warns...)
			if err != nil {
				return models.OptionModel{}, warnings, err
			}

			// these options will be the last options
			// we need to update the last options config names
			// currently they are the ios config names
			lastChilds := options.LastChilds()
			for _, child := range lastChilds {
				for _, child := range child.ChildOptionMap {
					if child.Config == "" {
						return models.OptionModel{}, warnings, fmt.Errorf("no config for option: %s", child.String())
					}

					descriptor := xcode.NewConfigDescriptorWithName(child.Config)
					configName := configName(scanner.androidScanner != nil, &descriptor)
					child.Config = configName
				}
			}
			// ---

			iosOptions = &options

			scanner.iosScanner = iosScanner
		}
	} else if scanner.androidScanner != nil {
		// no ios project detected, but android found
		// we did not updated the last option's config names
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
	}

	if androidOptions == nil && iosOptions == nil {
		return models.OptionModel{}, warnings, errors.New("no ios nor android config options found")
	}

	rootOption := models.NewOption("project_dir", "PROJECT_DIR")
	if androidOptions != nil {
		rootOption.AddOption(projectDir, androidOptions)
	}
	if iosOptions != nil {
		if androidOptions != nil {
			rootOption.AttachToLastChilds(iosOptions)
		} else {
			rootOption.AddOption(projectDir, iosOptions)
		}
	}

	// return *rootOption, warnings, nil
	return models.OptionModel{}, models.Warnings{}, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{}
}

func configName(hasAndroidProject bool, iosConfigDescriptor *xcode.ConfigDescriptor) string {
	name := "reactnative"
	if hasAndroidProject {
		name += "-android"
	}
	if iosConfigDescriptor != nil {
		name += "-" + iosConfigDescriptor.ConfigName(utility.XcodeProjectTypeIOS)
	}
	if !strings.HasSuffix(name, "-config") {
		name += "-config"
	}
	return name
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	configBuilder := models.NewDefaultConfigBuilder(true)

	if scanner.androidScanner != nil {
		androidConfigBuilder := android.GenerateConfigBuilder()
		configBuilder = &androidConfigBuilder
	}

	// ---
	if scanner.iosScanner != nil {
		descriptors := scanner.iosScanner.ConfigDescriptors
		descriptors = xcode.RemoveDuplicatedConfigDescriptors(descriptors, utility.XcodeProjectTypeIOS)

		for _, descriptor := range descriptors {
			iosConfigBuilder := xcode.GenerateConfigBuilder(utility.XcodeProjectTypeIOS, descriptor.HasPodfile, descriptor.HasTest, descriptor.MissingSharedSchemes, descriptor.CarthageCommand)
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
