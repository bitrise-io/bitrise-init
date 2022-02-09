package android

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/analytics"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/log"
)

// Scanner ...
type Scanner struct {
	Projects     []project
	ProjectRoots []string

	ExcludeTest    bool
	ExcludeAppIcon bool
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
	detected, projects, projectRoots, err := detect(searchDir)
	scanner.ProjectRoots = projectRoots
	scanner.Projects = projects

	return detected, err
}

func detect(searchDir string) (bool, []project, []string, error) {
	projectFiles := fileGroups{
		{"build.gradle", "build.gradle.kts"},
		{"settings.gradle", "settings.gradle.kts"},
	}
	skipDirs := []string{".git", "CordovaLib", "node_modules"}

	log.TInfof("Searching for android files")

	projectRoots, err := walkMultipleFileGroups(searchDir, projectFiles, skipDirs)
	if err != nil {
		return false, []project{}, []string{}, fmt.Errorf("failed to search for build.gradle files, error: %s", err)
	}

	log.TSuccessf("Platform detected")

	for _, file := range projectFiles {
		log.TPrintf("- %s", file)
	}

	log.TPrintf("%d android files detected", len(projectRoots))

	if len(projectRoots) == 0 {
		return false, []project{}, []string{}, err
	}

	projects, err := parseProjects(searchDir, projectRoots)

	return true, projects, projectRoots, err
}

func parseProjects(searchDir string, projectRoots []string) ([]project, error) {
	var (
		lastErr  error = nil
		projects       = []project{}
	)

	for _, projectRoot := range projectRoots {
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

		relProjectRoot, err := filepath.Rel(searchDir, projectRoot)
		if err != nil {
			lastErr = err
			continue
		}

		icons, err := LookupIcons(projectRoot, searchDir)
		if err != nil {
			analytics.LogInfo("android-icon-lookup", analytics.DetectorErrorData("android", err), "Failed to lookup android icon")
		}

		projects = append(projects, project{
			projectRelPath: relProjectRoot,
			icons:          icons,
			warnings:       warnings,
		})
	}

	if len(projects) == 0 && lastErr != nil {
		return []project{}, lastErr
	}

	return projects, nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputSummary, ProjectLocationInputEnvKey, models.TypeSelector)
	warnings := models.Warnings{}
	appIconsAllProjects := models.Icons{}

	for _, project := range scanner.Projects {
		warnings = append(warnings, project.warnings...)

		appIconsAllProjects = append(appIconsAllProjects, project.icons...)
		iconIDs := make([]string, len(project.icons))
		for i, icon := range project.icons {
			iconIDs[i] = icon.Filename
		}

		configOption := models.NewConfigOption(ConfigName, iconIDs)
		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputSummary, ModuleInputEnvKey, models.TypeUserInput)
		variantOption := models.NewOption(VariantInputTitle, VariantInputSummary, VariantInputEnvKey, models.TypeOptionalUserInput)

		projectLocationOption.AddOption(project.projectRelPath, moduleOption)
		moduleOption.AddOption("app", variantOption)
		variantOption.AddConfig("", configOption)
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
func (scanner *Scanner) Configs(isPrivateRepository bool) (models.BitriseConfigMap, error) {
	configBuilder := scanner.generateConfigBuilder(isPrivateRepository)

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
