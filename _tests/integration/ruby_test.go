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
		{
			Name:              "sample-ruby-rails-rspec",
			RepoURL:           "https://github.com/bitrise-io/sample-ruby-rails-rspec",
			RelativeSearchDir: ".",
			Branch:            "main",
			ExpectedResult:    rubyResultYML,
			ExpectedVersions:  rubyResultVersions,
		},
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
	steps.ScriptVersion,
	steps.ScriptVersion,
	steps.ScriptVersion,
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
      .:
        config: ruby-root-bundler-rspec-postgres-config
configs:
  ruby:
    ruby-root-bundler-rspec-postgres-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ruby
      containers:
        postgres:
          type: service
          image: postgres:17
          ports:
          - 5432:5432
          envs:
          - POSTGRES_PASSWORD: $DB_PASSWORD
          options: --health-cmd "pg_isready" --health-interval 10s --health-timeout 5s --health-retries
            5
      app:
        envs:
        - DB_HOST: postgres
        - DB_USERNAME: postgres
        - DB_PASSWORD: password
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
              title: Database setup
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rake db:create db:schema:load
              service_containers:
              - postgres
          - script@%s:
              title: Run tests
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rspec
              service_containers:
              - postgres
          - save-cache@%s: {}
          - deploy-to-bitrise-io@%s: {}
warnings:
  ruby: []
warnings_with_recommendations:
  ruby: []
`, rubyResultVersions...)
