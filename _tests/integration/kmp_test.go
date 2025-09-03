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
    title: The root directory of the Gradle project.
    summary: The root directory of the Gradle project, which contains all source files
      from your project, as well as Gradle files, including the Gradle Wrapper (`+"`gradlew`"+`)
      file.
    env_key: PROJECT_ROOT_DIR
    type: selector
    value_map:
      ./:
        title: Module
        summary: Modules provide a container for your Android project's source code,
          resource files, and app level settings, such as the module-level build file
          and Android manifest file. Each module can be independently built, tested,
          and debugged. You can add new modules to your Bitrise builds at any time.
        env_key: MODULE
        type: selector
        value_map:
          composeApp:
            title: Variant
            summary: Your Android build variant. You can add variants at any time,
              as well as further configure your existing variants later.
            env_key: VARIANT
            type: user_input_optional
            value_map:
              "":
                title: Project or Workspace path
                summary: The location of your Xcode project, Xcode workspace or SPM
                  project files stored as an Environment Variable. In your Workflows,
                  you can specify paths relative to this path.
                env_key: BITRISE_PROJECT_PATH
                type: selector
                value_map:
                  ./iosApp/iosApp.xcodeproj:
                    title: Scheme name
                    summary: An Xcode scheme defines a collection of targets to build,
                      a configuration to use when building, and a collection of tests
                      to execute. Only shared schemes are detected automatically but
                      you can use any scheme as a target on Bitrise. You can change
                      the scheme at any time in your Env Vars.
                    env_key: BITRISE_SCHEME
                    type: selector
                    value_map:
                      iosApp:
                        title: Distribution method
                        summary: The export method used to create an .ipa file in
                          your builds, stored as an Environment Variable. You can
                          change this at any time, or even create several .ipa files
                          with different export methods in the same build.
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
          - activate-build-cache-for-gradle@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
              - automatic_code_signing: api-key
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
	steps.ActivateBuildCacheForGradleVersion,
	steps.XcodeArchiveVersion,
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
