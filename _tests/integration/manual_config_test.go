package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/output"
	"github.com/bitrise-io/bitrise-init/scanner"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/fileutil"
)

func TestManualConfig(t *testing.T) {
	tmpDir := t.TempDir()
	testName := "manual-config"
	manualConfigDir := filepath.Join(tmpDir, testName)
	require.NoError(t, os.MkdirAll(manualConfigDir, 0777))
	fmt.Printf("manualConfigDir: %s\n", manualConfigDir)

	scanResult, err := scanner.ManualConfig()
	require.NoError(t, err)

	outputPth, err := output.WriteToFile(scanResult, output.YAMLFormat, filepath.Join(manualConfigDir, "result"))
	require.NoError(t, err)

	result, err := fileutil.ReadStringFromFile(outputPth)
	require.NoError(t, err)

	helper.ValidateConfigExpectation(t, testName, strings.TrimSpace(customConfigResultYML), strings.TrimSpace(result), customConfigVersions...)
}

// Expected results

var customConfigVersions = []interface{}{
	// default-android-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.ChangeAndroidVersionCodeAndVersionNameVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.AndroidBuildVersion,
	steps.SignAPKVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreGradleVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AvdManagerVersion,
	steps.WaitForAndroidEmulatorVersion,
	steps.GradleRunnerVersion,
	steps.CacheSaveGradleVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreGradleVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidUnitTestVersion,
	steps.CacheSaveGradleVersion,
	steps.DeployToBitriseIoVersion,

	// default-android-config-kts
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.ChangeAndroidVersionCodeAndVersionNameVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.AndroidBuildVersion,
	steps.SignAPKVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreGradleVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AvdManagerVersion,
	steps.WaitForAndroidEmulatorVersion,
	steps.GradleRunnerVersion,
	steps.CacheSaveGradleVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreGradleVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidUnitTestVersion,
	steps.CacheSaveGradleVersion,
	steps.DeployToBitriseIoVersion,

	// cordova
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.CordovaArchiveVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,

	// fastlane
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FastlaneVersion,
	steps.DeployToBitriseIoVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FastlaneVersion,
	steps.DeployToBitriseIoVersion,

	// flutter
	// flutter-config-test-android-2
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,

	// flutter-config-test-both-0
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,

	// flutter-config-test-ios-1
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,

	// ionic
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.IonicArchiveVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,

	// ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCocoapodsVersion,
	steps.CacheRestoreSPMVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeBuildForTestVersion,
	steps.CacheSaveCocoapodsVersion,
	steps.CacheSaveSPMVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCocoapodsVersion,
	steps.CacheRestoreSPMVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveCocoapodsVersion,
	steps.CacheSaveSPMVersion,
	steps.DeployToBitriseIoVersion,

	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,

	// java
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.GradleUnitTestVersion,
	steps.DeployToBitriseIoVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.DeployToBitriseIoVersion,

	// kotlin-multiplatform
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.GradleUnitTestVersion,
	steps.DeployToBitriseIoVersion,

	// macos
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestMacVersion,
	steps.XcodeArchiveMacVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCocoapodsVersion,
	steps.CacheRestoreSPMVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestMacVersion,
	steps.CacheSaveCocoapodsVersion,
	steps.CacheSaveSPMVersion,
	steps.DeployToBitriseIoVersion,

	// node-js
	// npm
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
	// yarn
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,

	// other
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.DeployToBitriseIoVersion,

	// react native
	// default-react-native-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// default-react-native-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,

	// default-react-native-expo-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.RunEASBuildVersion,
	steps.DeployToBitriseIoVersion,

	// default-react-native-expo-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var customConfigResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    summary: The root directory of your Android project, stored as an Environment
      Variable. In your Workflows, you can specify paths relative to this path. You
      can change this at any time.
    env_key: PROJECT_LOCATION
    type: user_input
    value_map:
      "":
        title: Module
        summary: Modules provide a container for your Android project's source code,
          resource files, and app level settings, such as the module-level build file
          and Android manifest file. Each module can be independently built, tested,
          and debugged. You can add new modules to your Bitrise builds at any time.
        env_key: MODULE
        type: user_input
        value_map:
          "":
            title: Variant
            summary: Your Android build variant. You can add variants at any time,
              as well as further configure your existing variants later.
            env_key: VARIANT
            type: user_input_optional
            value_map:
              "":
                title: Does your app use Kotlin build scripts?
                summary: The workflow configuration slightly differs based on what
                  language (Groovy or Kotlin) you used in your build scripts.
                type: selector
                value_map:
                  "no":
                    config: default-android-config
                  "yes":
                    config: default-android-config-kts
  cordova:
    title: Directory of the Cordova config.xml file
    summary: The working directory of your Cordova project is where you store your
      config.xml file. In your Workflows, you can specify paths relative to this path.
      You can change this at any time.
    env_key: CORDOVA_WORK_DIR
    type: user_input
    value_map:
      "":
        title: The platform to use in cordova-cli commands
        summary: The target platform for your build, stored as an Environment Variable.
          Your options are iOS, Android, or both. You can change this in your Env
          Vars at any time.
        env_key: CORDOVA_PLATFORM
        type: selector
        value_map:
          android:
            config: default-cordova-config
          ios:
            config: default-cordova-config
          ios,android:
            config: default-cordova-config
  fastlane:
    title: Working directory
    summary: The directory where your Fastfile is located.
    env_key: FASTLANE_WORK_DIR
    type: user_input
    value_map:
      "":
        title: Fastlane lane
        summary: The lane that will be used in your builds, stored as an Environment
          Variable. You can change this at any time.
        env_key: FASTLANE_LANE
        type: user_input
        value_map:
          "":
            title: Project type
            summary: The project type of the app you added to Bitrise.
            type: selector
            value_map:
              android:
                config: default-fastlane-android-config
              ios:
                config: default-fastlane-ios-config
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: user_input
    value_map:
      "":
        title: Platform
        summary: The target platform for your first build. Your options are iOS, Android,
          both, or neither. You can change this in your Env Vars at any time.
        type: selector
        value_map:
          android:
            config: flutter-config-test-android-2
          both:
            config: flutter-config-test-both-0
          ios:
            config: flutter-config-test-ios-1
  ionic:
    title: Directory of the Ionic config.xml file
    summary: The working directory of your Ionic project is where you store your config.xml
      file. This location is stored as an Environment Variable. In your Workflows,
      you can specify paths relative to this path. You can change this at any time.
    env_key: IONIC_WORK_DIR
    type: user_input
    value_map:
      "":
        title: The platform to use in ionic-cli commands
        summary: The target platform for your builds, stored as an Environment Variable.
          Your options are iOS, Android, or both. You can change this in your Env
          Vars at any time.
        env_key: IONIC_PLATFORM
        type: selector
        value_map:
          android:
            config: default-ionic-config
          ios:
            config: default-ionic-config
          ios,android:
            config: default-ionic-config
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: user_input
    value_map:
      "":
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: user_input
        value_map:
          "":
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: default-ios-config
              app-store:
                config: default-ios-config
              development:
                config: default-ios-config
              enterprise:
                config: default-ios-config
  java:
    title: Build tool
    summary: 'The build tool used in the project. Supported options: Gradle, Maven.'
    type: selector
    value_map:
      Gradle:
        title: The root directory of the Gradle project.
        summary: The root directory of the Gradle project, which contains all source
          files from your project, as well as Gradle files, including the Gradle Wrapper
          (`+"`gradlew`"+`) file.
        env_key: PROJECT_ROOT_DIR
        type: user_input
        value_map:
          "":
            config: default-java-gradle-config
      Maven:
        title: The root directory of the Maven project.
        summary: The root directory of the Maven project, which contains all source
          files from your project, as well as Maven files, including the Maven Wrapper
          (`+"`mvn`"+`) file.
        env_key: PROJECT_ROOT_DIR
        type: user_input
        value_map:
          "":
            config: default-java-maven-config
  kotlin-multiplatform:
    title: The root directory of the Gradle project.
    summary: The root directory of the Gradle project, which contains all source files
      from your project, as well as Gradle files, including the Gradle Wrapper (`+"`gradlew`"+`)
      file.
    env_key: PROJECT_ROOT_DIR
    type: user_input
    value_map:
      "":
        config: default-kotlin-multiplatform-config
  macos:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: user_input
    value_map:
      "":
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: user_input
        value_map:
          "":
            title: |-
              Application export method
              NOTE: `+"`none`"+` means: Export a copy of the application without re-signing.
            summary: The export method used to create an .app file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .app files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
            type: selector
            value_map:
              app-store:
                config: default-macos-config
              developer-id:
                config: default-macos-config
              development:
                config: default-macos-config
              none:
                config: default-macos-config
  node-js:
    title: Project Directory
    summary: The directory containing the package.json file
    env_key: NODEJS_PROJECT_DIR
    type: user_input
    value_map:
      "":
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          npm:
            config: default-node-js-npm-config
          yarn:
            config: default-node-js-yarn-config
  react-native:
    title: Is this an [Expo](https://expo.dev)-based React Native project?
    summary: |-
      Default deploy workflow runs builds on Expo Application Services (EAS) for Expo-based React Native projects.
      Otherwise native iOS and Android build steps will be used.
    type: selector
    value_map:
      "no":
        title: The root directory of an Android project
        summary: The root directory of your Android project, stored as an Environment
          Variable. In your Workflows, you can specify paths relative to this path.
          You can change this at any time.
        env_key: PROJECT_LOCATION
        type: user_input
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
                    summary: The location of your Xcode project, Xcode workspace or
                      SPM project files stored as an Environment Variable. In your
                      Workflows, you can specify paths relative to this path.
                    env_key: BITRISE_PROJECT_PATH
                    type: user_input
                    value_map:
                      "":
                        title: Scheme name
                        summary: An Xcode scheme defines a collection of targets to
                          build, a configuration to use when building, and a collection
                          of tests to execute. Only shared schemes are detected automatically
                          but you can use any scheme as a target on Bitrise. You can
                          change the scheme at any time in your Env Vars.
                        env_key: BITRISE_SCHEME
                        type: user_input
                        value_map:
                          "":
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
                                config: default-react-native-config
                              app-store:
                                config: default-react-native-config
                              development:
                                config: default-react-native-config
                              enterprise:
                                config: default-react-native-config
      "yes":
        title: Expo project directory
        summary: Path of the directory containing the project's  `+"`package.json`"+` and
          app configuration file (`+"`app.json`"+`, `+"`app.config.js`"+`, `+"`app.config.ts`"+`).
        env_key: WORKDIR
        type: user_input
        value_map:
          "":
            title: Platform to build
            summary: Which platform should be built by the deploy workflow?
            env_key: PLATFORM
            type: selector
            value_map:
              all:
                config: default-react-native-expo-config
              android:
                config: default-react-native-expo-config
              ios:
                config: default-react-native-expo-config
configs:
  android:
    default-android-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            run_instrumented_tests:
              parallel: $TEST_SHARD_COUNT
      workflows:
        build_apk:
          summary: Run your Android unit tests and create an APK file to install your app
            on a device or share it with your team.
          description: The workflow will first clone your Git repository, install Android
            tools, set the project's version code based on the build number, run Android
            lint and unit tests, build the project's APK file and save it.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - change-android-versioncode-and-versionname@%s:
              inputs:
              - build_gradle_path: $PROJECT_LOCATION/$MODULE/build.gradle
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
          - sign-apk@%s:
              run_if: '{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}'
          - deploy-to-bitrise-io@%s: {}
        run_instrumented_tests:
          summary: Run your Android instrumented tests and get the test report.
          description: The workflow will first clone your Git repository, cache your Gradle
            dependencies, install Android tools, run your Android instrumented tests and
            save the test report.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-gradle-cache@%s: {}
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - avd-manager@%s: {}
          - wait-for-android-emulator@%s: {}
          - gradle-runner@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
              - gradle_task: |-
                  connectedAndroidTest \
                    -Pandroid.testInstrumentationRunnerArguments.numShards=$BITRISE_IO_PARALLEL_TOTAL \
                    -Pandroid.testInstrumentationRunnerArguments.shardIndex=$BITRISE_IO_PARALLEL_INDEX
          - save-gradle-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
        run_tests:
          summary: Run your Android unit tests and get the test report.
          description: The workflow will first clone your Git repository, cache your Gradle
            dependencies, install Android tools, run your Android unit tests and save the
            test report.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-gradle-cache@%s: {}
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - variant: $VARIANT
              - cache_level: none
          - save-gradle-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
    default-android-config-kts: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            run_instrumented_tests:
              parallel: $TEST_SHARD_COUNT
      workflows:
        build_apk:
          summary: Run your Android unit tests and create an APK file to install your app
            on a device or share it with your team.
          description: The workflow will first clone your Git repository, install Android
            tools, set the project's version code based on the build number, run Android
            lint and unit tests, build the project's APK file and save it.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
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
          - sign-apk@%s:
              run_if: '{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}'
          - deploy-to-bitrise-io@%s: {}
        run_instrumented_tests:
          summary: Run your Android instrumented tests and get the test report.
          description: The workflow will first clone your Git repository, cache your Gradle
            dependencies, install Android tools, run your Android instrumented tests and
            save the test report.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-gradle-cache@%s: {}
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - avd-manager@%s: {}
          - wait-for-android-emulator@%s: {}
          - gradle-runner@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
              - gradle_task: |-
                  connectedAndroidTest \
                    -Pandroid.testInstrumentationRunnerArguments.numShards=$BITRISE_IO_PARALLEL_TOTAL \
                    -Pandroid.testInstrumentationRunnerArguments.shardIndex=$BITRISE_IO_PARALLEL_INDEX
          - save-gradle-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
        run_tests:
          summary: Run your Android unit tests and get the test report.
          description: The workflow will first clone your Git repository, cache your Gradle
            dependencies, install Android tools, run your Android unit tests and save the
            test report.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-gradle-cache@%s: {}
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - variant: $VARIANT
              - cache_level: none
          - save-gradle-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
  cordova:
    default-cordova-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: cordova
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              title: npm install
              inputs:
              - workdir: $CORDOVA_WORK_DIR
              - command: install
          - generate-cordova-build-configuration@%s: {}
          - cordova-archive@%s:
              inputs:
              - workdir: $CORDOVA_WORK_DIR
              - platform: $CORDOVA_PLATFORM
              - target: emulator
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
  fastlane:
    default-fastlane-android-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      app:
        envs:
        - FASTLANE_XCODE_LIST_TIMEOUT: "120"
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - fastlane@%s:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
              - enable_cache: "no"
          - deploy-to-bitrise-io@%s: {}
    default-fastlane-ios-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - FASTLANE_XCODE_LIST_TIMEOUT: "120"
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - fastlane@%s:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
              - enable_cache: "no"
          - deploy-to-bitrise-io@%s: {}
  flutter:
    flutter-config-test-android-2: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        deploy:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html) for signing and deployment options.
            - Check out the [Code signing guide](https://devcenter.bitrise.io/en/code-signing.html) for iOS and Android
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: android
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Builds project and runs tests.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html).
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - restore-dart-cache@%s: {}
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
    flutter-config-test-both-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        deploy:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html) for signing and deployment options.
            - Check out the [Code signing guide](https://devcenter.bitrise.io/en/code-signing.html) for iOS and Android
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: both
              - ios_output_type: archive
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Builds project and runs tests.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html).
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - restore-dart-cache@%s: {}
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
    flutter-config-test-ios-1: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        deploy:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html) for signing and deployment options.
            - Check out the [Code signing guide](https://devcenter.bitrise.io/en/code-signing.html) for iOS and Android
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: ios
              - ios_output_type: archive
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Builds project and runs tests.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html).
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - restore-dart-cache@%s: {}
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
  ionic:
    default-ionic-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ionic
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              title: npm install
              inputs:
              - workdir: $IONIC_WORK_DIR
              - command: install
          - generate-cordova-build-configuration@%s: {}
          - ionic-archive@%s:
              inputs:
              - workdir: $IONIC_WORK_DIR
              - platform: $IONIC_PLATFORM
              - target: emulator
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
  ios:
    default-ios-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        archive_and_export_app:
          summary: Run your Xcode tests and create an IPA file to install your app on a
            device or share it with your team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests, export an IPA file
            from the project and save it.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        build_for_testing:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-cocoapods-cache@%s: {}
          - restore-spm-cache@%s: {}
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - save-cocoapods-cache@%s: {}
          - save-spm-cache@%s: {}
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-cocoapods-cache@%s: {}
          - restore-spm-cache@%s: {}
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - save-cocoapods-cache@%s: {}
          - save-spm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
  java:
    default-java-gradle-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: java
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - gradle-unit-test@%s:
              inputs:
              - project_root_dir: $PROJECT_ROOT_DIR
          - deploy-to-bitrise-io@%s: {}
    default-java-maven-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: java
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Run Maven tests
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  # fail if any commands fails
                  set -e
                  # make pipelines' return status equal the last command to exit with a non-zero status, or zero if all commands exit successfully
                  set -o pipefail
                  # debug log
                  set -x

                  ./mvnw test
              - working_dir: $PROJECT_ROOT_DIR
          - deploy-to-bitrise-io@%s: {}
  kotlin-multiplatform:
    default-kotlin-multiplatform-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: kotlin-multiplatform
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - gradle-unit-test@%s:
              inputs:
              - project_root_dir: $PROJECT_ROOT_DIR
          - deploy-to-bitrise-io@%s: {}
  macos:
    default-macos-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: macos
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-cocoapods-cache@%s: {}
          - restore-spm-cache@%s: {}
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - save-cocoapods-cache@%s: {}
          - save-spm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
  node-js:
    default-node-js-npm-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Install Node.js
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
          - restore-npm-cache@%s: {}
          - npm@%s:
              title: npm install
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: install
          - npm@%s:
              title: npm run lint
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run lint
          - npm@%s:
              title: npm run test
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run test
          - save-npm-cache@%s: {}
    default-node-js-yarn-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Install Node.js
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: install
          - yarn@%s:
              title: yarn run lint
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run lint
          - yarn@%s:
              title: yarn run test
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run test
          - save-npm-cache@%s: {}
  other:
    other-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: other
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - deploy-to-bitrise-io@%s: {}
  react-native:
    default-react-native-config: |
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
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - yarn@%s:
              title: yarn test
              inputs:
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
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - yarn@%s:
              title: yarn test
              inputs:
              - command: test
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
    default-react-native-expo-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests the app and runs a build on Expo Application Services (EAS).

            Next steps:
            - Configure the `+"`Run Expo Application Services (EAS) build`"+` Step's `+"`Access Token`"+` input.
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
            - For an alternative deploy workflow checkout the [(React Native) Expo: Build using Turtle CLI recipe](https://github.com/bitrise-io/workflow-recipes/blob/main/recipes/rn-expo-turtle-build.md).
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - workdir: $WORKDIR
              - command: install
          - yarn@%s:
              title: yarn test
              inputs:
              - workdir: $WORKDIR
              - command: test
          - run-eas-build@%s:
              inputs:
              - platform: $PLATFORM
              - work_dir: $WORKDIR
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Runs tests.

            Next steps:
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - workdir: $WORKDIR
              - command: install
          - yarn@%s:
              title: yarn test
              inputs:
              - workdir: $WORKDIR
              - command: test
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
`, customConfigVersions...)
