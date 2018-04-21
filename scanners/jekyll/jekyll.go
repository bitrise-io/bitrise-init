package jekyll

import (
	"fmt"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
	"gopkg.in/yaml.v2"
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

	// typical Jekyll project structure:
	// "_config.yml" = mandatory
	// either "_posts" or "_layouts" dir = mandatory
	// "_includes"/"_data"/"_OTHER_DIRS" dirs = optional
	// "Gemfile" / "jekyll gem in Gemfile" = optional

	// Note: hexo.io also uses _config.yml, but package.json for dependencies and layout directory instead of _layouts.
	// This means having config file is not enough to detect Jekyll platform.
	// _config.yml + ("_posts" || "_layouts") = Jekyll platform
	// _config.yml + package.json + layout = hexo.io platform

	// Search for _config.yml file
	log.TInfof("Searching for %s file", configYmlFile)
	configYmlPath, err := filterProjectFile(configYmlFile, fileList)
	if err != nil {
		log.TWarnf("failed to search for %s file, error: %s", configYmlFile, err)
		return false, nil
	}
	log.TPrintf("%s found", configYmlFile)

	// Search _posts directory
	log.TInfof("Searching for %s directory", postsDirectory)
	postsDirExists := utility.DirExists(postsDirectory)
	if postsDirExists {
		log.TPrintf("%s found", postsDirectory)
	} else {
		log.TWarnf("failed to search for %s directory", postsDirectory)
	}

	// Search _includes directory
	log.TInfof("Searching for %s directory", includesDirectory)
	includesDirExists := utility.DirExists(includesDirectory)
	if includesDirExists {
		log.TPrintf("%s found", includesDirectory)
	} else {
		log.TWarnf("failed to search for %s directory", includesDirectory)
	}

	// config is mandatory
	if configYmlPath == "" {
		log.TPrintf("platform not detected")
		return false, nil
	}
	// at least one dir(_posts/_includes) is mandatory
	if !postsDirExists && !includesDirExists {
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
	return models.OptionModel{Config: ConfigName}, nil, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{Config: DefaultConfigName}
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configBuilder := GenerateConfigBuilder(true)

	config, err := configBuilder.Generate(ScannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		ConfigName: string(data),
	}, nil

}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := GenerateConfigBuilder(true)

	config, err := configBuilder.Generate(ScannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		DefaultConfigName: string(data),
	}, nil
}
