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
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.RecreateUserSchemesVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - recreate-user-schemes@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - recreate-user-schemes@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
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
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - cocoapods-install@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - cocoapods-install@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
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
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
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
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
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
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleAppsCarthageVersions...)

var sampleAppClipResultYML = `options:
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
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
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
      format_version: "8"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@4:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@4: {}
          - cache-pull@2: {}
          - script@1:
              title: Do anything with Script step
          - certificate-and-profile-installer@1: {}
          - xcode-archive@3:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - export-xcarchive@3:
              inputs:
              - product: app-clip
          - deploy-to-bitrise-io@1: {}
          - cache-push@2: {}
    ios-app-clip-app-store-config: |
      format_version: "8"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@4:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@4: {}
          - cache-pull@2: {}
          - script@1:
              title: Do anything with Script step
          - certificate-and-profile-installer@1: {}
          - xcode-archive@3:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@1: {}
          - cache-push@2: {}
    ios-app-clip-development-config: |
      format_version: "8"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@4:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@4: {}
          - cache-pull@2: {}
          - script@1:
              title: Do anything with Script step
          - certificate-and-profile-installer@1: {}
          - xcode-archive@3:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - export-xcarchive@3:
              inputs:
              - product: app-clip
          - deploy-to-bitrise-io@1: {}
          - cache-push@2: {}
    ios-app-clip-enterprise-config: |
      format_version: "8"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@4:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@4: {}
          - cache-pull@2: {}
          - script@1:
              title: Do anything with Script step
          - certificate-and-profile-installer@1: {}
          - xcode-archive@3:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@1: {}
          - cache-push@2: {}
warnings:
  ios: []
warnings_with_recommendations:
  ios: []`
