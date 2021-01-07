package integration

import (
	"fmt"
	"os"
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

func TestManualConfig(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__manual-config__")
	require.NoError(t, err)

	t.Log("manual-config")
	{
		manualConfigDir := filepath.Join(tmpDir, "manual-config")
		require.NoError(t, os.MkdirAll(manualConfigDir, 0777))
		fmt.Printf("manualConfigDir: %s\n", manualConfigDir)

		cmd := command.New(binPath(), "--ci", "manual-config", "--output-dir", manualConfigDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(manualConfigDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "manual-config", strings.TrimSpace(customConfigResultYML), strings.TrimSpace(result), customConfigVersions...)
	}
}

var customConfigVersions = []interface{}{
	// android
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.ChangeAndroidVersionCodeAndVersionNameVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.AndroidBuildVersion,
	steps.SignAPKVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidLintVersion,
	steps.AndroidUnitTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// cordova
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NpmVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.CordovaArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// fastlane
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FastlaneVersion,
	steps.DeployToBitriseIoVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.FastlaneVersion,
	steps.DeployToBitriseIoVersion,

	// flutter
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.FlutterInstallVersion,
	steps.CachePullVersion,
	steps.FlutterAnalyzeVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

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

	// ionic
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NpmVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.IonicArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// ios
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.RecreateUserSchemesVersion,
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
	steps.RecreateUserSchemesVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// macos
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.RecreateUserSchemesVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestMacVersion,
	steps.XcodeArchiveMacVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CachePullVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.RecreateUserSchemesVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeTestMacVersion,
	steps.DeployToBitriseIoVersion,
	steps.CachePushVersion,

	// other
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.DeployToBitriseIoVersion,

	// react native
	// default-react-native-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// default-react-native-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.DeployToBitriseIoVersion,

	// default-react-native-expo-config/deploy
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.YarnVersion,
	steps.ScriptVersion,
	steps.ExpoDetachVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	// default-react-native-expo-config/primary
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.DeployToBitriseIoVersion,

	// xamarin
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XamarinUserManagementVersion,
	steps.NugetRestoreVersion,
	steps.XamarinComponentsRestoreVersion,
	steps.XamarinArchiveVersion,
	steps.DeployToBitriseIoVersion,
}

var customConfigResultYML = fmt.Sprintf(`options:
  android:
    title: The root directory of an Android project
    summary: The root directory of your Android project, stored as an Environment
      Variable. In your Workflows, you can specify paths relative to this path. You
      can change this at any time.
    env_key: PROJECT_LOCATION
    type: user_input
    value_map:
      "":
        title: Module
        summary: Modules provide a container for your Android project's source code,
          resource files, and app level settings, such as the module-level build file
          and Android manifest file. Each module can be independently built, tested,
          and debugged. You can add new modules to your Bitrise builds at any time.
        env_key: MODULE
        type: user_input
        value_map:
          "":
            title: Variant
            summary: Your Android build variant. You can add variants at any time,
              as well as further configure your existing variants later.
            env_key: VARIANT
            type: user_input_optional
            value_map:
              "":
                config: default-android-config
  cordova:
    title: Directory of the Cordova config.xml file
    summary: The working directory of your Cordova project is where you store your
      config.xml file. In your Workflows, you can specify paths relative to this path.
      You can change this at any time.
    env_key: CORDOVA_WORK_DIR
    type: user_input
    value_map:
      "":
        title: The platform to use in cordova-cli commands
        summary: The target platform for your build, stored as an Environment Variable.
          Your options are iOS, Android, or both. You can change this in your Env
          Vars at any time.
        env_key: CORDOVA_PLATFORM
        type: selector
        value_map:
          android:
            config: default-cordova-config
          ios:
            config: default-cordova-config
          ios,android:
            config: default-cordova-config
  fastlane:
    title: Working directory
    summary: The directory where your Fastfile is located.
    env_key: FASTLANE_WORK_DIR
    type: user_input
    value_map:
      "":
        title: Fastlane lane
        summary: The lane that will be used in your builds, stored as an Environment
          Variable. You can change this at any time.
        env_key: FASTLANE_LANE
        type: user_input
        value_map:
          "":
            title: Project type
            summary: The project type of the app you added to Bitrise.
            type: selector
            value_map:
              android:
                config: default-fastlane-android-config
              ios:
                config: default-fastlane-ios-config
  flutter:
    title: Project location
    summary: The path to your Flutter project, stored as an Environment Variable.
      In your Workflows, you can specify paths relative to this path. You can change
      this at any time.
    env_key: BITRISE_FLUTTER_PROJECT_LOCATION
    type: user_input
    value_map:
      "":
        title: Run tests found in the project
        summary: Our Flutter Test Step can run the tests found in your project's repository.
        type: selector
        value_map:
          "no":
            title: Platform
            summary: The target platform for your first build. Your options are iOS,
              Android, both, or neither. You can change this in your Env Vars at any
              time.
            type: selector
            value_map:
              android:
                config: flutter-config-app-android
              both:
                title: Project or Workspace path
                summary: The location of your Xcode project or Xcode workspace files,
                  stored as an Environment Variable. In your Workflows, you can specify
                  paths relative to this path.
                env_key: BITRISE_PROJECT_PATH
                type: user_input
                value_map:
                  "":
                    title: Scheme name
                    summary: An Xcode scheme defines a collection of targets to build,
                      a configuration to use when building, and a collection of tests
                      to execute. Only shared schemes are detected automatically but
                      you can use any scheme as a target on Bitrise. You can change
                      the scheme at any time in your Env Vars.
                    env_key: BITRISE_SCHEME
                    type: user_input
                    value_map:
                      "":
                        title: ipa export method
                        summary: The export method used to create an .ipa file in
                          your builds, stored as an Environment Variable. You can
                          change this at any time, or even create several .ipa files
                          with different export methods in the same build.
                        env_key: BITRISE_EXPORT_METHOD
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
              ios:
                title: Project or Workspace path
                summary: The location of your Xcode project or Xcode workspace files,
                  stored as an Environment Variable. In your Workflows, you can specify
                  paths relative to this path.
                env_key: BITRISE_PROJECT_PATH
                type: user_input
                value_map:
                  "":
                    title: Scheme name
                    summary: An Xcode scheme defines a collection of targets to build,
                      a configuration to use when building, and a collection of tests
                      to execute. Only shared schemes are detected automatically but
                      you can use any scheme as a target on Bitrise. You can change
                      the scheme at any time in your Env Vars.
                    env_key: BITRISE_SCHEME
                    type: user_input
                    value_map:
                      "":
                        title: ipa export method
                        summary: The export method used to create an .ipa file in
                          your builds, stored as an Environment Variable. You can
                          change this at any time, or even create several .ipa files
                          with different export methods in the same build.
                        env_key: BITRISE_EXPORT_METHOD
                        type: selector
                        value_map:
                          ad-hoc:
                            config: flutter-config-app-ios
                          app-store:
                            config: flutter-config-app-ios
                          development:
                            config: flutter-config-app-ios
                          enterprise:
                            config: flutter-config-app-ios
              none:
                config: flutter-config
          "yes":
            title: Platform
            summary: The target platform for your first build. Your options are iOS,
              Android, both, or neither. You can change this in your Env Vars at any
              time.
            type: selector
            value_map:
              android:
                config: flutter-config-test-app-android
              both:
                title: Project or Workspace path
                summary: The location of your Xcode project or Xcode workspace files,
                  stored as an Environment Variable. In your Workflows, you can specify
                  paths relative to this path.
                env_key: BITRISE_PROJECT_PATH
                type: user_input
                value_map:
                  "":
                    title: Scheme name
                    summary: An Xcode scheme defines a collection of targets to build,
                      a configuration to use when building, and a collection of tests
                      to execute. Only shared schemes are detected automatically but
                      you can use any scheme as a target on Bitrise. You can change
                      the scheme at any time in your Env Vars.
                    env_key: BITRISE_SCHEME
                    type: user_input
                    value_map:
                      "":
                        title: ipa export method
                        summary: The export method used to create an .ipa file in
                          your builds, stored as an Environment Variable. You can
                          change this at any time, or even create several .ipa files
                          with different export methods in the same build.
                        env_key: BITRISE_EXPORT_METHOD
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
              ios:
                title: Project or Workspace path
                summary: The location of your Xcode project or Xcode workspace files,
                  stored as an Environment Variable. In your Workflows, you can specify
                  paths relative to this path.
                env_key: BITRISE_PROJECT_PATH
                type: user_input
                value_map:
                  "":
                    title: Scheme name
                    summary: An Xcode scheme defines a collection of targets to build,
                      a configuration to use when building, and a collection of tests
                      to execute. Only shared schemes are detected automatically but
                      you can use any scheme as a target on Bitrise. You can change
                      the scheme at any time in your Env Vars.
                    env_key: BITRISE_SCHEME
                    type: user_input
                    value_map:
                      "":
                        title: ipa export method
                        summary: The export method used to create an .ipa file in
                          your builds, stored as an Environment Variable. You can
                          change this at any time, or even create several .ipa files
                          with different export methods in the same build.
                        env_key: BITRISE_EXPORT_METHOD
                        type: selector
                        value_map:
                          ad-hoc:
                            config: flutter-config-test-app-ios
                          app-store:
                            config: flutter-config-test-app-ios
                          development:
                            config: flutter-config-test-app-ios
                          enterprise:
                            config: flutter-config-test-app-ios
              none:
                config: flutter-config-test
  ionic:
    title: Directory of the Ionic config.xml file
    summary: The working directory of your Ionic project is where you store your config.xml
      file. This location is stored as an Environment Variable. In your Workflows,
      you can specify paths relative to this path. You can change this at any time.
    env_key: IONIC_WORK_DIR
    type: user_input
    value_map:
      "":
        title: The platform to use in ionic-cli commands
        summary: The target platform for your builds, stored as an Environment Variable.
          Your options are iOS, Android, or both. You can change this in your Env
          Vars at any time.
        env_key: IONIC_PLATFORM
        type: selector
        value_map:
          android:
            config: default-ionic-config
          ios:
            config: default-ionic-config
          ios,android:
            config: default-ionic-config
  ios:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: user_input
    value_map:
      "":
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: user_input
        value_map:
          "":
            title: ipa export method
            summary: The export method used to create an .ipa file in your builds,
              stored as an Environment Variable. You can change this at any time,
              or even create several .ipa files with different export methods in the
              same build.
            env_key: BITRISE_EXPORT_METHOD
            type: selector
            value_map:
              ad-hoc:
                config: default-ios-config
              app-store:
                config: default-ios-config
              development:
                config: default-ios-config
              enterprise:
                config: default-ios-config
  macos:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: user_input
    value_map:
      "":
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: user_input
        value_map:
          "":
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
                config: default-macos-config
              developer-id:
                config: default-macos-config
              development:
                config: default-macos-config
              none:
                config: default-macos-config
  react-native:
    title: Was your React Native app created with the Expo CLI and using Managed Workflow?
    summary: Will include *Expo Eject** Step if using Expo Managed Workflow (https://docs.expo.io/introduction/managed-vs-bare/).
      If ios/android native projects are present in the repository, choose No.
    type: selector
    value_map:
      "no":
        title: The root directory of an Android project
        summary: The root directory of your Android project, stored as an Environment
          Variable. In your Workflows, you can specify paths relative to this path.
          You can change this at any time.
        env_key: PROJECT_LOCATION
        type: user_input
        value_map:
          "":
            title: Module
            summary: Modules provide a container for your Android project's source
              code, resource files, and app level settings, such as the module-level
              build file and Android manifest file. Each module can be independently
              built, tested, and debugged. You can add new modules to your Bitrise
              builds at any time.
            env_key: MODULE
            type: user_input
            value_map:
              "":
                title: Variant
                summary: Your Android build variant. You can add variants at any time,
                  as well as further configure your existing variants later.
                env_key: VARIANT
                type: user_input_optional
                value_map:
                  "":
                    title: Project or Workspace path
                    summary: The location of your Xcode project or Xcode workspace
                      files, stored as an Environment Variable. In your Workflows,
                      you can specify paths relative to this path.
                    env_key: BITRISE_PROJECT_PATH
                    type: user_input
                    value_map:
                      "":
                        title: Scheme name
                        summary: An Xcode scheme defines a collection of targets to
                          build, a configuration to use when building, and a collection
                          of tests to execute. Only shared schemes are detected automatically
                          but you can use any scheme as a target on Bitrise. You can
                          change the scheme at any time in your Env Vars.
                        env_key: BITRISE_SCHEME
                        type: user_input
                        value_map:
                          "":
                            title: ipa export method
                            summary: The export method used to create an .ipa file
                              in your builds, stored as an Environment Variable. You
                              can change this at any time, or even create several
                              .ipa files with different export methods in the same
                              build.
                            env_key: BITRISE_EXPORT_METHOD
                            type: selector
                            value_map:
                              ad-hoc:
                                config: default-react-native-config
                              app-store:
                                config: default-react-native-config
                              development:
                                config: default-react-native-config
                              enterprise:
                                config: default-react-native-config
      "yes":
        title: The iOS project path generated by the 'expo eject' process
        summary: |-
          Will add the Expo Eject Step to the Workflow to generate the native iOS project, so it can be built and archived.
          Run 'expo eject' in a local environment to determine this value. This experiment then can be undone by deleting the ios and android directories. See https://docs.expo.io/bare/customizing/ for more details.
          For example: './ios/myproject.xcworkspace'.
        env_key: BITRISE_PROJECT_PATH
        type: user_input_optional
        value_map:
          ./ios/< PROJECT NAME >.xcworkspace:
            title: iOS bundle identifier
            summary: |-
              Optional, only needs to be entered if the key expo/ios/bundleIdentifier is not set in 'app.json'.

              Will add the Expo Eject Step to the Workflow to generate the native iOS project, so the IPA can be exported.
              For your convenience, define it here temporarily. To set this value permanently run 'expo eject' in a local environment and commit 'app.json' changes.
              For example: 'com.sample.myapp'.
            env_key: EXPO_BARE_IOS_BUNLDE_ID
            type: user_input
            value_map:
              "":
                title: The iOS native project scheme name
                summary: |-
                  An Xcode scheme defines a collection of targets to build, a configuration to use when building, and a collection of tests to execute. You can change the scheme at any time.

                  Will add the Expo Eject Step to the Workflow to generate the native iOS project, so it can be built and archived.
                  Run 'expo eject' in a local environment to determine this value. This experiment then can be undone by deleting the ios and android directories.
                env_key: BITRISE_SCHEME
                type: user_input
                value_map:
                  "":
                    title: ipa export method
                    summary: The export method used to create an .ipa file in your
                      builds, stored as an Environment Variable. You can change this
                      at any time, or even create several .ipa files with different
                      export methods in the same build.
                    env_key: BITRISE_EXPORT_METHOD
                    type: selector
                    value_map:
                      ad-hoc:
                        title: Android package name
                        summary: |-
                          Optional, only needs to be entered if the key expo/android/package is not set in 'app.json'.

                          Will add the Expo Eject Step to the Workflow to generate the native Android project, so the bundle (AAB) can be built.
                          For your convenience, define it here temporarily. To set this value permanently run 'expo eject' in a local environment and commit 'app.json' changes.
                          For example: 'com.sample.myapp'.
                        env_key: EXPO_BARE_ANDROID_PACKAGE
                        type: user_input_optional
                        value_map:
                          "":
                            title: Project root directory
                            summary: The directory of the 'app.json' or 'package.json'
                              file of your React Native project.
                            env_key: WORKDIR
                            type: user_input
                            value_map:
                              "":
                                title: The root directory of an Android project
                                summary: The root directory of your Android project,
                                  stored as an Environment Variable. In your Workflows,
                                  you can specify paths relative to this path. You
                                  can change this at any time.
                                env_key: PROJECT_LOCATION
                                type: selector
                                value_map:
                                  ./android:
                                    title: Module
                                    summary: Modules provide a container for your
                                      Android project's source code, resource files,
                                      and app level settings, such as the module-level
                                      build file and Android manifest file. Each module
                                      can be independently built, tested, and debugged.
                                      You can add new modules to your Bitrise builds
                                      at any time.
                                    env_key: MODULE
                                    type: user_input
                                    value_map:
                                      app:
                                        title: Variant
                                        summary: Your Android build variant. You can
                                          add variants at any time, as well as further
                                          configure your existing variants later.
                                        env_key: VARIANT
                                        type: user_input_optional
                                        value_map:
                                          Release:
                                            config: default-react-native-expo-config
                      app-store:
                        title: Android package name
                        summary: |-
                          Optional, only needs to be entered if the key expo/android/package is not set in 'app.json'.

                          Will add the Expo Eject Step to the Workflow to generate the native Android project, so the bundle (AAB) can be built.
                          For your convenience, define it here temporarily. To set this value permanently run 'expo eject' in a local environment and commit 'app.json' changes.
                          For example: 'com.sample.myapp'.
                        env_key: EXPO_BARE_ANDROID_PACKAGE
                        type: user_input_optional
                        value_map:
                          "":
                            title: Project root directory
                            summary: The directory of the 'app.json' or 'package.json'
                              file of your React Native project.
                            env_key: WORKDIR
                            type: user_input
                            value_map:
                              "":
                                title: The root directory of an Android project
                                summary: The root directory of your Android project,
                                  stored as an Environment Variable. In your Workflows,
                                  you can specify paths relative to this path. You
                                  can change this at any time.
                                env_key: PROJECT_LOCATION
                                type: selector
                                value_map:
                                  ./android:
                                    title: Module
                                    summary: Modules provide a container for your
                                      Android project's source code, resource files,
                                      and app level settings, such as the module-level
                                      build file and Android manifest file. Each module
                                      can be independently built, tested, and debugged.
                                      You can add new modules to your Bitrise builds
                                      at any time.
                                    env_key: MODULE
                                    type: user_input
                                    value_map:
                                      app:
                                        title: Variant
                                        summary: Your Android build variant. You can
                                          add variants at any time, as well as further
                                          configure your existing variants later.
                                        env_key: VARIANT
                                        type: user_input_optional
                                        value_map:
                                          Release:
                                            config: default-react-native-expo-config
                      development:
                        title: Android package name
                        summary: |-
                          Optional, only needs to be entered if the key expo/android/package is not set in 'app.json'.

                          Will add the Expo Eject Step to the Workflow to generate the native Android project, so the bundle (AAB) can be built.
                          For your convenience, define it here temporarily. To set this value permanently run 'expo eject' in a local environment and commit 'app.json' changes.
                          For example: 'com.sample.myapp'.
                        env_key: EXPO_BARE_ANDROID_PACKAGE
                        type: user_input_optional
                        value_map:
                          "":
                            title: Project root directory
                            summary: The directory of the 'app.json' or 'package.json'
                              file of your React Native project.
                            env_key: WORKDIR
                            type: user_input
                            value_map:
                              "":
                                title: The root directory of an Android project
                                summary: The root directory of your Android project,
                                  stored as an Environment Variable. In your Workflows,
                                  you can specify paths relative to this path. You
                                  can change this at any time.
                                env_key: PROJECT_LOCATION
                                type: selector
                                value_map:
                                  ./android:
                                    title: Module
                                    summary: Modules provide a container for your
                                      Android project's source code, resource files,
                                      and app level settings, such as the module-level
                                      build file and Android manifest file. Each module
                                      can be independently built, tested, and debugged.
                                      You can add new modules to your Bitrise builds
                                      at any time.
                                    env_key: MODULE
                                    type: user_input
                                    value_map:
                                      app:
                                        title: Variant
                                        summary: Your Android build variant. You can
                                          add variants at any time, as well as further
                                          configure your existing variants later.
                                        env_key: VARIANT
                                        type: user_input_optional
                                        value_map:
                                          Release:
                                            config: default-react-native-expo-config
                      enterprise:
                        title: Android package name
                        summary: |-
                          Optional, only needs to be entered if the key expo/android/package is not set in 'app.json'.

                          Will add the Expo Eject Step to the Workflow to generate the native Android project, so the bundle (AAB) can be built.
                          For your convenience, define it here temporarily. To set this value permanently run 'expo eject' in a local environment and commit 'app.json' changes.
                          For example: 'com.sample.myapp'.
                        env_key: EXPO_BARE_ANDROID_PACKAGE
                        type: user_input_optional
                        value_map:
                          "":
                            title: Project root directory
                            summary: The directory of the 'app.json' or 'package.json'
                              file of your React Native project.
                            env_key: WORKDIR
                            type: user_input
                            value_map:
                              "":
                                title: The root directory of an Android project
                                summary: The root directory of your Android project,
                                  stored as an Environment Variable. In your Workflows,
                                  you can specify paths relative to this path. You
                                  can change this at any time.
                                env_key: PROJECT_LOCATION
                                type: selector
                                value_map:
                                  ./android:
                                    title: Module
                                    summary: Modules provide a container for your
                                      Android project's source code, resource files,
                                      and app level settings, such as the module-level
                                      build file and Android manifest file. Each module
                                      can be independently built, tested, and debugged.
                                      You can add new modules to your Bitrise builds
                                      at any time.
                                    env_key: MODULE
                                    type: user_input
                                    value_map:
                                      app:
                                        title: Variant
                                        summary: Your Android build variant. You can
                                          add variants at any time, as well as further
                                          configure your existing variants later.
                                        env_key: VARIANT
                                        type: user_input_optional
                                        value_map:
                                          Release:
                                            config: default-react-native-expo-config
  xamarin:
    title: Path to the Xamarin Solution file
    summary: Your solution file has to contain all the solution configurations you
      wish to use on Bitrise. A solution configuration specifies how projects in the
      solution are to be built and deployed.
    env_key: BITRISE_PROJECT_PATH
    type: user_input
    value_map:
      "":
        title: Xamarin solution configuration
        summary: The Xamarin solution configuration that you wish to run in your first
          build. You can change this at any time in your Workflows.
        env_key: BITRISE_XAMARIN_CONFIGURATION
        type: user_input
        value_map:
          "":
            title: Xamarin solution platform
            env_key: BITRISE_XAMARIN_PLATFORM
            type: user_input
            value_map:
              "":
                config: default-xamarin-config
configs:
  android:
    default-android-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          description: |
            ## How to get a signed APK

            This workflow contains the **Sign APK** step. To sign your APK all you have to do is to:

            1. Click on **Code Signing** tab
            1. Find the **ANDROID KEYSTORE FILE** section
            1. Click or drop your file on the upload file field
            1. Fill the displayed 3 input fields:
             1. **Keystore password**
             1. **Keystore alias**
             1. **Private key password**
            1. Click on **[Save metadata]** button

            That's it! From now on, **Sign APK** step will receive your uploaded files.

            ## To run this workflow

            If you want to run this workflow manually:

            1. Open the app's build list page
            2. Click on **[Start/Schedule a Build]** button
            3. Select **deploy** in **Workflow** dropdown input
            4. Click **[Start Build]** button

            Or if you need this workflow to be started by a GIT event:

            1. Click on **Triggers** tab
            2. Setup your desired event (push/tag/pull) and select **deploy** workflow
            3. Click on **[Done]** and then **[Save]** buttons

            The next change in your repository that matches any of your trigger map event will start **deploy** workflow.
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - cache-pull@%s: {}
          - script@%s:
              title: Do anything with Script step
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - change-android-versioncode-and-versionname@%s:
              inputs:
              - build_gradle_path: $PROJECT_LOCATION/$MODULE/build.gradle
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - sign-apk@%s:
              run_if: '{{getenv "BITRISEIO_ANDROID_KEYSTORE_URL" | ne ""}}'
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
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-lint@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - android-unit-test@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
  cordova:
    default-cordova-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: cordova
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
          - certificate-and-profile-installer@%s: {}
          - npm@%s:
              inputs:
              - command: install
              - workdir: $CORDOVA_WORK_DIR
          - generate-cordova-build-configuration@%s: {}
          - cordova-archive@%s:
              inputs:
              - workdir: $CORDOVA_WORK_DIR
              - platform: $CORDOVA_PLATFORM
              - target: emulator
          - deploy-to-bitrise-io@%s: {}
  fastlane:
    default-fastlane-android-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: android
      app:
        envs:
        - FASTLANE_XCODE_LIST_TIMEOUT: "120"
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
          - fastlane@%s:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
          - deploy-to-bitrise-io@%s: {}
    default-fastlane-ios-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ios
      app:
        envs:
        - FASTLANE_XCODE_LIST_TIMEOUT: "120"
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
          - certificate-and-profile-installer@%s: {}
          - fastlane@%s:
              inputs:
              - lane: $FASTLANE_LANE
              - work_dir: $FASTLANE_WORK_DIR
          - deploy-to-bitrise-io@%s: {}
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
              - export_method: $BITRISE_EXPORT_METHOD
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
              - export_method: $BITRISE_EXPORT_METHOD
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
              - export_method: $BITRISE_EXPORT_METHOD
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
              - export_method: $BITRISE_EXPORT_METHOD
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
  ionic:
    default-ionic-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ionic
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
          - certificate-and-profile-installer@%s: {}
          - npm@%s:
              inputs:
              - command: install
              - workdir: $IONIC_WORK_DIR
          - generate-cordova-build-configuration@%s: {}
          - ionic-archive@%s:
              inputs:
              - workdir: $IONIC_WORK_DIR
              - platform: $IONIC_PLATFORM
              - target: emulator
          - deploy-to-bitrise-io@%s: {}
  ios:
    default-ios-config: |
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
          - recreate-user-schemes@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
          - cocoapods-install@%s: {}
          - xcode-test@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
  macos:
    default-macos-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: macos
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
          - cocoapods-install@%s: {}
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
          - cocoapods-install@%s: {}
          - xcode-test-mac@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
          - deploy-to-bitrise-io@%s: {}
          - cache-push@%s: {}
  other:
    other-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: other
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
          - deploy-to-bitrise-io@%s: {}
  react-native:
    default-react-native-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          description: "## Configure Android part of the deploy workflow\n\nTo generate
            a signed APK:\n\n1. Open the **Workflow** tab of your project on Bitrise.io\n1.
            Add **Sign APK step right after Android Build step**\n1. Click on **Code Signing**
            tab\n1. Find the **ANDROID KEYSTORE FILE** section\n1. Click or drop your file
            on the upload file field\n1. Fill the displayed 3 input fields:\n1. **Keystore
            password**\n1. **Keystore alias**\n1. **Private key password**\n1. Click on
            **[Save metadata]** button\n\nThat's it! From now on, **Sign APK** step will
            receive your uploaded files.\n\n## Configure iOS part of the deploy workflow\n\nTo
            generate IPA:\n\n1. Open the **Workflow** tab of your project on Bitrise.io\n1.
            Click on **Code Signing** tab\n1. Find the **PROVISIONING PROFILE** section\n1.
            Click or drop your file on the upload file field\n1. Find the **CODE SIGNING
            IDENTITY** section\n1. Click or drop your file on the upload file field\n1.
            Click on **Workflows** tab\n1. Select deploy workflow\n1. Select **Xcode Archive
            & Export for iOS** step\n1. Open **Force Build Settings** input group\n1. Specify
            codesign settings\nSet **Force code signing with Development Team**, **Force
            code signing with Code Signing Identity**  \nand **Force code signing with Provisioning
            Profile** inputs regarding to the uploaded codesigning files\n1. Specify manual
            codesign style\nIf the codesigning files, are generated manually on the Apple
            Developer Portal,  \nyou need to explicitly specify to use manual coedsign settings
            \ \n(as ejected rn projects have xcode managed codesigning turned on).  \nTo
            do so, add 'CODE_SIGN_STYLE=\"Manual\"' to 'Additional options for xcodebuild
            call' input\n\n## To run this workflow\n\nIf you want to run this workflow manually:\n\n1.
            Open the app's build list page\n2. Click on **[Start/Schedule a Build]** button\n3.
            Select **deploy** in **Workflow** dropdown input\n4. Click **[Start Build]**
            button\n\nOr if you need this workflow to be started by a GIT event:\n\n1. Click
            on **Triggers** tab\n2. Setup your desired event (push/tag/pull) and select
            **deploy** workflow\n3. Click on **[Done]** and then **[Save]** buttons\n\nThe
            next change in your repository that matches any of your trigger map event will
            start **deploy** workflow.\n"
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - npm@%s:
              inputs:
              - command: install
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
          - certificate-and-profile-installer@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - npm@%s:
              inputs:
              - command: install
          - npm@%s:
              inputs:
              - command: test
          - deploy-to-bitrise-io@%s: {}
    default-react-native-expo-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          description: "## Configure Android part of the deploy workflow\n\nTo generate
            a signed APK:\n\n1. Open the **Workflow** tab of your project on Bitrise.io\n1.
            Add **Sign APK step right after Android Build step**\n1. Click on **Code Signing**
            tab\n1. Find the **ANDROID KEYSTORE FILE** section\n1. Click or drop your file
            on the upload file field\n1. Fill the displayed 3 input fields:\n1. **Keystore
            password**\n1. **Keystore alias**\n1. **Private key password**\n1. Click on
            **[Save metadata]** button\n\nThat's it! From now on, **Sign APK** step will
            receive your uploaded files.\n\n## Configure iOS part of the deploy workflow\n\nTo
            generate IPA:\n\n1. Open the **Workflow** tab of your project on Bitrise.io\n1.
            Click on **Code Signing** tab\n1. Find the **PROVISIONING PROFILE** section\n1.
            Click or drop your file on the upload file field\n1. Find the **CODE SIGNING
            IDENTITY** section\n1. Click or drop your file on the upload file field\n1.
            Click on **Workflows** tab\n1. Select deploy workflow\n1. Select **Xcode Archive
            & Export for iOS** step\n1. Open **Force Build Settings** input group\n1. Specify
            codesign settings\nSet **Force code signing with Development Team**, **Force
            code signing with Code Signing Identity**  \nand **Force code signing with Provisioning
            Profile** inputs regarding to the uploaded codesigning files\n1. Specify manual
            codesign style\nIf the codesigning files, are generated manually on the Apple
            Developer Portal,  \nyou need to explicitly specify to use manual coedsign settings
            \ \n(as ejected rn projects have xcode managed codesigning turned on).  \nTo
            do so, add 'CODE_SIGN_STYLE=\"Manual\"' to 'Additional options for xcodebuild
            call' input\n\n## To run this workflow\n\nIf you want to run this workflow manually:\n\n1.
            Open the app's build list page\n2. Click on **[Start/Schedule a Build]** button\n3.
            Select **deploy** in **Workflow** dropdown input\n4. Click **[Start Build]**
            button\n\nOr if you need this workflow to be started by a GIT event:\n\n1. Click
            on **Triggers** tab\n2. Setup your desired event (push/tag/pull) and select
            **deploy** workflow\n3. Click on **[Done]** and then **[Save]** buttons\n\nThe
            next change in your repository that matches any of your trigger map event will
            start **deploy** workflow.\n"
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - yarn@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - script@%s:
              title: Set bundleIdentifier, packageName for Expo Eject
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -ex

                  appJson="app.json"
                  tmp="/tmp/app.json"
                  jq '.expo.android |= if has("package") or env.EXPO_BARE_ANDROID_PACKAGE == "" or env.EXPO_BARE_ANDROID_PACKAGE == null then . else .package = env.EXPO_BARE_ANDROID_PACKAGE end |
                  .expo.ios |= if has("bundleIdentifier") or env.EXPO_BARE_IOS_BUNLDE_ID == "" or env.EXPO_BARE_IOS_BUNLDE_ID == null then . else .bundleIdentifier = env.EXPO_BARE_IOS_BUNLDE_ID end' <${appJson} >${tmp}
                  [[ $?==0 ]] && mv -f ${tmp} ${appJson}
          - expo-detach@%s:
              inputs:
              - project_path: $WORKDIR
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - certificate-and-profile-installer@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - yarn@%s:
              inputs:
              - workdir: $WORKDIR
              - command: install
          - yarn@%s:
              inputs:
              - workdir: $WORKDIR
              - command: test
          - deploy-to-bitrise-io@%s: {}
  xamarin:
    default-xamarin-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: xamarin
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
          - certificate-and-profile-installer@%s: {}
          - xamarin-user-management@%s:
              run_if: .IsCI
          - nuget-restore@%s: {}
          - xamarin-components-restore@%s: {}
          - xamarin-archive@%s:
              inputs:
              - xamarin_solution: $BITRISE_PROJECT_PATH
              - xamarin_configuration: $BITRISE_XAMARIN_CONFIGURATION
              - xamarin_platform: $BITRISE_XAMARIN_PLATFORM
          - deploy-to-bitrise-io@%s: {}
`, customConfigVersions...)
