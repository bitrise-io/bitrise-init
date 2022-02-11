package reactnative

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/scanners/android"
	"github.com/bitrise-io/bitrise-init/scanners/ios"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/bitrise-init/utility"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigName = "default-react-native-config"
)

type configDescriptor struct {
	hasIOS, hasAndroid bool
	hasTest            bool
	ios                ios.ConfigDescriptor
}

func (d configDescriptor) configName() string {
	name := "react-native"
	if d.hasAndroid {
		name += "-android"
	}
	if d.hasIOS {
		name += "-ios"
		if d.ios.MissingSharedSchemes {
			name += "-missing-shared-schemes"
		}
		if d.ios.HasPodfile {
			name += "-pod"
		}
		if d.ios.CarthageCommand != "" {
			name += "-carthage"
		}
	}
	if d.hasTest {
		name += "-test"
	}

	return name + "-config"
}

func generateIOSOptions(result ios.DetectResult, hasAndroid, hasTests bool) (*models.OptionNode, models.Warnings, []configDescriptor) {
	var (
		warnings    models.Warnings
		descriptors []configDescriptor
	)

	projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputSummary, ios.ProjectPathInputEnvKey, models.TypeSelector)
	for _, project := range result.Projects {
		warnings = append(warnings, project.Warnings...)

		schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputSummary, ios.SchemeInputEnvKey, models.TypeSelector)
		projectPathOption.AddOption(project.RelPath, schemeOption)

		for _, scheme := range project.Schemes {
			exportMethodOption := models.NewOption(ios.DistributionMethodInputTitle, ios.DistributionMethodInputSummary, ios.DistributionMethodEnvKey, models.TypeSelector)
			schemeOption.AddOption(scheme.Name, exportMethodOption)

			for _, exportMethod := range ios.IosExportMethods {
				iosConfig := ios.NewConfigDescriptor(project.IsPodWorkspace, project.CarthageCommand, scheme.HasXCTests, scheme.HasAppClip, exportMethod, scheme.Missing)
				descriptor := configDescriptor{
					hasIOS:     true,
					hasAndroid: hasAndroid,
					hasTest:    hasTests,
					ios:        iosConfig,
				}
				descriptors = append(descriptors, descriptor)

				exportMethodOption.AddConfig(exportMethod, models.NewConfigOption(descriptor.configName(), nil))
			}
		}
	}

	return projectPathOption, warnings, descriptors
}

// options implements ScannerInterface.Options function for plain React Native projects.
func (scanner *Scanner) options() (models.OptionNode, models.Warnings) {
	var (
		rootOption     models.OptionNode
		allDescriptors []configDescriptor
		warnings       = models.Warnings{}
	)

	// Android
	if len(scanner.androidProjects) > 0 {
		androidOptions := models.NewOption(android.ProjectLocationInputTitle, android.ProjectLocationInputSummary, android.ProjectLocationInputEnvKey, models.TypeSelector)
		rootOption = *androidOptions

		for _, project := range scanner.androidProjects {
			warnings = append(warnings, project.Warnings...)

			moduleOption := models.NewOption(android.ModuleInputTitle, android.ModuleInputSummary, android.ModuleInputEnvKey, models.TypeUserInput)
			variantOption := models.NewOption(android.VariantInputTitle, android.VariantInputSummary, android.VariantInputEnvKey, models.TypeOptionalUserInput)

			androidOptions.AddOption(project.RelPath, moduleOption)
			moduleOption.AddOption("app", variantOption)

			if len(scanner.iosProjects.Projects) == 0 {
				descriptor := configDescriptor{
					hasAndroid: true,
					hasTest:    scanner.hasTest,
				}
				allDescriptors = append(allDescriptors, descriptor)

				variantOption.AddConfig("", models.NewConfigOption(descriptor.configName(), nil))

				continue
			}

			iosOptions, iosWarnings, descriptors := generateIOSOptions(scanner.iosProjects, true, scanner.hasTest)
			warnings = append(warnings, iosWarnings...)
			allDescriptors = append(allDescriptors, descriptors...)

			variantOption.AddOption("", iosOptions)
		}
	} else {
		options, iosWarnings, descriptors := generateIOSOptions(scanner.iosProjects, false, scanner.hasTest)
		rootOption = *options
		warnings = append(warnings, iosWarnings...)
		allDescriptors = descriptors
	}

	scanner.configDescriptors = removeDuplicatedConfigDescriptors(allDescriptors)

	return rootOption, warnings
}

// defaultOptions implements ScannerInterface.DefaultOptions function for plain React Native projects.
func (scanner *Scanner) defaultOptions() models.OptionNode {
	androidOptions := (&android.Scanner{}).DefaultOptions()
	androidOptions.RemoveConfigs()

	iosOptions := (&ios.Scanner{}).DefaultOptions()
	for _, child := range iosOptions.LastChilds() {
		for _, child := range child.ChildOptionMap {
			child.Config = defaultConfigName
		}
	}

	androidOptions.AttachToLastChilds(&iosOptions)

	return androidOptions
}

// configs implements ScannerInterface.Configs function for plain React Native projects.
func (scanner *Scanner) configs(isPrivateRepo bool) (models.BitriseConfigMap, error) {
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

	for _, descriptor := range scanner.configDescriptors {
		configBuilder := models.NewDefaultConfigBuilder()

		// ci
		primaryDescription := primaryWorkflowNoTestsDescription
		if descriptor.hasTest {
			primaryDescription = primaryWorkflowDescription
		}

		configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, primaryDescription)
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
			ShouldIncludeCache:       false,
			ShouldIncludeActivateSSH: isPrivateRepo,
		})...)
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, scanner.getTestSteps(relPackageJSONDir)...)

		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepListV2(false)...)

		// cd
		configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
			ShouldIncludeCache:       false,
			ShouldIncludeActivateSSH: isPrivateRepo,
		})...)
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, scanner.getTestSteps(relPackageJSONDir)...)

		// android cd
		if descriptor.hasAndroid {
			projectLocationEnv := "$" + android.ProjectLocationInputEnvKey

			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
				envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
			))
			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
				envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: projectLocationEnv},
			))
		}

		// ios cd
		if descriptor.hasIOS {
			if descriptor.ios.MissingSharedSchemes {
				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.RecreateUserSchemesStepListItem(
					envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
				))
			}

			if descriptor.ios.HasPodfile {
				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CocoapodsInstallStepListItem())
			}

			if descriptor.ios.CarthageCommand != "" {
				configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CarthageStepListItem(
					envmanModels.EnvironmentItemModel{ios.CarthageCommandInputKey: descriptor.ios.CarthageCommand},
				))
			}

			configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
				envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
				envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
				envmanModels.EnvironmentItemModel{ios.DistributionMethodInputKey: "$" + ios.DistributionMethodEnvKey},
				envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
				envmanModels.EnvironmentItemModel{ios.AutomaticCodeSigningInputKey: ios.AutomaticCodeSigningInputAPIKeyValue},
			))
		}

		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepListV2(false)...)

		bitriseDataModel, err := configBuilder.Generate(scannerName)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(bitriseDataModel)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		configMap[descriptor.configName()] = string(data)
	}

	return configMap, nil
}

// defaultConfigs implements ScannerInterface.DefaultConfigs function for plain React Native projects.
func (scanner *Scanner) defaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	// primary
	configBuilder.SetWorkflowDescriptionTo(models.PrimaryWorkflowID, primaryWorkflowDescription)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
		ShouldIncludeCache:       false,
		ShouldIncludeActivateSSH: true,
	})...)
	// Assuming project uses yarn and has tests
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, getTestSteps("", true, true)...)
	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepListV2(false)...)

	// deploy
	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepListV2(steps.PrepareListParams{
		ShouldIncludeCache:       false,
		ShouldIncludeActivateSSH: true,
	})...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, getTestSteps("", true, true)...)

	// android
	projectLocationEnv := "$" + android.ProjectLocationInputEnvKey

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{android.GradlewPathInputKey: "$" + android.ProjectLocationInputEnvKey + "/gradlew"},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{android.ProjectLocationInputKey: projectLocationEnv},
	))

	// ios
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
		envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.DistributionMethodInputKey: "$" + ios.DistributionMethodEnvKey},
		envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: "Release"},
		envmanModels.EnvironmentItemModel{ios.AutomaticCodeSigningInputKey: ios.AutomaticCodeSigningInputAPIKeyValue},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepListV2(false)...)

	bitriseDataModel, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(bitriseDataModel)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	configName := defaultConfigName
	configMap := models.BitriseConfigMap{
		configName: string(data),
	}

	return configMap, nil
}

func getTestSteps(workDir string, hasYarnLockFile, hasTest bool) []bitriseModels.StepListItemModel {
	var testSteps []bitriseModels.StepListItemModel

	if hasYarnLockFile {
		testSteps = append(testSteps, steps.YarnStepListItem("install", workDir))
		if hasTest {
			testSteps = append(testSteps, steps.YarnStepListItem("test", workDir))
		}
	} else {
		testSteps = append(testSteps, steps.NpmStepListItem("install", workDir))
		if hasTest {
			testSteps = append(testSteps, steps.NpmStepListItem("test", workDir))
		}
	}

	return testSteps
}

func (scanner *Scanner) getTestSteps(workDir string) []bitriseModels.StepListItemModel {
	return getTestSteps(workDir, scanner.hasYarnLockFile, scanner.hasTest)
}

func removeDuplicatedConfigDescriptors(configDescriptors []configDescriptor) []configDescriptor {
	descritorNameMap := map[string]configDescriptor{}
	for _, descriptor := range configDescriptors {
		name := descriptor.configName()
		descritorNameMap[name] = descriptor
	}

	descriptors := []configDescriptor{}
	for _, descriptor := range descritorNameMap {
		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}
