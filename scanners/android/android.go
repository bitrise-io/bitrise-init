package android

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/bitrise-init/analytics"
	"github.com/bitrise-io/bitrise-init/detectors/gradle"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	bitriseModels "github.com/bitrise-io/bitrise/v2/models"
	envmanModels "github.com/bitrise-io/envman/v2/models"
	"github.com/bitrise-io/go-utils/log"
)

const (
	ScannerName                   = "android"
	ConfigName                    = "android-config"
	ConfigNameKotlinScript        = "android-config-kts"
	DefaultConfigName             = "default-android-config"
	DefaultConfigNameKotlinScript = "default-android-config-kts"

	testsWorkflowID         = "run_tests"
	testsWorkflowSummary    = "Run your Android unit tests and get the test report."
	testWorkflowDescription = "The workflow will first clone your Git repository, cache your Gradle dependencies, install Android tools, run your Android unit tests and save the test report."

	testPipelineID = "run_tests"

	runInstrumentedTestsWorkflowID          = "run_instrumented_tests"
	runInstrumentedTestsWorkflowSummary     = "Run your Android instrumented tests and get the test report."
	runInstrumentedTestsWorkflowDescription = "The workflow will first clone your Git repository, cache your Gradle dependencies, install Android tools, run your Android instrumented tests and save the test report."
	TestShardCountEnvKey                    = "TEST_SHARD_COUNT"
	TestShardCountEnvValue                  = 2
	ParallelTotalEnvKey                     = "BITRISE_IO_PARALLEL_TOTAL"
	ParallelIndexEnvKey                     = "BITRISE_IO_PARALLEL_INDEX"

	buildWorkflowID          = "build_apk"
	buildWorkflowSummary     = "Run your Android unit tests and create an APK file to install your app on a device or share it with your team."
	buildWorkflowDescription = "The workflow will first clone your Git repository, install Android tools, set the project's version code based on the build number, run Android lint and unit tests, build the project's APK file and save it."

	ProjectLocationInputKey     = "project_location"
	ProjectLocationInputEnvKey  = "PROJECT_LOCATION"
	ProjectLocationInputTitle   = "The root directory of an Android project"
	ProjectLocationInputSummary = "The root directory of your Android project, stored as an Environment Variable. In your Workflows, you can specify paths relative to this path. You can change this at any time."

	ModuleBuildGradlePathInputKey = "build_gradle_path"

	VariantInputKey     = "variant"
	VariantInputEnvKey  = "VARIANT"
	VariantInputTitle   = "Variant"
	VariantInputSummary = "Your Android build variant. You can add variants at any time, as well as further configure your existing variants later."

	ModuleInputKey     = "module"
	ModuleInputEnvKey  = "MODULE"
	ModuleInputTitle   = "Module"
	ModuleInputSummary = "Modules provide a container for your Android project's source code, resource files, and app level settings, such as the module-level build file and Android manifest file. Each module can be independently built, tested, and debugged. You can add new modules to your Bitrise builds at any time."

	BuildScriptInputTitle   = "Does your app use Kotlin build scripts?"
	BuildScriptInputSummary = "The workflow configuration slightly differs based on what language (Groovy or Kotlin) you used in your build scripts."

	GradlewPathInputKey       = "gradlew_path"
	GradlewGradleTaskInputKey = "gradle_task"

	CacheLevelInputKey = "cache_level"
	CacheLevelNone     = "none"

	gradleKotlinBuildFile    = "build.gradle.kts"
	gradleKotlinSettingsFile = "settings.gradle.kts"
)

// Scanner ...
type Scanner struct {
	GradleProject gradle.Project
	Icons         models.Icons
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (scanner *Scanner) Name() string {
	return ScannerName
}

// ExcludedScannerNames ...
func (scanner *Scanner) ExcludedScannerNames() []string {
	return nil
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (_ bool, err error) {
	log.TInfof("Searching for Gradle project files...")

	gradleProject, err := gradle.ScanProject(searchDir)
	if err != nil {
		return false, err
	}

	log.TDonef("Gradle project found: %v", gradleProject != nil)
	if gradleProject == nil {
		return false, nil
	}

	printGradleProject(*gradleProject)

	log.TInfof("Searching for Android dependencies...")
	androidDetected, err := gradleProject.DetectAnyDependencies([]string{
		"com.android.application",
	})
	if err != nil {
		return false, err
	}

	log.TDonef("Android dependencies found: %v", androidDetected)
	scanner.GradleProject = *gradleProject

	log.TInfof("Searching for project icons...")
	scanner.Icons, err = LookupIcons(scanner.GradleProject.RootDirEntry.AbsPath, searchDir)
	if err != nil {
		log.TWarnf("Failed to find icons: %v", err)
		analytics.LogInfo("android-icon-lookup", analytics.DetectorErrorData("android", err), "Failed to lookup android icon")
	}
	log.TDonef("%d icon(s) found", len(scanner.Icons))

	return androidDetected, nil
}

/*
generated config inputs:
- project root dir (gradlew dir) -> gradlew path
- app module's gradle build script -> app module's module and variant

- install-missing-android-tools@%s:
  inputs:
  - gradlew_path: $PROJECT_LOCATION/gradlew
- change-android-versioncode-and-versionname@%s:
  inputs:
  - build_gradle_path: $PROJECT_LOCATION/$MODULE/build.gradle.kts
- android-lint@%s:
  inputs:
  - project_location: $PROJECT_LOCATION
  - variant: $VARIANT
  - cache_level: none
- android-unit-test@%s:
  inputs:
  - project_location: $PROJECT_LOCATION
  - variant: $VARIANT
  - cache_level: none
- android-build@%s:
  inputs:
  - project_location: $PROJECT_LOCATION
  - module: $MODULE
  - variant: $VARIANT
  - cache_level: none
*/

// Options ...
// TODO: restore icon search support
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputSummary, ProjectLocationInputEnvKey, models.TypeSelector)
	moduleOption := models.NewOption(ModuleInputTitle, ModuleInputSummary, ModuleInputEnvKey, models.TypeUserInput)
	variantOption := models.NewOption(VariantInputTitle, VariantInputSummary, VariantInputEnvKey, models.TypeOptionalUserInput)

	possibleAppModuleBuildScriptPaths := scanner.listPossibleAppModuleBuildScriptPaths()
	if len(possibleAppModuleBuildScriptPaths) == 0 {
		// TODO: validate it at project type detection phase
		return models.OptionNode{}, nil, nil, fmt.Errorf("no Gradle build scripts found")
	}

	modulePathsToIsKotlinDSL := map[string]bool{}
	for _, buildScriptPath := range possibleAppModuleBuildScriptPaths {
		modulePath := modulePathFromBuildScriptPath(scanner.GradleProject.RootDirEntry.RelPath, buildScriptPath)
		if modulePath != "" {
			isKotlinDSL := strings.HasSuffix(scanner.GradleProject.RootDirEntry.RelPath, ".kts")
			modulePathsToIsKotlinDSL[modulePath] = isKotlinDSL
		} else {
			// TODO: remote log if no module name found
		}
	}

	iconIDs := make([]string, len(scanner.Icons))
	for i, icon := range scanner.Icons {
		iconIDs[i] = icon.Filename
	}

	for moduleName, isKotlinDSL := range modulePathsToIsKotlinDSL {
		var configOption *models.OptionNode
		if isKotlinDSL {
			configOption = models.NewConfigOption(ConfigNameKotlinScript, iconIDs)
		} else {
			configOption = models.NewConfigOption(ConfigName, iconIDs)
		}

		projectLocationOption.AddOption(scanner.GradleProject.RootDirEntry.RelPath, moduleOption)
		moduleOption.AddOption(moduleName, variantOption)
		variantOption.AddConfig("", configOption)
	}

	return *projectLocationOption, nil, scanner.Icons, nil

	//projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputSummary, ProjectLocationInputEnvKey, models.TypeSelector)
	//warnings := models.Warnings{}
	//appIconsAllProjects := models.Icons{}
	//
	//for _, project := range scanner.Projects {
	//	warnings = append(warnings, project.Warnings...)
	//	appIconsAllProjects = append(appIconsAllProjects, project.Icons...)
	//
	//	iconIDs := make([]string, len(project.Icons))
	//	for i, icon := range project.Icons {
	//		iconIDs[i] = icon.Filename
	//	}
	//
	//	name := ConfigName
	//	if project.UsesKotlinBuildScript {
	//		name = ConfigNameKotlinScript
	//	}
	//	configOption := models.NewConfigOption(name, iconIDs)
	//	moduleOption := models.NewOption(ModuleInputTitle, ModuleInputSummary, ModuleInputEnvKey, models.TypeUserInput)
	//	variantOption := models.NewOption(VariantInputTitle, VariantInputSummary, VariantInputEnvKey, models.TypeOptionalUserInput)
	//
	//	projectLocationOption.AddOption(project.RelPath, moduleOption)
	//	moduleOption.AddOption("app", variantOption)
	//	variantOption.AddConfig("", configOption)
	//}
	//
	//return *projectLocationOption, warnings, appIconsAllProjects, nil
}

// DefaultOptions ...
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputSummary, ProjectLocationInputEnvKey, models.TypeUserInput)
	moduleOption := models.NewOption(ModuleInputTitle, ModuleInputSummary, ModuleInputEnvKey, models.TypeUserInput)
	variantOption := models.NewOption(VariantInputTitle, VariantInputSummary, VariantInputEnvKey, models.TypeOptionalUserInput)

	buildScriptOption := models.NewOption(BuildScriptInputTitle, BuildScriptInputSummary, "", models.TypeSelector)
	regularConfigOption := models.NewConfigOption(DefaultConfigName, nil)
	kotlinScriptConfigOption := models.NewConfigOption(DefaultConfigNameKotlinScript, nil)

	projectLocationOption.AddOption(models.UserInputOptionDefaultValue, moduleOption)
	moduleOption.AddOption(models.UserInputOptionDefaultValue, variantOption)
	variantOption.AddOption(models.UserInputOptionDefaultValue, buildScriptOption)

	buildScriptOption.AddConfig("yes", kotlinScriptConfigOption)
	buildScriptOption.AddOption("no", regularConfigOption)

	return *projectLocationOption
}

type configBuildingParams struct {
	name            string
	useKotlinScript bool
}

// Configs ...
func (scanner *Scanner) Configs(sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	var params []configBuildingParams
	params = append(params, configBuildingParams{
		name:            ConfigName,
		useKotlinScript: false,
	})
	params = append(params, configBuildingParams{
		name:            ConfigNameKotlinScript,
		useKotlinScript: true,
	})
	return scanner.generateConfigs(sshKeyActivation, params)
}

// DefaultConfigs ...
func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	params := []configBuildingParams{
		{name: DefaultConfigName, useKotlinScript: false},
		{name: DefaultConfigNameKotlinScript, useKotlinScript: true},
	}
	return scanner.generateConfigs(models.SSHKeyActivationConditional, params)
}

func (scanner *Scanner) generateConfigs(sshKeyActivation models.SSHKeyActivation, params []configBuildingParams) (models.BitriseConfigMap, error) {
	bitriseDataMap := models.BitriseConfigMap{}

	for _, param := range params {
		configBuilder := scanner.generateConfigBuilder(sshKeyActivation, param.useKotlinScript)

		config, err := configBuilder.Generate(ScannerName,
			envmanModels.EnvironmentItemModel{TestShardCountEnvKey: TestShardCountEnvValue},
		)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		data, err := yaml.Marshal(config)
		if err != nil {
			return models.BitriseConfigMap{}, err
		}

		bitriseDataMap[param.name] = string(data)
	}

	return bitriseDataMap, nil
}

func (scanner *Scanner) generateConfigBuilder(sshKeyActivation models.SSHKeyActivation, useKotlinBuildScript bool) models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	projectLocationEnv, gradlewPath, moduleEnv, variantEnv := "$"+ProjectLocationInputEnvKey, "$"+ProjectLocationInputEnvKey+"/gradlew", "$"+ModuleInputEnvKey, "$"+VariantInputEnvKey

	//-- test
	configBuilder.AppendStepListItemsTo(testsWorkflowID, steps.DefaultPrepareStepList(steps.PrepareListParams{
		SSHKeyActivation: sshKeyActivation})...)
	configBuilder.AppendStepListItemsTo(testsWorkflowID, steps.RestoreGradleCache())
	configBuilder.AppendStepListItemsTo(testsWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{GradlewPathInputKey: gradlewPath},
	))
	configBuilder.AppendStepListItemsTo(testsWorkflowID, steps.AndroidUnitTestStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
		envmanModels.EnvironmentItemModel{
			CacheLevelInputKey: CacheLevelNone,
		},
	))
	configBuilder.AppendStepListItemsTo(testsWorkflowID, steps.SaveGradleCache())
	configBuilder.AppendStepListItemsTo(testsWorkflowID, steps.DefaultDeployStepList()...)
	configBuilder.SetWorkflowSummaryTo(testsWorkflowID, testsWorkflowSummary)
	configBuilder.SetWorkflowDescriptionTo(testsWorkflowID, testWorkflowDescription)

	//-- instrumented test
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.DefaultPrepareStepList(steps.PrepareListParams{
		SSHKeyActivation: sshKeyActivation,
	})...)
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.RestoreGradleCache())
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{GradlewPathInputKey: gradlewPath},
	))
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.AvdManagerStepListItem())
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.WaitForAndroidEmulatorStepListItem())
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.GradleRunnerStepListItem(
		gradlewPath,
		fmt.Sprintf("connectedAndroidTest \\\n  -Pandroid.testInstrumentationRunnerArguments.numShards=$%s \\\n  -Pandroid.testInstrumentationRunnerArguments.shardIndex=$%s",
			ParallelTotalEnvKey,
			ParallelIndexEnvKey,
		),
	))
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.SaveGradleCache())
	configBuilder.AppendStepListItemsTo(runInstrumentedTestsWorkflowID, steps.DefaultDeployStepList()...)
	configBuilder.SetWorkflowSummaryTo(runInstrumentedTestsWorkflowID, runInstrumentedTestsWorkflowSummary)
	configBuilder.SetWorkflowDescriptionTo(runInstrumentedTestsWorkflowID, runInstrumentedTestsWorkflowDescription)

	configBuilder.SetGraphPipelineWorkflowTo(testPipelineID, runInstrumentedTestsWorkflowID, bitriseModels.GraphPipelineWorkflowModel{
		Parallel: "$" + TestShardCountEnvKey,
	})

	//-- build
	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.DefaultPrepareStepList(steps.PrepareListParams{
		SSHKeyActivation: sshKeyActivation,
	})...)
	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.InstallMissingAndroidToolsStepListItem(
		envmanModels.EnvironmentItemModel{GradlewPathInputKey: gradlewPath},
	))

	basePath := filepath.Join(projectLocationEnv, moduleEnv)
	path := filepath.Join(basePath, "build.gradle")
	if useKotlinBuildScript {
		path = filepath.Join(basePath, gradleKotlinBuildFile)
	}
	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.ChangeAndroidVersionCodeAndVersionNameStepListItem(
		envmanModels.EnvironmentItemModel{ModuleBuildGradlePathInputKey: path},
	))

	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.AndroidLintStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
		envmanModels.EnvironmentItemModel{
			CacheLevelInputKey: CacheLevelNone,
		},
	))
	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.AndroidUnitTestStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
		envmanModels.EnvironmentItemModel{
			CacheLevelInputKey: CacheLevelNone,
		},
	))

	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{
			ProjectLocationInputKey: projectLocationEnv,
		},
		envmanModels.EnvironmentItemModel{
			ModuleInputKey: moduleEnv,
		},
		envmanModels.EnvironmentItemModel{
			VariantInputKey: variantEnv,
		},
		envmanModels.EnvironmentItemModel{
			CacheLevelInputKey: CacheLevelNone,
		},
	))
	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.SignAPKStepListItem())
	configBuilder.AppendStepListItemsTo(buildWorkflowID, steps.DefaultDeployStepList()...)

	configBuilder.SetWorkflowDescriptionTo(buildWorkflowID, buildWorkflowDescription)
	configBuilder.SetWorkflowSummaryTo(buildWorkflowID, buildWorkflowSummary)

	return *configBuilder
}

func printGradleProject(gradleProject gradle.Project) {
	log.TPrintf("Project root dir: %s", gradleProject.RootDirEntry.RelPath)
	log.TPrintf("Gradle wrapper script: %s", gradleProject.GradlewFileEntry.RelPath)
	if gradleProject.ConfigDirEntry != nil {
		log.TPrintf("Gradle config dir: %s", gradleProject.ConfigDirEntry.RelPath)
	}
	if gradleProject.VersionCatalogFileEntry != nil {
		log.TPrintf("Version catalog file: %s", gradleProject.VersionCatalogFileEntry.RelPath)
	}
	if gradleProject.SettingsGradleFileEntry != nil {
		log.TPrintf("Gradle settings file: %s", gradleProject.SettingsGradleFileEntry.RelPath)
	}
	if len(gradleProject.IncludedProjects) > 0 {
		log.TPrintf("Included projects:")
		for _, includedProject := range gradleProject.IncludedProjects {
			log.TPrintf("- %s: %s", includedProject.Name, includedProject.BuildScriptFileEntry.RelPath)
		}
	}
}

func (scanner *Scanner) listPossibleAppModuleBuildScriptPaths() []string {
	var appModuleBuildScriptPath string
	var possibleAppModuleBuildScriptPaths []string

	for _, includedProject := range scanner.GradleProject.IncludedProjects {
		if includedProject.Name == "app" {
			appModuleBuildScriptPath = includedProject.BuildScriptFileEntry.RelPath
			continue
		}

		possibleAppModuleBuildScriptPaths = append(possibleAppModuleBuildScriptPaths, includedProject.BuildScriptFileEntry.RelPath)
	}

	if len(possibleAppModuleBuildScriptPaths) == 0 {
		// TODO: remote log if no included projects found
		for _, buildScript := range scanner.GradleProject.AllBuildScriptFileEntries {
			possibleAppModuleBuildScriptPaths = append(possibleAppModuleBuildScriptPaths, buildScript.RelPath)
		}
	}

	if appModuleBuildScriptPath != "" {
		possibleAppModuleBuildScriptPaths = append([]string{appModuleBuildScriptPath}, possibleAppModuleBuildScriptPaths...)
	} else {
		// TODO: remote log if no app module build script found
	}

	return possibleAppModuleBuildScriptPaths
}

// :backend:datastore: ./backend/datastore/build.gradle.kts
// modulePathFromBuildScriptPath returns the module path from the build script path
func modulePathFromBuildScriptPath(projectRootDir, buildScriptPth string) string {
	relBuildScriptPath := strings.TrimPrefix(buildScriptPth, projectRootDir)
	relBuildScriptPath = strings.TrimPrefix(relBuildScriptPath, "/")
	pathComponents := strings.Split(relBuildScriptPath, "/")
	if len(pathComponents) < 2 {
		return ""
	}

	return strings.Join(pathComponents[:len(pathComponents)-1], "/")
}
