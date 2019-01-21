package scanner

import (
	"fmt"
	"os"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
)

const otherProjectType = "other"

type scanResultStatus string

const (
	scanResultNotDetected        scanResultStatus = "scanResultNotDetected"
	scanResultDetectedWithErrors scanResultStatus = "scanResultDetectedWithErrors"
	scanResultDetected           scanResultStatus = "scanResultDetected"
)

type scannerOutput struct {
	scanResult scanResultStatus

	warnings models.Warnings

	errors models.Errors

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

	scannerToDetectResults := map[string]scannerOutput{}
	{
		projectScannerMatchResults := mapScannerOutput(scanners.ProjectScanners, searchDir)
		detectedProjectTypes := make([]string, 0, len(projectScannerMatchResults))
		for scannerKey, scannerOutput := range projectScannerMatchResults {
			if scannerOutput.scanResult == scanResultDetected {
				detectedProjectTypes = append(detectedProjectTypes, scannerKey)
			}
		}
		log.Printf("Detected project types: %s", detectedProjectTypes)
		fmt.Println()

		if len(detectedProjectTypes) == 0 {
			detectedProjectTypes = []string{otherProjectType}
		}
		for _, toolScanner := range scanners.AutomationToolScanners {
			toolScanner.(scanners.AutomationToolScanner).SetDetectedProjectTypes(detectedProjectTypes)
		}

		toolScannerResults := mapScannerOutput(scanners.AutomationToolScanners, searchDir)
		detectedAutomationToolScanners := make([]string, 0, len(toolScannerResults))
		for scannerKey, scannerOutput := range toolScannerResults {
			if scannerOutput.scanResult == scanResultDetected {
				detectedAutomationToolScanners = append(detectedAutomationToolScanners, scannerKey)
			}
		}
		log.Printf("Detected automation tools: %s", detectedAutomationToolScanners)
		fmt.Println()

		scannerToDetectResults = toolScannerResults
		for k, v := range projectScannerMatchResults {
			scannerToDetectResults[k] = v
		}
	}

	scannerToWarnings := map[string]models.Warnings{}
	scannerToErrors := map[string]models.Errors{}
	scannerToOptions := map[string]models.OptionNode{}
	scannerToConfigMap := map[string]models.BitriseConfigMap{}
	for k, v := range scannerToDetectResults {
		if v.scanResult == scanResultNotDetected && v.warnings != nil ||
			v.scanResult != scanResultNotDetected {
			scannerToWarnings[k] = v.warnings
		}
		if v.scanResult == scanResultDetected || v.scanResult == scanResultDetectedWithErrors {
			if v.errors != nil {
				scannerToErrors[k] = v.errors
			}
		}
		if v.scanResult == scanResultDetected {
			if v.configMap != nil {
				scannerToOptions[k] = v.optionModel
				scannerToConfigMap[k] = v.configMap
			}
		}
	}
	return models.ScanResultModel{
		ScannerToOptionRoot:       scannerToOptions,
		ScannerToBitriseConfigMap: scannerToConfigMap,
		ScannerToWarnings:         scannerToWarnings,
		ScannerToErrors:           scannerToErrors,
	}
}

func mapScannerOutput(scannerList []scanners.ScannerInterface, searchDir string) map[string]scannerOutput {
	scannerOutputs := map[string]scannerOutput{}
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
		scannerOutput := runScanner(scanner, searchDir)
		log.TPrintf("|                                                                              |")
		log.TPrintf("+------------------------------------------------------------------------------+")
		fmt.Println()

		scannerOutputs[scanner.Name()] = scannerOutput
		excludedScannerNames = append(excludedScannerNames, scannerOutput.excludedScanners...)
	}
	return scannerOutputs
}

func runScanner(detector scanners.ScannerInterface, searchDir string) scannerOutput {
	var detectorWarnings models.Warnings
	var detectorErrors []string

	if detected, err := detector.DetectPlatform(searchDir); err != nil {
		log.TErrorf("Scanner failed, error: %s", err)
		return scannerOutput{
			scanResult: scanResultNotDetected,
			warnings:   models.Warnings{err.Error()},
		}
	} else if !detected {
		return scannerOutput{
			scanResult: scanResultNotDetected,
		}
	}

	options, projectWarnings, err := detector.Options()
	detectorWarnings = append(detectorWarnings, projectWarnings...)

	if err != nil {
		log.TErrorf("Analyzer failed, error: %s", err)
		detectorWarnings = append(detectorWarnings, err.Error())
		return scannerOutput{
			scanResult: scanResultDetectedWithErrors,
			warnings:   detectorWarnings,
			errors:     detectorErrors,
		}
	}

	// Generate configs
	configs, err := detector.Configs()
	if err != nil {
		log.TErrorf("Failed to generate config, error: %s", err)
		detectorErrors = append(detectorErrors, err.Error())
		return scannerOutput{
			scanResult: scanResultDetectedWithErrors,
			warnings:   detectorWarnings,
			errors:     detectorErrors,
		}
	}

	scannerExcludedScanners := detector.ExcludedScannerNames()
	if len(scannerExcludedScanners) > 0 {
		log.TWarnf("Scanner will exclude scanners: %v", scannerExcludedScanners)
	}

	return scannerOutput{
		scanResult:       scanResultDetected,
		warnings:         detectorWarnings,
		errors:           detectorErrors,
		optionModel:      options,
		configMap:        configs,
		excludedScanners: scannerExcludedScanners,
	}
}
