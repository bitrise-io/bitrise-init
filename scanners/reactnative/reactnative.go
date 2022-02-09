package reactnative

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/serialized"
)

const scannerName = "react-native"

const (
	// workDirInputKey is a key of the working directory step input.
	workDirInputKey = "workdir"
)

const (
	isExpoCLIInputTitle   = "Was your React Native app created with the Expo CLI and using Managed Workflow?"
	isExpoCLIInputSummary = "Will include *Expo Eject** Step if using Expo Managed Workflow (https://docs.expo.io/introduction/managed-vs-bare/). If ios/android native projects are present in the repository, choose No."
)

// Scanner implements the project scanner for plain React Native and Expo based projects.
type Scanner struct {
	searchDir       string
	iosScanner      *ios.Scanner
	androidProjects []android.Project

	hasTest         bool
	hasYarnLockFile bool
	packageJSONPth  string

	expoSettings *expoSettings
}

// NewScanner creates a new scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name implements ScannerInterface.Name function.
func (Scanner) Name() string {
	return scannerName
}

type expoSettings struct {
	name                string
	isIOS, isAndroid    bool
	bundleIdentifierIOS string
	packageNameAndroid  string
}

func (settings *expoSettings) isAllIdentifierPresent() bool {
	return !(settings.isAndroid && settings.packageNameAndroid == "" ||
		settings.isIOS && settings.bundleIdentifierIOS == "")
}

// parseExpoProjectSettings reports whether a project is Expo based and it's settings, like targeted platforms
func parseExpoProjectSettings(packageJSONPth string) (*expoSettings, error) {
	packages, err := utility.ParsePackagesJSON(packageJSONPth)
	if err != nil {
		return nil, fmt.Errorf("failed to parse package json file (%s): %s", packageJSONPth, err)
	}

	if _, found := packages.Dependencies["expo"]; !found {
		return nil, nil
	}

	// app.json file is a required part of an expo projects and should be placed next to the root package.json file
	appJSONPth := filepath.Join(filepath.Dir(packageJSONPth), "app.json")
	exist, err := pathutil.IsPathExists(appJSONPth)
	if err != nil {
		return nil, fmt.Errorf("failed to check if app.json file (%s) exist: %s", appJSONPth, err)
	}
	if !exist {
		return nil, nil
	}

	appJSON, err := fileutil.ReadStringFromFile(appJSONPth)
	if err != nil {
		return nil, err
	}
	var app serialized.Object
	if err := json.Unmarshal([]byte(appJSON), &app); err != nil {
		return nil, err
	}

	expoObj, err := app.Object("expo")
	if err != nil {
		log.Warnf("%s", fmt.Errorf("app.json file (%s) has no 'expo' entry, not an Expo project", appJSONPth))
		return nil, nil
	}
	projectName, err := expoObj.String("name")
	if err != nil || projectName == "" {
		log.Debugf("%s", fmt.Errorf("app.json file (%s) has no 'expo/name' entry, can not guess iOS project path, will ask for it during project configuration", appJSONPth))
	}
	iosObj, err := expoObj.Object("ios")
	if err != nil {
		log.TDebugf("%s", fmt.Errorf("app.json file (%s) has no no 'expo/ios entry', assuming iOS is targeted by Expo", appJSONPth))
	}
	bundleID, err := iosObj.String("bundleIdentifier")
	if err != nil || bundleID == "" {
		log.TDebugf("%s", fmt.Errorf("app.json file (%s) has no no 'expo/ios/bundleIdentifier' entry, will ask for it during project configuration", appJSONPth))
	}
	androidObj, err := expoObj.Object("android")
	if err != nil {
		log.TDebugf("%s", fmt.Errorf("app.json file (%s) has no 'expo/android' entry, assuming Android is targeted by Expo", appJSONPth))
	}
	packageName, err := androidObj.String("package")
	if err != nil || packageName == "" {
		log.TDebugf("%s", fmt.Errorf("app.json file (%s) has no no 'expo/android/package' entry, will ask for it during project configuration", appJSONPth))
	}

	// expo/ios and expo/android entry is optional
	return &expoSettings{
		name:                projectName,
		isIOS:               true,
		isAndroid:           true,
		packageNameAndroid:  packageName,
		bundleIdentifierIOS: bundleID,
	}, nil
}

func hasNativeIOSProject(searchDir, projectDir string, iosScanner *ios.Scanner) (bool, error) {
	absProjectDir, err := pathutil.AbsPath(projectDir)
	if err != nil {
		return false, err
	}

	iosDir := filepath.Join(absProjectDir, "ios")
	if exist, err := pathutil.IsDirExists(iosDir); err != nil || !exist {
		return false, err
	}

	return iosScanner.DetectPlatform(searchDir)
}

func hasNativeAndroidProject(searchDir, projectDir string, androidScanner *android.Scanner) (bool, []android.Project, error) {
	absProjectDir, err := pathutil.AbsPath(projectDir)
	if err != nil {
		return false, nil, err
	}

	androidDir := filepath.Join(absProjectDir, "android")
	if exist, err := pathutil.IsDirExists(androidDir); err != nil || !exist {
		return false, nil, err
	}

	if detected, err := androidScanner.DetectPlatform(searchDir); err != nil || !detected {
		return false, nil, err
	}

	return true, androidScanner.Projects, nil
}

// DetectPlatform implements ScannerInterface.DetectPlatform function.
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.searchDir = searchDir

	log.TInfof("Collect package.json files")

	packageJSONPths, err := CollectPackageJSONFiles(searchDir)
	if err != nil {
		return false, err
	}

	log.TPrintf("%d package.json file detected", len(packageJSONPths))
	log.TPrintf("Filter relevant package.json files")

	var expoSettings *expoSettings
	var packageFile string

	for _, packageJSONPth := range packageJSONPths {
		log.TPrintf("Checking: %s", packageJSONPth)

		expoPrefs, err := parseExpoProjectSettings(packageJSONPth)
		if err != nil {
			log.TWarnf("failed to check if project uses Expo: %s", err)
		}

		log.TPrintf("Project uses expo: %v", expoPrefs != nil)
		if expoPrefs != nil {
			// Treating the project as an Expo based React Native project
			log.TPrintf("Expo configuration: %+v", expoPrefs)
			expoSettings = expoPrefs
			packageFile = packageJSONPth
			break
		}

		if scanner.iosScanner == nil {
			scanner.iosScanner = ios.NewScanner()
			scanner.iosScanner.ExcludeAppIcon = true
		}
		androidScanner := android.NewScanner()

		projectDir := filepath.Dir(packageJSONPth)
		isIOSProject, err := hasNativeIOSProject(searchDir, projectDir, scanner.iosScanner)
		if err != nil {
			log.TWarnf("failed to check native iOS projects: %s", err)
		}
		log.TPrintf("Found native ios project: %v", isIOSProject)
		if !isIOSProject {
			scanner.iosScanner = nil
		}

		isAndroidProject, androidProjects, err := hasNativeAndroidProject(searchDir, projectDir, androidScanner)
		if err != nil {
			log.TWarnf("failed to check native Android projects: %s", err)
		}
		log.TPrintf("Found native android project: %v", isAndroidProject)
		scanner.androidProjects = androidProjects

		if isIOSProject || isAndroidProject {
			// Treating the project as a plain React Native project
			packageFile = packageJSONPth
			break
		}
	}

	if packageFile == "" {
		return false, nil
	}

	scanner.expoSettings = expoSettings
	scanner.packageJSONPth = packageFile

	// determine Js dependency manager
	if scanner.hasYarnLockFile, err = containsYarnLock(filepath.Dir(scanner.packageJSONPth)); err != nil {
		return false, err
	}
	log.TPrintf("Js dependency manager for %s is yarn: %t", scanner.packageJSONPth, scanner.hasYarnLockFile)

	packages, err := utility.ParsePackagesJSON(scanner.packageJSONPth)
	if err != nil {
		return false, err
	}

	if _, found := packages.Scripts["test"]; found {
		scanner.hasTest = true
	}
	log.TPrintf("Test script found in package.json: %v", scanner.hasTest)

	return true, nil
}

// Options implements ScannerInterface.Options function.
func (scanner *Scanner) Options() (options models.OptionNode, warnings models.Warnings, icons models.Icons, err error) {
	if scanner.expoSettings != nil {
		options, warnings, err = scanner.expoOptions()
	} else {
		options, warnings, err = scanner.options()
	}

	return
}

// Configs implements ScannerInterface.Configs function.
func (scanner *Scanner) Configs(isPrivateRepo bool) (models.BitriseConfigMap, error) {
	if scanner.expoSettings != nil {
		return scanner.expoConfigs(isPrivateRepo)
	}

	return scanner.configs(isPrivateRepo)
}

// DefaultOptions implements ScannerInterface.DefaultOptions function.
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	expoOption := models.NewOption(isExpoCLIInputTitle, isExpoCLIInputSummary, "", models.TypeSelector)

	expoDefaultOptions := scanner.expoDefaultOptions()
	expoOption.AddOption("yes", &expoDefaultOptions)

	defaultOptions := scanner.defaultOptions()
	expoOption.AddOption("no", &defaultOptions)

	return *expoOption
}

// DefaultConfigs implements ScannerInterface.DefaultConfigs function.
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	configs, err := scanner.defaultConfigs()
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	for k, v := range configs {
		configMap[k] = v
	}

	expoConfigs, err := scanner.expoDefaultConfigs()
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	for k, v := range expoConfigs {
		configMap[k] = v
	}

	return configMap, nil
}

// ExcludedScannerNames implements ScannerInterface.ExcludedScannerNames function.
func (Scanner) ExcludedScannerNames() []string {
	return []string{
		string(ios.XcodeProjectTypeIOS),
		string(ios.XcodeProjectTypeMacOS),
		android.ScannerName,
	}
}
