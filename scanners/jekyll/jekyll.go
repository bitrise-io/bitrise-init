package jekyll

import (
	"fmt"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"strings"
)

// Scanner ...
type Scanner struct {
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return ScannerName
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return false, fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}

	// Search for _config.yml file
	// Note: hexo.io also uses _config.yml, but package.json for dependencies.
	// This means having config file is not enough to detect Jekyll platform.
	// _config.yml + Gemfile (with jekyll gem) = Jekyll platform
	log.TInfof("Searching for %s file", configYmlFile)
	configYmlPath, err := filterProjectFile(configYmlFile, fileList)
	if err != nil {
		return false, fmt.Errorf("failed to search for %s file, error: %s", configYmlFile, err)
	}
	log.TPrintf("%s found", configYmlFile)

	// Search for Jekyll in Gemfile
	log.TInfof("Searching for %s file", gemfileFile)
	gemfilePath, err := filterProjectFile(gemfileFile, fileList)
	if err != nil {
		return false, fmt.Errorf("failed to search for %s file, error: %s", gemfileFile, err)
	}
	log.TPrintf("%s found", gemfileFile)

	if configYmlPath == "" || gemfilePath == "" {
		log.TPrintf("platform not detected")
		return false, nil
	}

	gemfileContent, err := readGemfileToString()
	if err != nil {
		log.TPrintf("can not read Gemfile, error: %s", err)
		log.TPrintf("platform not detected")
		return false, nil
	}

	// ensure jekyll gem is in Gemfile
	if !strings.Contains(gemfileContent, jekyllGemName) {
		log.TPrintf("%s does not contain %s gem", gemfileFile, jekyllGemName)
		log.TPrintf("platform not detected")
		return false, nil
	}

	log.TSuccessf("Platform detected")

	return true, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	return models.OptionModel{}, nil, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{}
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}
