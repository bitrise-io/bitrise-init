package flutter

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/v2/fileutil"
	pathutilv2 "github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-io/go-xcode/pathfilters"
	"github.com/godrei/go-flutter/flutterproject"
	"gopkg.in/yaml.v2"
)

const (
	scannerName                 = "flutter"
	configName                  = "flutter-config"
	projectLocationInputKey     = "project_location"
	projectLocationInputEnvKey  = "BITRISE_FLUTTER_PROJECT_LOCATION"
	projectLocationInputTitle   = "Project location"
	projectLocationInputSummary = "The path to your Flutter project, stored as an Environment Variable. In your Workflows, you can specify paths relative to this path. You can change this at any time."
	platformInputKey            = "platform"
	platformInputTitle          = "Platform"
	platformInputSummary        = "The target platform for your first build. Your options are iOS, Android, both, or neither. You can change this in your Env Vars at any time."
	iosOutputTypeKey            = "ios_output_type"
	iosOutputTypeArchive        = "archive"
)

var (
	platforms = []string{
		"android",
		"ios",
		"both",
	}
)

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	projects []project
}

type project struct {
	flutterproject.Project
	id                  int
	flutterVersionToUse string
}

func (proj project) platform() string {
	switch {
	case proj.AndroidProjectPth() != "" && proj.IOSProjectPth() != "":
		return "android"
	case proj.AndroidProjectPth() == "" && proj.IOSProjectPth() != "":
		return "ios"
	default:
		return "both"
	}
}

type pubspec struct {
	Name string `yaml:"name"`
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (scanner *Scanner) Name() string {
	return scannerName
}

func findProjectLocations(searchDir string) ([]string, error) {
	fileList, err := pathutil.ListPathInDirSortedByComponents(searchDir, true)
	if err != nil {
		return nil, err
	}

	filters := []pathutil.FilterFunc{
		pathutil.BaseFilter("pubspec.yaml", true),
		pathutil.ComponentFilter("node_modules", false),
	}

	paths, err := pathutil.FilterPaths(fileList, filters...)
	if err != nil {
		return nil, err
	}

	for i, path := range paths {
		paths[i] = filepath.Dir(path)
	}

	return paths, nil
}

func findWorkspaceLocations(projectLocation string) ([]string, error) {
	fileList, err := pathutil.ListPathInDirSortedByComponents(projectLocation, true)
	if err != nil {
		return nil, err
	}

	for i, file := range fileList {
		fileList[i] = filepath.Join(projectLocation, file)
	}

	filters := []pathutil.FilterFunc{
		pathfilters.AllowXCWorkspaceExtFilter,
		pathfilters.AllowIsDirectoryFilter,
		pathfilters.ForbidEmbeddedWorkspaceRegexpFilter,
		pathfilters.ForbidGitDirComponentFilter,
		pathfilters.ForbidPodsDirComponentFilter,
		pathfilters.ForbidCarthageDirComponentFilter,
		pathfilters.ForbidFramworkComponentWithExtensionFilter,
		pathfilters.ForbidCordovaLibDirComponentFilter,
		pathfilters.ForbidNodeModulesComponentFilter,
	}

	return pathutil.FilterPaths(fileList, filters...)
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	log.TInfof("Search for project(s)")
	projectLocations, err := findProjectLocations(searchDir)
	if err != nil {
		return false, err
	}

	log.TPrintf("Paths containing pubspec.yaml(%d):", len(projectLocations))
	for _, p := range projectLocations {
		log.TPrintf("- %s", p)
	}
	log.TPrintf("")

	log.TInfof("Fetching pubspec.yaml files")

	currentID := -1
	for _, projectLocation := range projectLocations {
		flutterProj, err := flutterproject.New(projectLocation, fileutil.NewFileManager(), pathutilv2.NewPathChecker())
		if err != nil {
			log.TErrorf(err.Error())
			continue
		}

		flutterVersion, err := flutterProj.FlutterSDKVersionToUse()
		if err != nil {
			log.Warnf(err.Error())
		}

		currentID++
		proj := project{
			Project:             *flutterProj,
			id:                  currentID,
			flutterVersionToUse: flutterVersion,
		}

		scanner.projects = append(scanner.projects, proj)

		log.TPrintf("- Project name: %s", proj.Pubspec().Name)
		log.TPrintf("  Path: %s", proj.RootDir())
		log.TPrintf("  HasTest: %s", proj.TestDirPth() != "")
		log.TPrintf("  HasAndroidProject: %s", proj.AndroidProjectPth() != "")
		log.TPrintf("  HasIosProject: %s", proj.IOSProjectPth() != "")
		if flutterVersion != "" {
			log.TPrintf("  Flutter version to use: %s", proj.flutterVersionToUse)
		}
	}

	return len(scanner.projects) > 0, nil
}

// ExcludedScannerNames ...
func (scanner *Scanner) ExcludedScannerNames() []string {
	return []string{
		string(ios.XcodeProjectTypeIOS),
		android.ScannerName,
	}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	flutterProjectLocationOption := models.NewOption(projectLocationInputTitle, projectLocationInputSummary, projectLocationInputEnvKey, models.TypeSelector)

	for _, project := range scanner.projects {
		configOption := models.NewConfigOption(configNameForProject(project), nil)
		flutterProjectLocationOption.AddOption(project.RootDir(), configOption)
	}

	return *flutterProjectLocationOption, models.Warnings{}, nil, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	flutterProjectLocationOption := models.NewOption(projectLocationInputTitle, projectLocationInputSummary, projectLocationInputEnvKey, models.TypeUserInput)

	cfg := configName + "-test"

	flutterPlatformOption := models.NewOption(platformInputTitle, platformInputSummary, "", models.TypeSelector)
	flutterProjectLocationOption.AddOption("", flutterPlatformOption)

	for _, platform := range platforms {
		configOption := models.NewConfigOption(cfg+"-app-"+platform, nil)
		flutterPlatformOption.AddConfig(platform, configOption)
	}

	return *flutterProjectLocationOption
}

func (scanner *Scanner) Configs(repoAccess models.RepoAccess) (models.BitriseConfigMap, error) {
	configs := models.BitriseConfigMap{}

	for _, proj := range scanner.projects {
		config, err := generateConfig(repoAccess, proj)
		if err != nil {
			return nil, err
		}

		configs[configNameForProject(proj)] = config
	}

	return configs, nil
}

func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return generateDefaultConfigMap(models.RepoAccessUnknown)
}

func generateConfig(repoAccess models.RepoAccess, proj project) (string, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	// Common steps to all workflows
	prepareSteps := steps.DefaultPrepareStepList(steps.PrepareListParams{RepoAccess: repoAccess})
	flutterInstallStep := steps.FlutterInstallStepListItem(proj.flutterVersionToUse, false)
	deploySteps := steps.DefaultDeployStepList()

	// primary
	configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, primaryWorkflowDescription)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, prepareSteps...)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, flutterInstallStep)

	// restore cache is after flutter-installer, to prevent removal of pub system cache
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.RestoreDartCache())

	if proj.TestDirPth() != "" {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterTestStepListItem(
			envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
		))
	}

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.SaveDartCache())

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, deploySteps...)

	// deploy
	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, prepareSteps...)

	if proj.IOSProjectPth() != "" {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	}

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, flutterInstallStep)

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterAnalyzeStepListItem(
		envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
	))

	if proj.TestDirPth() != "" {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterTestStepListItem(
			envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
		))
	}

	flutterBuildInputs := []envmanModels.EnvironmentItemModel{
		{projectLocationInputKey: "$" + projectLocationInputEnvKey},
		{platformInputKey: proj.platform()},
	}
	if proj.platform() != "android" {
		flutterBuildInputs = append(flutterBuildInputs, envmanModels.EnvironmentItemModel{iosOutputTypeKey: iosOutputTypeArchive})
	}
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterBuildStepListItem(flutterBuildInputs...))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, deploySteps...)

	config, err := configBuilder.Generate(scannerName)
	if err != nil {
		return "", err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// TODO: merge with generateConfig
func generateDefaultConfigMap(repoAccess models.RepoAccess) (models.BitriseConfigMap, error) {
	configs := models.BitriseConfigMap{}

	for _, variant := range []struct {
		configID string
		test     bool
		platform string
	}{
		{test: false, platform: "both", configID: configName + "-notest-app-both"},
		{test: true, platform: "both", configID: configName + "-test-app-both"},
		{test: false, platform: "android", configID: configName + "-notest-app-android"},
		{test: true, platform: "android", configID: configName + "-test-app-android"},
		{test: false, platform: "ios", configID: configName + "-notest-app-ios"},
		{test: true, platform: "ios", configID: configName + "-test-app-ios"},
	} {
		configBuilder := models.NewDefaultConfigBuilder()

		// Common steps to all workflows
		prepareSteps := steps.DefaultPrepareStepList(steps.PrepareListParams{RepoAccess: repoAccess})
		flutterInstallStep := steps.FlutterInstallStepListItem("", false)
		deploySteps := steps.DefaultDeployStepList()

		// primary
		configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, primaryWorkflowDescription)

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, prepareSteps...)

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, flutterInstallStep)

		// restore cache is after flutter-installer, to prevent removal of pub system cache
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.RestoreDartCache())

		if variant.test {
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterTestStepListItem(
				envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
			))
		}

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.SaveDartCache())

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, deploySteps...)

		// deploy
		configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, prepareSteps...)

		if variant.platform != "android" {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
		}

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, flutterInstallStep)

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterAnalyzeStepListItem(
			envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
		))

		if variant.test {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterTestStepListItem(
				envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
			))
		}

		flutterBuildInputs := []envmanModels.EnvironmentItemModel{
			{projectLocationInputKey: "$" + projectLocationInputEnvKey},
			{platformInputKey: variant.platform},
		}
		if variant.platform != "android" {
			flutterBuildInputs = append(flutterBuildInputs, envmanModels.EnvironmentItemModel{iosOutputTypeKey: iosOutputTypeArchive})
		}
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterBuildStepListItem(flutterBuildInputs...))

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, deploySteps...)

		config, err := configBuilder.Generate(scannerName)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configs[variant.configID] = string(data)
	}

	return configs, nil
}

func configNameForProject(proj project) string {
	name := configName
	if proj.TestDirPth() != "" {
		name += "-test"
	} else {
		name += "-notest"
	}

	switch proj.platform() {
	case "android":
		name += "-android"
	case "ios":
		name += "-ios"
	default:
		name += "-both"
	}

	name += fmt.Sprintf("-%d", proj.id)

	return name
}
