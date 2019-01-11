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
	Warnings             models.Warnings
	Errors               models.Errors
	OptionModel          models.OptionNode
	ConfigMap            models.BitriseConfigMap
	ExcludedScannerNames []string
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
		var excludedScannerNames []string
		projectScannerWarnings, projectScannerMatchResults := mapScannersToOutput(scanners.ProjectScanners, searchDir, excludedScannerNames)
		var matchingProjectTypes []string
		for scannerName := range projectScannerMatchResults {
			matchingProjectTypes = append(matchingProjectTypes, scannerName)
		}

		toolScannerWarnings, toolScannerResults := mapScannersToOutput(scanners.AutomationToolScanners, searchDir, excludedScannerNames)
		// Add project_type property option to tool scanner's as they do not dertect project/platform
		for _, detectResult := range toolScannerResults {
			detectResult.OptionModel = toolscanner.AddProjectTypeToToolScanner(detectResult.OptionModel, matchingProjectTypes)
		}

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
		scannerToOptions[k] = v.OptionModel
		scannerToConfigMap[k] = v.ConfigMap
		scannerToWarnings[k] = v.Warnings
		if v.Errors != nil {
			scannerToErrors[k] = v.Errors
		}
	}
	return models.ScanResultModel{
		ScannerToOptionRoot:       scannerToOptions,
		ScannerToBitriseConfigMap: scannerToConfigMap,
		ScannerToWarnings:         scannerToWarnings,
		ScannerToErrors:           scannerToErrors,
	}
}

func mapScannersToOutput(scannerList []scanners.ScannerInterface, searchDir string, excludedScannerNames []string) (map[string]models.Warnings, map[string]scannerDetectResult) {
	scannerToMatchResult := map[string]scannerDetectResult{}
	scannerToWarnings := map[string]models.Warnings{}
	for _, scanner := range scannerList {
		warnings, matchResult := checkScannerMatchAndReturnOutput(scanner, searchDir, excludedScannerNames)
		if warnings != nil {
			scannerToWarnings[scanner.Name()] = *warnings
		}
		if matchResult != nil {
			scannerToMatchResult[scanner.Name()] = *matchResult
			excludedScannerNames = append(excludedScannerNames, (*matchResult).ExcludedScannerNames...)
		}
	}
	return scannerToWarnings, scannerToMatchResult
}

func checkScannerMatchAndReturnOutput(detector scanners.ScannerInterface, searchDir string, excludedScannerNamesPrevious []string) (*models.Warnings, *scannerDetectResult) {
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
		return nil, &scannerDetectResult{
			Warnings: detectorWarnings,
			Errors:   detectorErrors,
		}
	}

	// Generate configs
	configs, err := detector.Configs()
	if err != nil {
		log.TErrorf("Failed to generate config, error: %s", err)
		detectorErrors = append(detectorErrors, err.Error())
		return nil, &scannerDetectResult{
			Warnings: detectorWarnings,
			Errors:   detectorErrors,
		}
	}

	log.TPrintf("|                                                                              |")
	log.TPrintf("+------------------------------------------------------------------------------+")

	excludedScannerNamesCurrent := detector.ExcludedScannerNames()
	if len(excludedScannerNamesCurrent) > 0 {
		log.TWarnf("Scanner will exclude scanners: %v", excludedScannerNamesCurrent)
		excludedScannerNamesCurrent = append(excludedScannerNamesPrevious, excludedScannerNamesCurrent...)
	}

	fmt.Println()
	return &models.Warnings{}, &scannerDetectResult{
		Warnings:             detectorWarnings,
		Errors:               detectorErrors,
		OptionModel:          options,
		ConfigMap:            configs,
		ExcludedScannerNames: excludedScannerNamesCurrent,
	}
}
