package fastlane

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"bufio"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
)

const scannerName = "fastlane"

const (
	fastfileBasePath = "Fastfile"
)

const (
	laneInputKey    = "lane"
	laneInputTitle  = "Fastlane lane"
	laneInputEnvKey = "FASTLANE_LANE"
)

const (
	workDirInputKey    = "work_dir"
	workDirInputTitle  = "Working directory"
	workDirInputEnvKey = "FASTLANE_WORK_DIR"
)

const (
	fastlaneXcodeListTimeoutEnvKey   = "FASTLANE_XCODE_LIST_TIMEOUT"
	fastlaneXcodeListTimeoutEnvValue = "120"
)

//--------------------------------------------------
// Utility
//--------------------------------------------------

func filterFastfiles(fileList []string) ([]string, error) {
	allowFastfileBaseFilter := utility.BaseFilter(fastfileBasePath, true)
	fastfiles, err := utility.FilterPaths(fileList, allowFastfileBaseFilter)
	if err != nil {
		return []string{}, err
	}

	return utility.SortPathsByComponents(fastfiles)
}

func inspectFastfileContent(content string) ([]string, error) {
	commonLanes := []string{}
	laneMap := map[string][]string{}

	// platform :ios do ...
	platformSectionStartRegexp := regexp.MustCompile(`platform\s+:(?P<platform>.*)\s+do`)
	platformSectionEndPattern := "end"
	platform := ""

	// lane :test_and_snapshot do
	laneRegexp := regexp.MustCompile(`^[\s]*lane\s+:(?P<lane>.*)\s+do`)

	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " ")

		if platform != "" && line == platformSectionEndPattern {
			platform = ""
			continue
		}

		if platform == "" {
			if match := platformSectionStartRegexp.FindStringSubmatch(line); len(match) == 2 {
				platform = match[1]
				continue
			}
		}

		if match := laneRegexp.FindStringSubmatch(line); len(match) == 2 {
			lane := match[1]

			if platform != "" {
				lanes, found := laneMap[platform]
				if !found {
					lanes = []string{}
				}
				lanes = append(lanes, lane)
				laneMap[platform] = lanes
			} else {
				commonLanes = append(commonLanes, lane)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	lanes := commonLanes
	for platform, platformLanes := range laneMap {
		for _, lane := range platformLanes {
			lanes = append(lanes, platform+" "+lane)
		}
	}

	return lanes, nil
}

func inspectFastfile(fastFile string) ([]string, error) {
	content, err := fileutil.ReadStringFromFile(fastFile)
	if err != nil {
		return []string{}, err
	}

	return inspectFastfileContent(content)
}

// Returns:
//  - fastlane dir's parent, if Fastfile is in fastlane dir (test/fastlane/Fastfile)
//  - Fastfile's dir, if Fastfile is NOT in fastlane dir (test/Fastfile)
func fastlaneWorkDir(fastfilePth string) string {
	dirPth := filepath.Dir(fastfilePth)
	dirName := filepath.Base(dirPth)
	if dirName == "fastlane" {
		return filepath.Dir(dirPth)
	}
	return dirPth
}

func configName() string {
	return "fastlane-config"
}

func defaultConfigName() string {
	return "default-fastlane-config"
}

//--------------------------------------------------
// Scanner
//--------------------------------------------------

// Scanner ...
type Scanner struct {
	Fastfiles []string
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

	// Search for Fastfile
	log.Infoft("Searching for Fastfiles")

	fastfiles, err := filterFastfiles(fileList)
	if err != nil {
		return false, fmt.Errorf("failed to search for Fastfile in (%s), error: %s", searchDir, err)
	}

	scanner.Fastfiles = fastfiles

	log.Printft("%d Fastfiles detected", len(fastfiles))
	for _, file := range fastfiles {
		log.Printft("- %s", file)
	}

	if len(fastfiles) == 0 {
		log.Printft("platform not detected")
		return false, nil
	}

	log.Doneft("Platform detected")

	return true, nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}

	isValidFastfileFound := false

	// Inspect Fastfiles

	workDirOption := models.NewOption(workDirInputTitle, workDirInputEnvKey)

	for _, fastfile := range scanner.Fastfiles {
		log.Infoft("Inspecting Fastfile: %s", fastfile)

		workDir := fastlaneWorkDir(fastfile)
		log.Printft("fastlane work dir: %s", workDir)

		lanes, err := inspectFastfile(fastfile)
		if err != nil {
			log.Warnft("Failed to inspect Fastfile, error: %s", err)
			warnings = append(warnings, fmt.Sprintf("Failed to inspect Fastfile (%s), error: %s", fastfile, err))
			continue
		}

		log.Printft("%d lanes found", len(lanes))

		if len(lanes) == 0 {
			log.Warnft("No lanes found")
			warnings = append(warnings, fmt.Sprintf("No lanes found for Fastfile: %s", fastfile))
			continue
		}

		isValidFastfileFound = true

		laneOption := models.NewOption(laneInputTitle, laneInputEnvKey)
		workDirOption.AddOption(workDir, laneOption)

		for _, lane := range lanes {
			log.Printft("- %s", lane)

			configOption := models.NewConfigOption(configName())
			laneOption.AddConfig(lane, configOption)
		}
	}

	if !isValidFastfileFound {
		log.Errorft("No valid Fastfile found")
		warnings = append(warnings, "No valid Fastfile found")
		return models.OptionModel{}, warnings, nil
	}

	return *workDirOption, warnings, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionModel {
	workDirOption := models.NewOption(workDirInputTitle, workDirInputEnvKey)

	laneOption := models.NewOption(laneInputTitle, laneInputEnvKey)
	workDirOption.AddOption("_", laneOption)

	configOption := models.NewConfigOption(defaultConfigName())
	laneOption.AddConfig("_", configOption)

	return *workDirOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	configBuilder.AppendPreparStepList(steps.CertificateAndProfileInstallerStepListItem())

	configBuilder.AppendMainStepList(steps.FastlaneStepListItem(
		envmanModels.EnvironmentItemModel{laneInputKey: "$" + laneInputEnvKey},
		envmanModels.EnvironmentItemModel{workDirInputKey: "$" + workDirInputEnvKey},
	))

	config, err := configBuilder.Generate(envmanModels.EnvironmentItemModel{fastlaneXcodeListTimeoutEnvKey: fastlaneXcodeListTimeoutEnvValue})
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		configName(): string(data),
	}, nil
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	configBuilder.AppendPreparStepList(steps.CertificateAndProfileInstallerStepListItem())

	configBuilder.AppendMainStepList(steps.FastlaneStepListItem(
		envmanModels.EnvironmentItemModel{laneInputKey: "$" + laneInputEnvKey},
		envmanModels.EnvironmentItemModel{workDirInputKey: "$" + workDirInputEnvKey},
	))

	config, err := configBuilder.Generate(envmanModels.EnvironmentItemModel{fastlaneXcodeListTimeoutEnvKey: fastlaneXcodeListTimeoutEnvValue})
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		defaultConfigName(): string(data),
	}, nil
}
