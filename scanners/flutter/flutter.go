package flutter

import (
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/scanners/android"
	"github.com/bitrise-core/bitrise-init/scanners/cordova"
	"github.com/bitrise-core/bitrise-init/scanners/ios"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-core/bitrise-init/utility"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xcode-project/xcworkspace"
	yaml "gopkg.in/yaml.v2"
)

const (
	scannerName                = "flutter"
	defaultConfigName          = "default-flutter-config"
	projectLocationInputKey    = "project_location"
	defaultIOSConfiguration    = "Release"
	projectLocationInputEnvKey = "BITRISE_FLUTTER_PROJECT_LOCATION"
	projectLocationInputTitle  = "Project Location"
)

//------------------
// ScannerInterface
//------------------

// Scanner ...
type Scanner struct {
	projectLocations  []string
	xcodeProjectPaths []string
	sharedSchemes     map[string][]string
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return scannerName
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	log.TPrintf("Looking for pubspec.yaml files")
	{
		fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
		if err != nil {
			return false, err
		}

		filters := []utility.FilterFunc{
			utility.BaseFilter("pubspec.yaml", true),
		}

		pubspecLocations, err := utility.FilterPaths(fileList, filters...)
		if err != nil {
			return false, err
		}

		for _, path := range pubspecLocations {
			scanner.projectLocations = append(scanner.projectLocations, filepath.Dir(path))
		}

		if len(scanner.projectLocations) == 0 {
			log.TErrorf("Couldn't find pubspec.yaml files")
			return false, nil
		}
		log.TDonef("Found")
		log.TPrintf("")
	}

	log.TPrintf("Looking for iOS workspace files")
	{
		fileList, err := utility.ListPathInDirSortedByComponents(searchDir, true)
		if err != nil {
			return false, err
		}

		filters := []utility.FilterFunc{
			ios.AllowXCWorkspaceExtFilter,
			ios.AllowIsDirectoryFilter,
			ios.ForbidEmbeddedWorkspaceRegexpFilter,
			ios.ForbidGitDirComponentFilter,
			ios.ForbidPodsDirComponentFilter,
			ios.ForbidCarthageDirComponentFilter,
			ios.ForbidFramworkComponentWithExtensionFilter,
			ios.ForbidCordovaLibDirComponentFilter,
			ios.ForbidNodeModulesComponentFilter,
		}

		xcodeProjectPaths, err := utility.FilterPaths(fileList, filters...)
		if err != nil {
			return false, err
		}

		scanner.xcodeProjectPaths = xcodeProjectPaths

		if len(scanner.xcodeProjectPaths) == 0 {
			log.TErrorf("Couldn't find workspace files")
			return false, nil
		}
		log.TDonef("Found")
		log.TPrintf("")
	}

	log.TPrintf("Looking for iOS shared schemes")
	{
		scanner.sharedSchemes = map[string][]string{}
		for _, workspacePath := range scanner.xcodeProjectPaths {
			ws, err := xcworkspace.Open(workspacePath)
			if err != nil {
				log.TErrorf("Couldn't open workspace(%s), error: %s", workspacePath, err)
				return false, nil
			}
			schemeMap, err := ws.Schemes()
			if err != nil {
				log.TErrorf("Couldn't find schemes in workspace(%s), error: %s", workspacePath, err)
				return false, nil
			}

			for _, schemes := range schemeMap {
				for _, scheme := range schemes {
					scanner.sharedSchemes[workspacePath] = append(scanner.sharedSchemes[workspacePath], scheme.Name)
				}
			}
		}

		if len(scanner.sharedSchemes) == 0 {
			log.TErrorf("Couldn't find schemes")
			return false, nil
		}

		for wsPath, s := range scanner.sharedSchemes {
			if len(s) == 0 {
				log.TErrorf("Couldn't find scheme in: %s", wsPath)
				return false, nil
			}
		}

		log.TDonef("Found")
		log.TPrintf("")
	}

	log.TDonef("Detected")

	return true, nil
}

// ExcludedScannerNames ...
func (Scanner) ExcludedScannerNames() []string {
	return []string{
		string(ios.XcodeProjectTypeIOS),
		string(ios.XcodeProjectTypeMacOS),
		cordova.ScannerName,
		android.ScannerName,
	}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	flutterProjectLocationOption := models.NewOption(projectLocationInputTitle, projectLocationInputEnvKey)

	for _, location := range scanner.projectLocations {
		projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
		flutterProjectLocationOption.AddOption(location, projectPathOption)

		for _, xcodeWorkspacePath := range scanner.xcodeProjectPaths {
			schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)
			projectPathOption.AddOption(xcodeWorkspacePath, schemeOption)

			for _, scheme := range scanner.sharedSchemes[xcodeWorkspacePath] {
				exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
				schemeOption.AddOption(scheme, exportMethodOption)

				for _, exportMethod := range ios.IosExportMethods {
					configOption := models.NewConfigOption(defaultConfigName)
					exportMethodOption.AddConfig(exportMethod, configOption)
				}
			}
		}
	}

	return *flutterProjectLocationOption, nil, nil
}

// DefaultOptions ...
func (Scanner) DefaultOptions() models.OptionModel {
	flutterProjectLocationOption := models.NewOption(projectLocationInputTitle, projectLocationInputEnvKey)

	projectPathOption := models.NewOption(ios.ProjectPathInputTitle, ios.ProjectPathInputEnvKey)
	flutterProjectLocationOption.AddOption("_", projectPathOption)

	schemeOption := models.NewOption(ios.SchemeInputTitle, ios.SchemeInputEnvKey)
	projectPathOption.AddOption("_", schemeOption)

	exportMethodOption := models.NewOption(ios.IosExportMethodInputTitle, ios.ExportMethodInputEnvKey)
	schemeOption.AddOption("_", exportMethodOption)

	for _, exportMethod := range ios.IosExportMethods {
		configOption := models.NewConfigOption(defaultConfigName)
		exportMethodOption.AddConfig(exportMethod, configOption)
	}

	return *flutterProjectLocationOption
}

// Configs ...
func (scanner *Scanner) Configs() (models.BitriseConfigMap, error) {
	return scanner.DefaultConfigs()
}

// DefaultConfigs ...
func (Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	// primary

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterInstallStepListItem())

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.FlutterTestStepListItem(
		envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
	))

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(false)...)

	// deploy

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(false)...)

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.CertificateAndProfileInstallerStepListItem())

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterInstallStepListItem())

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterTestStepListItem(
		envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.FlutterBuildStepListItem(
		envmanModels.EnvironmentItemModel{projectLocationInputKey: "$" + projectLocationInputEnvKey},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.XcodeArchiveStepListItem(
		envmanModels.EnvironmentItemModel{ios.ProjectPathInputKey: "$" + ios.ProjectPathInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.SchemeInputKey: "$" + ios.SchemeInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ExportMethodInputKey: "$" + ios.ExportMethodInputEnvKey},
		envmanModels.EnvironmentItemModel{ios.ConfigurationInputKey: defaultIOSConfiguration},
	))

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(false)...)

	//

	config, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		defaultConfigName: string(data),
	}, nil
}
