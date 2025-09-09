package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestKotlinMultiplatform(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:             "Taskman",
			RepoURL:          "https://github.com/bitrise-io/kotlin-multiplatform-sample-taskman.git",
			Branch:           "main",
			ExpectedResult:   kmpTaskmanResultYaml,
			ExpectedVersions: kmpTaskmanResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

var kmpTaskmanResultYaml = fmt.Sprintf(`options:
  kotlin-multiplatform:
    title: The root directory of the Kotlin Multiplatform project.
    summary: The root directory of the Kotlin Multiplatform project, which contains
      all source files from your project, as well as Gradle files, including the Gradle
      Wrapper (gradlew) file.
    env_key: PROJECT_ROOT_DIR
    type: selector
    value_map:
      ./:
        title: Android Application Module
        summary: The name of the Android application module to build.
        env_key: MODULE
        type: selector
        value_map:
          composeApp:
            title: Android Application Variant
            summary: The name of the Android application variant to build.
            env_key: VARIANT
            type: user_input_optional
            value_map:
              "":
                title: iOS Application Project or Workspace path
                summary: The path of iOS application Xcode project or workspace to
                  build.
                env_key: BITRISE_PROJECT_PATH
                type: selector
                value_map:
                  ./iosApp/iosApp.xcodeproj:
                    title: iOS Application Scheme
                    summary: The name of the iOS application scheme to build.
                    env_key: BITRISE_SCHEME
                    type: selector
                    value_map:
                      iosApp:
                        title: iOS Application Distribution method
                        summary: The export method to use to build the iOS application
                          IPA file.
                        env_key: BITRISE_DISTRIBUTION_METHOD
                        type: selector
                        value_map:
                          ad-hoc:
                            config: kotlin-multiplatform-config
                          app-store:
                            config: kotlin-multiplatform-config
                          development:
                            config: kotlin-multiplatform-config
                          enterprise:
                            config: kotlin-multiplatform-config
configs:
  kotlin-multiplatform:
    kotlin-multiplatform-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: kotlin-multiplatform
      pipelines:
        build:
          workflows:
            android_build: {}
            ios_build: {}
      workflows:
        android_build:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-gradle-cache@%s: {}
          - activate-build-cache-for-gradle@%s: {}
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_ROOT_DIR
              - module: $MODULE
              - variant: $VARIANT
          - sign-apk@%s:
              run_if: '{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}'
          - save-gradle-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
        ios_build:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-gradle-cache@%s: {}
          - activate-build-cache-for-gradle@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
              - automatic_code_signing: api-key
          - save-gradle-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-gradle-cache@%s: {}
          - activate-build-cache-for-gradle@%s: {}
          - gradle-unit-test@%s:
              inputs:
              - project_root_dir: $PROJECT_ROOT_DIR
          - save-gradle-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  kotlin-multiplatform: []
warnings_with_recommendations:
  kotlin-multiplatform: []
`, kmpTaskmanResultVersions...)

var kmpTaskmanResultVersions = []interface{}{
	models.FormatVersion,
	// android_build
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreGradleVersion,
	steps.ActivateBuildCacheForGradleVersion,
	steps.AndroidBuildVersion,
	steps.SignAPKVersion,
	steps.CacheSaveGradleVersion,
	steps.DeployToBitriseIoVersion,
	// ios_build
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreGradleVersion,
	steps.ActivateBuildCacheForGradleVersion,
	steps.XcodeArchiveVersion,
	steps.CacheSaveGradleVersion,
	steps.DeployToBitriseIoVersion,
	// run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreGradleVersion,
	steps.ActivateBuildCacheForGradleVersion,
	steps.GradleUnitTestVersion,
	steps.CacheSaveGradleVersion,
	steps.DeployToBitriseIoVersion,
}
