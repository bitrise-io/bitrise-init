package cordova

import (
	"encoding/xml"
	"fmt"

	yaml "gopkg.in/yaml.v1"

	"path/filepath"

	"strings"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/xcode"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
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
	pathInputKey              = "path"
	projectRootDirInputEnvKey = "PROJECT_ROOT_DIR"
	projectRootDirInputTitle  = "Path to the project root direcotry"
)

const (
	forceTeamIDInputKey              = "force_team_id"
	forceCodeSignIdentityInputKey    = "force_code_sign_identity"
	forceProvisioningProfileInputKey = "force_provisioning_profile"
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

// EngineModel ...
type EngineModel struct {
	Name string `xml:"name,attr"`
	Spec string `xml:"spec,attr"`
}

// WidgetModel ...
type WidgetModel struct {
	ID      string        `xml:"id,attr"`
	Version string        `xml:"version,attr"`
	Name    string        `xml:"name"`
	Engines []EngineModel `xml:"engine"`
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
func ConfigName(iosConfigName, androidConfigName string) string {
	configName := "cordova"
	if iosConfigName != "" {
		configName += ("-" + strings.TrimSuffix(iosConfigName, "-config"))
	}
	if androidConfigName != "" {
		configName += ("-" + strings.TrimSuffix(androidConfigName, "-config"))
	}
	return configName + "-config"
}

// Scanner ...
type Scanner struct {
	configXMLPth string
	widget       WidgetModel
	platformsDir string

	configDescriptor ConfigDescriptor
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (scanner Scanner) Name() string {
	return scannerName
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

	if configXMLPth == "" {
		log.Printft("platform not detected")
		return false, nil
	}

	log.Printft("config.xml: %s", configXMLPth)

	widget, err := parseConfigXML(configXMLPth)
	if err != nil {
		log.Printft("can not parse config.xml as a Cordova widget, error: %s", err)
		log.Printft("platform not detected")
		return false, nil
	}

	if len(widget.Engines) == 0 {
		log.Printft("no engines found in config.xml")
		log.Printft("platform not detected")
		return false, nil
	}

	log.Doneft("Platform detected")

	scanner.configXMLPth = configXMLPth
	scanner.widget = widget

	return true, nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}
	projectRoot := filepath.Dir(scanner.configXMLPth)

	projectTypes := []string{}
	hasIosProject := false
	hasAndroidProject := false

	log.Printft("available project types:")

	for _, engine := range scanner.widget.Engines {
		log.Printft("- %s", engine.Name)

		projectTypes = append(projectTypes, engine.Name)

		if engine.Name == "ios" {
			hasIosProject = true
		}

		if engine.Name == "android" {
			hasAndroidProject = true
		}
	}
	if hasIosProject && hasAndroidProject {
		projectTypes = append(projectTypes, "ios+android")
	}

	log.Infoft("Checking for platforms directory")

	isPrepareRequired := false
	platformsDir := filepath.Join(projectRoot, platformsDirName)
	platformsDirExist, err := pathutil.IsPathExists(platformsDir)
	if err != nil {
		return models.OptionModel{}, warnings, fmt.Errorf("failed to check if path (%s) exists, error: %s", platformsDir, err)
	}

	log.Printft("platforms directory exists: %v", platformsDirExist)

	if platformsDirExist {
		for _, engine := range scanner.widget.Engines {
			if engine.Name == "ios" {
				iosPlatformDir := filepath.Join(platformsDir, "ios")
				iosPlatformDirExist, err := pathutil.IsPathExists(iosPlatformDir)
				if err != nil {
					return models.OptionModel{}, warnings, fmt.Errorf("failed to check if path (%s) exists, error: %s", iosPlatformDir, err)
				}

				log.Printft("platforms/ios directory exists: %v", iosPlatformDirExist)

				if !iosPlatformDirExist {
					log.Printft("platforms directory exists: %v", platformsDirExist)
					isPrepareRequired = true
				}
			}

			if engine.Name == "android" {
				androidPlatformDir := filepath.Join(platformsDir, "android")
				androidPlatformDirExist, err := pathutil.IsPathExists(androidPlatformDir)
				if err != nil {
					return models.OptionModel{}, warnings, fmt.Errorf("failed to check if path (%s) exists, error: %s", androidPlatformDir, err)
				}

				log.Printft("platforms/android directory exists: %v", androidPlatformDirExist)

				if !androidPlatformDirExist {
					isPrepareRequired = true
				}
			}
		}
	} else {
		isPrepareRequired = true
	}

	scanner.platformsDir = platformsDir

	if isPrepareRequired {
		log.Infoft("Prepareing project")

		whichCordovaCmd := command.New("which", "cordova")
		out, err := whichCordovaCmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil || out == "" {
			log.Printft("cordova not installed, installing...")

			installCordovaCmd := command.New("npm", "install", "-g", "cordova")
			if err := installCordovaCmd.Run(); err != nil {
				return models.OptionModel{}, warnings, err
			}
		} else {
			log.Printft("cordova installed")
		}

		prepareCmd := command.NewWithStandardOuts("cordova", "prepare")

		log.Printft("$ %s", prepareCmd.PrintableCommandArgs())

		if err := prepareCmd.Run(); err != nil {
			return models.OptionModel{}, warnings, err
		}

		exist, err := pathutil.IsPathExists(platformsDir)
		if err != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to check if path (%s) exists, error: %s", platformsDir, err)
		} else if !exist {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to generate platforms")
		}
	}

	// ---

	projectTypeOption := models.NewOption(projectTypeInputTitle, projectTypeInputEnvKey)

	// ---

	configDescrptor := NewConfigDescriptor()

	iosOptions := new(models.OptionModel)
	iosConfigDescriptors := []xcode.ConfigDescriptor{}
	if hasIosProject {
		platformDir := filepath.Join(platformsDir, "ios")

		var detectorErr error
		if changeDirForFunctionErr := pathutil.ChangeDirForFunction(platformDir, func() {
			log.Printft("")

			detected, err := xcode.Detect(utility.XcodeProjectTypeIOS, "./")
			if err != nil {
				detectorErr = fmt.Errorf("failed to detect ios platform, error: %s", err)
				return
			}
			if !detected {
				detectorErr = fmt.Errorf("config.xml contains ios project, but ios scanner does not detect platform")
				return
			}

			options, configDescriptors, warnings, err := xcode.GenerateOptions(utility.XcodeProjectTypeIOS, "./")
			if err != nil {
				detectorErr = fmt.Errorf("failed to create ios project options, error: %s", err)
				return
			}
			for _, warning := range warnings {
				log.Warnft(warning)
			}

			iosOptions = &options
			iosConfigDescriptors = configDescriptors
		}); changeDirForFunctionErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to change dir to: %s, error: %s", platformsDir, changeDirForFunctionErr)
		}

		if detectorErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to create options, error: %s", detectorErr)
		}

		iosOptionsCopy := iosOptions.Copy()
		iosConfigOptions := iosOptionsCopy.LastChilds()
		for _, iosConfigOption := range iosConfigOptions {
			iosConfigOption.Config = ConfigName(iosConfigOption.Config, "")
		}

		projectTypeOption.AddOption("ios", iosOptionsCopy)

		configDescrptor.iosConfigDescriptors = iosConfigDescriptors
	}

	androidOptions := new(models.OptionModel)
	androidConfigDescriptors := []android.ConfigDescriptor{}
	if hasAndroidProject {
		platformDir := filepath.Join(platformsDir, "android")

		var detectorErr error
		if changeDirForFunctionErr := pathutil.ChangeDirForFunction(platformDir, func() {
			log.Printft("")

			androidScanner := android.NewScanner()

			detected, err := androidScanner.DetectPlatform(".")
			if err != nil {
				detectorErr = fmt.Errorf("failed to detect android platform, error: %s", err)
				return
			}
			if !detected {
				detectorErr = fmt.Errorf("config.xml contains android project, but android scanner does not detect platform")
				return
			}

			options, configDescriptors, warnings, err := androidScanner.GenerateOption(true, false)
			if err != nil {
				detectorErr = fmt.Errorf("failed to create android project options, error: %s", err)
				return
			}
			for _, warning := range warnings {
				log.Warnft(warning)
			}

			androidOptions = &options
			androidConfigDescriptors = configDescriptors
		}); changeDirForFunctionErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to change dir to: %s, error: %s", platformsDir, changeDirForFunctionErr)
		}

		if detectorErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to create options, error: %s", detectorErr)
		}

		androidOptionsCopy := androidOptions.Copy()
		androidConfigOptions := androidOptionsCopy.LastChilds()
		for _, androidConfigOption := range androidConfigOptions {
			androidConfigOption.Config = ConfigName("", androidConfigOption.Config)
		}

		projectTypeOption.AddOption("android", androidOptionsCopy)

		configDescrptor.androidConfigDescriptors = androidConfigDescriptors
	}

	if hasIosProject && hasAndroidProject {
		iosOptionsCopy := iosOptions.Copy()
		androidOptionsCopy := androidOptions.Copy()

		iosConfigOptions := iosOptions.LastChilds()
		for _, iosConfigOption := range iosConfigOptions {
			androidConfigOptions := androidOptionsCopy.LastChilds()
			for _, androidConfigOption := range androidConfigOptions {
				androidConfigOption.Config = ConfigName(iosConfigOption.Config, androidConfigOption.Config)
			}

			iosLastOption, underKey, ok := iosConfigOption.Parent()
			if !ok {
				return models.OptionModel{}, warnings, fmt.Errorf("invalid config: %s", iosConfigOption)
			}

			iosLastOptionCopy, ok := iosOptionsCopy.Child(iosLastOption.Components...)
			if !ok {
				return models.OptionModel{}, warnings, fmt.Errorf("invalid config: %s", iosOptionsCopy)
			}
			iosLastOptionCopy.AddOption(underKey, androidOptionsCopy)
		}

		projectTypeOption.AddOption("ios+android", iosOptionsCopy)
	}

	scanner.configDescriptor = configDescrptor

	return *projectTypeOption, warnings, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{}
}

func cordovaIosStepList(iosPlatformDir string, missingSharedSchemes bool, hasPodfile bool, carthageCommand string, hasTest bool) []bitriseModels.StepListItemModel {
	iosStepList := []bitriseModels.StepListItemModel{}
	iosStepList = append(iosStepList, steps.ChangeWorkDirStepListItem(
		envmanModels.EnvironmentItemModel{steps.ChangeWorkDirInputPathKey: iosPlatformDir}),
	)
	iosStepList = append(iosStepList, steps.CertificateAndProfileInstallerStepListItem())

	if missingSharedSchemes {
		iosStepList = append(iosStepList, steps.RecreateUserSchemesStepListItem(
			envmanModels.EnvironmentItemModel{xcode.ProjectPathInputKey: "$" + xcode.ProjectPathInputEnvKey},
		))
	}

	if hasPodfile {
		iosStepList = append(iosStepList, steps.CocoapodsInstallStepListItem())
	}

	if carthageCommand != "" {
		iosStepList = append(iosStepList, steps.CarthageStepListItem(
			envmanModels.EnvironmentItemModel{xcode.CarthageCommandInputKey: carthageCommand},
		))
	}

	xcodeTestAndArchiveStepInputModels := []envmanModels.EnvironmentItemModel{
		envmanModels.EnvironmentItemModel{forceTeamIDInputKey: "72SA8V3WYL"},
		envmanModels.EnvironmentItemModel{forceCodeSignIdentityInputKey: "iPhone Developer"},
		envmanModels.EnvironmentItemModel{forceProvisioningProfileInputKey: "BitriseBot-Wildcard"},
		envmanModels.EnvironmentItemModel{xcode.ProjectPathInputKey: "$" + xcode.ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{xcode.SchemeInputKey: "$" + xcode.SchemeInputEnvKey},
	}

	if hasTest {
		iosStepList = append(iosStepList, steps.XcodeTestStepListItem(xcodeTestAndArchiveStepInputModels...))
	}

	iosStepList = append(iosStepList, steps.XcodeArchiveStepListItem(xcodeTestAndArchiveStepInputModels...))

	return iosStepList
}

func cordovaAndroidStepList(androidPlatformDir string, missingGraldew bool) []bitriseModels.StepListItemModel {
	androidStepList := []bitriseModels.StepListItemModel{}
	androidStepList = append(androidStepList, steps.ChangeWorkDirStepListItem(envmanModels.EnvironmentItemModel{pathInputKey: androidPlatformDir}))
	if missingGraldew {
		androidStepList = append(androidStepList, steps.GenerateGradleWrapperStepListItem())
	}
	androidStepList = append(androidStepList, steps.InstallMissingAndroidToolsStepListItem())
	androidStepList = append(androidStepList, steps.GradleRunnerStepListItem(
		envmanModels.EnvironmentItemModel{android.GradleFileInputKey: "$" + android.GradleFileInputEnvKey},
		envmanModels.EnvironmentItemModel{android.GradleTaskInputKey: "$" + android.GradleTaskInputEnvKey},
		envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.GradlewPathInputEnvKey},
	))
	return androidStepList
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configMap := models.BitriseConfigMap{}

	// common steps
	commonStepList := steps.DefaultPrepareStepList()

	if len(scanner.configDescriptor.iosConfigDescriptors) > 0 {
		for _, descriptor := range scanner.configDescriptor.iosConfigDescriptors {
			// ios
			iosStepList := cordovaIosStepList(filepath.Join("$PROJECT_ROOT_DIR", scanner.platformsDir, "ios"), descriptor.MissingSharedSchemes, descriptor.HasPodfile, descriptor.CarthageCommand, descriptor.HasTest)

			iosConfig, err := models.NewConfigBuilder(append(commonStepList, iosStepList...)).Generate(envmanModels.EnvironmentItemModel{"PROJECT_ROOT_DIR": "$BITRISE_SOURCE_DIR"})
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			iosConfigData, err := yaml.Marshal(iosConfig)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			iosConfigName := descriptor.ConfigName(utility.XcodeProjectTypeIOS)
			cordovaIosConfigName := ConfigName(iosConfigName, "")
			configMap[cordovaIosConfigName] = string(iosConfigData)

			for _, descriptor := range scanner.configDescriptor.androidConfigDescriptors {
				// android
				androidStepList := cordovaAndroidStepList(filepath.Join("$PROJECT_ROOT_DIR", scanner.platformsDir, "android"), descriptor.MissingGradlew)

				androidConfig, err := models.NewConfigBuilder(append(commonStepList, androidStepList...)).Generate(envmanModels.EnvironmentItemModel{"PROJECT_ROOT_DIR": "$BITRISE_SOURCE_DIR"})
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				androidConfigData, err := yaml.Marshal(androidConfig)
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				androidConfigName := android.ConfigName
				cordovaAndroidConfigName := ConfigName("", androidConfigName)
				configMap[cordovaAndroidConfigName] = string(androidConfigData)

				// ios + android
				iosAndroidSteplist := append(iosStepList, androidStepList...)

				iosAndroidConfig, err := models.NewConfigBuilder(append(commonStepList, iosAndroidSteplist...)).Generate(envmanModels.EnvironmentItemModel{"PROJECT_ROOT_DIR": "$BITRISE_SOURCE_DIR"})
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				iosAndroidConfigData, err := yaml.Marshal(iosAndroidConfig)
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				cordovaIosAndroidConfigName := ConfigName(iosConfigName, androidConfigName)
				configMap[cordovaIosAndroidConfigName] = string(iosAndroidConfigData)
			}
		}
	} else if len(scanner.configDescriptor.androidConfigDescriptors) > 0 {
		for _, descriptor := range scanner.configDescriptor.androidConfigDescriptors {
			androidStepList := cordovaAndroidStepList(filepath.Join("$PROJECT_ROOT_DIR", scanner.platformsDir, "android"), descriptor.MissingGradlew)

			androidConfig, err := models.NewConfigBuilder(append(commonStepList, androidStepList...)).Generate(envmanModels.EnvironmentItemModel{"PROJECT_ROOT_DIR": "$BITRISE_SOURCE_DIR"})
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			androidConfigData, err := yaml.Marshal(androidConfig)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			androidConfigName := android.ConfigName
			cordovaAndroidConfigName := ConfigName("", androidConfigName)
			configMap[cordovaAndroidConfigName] = string(androidConfigData)
		}
	}

	return configMap, nil
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}
