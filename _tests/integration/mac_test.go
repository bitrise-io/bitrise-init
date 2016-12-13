package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func _TestMacOS(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__macos__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("sample-apps-osx-10-11")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-osx-10-11")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-osx-10-11.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsOSX1011ResultYML), strings.TrimSpace(result))
	}
}

var sampleAppsOSX1011ResultYML = fmt.Sprintf(`options:
  macos:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      sample-apps-osx-10-11.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          sample-apps-osx-10-11:
            config: macos-test-config
configs:
  macos:
    macos-test-config: |
      format_version: %s
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      app: {}
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
warnings:
  macos: []
`, models.FormatVersion,
	steps.ActivateSSHKeyVersion, steps.GitCloneVersion, steps.ScriptVersion, steps.CertificateAndProfileInstallerVersion, steps.XcodeTestMacVersion, steps.XcodeArchiveMacVersion, steps.DeployToBitriseIoVersion,
	steps.ActivateSSHKeyVersion, steps.GitCloneVersion, steps.ScriptVersion, steps.CertificateAndProfileInstallerVersion, steps.XcodeTestMacVersion, steps.DeployToBitriseIoVersion)
