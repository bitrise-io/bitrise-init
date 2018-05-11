package android

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-tools/go-android/gradle"
	"github.com/bitrise-tools/go-android/sdk"
	"github.com/bitrise-tools/go-android/sdkcomponent"
	"github.com/bitrise-tools/go-android/sdkmanager"
)

// Constants ...
const (
	ScannerName       = "android"
	ConfigName        = "android-config"
	DefaultConfigName = "default-android-config"

	ProjectLocationInputKey    = "project_location"
	ProjectLocationInputEnvKey = "PROJECT_LOCATION"
	ProjectLocationInputTitle  = "The root directory of an Android project"

	ModuleBuildGradlePathInputKey = "build_gradle_path"

	ModuleInputKey    = "module"
	ModuleInputEnvKey = "MODULE"
	ModuleInputTitle  = "Module"

	VariantInputKey         = "variant"
	TestVariantInputEnvKey  = "TEST_VARIANT"
	BuildVariantInputEnvKey = "BUILD_VARIANT"
	TestVariantInputTitle   = "Variant for testing"
	BuildVariantInputTitle  = "Variant for building"

	GradlewPathInputKey    = "gradlew_path"
	GradlewPathInputEnvKey = "GRADLEW_PATH"
	GradlewPathInputTitle  = "Gradlew file path"
)

func walk(src string, fn func(path string, info os.FileInfo) error) error {
	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			if err := walk(filepath.Join(src, info.Name()), fn); err != nil {
				return err
			}
		}
		if err := fn(filepath.Join(src, info.Name()), info); err != nil {
			return err
		}
	}
	return nil
}

func checkFiles(path string, files ...string) (bool, error) {
	for _, file := range files {
		exists, err := pathutil.IsPathExists(filepath.Join(path, file))
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}
	return true, nil
}

func walkMultipleFiles(searchDir string, files ...string) (matches []string, err error) {
	match, err := checkFiles(searchDir, files...)
	if err != nil {
		return nil, err
	}
	if match {
		matches = append(matches, searchDir)
	}
	return matches, walk(searchDir, func(path string, info os.FileInfo) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			match, err := checkFiles(path, files...)
			if err != nil {
				return err
			}
			if match {
				matches = append(matches, path)
			}
		}
		return nil
	})
}

func checkLocalProperties(projectDir string) string {
	localPropertiesPth := filepath.Join(projectDir, "local.properties")
	exist, err := pathutil.IsPathExists(localPropertiesPth)
	if err == nil && exist {
		return fmt.Sprintf(`The local.properties file must NOT be checked into Version Control Systems, as it contains information specific to your local configuration.
The location of the file is: %s`, localPropertiesPth)
	}
	return ""
}

func checkGradlew(projectDir string) error {
	gradlewPth := filepath.Join(projectDir, "gradlew")
	exist, err := pathutil.IsPathExists(gradlewPth)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New(`<b>No Gradle Wrapper (gradlew) found.</b> 
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`)
	}
	return nil
}

func ensureAndroidLicences(androidHome string, isLegacySDK bool) error {
	if !isLegacySDK {
		licensesCmd := command.New(filepath.Join(androidHome, "tools/bin/sdkmanager"), "--licenses")
		licensesCmd.SetStdin(bytes.NewReader([]byte(strings.Repeat("y\n", 1000))))
		if err := licensesCmd.Run(); err == nil {
			return nil
		}
	}

	licenceMap := map[string]string{
		"android-sdk-license":           "8933bad161af4178b1185d1a37fbf41ea5269c55\n\nd56f5187479451eabf01fb78af6dfcb131a6481e",
		"android-googletv-license":      "\n601085b94cd77f0b54ff86406957099ebe79c4d6",
		"android-sdk-preview-license":   "\n84831b9409646a918e30573bab4c9c91346d8abd",
		"intel-android-extra-license":   "\nd975f751698a77b662f1254ddbeed3901e976f5a",
		"google-gdk-license":            "\n33b6a2b64607f11b759f320ef9dff4ae5c47d97a",
		"mips-android-sysimage-license": "\ne9acab5b5fbb560a72cfaecce8946896ff6aab9d",
	}

	licencesDir := filepath.Join(androidHome, "licenses")
	if exist, err := pathutil.IsDirExists(licencesDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(licencesDir, os.ModePerm); err != nil {
			return err
		}
	}

	for name, content := range licenceMap {
		pth := filepath.Join(licencesDir, name)

		if err := fileutil.WriteStringToFile(pth, content); err != nil {
			return err
		}
	}

	return nil
}

func installMissingAndroidTools(srcRoot, androidHome string) error {
	gradlewPath := filepath.Join(srcRoot, "gradlew")
	if err := os.Chmod(gradlewPath, 0770); err != nil {
		return fmt.Errorf("failed to set executable permission for gradlew, error: %s", err)
	}

	androidSdk, err := sdk.New(androidHome)
	if err != nil {
		return fmt.Errorf("failed to initialize Android SDK, error: %s", err)
	}

	sdkManager, err := sdkmanager.New(androidSdk)
	if err != nil {
		return fmt.Errorf("failed to create SDK manager, error: %s", err)
	}

	if err := ensureAndroidLicences(androidHome, sdkManager.IsLegacySDK()); err != nil {
		return fmt.Errorf("failed to ensure android licences, error: %s", err)
	}

	retryCount := 0
	for true {
		gradleCmd := command.New("./gradlew", "dependencies")
		gradleCmd.SetStdin(strings.NewReader("y"))
		gradleCmd.SetDir(filepath.Dir(gradlewPath))

		if out, err := gradleCmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			reader := strings.NewReader(out)
			scanner := bufio.NewScanner(reader)

			missingSDKComponentFound := false

			for scanner.Scan() {
				line := scanner.Text()
				{
					// failed to find target with hash string 'android-22'
					targetPattern := `failed to find target with hash string 'android-(?P<version>.*)'\s*`
					targetRe := regexp.MustCompile(targetPattern)
					if matches := targetRe.FindStringSubmatch(line); len(matches) == 2 {
						missingSDKComponentFound = true

						targetVersion := "android-" + matches[1]

						platformComponent := sdkcomponent.Platform{
							Version: targetVersion,
						}

						cmd := sdkManager.InstallCommand(platformComponent)
						cmd.SetStdin(strings.NewReader("y"))

						if err := retry.Times(1).Wait(time.Second).Try(func(attempt uint) error {
							if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
								if attempt > 0 {
									return fmt.Errorf("output: %s, error: %s", out, err)
								}
								return err
							}

							return nil
						}); err != nil {
							return fmt.Errorf("failed to install platform:\n%s", err)
						}
					}
				}

				{
					// failed to find Build Tools revision 22.0.1
					buildToolsPattern := `failed to find Build Tools revision (?P<version>[0-9.]*)\s*`
					buildToolsRe := regexp.MustCompile(buildToolsPattern)
					if matches := buildToolsRe.FindStringSubmatch(line); len(matches) == 2 {
						missingSDKComponentFound = true

						buildToolsVersion := matches[1]

						buildToolsComponent := sdkcomponent.BuildTool{
							Version: buildToolsVersion,
						}

						cmd := sdkManager.InstallCommand(buildToolsComponent)
						cmd.SetStdin(strings.NewReader("y"))

						if err := retry.Times(1).Wait(time.Second).Try(func(attempt uint) error {
							if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
								if attempt > 0 {
									return fmt.Errorf("output: %s, error: %s", out, err)
								}
								return err
							}

							return nil
						}); err != nil {
							return fmt.Errorf("failed to install build tools:\n%s", err)
						}
					}
				}

				{
					// Example: "Could not find com.android.support.constraint:constraint-layout:1.0.2."
					extrasPattern := `Could not find (?P<package>com\.android\.support\..*)\.`
					extrasRe := regexp.MustCompile(extrasPattern)
					if matches := extrasRe.FindStringSubmatch(line); len(matches) == 2 {
						missingSDKComponentFound = true

						lib := matches[1]
						firstColon := strings.Index(lib, ":")
						lib = strings.Replace(lib[:firstColon], ".", ";", -1) + strings.Replace(lib[firstColon:], ":", ";", -1)

						extrasComponents := sdkcomponent.SupportLibraryInstallComponents()
						extrasComponents = append(extrasComponents, sdkcomponent.Extras{
							Provider:    "m2repository",
							PackageName: lib,
						})
						for _, e := range extrasComponents {
							cmd := sdkManager.InstallCommand(e)
							cmd.SetStdin(strings.NewReader("y"))

							if err := retry.Times(1).Wait(time.Second).Try(func(attempt uint) error {
								if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
									if attempt > 0 {
										return fmt.Errorf("output: %s, error: %s", out, err)
									}
									return err
								}

								return nil
							}); err != nil {
								return fmt.Errorf("Failed to install support library dependency:\n%s", err)
							}
						}
					}
				}
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to analyze gradle output, error: %s", err)
			}

			if !missingSDKComponentFound {
				if retryCount < 2 {
					retryCount++
					continue
				}
				fmt.Println(out)
				return fmt.Errorf("%s", err)
			}
		} else {
			break
		}
	}
	return nil
}

func (scanner *Scanner) generateOptions(searchDir string) (models.OptionModel, models.Warnings, error) {
	warnings := models.Warnings{}

	projectLocationOption := models.NewOption(ProjectLocationInputTitle, ProjectLocationInputEnvKey)

	for _, projectRoot := range scanner.ProjectRoots {
		if warning := checkLocalProperties(projectRoot); warning != "" {
			warnings = append(warnings, warning)
		}

		if err := checkGradlew(projectRoot); err != nil {
			return models.OptionModel{}, warnings, err
		}

		if err := installMissingAndroidTools(projectRoot, os.Getenv("ANDROID_HOME")); err != nil {
			return models.OptionModel{}, warnings, err
		}

		proj, err := gradle.NewProject(projectRoot)
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
		testVariantsMap, err := proj.GetTask("test").GetVariants()
		if err != nil {
			return models.OptionModel{}, warnings, err
		}
		buildVariantsMap, err := proj.GetTask("assemble").GetVariants()
		if err != nil {
			return models.OptionModel{}, warnings, err
		}

		moduleOption := models.NewOption(ModuleInputTitle, ModuleInputEnvKey)

		for module, variants := range buildVariantsMap {
			testVariantOption := models.NewOption(TestVariantInputTitle, TestVariantInputEnvKey)
			buildVariantOption := models.NewOption(BuildVariantInputTitle, BuildVariantInputEnvKey)

			configOption := models.NewConfigOption(ConfigName)

			if !scanner.ExcludeTest {
				for _, variant := range testVariantsMap[module] {
					variant = strings.TrimSuffix(variant, "UnitTest")
					testVariantOption.AddOption(variant, configOption)
				}
			}

			for _, variant := range variants {
				if !scanner.ExcludeTest {
					configOption = testVariantOption
				}
				buildVariantOption.AddOption(variant, configOption)
			}

			moduleOption.AddOption(module, buildVariantOption)
		}

		relProjectRoot, err := filepath.Rel(scanner.SearchDir, projectRoot)
		if err != nil {
			return models.OptionModel{}, warnings, err
		}

		gradlewPthOption := models.NewOption(GradlewPathInputTitle, GradlewPathInputEnvKey)
		gradlewPthOption.AddOption(filepath.Join(relProjectRoot, "gradlew"), moduleOption)

		projectLocationOption.AddOption(relProjectRoot, gradlewPthOption)
	}
	return *projectLocationOption, warnings, nil
}

func (scanner *Scanner) generateConfigBuilder(isIncludeCache bool) models.ConfigBuilderModel {
	configBuilder := models.NewDefaultConfigBuilder()

	projectLocationEnv, moduleEnv, testVariantEnv, buildVariantEnv := "$"+ProjectLocationInputEnvKey, "$"+ModuleInputEnvKey, "$"+TestVariantInputEnvKey, "$"+BuildVariantInputEnvKey

	//-- primary
	if !scanner.ExcludeTest {
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.InstallMissingAndroidToolsStepListItem())
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidLintStepListItem(
			envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
			envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
			envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
		))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.AndroidUnitTestStepListItem(
			envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
			envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
			envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
		))
		configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)
	}
	//-- deploy
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultPrepareStepList(isIncludeCache)...)
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.InstallMissingAndroidToolsStepListItem())

	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.ChangeAndroidVersionCodeAndVersionNameStepListItem(
		envmanModels.EnvironmentItemModel{ModuleBuildGradlePathInputKey: filepath.Join(projectLocationEnv, moduleEnv, "build.gradle")},
	))
	if !scanner.ExcludeTest {
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidLintStepListItem(
			envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
			envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
			envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
		))
		configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidUnitTestStepListItem(
			envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
			envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
			envmanModels.EnvironmentItemModel{VariantInputKey: testVariantEnv},
		))
	}
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.AndroidBuildStepListItem(
		envmanModels.EnvironmentItemModel{ProjectLocationInputKey: projectLocationEnv},
		envmanModels.EnvironmentItemModel{ModuleInputKey: moduleEnv},
		envmanModels.EnvironmentItemModel{VariantInputKey: buildVariantEnv},
	))
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.SignAPKStepListItem())
	configBuilder.AppendStepListItemsTo(models.DeployWorkflowID, steps.DefaultDeployStepList(isIncludeCache)...)

	configBuilder.SetWorkflowDescriptionTo(models.DeployWorkflowID, deployWorkflowDescription)

	return *configBuilder
}
