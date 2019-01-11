package scanner

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners"
	"github.com/bitrise-core/bitrise-init/toolscanner"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
)

type scannerDetectResult struct {
	warnings         models.Warnings
	errors           models.Errors
	optionModel      models.OptionNode
	configMap        models.BitriseConfigMap
	excludedScanners []string
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
	log.TInfof(colorstring.Blue("Running scanners:"))
	fmt.Println()

	var scannerToNoDetectWarnings map[string]models.Warnings
	var scannerToDetectResults map[string]scannerDetectResult
	{
		projectScannerWarnings, projectScannerMatchResults := mapScannerOutput(scanners.ProjectScanners, searchDir)
		detectedProjectTypes := make([]string, 0, len(projectScannerMatchResults))
		for scannerKey := range projectScannerMatchResults {
			detectedProjectTypes = append(detectedProjectTypes, scannerKey)
		}
		log.Printf("Detected project types: %s", detectedProjectTypes)
		fmt.Println()

		toolScannerWarnings, toolScannerResults := mapScannerOutput(scanners.AutomationToolScanners, searchDir)
		detectedAutomationToolScanners := make([]string, 0, len(toolScannerResults))
		for scannerKey := range toolScannerResults {
			detectedAutomationToolScanners = append(detectedAutomationToolScanners, scannerKey)
		}
		log.Printf("Detected automation tools: %s", detectedAutomationToolScanners)
		fmt.Println()

		// Add project_type property option to tool scanner's as they do not detect project/platform
		toolScannerResults = mapAddProjectType(toolScannerResults, detectedProjectTypes)

		for k, v := range toolScannerWarnings {
			projectScannerWarnings[k] = v
		}
		scannerToNoDetectWarnings = projectScannerWarnings
		for k, v := range toolScannerResults {
			projectScannerMatchResults[k] = v
		}
		scannerToDetectResults = projectScannerMatchResults
	}

	scannerToOptions := map[string]models.OptionNode{}
	scannerToConfigMap := map[string]models.BitriseConfigMap{}
	scannerToWarnings := scannerToNoDetectWarnings
	scannerToErrors := map[string]models.Errors{}
	for k, v := range scannerToDetectResults {
		scannerToOptions[k] = v.optionModel
		scannerToConfigMap[k] = v.configMap
		scannerToWarnings[k] = v.warnings
		if v.errors != nil {
			scannerToErrors[k] = v.errors
		}
	}
	return models.ScanResultModel{
		ScannerToOptionRoot:       scannerToOptions,
		ScannerToBitriseConfigMap: scannerToConfigMap,
		ScannerToWarnings:         scannerToWarnings,
		ScannerToErrors:           scannerToErrors,
	}
}

func mapScannerOutput(scannerList []scanners.ScannerInterface, searchDir string) (map[string]models.Warnings, map[string]scannerDetectResult) {
	scannerToDetectResult := map[string]scannerDetectResult{}
	scannerToNoDetectWarnings := map[string]models.Warnings{}
	var excludedScannerNames []string
	for _, scanner := range scannerList {
		log.TInfof("Scanner: %s", colorstring.Blue(scanner.Name()))
		if sliceutil.IsStringInSlice(scanner.Name(), excludedScannerNames) {
			log.TWarnf("scanner is marked as excluded, skipping...")
			fmt.Println()
			continue
		}

		log.TPrintf("+------------------------------------------------------------------------------+")
		log.TPrintf("|                                                                              |")

		warnings, matchResult := checkScannerDetectAndReturnOutput(scanner, searchDir)

		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()

		if warnings != nil {
			scannerToNoDetectWarnings[scanner.Name()] = *warnings
		}
		if matchResult != nil {
			scannerToDetectResult[scanner.Name()] = *matchResult
			excludedScannerNames = append(excludedScannerNames, (*matchResult).excludedScanners...)
		}
	}
	return scannerToNoDetectWarnings, scannerToDetectResult
}

func checkScannerDetectAndReturnOutput(detector scanners.ScannerInterface, searchDir string) (*models.Warnings, *scannerDetectResult) {
	var detectorWarnings models.Warnings
	var detectorErrors []string

	if detected, err := detector.DetectPlatform(searchDir); err != nil {
		log.TErrorf("Scanner failed, error: %s", err)
		return &models.Warnings{err.Error()}, nil
	} else if !detected {
		return nil, nil
	}

	options, projectWarnings, err := detector.Options()
	detectorWarnings = append(detectorWarnings, projectWarnings...)

	if err != nil {
		log.TErrorf("Analyzer failed, error: %s", err)
		detectorWarnings = append(detectorWarnings, err.Error())
		return nil, &scannerDetectResult{
			warnings: detectorWarnings,
			errors:   detectorErrors,
		}
	}

	// Generate configs
	configs, err := detector.Configs()
	if err != nil {
		log.TErrorf("Failed to generate config, error: %s", err)
		detectorErrors = append(detectorErrors, err.Error())
		return nil, &scannerDetectResult{
			warnings: detectorWarnings,
			errors:   detectorErrors,
		}
	}

	scannerExcludedScanners := detector.ExcludedScannerNames()
	if len(scannerExcludedScanners) > 0 {
		log.TWarnf("Scanner will exclude scanners: %v", scannerExcludedScanners)
	}

	return &models.Warnings{}, &scannerDetectResult{
		warnings:         detectorWarnings,
		errors:           detectorErrors,
		optionModel:      options,
		configMap:        configs,
		excludedScanners: scannerExcludedScanners,
	}
}

func mapAddProjectType(toolScannerResults map[string]scannerDetectResult, detectedProjectTypes []string) map[string]scannerDetectResult {
	toolScannerResultsWithProjectType := map[string]scannerDetectResult{}
	for scannerKey, detectResult := range toolScannerResults {
		detectResult.optionModel = toolscanner.AddProjectTypeToToolScanner(detectResult.optionModel, detectedProjectTypes)
		toolScannerResultsWithProjectType[scannerKey] = detectResult
	}
	return toolScannerResultsWithProjectType
}
