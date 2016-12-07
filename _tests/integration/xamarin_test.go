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

func TestXamarin(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__xamarin__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("xamarin-sample-app")
	{
		sampleAppDir := filepath.Join(tmpDir, "xamarin-sample-app")
		sampleAppURL := "https://github.com/bitrise-samples/xamarin-sample-app.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(xamarinSampleAppResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-xamarin-ios")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-xamarin-ios")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-xamarin-ios.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsXamarinIosResultYML), strings.TrimSpace(result))
	}

	t.Log("sample-apps-xamarin-android")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-xamarin-android")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-xamarin-android.git"
		require.NoError(t, cmdex.GitClone(sampleAppURL, sampleAppDir))

		cmd := cmdex.NewCommand(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(sampleAppsXamarinAndroidResultYML), strings.TrimSpace(result))
	}
}

const sampleAppsXamarinAndroidResultYML = `options:
  xamarin:
    title: Path to the Xamarin Solution file
    env_key: BITRISE_PROJECT_PATH
    value_map:
      CreditCardValidator.Droid.sln:
        title: Xamarin solution configuration
        env_key: BITRISE_XAMARIN_CONFIGURATION
        value_map:
          Debug:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
          Release:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
configs:
  xamarin:
    xamarin-nuget-config: |
      format_version: 1.3.1
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
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
          - xamarin-user-management@1.0.3:
              run_if: .IsCI
          - nuget-restore@1.0.1: {}
          - xamarin-archive@1.1.1:
              inputs:
              - xamarin_solution: $BITRISE_PROJECT_PATH
              - xamarin_configuration: $BITRISE_XAMARIN_CONFIGURATION
              - xamarin_platform: $BITRISE_XAMARIN_PLATFORM
          - deploy-to-bitrise-io@1.2.5: {}
warnings:
  xamarin: []
`

const sampleAppsXamarinIosResultYML = `options:
  xamarin:
    title: Path to the Xamarin Solution file
    env_key: BITRISE_PROJECT_PATH
    value_map:
      CreditCardValidator.iOS.sln:
        title: Xamarin solution configuration
        env_key: BITRISE_XAMARIN_CONFIGURATION
        value_map:
          Debug:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
              iPhone:
                config: xamarin-nuget-config
              iPhoneSimulator:
                config: xamarin-nuget-config
          Release:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-config
              iPhone:
                config: xamarin-nuget-config
              iPhoneSimulator:
                config: xamarin-nuget-config
configs:
  xamarin:
    xamarin-nuget-config: |
      format_version: 1.3.1
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
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
          - xamarin-user-management@1.0.3:
              run_if: .IsCI
          - nuget-restore@1.0.1: {}
          - xamarin-archive@1.1.1:
              inputs:
              - xamarin_solution: $BITRISE_PROJECT_PATH
              - xamarin_configuration: $BITRISE_XAMARIN_CONFIGURATION
              - xamarin_platform: $BITRISE_XAMARIN_PLATFORM
          - deploy-to-bitrise-io@1.2.5: {}
warnings:
  xamarin: []
`

const xamarinSampleAppResultYML = `options:
  xamarin:
    title: Path to the Xamarin Solution file
    env_key: BITRISE_PROJECT_PATH
    value_map:
      XamarinSampleApp.sln:
        title: Xamarin solution configuration
        env_key: BITRISE_XAMARIN_CONFIGURATION
        value_map:
          Debug:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-components-config
              iPhone:
                config: xamarin-nuget-components-config
              iPhoneSimulator:
                config: xamarin-nuget-components-config
          Release:
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            value_map:
              Any CPU:
                config: xamarin-nuget-components-config
              iPhone:
                config: xamarin-nuget-components-config
              iPhoneSimulator:
                config: xamarin-nuget-components-config
configs:
  xamarin:
    xamarin-nuget-components-config: |
      format_version: 1.3.1
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
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
          - xamarin-user-management@1.0.3:
              run_if: .IsCI
          - nuget-restore@1.0.1: {}
          - xamarin-components-restore@0.9.0: {}
          - xamarin-archive@1.1.1:
              inputs:
              - xamarin_solution: $BITRISE_PROJECT_PATH
              - xamarin_configuration: $BITRISE_XAMARIN_CONFIGURATION
              - xamarin_platform: $BITRISE_XAMARIN_PLATFORM
          - deploy-to-bitrise-io@1.2.5: {}
warnings:
  xamarin: []
`
