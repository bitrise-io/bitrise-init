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

func TestIOS(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__ios__")
	require.NoError(t, err)

	type testCase struct {
		name           string
		repoURL        string
		expectedResult string
	}

	testCases := []testCase{
		{
			name:           "ios-no-shared-schemes",
			repoURL:        "https://github.com/bitrise-samples/ios-no-shared-schemes.git",
			expectedResult: iosNoSharedSchemesResultYML,
		},
		{
			name:           "ios-cocoapods-at-root",
			repoURL:        "https://github.com/bitrise-samples/ios-cocoapods-at-root.git",
			expectedResult: iosCocoapodsAtRootResultYML,
		},
		{
			name:           "sample-apps-ios-watchkit",
			repoURL:        "https://github.com/bitrise-io/sample-apps-ios-watchkit.git",
			expectedResult: sampleAppsIosWatchkitResultYML,
		},
		{
			name:           "sample-apps-carthage",
			repoURL:        "https://github.com/bitrise-samples/sample-apps-carthage.git",
			expectedResult: sampleAppsCarthageResultYML,
		},
		{
			name:           "sample-apps-appclip",
			repoURL:        "https://github.com/bitrise-io/sample-apps-ios-with-appclip.git",
			expectedResult: sampleAppClipResultYML,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		{
			sampleAppDir := filepath.Join(tmpDir, testCase.name)
			gitClone(t, sampleAppDir, testCase.repoURL)

			cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
			out, err := cmd.RunAndReturnTrimmedCombinedOutput()
			require.NoError(t, err, out)

			scanResultPth := filepath.Join(sampleAppDir, "result.yml")

			result, err := fileutil.ReadStringFromFile(scanResultPth)
			require.NoError(t, err)
			require.Equal(t, strings.TrimSpace(testCase.expectedResult), strings.TrimSpace(result))
		}
	}
}

var iosNoSharedSchemesVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,
}

var iosNoSharedSchemesResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      BitriseXcode7Sample.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          BitriseXcode7Sample:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-test-missing-shared-schemes-config
              app-store:
                config: ios-test-missing-shared-schemes-config
              development:
                config: ios-test-missing-shared-schemes-config
              enterprise:
                config: ios-test-missing-shared-schemes-config
configs:
  ios:
    ios-test-missing-shared-schemes-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - recreate-user-schemes@%s: {}
          - xcode-test@%s:
              inputs:
              - test_repetition_mode: retry_on_failure
          - xcode-archive@%s:
              inputs:
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
          - recreate-user-schemes@%s: {}
          - xcode-test@%s:
              inputs:
              - test_repetition_mode: retry_on_failure
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios:
  - error: |-
      No shared schemes found for project: BitriseXcode7Sample.xcodeproj.
      Automatically generated schemes may differ from the ones in your project.
      Make sure to <a href="http://devcenter.bitrise.io/ios/frequent-ios-issues/#xcode-scheme-not-found">share your schemes</a> for the expected behaviour.
    recommendations:
      DetailedError:
        title: We couldn’t parse your project files.
        description: |-
          You can fix the problem and try again, or skip auto-configuration and set up your project manually. Our auto-configurator returned the following error:
          No shared schemes found for project: BitriseXcode7Sample.xcodeproj.
          Automatically generated schemes may differ from the ones in your project.
          Make sure to <a href="http://devcenter.bitrise.io/ios/frequent-ios-issues/#xcode-scheme-not-found">share your schemes</a> for the expected behaviour.
`, iosNoSharedSchemesVersions...)

var iosCocoapodsAtRootVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,
}

var iosCocoapodsAtRootResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      iOSMinimalCocoaPodsSample.xcworkspace:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          iOSMinimalCocoaPodsSample:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-pod-test-config
              app-store:
                config: ios-pod-test-config
              development:
                config: ios-pod-test-config
              enterprise:
                config: ios-pod-test-config
configs:
  ios:
    ios-pod-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - cocoapods-install@%s: {}
          - xcode-test@%s:
              inputs:
              - test_repetition_mode: retry_on_failure
          - xcode-archive@%s:
              inputs:
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
          - cocoapods-install@%s: {}
          - xcode-test@%s:
              inputs:
              - test_repetition_mode: retry_on_failure
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, iosCocoapodsAtRootVersions...)

var sampleAppsIosWatchkitVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeBuildForTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

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

var sampleAppsIosWatchkitResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      watch-test.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          Complication - watch-test WatchKit App:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
          Glance - watch-test WatchKit App:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
          Notification - watch-test WatchKit App:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
          watch-test:
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
          watch-test WatchKit App:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-config
              app-store:
                config: ios-config
              development:
                config: ios-config
              enterprise:
                config: ios-config
configs:
  ios:
    ios-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - xcode-archive@%s:
              inputs:
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
    ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
              - test_repetition_mode: retry_on_failure
          - xcode-archive@%s:
              inputs:
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
              - test_repetition_mode: retry_on_failure
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleAppsIosWatchkitVersions...)

var sampleAppsCarthageVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsCarthageResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      sample-apps-carthage.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          sample-apps-carthage:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-carthage-test-config
              app-store:
                config: ios-carthage-test-config
              development:
                config: ios-carthage-test-config
              enterprise:
                config: ios-carthage-test-config
configs:
  ios:
    ios-carthage-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-test@%s:
              inputs:
              - test_repetition_mode: retry_on_failure
          - xcode-archive@%s:
              inputs:
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
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-test@%s:
              inputs:
              - test_repetition_mode: retry_on_failure
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleAppsCarthageVersions...)

var sampleAppClipResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      Sample.xcworkspace:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          SampleAppClipApp:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-app-clip-ad-hoc-config
              app-store:
                config: ios-app-clip-app-store-config
              development:
                config: ios-app-clip-development-config
              enterprise:
                config: ios-app-clip-enterprise-config
configs:
  ios:
    ios-app-clip-ad-hoc-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - xcode-archive@%s:
              inputs:
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - export-xcarchive@%s:
              inputs:
              - product: app-clip
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
    ios-app-clip-app-store-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - xcode-archive@%s:
              inputs:
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
    ios-app-clip-development-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - xcode-archive@%s:
              inputs:
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - export-xcarchive@%s:
              inputs:
              - product: app-clip
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
    ios-app-clip-enterprise-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
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
          - xcode-archive@%s:
              inputs:
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
          - cache-push@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []`,
	// ios-app-clip-ad-hoc-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeArchiveVersion,
	steps.ExportXCArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-ad-hoc-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeBuildForTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-app-store-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-app-store-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeBuildForTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-development-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeArchiveVersion,
	steps.ExportXCArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-development-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeBuildForTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-enterprise-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeArchiveVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-enterprise-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.XcodeBuildForTestVersion,
	steps.CachePushVersion,
	steps.DeployToBitriseIoVersion,
)
