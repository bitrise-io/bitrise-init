package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestRuby(t *testing.T) {
	var testCases = []helper.TestCase{
		// Add test cases here when you have sample Ruby projects
		// Example:
		// {
		// 	Name:              "ruby-sample",
		// 	RepoURL:           "https://github.com/example/ruby-sample.git",
		// 	RelativeSearchDir: ".",
		// 	Branch:            "main",
		// 	ExpectedResult:    rubyResultYML,
		// 	ExpectedVersions:  rubyResultVersions,
		// },
	}

	if len(testCases) > 0 {
		helper.Execute(t, testCases)
	} else {
		t.Skip("No Ruby integration test cases defined yet")
	}
}

var rubyResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreGemVersion,
	steps.CacheSaveGemVersion,
	steps.DeployToBitriseIoVersion,
}

var rubyResultYML = fmt.Sprintf(`options:
  ruby:
    title: Project Directory
    summary: The directory containing the Gemfile
    env_key: RUBY_PROJECT_DIR
    type: selector
    value_map:
      ./:
        config: ruby-root-bundler-rspec-config
configs:
  ruby:
    ruby-root-bundler-rspec-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ruby
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              title: Install Ruby
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  pushd "${RUBY_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Ruby version in these files: .tool-versions, .ruby-version
                  # See: https://github.com/asdf-vm/asdf-ruby
                  asdf install ruby

                  popd > /dev/null
          - restore-cache@%s: {}
          - script@%s:
              title: Install dependencies
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  pushd "${RUBY_PROJECT_DIR:-.}" > /dev/null

                  bundle install

                  popd > /dev/null
          - script@%s:
              title: Run tests
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rspec
          - save-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ruby: []
warnings_with_recommendations:
  ruby: []
`, rubyResultVersions...)
