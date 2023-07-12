package android

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/analytics"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/log"
)

const (
	gradleKotlinBuildFile    = "build.gradle.kts"
	gradleKotlinSettingsFile = "settings.gradle.kts"
)

// Scanner ...
type Scanner struct {
	Projects []Project
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
	projects, err := detect(searchDir)
	scanner.Projects = projects

	detected := len(projects) > 0
	return detected, err
}

func detect(searchDir string) ([]Project, error) {
	projectFiles := fileGroups{
		{"build.gradle", gradleKotlinBuildFile},
		{"settings.gradle", gradleKotlinSettingsFile},
	}
	skipDirs := []string{".git", "CordovaLib", "node_modules"}

	log.TInfof("Searching for android files")

	projectRoots, err := walkMultipleFileGroups(searchDir, projectFiles, skipDirs)
	if err != nil {
		return nil, fmt.Errorf("failed to search for build.gradle files, error: %s", err)
	}

	log.TPrintf("%d android files detected", len(projectRoots))
	for _, file := range projectRoots {
		log.TPrintf("- %s", file)
	}

	if len(projectRoots) == 0 {
		return nil, nil
	}
	log.TSuccessf("Platform detected")

	projects, err := parseProjects(searchDir, projectRoots)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func parseProjects(searchDir string, projectRoots []string) ([]Project, error) {
	var (
		lastErr  error
		projects []Project
	)

	for _, projectRoot := range projectRoots {
		var warnings models.Warnings

		log.TInfof("Investigating Android project: %s", projectRoot)

		exists, err := containsLocalProperties(projectRoot)
		if err != nil {
			lastErr = err
			log.TWarnf("%s", err)

			continue
		}
		if exists {
			containsLocalPropertiesWarning := fmt.Sprintf("the local.properties file should NOT be checked into Version Control Systems, as it contains information specific to your local configuration, the location of the file is: %s", filepath.Join(projectRoot, "local.properties"))
			warnings = []string{containsLocalPropertiesWarning}
		}

		if err := checkGradlew(projectRoot); err != nil {
			lastErr = err
			log.TWarnf("%s", err)

			continue
		}

		relProjectRoot, err := filepath.Rel(searchDir, projectRoot)
		if err != nil {
			lastErr = err
			log.TWarnf("%s", err)

			continue
		}

		icons, err := LookupIcons(projectRoot, searchDir)
		if err != nil {
			analytics.LogInfo("android-icon-lookup", analytics.DetectorErrorData("android", err), "Failed to lookup android icon")
		}

		kotlinBuildScriptBased := usesKotlinBuildScripts(projectRoot)
		projects = append(projects, Project{
			RelPath:               relProjectRoot,
			UsesKotlinBuildScript: kotlinBuildScriptBased,
			Icons:                 icons,
			Warnings:              warnings,
		})
	}

	if len(projects) == 0 {
		return []Project{}, lastErr
	}

	return projects, nil
}

func usesKotlinBuildScripts(projectRoot string) bool {
	return fileExists(filepath.Join(projectRoot, gradleKotlinBuildFile)) && fileExists(filepath.Join(projectRoot, gradleKotlinSettingsFile))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputSummary, ProjectLocationInputEnvKey, models.TypeSelector)
	warnings := models.Warnings{}
	appIconsAllProjects := models.Icons{}

	for _, project := range scanner.Projects {
		warnings = append(warnings, project.Warnings...)
		appIconsAllProjects = append(appIconsAllProjects, project.Icons...)

		iconIDs := make([]string, len(project.Icons))
		for i, icon := range project.Icons {
			iconIDs[i] = icon.Filename
		}

		name := ConfigName
		if project.UsesKotlinBuildScript {
			name = ConfigNameKotlinScript
		}
		configOption := models.NewConfigOption(name, iconIDs)
		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputSummary, ModuleInputEnvKey, models.TypeUserInput)
		variantOption := models.NewOption(VariantInputTitle, VariantInputSummary, VariantInputEnvKey, models.TypeOptionalUserInput)

		projectLocationOption.AddOption(project.RelPath, moduleOption)
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

	projectLocationOption.AddOption(models.UserInputOptionDefaultValue, moduleOption)
	moduleOption.AddOption(models.UserInputOptionDefaultValue, variantOption)
	variantOption.AddConfig("", configOption)

	return *projectLocationOption
}

// Configs ...
func (scanner *Scanner) Configs(repoAccess models.RepoAccess) (models.BitriseConfigMap, error) {
	bitriseDataMap := models.BitriseConfigMap{}
	params := configBuildingParameters(scanner.Projects)
	for _, param := range params {
		configBuilder := scanner.generateConfigBuilder(repoAccess, param.useKotlinScript)

		config, err := configBuilder.Generate(ScannerName)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		bitriseDataMap[param.name] = string(data)
	}

	return bitriseDataMap, nil
}

type configBuildingParams struct {
	name            string
	useKotlinScript bool
}

func configBuildingParameters(projects []Project) []configBuildingParams {
	regularProjectCount := 0
	kotlinBuildScriptProjectCount := 0

	for _, project := range projects {
		if project.UsesKotlinBuildScript {
			kotlinBuildScriptProjectCount += 1
		} else {
			regularProjectCount += 1
		}
	}

	var params []configBuildingParams
	if 0 < regularProjectCount {
		params = append(params, configBuildingParams{
			name:            ConfigName,
			useKotlinScript: false,
		})
	}
	if 0 < kotlinBuildScriptProjectCount {
		params = append(params, configBuildingParams{
			name:            ConfigNameKotlinScript,
			useKotlinScript: true,
		})
	}
	return params
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := scanner.generateConfigBuilder(models.RepoAccessUnknown, false)

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
