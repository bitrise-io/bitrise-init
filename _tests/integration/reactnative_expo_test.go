package integration

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestReactNativeExpo(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__reactnative_expo__")
	require.NoError(t, err)

	t.Log("Managed workflow, no tests")
	{
		sampleAppDir := filepath.Join(tmpDir, "managed-notest")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-expo.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "Managed Expo Workflow, no tests", strings.TrimSpace(managedWorkflowNoTestsResultsYML), strings.TrimSpace(result))
	}

	t.Log("Managed workflow with tests")
	{
		sampleAppDir := filepath.Join(tmpDir, "managed-test")
		sampleAppURL := "https://github.com/bitrise-io/Bitrise-React-Native-Expo-Sample.git"

		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "Managed Expo Workflow with tests", strings.TrimSpace(managedWorkflowResultsYML), strings.TrimSpace(result))
	}

	t.Log("Bare workflow")
	{
		sampleAppDir := filepath.Join(tmpDir, "bare")
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-expo.git"
		gitCloneBranch(t, sampleAppDir, sampleAppURL, "bare")

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "Bare Expo Workflow", strings.TrimSpace(bareWorkflowResultYML), strings.TrimSpace(result))
	}
}

var managedExpoVersions = []interface{}{
	models.FormatVersion,
	// deploy
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.RunEASBuildVersion,
	steps.DeployToBitriseIoVersion,
	// primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.DeployToBitriseIoVersion,
}

var managedWorkflowNoTestsResultsYML = fmt.Sprintf(`options:
  react-native:
    title: Platform to build
    summary: Which platform should be built by the deploy workflow?
    env_key: PLATFORM
    type: selector
    value_map:
      all:
        config: react-native-expo-config
      android:
        config: react-native-expo-config
      ios:
        config: react-native-expo-config
configs:
  react-native:
    react-native-expo-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Runs a build on Expo Application Services (EAS).

            Next steps:
            - Configure the `+"`Run Expo Application Services (EAS) build`"+` Step's `+"`Access Token`"+` input.
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
            - For an alternative deploy workflow checkout the [(React Native) Expo: Build using Turtle CLI recipe](https://github.com/bitrise-io/workflow-recipes/blob/main/recipes/rn-expo-turtle-build.md).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              inputs:
              - command: install
          - run-eas-build@%s:
              inputs:
              - platform: $PLATFORM
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Installs dependencies.

            Next steps:
            - Add tests to your project and configure the workflow to run them.
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              inputs:
              - command: install
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, managedExpoVersions...)

var managedExpo2Versions = []interface{}{
	models.FormatVersion,
	// deploy
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.RunEASBuildVersion,
	steps.DeployToBitriseIoVersion,
	// primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.DeployToBitriseIoVersion,
}

var managedWorkflowResultsYML = fmt.Sprintf(`options:
  react-native:
    title: Platform to build
    summary: Which platform should be built by the deploy workflow?
    env_key: PLATFORM
    type: selector
    value_map:
      all:
        config: react-native-expo-config
      android:
        config: react-native-expo-config
      ios:
        config: react-native-expo-config
configs:
  react-native:
    react-native-expo-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests the app and runs a build on Expo Application Services (EAS).

            Next steps:
            - Configure the `+"`Run Expo Application Services (EAS) build`"+` Step's `+"`Access Token`"+` input.
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
            - For an alternative deploy workflow checkout the [(React Native) Expo: Build using Turtle CLI recipe](https://github.com/bitrise-io/workflow-recipes/blob/main/recipes/rn-expo-turtle-build.md).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              inputs:
              - command: install
          - yarn@%s:
              inputs:
              - command: test
          - run-eas-build@%s:
              inputs:
              - platform: $PLATFORM
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Runs tests.

            Next steps:
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              inputs:
              - command: install
          - yarn@%s:
              inputs:
              - command: test
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, managedExpo2Versions...)

// Bare workflow is the same as react-native with native projects
var sampleAppsExpoBareVersions = []interface{}{
	models.FormatVersion,
	// deploy
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.RunEASBuildVersion,
	steps.DeployToBitriseIoVersion,
	// primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.DeployToBitriseIoVersion,
}

var bareWorkflowResultYML = fmt.Sprintf(`options:
  react-native:
    title: Platform to build
    summary: Which platform should be built by the deploy workflow?
    env_key: PLATFORM
    type: selector
    value_map:
      all:
        config: react-native-expo-config
      android:
        config: react-native-expo-config
      ios:
        config: react-native-expo-config
configs:
  react-native:
    react-native-expo-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      workflows:
        deploy:
          description: |
            Tests the app and runs a build on Expo Application Services (EAS).

            Next steps:
            - Configure the `+"`Run Expo Application Services (EAS) build`"+` Step's `+"`Access Token`"+` input.
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
            - For an alternative deploy workflow checkout the [(React Native) Expo: Build using Turtle CLI recipe](https://github.com/bitrise-io/workflow-recipes/blob/main/recipes/rn-expo-turtle-build.md).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              inputs:
              - command: install
          - yarn@%s:
              inputs:
              - command: test
          - run-eas-build@%s:
              inputs:
              - platform: $PLATFORM
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Runs tests.

            Next steps:
            - Check out [Getting started with Expo apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-expo-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              inputs:
              - command: install
          - yarn@%s:
              inputs:
              - command: test
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, sampleAppsExpoBareVersions...)
