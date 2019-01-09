package scanner

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-init/toolscanner"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
)

type scannerErrors struct {
	errors   models.Errors
	warnings models.Warnings
}

type scannerOutput struct {
	optionModel          models.OptionModel
	configMap            models.BitriseConfigMap
	excludedScannerNames []string
}

// Config ...
func Config(searchDir string) models.ScanResultModel {
	result := models.ScanResultModel{}

	//
	// Setup
	currentDir, err := os.Getwd()
	if err != nil {
		result.AddError("general", fmt.Sprintf("Failed to expand current directory path, error: %s", err))
		return result
	}

	if searchDir == "" {
		searchDir = currentDir
	} else {
		absScerach, err := pathutil.AbsPath(searchDir)
		if err != nil {
			result.AddError("general", fmt.Sprintf("Failed to expand path (%s), error: %s", searchDir, err))
			return result
		}
		searchDir = absScerach
	}

	if searchDir != currentDir {
		if err := os.Chdir(searchDir); err != nil {
			result.AddError("general", fmt.Sprintf("Failed to change dir, to (%s), error: %s", searchDir, err))
			return result
		}
		defer func() {
			if err := os.Chdir(currentDir); err != nil {
				log.TWarnf("Failed to change dir, to (%s), error: %s", searchDir, err)
			}
		}()
	}
	// ---

	//
	// Scan
	projectScanners := scanners.ProjectScanners

	projectTypeErrorMap := map[string]models.Errors{}
	projectTypeWarningMap := map[string]models.Warnings{}
	projectTypeOptionMap := map[string]models.OptionModel{}
	projectTypeConfigMap := map[string]models.BitriseConfigMap{}

	excludedScannerNames := []string{}

	log.TInfof(colorstring.Blue("Running scanners:"))
	fmt.Println()

	for _, detector := range projectScanners {
		detectorName, detectorErrors, detectorWarnings, scannerOutput := checkScannerMatchAndReturnOutput(detector, searchDir, excludedScannerNames)
		if scannerOutput != nil {
			projectTypeOptionMap[detectorName] = scannerOutput.optionModel
			projectTypeConfigMap[detectorName] = scannerOutput.configMap
			excludedScannerNames = append(excludedScannerNames, scannerOutput.excludedScannerNames...)
			projectTypeErrorMap[detectorName] = append(projectTypeErrorMap[detectorName], detectorErrors...)
			projectTypeWarningMap[detectorName] = append(projectTypeWarningMap[detectorName], detectorWarnings...)
		}
	}
	// ---

	return models.ScanResultModel{
		PlatformOptionMap:    projectTypeOptionMap,
		PlatformConfigMapMap: projectTypeConfigMap,
		PlatformWarningsMap:  projectTypeWarningMap,
		PlatformErrorsMap:    projectTypeErrorMap,
	}
}

func checkScannerMatchAndReturnOutput(detector scanners.ScannerInterface, searchDir string, excludedScannerNamesPrevious []string) (string, models.Errors, models.Warnings, *scannerOutput) {
	detectorName := detector.Name()
	var detectorWarnings []string
	var detectorErrors []string

	log.TInfof("Scanner: %s", colorstring.Blue(detectorName))

	if sliceutil.IsStringInSlice(detectorName, excludedScannerNamesPrevious) {
		log.TWarnf("scanner is marked as excluded, skipping...")
		fmt.Println()
		return detectorName, detectorErrors, detectorWarnings, nil
	}

	log.TPrintf("+------------------------------------------------------------------------------+")
	log.TPrintf("|                                                                              |")

	detected, err := detector.DetectPlatform(searchDir)
	if err != nil {
		log.TErrorf("Scanner failed, error: %s", err)
		detectorWarnings = append(detectorWarnings, err.Error())
		detected = false
	}

	if !detected {
		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()
		return detectorName, detectorErrors, detectorWarnings, nil
	}

	options, projectWarnings, err := detector.Options()
	detectorWarnings = append(detectorWarnings, projectWarnings...)

	if err != nil {
		log.TErrorf("Analyzer failed, error: %s", err)
		detectorWarnings = append(detectorWarnings, err.Error())

		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()
		return detectorName, detectorErrors, detectorWarnings, nil
	}

	// Generate configs
	configs, err := detector.Configs()
	if err != nil {
		log.TErrorf("Failed to generate config, error: %s", err)
		detectorErrors = append(detectorErrors, err.Error())
		return detectorName, detectorErrors, detectorWarnings, nil
	}

	log.TPrintf("|                                                                              |")
	log.TPrintf("+------------------------------------------------------------------------------+")

	excludedScannerNamesCurrent := detector.ExcludedScannerNames()
	if len(excludedScannerNamesCurrent) > 0 {
		log.TWarnf("Scanner will exclude scanners: %v", excludedScannerNamesCurrent)
		excludedScannerNamesCurrent = append(excludedScannerNamesPrevious, excludedScannerNamesCurrent...)
	}

	fmt.Println()
	return detectorName, detectorErrors, detectorWarnings, &scannerOutput{options, configs, excludedScannerNamesCurrent}
}

// addProjectTypeToToolScanner is used to add a project type for automation tool scanners's option map
func addProjectTypeToToolScanner(toolScannerOptionModel models.OptionModel, detectedProjectTypes []string) *models.OptionModel {
	projectTypeOption := models.NewOption(toolscanner.ProjectTypeUserTitle, toolscanner.ProjectTypeEnvKey)
	for _, projectType := range detectedProjectTypes {
		projectTypeOption.AddOption(projectType, &toolScannerOptionModel)
	}
	return projectTypeOption
}
