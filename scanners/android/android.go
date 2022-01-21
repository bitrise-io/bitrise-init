package android

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/analytics"
	"github.com/bitrise-io/bitrise-init/builder"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/log"
)

// Scanner ...
type Scanner struct {
	SearchDir      string
	ProjectRoots   []string
	ExcludeTest    bool
	ExcludeAppIcon bool
}

// Template is the v2 implementation
type Template struct {
	SearchDir    string
	ProjectRoots []string
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

	projectFiles := fileGroups{
		{"build.gradle", "build.gradle.kts"},
		{"settings.gradle", "settings.gradle.kts"},
	}
	skipDirs := []string{".git", "CordovaLib", "node_modules"}

	log.TInfof("Searching for android files")

	scanner.ProjectRoots, err = walkMultipleFileGroups(searchDir, projectFiles, skipDirs)
	if err != nil {
		return false, fmt.Errorf("failed to search for build.gradle files, error: %s", err)
	}

	log.TSuccessf("Platform detected")

	for _, file := range projectFiles {
		log.TPrintf("- %s", file)
	}
	countDetected := len(scanner.ProjectRoots)
	log.TPrintf("%d android files detected", countDetected)
	nonZeroCount := countDetected > 0
	return nonZeroCount, err
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputSummary, ProjectLocationInputEnvKey, models.TypeSelector)
	warnings := models.Warnings{}
	appIconsAllProjects := models.Icons{}

	foundOptions := false
	var lastErr error = nil
	for _, projectRoot := range scanner.ProjectRoots {
		exists, err := containsLocalProperties(projectRoot)
		if err != nil {
			lastErr = err
			continue
		}
		if exists {
			containsLocalPropertiesWarning := fmt.Sprintf("the local.properties file should NOT be checked into Version Control Systems, as it contains information specific to your local configuration, the location of the file is: %s", filepath.Join(projectRoot, "local.properties"))
			warnings = append(warnings, containsLocalPropertiesWarning)
		}

		if err := checkGradlew(projectRoot); err != nil {
			lastErr = err
			continue
		}

		relProjectRoot, err := filepath.Rel(scanner.SearchDir, projectRoot)
		if err != nil {
			lastErr = err
			continue
		}

		icons, err := LookupIcons(projectRoot, scanner.SearchDir)
		if err != nil {
			analytics.LogInfo("android-icon-lookup", analytics.DetectorErrorData("android", err), "Failed to lookup android icon")
		}
		appIconsAllProjects = append(appIconsAllProjects, icons...)
		iconIDs := make([]string, len(icons))
		for i, icon := range icons {
			iconIDs[i] = icon.Filename
		}

		configOption := models.NewConfigOption(ConfigName, iconIDs)
		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputSummary, ModuleInputEnvKey, models.TypeUserInput)
		variantOption := models.NewOption(VariantInputTitle, VariantInputSummary, VariantInputEnvKey, models.TypeOptionalUserInput)

		projectLocationOption.AddOption(relProjectRoot, moduleOption)
		moduleOption.AddOption("app", variantOption)
		variantOption.AddConfig("", configOption)
		foundOptions = true
	}
	if !foundOptions && lastErr != nil {
		return models.OptionNode{}, warnings, nil, lastErr
	}

	return *projectLocationOption, warnings, appIconsAllProjects, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputSummary, ProjectLocationInputEnvKey, models.TypeUserInput)
	moduleOption := models.NewOption(ModuleInputTitle, ModuleInputSummary, ModuleInputEnvKey, models.TypeUserInput)
	variantOption := models.NewOption(VariantInputTitle, VariantInputSummary, VariantInputEnvKey, models.TypeOptionalUserInput)
	configOption := models.NewConfigOption(DefaultConfigName, nil)

	projectLocationOption.AddOption("", moduleOption)
	moduleOption.AddOption("", variantOption)
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

// NewTemplate ...
func NewTemplate() *Template {
	return &Template{}
}

// Name ...
func (*Template) Name() string {
	return TemplateName
}

// ExcludedScannerNames ...
func (*Template) ExcludedScannerNames() []string {
	return nil
}

// DetectPlatform ...
func (template *Template) DetectPlatform(searchDir string) (_ bool, err error) {
	scanner := NewScanner()
	ok, err := scanner.DetectPlatform(searchDir)
	template.ProjectRoots = scanner.ProjectRoots
	template.SearchDir = scanner.SearchDir

	return ok, err
}

func (template *Template) collectAndroidProjects() ([]androidProject, error) {
	var (
		lastErr  error = nil
		projects       = []androidProject{}
	)

	for _, projectRoot := range template.ProjectRoots {
		var warnings models.Warnings

		exists, err := containsLocalProperties(projectRoot)
		if err != nil {
			lastErr = err
			continue
		}
		if exists {
			containsLocalPropertiesWarning := fmt.Sprintf("the local.properties file should NOT be checked into Version Control Systems, as it contains information specific to your local configuration, the location of the file is: %s", filepath.Join(projectRoot, "local.properties"))
			warnings = []string{containsLocalPropertiesWarning}
		}

		if err := checkGradlew(projectRoot); err != nil {
			lastErr = err
			continue
		}

		relProjectRoot, err := filepath.Rel(template.SearchDir, projectRoot)
		if err != nil {
			lastErr = err
			continue
		}

		icons, err := LookupIcons(projectRoot, template.SearchDir)
		if err != nil {
			analytics.LogInfo("android-icon-lookup", analytics.DetectorErrorData("android", err), "Failed to lookup android icon")
		}

		projects = append(projects, androidProject{
			projectRelPath: relProjectRoot,
			icons:          icons,
			warnings:       warnings,
		})
	}
	if len(projects) == 0 && lastErr != nil {
		return []androidProject{}, lastErr
	}

	return projects, nil
}

func (template *Template) Get() (builder.TemplateNode, error) {
	return template.getTemplate()
}

func (template *Template) GetManual() (builder.TemplateNode, error) {
	panic("not implemented")
}
