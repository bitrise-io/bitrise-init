package expo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/cordova"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/scanners/reactnative"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/xcode-project/serialized"
	yaml "gopkg.in/yaml.v2"
)

const (
	configName        = "react-native-expo-config"
	defaultConfigName = "default-" + configName
)

// Name ...
const Name = "react-native-expo"

// Scanner ...
type Scanner struct {
	searchDir      string
	packageJSONPth string
	usesExpoKit    bool
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return Name
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.searchDir = searchDir

	log.TInfof("Collect package.json files")

	packageJSONPths, err := reactnative.CollectPackageJSONFiles(searchDir)
	if err != nil {
		return false, err
	}

	if len(packageJSONPths) == 0 {
		return false, nil
	}

	log.TPrintf("%d package.json file detected", len(packageJSONPths))
	for _, pth := range packageJSONPths {
		log.TPrintf("- %s", pth)
	}
	log.TPrintf("")

	log.TInfof("Filter package.json files with expo dependency")

	relevantPackageJSONPths := []string{}
	for _, packageJSONPth := range packageJSONPths {
		packages, err := cordova.ParsePackagesJSON(packageJSONPth)
		if err != nil {
			log.Warnf("Failed to parse package json file: %s, skipping...", packageJSONPth)
			continue
		}

		_, found := packages.Dependencies["expo"]
		if !found {
			continue
		}

		// app.json file is a required part of react native projects and it exists next to the root package.json file
		appJSONPth := filepath.Join(filepath.Dir(packageJSONPth), "app.json")
		if exist, err := pathutil.IsPathExists(appJSONPth); err != nil {
			log.Warnf("Failed to check if app.json file exist at: %s, skipping package json file: %s, error: %s", appJSONPth, packageJSONPth, err)
			continue
		} else if !exist {
			log.Warnf("No app.json file exist at: %s, skipping package json file: %s", appJSONPth, packageJSONPth)
			continue
		}

		relevantPackageJSONPths = append(relevantPackageJSONPths, packageJSONPth)
	}

	log.TPrintf("%d package.json file detected with expo dependency", len(relevantPackageJSONPths))
	for _, pth := range relevantPackageJSONPths {
		log.TPrintf("- %s", pth)
	}
	log.TPrintf("")

	if len(relevantPackageJSONPths) == 0 {
		return false, nil
	} else if len(relevantPackageJSONPths) > 1 {
		log.TWarnf("Multiple package.json file found, using: %s\n", relevantPackageJSONPths[0])
	}

	scanner.packageJSONPth = relevantPackageJSONPths[0]
	return true, nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}

	// we need to know if the project uses the Expo Kit,
	// since its usage differentiates the eject process
	// and the config options
	usesExpoKit := false

	fileList, err := utility.ListPathInDirSortedByComponents(scanner.searchDir, true)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}

	filters := []utility.FilterFunc{
		utility.ExtensionFilter(".js", true),
		utility.ComponentFilter("node_modules", false),
	}
	sourceFiles, err := utility.FilterPaths(fileList, filters...)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}

	re := regexp.MustCompile(`import .* from 'expo'`)

SOURCE_FILE_LOOP:
	for _, sourceFile := range sourceFiles {
		f, err := os.Open(sourceFile)
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if match := re.FindString(scanner.Text()); match != "" {
				usesExpoKit = true
				break SOURCE_FILE_LOOP
			}
		}
		if err := scanner.Err(); err != nil {
			return models.OptionModel{}, warnings, err
		}
	}

	scanner.usesExpoKit = usesExpoKit
	log.TPrintf("Uses ExpoKit: %v", usesExpoKit)

	// ensure app.json contains the required information (for non interactive eject)
	// and determine the ejected project name
	var projectName string

	rootDir := filepath.Dir(scanner.packageJSONPth)
	appJSONPth := filepath.Join(rootDir, "app.json")
	appJSON, err := fileutil.ReadStringFromFile(appJSONPth)
	if err != nil {
		return models.OptionModel{}, warnings, err
	}
	var app serialized.Object
	if err := json.Unmarshal([]byte(appJSON), &app); err != nil {
		return models.OptionModel{}, warnings, err
	}

	if usesExpoKit {
		// if the project uses Expo Kit app.json needs to contain expo/ios/bundleIdentifier and expo/android/package entries
		// to be able to eject in non interactive mode
		expoObj, err := app.Object("expo")
		if err != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain expo key", appJSONPth)
		}
		projectName, err = expoObj.String("name")
		if err != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain expo/name key", appJSONPth)
		}

		iosObj, err := expoObj.Object("ios")
		if err != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain expo/ios key", appJSONPth)
		}
		bundleID, err := iosObj.String("bundleIdentifier")
		if err != nil || bundleID == "" {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain expo/ios/bundleIdentifier key or its value is empty", appJSONPth)
		}

		androidObj, err := expoObj.Object("android")
		if err != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain expo/android key", appJSONPth)
		}
		packageName, err := androidObj.String("package")
		if err != nil || packageName == "" {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain expo/android/package key or its value is empty", appJSONPth)
		}
	} else {
		// if the project does not use Expo Kit app.json needs to contain name and displayName entries
		// to be able to eject in non interactive mode
		projectName, err = app.String("name")
		if err != nil || projectName == "" {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain name key", appJSONPth)
		}
		displayName, err := app.String("displayName")
		if err != nil || displayName == "" {
			return models.OptionModel{}, warnings, fmt.Errorf("app.json file (%s), does not contain displayName key", appJSONPth)
		}
	}

	log.TPrintf("Project name: %v", projectName)

	// ios options
	projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
	schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)

	if usesExpoKit {
		projectName = strings.ToLower(regexp.MustCompile(`(?i:[^a-z0-9_\-])`).ReplaceAllString(projectName, "-"))
		projectPathOption.AddOption(filepath.Join("./", "ios", projectName+".xcworkspace"), schemeOption)
	} else {
		projectPathOption.AddOption(filepath.Join("./", "ios", projectName+".xcodeproj"), schemeOption)
	}

	exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
	schemeOption.AddOption(projectName, exportMethodOption)

	// android options
	projectLocationOption := models.NewOption(android.ProjectLocationInputTitle, android.ProjectLocationInputEnvKey)
	for _, exportMethod := range ios.IosExportMethods {
		exportMethodOption.AddOption(exportMethod, projectLocationOption)
	}

	moduleOption := models.NewOption(android.ModuleInputTitle, android.ModuleInputEnvKey)
	projectLocationOption.AddOption("./android", moduleOption)

	buildVariantOption := models.NewOption(android.BuildVariantInputTitle, android.BuildVariantInputEnvKey)
	moduleOption.AddOption("app", buildVariantOption)

	// expo options
	if scanner.usesExpoKit {
		userNameOption := models.NewOption("Expo username", "EXPO_USERNAME")
		buildVariantOption.AddOption("Release", userNameOption)

		passwordOption := models.NewOption("Expo password", "EXPO_PASSWORD")
		userNameOption.AddOption("_", passwordOption)

		configOption := models.NewConfigOption(configName)
		passwordOption.AddConfig("_", configOption)
	} else {
		configOption := models.NewConfigOption(configName)
		buildVariantOption.AddConfig("Release", configOption)
	}

	return *projectPathOption, warnings, nil
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	packageJSONDir := filepath.Dir(scanner.packageJSONPth)
	relPackageJSONDir, err := utility.RelPath(scanner.searchDir, packageJSONDir)
	if err != nil {
		return models.BitriseConfigMap{}, fmt.Errorf("Failed to get relative config.xml dir path, error: %s", err)
	}
	if relPackageJSONDir == "." {
		// config.xml placed in the search dir, no need to change-dir in the workflows
		relPackageJSONDir = ""
	}

	workdirEnvList := []envmanModels.EnvironmentItemModel{}
	if relPackageJSONDir != "" {
		workdirEnvList = append(workdirEnvList, envmanModels.EnvironmentItemModel{reactnative.WorkDirInputKey: relPackageJSONDir})
	}

	// primary workflow
	configBuilder := models.NewDefaultConfigBuilder()
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "test"})...))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

	// deploy workflow
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))

	if scanner.usesExpoKit {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ExpoDetachStepListItem(
			envmanModels.EnvironmentItemModel{"project_path": relPackageJSONDir},
			envmanModels.EnvironmentItemModel{"user_name": "$EXPO_USERNAME"},
			envmanModels.EnvironmentItemModel{"password": "$EXPO_PASSWORD"},
		))
	} else {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ExpoDetachStepListItem(
			envmanModels.EnvironmentItemModel{"project_path": relPackageJSONDir},
		))
	}

	// android build
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: "$" + android.ProjectLocationInputEnvKey},
		envmanModels.EnvironmentItemModel{android.ModuleInputKey: "$" + android.ModuleInputEnvKey},
		envmanModels.EnvironmentItemModel{android.VariantInputKey: "$" + android.BuildVariantInputEnvKey},
	))

	// ios build
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

	if scanner.usesExpoKit {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CocoapodsInstallStepListItem())
	}

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
		envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

	bitriseDataModel, err := configBuilder.Generate(Name)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(bitriseDataModel)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	configMap[configName] = string(data)

	return configMap, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{}
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{
		reactnative.Name,
		string(ios.XcodeProjectTypeIOS),
		string(ios.XcodeProjectTypeMacOS),
		android.ScannerName,
	}
}
