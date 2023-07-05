package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestIOS(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"ios-no-shared-schemes",
			"https://github.com/bitrise-samples/ios-no-shared-schemes.git",
			"",
			iosNoSharedSchemesResultYML,
			iosNoSharedSchemesVersions,
		},
		{
			"ios-cocoapods-at-root",
			"https://github.com/bitrise-samples/ios-cocoapods-at-root.git",
			"",
			iosCocoapodsAtRootResultYML,
			iosCocoapodsAtRootVersions,
		},
		{
			"sample-apps-ios-watchkit",
			"https://github.com/bitrise-io/sample-apps-ios-watchkit.git",
			"",
			sampleAppsIosWatchkitResultYML,
			sampleAppsIosWatchkitVersions,
		},
		{
			"sample-apps-carthage",
			"https://github.com/bitrise-samples/sample-apps-carthage.git",
			"",
			sampleAppsCarthageResultYML,
			sampleAppsCarthageVersions,
		},
		{
			"sample-apps-appclip",
			"https://github.com/bitrise-io/sample-apps-ios-with-appclip.git",
			"",
			sampleAppClipResultYML,
			sampleAppClipVersions,
		},
		{
			"sample-apps-ios-swiftpm",
			"https://github.com/bitrise-io/aci-xcode-spm-sample",
			"",
			sampleSPMResultYML,
			sampleSPMVersions,
		},
		{
			"sample-spm-project",
			"https://github.com/bitrise-io/sample-spm-project.git",
			"",
			sampleSPMProjectResultYML,
			sampleSPMProjectVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var iosNoSharedSchemesVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
}

var iosNoSharedSchemesResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
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
          - recreate-user-schemes@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - recreate-user-schemes@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios:
  - error: |-
      No shared schemes found for project: BitriseXcode7Sample.xcodeproj.
      Automatically generated schemes may differ from the ones in your project.
      Make sure to <a href="https://support.bitrise.io/hc/en-us/articles/4405779956625">share your schemes</a> for the expected behaviour.
    recommendations:
      DetailedError:
        title: We couldn’t parse your project files.
        description: |-
          You can fix the problem and try again, or skip auto-configuration and set up your project manually. Our auto-configurator returned the following error:
          No shared schemes found for project: BitriseXcode7Sample.xcodeproj.
          Automatically generated schemes may differ from the ones in your project.
          Make sure to <a href="https://support.bitrise.io/hc/en-us/articles/4405779956625">share your schemes</a> for the expected behaviour.
`, iosNoSharedSchemesVersions...)

var iosCocoapodsAtRootVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCocoapodsVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveCocoapodsVersion,
	steps.DeployToBitriseIoVersion,
}

var iosCocoapodsAtRootResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
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
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-cocoapods-cache@%s: {}
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - save-cocoapods-cache@%s: {}
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
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.DeployToBitriseIoVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsIosWatchkitResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
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
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
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
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleAppsIosWatchkitVersions...)

var sampleAppsCarthageVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCarthageVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveCarthageVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsCarthageResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
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
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-carthage-cache@%s: {}
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - save-carthage-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleAppsCarthageVersions...)

var sampleAppClipVersions = []interface{}{
	// ios-app-clip-ad-hoc-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeArchiveVersion,
	steps.ExportXCArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-ad-hoc-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-app-store-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-app-store-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-development-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeArchiveVersion,
	steps.ExportXCArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-development-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-enterprise-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-enterprise-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppClipResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - export-xcarchive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - product: app-clip
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
    ios-app-clip-app-store-config: |
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
    ios-app-clip-development-config: |
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - export-xcarchive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - product: app-clip
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
    ios-app-clip-enterprise-config: |
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow only builds the project because the project scanner could not find any tests.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: platform=iOS Simulator,name=iPhone 8 Plus,OS=latest
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleAppClipVersions...)

var sampleSPMVersions = []interface{}{
	models.FormatVersion,

	// ios-spm-test-config/deploy
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-spm-test-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreSPMVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveSPMVersion,
	steps.DeployToBitriseIoVersion,
}
var sampleSPMResultYML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      aci-xcode-spm-sample.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          aci-xcode-spm-sample:
            title: Distribution method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_DISTRIBUTION_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: ios-spm-test-config
              app-store:
                config: ios-spm-test-config
              development:
                config: ios-spm-test-config
              enterprise:
                config: ios-spm-test-config
configs:
  ios:
    ios-spm-test-config: |
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
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - automatic_code_signing: api-key
              - cache_level: none
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-spm-cache@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - save-spm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleSPMVersions...)

var sampleSPMProjectVersions = []interface{}{
	models.FormatVersion,

	// ios-spm-spm-project-test-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreSPMVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveSPMVersion,
	steps.DeployToBitriseIoVersion,
}
var sampleSPMProjectResultYML = fmt.Sprintf(`options:
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
            config: ios-spm-spm-project-test-config
configs:
  ios:
    ios-spm-spm-project-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      workflows:
        primary:
          description: |
            The workflow executes the tests. The *retry_on_failure* test repetition mode is enabled.

            Next steps:
            - Check out [Getting started with iOS apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-ios-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-spm-cache@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - test_repetition_mode: retry_on_failure
              - cache_level: none
          - save-spm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleSPMProjectVersions...)
