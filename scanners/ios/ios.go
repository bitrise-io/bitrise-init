package ios

import (
	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
)

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	SearchDir         string
	ConfigDescriptors []ConfigDescriptor
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return string(utility.XcodeProjectTypeIOS)
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.SearchDir = searchDir

	detected, err := Detect(utility.XcodeProjectTypeIOS, searchDir)
	if err != nil {
		return false, err
	}

	return detected, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	options, configDescriptors, warnings, err := GenerateOptions(utility.XcodeProjectTypeIOS, scanner.SearchDir)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}

	scanner.ConfigDescriptors = configDescriptors

	return options, warnings, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionModel {
	return GenerateDefaultOptions(utility.XcodeProjectTypeIOS)
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return GenerateConfig(utility.XcodeProjectTypeIOS, scanner.ConfigDescriptors, true)
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return GenerateDefaultConfig(utility.XcodeProjectTypeIOS, true)
}
