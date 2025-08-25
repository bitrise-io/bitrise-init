package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestAndroid(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:             "bitrise-init-android-test-apps",
			RepoURL:          "https://github.com/bitrise-io/bitrise-init-android-test-apps",
			ExpectedResult:   monoRepoResultYML,
			ExpectedVersions: monoRepoVersions,
		},
		{
			Name:             "sample-apps-android-sdk22",
			RepoURL:          "https://github.com/bitrise-samples/sample-apps-android-sdk22.git",
			ExpectedResult:   sampleAppsAndroid22ResultYML,
			ExpectedVersions: sampleAppsAndroid22Versions,
		},
		{
			Name:             "android-gradle-kotlin-dsl",
			RepoURL:          "https://github.com/bitrise-samples/android-gradle-kotlin-dsl",
			ExpectedResult:   sampleAppsKotlinDSLResultYML,
			ExpectedVersions: sampleAppsKotlinDSLVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var monoRepoVersions = []interface{}{
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

var monoRepoResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    summary: The root directory of your Android project, stored as an Environment
      Variable. In your Workflows, you can specify paths relative to this path. You
      can change this at any time.
    env_key: PROJECT_LOCATION
    type: selector
    value_map:
      ./GroovyResponsiveViewsActivity:
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
      ./KotlinResponsiveViewsActivity:
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
              - build_root_directory: $PROJECT_LOCATION
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
              - build_root_directory: $PROJECT_LOCATION
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
`, monoRepoVersions...)

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
      ./:
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
              - build_root_directory: $PROJECT_LOCATION
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
      ./:
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
              - build_root_directory: $PROJECT_LOCATION
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
