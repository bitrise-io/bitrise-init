package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestFastlane(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("fastlane")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("fastlane")
	{
		sampleAppDir := filepath.Join(tmpDir, "fastlane")
		sampleAppURL := "https://github.com/bitrise-samples/fastlane.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(fastlaneResultYML), strings.TrimSpace(result))
	}
}

const fastlaneResultYML = `options:
  fastlane:
    title: Working directory
    env_key: FASTLANE_WORK_DIR
    value_map:
      BitriseFastlaneSample:
        title: Fastlane lane
        env_key: FASTLANE_LANE
        value_map:
          test:
            config: fastlane-config
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      BitriseFastlaneSample/BitriseFastlaneSample.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          BitriseFastlaneSample:
            config: ios-test-config
configs:
  fastlane:
    fastlane-config: |
      format_version: 1.3.1
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      app:
        envs:
        - FASTLANE_XCODE_LIST_TIMEOUT: "120"
      trigger_map:
      - workflow: primary
        pattern: '*'
        is_pull_request_allowed: true
      workflows:
        primary:
          steps:
          - activate-ssh-key@3.1.1:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@3.4.1: {}
          - script@1.1.3:
              title: Do anything with Script step
          - certificate-and-profile-installer@1.8.1: {}
          - fastlane@2.2.0:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
          - deploy-to-bitrise-io@1.2.5: {}
  ios:
    ios-test-config: |
      format_version: 1.3.1
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          steps:
          - activate-ssh-key@3.1.1:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@3.4.1: {}
          - script@1.1.3:
              title: Do anything with Script step
          - certificate-and-profile-installer@1.8.1: {}
          - xcode-test@1.17.1:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive@1.10.1:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@1.2.5: {}
        primary:
          steps:
          - activate-ssh-key@3.1.1:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@3.4.1: {}
          - script@1.1.3:
              title: Do anything with Script step
          - certificate-and-profile-installer@1.8.1: {}
          - xcode-test@1.17.1:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@1.2.5: {}
warnings:
  fastlane: []
  ios: []
`
