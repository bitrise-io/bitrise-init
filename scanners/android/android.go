package android

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-core/bitrise-init/models"
)

// Scanner ...
type Scanner struct {
	SearchDir    string
	ProjectRoots []string
	ExcludeTest  bool
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return ScannerName
}

// ExcludedScannerNames ...
func (*Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (b bool, err error) {
	scanner.SearchDir = searchDir

	scanner.ProjectRoots, err = walkMultipleFiles(searchDir, "build.gradle", "settings.gradle")
	if err != nil {
		return false, fmt.Errorf("failed to search for build.gradle files, error: %s", err)
	}

	return len(scanner.ProjectRoots) > 0, err
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	return scanner.generateOptions(scanner.SearchDir)
}

// DefaultOptions ...
func (*Scanner) DefaultOptions() models.OptionModel {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputEnvKey)

	moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey)
	projectLocationOption.AddOption("_", moduleOption)

	testVariantOption := models.NewOption(TestVariantInputTitle, TestVariantInputEnvKey)
	moduleOption.AddOption("_", testVariantOption)

	buildVariantOption := models.NewOption(BuildVariantInputTitle, BuildVariantInputEnvKey)
	testVariantOption.AddOption("_", buildVariantOption)

	configOption := models.NewConfigOption(DefaultConfigName)
	buildVariantOption.AddConfig("_", configOption)

	return *projectLocationOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configBuilder := scanner.generateConfigBuilder(true)

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
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := scanner.generateConfigBuilder(true)

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
