package integration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/output"
	"github.com/bitrise-io/bitrise-init/scanner"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/stretchr/testify/require"
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
			Name:              "spring-boot-maven-sample",
			RepoURL:           "https://github.com/godrei/java-sample-apps.git",
			RelativeSearchDir: "spring-boot-maven-sample",
			Branch:            "main",
			ExpectedResult:    javaMavenResultYML,
			ExpectedVersions:  javaMavenResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

func TestMissingMavenWrapper(t *testing.T) {
	tmpDir := t.TempDir()
	testName := "java-maven-sample"
	sampleAppDir := filepath.Join(tmpDir, testName)
	sampleAppURL := "https://github.com/godrei/java-sample-apps.git"
	helper.GitClone(t, sampleAppDir, sampleAppURL)

	searchDir := filepath.Join(sampleAppDir, "java-maven-sample")
	_, err := scanner.GenerateAndWriteResults(searchDir, searchDir, output.YAMLFormat)
	require.EqualError(t, err, "No known platform detected")
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
    title: The root directory of the Gradle project.
    summary: The root directory of the Gradle project, which contains all source files
      from your project, as well as Gradle files, including the Gradle Wrapper (`+"`gradlew`"+`)
      file.
    env_key: PROJECT_ROOT_DIR
    type: selector
    value_map:
      ./:
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
              - project_root_dir: $PROJECT_ROOT_DIR
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
	steps.DeployToBitriseIoVersion,
}

var javaMavenResultYML = fmt.Sprintf(`options:
  java:
    title: The root directory of the Maven project.
    summary: The root directory of the Maven project, which contains all source files
      from your project, as well as Maven files, including the Maven Wrapper (`+"`mvn`"+`)
      file.
    env_key: PROJECT_ROOT_DIR
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

                  ./mvnw test
              - working_dir: $PROJECT_ROOT_DIR
          - deploy-to-bitrise-io@%s: {}
warnings:
  java: []
warnings_with_recommendations:
  java: []
`, javaMavenResultVersions...)
