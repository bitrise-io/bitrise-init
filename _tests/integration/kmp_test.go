package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestKotlinMultiplatform(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"Taskman",
			"https://github.com/bitrise-io/kotlin-multiplatform-sample-taskman.git",
			"main",
			kmpTaskmanResultYaml,
			kmpTaskmanResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

var kmpTaskmanResultYaml = fmt.Sprintf(`options:
  kotlin-multiplatform:
    title: The project's Gradle Wrapper script (gradlew) path.
    summary: The project's Gradle Wrapper script (gradlew) path.
    env_key: GRADLEW_PATH
    type: selector
    value_map:
      ./gradlew:
        config: kotlin-multiplatform-config
configs:
  kotlin-multiplatform:
    kotlin-multiplatform-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: kotlin-multiplatform
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - gradle-unit-test@%s:
              inputs:
              - gradlew_path: $GRADLEW_PATH
          - deploy-to-bitrise-io@%s: {}
warnings:
  kotlin-multiplatform: []
warnings_with_recommendations:
  kotlin-multiplatform: []
`, kmpTaskmanResultVersions...)

var kmpTaskmanResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.GradleUnitTestVersion,
	steps.DeployToBitriseIoVersion,
}
