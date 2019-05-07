package reactnative

import (
	"errors"
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

// scannerName is the name of the scanner.
const scannerName = "react-native"

const (
	// WorkDirInputKey is a key of the working directory step input.
	WorkDirInputKey = "workdir"
)

// Scanner implements the React Native project scanner.
type Scanner struct {
	searchDir      string
	iosScanner     *ios.Scanner
	androidScanner *android.Scanner

	hasNPMTest     bool
	packageJSONPth string

	usesExpo    bool
	usesExpoKit bool
}

// NewScanner creates a new scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name exposes the scanner name.
func (Scanner) Name() string {
	return scannerName
}

// isExpoProject reports whether a project is Expo based or not.
func isExpoProject(packageJSONPth string) (bool, error) {
	packages, err := utility.ParsePackagesJSON(packageJSONPth)
	if err != nil {
		return false, fmt.Errorf("failed to parse package json file (%s): %s", packageJSONPth, err)
	}

	if _, found := packages.Dependencies["expo"]; !found {
		return false, nil
	}

	// app.json file is a required part of an expo projects and shoulb be placed next to the root package.json file
	appJSONPth := filepath.Join(filepath.Dir(packageJSONPth), "app.json")
	exist, err := pathutil.IsPathExists(appJSONPth)
	if err != nil {
		return false, fmt.Errorf("failed to check if app.json file (%s) exist: %s", appJSONPth, err)
	}
	return exist, nil
}

// hasNativeProjects reports whether the project directory contains ios and android native projects or not.
func hasNativeProjects(packageJSONPth string, iosScanner *ios.Scanner, androidScanner *android.Scanner) (bool, bool, error) {
	projectDir := filepath.Dir(packageJSONPth)

	iosProjectDetected := false
	iosDir := filepath.Join(projectDir, "ios")
	if exist, err := pathutil.IsDirExists(iosDir); err != nil {
		return false, false, err
	} else if exist {
		if detected, err := iosScanner.DetectPlatform(projectDir); err != nil {
			return false, false, err
		} else if detected {
			iosProjectDetected = true
		}
	}

	androidProjectDetected := false
	androidDir := filepath.Join(projectDir, "android")
	if exist, err := pathutil.IsDirExists(androidDir); err != nil {
		return false, false, err
	} else if exist {
		if detected, err := androidScanner.DetectPlatform(projectDir); err != nil {
			return false, false, err
		} else if detected {
			androidProjectDetected = true
		}
	}

	return iosProjectDetected, androidProjectDetected, nil
}

// DetectPlatform implements the ScannerInterface.
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	scanner.searchDir = searchDir

	log.TInfof("Collect package.json files")

	packageJSONPths, err := CollectPackageJSONFiles(searchDir)
	if err != nil {
		return false, err
	}

	log.TPrintf("%d package.json file detected", len(packageJSONPths))
	log.TInfof("Filter relevant package.json files")

	usesExpo := false
	var packageFile string

	for _, packageJSONPth := range packageJSONPths {
		log.TPrintf("checking: %s", packageJSONPth)

		expo, err := isExpoProject(packageJSONPth)
		if err != nil {
			log.TWarnf("failed to check if project uses Expo: %s", err)
		} else {
			log.TDonef("project uses Expo: %s", expo)
		}

		if expo {
			usesExpo = true
			packageFile = packageJSONPth
			break
		}

		if scanner.iosScanner == nil {
			scanner.iosScanner = ios.NewScanner()
		}
		if scanner.androidScanner == nil {
			scanner.androidScanner = android.NewScanner()
		}

		ios, android, err := hasNativeProjects(packageJSONPth, scanner.iosScanner, scanner.androidScanner)
		if err != nil {
			log.TWarnf("failed to check native projects: %s", err)
		} else {
			log.TDonef("has native ios project: %s", ios)
			log.TDonef("has native android project: %s", android)
		}

		if ios || android {
			packageFile = packageJSONPth
			break
		}
	}

	if packageFile == "" {
		return false, nil
	}

	scanner.usesExpo = usesExpo
	scanner.packageJSONPth = packageFile

	return true, nil
}

// Options ...
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, error) {
	if scanner.usesExpo {
		return scanner.expoOptions()
	}

	warnings := models.Warnings{}

	var rootOption models.OptionNode

	// react options
	packages, err := utility.ParsePackagesJSON(scanner.packageJSONPth)
	if err != nil {
		return models.OptionNode{}, warnings, err
	}

	hasNPMTest := false
	if _, found := packages.Scripts["test"]; found {
		hasNPMTest = true
		scanner.hasNPMTest = true
	}

	projectDir := filepath.Dir(scanner.packageJSONPth)

	// android options
	var androidOptions *models.OptionNode
	androidDir := filepath.Join(projectDir, "android")
	if exist, err := pathutil.IsDirExists(androidDir); err != nil {
		return models.OptionNode{}, warnings, err
	} else if exist {
		androidScanner := android.NewScanner()

		if detected, err := androidScanner.DetectPlatform(scanner.searchDir); err != nil {
			return models.OptionNode{}, warnings, err
		} else if detected {
			// only the first match we need
			androidScanner.ExcludeTest = true
			androidScanner.ProjectRoots = []string{androidScanner.ProjectRoots[0]}

			npmCmd := command.New("npm", "install")
			npmCmd.SetDir(projectDir)
			if out, err := npmCmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
				return models.OptionNode{}, warnings, fmt.Errorf("failed to npm install react-native in: %s\noutput: %s\nerror: %s", projectDir, out, err)
			}

			options, warns, err := androidScanner.Options()
			warnings = append(warnings, warns...)
			if err != nil {
				return models.OptionNode{}, warnings, err
			}

			androidOptions = &options
			scanner.androidScanner = androidScanner
		}
	}

	// ios options
	var iosOptions *models.OptionNode
	iosDir := filepath.Join(projectDir, "ios")
	if exist, err := pathutil.IsDirExists(iosDir); err != nil {
		return models.OptionNode{}, warnings, err
	} else if exist {
		iosScanner := ios.NewScanner()

		if detected, err := iosScanner.DetectPlatform(scanner.searchDir); err != nil {
			return models.OptionNode{}, warnings, err
		} else if detected {
			options, warns, err := iosScanner.Options()
			warnings = append(warnings, warns...)
			if err != nil {
				return models.OptionNode{}, warnings, err
			}

			iosOptions = &options
			scanner.iosScanner = iosScanner
		}
	}

	if androidOptions == nil && iosOptions == nil {
		return models.OptionNode{}, warnings, errors.New("no ios nor android project detected")
	}
	// ---

	if androidOptions != nil {
		if iosOptions == nil {
			// we only found an android project
			// we need to update the config names
			lastChilds := androidOptions.LastChilds()
			for _, child := range lastChilds {
				for _, child := range child.ChildOptionMap {
					if child.Config == "" {
						return models.OptionNode{}, warnings, fmt.Errorf("no config for option: %s", child.String())
					}

					configName := configName(true, false, hasNPMTest)
					child.Config = configName
				}
			}
		} else {
			// we have both ios and android projects
			// we need to remove the android option's config names,
			// since ios options will hold them
			androidOptions.RemoveConfigs()
		}

		rootOption = *androidOptions
	}

	if iosOptions != nil {
		lastChilds := iosOptions.LastChilds()
		for _, child := range lastChilds {
			for _, child := range child.ChildOptionMap {
				if child.Config == "" {
					return models.OptionNode{}, warnings, fmt.Errorf("no config for option: %s", child.String())
				}

				configName := configName(scanner.androidScanner != nil, true, hasNPMTest)
				child.Config = configName
			}
		}

		if androidOptions == nil {
			// we only found an ios project
			rootOption = *iosOptions
		} else {
			// we have both ios and android projects
			// we attach ios options to the android options
			rootOption.AttachToLastChilds(iosOptions)
		}

	}

	return rootOption, warnings, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	expoOption := models.NewOption("Project was created using Expo CLI?", "")

	expoDefaultOptions := scanner.expoDefaultOptions()
	expoOption.AddOption("yes", &expoDefaultOptions)

	androidOptions := (&android.Scanner{ExcludeTest: true}).DefaultOptions()
	androidOptions.RemoveConfigs()
	expoOption.AddOption("no", &androidOptions)

	iosOptions := (&ios.Scanner{}).DefaultOptions()
	for _, child := range iosOptions.LastChilds() {
		for _, child := range child.ChildOptionMap {
			child.Config = defaultConfigName()
		}
	}

	androidOptions.AttachToLastChilds(&iosOptions)

	return *expoOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	if scanner.usesExpo {
		return scanner.expoConfigs()
	}

	configMap := models.BitriseConfigMap{}

	packageJSONDir := filepath.Dir(scanner.packageJSONPth)
	relPackageJSONDir, err := utility.RelPath(scanner.searchDir, packageJSONDir)
	if err != nil {
		return models.BitriseConfigMap{}, fmt.Errorf("Failed to get relative config.xml dir path, error: %s", err)
	}
	if relPackageJSONDir == "." {
		// config.xml placed in the search dir, no need to change-dir in the workflows
		relPackageJSONDir = ""
	}

	workdirEnvList := []envmanModels.EnvironmentItemModel{}
	if relPackageJSONDir != "" {
		workdirEnvList = append(workdirEnvList, envmanModels.EnvironmentItemModel{WorkDirInputKey: relPackageJSONDir})
	}

	if scanner.hasNPMTest {
		configBuilder := models.NewDefaultConfigBuilder()

		// ci
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "test"})...))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

		// cd
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))

		// android cd
		if scanner.androidScanner != nil {
			projectLocationEnv := "$" + android.ProjectLocationInputEnvKey

			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
				envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
			))
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
				envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: projectLocationEnv},
			))
		}

		// ios cd
		if scanner.iosScanner != nil {
			for _, descriptor := range scanner.iosScanner.ConfigDescriptors {
				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

				if descriptor.MissingSharedSchemes {
					configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.RecreateUserSchemesStepListItem(
						envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
					))
				}

				if descriptor.HasPodfile {
					configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CocoapodsInstallStepListItem())
				}

				if descriptor.CarthageCommand != "" {
					configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CarthageStepListItem(
						envmanModels.EnvironmentItemModel{ios.CarthageCommandInputKey: descriptor.CarthageCommand},
					))
				}

				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
					envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
				))

				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

				bitriseDataModel, err := configBuilder.Generate(scannerName)
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				data, err := yaml.Marshal(bitriseDataModel)
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				configName := configName(scanner.androidScanner != nil, true, true)
				configMap[configName] = string(data)
			}
		} else {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

			bitriseDataModel, err := configBuilder.Generate(scannerName)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			data, err := yaml.Marshal(bitriseDataModel)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			configName := configName(scanner.androidScanner != nil, false, true)
			configMap[configName] = string(data)
		}
	} else {
		configBuilder := models.NewDefaultConfigBuilder()

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(append(workdirEnvList, envmanModels.EnvironmentItemModel{"command": "install"})...))

		if scanner.androidScanner != nil {
			projectLocationEnv := "$" + android.ProjectLocationInputEnvKey

			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
				envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
			))
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
				envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: projectLocationEnv},
			))
		}

		if scanner.iosScanner != nil {
			for _, descriptor := range scanner.iosScanner.ConfigDescriptors {
				configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

				if descriptor.MissingSharedSchemes {
					configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.RecreateUserSchemesStepListItem(
						envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
					))
				}

				if descriptor.HasPodfile {
					configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CocoapodsInstallStepListItem())
				}

				if descriptor.CarthageCommand != "" {
					configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.CarthageStepListItem(
						envmanModels.EnvironmentItemModel{ios.CarthageCommandInputKey: descriptor.CarthageCommand},
					))
				}

				configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.XcodeArchiveStepListItem(
					envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
					envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
				))

				configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

				bitriseDataModel, err := configBuilder.Generate(scannerName)
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				data, err := yaml.Marshal(bitriseDataModel)
				if err != nil {
					return models.BitriseConfigMap{}, err
				}

				configName := configName(scanner.androidScanner != nil, true, false)
				configMap[configName] = string(data)
			}
		} else {
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

			bitriseDataModel, err := configBuilder.Generate(scannerName)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			data, err := yaml.Marshal(bitriseDataModel)
			if err != nil {
				return models.BitriseConfigMap{}, err
			}

			configName := configName(scanner.androidScanner != nil, false, false)
			configMap[configName] = string(data)
		}
	}

	return configMap, nil
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	// ci
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{"command": "install"}))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{"command": "test"}))
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

	// cd
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.NpmStepListItem(envmanModels.EnvironmentItemModel{"command": "install"}))

	// android
	projectLocationEnv := "$" + android.ProjectLocationInputEnvKey

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: projectLocationEnv},
	))

	// ios
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
		envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

	bitriseDataModel, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(bitriseDataModel)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	configName := defaultConfigName()
	configMap := models.BitriseConfigMap{
		configName: string(data),
	}
	return configMap, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{
		string(ios.XcodeProjectTypeIOS),
		string(ios.XcodeProjectTypeMacOS),
		android.ScannerName,
	}
}
