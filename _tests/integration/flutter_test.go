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

func TestFlutter(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__flutter__")
	require.NoError(t, err)

	t.Log("sample-apps-flutter-ios-android")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-flutter-ios-android")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-flutter-ios-android.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "sample-apps-flutter-ios-android", strings.TrimSpace(flutterSampleAppResultYML), strings.TrimSpace(result), flutterSampleAppVersions...)
	}

	t.Log("sample-apps-flutter-ios-android-package")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-flutter-ios-android-package")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-flutter-ios-android-package.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "sample-apps-flutter-ios-android-package", strings.TrimSpace(flutterSamplePackageResultYML), strings.TrimSpace(result), flutterSamplePackageVersions...)
	}

	t.Log("sample-apps-flutter-ios-android-plugin")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-flutter-ios-android-plugin")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-flutter-ios-android-plugin.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "sample-apps-flutter-ios-android-plugin", strings.TrimSpace(flutterSamplePluginResultYML), strings.TrimSpace(result), flutterSamplePluginVersions...)
	}
}

var flutterSampleAppVersions = []interface{}{
	// flutter-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-android
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-both
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-android
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-both
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
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
        title: Run tests found in the project
        summary: Our Flutter Test Step can run the tests found in your project's repository.
        type: selector
        value_map:
          "no":
            title: Project or Workspace path
            summary: The location of your Xcode project or Xcode workspace files,
              stored as an Environment Variable. In your Workflows, you can specify
              paths relative to this path.
            env_key: BITRISE_PROJECT_PATH
            type: selector
            value_map:
              ios/Runner.xcworkspace:
                title: Scheme name
                summary: An Xcode scheme defines a collection of targets to build,
                  a configuration to use when building, and a collection of tests
                  to execute. Only shared schemes are detected automatically but you
                  can use any scheme as a target on Bitrise. You can change the scheme
                  at any time in your Env Vars.
                env_key: BITRISE_SCHEME
                type: selector
                value_map:
                  Runner:
                    title: Distribution method
                    summary: The export method used to create an .ipa file in your
                      builds, stored as an Environment Variable. You can change this
                      at any time, or even create several .ipa files with different
                      export methods in the same build.
                    env_key: BITRISE_DISTRIBUTION_METHOD
                    type: selector
                    value_map:
                      ad-hoc:
                        config: flutter-config-app-both
                      app-store:
                        config: flutter-config-app-both
                      development:
                        config: flutter-config-app-both
                      enterprise:
                        config: flutter-config-app-both
          "yes":
            title: Project or Workspace path
            summary: The location of your Xcode project or Xcode workspace files,
              stored as an Environment Variable. In your Workflows, you can specify
              paths relative to this path.
            env_key: BITRISE_PROJECT_PATH
            type: selector
            value_map:
              ios/Runner.xcworkspace:
                title: Scheme name
                summary: An Xcode scheme defines a collection of targets to build,
                  a configuration to use when building, and a collection of tests
                  to execute. Only shared schemes are detected automatically but you
                  can use any scheme as a target on Bitrise. You can change the scheme
                  at any time in your Env Vars.
                env_key: BITRISE_SCHEME
                type: selector
                value_map:
                  Runner:
                    title: Distribution method
                    summary: The export method used to create an .ipa file in your
                      builds, stored as an Environment Variable. You can change this
                      at any time, or even create several .ipa files with different
                      export methods in the same build.
                    env_key: BITRISE_DISTRIBUTION_METHOD
                    type: selector
                    value_map:
                      ad-hoc:
                        config: flutter-config-test-app-both
                      app-store:
                        config: flutter-config-test-app-both
                      development:
                        config: flutter-config-test-app-both
                      enterprise:
                        config: flutter-config-test-app-both
configs:
  flutter:
    flutter-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-android: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: android
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-both: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: both
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: ios
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-android: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
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
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-both: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: ios
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  flutter: []
warnings_with_recommendations:
  flutter: []
`, flutterSampleAppVersions...)

var flutterSamplePackageVersions = []interface{}{
	// flutter-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-android
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-both
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-android
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-both
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
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
        title: Run tests found in the project
        summary: Our Flutter Test Step can run the tests found in your project's repository.
        type: selector
        value_map:
          "no":
            config: flutter-config
          "yes":
            config: flutter-config-test
configs:
  flutter:
    flutter-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-android: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: android
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-both: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: both
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: ios
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-android: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
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
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-both: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: ios
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  flutter: []
warnings_with_recommendations:
  flutter: []
`, flutterSamplePackageVersions...)

var flutterSamplePluginVersions = []interface{}{
	// flutter-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-android
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-both
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-app-ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-android
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-both
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// flutter-config-test-app-ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.FlutterBuildVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.FlutterTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
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
        config: flutter-config-app-android
      example:
        title: Run tests found in the project
        summary: Our Flutter Test Step can run the tests found in your project's repository.
        type: selector
        value_map:
          "no":
            title: Project or Workspace path
            summary: The location of your Xcode project or Xcode workspace files,
              stored as an Environment Variable. In your Workflows, you can specify
              paths relative to this path.
            env_key: BITRISE_PROJECT_PATH
            type: selector
            value_map:
              example/ios/Runner.xcworkspace:
                title: Scheme name
                summary: An Xcode scheme defines a collection of targets to build,
                  a configuration to use when building, and a collection of tests
                  to execute. Only shared schemes are detected automatically but you
                  can use any scheme as a target on Bitrise. You can change the scheme
                  at any time in your Env Vars.
                env_key: BITRISE_SCHEME
                type: selector
                value_map:
                  Runner:
                    title: Distribution method
                    summary: The export method used to create an .ipa file in your
                      builds, stored as an Environment Variable. You can change this
                      at any time, or even create several .ipa files with different
                      export methods in the same build.
                    env_key: BITRISE_DISTRIBUTION_METHOD
                    type: selector
                    value_map:
                      ad-hoc:
                        config: flutter-config-app-both
                      app-store:
                        config: flutter-config-app-both
                      development:
                        config: flutter-config-app-both
                      enterprise:
                        config: flutter-config-app-both
          "yes":
            title: Project or Workspace path
            summary: The location of your Xcode project or Xcode workspace files,
              stored as an Environment Variable. In your Workflows, you can specify
              paths relative to this path.
            env_key: BITRISE_PROJECT_PATH
            type: selector
            value_map:
              example/ios/Runner.xcworkspace:
                title: Scheme name
                summary: An Xcode scheme defines a collection of targets to build,
                  a configuration to use when building, and a collection of tests
                  to execute. Only shared schemes are detected automatically but you
                  can use any scheme as a target on Bitrise. You can change the scheme
                  at any time in your Env Vars.
                env_key: BITRISE_SCHEME
                type: selector
                value_map:
                  Runner:
                    title: Distribution method
                    summary: The export method used to create an .ipa file in your
                      builds, stored as an Environment Variable. You can change this
                      at any time, or even create several .ipa files with different
                      export methods in the same build.
                    env_key: BITRISE_DISTRIBUTION_METHOD
                    type: selector
                    value_map:
                      ad-hoc:
                        config: flutter-config-test-app-both
                      app-store:
                        config: flutter-config-test-app-both
                      development:
                        config: flutter-config-test-app-both
                      enterprise:
                        config: flutter-config-test-app-both
configs:
  flutter:
    flutter-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-android: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: android
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-both: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: both
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-app-ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: ios
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-android: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
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
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-both: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
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
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
    flutter-config-test-app-ios: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: flutter
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
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-build@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
              - platform: ios
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - distribution_method: $BITRISE_DISTRIBUTION_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - flutter-installer@%s:
              inputs:
              - is_update: "false"
          - cache-pull@%s: {}
          - flutter-analyze@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - flutter-test@%s:
              inputs:
              - project_location: $BITRISE_FLUTTER_PROJECT_LOCATION
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  flutter: []
warnings_with_recommendations:
  flutter: []
`, flutterSamplePluginVersions...)
