package scanner

import (
	"fmt"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/fastlane"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/scanners/xamarin"
)

// ManualConfig ...
func ManualConfig() (models.ScanResultModel, error) {
	projectScanners := []scanners.ScannerInterface{
		new(android.Scanner),
		new(xamarin.Scanner),
		new(ios.Scanner),
		new(fastlane.Scanner),
	}

	projectTypeOptionMap := map[string]models.OptionModel{}
	projectTypeConfigMap := map[string]models.BitriseConfigMap{}

	for _, detector := range projectScanners {
		detectorName := detector.Name()

		option := detector.DefaultOptions()
		projectTypeOptionMap[detectorName] = option

		configs, err := detector.DefaultConfigs()
		if err != nil {
			return models.ScanResultModel{}, fmt.Errorf("Failed create default configs, error: %s", err)
		}
		projectTypeConfigMap[detectorName] = configs
	}

	customConfig, err := scanners.CustomConfig()
	if err != nil {
		return models.ScanResultModel{}, fmt.Errorf("Failed create default custom configs, error: %s", err)
	}

	projectTypeConfigMap["custom"] = customConfig

	return models.ScanResultModel{
		OptionsMap: projectTypeOptionMap,
		ConfigsMap: projectTypeConfigMap,
	}, nil
}
