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
			ExpectedResult:    rubyRspecPostgresRedisResultYML,
			ExpectedVersions:  rubyRspecPostgresRedisResultVersions,
		},
		{
			Name:              "sample-ruby-on-rails-minitest-sqlite-mongodb",
			RepoURL:           "https://github.com/bitrise-io/ruby-samples.git",
			RelativeSearchDir: "sample-ruby-on-rails-minitest-sqlite-mongodb",
			Branch:            "main",
			ExpectedResult:    rubyMinitestSqliteMongoDBResultYML,
			ExpectedVersions:  rubyMinitestSqliteMongoDBResultVersions,
		},
		{
			Name:    "ruby-samples-monorepo",
			RepoURL: "https://github.com/bitrise-io/ruby-samples.git",
			Branch:  "main",
			ExpectedResult:   rubyMonorepoResultYML,
			ExpectedVersions: rubyMonorepoResultVersions,
		},
	}

	if len(testCases) > 0 {
		helper.Execute(t, testCases)
	} else {
		t.Skip("No Ruby integration test cases defined yet")
	}
}

var rubyRspecPostgresRedisResultVersions = []interface{}{
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

var rubyRspecPostgresRedisResultYML = fmt.Sprintf(`options:
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
        - REDIS_URL: redis://redis:6379/0
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
`, rubyRspecPostgresRedisResultVersions...)

var rubyMinitestSqliteMongoDBResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreVersion,
	steps.ScriptVersion,
	steps.ScriptVersion,
	steps.CacheSaveVersion,
	steps.DeployToBitriseIoVersion,
}

var rubyMinitestSqliteMongoDBResultYML = fmt.Sprintf(`options:
  ruby:
    title: Project Directory
    summary: The directory containing the Gemfile
    env_key: RUBY_PROJECT_DIR
    type: selector
    value_map:
      .:
        config: ruby-root-bundler-minitest-mongodb-config
configs:
  ruby:
    ruby-root-bundler-minitest-mongodb-config: |
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
              title: Run tests
              service_containers:
              - mongodb
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rake test
          - save-cache@%s:
              inputs:
              - key: gem-{{ checksum "Gemfile.lock" }}
              - paths: vendor/bundle
          - deploy-to-bitrise-io@%s: {}
      containers:
        mongodb:
          type: service
          image: mongo:8
          ports:
          - 27017:27017
          options: --health-cmd "mongosh --eval 'db.runCommand({ping:1})'" --health-interval
            10s --health-timeout 5s --health-retries 5
warnings:
  ruby: []
warnings_with_recommendations:
  ruby: []
`, rubyMinitestSqliteMongoDBResultVersions...)

var rubyMonorepoResultVersions = []interface{}{
	// ruby-bundler-minitest-mongodb-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion, // Install Ruby
	steps.CacheRestoreVersion,
	steps.ScriptVersion, // Install dependencies
	steps.ScriptVersion, // Run tests
	steps.CacheSaveVersion,
	steps.DeployToBitriseIoVersion,
	// ruby-bundler-rspec-mysql-redis-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion, // Install Ruby
	steps.CacheRestoreVersion,
	steps.ScriptVersion, // Install dependencies
	steps.ScriptVersion, // Database setup
	steps.ScriptVersion, // Run tests
	steps.CacheSaveVersion,
	steps.DeployToBitriseIoVersion,
	// ruby-bundler-rspec-postgres-redis-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion, // Install Ruby
	steps.CacheRestoreVersion,
	steps.ScriptVersion, // Install dependencies
	steps.ScriptVersion, // Database setup
	steps.ScriptVersion, // Run tests
	steps.CacheSaveVersion,
	steps.DeployToBitriseIoVersion,
}

var rubyMonorepoResultYML = fmt.Sprintf(`options:
  ruby:
    title: Project Directory
    summary: The directory containing the Gemfile
    env_key: RUBY_PROJECT_DIR
    type: selector
    value_map:
      sample-ruby-on-rails-minitest-sqlite-mongodb:
        config: ruby-bundler-minitest-mongodb-config
      sample-ruby-on-rails-rspec-mysql-redis:
        config: ruby-bundler-rspec-mysql-redis-config
      sample-ruby-on-rails-rspec-postgres-redis:
        config: ruby-bundler-rspec-postgres-redis-config
configs:
  ruby:
    ruby-bundler-minitest-mongodb-config: |
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

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Ruby version in these files: .tool-versions, .ruby-version
                  # See: https://github.com/asdf-vm/asdf-ruby
                  asdf install ruby
              - working_dir: $RUBY_PROJECT_DIR
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
              - working_dir: $RUBY_PROJECT_DIR
          - script@%s:
              title: Run tests
              service_containers:
              - mongodb
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rake test
              - working_dir: $RUBY_PROJECT_DIR
          - save-cache@%s:
              inputs:
              - key: gem-{{ checksum "Gemfile.lock" }}
              - paths: vendor/bundle
          - deploy-to-bitrise-io@%s: {}
      containers:
        mongodb:
          type: service
          image: mongo:8
          ports:
          - 27017:27017
          options: --health-cmd "mongosh --eval 'db.runCommand({ping:1})'" --health-interval
            10s --health-timeout 5s --health-retries 5
    ruby-bundler-rspec-mysql-redis-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ruby
      app:
        envs:
        - DB_HOST: mysql
        - DB_USERNAME: root
        - DB_PASSWORD: password
        - REDIS_URL: redis://redis:6379/0
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
              - working_dir: $RUBY_PROJECT_DIR
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
              - working_dir: $RUBY_PROJECT_DIR
          - script@%s:
              title: Database setup
              service_containers:
              - mysql
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rake db:create db:schema:load
              - working_dir: $RUBY_PROJECT_DIR
          - script@%s:
              title: Run tests
              service_containers:
              - mysql
              - redis
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rspec
              - working_dir: $RUBY_PROJECT_DIR
          - save-cache@%s:
              inputs:
              - key: gem-{{ checksum "Gemfile.lock" }}
              - paths: vendor/bundle
          - deploy-to-bitrise-io@%s: {}
      containers:
        mysql:
          type: service
          image: mysql:8
          ports:
          - 3306:3306
          envs:
          - MYSQL_ROOT_PASSWORD: $DB_PASSWORD
          options: --health-cmd "mysqladmin ping -h 127.0.0.1 -u root --password=$$MYSQL_ROOT_PASSWORD"
            --health-interval 10s --health-timeout 5s --health-retries 5
        redis:
          type: service
          image: redis:7
          ports:
          - 6379:6379
          options: --health-cmd "redis-cli ping" --health-interval 10s --health-timeout
            5s --health-retries 5
    ruby-bundler-rspec-postgres-redis-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ruby
      app:
        envs:
        - DB_HOST: postgres
        - DB_USERNAME: postgres
        - DB_PASSWORD: password
        - REDIS_URL: redis://redis:6379/0
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
              - working_dir: $RUBY_PROJECT_DIR
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
              - working_dir: $RUBY_PROJECT_DIR
          - script@%s:
              title: Database setup
              service_containers:
              - postgres
              inputs:
              - content: |-
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  bundle exec rake db:create db:schema:load
              - working_dir: $RUBY_PROJECT_DIR
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
              - working_dir: $RUBY_PROJECT_DIR
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
`, rubyMonorepoResultVersions...)
