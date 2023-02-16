package integration

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/output"
	"github.com/bitrise-io/bitrise-init/scanner"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/stretchr/testify/require"
)

const noTestPackageJSON = `{
  "name": "SampleAppsReactNativeAndroid",
  "version": "0.0.1",
  "private": true,
  "scripts": {
    "start": "node node_modules/react-native/local-cli/cli.js start"
  },
  "dependencies": {
    "react": "15.4.2",
    "react-native": "0.42.0"
  },
  "devDependencies": {
    "babel-jest": "19.0.0",
    "babel-preset-react-native": "1.9.1",
    "jest": "19.0.2",
    "react-test-renderer": "15.4.2"
  },
  "jest": {
    "preset": "react-native"
  }
}`

const simpleSample = "https://github.com/bitrise-samples/sample-apps-react-native-ios-and-android.git"

func TestReactNative(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"joplin",
			"https://github.com/bitrise-io/joplin.git",
			"",
			sampleAppsReactNativeJoplinResultYML,
			sampleAppsReactNativeJoplinVersions,
		},
		{
			"sample-apps-react-native-ios-and-android",
			simpleSample,
			"",
			sampleAppsReactNativeIosAndAndroidResultYML,
			sampleAppsReactNativeIosAndAndroidVersions,
		},
		{
			"sample-apps-react-native-subdir",
			"https://github.com/bitrise-samples/sample-apps-react-native-subdir.git",
			"",
			sampleAppsReactNativeSubdirResultYML,
			sampleAppsReactNativeSubdirVersions,
		},
	}

	helper.Execute(t, testCases)
}

func TestNoTests(t *testing.T) {
	testName := "sample-apps-react-native-ios-and-android-no-test"
	dir := setupSample(t, testName, simpleSample)

	err := ioutil.WriteFile(filepath.Join(dir, "package.json"), []byte(noTestPackageJSON), 0600)
	require.NoError(t, err)

	generateAndValidateResult(t, testName, dir, sampleAppsReactNativeIosAndAndroidNoTestResultYML, sampleAppsReactNativeIosAndAndroidNoTestVersions)
}

func TestYarn(t *testing.T) {
	testName := "sample-apps-react-native-ios-and-android-yarn"
	dir := setupSample(t, testName, simpleSample)

	yarnCommand := command.New("yarn", "install")
	yarnCommand.SetDir(dir)
	out, err := yarnCommand.RunAndReturnTrimmedCombinedOutput()
	require.NoError(t, err, out)

	generateAndValidateResult(t, testName, dir, sampleAppsReactNativeIosAndAndroidYarnResultYML, sampleAppsReactNativeIosAndAndroidYarnVersions)
}

// Helpers

func setupSample(t *testing.T, name, repoURL string) string {
	tmpDir := t.TempDir()
	sampleAppDir := filepath.Join(tmpDir, name)
	helper.GitClone(t, sampleAppDir, repoURL)

	return sampleAppDir
}

func generateAndValidateResult(t *testing.T, name, dir, expectedResult string, expectedVersions []interface{}) {
	_, err := scanner.GenerateAndWriteResults(dir, dir, output.YAMLFormat)
	require.NoError(t, err)

	scanResultPth := filepath.Join(dir, "result.yml")

	result, err := fileutil.ReadStringFromFile(scanResultPth)
	require.NoError(t, err)

	helper.ValidateConfigExpectation(t, name, strings.TrimSpace(expectedResult), strings.TrimSpace(result), expectedVersions...)
}

// Expected results

var sampleAppsReactNativeSubdirVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsReactNativeSubdirResultYML = fmt.Sprintf(`options:
  react-native:
    title: React Native project directory
    summary: Path of the directory containing the project's `+"`package.json`"+` file.
    env_key: WORKDIR
    type: selector
    value_map:
      project:
        title: The root directory of an Android project
        summary: The root directory of your Android project, stored as an Environment
          Variable. In your Workflows, you can specify paths relative to this path.
          You can change this at any time.
        env_key: PROJECT_LOCATION
        type: selector
        value_map:
          project/android:
            title: Module
            summary: Modules provide a container for your Android project's source
              code, resource files, and app level settings, such as the module-level
              build file and Android manifest file. Each module can be independently
              built, tested, and debugged. You can add new modules to your Bitrise
              builds at any time.
            env_key: MODULE
            type: user_input
            value_map:
              app:
                title: Variant
                summary: Your Android build variant. You can add variants at any time,
                  as well as further configure your existing variants later.
                env_key: VARIANT
                type: user_input_optional
                value_map:
                  Debug:
                    title: Project or Workspace path
                    summary: The location of your Xcode project or Xcode workspace
                      files, stored as an Environment Variable. In your Workflows,
                      you can specify paths relative to this path.
                    env_key: BITRISE_PROJECT_PATH
                    type: selector
                    value_map:
                      project/ios/SampleAppsReactNativeAndroid.xcodeproj:
                        title: Scheme name
                        summary: An Xcode scheme defines a collection of targets to
                          build, a configuration to use when building, and a collection
                          of tests to execute. Only shared schemes are detected automatically
                          but you can use any scheme as a target on Bitrise. You can
                          change the scheme at any time in your Env Vars.
                        env_key: BITRISE_SCHEME
                        type: selector
                        value_map:
                          SampleAppsReactNativeAndroid:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-test-config
                              app-store:
                                config: react-native-android-ios-test-config
                              development:
                                config: react-native-android-ios-test-config
                              enterprise:
                                config: react-native-android-ios-test-config
                          SampleAppsReactNativeAndroid-tvOS:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-test-config
                              app-store:
                                config: react-native-android-ios-test-config
                              development:
                                config: react-native-android-ios-test-config
                              enterprise:
                                config: react-native-android-ios-test-config
configs:
  react-native:
    react-native-android-ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests, builds and deploys the app using *Deploy to bitrise.io* Step.

            Next steps:
            - Set up an [Apple service with API key](https://devcenter.bitrise.io/en/accounts/connecting-to-services/connecting-to-an-apple-service-with-api-key.html).
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: test
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
              - automatic_code_signing: api-key
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Runs tests.

            Next steps:
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: test
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, sampleAppsReactNativeSubdirVersions...)

var sampleAppsReactNativeIosAndAndroidVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsReactNativeIosAndAndroidResultYML = fmt.Sprintf(`options:
  react-native:
    title: React Native project directory
    summary: Path of the directory containing the project's `+"`package.json`"+` file.
    env_key: WORKDIR
    type: selector
    value_map:
      .:
        title: The root directory of an Android project
        summary: The root directory of your Android project, stored as an Environment
          Variable. In your Workflows, you can specify paths relative to this path.
          You can change this at any time.
        env_key: PROJECT_LOCATION
        type: selector
        value_map:
          android:
            title: Module
            summary: Modules provide a container for your Android project's source
              code, resource files, and app level settings, such as the module-level
              build file and Android manifest file. Each module can be independently
              built, tested, and debugged. You can add new modules to your Bitrise
              builds at any time.
            env_key: MODULE
            type: user_input
            value_map:
              app:
                title: Variant
                summary: Your Android build variant. You can add variants at any time,
                  as well as further configure your existing variants later.
                env_key: VARIANT
                type: user_input_optional
                value_map:
                  Debug:
                    title: Project or Workspace path
                    summary: The location of your Xcode project or Xcode workspace
                      files, stored as an Environment Variable. In your Workflows,
                      you can specify paths relative to this path.
                    env_key: BITRISE_PROJECT_PATH
                    type: selector
                    value_map:
                      ios/SampleAppsReactNativeAndroid.xcodeproj:
                        title: Scheme name
                        summary: An Xcode scheme defines a collection of targets to
                          build, a configuration to use when building, and a collection
                          of tests to execute. Only shared schemes are detected automatically
                          but you can use any scheme as a target on Bitrise. You can
                          change the scheme at any time in your Env Vars.
                        env_key: BITRISE_SCHEME
                        type: selector
                        value_map:
                          SampleAppsReactNativeAndroid:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-test-config
                              app-store:
                                config: react-native-android-ios-test-config
                              development:
                                config: react-native-android-ios-test-config
                              enterprise:
                                config: react-native-android-ios-test-config
                          SampleAppsReactNativeAndroid-tvOS:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-test-config
                              app-store:
                                config: react-native-android-ios-test-config
                              development:
                                config: react-native-android-ios-test-config
                              enterprise:
                                config: react-native-android-ios-test-config
configs:
  react-native:
    react-native-android-ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests, builds and deploys the app using *Deploy to bitrise.io* Step.

            Next steps:
            - Set up an [Apple service with API key](https://devcenter.bitrise.io/en/accounts/connecting-to-services/connecting-to-an-apple-service-with-api-key.html).
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: test
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
              - automatic_code_signing: api-key
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Runs tests.

            Next steps:
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: test
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, sampleAppsReactNativeIosAndAndroidVersions...)

var sampleAppsReactNativeIosAndAndroidNoTestVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.NpmVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsReactNativeIosAndAndroidNoTestResultYML = fmt.Sprintf(`options:
  react-native:
    title: React Native project directory
    summary: Path of the directory containing the project's `+"`package.json`"+` file.
    env_key: WORKDIR
    type: selector
    value_map:
      .:
        title: The root directory of an Android project
        summary: The root directory of your Android project, stored as an Environment
          Variable. In your Workflows, you can specify paths relative to this path.
          You can change this at any time.
        env_key: PROJECT_LOCATION
        type: selector
        value_map:
          android:
            title: Module
            summary: Modules provide a container for your Android project's source
              code, resource files, and app level settings, such as the module-level
              build file and Android manifest file. Each module can be independently
              built, tested, and debugged. You can add new modules to your Bitrise
              builds at any time.
            env_key: MODULE
            type: user_input
            value_map:
              app:
                title: Variant
                summary: Your Android build variant. You can add variants at any time,
                  as well as further configure your existing variants later.
                env_key: VARIANT
                type: user_input_optional
                value_map:
                  Debug:
                    title: Project or Workspace path
                    summary: The location of your Xcode project or Xcode workspace
                      files, stored as an Environment Variable. In your Workflows,
                      you can specify paths relative to this path.
                    env_key: BITRISE_PROJECT_PATH
                    type: selector
                    value_map:
                      ios/SampleAppsReactNativeAndroid.xcodeproj:
                        title: Scheme name
                        summary: An Xcode scheme defines a collection of targets to
                          build, a configuration to use when building, and a collection
                          of tests to execute. Only shared schemes are detected automatically
                          but you can use any scheme as a target on Bitrise. You can
                          change the scheme at any time in your Env Vars.
                        env_key: BITRISE_SCHEME
                        type: selector
                        value_map:
                          SampleAppsReactNativeAndroid:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-config
                              app-store:
                                config: react-native-android-ios-config
                              development:
                                config: react-native-android-ios-config
                              enterprise:
                                config: react-native-android-ios-config
                          SampleAppsReactNativeAndroid-tvOS:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-config
                              app-store:
                                config: react-native-android-ios-config
                              development:
                                config: react-native-android-ios-config
                              enterprise:
                                config: react-native-android-ios-config
configs:
  react-native:
    react-native-android-ios-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests, builds and deploys the app using *Deploy to bitrise.io* Step.

            Next steps:
            - Set up an [Apple service with API key](https://devcenter.bitrise.io/en/accounts/connecting-to-services/connecting-to-an-apple-service-with-api-key.html).
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
              - automatic_code_signing: api-key
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Installs dependencies.

            Next steps:
            - Add tests to your project and configure the workflow to run them.
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, sampleAppsReactNativeIosAndAndroidNoTestVersions...)

var sampleAppsReactNativeIosAndAndroidYarnVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsReactNativeIosAndAndroidYarnResultYML = fmt.Sprintf(`options:
  react-native:
    title: React Native project directory
    summary: Path of the directory containing the project's `+"`package.json`"+` file.
    env_key: WORKDIR
    type: selector
    value_map:
      .:
        title: The root directory of an Android project
        summary: The root directory of your Android project, stored as an Environment
          Variable. In your Workflows, you can specify paths relative to this path.
          You can change this at any time.
        env_key: PROJECT_LOCATION
        type: selector
        value_map:
          android:
            title: Module
            summary: Modules provide a container for your Android project's source
              code, resource files, and app level settings, such as the module-level
              build file and Android manifest file. Each module can be independently
              built, tested, and debugged. You can add new modules to your Bitrise
              builds at any time.
            env_key: MODULE
            type: user_input
            value_map:
              app:
                title: Variant
                summary: Your Android build variant. You can add variants at any time,
                  as well as further configure your existing variants later.
                env_key: VARIANT
                type: user_input_optional
                value_map:
                  Debug:
                    title: Project or Workspace path
                    summary: The location of your Xcode project or Xcode workspace
                      files, stored as an Environment Variable. In your Workflows,
                      you can specify paths relative to this path.
                    env_key: BITRISE_PROJECT_PATH
                    type: selector
                    value_map:
                      ios/SampleAppsReactNativeAndroid.xcodeproj:
                        title: Scheme name
                        summary: An Xcode scheme defines a collection of targets to
                          build, a configuration to use when building, and a collection
                          of tests to execute. Only shared schemes are detected automatically
                          but you can use any scheme as a target on Bitrise. You can
                          change the scheme at any time in your Env Vars.
                        env_key: BITRISE_SCHEME
                        type: selector
                        value_map:
                          SampleAppsReactNativeAndroid:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-test-yarn-config
                              app-store:
                                config: react-native-android-ios-test-yarn-config
                              development:
                                config: react-native-android-ios-test-yarn-config
                              enterprise:
                                config: react-native-android-ios-test-yarn-config
                          SampleAppsReactNativeAndroid-tvOS:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-test-yarn-config
                              app-store:
                                config: react-native-android-ios-test-yarn-config
                              development:
                                config: react-native-android-ios-test-yarn-config
                              enterprise:
                                config: react-native-android-ios-test-yarn-config
configs:
  react-native:
    react-native-android-ios-test-yarn-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests, builds and deploys the app using *Deploy to bitrise.io* Step.

            Next steps:
            - Set up an [Apple service with API key](https://devcenter.bitrise.io/en/accounts/connecting-to-services/connecting-to-an-apple-service-with-api-key.html).
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - yarn@%s:
              inputs:
              - workdir: $WORKDIR
              - command: test
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
              - automatic_code_signing: api-key
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Runs tests.

            Next steps:
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - yarn@%s:
              inputs:
              - workdir: $WORKDIR
              - command: test
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, sampleAppsReactNativeIosAndAndroidYarnVersions...)

var sampleAppsReactNativeJoplinVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.NpmVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsReactNativeJoplinResultYML = fmt.Sprintf(`options:
  react-native:
    title: React Native project directory
    summary: Path of the directory containing the project's `+"`package.json`"+` file.
    env_key: WORKDIR
    type: selector
    value_map:
      packages/app-mobile:
        title: The root directory of an Android project
        summary: The root directory of your Android project, stored as an Environment
          Variable. In your Workflows, you can specify paths relative to this path.
          You can change this at any time.
        env_key: PROJECT_LOCATION
        type: selector
        value_map:
          packages/app-mobile/android:
            title: Module
            summary: Modules provide a container for your Android project's source
              code, resource files, and app level settings, such as the module-level
              build file and Android manifest file. Each module can be independently
              built, tested, and debugged. You can add new modules to your Bitrise
              builds at any time.
            env_key: MODULE
            type: user_input
            value_map:
              app:
                title: Variant
                summary: Your Android build variant. You can add variants at any time,
                  as well as further configure your existing variants later.
                env_key: VARIANT
                type: user_input_optional
                value_map:
                  Debug:
                    title: Project or Workspace path
                    summary: The location of your Xcode project or Xcode workspace
                      files, stored as an Environment Variable. In your Workflows,
                      you can specify paths relative to this path.
                    env_key: BITRISE_PROJECT_PATH
                    type: selector
                    value_map:
                      packages/app-mobile/ios/Joplin.xcworkspace:
                        title: Scheme name
                        summary: An Xcode scheme defines a collection of targets to
                          build, a configuration to use when building, and a collection
                          of tests to execute. Only shared schemes are detected automatically
                          but you can use any scheme as a target on Bitrise. You can
                          change the scheme at any time in your Env Vars.
                        env_key: BITRISE_SCHEME
                        type: selector
                        value_map:
                          Joplin:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-pod-config
                              app-store:
                                config: react-native-android-ios-pod-config
                              development:
                                config: react-native-android-ios-pod-config
                              enterprise:
                                config: react-native-android-ios-pod-config
                          Joplin-tvOS:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-pod-config
                              app-store:
                                config: react-native-android-ios-pod-config
                              development:
                                config: react-native-android-ios-pod-config
                              enterprise:
                                config: react-native-android-ios-pod-config
                          ShareExtension:
                            title: Distribution method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_DISTRIBUTION_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: react-native-android-ios-pod-config
                              app-store:
                                config: react-native-android-ios-pod-config
                              development:
                                config: react-native-android-ios-pod-config
                              enterprise:
                                config: react-native-android-ios-pod-config
configs:
  react-native:
    react-native-android-ios-pod-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests, builds and deploys the app using *Deploy to bitrise.io* Step.

            Next steps:
            - Set up an [Apple service with API key](https://devcenter.bitrise.io/en/accounts/connecting-to-services/connecting-to-an-apple-service-with-api-key.html).
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
              - automatic_code_signing: api-key
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Installs dependencies.

            Next steps:
            - Add tests to your project and configure the workflow to run them.
            - Check out [Getting started with React Native apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-react-native-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, sampleAppsReactNativeJoplinVersions...)
