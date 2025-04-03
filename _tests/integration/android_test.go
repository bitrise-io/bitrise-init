package integration

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/output"
	"github.com/bitrise-io/bitrise-init/scanner"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/stretchr/testify/require"
)

func TestAndroid(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"sample-apps-android-sdk22",
			"https://github.com/bitrise-samples/sample-apps-android-sdk22.git",
			"",
			sampleAppsAndroid22ResultYML,
			sampleAppsAndroid22Versions,
		},
		{
			"android-non-executable-gradlew",
			"https://github.com/bitrise-samples/android-non-executable-gradlew.git",
			"",
			androidNonExecutableGradlewResultYML,
			androidNonExecutableGradlewVersions,
		},
		{
			"android-sdk22-subdir",
			"https://github.com/bitrise-samples/sample-apps-android-sdk22-subdir",
			"",
			sampleAppsAndroidSDK22SubdirResultYML,
			sampleAppsAndroidSDK22SubdirVersions,
		},
		{
			"android-gradle-kotlin-dsl",
			"https://github.com/bitrise-samples/android-gradle-kotlin-dsl",
			"",
			sampleAppsKotlinDSLResultYML,
			sampleAppsKotlinDSLVersions,
		},
	}

	helper.Execute(t, testCases)
}

func TestMissingGradlewWrapper(t *testing.T) {
	tmpDir := t.TempDir()
	testName := "android-sdk22-no-gradlew"
	sampleAppDir := filepath.Join(tmpDir, testName)
	sampleAppURL := "https://github.com/bitrise-samples/android-sdk22-no-gradlew.git"
	helper.GitClone(t, sampleAppDir, sampleAppURL)

	_, err := scanner.GenerateAndWriteResults(sampleAppDir, sampleAppDir, output.YAMLFormat)
	require.EqualError(t, err, "No known platform detected")

	scanResultPth := filepath.Join(sampleAppDir, "result.yml")

	result, err := fileutil.ReadStringFromFile(scanResultPth)
	require.NoError(t, err)

	helper.ValidateConfigExpectation(t, testName, strings.TrimSpace(sampleAppsSDK22NoGradlewResultYML), strings.TrimSpace(result))
}

// Expected results

var sampleAppsAndroidSDK22SubdirVersions = []interface{}{
	// android-config
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
}

var sampleAppsAndroidSDK22SubdirResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    summary: The root directory of your Android project, stored as an Environment
      Variable. In your Workflows, you can specify paths relative to this path. You
      can change this at any time.
    env_key: PROJECT_LOCATION
    type: selector
    value_map:
      src:
        title: Module
        summary: Modules provide a container for your Android project's source code,
          resource files, and app level settings, such as the module-level build file
          and Android manifest file. Each module can be independently built, tested,
          and debugged. You can add new modules to your Bitrise builds at any time.
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
              "":
                config: android-config
                icons:
                - 5d50523f459dfaf760b7adeb5113216474b5d659a5ef66695239626376be7c89.png
configs:
  android:
    android-config: |
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
warnings:
  android: []
warnings_with_recommendations:
  android: []
`, sampleAppsAndroidSDK22SubdirVersions...)

var sampleAppsSDK22NoGradlewResultYML = `warnings:
  android: []
errors_with_recommendations:
  general:
  - error: No known platform detected
    recommendations:
      DetailedError:
        title: We couldn't recognize your platform.
        description: Our auto-configurator supports kotlin-multiplatform, react-native, flutter, ionic,
          cordova, ios, macos, android, node-js, fastlane projects. If you're adding
          something else, skip this step and configure your Workflow manually.
      NoPlatformDetected: true
warnings_with_recommendations:
  android:
  - error: |-
      <b>No Gradle Wrapper (gradlew) found.</b>
      Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure
      that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>
    recommendations:
      DetailedError:
        title: We couldn't find your Gradle Wrapper. Please make sure there is a gradlew
          file in your project's root directory.
        description: The Gradle Wrapper ensures that the right Gradle version is installed
          and used for the build. You can find out more about <a target="_blank" href="https://docs.gradle.org/current/userguide/gradle_wrapper.html">the
          Gradle Wrapper in the Gradle docs</a>.
`

var sampleAppsAndroid22Versions = []interface{}{
	// android-config
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
}

var sampleAppsAndroid22ResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    summary: The root directory of your Android project, stored as an Environment
      Variable. In your Workflows, you can specify paths relative to this path. You
      can change this at any time.
    env_key: PROJECT_LOCATION
    type: selector
    value_map:
      .:
        title: Module
        summary: Modules provide a container for your Android project's source code,
          resource files, and app level settings, such as the module-level build file
          and Android manifest file. Each module can be independently built, tested,
          and debugged. You can add new modules to your Bitrise builds at any time.
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
              "":
                config: android-config
                icons:
                - 81af22c35b03b30a1931a6283349eae094463aa69c52af3afe804b40dbe6dc12.png
configs:
  android:
    android-config: |
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
warnings:
  android: []
warnings_with_recommendations:
  android: []
`, sampleAppsAndroid22Versions...)

var androidNonExecutableGradlewVersions = []interface{}{
	// android-config
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
}

var androidNonExecutableGradlewResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    summary: The root directory of your Android project, stored as an Environment
      Variable. In your Workflows, you can specify paths relative to this path. You
      can change this at any time.
    env_key: PROJECT_LOCATION
    type: selector
    value_map:
      .:
        title: Module
        summary: Modules provide a container for your Android project's source code,
          resource files, and app level settings, such as the module-level build file
          and Android manifest file. Each module can be independently built, tested,
          and debugged. You can add new modules to your Bitrise builds at any time.
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
              "":
                config: android-config
                icons:
                - 81af22c35b03b30a1931a6283349eae094463aa69c52af3afe804b40dbe6dc12.png
configs:
  android:
    android-config: |
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
warnings:
  android: []
warnings_with_recommendations:
  android: []
`, androidNonExecutableGradlewVersions...)

var sampleAppsKotlinDSLVersions = []interface{}{
	// android-config-kts
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
}

var sampleAppsKotlinDSLResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    summary: The root directory of your Android project, stored as an Environment
      Variable. In your Workflows, you can specify paths relative to this path. You
      can change this at any time.
    env_key: PROJECT_LOCATION
    type: selector
    value_map:
      .:
        title: Module
        summary: Modules provide a container for your Android project's source code,
          resource files, and app level settings, such as the module-level build file
          and Android manifest file. Each module can be independently built, tested,
          and debugged. You can add new modules to your Bitrise builds at any time.
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
              "":
                config: android-config-kts
                icons:
                - 81af22c35b03b30a1931a6283349eae094463aa69c52af3afe804b40dbe6dc12.png
                - 3a50cbe24812ec6ef995f7142267bf67059d3e73e6b042873043b00354dbfde0.png
configs:
  android:
    android-config-kts: |
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
          - activate-ssh-key@%s: {}
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
warnings:
  android: []
warnings_with_recommendations:
  android: []

`, sampleAppsKotlinDSLVersions...)
