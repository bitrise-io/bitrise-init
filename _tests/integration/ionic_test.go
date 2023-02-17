package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestIonic(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"ionic-2",
			"https://github.com/bitrise-samples/ionic-2.git",
			"",
			ionic2ResultYML,
			ionic2Versions,
		},
	}

	helper.Execute(t, testCases)
}

// Expected results

var ionic2Versions = []interface{}{
	models.FormatVersion,
	steps.GitCloneVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.IonicArchiveVersion,
	steps.CacheSaveNPMVersion,
	steps.DeployToBitriseIoVersion,
}

var ionic2ResultYML = fmt.Sprintf(`options:
  ionic:
    title: Directory of the Ionic config.xml file
    summary: The working directory of your Ionic project is where you store your config.xml
      file. This location is stored as an Environment Variable. In your Workflows,
      you can specify paths relative to this path. You can change this at any time.
    env_key: IONIC_WORK_DIR
    type: selector
    value_map:
      cutePuppyPics:
        title: The platform to use in ionic-cli commands
        summary: The target platform for your builds, stored as an Environment Variable.
          Your options are iOS, Android, or both. You can change this in your Env
          Vars at any time.
        env_key: IONIC_PLATFORM
        type: selector
        value_map:
          android:
            config: ionic-config
          ios:
            config: ionic-config
          ios,android:
            config: ionic-config
configs:
  ionic:
    ionic-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ionic
      workflows:
        primary:
          steps:
          - git-clone@%s: {}
          - certificate-and-profile-installer@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              inputs:
              - workdir: $IONIC_WORK_DIR
              - command: install
          - generate-cordova-build-configuration@%s: {}
          - ionic-archive@%s:
              inputs:
              - platform: $IONIC_PLATFORM
              - target: emulator
              - workdir: $IONIC_WORK_DIR
          - save-npm-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ionic: []
warnings_with_recommendations:
  ionic: []
`, ionic2Versions...)
