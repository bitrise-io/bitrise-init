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

// type scannerWarnings struct {
// 	Name     string
// 	warnings models.Warnings
// }

type scannerMatchResult struct {
	Name                 string
	Warnings             models.Warnings
	Errors               models.Errors
	OptionModel          models.OptionModel
	ConfigMap            models.BitriseConfigMap
	ExcludedScannerNames []string
}

// type scannerOutput struct {
// 	Warnings    *scannerWarnings
// 	MatchResult *scannerMatchOutput
// }

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
	projectTypeErrorMap := map[string]models.Errors{}
	projectTypeWarningMap := map[string]models.Warnings{}
	projectTypeOptionMap := map[string]models.OptionModel{}
	projectTypeConfigMap := map[string]models.BitriseConfigMap{}

	var excludedScannerNames []string

	log.TInfof(colorstring.Blue("Running scanners:"))
	fmt.Println()

	projectScannerWarnings, projectScannerMatchResults := mapScannersToOutput(scanners.ProjectScanners, searchDir, excludedScannerNames)
	var matchingProjectTypes []string
	for _, match := range projectScannerMatchResults {
		matchingProjectTypes = append(matchingProjectTypes, match.Name)
	}

	toolScannerWarnings, toolScannerResults := mapScannersToOutput(scanners.AutomationToolScanners, searchDir, excludedScannerNames)
	// Add project_type property option to tool scanner's as they do not dertect project/platform
	for _, match := range toolScannerResults {
		match.OptionModel = addProjectTypeToToolScanner(match.OptionModel, matchingProjectTypes)
	}

	for k, warnings := range toolScannerWarnings {
		projectTypeWarningMap[k] = warnings
	}
	// for {
	// 	if scannerOutput != nil {
	// 		projectTypeOptionMap[detectorName] = scannerOutput.optionModel
	// 		projectTypeConfigMap[detectorName] = scannerOutput.configMap
	// 		excludedScannerNames = append(excludedScannerNames, scannerOutput.excludedScannerNames...)
	// 		if len(detectorErrors) > 0 {
	// 			projectTypeErrorMap[detectorName] = append(projectTypeErrorMap[detectorName], detectorErrors...)
	// 		}
	// 		projectTypeWarningMap[detectorName] = append(projectTypeWarningMap[detectorName], detectorWarnings...)
	// 	} else {
	// 		if len(detectorWarnings) > 0 {
	// 			projectTypeWarningMap[detectorName] = append(projectTypeWarningMap[detectorName], detectorWarnings...)
	// 		}
	// 	}
	// }

	// ---

	return models.ScanResultModel{
		PlatformOptionMap:    projectTypeOptionMap,
		PlatformConfigMapMap: projectTypeConfigMap,
		PlatformWarningsMap:  projectTypeWarningMap,
		PlatformErrorsMap:    projectTypeErrorMap,
	}
}

func mapScannersToOutput(scannerList []scanners.ScannerInterface, searchDir string, excludedScannerNames []string) (map[string]models.Warnings, map[string]scannerMatchResult) {
	scannerToMatchResult := map[string]scannerMatchOutput{}
	scannerToWarnings := map[string]models.Warnings{}
	for _, scanner := range scannerList {
		warnings, matchResult := checkScannerMatchAndReturnOutput(scanner, searchDir, excludedScannerNames)
		if detectorWarnings != nil {
			scannerToWarnings[scanner.Name()] = *warnings
		}
		if scannerOutput != nil {
			scannerToMatchResult[scanner.Name()] = matchResult
		}
	}
	return scannerToWarnings, scannerToMatchResult
}

func checkScannerMatchAndReturnOutput(detector scanners.ScannerInterface, searchDir string, excludedScannerNamesPrevious []string) (*models.Warnings, *scannerMatchOutput) {
	detectorName := detector.Name()
	var detectorWarnings models.Warnings
	var detectorErrors []string

	log.TInfof("Scanner: %s", colorstring.Blue(detectorName))

	if sliceutil.IsStringInSlice(detectorName, excludedScannerNamesPrevious) {
		log.TWarnf("scanner is marked as excluded, skipping...")
		fmt.Println()
		return nil, nil
	}

	log.TPrintf("+------------------------------------------------------------------------------+")
	log.TPrintf("|                                                                              |")

	if detected, err := detector.DetectPlatform(searchDir); err != nil {
		log.TErrorf("Scanner failed, error: %s", err)
		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()
		return &models.Warnings{err.Error()}, nil
	} else if !detected {
		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()
		return nil, nil
	}

	options, projectWarnings, err := detector.Options()
	detectorWarnings = append(detectorWarnings, projectWarnings...)

	if err != nil {
		log.TErrorf("Analyzer failed, error: %s", err)
		detectorWarnings = append(detectorWarnings, err.Error())

		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()
		return &detectorWarnings, nil
	}

	// Generate configs
	configs, err := detector.Configs()
	if err != nil {
		log.TErrorf("Failed to generate config, error: %s", err)
		detectorErrors = append(detectorErrors, err.Error())
		return &detectorWarnings, nil
	}

	log.TPrintf("|                                                                              |")
	log.TPrintf("+------------------------------------------------------------------------------+")

	excludedScannerNamesCurrent := detector.ExcludedScannerNames()
	if len(excludedScannerNamesCurrent) > 0 {
		log.TWarnf("Scanner will exclude scanners: %v", excludedScannerNamesCurrent)
		excludedScannerNamesCurrent = append(excludedScannerNamesPrevious, excludedScannerNamesCurrent...)
	}

	fmt.Println()
	return &models.Warnings{}, &scannerMatchOutput{
		Name:                 detectorName,
		Warnings:             detectorWarnings,
		Errors:               detectorErrors,
		OptionModel:          options,
		ConfigMap:            configs,
		ExcludedScannerNames: excludedScannerNamesCurrent,
	}
}

// addProjectTypeToToolScanner is used to add a project type for automation tool scanners's option map
func addProjectTypeToToolScanner(toolScannerOptionModel models.OptionModel, detectedProjectTypes []string) models.OptionModel {
	projectTypeOption := models.NewOption(toolscanner.ProjectTypeUserTitle, toolscanner.ProjectTypeEnvKey)
	for _, projectType := range detectedProjectTypes {
		projectTypeOption.AddOption(projectType, &toolScannerOptionModel)
	}
	return *projectTypeOption
}
