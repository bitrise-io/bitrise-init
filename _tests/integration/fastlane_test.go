package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestFastlane(t *testing.T) {
	tmpDir := t.TempDir()

	var testCases = []helper.TestCase{
		{
			"fastlane",
			"https://github.com/bitrise-samples/fastlane.git",
			"",
			fastlaneResultYML,
			fastlaneVersions,
		},
	}

	helper.Execute(t, tmpDir, testCases)
}

// Expected results

var fastlaneVersions = []interface{}{
	// fastlane
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FastlaneVersion,
	steps.DeployToBitriseIoVersion,

	// ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,
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
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
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
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - fastlane@%s:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
          - deploy-to-bitrise-io@%s: {}
  ios:
    ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      workflows:
        deploy:
          description: |
            The workflow tests, builds and deploys the app using *Deploy to bitrise.io* step.

            For testing the *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Set up [Connecting to an Apple service with API key](https://devcenter.bitrise.io/en/accounts/connecting-to-services/connecting-to-an-apple-service-with-api-key.html##).
            - Or further customise code signing following our [iOS code signing](https://devcenter.bitrise.io/en/code-signing/ios-code-signing.html) guide.
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  fastlane: []
  ios: []
warnings_with_recommendations:
  fastlane: []
  ios: []
`, fastlaneVersions...)
