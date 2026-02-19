package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestReactNativeExpo(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:             "managed-workflow-no-tests",
			RepoURL:          "https://github.com/bitrise-io/sample-apps-expo.git",
			Branch:           "",
			ExpectedResult:   managedWorkflowNoTestsResultsYML,
			ExpectedVersions: managedExpoVersions,
		},
		{
			Name:             "managed-workflow-with-tests",
			RepoURL:          "https://github.com/bitrise-io/Bitrise-React-Native-Expo-Sample.git",
			Branch:           "",
			ExpectedResult:   managedWorkflowResultsYML,
			ExpectedVersions: managedExpo2Versions,
		},
		{
			Name:             "bare-workflow",
			RepoURL:          "https://github.com/bitrise-io/sample-apps-expo.git",
			Branch:           "bare",
			ExpectedResult:   bareWorkflowResultYML,
			ExpectedVersions: sampleAppsExpoBareVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

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
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
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
            - Check out [Getting started with Expo apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-expo-projects.html).
            - For an alternative deploy workflow checkout the [(React Native) Expo: Build using Turtle CLI recipe](https://github.com/bitrise-io/workflow-recipes/blob/main/recipes/rn-expo-turtle-build.md).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              title: yarn install
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
            - Check out [Getting started with Expo apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-expo-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - save-npm-cache@%s: {}
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
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
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
            - Check out [Getting started with Expo apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-expo-projects.html).
            - For an alternative deploy workflow checkout the [(React Native) Expo: Build using Turtle CLI recipe](https://github.com/bitrise-io/workflow-recipes/blob/main/recipes/rn-expo-turtle-build.md).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - yarn@%s:
              title: yarn test
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
            - Check out [Getting started with Expo apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-expo-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - yarn@%s:
              title: yarn test
              inputs:
              - command: test
          - save-npm-cache@%s: {}
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
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
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
            - Check out [Getting started with Expo apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-expo-projects.html).
            - For an alternative deploy workflow checkout the [(React Native) Expo: Build using Turtle CLI recipe](https://github.com/bitrise-io/workflow-recipes/blob/main/recipes/rn-expo-turtle-build.md).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - yarn@%s:
              title: yarn test
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
            - Check out [Getting started with Expo apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-expo-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - yarn@%s:
              title: yarn test
              inputs:
              - command: test
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
warnings_with_recommendations:
  react-native: []
`, sampleAppsExpoBareVersions...)
