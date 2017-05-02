package cordova

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/xcode"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

const scannerName = "cordova"

const (
	configXMLBasePath = "config.xml"
	platformsDirName  = "platforms"
)

const (
	projectTypeInputTitle  = "Project type"
	projectTypeInputEnvKey = "_"
)

const (
	workDirInputKey    = "work_dir"
	workDirInputTitle  = "Working directory"
	workDirInputEnvKey = "CORDOVA_WORK_DIR"
)

// ConfigDescriptor ...
type ConfigDescriptor struct {
	iosConfigDescriptors     []xcode.ConfigDescriptor
	androidConfigDescriptors []android.ConfigDescriptor
}

// NewConfigDescriptor ...
func NewConfigDescriptor() ConfigDescriptor {
	return ConfigDescriptor{}
}

// ConfigName ...
func (descriptor ConfigDescriptor) ConfigName() string {
	return ""
}

// WidgetModel ...
type WidgetModel struct {
	ID       string `xml:"id,attr"`
	Version  string `xml:"version,attr"`
	XMLNS    string `xml:"xmlns,attr"`
	XMLNSCDV string `xml:"xmlns cdv,attr"`
}

// ProjectConfigModel ...
type ProjectConfigModel struct {
	pth    string
	widget WidgetModel
}

func parseConfigXMLContent(content string) (WidgetModel, error) {
	widget := WidgetModel{}
	if err := xml.Unmarshal([]byte(content), &widget); err != nil {
		return WidgetModel{}, err
	}
	return widget, nil
}

func parseConfigXML(pth string) (WidgetModel, error) {
	content, err := fileutil.ReadStringFromFile(pth)
	if err != nil {
		return WidgetModel{}, err
	}
	return parseConfigXMLContent(content)
}

func filterRootConfigXMLFile(fileList []string) (string, error) {
	allowConfigXMLBaseFilter := utility.BaseFilter(configXMLBasePath, true)
	configXMLs, err := utility.FilterPaths(fileList, allowConfigXMLBaseFilter)
	if err != nil {
		return "", err
	}

	if len(configXMLs) == 0 {
		return "", nil
	}

	return configXMLs[0], nil
}

// ConfigName ...
func ConfigName(hasIosProject, hasAndroidProject bool) string {
	configName := "cordova"
	if hasIosProject {
		configName += "-ios"
	}
	if hasAndroidProject {
		configName += "-android"
	}
	return configName + "-config"
}

// Scanner ...
type Scanner struct {
	projectConfig       ProjectConfigModel
	searchDir           string
	hasKarmaJasmineTest bool
	hasJasmineTest      bool
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (scanner Scanner) Name() string {
	return scannerName
}

func pathsEquals(pth1, pth2 string) (bool, error) {
	absPth1, err := pathutil.AbsPath(pth1)
	if err != nil {
		return false, err
	}

	absPth2, err := pathutil.AbsPath(pth2)
	if err != nil {
		return false, err
	}

	return (absPth1 == absPth2), nil
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return false, fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}

	// Search for config.xml file
	log.Infoft("Searching for config.xml file")

	configXMLPth, err := filterRootConfigXMLFile(fileList)
	if err != nil {
		return false, fmt.Errorf("failed to search for config.xml file, error: %s", err)
	}

	log.Printft("config.xml: %s", configXMLPth)

	if configXMLPth == "" {
		log.Printft("platform not detected")
		return false, nil
	}

	widget, err := parseConfigXML(configXMLPth)
	if err != nil {
		log.Printft("can not parse config.xml as a Cordova widget, error: %s", err)
		log.Printft("platform not detected")
		return false, nil
	}

	// ensure it is a cordova widget
	if !strings.Contains(widget.XMLNSCDV, "cordova.apache.org") {
		log.Printft("config.xml propert: xmlns:cdv does not contain cordova.apache.org")
		log.Printft("platform not detected")
		return false, nil
	}

	// ensure it is not an ionic project
	projectBaseDir := filepath.Dir(configXMLPth)

	if exist, err := pathutil.IsPathExists(filepath.Join(projectBaseDir, "ionic.project")); err != nil {
		return false, fmt.Errorf("failed to check if project is an ionic project, error: %s", err)
	} else if exist {
		log.Printft("ionic.project file found seems to be an ionic project")
		return false, nil
	}

	if exist, err := pathutil.IsPathExists(filepath.Join(projectBaseDir, "ionic.config.json")); err != nil {
		return false, fmt.Errorf("failed to check if project is an ionic project, error: %s", err)
	} else if exist {
		log.Printft("ionic.config.json file found seems to be an ionic project")
		return false, nil
	}

	log.Doneft("Platform detected")

	scanner.projectConfig = ProjectConfigModel{
		pth:    configXMLPth,
		widget: widget,
	}

	scanner.searchDir = searchDir

	return true, nil
}

// ExcludedScannerNames ...
func (scanner *Scanner) ExcludedScannerNames() []string {
	return []string{
		string(utility.XcodeProjectTypeIOS),
		string(utility.XcodeProjectTypeMacOS),
		android.ScannerName,
	}
}

func detectPlatforms(platformsDir string) ([]string, error) {
	platformsJSONPth := filepath.Join(platformsDir, "platforms.json")
	if exist, err := pathutil.IsPathExists(platformsJSONPth); err != nil {
		return []string{}, err
	} else if !exist {
		return []string{}, nil
	}

	bytes, err := fileutil.ReadBytesFromFile(platformsJSONPth)
	if err != nil {
		return []string{}, err
	}

	type PlatformsModel struct {
		IOS     string `json:"ios"`
		Android string `json:"android"`
	}

	var platformsModel PlatformsModel
	if err := json.Unmarshal(bytes, &platformsModel); err != nil {
		return []string{}, err
	}

	platforms := []string{}
	if platformsModel.IOS != "" {
		platforms = append(platforms, "ios")
	}
	if platformsModel.Android != "" {
		platforms = append(platforms, "android")
	}

	return platforms, nil
}

// PackagesModel ...
type PackagesModel struct {
	Platforms       []string          `json:"cordovaPlatforms"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func parsePackagesJSONContent(content string) (PackagesModel, error) {
	var packages PackagesModel
	if err := json.Unmarshal([]byte(content), &packages); err != nil {
		return PackagesModel{}, err
	}
	return packages, nil
}

func parsePackagesJSON(packagesJSONPth string) (PackagesModel, error) {
	content, err := fileutil.ReadStringFromFile(packagesJSONPth)
	if err != nil {
		return PackagesModel{}, err
	}
	return parsePackagesJSONContent(content)
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}
	projectRootDir := filepath.Dir(scanner.projectConfig.pth)

	packagesJSONPth := filepath.Join(projectRootDir, "package.json")
	packages, err := parsePackagesJSON(packagesJSONPth)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}

	// Search for karma/jasmine tests
	log.Printft("Searching for karma/jasmine test")

	karmaTestDetected := false

	karmaJasmineDependencyFound := false
	for dependency := range packages.Dependencies {
		if strings.Contains(dependency, "karma-jasmine") {
			karmaJasmineDependencyFound = true
		}
	}
	if !karmaJasmineDependencyFound {
		for dependency := range packages.DevDependencies {
			if strings.Contains(dependency, "karma-jasmine") {
				karmaJasmineDependencyFound = true
			}
		}
	}
	log.Printft("karma-jasmine dependency found: %v", karmaJasmineDependencyFound)

	if karmaJasmineDependencyFound {
		karmaConfigJSONPth := filepath.Join(projectRootDir, "karma.conf.js")
		if exist, err := pathutil.IsPathExists(karmaConfigJSONPth); err != nil {
			return models.OptionModel{}, warnings, err
		} else if exist {
			karmaTestDetected = true
		}
	}
	log.Printft("karma.conf.js found: %v", karmaTestDetected)

	scanner.hasKarmaJasmineTest = karmaTestDetected

	// Search for jasmine tests
	jasminTestDetected := false

	if !karmaTestDetected {
		log.Printft("Searching for jasmine test")

		jasmineDependencyFound := false
		for dependency := range packages.Dependencies {
			if strings.Contains(dependency, "jasmine") {
				jasmineDependencyFound = true
				break
			}
		}
		if !jasmineDependencyFound {
			for dependency := range packages.DevDependencies {
				if strings.Contains(dependency, "jasmine") {
					jasmineDependencyFound = true
					break
				}
			}
		}
		log.Printft("jasmine dependency found: %v", jasmineDependencyFound)

		if jasmineDependencyFound {
			jasmineConfigJSONPth := filepath.Join(projectRootDir, "spec", "support", "jasmine.json")
			if exist, err := pathutil.IsPathExists(jasmineConfigJSONPth); err != nil {
				return models.OptionModel{}, warnings, err
			} else if exist {
				jasminTestDetected = true
			}
		}

		log.Printft("jasmine.json found: %v", jasminTestDetected)

		scanner.hasJasmineTest = jasminTestDetected
	}

	projectTypeOption := models.NewOption(projectTypeInputTitle, projectTypeInputEnvKey)

	iosConfigOption := models.NewConfigOption(ConfigName(true, false))
	projectTypeOption.AddConfig("ios", iosConfigOption)

	androidConfigOption := models.NewConfigOption(ConfigName(true, true))
	projectTypeOption.AddConfig("android", androidConfigOption)

	iosAndroidConfigOption := models.NewConfigOption(ConfigName(true, true))
	projectTypeOption.AddConfig("ios+android", iosAndroidConfigOption)

	return *projectTypeOption, warnings, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{}
}

func relCordovaWorkDir(baseDir, cordovaConfigPth string) (string, error) {
	absBaseDir, err := pathutil.AbsPath(baseDir)
	if err != nil {
		return "", err
	}

	absCordovaConfigPth, err := pathutil.AbsPath(cordovaConfigPth)
	if err != nil {
		return "", err
	}

	absCordovaWorkDir := filepath.Dir(absCordovaConfigPth)
	if absBaseDir == absCordovaWorkDir {
		return "", nil
	}

	cordovaWorkdir, err := filepath.Rel(absBaseDir, absCordovaWorkDir)
	if err != nil {
		return "", err
	}

	return cordovaWorkdir, nil
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	workdir, err := relCordovaWorkDir(scanner.searchDir, scanner.projectConfig.pth)
	if err != nil {
		return models.BitriseConfigMap{}, fmt.Errorf("Failed to check if search dir is the work dir, error: %s", err)
	}

	platforms := []string{
		"ios",
		"android",
		"ios,android",
	}

	bitriseDataMap := models.BitriseConfigMap{}
	for _, platform := range platforms {
		iosProject := strings.Contains(platform, "ios")
		androidProject := strings.Contains(platform, "android")

		configBuilder := models.NewDefaultConfigBuilder()

		workdirEnvList := []envmanModels.EnvironmentItemModel{}
		if workdir != "" {
			workdirEnvList = append(workdirEnvList, envmanModels.EnvironmentItemModel{"workdir": "$" + workDirInputEnvKey})
		}

		if scanner.hasJasmineTest || scanner.hasKarmaJasmineTest {
			// CI
			if scanner.hasKarmaJasmineTest {
				configBuilder.AppendMainStepList(steps.KarmaJasmineTestRunnerStepListItem(workdirEnvList...))
			} else if scanner.hasJasmineTest {
				configBuilder.AppendMainStepList(steps.JasmineTestRunnerStepListItem(workdirEnvList...))
			}

			// CD
			configBuilder.AddDefaultWorkflowBuilder(models.DeployWorkflowID)

			if scanner.hasKarmaJasmineTest {
				configBuilder.AppendMainStepListTo(models.DeployWorkflowID, steps.KarmaJasmineTestRunnerStepListItem(workdirEnvList...))
			} else if scanner.hasJasmineTest {
				configBuilder.AppendMainStepListTo(models.DeployWorkflowID, steps.JasmineTestRunnerStepListItem(workdirEnvList...))
			}

			configBuilder.AppendMainStepListTo(models.DeployWorkflowID, steps.GenerateCordovaBuildConfigStepListItem())

			cordovaArchiveEnvs := []envmanModels.EnvironmentItemModel{
				envmanModels.EnvironmentItemModel{"platform": platform},
				envmanModels.EnvironmentItemModel{"target": "emulator"},
			}
			if workdir != "" {
				cordovaArchiveEnvs = append(cordovaArchiveEnvs, envmanModels.EnvironmentItemModel{"workdir": "$" + workDirInputEnvKey})
			}
			configBuilder.AppendMainStepListTo(models.DeployWorkflowID, steps.CordovaArchiveStepListItem(cordovaArchiveEnvs...))

			appEnvs := []envmanModels.EnvironmentItemModel{}
			if workdir != "" {
				appEnvs = append(appEnvs, envmanModels.EnvironmentItemModel{workDirInputEnvKey: workdir})
			}
			config, err := configBuilder.Generate(appEnvs...)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			data, err := yaml.Marshal(config)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			bitriseDataMap[ConfigName(iosProject, androidProject)] = string(data)
		} else {
			configBuilder.AppendMainStepList(steps.GenerateCordovaBuildConfigStepListItem())

			cordovaArchiveEnvs := []envmanModels.EnvironmentItemModel{
				envmanModels.EnvironmentItemModel{"platform": platform},
				envmanModels.EnvironmentItemModel{"target": "emulator"},
			}
			if workdir != "" {
				cordovaArchiveEnvs = append(cordovaArchiveEnvs, envmanModels.EnvironmentItemModel{"workdir": "$" + workDirInputEnvKey})
			}
			configBuilder.AppendMainStepList(steps.CordovaArchiveStepListItem(
				envmanModels.EnvironmentItemModel{"platform": platform},
				envmanModels.EnvironmentItemModel{"target": "emulator"},
			))

			appEnvs := []envmanModels.EnvironmentItemModel{}
			if workdir != "" {
				appEnvs = append(appEnvs, envmanModels.EnvironmentItemModel{workDirInputEnvKey: workdir})
			}
			config, err := configBuilder.Generate(appEnvs...)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			data, err := yaml.Marshal(config)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			bitriseDataMap[ConfigName(iosProject, androidProject)] = string(data)
		}
	}

	return bitriseDataMap, nil
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}
