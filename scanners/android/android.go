package android

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android/icon"
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
	return nil
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (_ bool, err error) {
	scanner.SearchDir = searchDir

	scanner.ProjectRoots, err = walkMultipleFiles(searchDir, "build.gradle", "settings.gradle")
	if err != nil {
		return false, fmt.Errorf("failed to search for build.gradle files, error: %s", err)
	}

	kotlinRoots, err := walkMultipleFiles(searchDir, "build.gradle.kts", "settings.gradle.kts")
	if err != nil {
		return false, fmt.Errorf("failed to search for build.gradle files, error: %s", err)
	}

	scanner.ProjectRoots = append(scanner.ProjectRoots, kotlinRoots...)

	return len(scanner.ProjectRoots) > 0, err
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputEnvKey)
	warnings := models.Warnings{}
	appIconsAllProjects := models.Icons{}

	for _, projectRoot := range scanner.ProjectRoots {
		if err := checkGradlew(projectRoot); err != nil {
			return models.OptionNode{}, warnings, models.Icons{}, err
		}

		relProjectRoot, err := filepath.Rel(scanner.SearchDir, projectRoot)
		if err != nil {
			return models.OptionNode{}, warnings, models.Icons{}, err
		}

		appIcons, err := icon.GetAllIcons(projectRoot, scanner.SearchDir)
		if err != nil {
			return models.OptionNode{}, warnings, models.Icons{}, err
		}
		for iconID, iconPath := range appIcons {
			appIconsAllProjects[iconID] = iconPath
		}

		configOption := models.NewConfigOption(ConfigName)
		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey)
		variantOption := models.NewOption(VariantInputTitle, VariantInputEnvKey)

		projectLocationOption.AddOption(relProjectRoot, moduleOption)
		moduleOption.AddOption("app", variantOption)
		variantOption.AddConfig("", configOption)

		if len(appIcons) == 0 {
			variantOption.AddConfig("", configOption)
		} else {
			iconOption := models.NewOption(appIconTitle, "")
			iconOption.SetSelectorType(models.IconSelector)
			variantOption.AddConfig("", iconOption)
			for iconID := range appIcons {
				iconOption.AddConfig(iconID, configOption)
			}
		}
	}

	return *projectLocationOption, warnings, models.Icons{}, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputEnvKey)
	moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey)
	variantOption := models.NewOption(VariantInputTitle, VariantInputEnvKey)
	configOption := models.NewConfigOption(DefaultConfigName)

	projectLocationOption.AddOption("_", moduleOption)
	moduleOption.AddOption("_", variantOption)
	variantOption.AddConfig("", configOption)

	return *projectLocationOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configBuilder := scanner.generateConfigBuilder()

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
	configBuilder := scanner.generateConfigBuilder()

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
