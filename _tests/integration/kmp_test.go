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
			Name:             "Taskman",
			RepoURL:          "https://github.com/bitrise-io/kotlin-multiplatform-sample-taskman.git",
			Branch:           "main",
			ExpectedResult:   kmpTaskmanResultYaml,
			ExpectedVersions: kmpTaskmanResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

var kmpTaskmanResultYaml = fmt.Sprintf(`options:
  kotlin-multiplatform:
    title: The root directory of the Gradle project.
    summary: The root directory of the Gradle project, which contains all source files
      from your project, as well as Gradle files, including the Gradle Wrapper (`+"`gradlew`"+`)
      file.
    env_key: PROJECT_ROOT_DIR
    type: selector
    value_map:
      ./:
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
              - project_root_dir: $PROJECT_ROOT_DIR
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
