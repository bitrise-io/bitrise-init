package cordova

import (
	"encoding/xml"
	"fmt"

	"path/filepath"

	"encoding/json"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/xcode"
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

const scannerName = "cordova"

const (
	configXMLBasePath = "config.xml"
	platformsDirName  = "platforms"

	projectTypeKey    = "project_type"
	projectTypeTitle  = "Project type"
	projectTypeEnvKey = "PROJECT_TYPE"
)

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

func configName(projectType string) string {
	return fmt.Sprintf("cordova-%s", projectType)
}

// Scanner ...
type Scanner struct {
	fileList     []string
	configXMLPth string
	widget       WidgetModel
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

	scanner.fileList = fileList
	scanner.configXMLPth = configXMLPth
	scanner.widget = widget

	return true, nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}
	projectRoot := filepath.Dir(scanner.configXMLPth)

	log.Infoft("Checking for platforms directory")

	platformsDir := filepath.Join(projectRoot, platformsDirName)
	exist, err := pathutil.IsPathExists(platformsDir)
	if err != nil {
		return models.OptionModel{}, warnings, fmt.Errorf("failed to check if path (%s) exists, error: %s", platformsDir, err)
	}

	log.Printft("platforms directory exists: %v", exist)

	if !exist {
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

	projectTypeOption := models.NewOption(projectTypeTitle, projectTypeEnvKey)

	// ---

	iosOptions := new(models.OptionModel)
	if hasIosProject {
		platformDir := filepath.Join(platformsDir, "ios")

		var detectorErr error
		if changeDirForFunctionErr := pathutil.ChangeDirForFunction(platformDir, func() {
			log.Printft("")

			iosScanner := &xcode.Scanner{ProjectType: xcode.ProjectTypeIOS}

			detected, err := iosScanner.CommonDetectPlatform("./")
			if err != nil {
				detectorErr = fmt.Errorf("failed to detect ios platform, error: %s", err)
				return
			}
			if !detected {
				detectorErr = fmt.Errorf("config.xml contains ios project, but ios scannern does not detect platform")
				return
			}

			options, warnings, err := iosScanner.GenerateOptions(false)
			if err != nil {
				detectorErr = fmt.Errorf("failed to create ios project options, error: %s", err)
				return
			}
			for _, warning := range warnings {
				log.Warnft(warning)
			}

			iosOptions = &options
		}); changeDirForFunctionErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to change dir to: %s, error: %s", platformsDir, err)
		}

		if detectorErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to create options, error: %s", detectorErr)
		}

		projectTypeOption.AddOption("ios", iosOptions)

		bytes, err := json.MarshalIndent(iosOptions, "", "  ")
		if err != nil {
			log.Errorft("Failed to marshal, error: %s", err)
		}
		log.Doneft("\niosOptions: %s\n", string(bytes))
	}

	androidOptions := new(models.OptionModel)
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
				detectorErr = fmt.Errorf("config.xml contains android project, but android scannern does not detect platform")
				return
			}

			options, warnings, err := androidScanner.GenerateOption(false, true)
			if err != nil {
				detectorErr = fmt.Errorf("failed to create android project options, error: %s", err)
				return
			}
			for _, warning := range warnings {
				log.Warnft(warning)
			}

			androidOptions = &options
		}); changeDirForFunctionErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to change dir to: %s, error: %s", platformsDir, err)
		}

		if detectorErr != nil {
			return models.OptionModel{}, warnings, fmt.Errorf("failed to create options, error: %s", detectorErr)
		}

		projectTypeOption.AddOption("android", androidOptions)

		bytes, err := json.MarshalIndent(androidOptions, "", "  ")
		if err != nil {
			log.Errorft("Failed to marshal, error: %s", err)
		}
		log.Doneft("\nandroidOptions: %s\n", string(bytes))
	}

	if hasIosProject && hasAndroidProject {
		iosOptionsCopy := iosOptions.Copy()

		lastOptions := iosOptionsCopy.LastOptions()
		for _, lastOption := range lastOptions {
			for value, childOption := range lastOption.ChildOptionMap {
				if childOption != nil {
					log.Errorft("Child should be nil")
				}
				lastOption.AddOption(value, androidOptions)
			}
		}
		projectTypeOption.AddOption("ios+android", iosOptionsCopy)
	}

	bytes, err := json.MarshalIndent(projectTypeOption, "", "  ")
	if err != nil {
		log.Errorft("Failed to marshal, error: %s", err)
	}
	log.Doneft("\nprojectTypeOption: %s\n", string(bytes))

	return *projectTypeOption, warnings, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{}
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return models.BitriseConfigMap{}, nil
}
