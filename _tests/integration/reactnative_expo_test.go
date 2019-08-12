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

func TestReactNativeExpoWithExpoKit(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__reactnative_expo_with_expo_kit__")
	require.NoError(t, err)

	t.Log("BitriseExpoKit")
	{
		sampleAppDir := filepath.Join(tmpDir, "BitriseExpoKit")
		sampleAppURL := "https://github.com/bitrise-samples/BitriseExpoKit.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "BitriseExpoKit", strings.TrimSpace(bitriseExpoKitResultYML), strings.TrimSpace(result), bitriseExpoKitVersions...)
	}
}

func TestReactNativeExpo(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__reactnative_expo__")
	require.NoError(t, err)

	t.Log("BitriseCRNA")
	{
		sampleAppDir := filepath.Join(tmpDir, "BitriseCRNA")
		sampleAppURL := "https://github.com/bitrise-samples/BitriseCRNA.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "BitriseCRNA", strings.TrimSpace(bitriseCRNAResultYML), strings.TrimSpace(result), bitriseCRNAVersions...)
	}
}

var bitriseCRNAVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.ExpoDetachVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,
}

var bitriseCRNAResultYML = fmt.Sprintf(`options:
  react-native:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      ios/BitriseCRNA.xcodeproj:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          BitriseCRNA:
            title: iOS Development team
            summary: The Apple Development Team that the iOS version of the app belongs
              to.
            env_key: BITRISE_IOS_DEVELOPMENT_TEAM
            type: user_input
            value_map:
              "":
                title: ipa export method
                summary: The export method used to create an .ipa file in your builds,
                  stored as an Environment Variable. You can change this at any time,
                  or even create several .ipa files with different export methods
                  in the same build.
                env_key: BITRISE_EXPORT_METHOD
                type: selector
                value_map:
                  ad-hoc:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                config: react-native-expo-config
                  app-store:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                config: react-native-expo-config
                  development:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                config: react-native-expo-config
                  enterprise:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                config: react-native-expo-config
configs:
  react-native:
    react-native-expo-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
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
          - expo-detach@%s:
              inputs:
              - project_path: ./
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
              - configuration: Release
              - export_method: $BITRISE_EXPORT_METHOD
              - force_team_id: $BITRISE_IOS_DEVELOPMENT_TEAM
              - xcodebuild_options: -UseModernBuildSystem=NO
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native: []
`, bitriseCRNAVersions...)

var bitriseExpoKitVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.ExpoDetachVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CocoapodsInstallVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.DeployToBitriseIoVersion,
}

var bitriseExpoKitResultYML = fmt.Sprintf(`options:
  react-native:
    title: Project or Workspace path
    summary: The location of your Xcode project or Xcode workspace files, stored as
      an Environment Variable. In your Workflows, you can specify paths relative to
      this path.
    env_key: BITRISE_PROJECT_PATH
    type: selector
    value_map:
      ios/bitriseexpokit.xcworkspace:
        title: Scheme name
        summary: An Xcode scheme defines a collection of targets to build, a configuration
          to use when building, and a collection of tests to execute. Only shared
          schemes are detected automatically but you can use any scheme as a target
          on Bitrise. You can change the scheme at any time in your Env Vars.
        env_key: BITRISE_SCHEME
        type: selector
        value_map:
          bitriseexpokit:
            title: iOS Development team
            summary: The Apple Development Team that the iOS version of the app belongs
              to.
            env_key: BITRISE_IOS_DEVELOPMENT_TEAM
            type: user_input
            value_map:
              "":
                title: ipa export method
                summary: The export method used to create an .ipa file in your builds,
                  stored as an Environment Variable. You can change this at any time,
                  or even create several .ipa files with different export methods
                  in the same build.
                env_key: BITRISE_EXPORT_METHOD
                type: selector
                value_map:
                  ad-hoc:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                title: Expo username
                                summary: 'Your Expo account username: only required
                                  if you use ExpoKit.'
                                env_key: EXPO_USERNAME
                                type: user_input
                                value_map:
                                  "":
                                    title: Expo password
                                    summary: 'Your Expo account password: only required
                                      if you use ExpoKit.'
                                    env_key: EXPO_PASSWORD
                                    type: user_input
                                    value_map:
                                      "":
                                        config: react-native-expo-config
                  app-store:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                title: Expo username
                                summary: 'Your Expo account username: only required
                                  if you use ExpoKit.'
                                env_key: EXPO_USERNAME
                                type: user_input
                                value_map:
                                  "":
                                    title: Expo password
                                    summary: 'Your Expo account password: only required
                                      if you use ExpoKit.'
                                    env_key: EXPO_PASSWORD
                                    type: user_input
                                    value_map:
                                      "":
                                        config: react-native-expo-config
                  development:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                title: Expo username
                                summary: 'Your Expo account username: only required
                                  if you use ExpoKit.'
                                env_key: EXPO_USERNAME
                                type: user_input
                                value_map:
                                  "":
                                    title: Expo password
                                    summary: 'Your Expo account password: only required
                                      if you use ExpoKit.'
                                    env_key: EXPO_PASSWORD
                                    type: user_input
                                    value_map:
                                      "":
                                        config: react-native-expo-config
                  enterprise:
                    title: The root directory of an Android project
                    summary: The root directory of your Android project, stored as an
                      Environment Variable. In your Workflows, you can specify paths
                      relative to this path. You can change this at any time.
                    env_key: PROJECT_LOCATION
                    type: selector
                    value_map:
                      ./android:
                        title: Module
                        summary: Modules provide a container for your Android project's
                          source code, resource files, and app level settings, such
                          as the module-level build file and Android manifest file.
                          Each module can be independently built, tested, and debugged.
                          You can add new modules to your Bitrise builds at any time.
                        env_key: MODULE
                        type: user_input
                        value_map:
                          app:
                            title: Variant
                            summary: Your Android build variant. You can add
                              variants at any time, as well as further configure your
                              existing variants later.
                            env_key: VARIANT
                            type: user_input_optional
                            value_map:
                              Release:
                                title: Expo username
                                summary: 'Your Expo account username: only required
                                  if you use ExpoKit.'
                                env_key: EXPO_USERNAME
                                type: user_input
                                value_map:
                                  "":
                                    title: Expo password
                                    summary: 'Your Expo account password: only required
                                      if you use ExpoKit.'
                                    env_key: EXPO_PASSWORD
                                    type: user_input
                                    value_map:
                                      "":
                                        config: react-native-expo-config
configs:
  react-native:
    react-native-expo-config: |
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
          - expo-detach@%s:
              inputs:
              - project_path: ./
              - user_name: $EXPO_USERNAME
              - password: $EXPO_PASSWORD
              - run_publish: "yes"
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $VARIANT
          - certificate-and-profile-installer@%s: {}
          - cocoapods-install@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - configuration: Release
              - export_method: $BITRISE_EXPORT_METHOD
              - force_team_id: $BITRISE_IOS_DEVELOPMENT_TEAM
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
warnings:
  react-native: []
`, bitriseExpoKitVersions...)
