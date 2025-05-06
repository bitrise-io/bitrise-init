package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestMacOS(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:             "sample-apps-osx-10-11",
			RepoURL:          "https://github.com/bitrise-samples/sample-apps-osx-10-11.git",
			ExpectedResult:   sampleAppsOSX1011ResultYML,
			ExpectedVersions: sampleAppsOSX1011Versions,
		},
		{
			Name:             "sample-spm-mac-project",
			RepoURL:          "https://github.com/bitrise-io/sample-spm-project.git",
			ExpectedResult:   sampleSPMMacProjectResultYML,
			ExpectedVersions: sampleSPMMacProjectVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var sampleAppsOSX1011Versions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestMacVersion,
	steps.XcodeArchiveMacVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestMacVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsOSX1011ResultYML = fmt.Sprintf(`options:
  macos:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      sample-apps-osx-10-11.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          sample-apps-osx-10-11:
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
                config: macos-test-config
              developer-id:
                config: macos-test-config
              development:
                config: macos-test-config
              none:
                config: macos-test-config
configs:
  macos:
    macos-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: macos
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
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
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
warnings:
  macos: []
warnings_with_recommendations:
  macos: []
`, sampleAppsOSX1011Versions...)

var sampleSPMMacProjectVersions = []interface{}{
	// iOS
	models.FormatVersion,

	// ios-spm-project-test-config
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

	// macOS
	models.FormatVersion,

	// macos-spm-project-test-config
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestMacVersion,
	steps.DeployToBitriseIoVersion,
}
var sampleSPMMacProjectResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      Package.swift:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          CoolFeature-Package:
            config: ios-spm-project-test-config
  macos:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      Package.swift:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          CoolFeature-Package:
            config: macos-spm-project-test-config
configs:
  ios:
    ios-spm-project-test-config: |
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
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
  macos:
    macos-spm-project-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: macos
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
  macos: []
warnings_with_recommendations:
  ios: []
  macos: []
 `, sampleSPMMacProjectVersions...)
