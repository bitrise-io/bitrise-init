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
			Name:              "sample-ruby-on-rails-rspec-postgres-redis",
			RepoURL:           "https://github.com/bitrise-io/ruby-samples.git",
			RelativeSearchDir: "sample-ruby-on-rails-rspec-postgres-redis",
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
	steps.CacheRestoreVersion,
	steps.ScriptVersion,
	steps.ScriptVersion,
	steps.ScriptVersion,
	steps.CacheSaveVersion,
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
        config: ruby-root-bundler-rspec-postgres-redis-config
configs:
  ruby:
    ruby-root-bundler-rspec-postgres-redis-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ruby
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

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Ruby version in these files: .tool-versions, .ruby-version
                  # See: https://github.com/asdf-vm/asdf-ruby
                  asdf install ruby
          - restore-cache@%s:
              inputs:
              - key: gem-{{ checksum "Gemfile.lock" }}
          - script@%s:
              title: Install dependencies
              inputs:
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle config set --local path vendor/bundle
                  bundle install
          - script@%s:
              title: Database setup
              service_containers:
              - postgres
              - redis
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rake db:create db:schema:load
          - script@%s:
              title: Run tests
              service_containers:
              - postgres
              - redis
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rspec
          - save-cache@%s:
              inputs:
              - key: gem-{{ checksum "Gemfile.lock" }}
              - paths: vendor/bundle
          - deploy-to-bitrise-io@%s: {}
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
        redis:
          type: service
          image: redis:7
          ports:
          - 6379:6379
          options: --health-cmd "redis-cli ping" --health-interval 10s --health-timeout
            5s --health-retries 5
warnings:
  ruby: []
warnings_with_recommendations:
  ruby: []
`, rubyResultVersions...)
