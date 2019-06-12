package android

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/models"
)

// Scanner ...
type Scanner struct {
	SearchDir      string
	ProjectRoots   []string
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
			return models.OptionNode{}, warnings, nil, err
		}

		relProjectRoot, err := filepath.Rel(scanner.SearchDir, projectRoot)
		if err != nil {
			return models.OptionNode{}, warnings, nil, err
		}

		icons, err := LookupIcons(projectRoot, scanner.SearchDir)
		if err != nil {
			return models.OptionNode{}, warnings, nil, err
		}
		appIconsAllProjects = append(appIconsAllProjects, icons...)
		iconIDs := make([]string, len(icons))
		for i, icon := range icons {
			iconIDs[i] = icon.Filename
		}

		configOption := models.NewConfigOption(ConfigName, iconIDs)
		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey).SetType(models.TypeUserInput)
		variantOption := models.NewOption(VariantInputTitle, VariantInputEnvKey).SetType(models.TypeOptionalUserInput)

		projectLocationOption.AddOption(relProjectRoot, moduleOption)
		moduleOption.AddOption("app", variantOption)
		variantOption.AddConfig("", configOption)
	}

	return *projectLocationOption, warnings, appIconsAllProjects, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	// projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputEnvKey).SetType(models.TypeUserInput)
	// moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey).SetType(models.TypeUserInput)
	// variantOption := models.NewOption(VariantInputTitle, VariantInputEnvKey).SetType(models.TypeOptionalUserInput)
	// configOption := models.NewConfigOption(DefaultConfigName, nil)

	// projectLocationOption.AddOption("", moduleOption)
	// moduleOption.AddOption("", variantOption)
	// variantOption.AddConfig("", configOption)

	// return *projectLocationOption

	// testing options

	/*

		The section below is for client side testing only, should be removed for production use.

	*/

	q1 := models.NewOption("title1-TypeUserInput", "envkey1").SetType(models.TypeUserInput)

	q2 := models.NewOption("title2-TypeOptionalUserInput", "envkey2").SetType(models.TypeOptionalUserInput)

	q3 := models.NewOption("title3-TypeUserInput-default-value", "envkey3").SetType(models.TypeUserInput)

	q4 := models.NewOption("title4-TypeOptionalUserInput-default-value", "envkey4").SetType(models.TypeOptionalUserInput)

	q5 := models.NewOption("title5-TypeSelectore-one-item", "envkey5").SetType(models.TypeSelector)

	q6 := models.NewOption("title6-TypeOptionalSelector-one-item", "envkey6").SetType(models.TypeOptionalSelector)

	q7 := models.NewOption("title7-TypeSelector-3-item", "envkey7").SetType(models.TypeSelector)
	q8 := models.NewOption("title8-TypeSelector-3-item", "envkey8").SetType(models.TypeSelector)
	q9 := models.NewOption("title9-TypeSelector-3-item", "envkey9").SetType(models.TypeSelector)

	q10 := models.NewOption("title10-TypeOptionalSelector-3-item", "envkey10").SetType(models.TypeOptionalSelector)
	q11 := models.NewOption("title11-TypeOptionalSelector-3-item", "envkey11").SetType(models.TypeOptionalSelector)
	q12 := models.NewOption("title12-TypeOptionalSelector-3-item", "envkey12").SetType(models.TypeOptionalSelector)

	q10a := models.NewOption("title10a-TypeOptionalSelector-3-item", "envkey10a").SetType(models.TypeOptionalSelector)
	q11a := models.NewOption("title11a-TypeOptionalSelector-3-item", "envkey11a").SetType(models.TypeOptionalSelector)
	q12a := models.NewOption("title12a-TypeOptionalSelector-3-item", "envkey12a").SetType(models.TypeOptionalSelector)

	q10b := models.NewOption("title10b-TypeOptionalSelector-3-item", "envkey10b").SetType(models.TypeOptionalSelector)
	q11b := models.NewOption("title11b-TypeOptionalSelector-3-item", "envkey11b").SetType(models.TypeOptionalSelector)
	q12b := models.NewOption("title12b-TypeOptionalSelector-3-item", "envkey12b").SetType(models.TypeOptionalSelector)

	q1.AddOption("", q2)
	q2.AddOption("", q3)
	q3.AddOption("my default value", q4)
	q4.AddOption("my default value", q5)
	q5.AddOption("my only option", q6)
	q6.AddOption("my only option", q7)

	//
	q10.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q10.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q10.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q7.AddOption("first option", q10)

	q11.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q11.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q11.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q7.AddOption("second option", q11)

	q12.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q12.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q12.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q7.AddOption("third option", q12)

	//
	q10a.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q10a.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q10a.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q8.AddOption("first option", q10a)

	q11a.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q11a.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q11a.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q8.AddOption("second option", q11a)

	q12a.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q12a.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q12a.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q8.AddOption("third option", q12a)

	//
	q10b.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q10b.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q10b.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q9.AddOption("first option", q10b)

	q11b.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q11b.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q11b.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q9.AddOption("second option", q11b)

	q12b.AddConfig("first option", models.NewConfigOption("myconf", nil))
	q12b.AddConfig("second option", models.NewConfigOption("myconf", nil))
	q12b.AddConfig("third option", models.NewConfigOption("myconf", nil))
	q9.AddOption("third option", q12b)

	return *q1
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
