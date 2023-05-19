package flutter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android"
	"github.com/bitrise-io/bitrise-init/scanners/flutter/flutterproject"
	"github.com/bitrise-io/bitrise-init/scanners/flutter/tracker"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	v2log "github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/pathfilters"
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
	installerUpdateFlutterKey   = "is_update"
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
	tracker  tracker.FlutterTracker
}

type project struct {
	id   int
	path string
	// TODO: merge project and flutterproject.Project
	flutterProject    flutterproject.Project
	hasTest           bool
	hasIosProject     bool
	hasAndroidProject bool
}

func (proj project) platform() string {
	switch {
	case proj.hasAndroidProject && !proj.hasIosProject:
		return "android"
	case !proj.hasAndroidProject && proj.hasIosProject:
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
	return &Scanner{
		tracker: tracker.NewFlutterTracker(v2log.NewLogger()),
	}
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
		var proj project

		pubspecPath := filepath.Join(projectLocation, "pubspec.yaml")
		pubspecFile, err := os.Open(pubspecPath)
		if err != nil {
			log.TErrorf("Failed to open pubspec.yaml file at: %s, error: %s", pubspecPath, err)
			return false, err
		}

		var ps pubspec
		if err := yaml.NewDecoder(pubspecFile).Decode(&ps); err != nil {
			log.TErrorf("Failed to decode yaml pubspec.yaml file at: %s, error: %s", pubspecPath, err)
			return false, err
		}

		testsDirPath := filepath.Join(projectLocation, "test")
		if exists, err := pathutil.IsDirExists(testsDirPath); err == nil && exists {
			if files, err := ioutil.ReadDir(testsDirPath); err == nil && len(files) > 0 {
				for _, file := range files {
					if strings.HasSuffix(file.Name(), "_test.dart") {
						proj.hasTest = true
						break
					}
				}
			}
		}

		iosProjPath := filepath.Join(projectLocation, "ios", "Runner.xcworkspace")
		if exists, err := pathutil.IsPathExists(iosProjPath); err == nil && exists {
			proj.hasIosProject = true
		}

		androidProjPath := filepath.Join(projectLocation, "android", "build.gradle")
		if exists, err := pathutil.IsPathExists(androidProjPath); err == nil && exists {
			proj.hasAndroidProject = true
		}

		if !proj.hasAndroidProject {
			androidProjPath := filepath.Join(projectLocation, "android", "build.gradle.kts")
			if exists, err := pathutil.IsPathExists(androidProjPath); err == nil && exists {
				proj.hasAndroidProject = true
			}
		}

		log.TPrintf("- Project name: %s", ps.Name)
		log.TPrintf("  Path: %s", projectLocation)
		log.TPrintf("  HasTest: %t", proj.hasTest)
		log.TPrintf("  HasAndroidProject: %t", proj.hasAndroidProject)
		log.TPrintf("  HasIosProject: %t", proj.hasIosProject)

		proj.path = projectLocation

		currentID++
		proj.id = currentID
		scanner.projects = append(scanner.projects, proj)
	}

	for i, project := range scanner.projects {
		proj, err := flutterproject.New(project.path, flutterproject.NewFileOpener())
		if err == nil {
			sdkVersions := proj.FlutterAndDartSDKVersions()
			scanner.projects[i].flutterProject = *proj
			scanner.tracker.LogSDKVersions(sdkVersions)
		}
	}
	scanner.tracker.Wait()

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
		flutterProjectLocationOption.AddOption(project.path, configOption)
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
	version := proj.flutterProject.FlutterSDKVersionToUse()
	flutterInstallStep := steps.FlutterInstallStepListItem(version, false)
	deploySteps := steps.DefaultDeployStepList()

	// primary
	configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, primaryWorkflowDescription)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, prepareSteps...)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, flutterInstallStep)

	// restore cache is after flutter-installer, to prevent removal of pub system cache
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.RestoreDartCache())

	if proj.hasTest {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterTestStepListItem(
			envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
		))
	}

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.SaveDartCache())

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, deploySteps...)

	// deploy
	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, prepareSteps...)

	if proj.hasIosProject {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	}

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, flutterInstallStep)

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterAnalyzeStepListItem(
		envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
	))

	if proj.hasTest {
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
	if proj.hasTest {
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
