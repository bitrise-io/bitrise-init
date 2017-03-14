package ios

import "github.com/bitrise-core/bitrise-init/models"
import "github.com/bitrise-core/bitrise-init/scanners/xcode"

// ScannerName ...
const ScannerName = "ios"

// Scanner ...
type Scanner struct {
	*xcode.Scanner
}

var wrapperScanner = *xcode.NewScanner(xcode.ProjectTypeIOS)

// Name ...
func (scanner *Scanner) Name() string { return ScannerName }

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	return wrapperScanner.DetectPlatform(searchDir)
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	return wrapperScanner.Options()
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	return wrapperScanner.DefaultOptions()
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return wrapperScanner.Configs()
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return wrapperScanner.DefaultConfigs()
}
