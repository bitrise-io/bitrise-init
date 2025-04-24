package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestJava(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:              "java-gradle-sample",
			RepoURL:           "https://github.com/godrei/java-sample-apps.git",
			RelativeSearchDir: "java-gradle-sample",
			Branch:            "main",
			ExpectedResult:    javaGradleResultYML,
			ExpectedVersions:  javaGradleResultVersions,
		},
		{
			Name:              "ktor-sample",
			RepoURL:           "https://github.com/godrei/java-sample-apps.git",
			RelativeSearchDir: "ktor-sample",
			Branch:            "main",
			ExpectedResult:    javaGradleResultYML,
			ExpectedVersions:  javaGradleResultVersions,
		},
		{
			Name:              "java-maven-sample",
			RepoURL:           "https://github.com/godrei/java-sample-apps.git",
			RelativeSearchDir: "java-maven-sample",
			Branch:            "main",
			ExpectedResult:    javaMavenResultYML,
			ExpectedVersions:  javaMavenResultVersions,
		},
		{
			Name:              "maven-sample",
			RepoURL:           "https://github.com/godrei/java-sample-apps.git",
			RelativeSearchDir: "maven-sample",
			Branch:            "main",
			ExpectedResult:    javaMavenResultYML,
			ExpectedVersions:  javaMavenResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

var javaGradleResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.GradleUnitTestVersion,
	steps.DeployToBitriseIoVersion,
}

var javaGradleResultYML = fmt.Sprintf(`options:
  java:
    title: The project's Gradle Wrapper script (gradlew) path.
    summary: The project's Gradle Wrapper script (gradlew) path.
    env_key: GRADLEW_PATH
    type: selector
    value_map:
      ./gradlew:
        config: java-gradle-config
configs:
  java:
    java-gradle-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: java
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
  java: []
warnings_with_recommendations:
  java: []
`, javaGradleResultVersions...)

var javaMavenResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.ScriptVersion,
	steps.DeployToBitriseIoVersion,
}

var javaMavenResultYML = fmt.Sprintf(`options:
  java:
    title: The root directory of the Maven project (where the pom.xml file is located).
    summary: The root directory of the Maven project (where the pom.xml file is located).
    env_key: MAVEN_PROJECT_ROOT_DIR
    type: selector
    value_map:
      ./:
        config: java-maven-config
configs:
  java:
    java-maven-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: java
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              title: Install Maven
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  # fail if any commands fails
                  set -e
                  # make pipelines' return status equal the last command to exit with a non-zero status, or zero if all commands exit successfully
                  set -o pipefail
                  # debug log
                  set -x

                  sudo apt install maven
          - script@%s:
              title: Run Maven tests
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  # fail if any commands fails
                  set -e
                  # make pipelines' return status equal the last command to exit with a non-zero status, or zero if all commands exit successfully
                  set -o pipefail
                  # debug log
                  set -x

                  mvm test
              - working_dir: $MAVEN_PROJECT_ROOT_DIR
          - deploy-to-bitrise-io@%s: {}
warnings:
  java: []
warnings_with_recommendations:
  java: []
`, javaMavenResultVersions...)
