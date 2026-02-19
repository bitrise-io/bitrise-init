package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestFlutter(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:             "sample-apps-flutter-ios-android",
			RepoURL:          "https://github.com/bitrise-samples/sample-apps-flutter-ios-android.git",
			ExpectedResult:   flutterSampleAppResultYML,
			ExpectedVersions: flutterSampleAppVersions,
		},
		{
			Name:             "sample-apps-flutter-ios-android-package",
			RepoURL:          "https://github.com/bitrise-samples/sample-apps-flutter-ios-android-package.git",
			ExpectedResult:   flutterSamplePackageResultYML,
			ExpectedVersions: flutterSamplePackageVersions,
		},
		{
			Name:             "sample-apps-flutter-ios-android-plugin",
			RepoURL:          "https://github.com/bitrise-samples/sample-apps-flutter-ios-android-plugin.git",
			ExpectedResult:   flutterSamplePluginResultYML,
			ExpectedVersions: flutterSamplePluginVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var flutterSampleAppVersions = []interface{}{
	// flutter-config-test-app-both-0
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterSampleAppResultYML = fmt.Sprintf(`options:
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: selector
    value_map:
      .:
        config: flutter-config-test-both-0
configs:
  flutter:
    flutter-config-test-both-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        deploy:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html) for signing and deployment options.
            - Check out the [Code signing guide](https://devcenter.bitrise.io/en/code-signing.html) for iOS and Android
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.7.12
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: both
              - ios_output_type: archive
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Builds project and runs tests.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.7.12
          - restore-dart-cache@%s: {}
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  flutter: []
warnings_with_recommendations:
  flutter: []
`, flutterSampleAppVersions...)

var flutterSamplePackageVersions = []interface{}{
	// flutter-config-test-0
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterSamplePackageResultYML = fmt.Sprintf(`options:
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: selector
    value_map:
      .:
        config: flutter-config-test-0
configs:
  flutter:
    flutter-config-test-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        primary:
          description: |
            Builds project and runs tests.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.7.12
          - restore-dart-cache@%s: {}
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  flutter: []
warnings_with_recommendations:
  flutter: []
`, flutterSamplePackageVersions...)

var flutterSamplePluginVersions = []interface{}{
	// flutter-config-notest-android-0
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterAnalyzeVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,

	// flutter-config-test-both-1
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterSamplePluginResultYML = fmt.Sprintf(`options:
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: selector
    value_map:
      .:
        config: flutter-config-notest-android-0
      example:
        config: flutter-config-test-both-1
configs:
  flutter:
    flutter-config-notest-android-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        deploy:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html) for signing and deployment options.
            - Check out the [Code signing guide](https://devcenter.bitrise.io/en/code-signing.html) for iOS and Android
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.7.12
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: android
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Builds project and runs tests.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.7.12
          - restore-dart-cache@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
    flutter-config-test-both-1: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        deploy:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html) for signing and deployment options.
            - Check out the [Code signing guide](https://devcenter.bitrise.io/en/code-signing.html) for iOS and Android
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.7.12
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: both
              - ios_output_type: archive
          - deploy-to-bitrise-io@%s: {}
        primary:
          description: |
            Builds project and runs tests.

            Next steps:
            - Check out [Getting started with Flutter apps](https://devcenter.bitrise.io/en/getting-started/getting-started-with-flutter-apps.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.7.12
          - restore-dart-cache@%s: {}
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  flutter: []
warnings_with_recommendations:
  flutter: []
`, flutterSamplePluginVersions...)
