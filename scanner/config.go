package scanner

import (
	"fmt"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners"
	"github.com/bitrise-core/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/log"
)

// Config ...
func Config(searchDir string) models.ScanResultModel {
	//
	// Scan
	projectScanners := scanners.ActiveScanners
	ignoreScanners := []string{}

	projectTypeErrorMap := map[string]models.Errors{}
	projectTypeWarningMap := map[string]models.Warnings{}
	projectTypeOptionMap := map[string]models.OptionModel{}
	projectTypeConfigMap := map[string]models.BitriseConfigMap{}

	log.Infoft(colorstring.Blue("Running scanners:"))
	fmt.Println()

	for _, detector := range projectScanners {
		detectorName := detector.Name()
		detectorWarnings := []string{}
		detectorErrors := []string{}

		log.Infoft("Scanner: %s", colorstring.Blue(detectorName))
		log.Printft("+------------------------------------------------------------------------------+")
		log.Printft("|                                                                              |")

		if utility.IgnoreScanner(detectorName, ignoreScanners) {
			log.Warnft("%s scanner is marked to be ignored", detectorName)
			log.Printft("|                                                                              |")
			log.Printft("+------------------------------------------------------------------------------+")
			fmt.Println()
			continue
		}

		detected, err := detector.DetectPlatform(searchDir)
		if err != nil {
			log.Errorft("Scanner failed, error: %s", err)
			detectorWarnings = append(detectorWarnings, err.Error())
			projectTypeWarningMap[detectorName] = detectorWarnings
			detected = false
		}

		if !detected {
			log.Printft("|                                                                              |")
			log.Printft("+------------------------------------------------------------------------------+")
			fmt.Println()
			continue
		}

		ignoreScanners = append(ignoreScanners, detector.IgnoreScanners()...)

		options, projectWarnings, err := detector.Options()
		detectorWarnings = append(detectorWarnings, projectWarnings...)

		if err != nil {
			log.Errorft("Analyzer failed, error: %s", err)
			detectorWarnings = append(detectorWarnings, err.Error())
			projectTypeWarningMap[detectorName] = detectorWarnings

			log.Printft("|                                                                              |")
			log.Printft("+------------------------------------------------------------------------------+")
			fmt.Println()
			continue
		}

		projectTypeWarningMap[detectorName] = detectorWarnings
		projectTypeOptionMap[detectorName] = options

		// Generate configs
		configs, err := detector.Configs()
		if err != nil {
			log.Errorft("Failed to generate config, error: %s", err)
			detectorErrors = append(detectorErrors, err.Error())
			projectTypeErrorMap[detectorName] = detectorErrors
			continue
		}

		projectTypeConfigMap[detectorName] = configs

		log.Printft("|                                                                              |")
		log.Printft("+------------------------------------------------------------------------------+")
		fmt.Println()
	}
	// ---

	return models.ScanResultModel{
		OptionsMap:  projectTypeOptionMap,
		ConfigsMap:  projectTypeConfigMap,
		WarningsMap: projectTypeWarningMap,
	}
}
