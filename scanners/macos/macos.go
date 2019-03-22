package macos

import (
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
)

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	searchDir         string
	configDescriptors []ios.ConfigDescriptor
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return string(ios.XcodeProjectTypeMacOS)
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.searchDir = searchDir

	detected, err := ios.Detect(ios.XcodeProjectTypeMacOS, searchDir)
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
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	options, configDescriptors, _, warnings, err := ios.GenerateOptions(ios.XcodeProjectTypeMacOS, scanner.searchDir, true)
	if err != nil {
		return models.OptionNode{}, warnings, models.Icons{}, err
	}

	scanner.configDescriptors = configDescriptors

	return options, warnings, models.Icons{}, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionNode {
	return ios.GenerateDefaultOptions(ios.XcodeProjectTypeMacOS)
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return ios.GenerateConfig(ios.XcodeProjectTypeMacOS, scanner.configDescriptors, true)
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return ios.GenerateDefaultConfig(ios.XcodeProjectTypeMacOS, true)
}
