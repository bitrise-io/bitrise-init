package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestFastlane(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"fastlane",
			"https://github.com/bitrise-samples/fastlane.git",
			"",
			fastlaneResultYML,
			fastlaneVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var fastlaneVersions = []interface{}{
	// fastlane
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FastlaneVersion,
	steps.DeployToBitriseIoVersion,

	// ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,

	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,
}

var fastlaneResultYML = fmt.Sprintf(`options:
  fastlane:
    title: Project type
    summary: The type of your project. This determines what Steps are added to your
      automatically configured Workflows. You can, however, add any Steps to your
      Workflows at any time.
    type: selector
    value_map:
      ios:
        title: Working directory
        summary: The directory where your Fastfile is located.
        env_key: FASTLANE_WORK_DIR
        type: selector
        value_map:
          BitriseFastlaneSample:
            title: Fastlane lane
            summary: The lane that will be used in your builds, stored as an Environment
              Variable. You can change this at any time.
            env_key: FASTLANE_LANE
            type: selector
            value_map:
              ios test:
                config: fastlane-config_ios
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      BitriseFastlaneSample/BitriseFastlaneSample.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          BitriseFastlaneSample:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-test-config
              app-store:
                config: ios-test-config
              development:
                config: ios-test-config
              enterprise:
                config: ios-test-config
configs:
  fastlane:
    fastlane-config_ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - FASTLANE_XCODE_LIST_TIMEOUT: "120"
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - fastlane@%s:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
              - enable_cache: "no"
          - deploy-to-bitrise-io@%s: {}
  ios:
    ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests_parallel:
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
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
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
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
                  BITRISE_XCTESTRUN_FILE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
warnings:
  fastlane: []
  ios: []
warnings_with_recommendations:
  fastlane: []
  ios: []
`, fastlaneVersions...)
