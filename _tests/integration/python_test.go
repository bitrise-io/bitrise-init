package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestPython(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			Name:              "fastapi-sample",
			RepoURL:           "https://github.com/bitrise-io/python-samples.git",
			RelativeSearchDir: "fastapi-sample",
			Branch:            "main",
			ExpectedResult:    pythonFastapiResultYML,
			ExpectedVersions:  pythonFastapiResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

var pythonFastapiResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreVersion,
	steps.ScriptVersion,
	steps.ScriptVersion,
	steps.CacheSaveVersion,
	steps.DeployToBitriseIoVersion,
}

var pythonFastapiResultYML = fmt.Sprintf(`options:
  python:
    title: Python Project Directory
    summary: The directory containing the Python project files (requirements.txt,
      pyproject.toml, etc.)
    env_key: PYTHON_PROJECT_DIR
    type: selector
    value_map:
      .:
        config: python-root-pip-pytest-config
configs:
  python:
    python-root-pip-pytest-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: python
      meta:
        bitrise.io/framework: fastapi
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-cache@%s:
              inputs:
              - key: pip-{{ checksum "requirements.txt" }}
          - script@%s:
              title: Install dependencies
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  pip install -r requirements.txt
          - script@%s:
              title: Run tests
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  pytest
          - save-cache@%s:
              inputs:
              - key: pip-{{ checksum "requirements.txt" }}
              - paths: ~/.cache/pip
          - deploy-to-bitrise-io@%s: {}
      tools:
        python: "3.12"
warnings:
  python: []
warnings_with_recommendations:
  python: []`,
	pythonFastapiResultVersions...)
