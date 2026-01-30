package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

const workspaceSettingsWithAutocreateSchemesDisabledContent = `<?xml version="1.0" encoding="UTF-8"?>
 <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
 <plist version="1.0">
 <dict>
 	<key>IDEWorkspaceSharedSettings_AutocreateContextsIfNeeded</key>
 	<false/>
 </dict>
 </plist>
 `

 // TODO: broken test, needs investigation
 // failed to list Schemes in Project (/var/folders/f6/wf2hj3cj75qdwmt5rn814r_00000gn/T/TestIOSNoSchemes3008879754/001/BitriseXcode7Sample.xcodeproj): no schemes found and the Xcode project's 'Autocreate schemes' option is disabled
// func TestIOSNoSchemes(t *testing.T) {
// 	sampleAppDir := t.TempDir()

// 	helper.GitClone(t, sampleAppDir, "https://github.com/bitrise-samples/ios-no-shared-schemes.git")
// 	xcodeProjectPath := filepath.Join(sampleAppDir, "BitriseXcode7Sample.xcodeproj")
// 	projectEmbeddedWorksaceSettingsPth := filepath.Join(xcodeProjectPath, "project.xcworkspace/xcshareddata/WorkspaceSettings.xcsettings")
// 	require.NoError(t, os.MkdirAll(filepath.Dir(projectEmbeddedWorksaceSettingsPth), os.ModePerm))
// 	require.NoError(t, fileutil.WriteStringToFile(projectEmbeddedWorksaceSettingsPth, workspaceSettingsWithAutocreateSchemesDisabledContent))

// 	result, err := scanner.GenerateAndWriteResults(sampleAppDir, sampleAppDir, output.YAMLFormat)
// 	require.Error(t, err)

// 	iosWarnings := result.ScannerToWarningsWithRecommendations["ios"]
// 	require.Equal(t, 1, len(iosWarnings))
// 	require.True(t, strings.Contains(iosWarnings[0].Error, "no schemes found and the Xcode project's 'Autocreate schemes' option is disabled"))

// 	generalErrors := result.ScannerToErrorsWithRecommendations["general"]
// 	require.Equal(t, 1, len(generalErrors))
// 	require.Equal(t, "No known platform detected", generalErrors[0].Error)
// }

func TestIOS(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:             "ios-no-shared-schemes",
			RepoURL:          "https://github.com/bitrise-samples/ios-no-shared-schemes.git",
			ExpectedResult:   iosNoSharedSchemesResultYML,
			ExpectedVersions: iosNoSharedSchemesVersions,
		},
		{
			Name:             "ios-cocoapods-at-root",
			RepoURL:          "https://github.com/bitrise-samples/ios-cocoapods-at-root.git",
			ExpectedResult:   iosCocoapodsAtRootResultYML,
			ExpectedVersions: iosCocoapodsAtRootVersions,
		},
		{
			Name:             "sample-apps-ios-watchkit",
			RepoURL:          "https://github.com/bitrise-io/sample-apps-ios-watchkit.git",
			ExpectedResult:   sampleAppsIosWatchkitResultYML,
			ExpectedVersions: sampleAppsIosWatchkitVersions,
		},
		{
			Name:             "sample-apps-carthage",
			RepoURL:          "https://github.com/bitrise-samples/sample-apps-carthage.git",
			ExpectedResult:   sampleAppsCarthageResultYML,
			ExpectedVersions: sampleAppsCarthageVersions,
		},
		{
			Name:             "sample-apps-appclip",
			RepoURL:          "https://github.com/bitrise-io/sample-apps-ios-with-appclip.git",
			ExpectedResult:   sampleAppClipResultYML,
			ExpectedVersions: sampleAppClipVersions,
		},
		{
			Name:             "sample-apps-ios-swiftpm",
			RepoURL:          "https://github.com/bitrise-io/aci-xcode-spm-sample",
			ExpectedResult:   sampleSPMResultYML,
			ExpectedVersions: sampleSPMVersions,
		},
		{
			Name:             "samples-ios-swiftui-bitrise-todos",
			RepoURL:          "https://github.com/bitrise-io/samples-ios-swiftui-bitrise-todos",
			ExpectedResult:   appleMultiplatformAppResultYAML,
			ExpectedVersions: appleMultiplatformAppVersions,
		},
		{
			Name:             "sample-spm-project",
			RepoURL:          "https://github.com/bitrise-io/sample-spm-project.git",
			ExpectedResult:   sampleSPMProjectResultYML,
			ExpectedVersions: sampleSPMProjectVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var iosNoSharedSchemesVersions = []interface{}{
	models.FormatVersion,

	// ios-no-shared-schemes/archive_and_export_app
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-no-shared-schemes/build_for_testing
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	// ios-no-shared-schemes/run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-no-shared-schemes/test_without_building
	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,
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
                config: ios-test-config
              app-store:
                config: ios-test-config
              development:
                config: ios-test-config
              enterprise:
                config: ios-test-config
configs:
  ios:
    ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        archive_and_export_app:
          summary: Run your Xcode tests and create an IPA file to install your app on a
            device or share it with your team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests, export an IPA file
            from the project and save it.
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
        build_for_testing:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
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
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, iosNoSharedSchemesVersions...)

var iosCocoapodsAtRootVersions = []interface{}{
	models.FormatVersion,

	// ios-cocoapods-at-root/archive_and_export_app
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-cocoapods-at-root/build_for_testing
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCocoapodsVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeBuildForTestVersion,
	steps.CacheSaveCocoapodsVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	// ios-cocoapods-at-root/run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCocoapodsVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveCocoapodsVersion,
	steps.DeployToBitriseIoVersion,

	// ios-cocoapods-at-root/test_without_building
	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,
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
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        archive_and_export_app:
          summary: Run your Xcode tests and create an IPA file to install your app on a
            device or share it with your team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests, export an IPA file
            from the project and save it.
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
        build_for_testing:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-cocoapods-cache@%s: {}
          - cocoapods-install@%s:
              inputs:
              - is_cache_disabled: "true"
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - save-cocoapods-cache@%s: {}
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
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
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, iosCocoapodsAtRootVersions...)

var sampleAppsIosWatchkitVersions = []interface{}{
	// ios-config
	models.FormatVersion,

	// ios-app-clip-ad-hoc-config/ios-config/archive_and_export_app
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-ad-hoc-config/ios-config/build
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-test-config
	models.FormatVersion,

	// ios-app-clip-ad-hoc-config/ios-test-config/archive_and_export_app
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-ad-hoc-config/ios-test-config/build_for_testing
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-ad-hoc-config/ios-test-config/run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-app-clip-ad-hoc-config/ios-test-config/test_without_building
	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,
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
        archive_and_export_app:
          summary: Create an IPA file to install your app on a device or share it with your
            team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, export an IPA file from the project and
            save it.
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
        build:
          summary: Build your Xcode project.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any and build your project.
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
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        archive_and_export_app:
          summary: Run your Xcode tests and create an IPA file to install your app on a
            device or share it with your team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests, export an IPA file
            from the project and save it.
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
        build_for_testing:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
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
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleAppsIosWatchkitVersions...)

var sampleAppsCarthageVersions = []interface{}{
	models.FormatVersion,

	// ios-carthage-test-config/archive_and_export_app
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-carthage-test-config/build_for_testing
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCarthageVersion,
	steps.CarthageVersion,
	steps.XcodeBuildForTestVersion,
	steps.CacheSaveCarthageVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	// ios-carthage-test-config/run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreCarthageVersion,
	steps.CarthageVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveCarthageVersion,
	steps.DeployToBitriseIoVersion,

	// ios-carthage-test-config/test_without_building
	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,
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
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        archive_and_export_app:
          summary: Run your Xcode tests and create an IPA file to install your app on a
            device or share it with your team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests, export an IPA file
            from the project and save it.
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
        build_for_testing:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-carthage-cache@%s: {}
          - carthage@%s:
              inputs:
              - carthage_command: bootstrap
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - save-carthage-cache@%s: {}
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
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
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
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
        archive_and_export_app:
          summary: Create an IPA file to install your app on a device or share it with your
            team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, export an IPA file from the project and
            save it.
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
        build:
          summary: Build your Xcode project.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any and build your project.
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
        archive_and_export_app:
          summary: Create an IPA file to install your app on a device or share it with your
            team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, export an IPA file from the project and
            save it.
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
        build:
          summary: Build your Xcode project.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any and build your project.
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
        archive_and_export_app:
          summary: Create an IPA file to install your app on a device or share it with your
            team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, export an IPA file from the project and
            save it.
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
        build:
          summary: Build your Xcode project.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any and build your project.
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
        archive_and_export_app:
          summary: Create an IPA file to install your app on a device or share it with your
            team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, export an IPA file from the project and
            save it.
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
        build:
          summary: Build your Xcode project.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any and build your project.
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

var appleMultiplatformAppVersions = []interface{}{
	// iOS
	models.FormatVersion,

	// ios-test-config/archive_and_export_app
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-test-config/build_for_testing
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	// ios-test-config/run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-test-config/test_without_building
	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,

	// macOS
	models.FormatVersion,

	// ios-test-config/deploy
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeTestMacVersion,
	steps.XcodeArchiveMacVersion,
	steps.DeployToBitriseIoVersion,

	// ios-test-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestMacVersion,
	steps.DeployToBitriseIoVersion,
}

var appleMultiplatformAppResultYAML = fmt.Sprintf(`options:
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      Bitrise TODOs Sample.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          Bitrise TODOs Sample:
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
  macos:
    title: Project or Workspace path
    summary: The location of your Xcode project, Xcode workspace or SPM project files
      stored as an Environment Variable. In your Workflows, you can specify paths
      relative to this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      Bitrise TODOs Sample.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          Bitrise TODOs Sample:
            title: |-
              Application export method
              NOTE: `+"`none`"+` means: Export a copy of the application without re-signing.
            summary: The export method used to create an .app file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .app files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
            type: selector
            value_map:
              app-store:
                config: macos-test-config
              developer-id:
                config: macos-test-config
              development:
                config: macos-test-config
              none:
                config: macos-test-config
configs:
  ios:
    ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        archive_and_export_app:
          summary: Run your Xcode tests and create an IPA file to install your app on a
            device or share it with your team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests, export an IPA file
            from the project and save it.
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
        build_for_testing:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
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
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
  macos:
    macos-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: macos
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - xcode-archive-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
  macos: []
warnings_with_recommendations:
  ios: []
  macos: []
`, appleMultiplatformAppVersions...)

var sampleSPMVersions = []interface{}{
	// iOS
	models.FormatVersion,

	// ios-spm-test-config/archive_and_export_app
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios-spm-test-config/build_for_testing
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreSPMVersion,
	steps.XcodeBuildForTestVersion,
	steps.CacheSaveSPMVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	// ios-spm-test-config/run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreSPMVersion,
	steps.XcodeTestVersion,
	steps.CacheSaveSPMVersion,
	steps.DeployToBitriseIoVersion,

	// ios-spm-test-config/test_without_building
	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,
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
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        archive_and_export_app:
          summary: Run your Xcode tests and create an IPA file to install your app on a
            device or share it with your team.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests, export an IPA file
            from the project and save it.
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
        build_for_testing:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-spm-cache@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - save-spm-cache@%s: {}
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
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
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
warnings:
  ios: []
warnings_with_recommendations:
  ios: []
`, sampleSPMVersions...)

var sampleSPMProjectVersions = []interface{}{
	// iOS
	models.FormatVersion,

	// ios-spm-project-test-config/build_for_testing
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeBuildForTestVersion,
	steps.XcodeTestShardCalculationVersion,
	steps.DeployToBitriseIoVersion,

	// ios-spm-project-test-config/run_tests
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,

	// ios-spm-project-test-config/test_without_building
	steps.PullIntermediateFilesVersion,
	steps.XcodeTestWithoutBuildingVersion,

	// macOS
	models.FormatVersion,

	// macos-spm-project-test-config
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.XcodeTestMacVersion,
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
            config: ios-spm-project-test-config
  macos:
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
            config: macos-spm-project-test-config
configs:
  ios:
    ios-spm-project-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - TEST_SHARD_COUNT: 2
      pipelines:
        run_tests:
          workflows:
            build_for_testing: {}
            test_without_building:
              depends_on:
              - build_for_testing
              parallel: $TEST_SHARD_COUNT
      workflows:
        build_for_testing:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-build-for-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - destination: generic/platform=iOS Simulator
              - cache_level: none
          - xcode-test-shard-calculation@%s:
              inputs:
              - shard_count: $TEST_SHARD_COUNT
              - product_path: $BITRISE_XCTESTRUN_FILE_PATH
          - deploy-to-bitrise-io@%s:
              inputs:
              - pipeline_intermediate_files: |-
                  BITRISE_TEST_SHARDS_PATH
                  BITRISE_TEST_BUNDLE_PATH
        run_tests:
          summary: Run your Xcode tests and get the test report.
          description: The workflow will first clone your Git repository, cache and install
            your project's dependencies if any, run your Xcode tests and save the test results.
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
        test_without_building:
          steps:
          - pull-intermediate-files@%s: {}
          - xcode-test-without-building@%s:
              inputs:
              - only_testing: $BITRISE_TEST_SHARDS_PATH/$BITRISE_IO_PARALLEL_INDEX
              - xctestrun: $BITRISE_TEST_BUNDLE_PATH/all_tests.xctestrun
  macos:
    macos-spm-project-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: macos
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
warnings:
  ios: []
  macos: []
warnings_with_recommendations:
  ios: []
  macos: []
 `, sampleSPMProjectVersions...)
