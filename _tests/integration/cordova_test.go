package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestCordova(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"sample-apps-cordova-with-jasmine",
			"https://github.com/bitrise-samples/sample-apps-cordova-with-jasmine.git",
			"",
			sampleAppsCordovaWithJasmineResultYML,
			sampleAppsCordovaWithJasmineVersions,
		},
		{
			"sample-apps-cordova-with-karma-jasmine",
			"https://github.com/bitrise-samples/sample-apps-cordova-with-karma-jasmine.git",
			"",
			sampleAppsCordovaWithKarmaJasmineResultYML,
			sampleAppsCordovaWithKarmaJasmineVersions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var sampleAppsCordovaWithJasmineVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NpmVersion,
	steps.JasmineTestRunnerVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.CordovaArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.JasmineTestRunnerVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsCordovaWithJasmineResultYML = fmt.Sprintf(`options:
  cordova:
    title: The platform to use in cordova-cli commands
    summary: The target platform for your build, stored as an Environment Variable.
      Your options are iOS, Android, or both. You can change this in your Env Vars
      at any time.
    env_key: CORDOVA_PLATFORM
    type: selector
    value_map:
      android:
        config: cordova-config
      ios:
        config: cordova-config
      ios,android:
        config: cordova-config
configs:
  cordova:
    cordova-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: cordova
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - npm@%s:
              inputs:
              - command: install
          - jasmine-runner@%s: {}
          - generate-cordova-build-configuration@%s: {}
          - cordova-archive@%s:
              inputs:
              - platform: $CORDOVA_PLATFORM
              - target: emulator
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              inputs:
              - command: install
          - jasmine-runner@%s: {}
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  cordova: []
warnings_with_recommendations:
  cordova: []
`, sampleAppsCordovaWithJasmineVersions...)

var sampleAppsCordovaWithKarmaJasmineVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NpmVersion,
	steps.KarmaJasmineTestRunnerVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.CordovaArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.KarmaJasmineTestRunnerVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsCordovaWithKarmaJasmineResultYML = fmt.Sprintf(`options:
  cordova:
    title: The platform to use in cordova-cli commands
    summary: The target platform for your build, stored as an Environment Variable.
      Your options are iOS, Android, or both. You can change this in your Env Vars
      at any time.
    env_key: CORDOVA_PLATFORM
    type: selector
    value_map:
      android:
        config: cordova-config
      ios:
        config: cordova-config
      ios,android:
        config: cordova-config
configs:
  cordova:
    cordova-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: cordova
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - npm@%s:
              inputs:
              - command: install
          - karma-jasmine-runner@%s: {}
          - generate-cordova-build-configuration@%s: {}
          - cordova-archive@%s:
              inputs:
              - platform: $CORDOVA_PLATFORM
              - target: emulator
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              inputs:
              - command: install
          - karma-jasmine-runner@%s: {}
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  cordova: []
warnings_with_recommendations:
  cordova: []
  `, sampleAppsCordovaWithKarmaJasmineVersions...)
