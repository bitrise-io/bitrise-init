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
			"https://github.com/godrei/Taskman.git",
			"main",
			kmpTaskmanResultYaml,
			kmpTaskmanResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

var kmpTaskmanResultYaml = fmt.Sprintf(`options:
  kmp:
    title: The project's Gradle Wrapper script (gradlew) path.
    summary: The project's Gradle Wrapper script (gradlew) path.
    env_key: GRADLEW_PATH
    type: selector
    value_map:
      ./gradlew:
        config: kotlin-multiplatform-config
configs:
  kmp:
    kotlin-multiplatform-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: kmp
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - gradle-runner@%s:
              inputs:
              - gradlew_path: $GRADLEW_PATH
              - gradle_task: test
warnings:
  kmp: []
warnings_with_recommendations:
  kmp: []
`, kmpTaskmanResultVersions...)

var kmpTaskmanResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.GradleRunnerVersion,
}
