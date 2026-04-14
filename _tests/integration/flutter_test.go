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
			Name:              "flutter-ios-android",
			RepoURL:           "https://github.com/bitrise-io/flutter-samples.git",
			RelativeSearchDir: "flutter-ios-android",
			Branch:            "main",
			ExpectedResult:    flutterIosAndroidResultYML,
			ExpectedVersions:  flutterIosAndroidVersions,
		},
		{
			Name:              "flutter-package",
			RepoURL:           "https://github.com/bitrise-io/flutter-samples.git",
			RelativeSearchDir: "flutter-package",
			Branch:            "main",
			ExpectedResult:    flutterPackageResultYML,
			ExpectedVersions:  flutterPackageVersions,
		},
		{
			Name:              "flutter-plugin",
			RepoURL:           "https://github.com/bitrise-io/flutter-samples.git",
			RelativeSearchDir: "flutter-plugin",
			Branch:            "main",
			ExpectedResult:    flutterPluginResultYML,
			ExpectedVersions:  flutterPluginVersions,
		},
		{
			Name:              "flutter-web",
			RepoURL:           "https://github.com/bitrise-io/flutter-samples.git",
			RelativeSearchDir: "flutter-web",
			Branch:            "main",
			ExpectedResult:    flutterWebResultYML,
			ExpectedVersions:  flutterWebVersions,
		},
		{
			Name:              "flutter-ios-android-web",
			RepoURL:           "https://github.com/bitrise-io/flutter-samples.git",
			RelativeSearchDir: "flutter-ios-android-web",
			Branch:            "main",
			ExpectedResult:    flutterIosAndroidWebResultYML,
			ExpectedVersions:  flutterIosAndroidWebVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var flutterIosAndroidVersions = []interface{}{
	// flutter-config-test-ios-android-0
	models.FormatVersion,
	// build_app workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	// run_tests workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterIosAndroidResultYML = fmt.Sprintf(`options:
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: selector
    value_map:
      .:
        config: flutter-config-test-ios-android-0
configs:
  flutter:
    flutter-config-test-ios-android-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        build_app:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html) for signing and deployment options.
            - Check out the Code signing guide for [iOS](https://docs.bitrise.io/en/bitrise-ci/code-signing/ios-code-signing.html) and [Android](https://docs.bitrise.io/en/bitrise-ci/code-signing/android-code-signing.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
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
        run_tests:
          description: |
            Runs tests or analysis.

            Runs flutter-test if a test directory is present, otherwise runs flutter-analyze.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
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
`, flutterIosAndroidVersions...)

var flutterPackageVersions = []interface{}{
	// flutter-config-test-0
	models.FormatVersion,
	// run_tests workflow only
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterPackageResultYML = fmt.Sprintf(`options:
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
        run_tests:
          description: |
            Runs tests or analysis.

            Runs flutter-test if a test directory is present, otherwise runs flutter-analyze.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
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
`, flutterPackageVersions...)

var flutterPluginVersions = []interface{}{
	// flutter-config-test-android-0 (root plugin: Android only, has tests)
	models.FormatVersion,
	// build_app workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	// run_tests workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,

	// flutter-config-test-ios-android-1 (example app: iOS + Android, has tests)
	models.FormatVersion,
	// build_app workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	// run_tests workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterPluginResultYML = fmt.Sprintf(`options:
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: selector
    value_map:
      .:
        config: flutter-config-test-android-0
      example:
        config: flutter-config-test-ios-android-1
configs:
  flutter:
    flutter-config-test-android-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        build_app:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html) for signing and deployment options.
            - Check out the Code signing guide for [iOS](https://docs.bitrise.io/en/bitrise-ci/code-signing/ios-code-signing.html) and [Android](https://docs.bitrise.io/en/bitrise-ci/code-signing/android-code-signing.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: android
          - deploy-to-bitrise-io@%s: {}
        run_tests:
          description: |
            Runs tests or analysis.

            Runs flutter-test if a test directory is present, otherwise runs flutter-analyze.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
          - restore-dart-cache@%s: {}
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - save-dart-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
    flutter-config-test-ios-android-1: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        build_app:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html) for signing and deployment options.
            - Check out the Code signing guide for [iOS](https://docs.bitrise.io/en/bitrise-ci/code-signing/ios-code-signing.html) and [Android](https://docs.bitrise.io/en/bitrise-ci/code-signing/android-code-signing.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
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
        run_tests:
          description: |
            Runs tests or analysis.

            Runs flutter-test if a test directory is present, otherwise runs flutter-analyze.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
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
`, flutterPluginVersions...)

var flutterWebVersions = []interface{}{
	// flutter-config-test-web-0
	models.FormatVersion,
	// run_tests workflow only
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterWebResultYML = fmt.Sprintf(`options:
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: selector
    value_map:
      .:
        config: flutter-config-test-web-0
configs:
  flutter:
    flutter-config-test-web-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        run_tests:
          description: |
            Runs tests or analysis.

            Runs flutter-test if a test directory is present, otherwise runs flutter-analyze.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.41.6
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
`, flutterWebVersions...)

var flutterIosAndroidWebVersions = []interface{}{
	// flutter-config-test-ios-android-0
	models.FormatVersion,
	// build_app workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	// run_tests workflow
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.FlutterInstallVersion,
	steps.CacheRestoreDartVersion,
	steps.FlutterTestVersion,
	steps.CacheSaveDartVersion,
	steps.DeployToBitriseIoVersion,
}

var flutterIosAndroidWebResultYML = fmt.Sprintf(`options:
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: selector
    value_map:
      .:
        config: flutter-config-test-ios-android-web-0
configs:
  flutter:
    flutter-config-test-ios-android-web-0: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      workflows:
        build_app:
          description: |
            Builds and deploys app using [Deploy to bitrise.io Step](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html#deploying-a-flutter-app).

            If you build for iOS, make sure to set up code signing secrets on Bitrise for a successful build.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html) for signing and deployment options.
            - Check out the Code signing guide for [iOS](https://docs.bitrise.io/en/bitrise-ci/code-signing/ios-code-signing.html) and [Android](https://docs.bitrise.io/en/bitrise-ci/code-signing/android-code-signing.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.29.2
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
        run_tests:
          description: |
            Runs tests or analysis.

            Runs flutter-test if a test directory is present, otherwise runs flutter-analyze.

            Next steps:
            - Check out [Getting started with Flutter apps](https://docs.bitrise.io/en/bitrise-ci/getting-started/quick-start-guides/getting-started-with-flutter-projects.html).
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - flutter-installer@%s:
              inputs:
              - version: 3.29.2
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
`, flutterIosAndroidVersions...)
