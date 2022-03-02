package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestMacOS(t *testing.T) {
	tmpDir := t.TempDir()

	var testCases = []helper.TestCase{
		{
			"sample-apps-osx-10-11",
			"https://github.com/bitrise-samples/sample-apps-osx-10-11.git",
			"",
			sampleAppsOSX1011ResultYML,
			sampleAppsOSX1011Versions,
		},
	}

	helper.Execute(t, tmpDir, testCases)
}

// Expected results

var sampleAppsOSX1011Versions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestMacVersion,
	steps.XcodeArchiveMacVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeTestMacVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsOSX1011ResultYML = fmt.Sprintf(`options:
  macos:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
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
          - cache-pull@%s: {}
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
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  macos: []
warnings_with_recommendations:
  macos: []
`, sampleAppsOSX1011Versions...)
