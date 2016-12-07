package integration

import (
	"os"
	"path/filepath"
	"testing"

	"strings"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestIOS(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__ios__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("ios-no-shared-schemes")
	{
		sampleAppDir := filepath.Join(tmpDir, "ios-no-shared-scheme")
		sampleAppURL := "https://github.com/bitrise-samples/ios-no-shared-schemes.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(iosNoSharedSchemesResultYML), strings.TrimSpace(result))
	}

	t.Log("ios-cocoapods-at-root")
	{
		sampleAppDir := filepath.Join(tmpDir, "ios-cocoapods-at-root")
		sampleAppURL := "https://github.com/bitrise-samples/ios-cocoapods-at-root.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(iosCocoapodsAtRootResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-ios-simple-objc")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-ios-simple-objc")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-ios-simple-objc.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsIosSimpleObjcResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-ios-watchkit")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-ios-watchkit")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-ios-watchkit.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsIosWatchkitResultYML), strings.TrimSpace(result))
	}
}

const sampleAppsIosWatchkitResultYML = `options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      watch-test.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          Complication - watch-test WatchKit App:
            config: ios-config
          Glance - watch-test WatchKit App:
            config: ios-config
          Notification - watch-test WatchKit App:
            config: ios-config
          watch-test:
            config: ios-test-config
          watch-test WatchKit App:
            config: ios-config
configs:
  ios:
    ios-config: |
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
          - deploy-to-bitrise-io@1.2.5: {}
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
  ios: []
`

const sampleAppsIosSimpleObjcResultYML = `options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      ios-simple-objc/ios-simple-objc.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          ios-simple-objc:
            config: ios-test-config
configs:
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
  ios: []
`

const iosCocoapodsAtRootResultYML = `options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      iOSMinimalCocoaPodsSample.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          iOSMinimalCocoaPodsSample:
            config: ios-test-config
configs:
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
  ios: []
`

const iosNoSharedSchemesResultYML = `options:
  ios:
    title: Project (or Workspace) path
    env_key: BITRISE_PROJECT_PATH
    value_map:
      BitriseXcode7Sample.xcodeproj:
        title: Scheme name
        env_key: BITRISE_SCHEME
        value_map:
          BitriseXcode7Sample:
            config: ios-test-missing-shared-schemes-config
configs:
  ios:
    ios-test-missing-shared-schemes-config: |
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
          - recreate-user-schemes@0.9.4:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
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
          - recreate-user-schemes@0.9.4:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
          - xcode-test@1.17.1:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@1.2.5: {}
warnings:
  ios:
  - "No shared schemes found for project: BitriseXcode7Sample.xcodeproj.\n\tAutomatically
    generated schemes for this project.\n\tThese schemes may differ from the ones
    in your project.\n\tMake sure to <a href=\"https://developer.apple.com/library/ios/recipes/xcode_help-scheme_editor/Articles/SchemeManage.html\">share
    your schemes</a> for the expected behaviour."
`
