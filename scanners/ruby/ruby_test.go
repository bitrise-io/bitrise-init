package ruby

import (
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigName(t *testing.T) {
	tests := []struct {
		name       string
		descriptor configDescriptor
		want       string
	}{
		{
			name: "default config",
			descriptor: configDescriptor{
				isDefault: true,
			},
			want: "default-ruby-config",
		},
		{
			name: "root with bundler and rspec",
			descriptor: configDescriptor{
				workdir:       "",
				hasBundler:    true,
				testFramework: "rspec",
			},
			want: "ruby-root-bundler-rspec-config",
		},
		{
			name: "subdirectory with bundler and minitest",
			descriptor: configDescriptor{
				workdir:       "$RUBY_PROJECT_DIR",
				hasBundler:    true,
				testFramework: "minitest",
			},
			want: "ruby-bundler-minitest-config",
		},
		{
			name: "no bundler, no test framework",
			descriptor: configDescriptor{
				workdir:    "$RUBY_PROJECT_DIR",
				hasBundler: false,
			},
			want: "ruby-config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, configName(tt.descriptor))
		})
	}
}

func TestGenerateTestScript(t *testing.T) {
	tests := []struct {
		name       string
		descriptor configDescriptor
		wantEmpty  bool
		contains   string
	}{
		{
			name: "rspec with bundler",
			descriptor: configDescriptor{
				hasBundler:    true,
				testFramework: "rspec",
			},
			contains: "bundle exec rspec",
		},
		{
			name: "rspec without bundler",
			descriptor: configDescriptor{
				hasBundler:    false,
				testFramework: "rspec",
			},
			contains: "rspec",
		},
		{
			name: "minitest with rakefile and bundler",
			descriptor: configDescriptor{
				hasBundler:    true,
				hasRakefile:   true,
				testFramework: "minitest",
			},
			contains: "bundle exec rake test",
		},
		{
			name: "minitest without rakefile",
			descriptor: configDescriptor{
				hasBundler:    true,
				hasRakefile:   false,
				testFramework: "minitest",
			},
			contains: "bundle exec ruby -Itest test/**/*_test.rb",
		},
		{
			name: "no test framework but has rakefile",
			descriptor: configDescriptor{
				hasBundler:  true,
				hasRakefile: true,
			},
			contains: "bundle exec rake test",
		},
		{
			name: "no test framework, no rakefile",
			descriptor: configDescriptor{
				hasBundler:  true,
				hasRakefile: false,
			},
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTestScript(tt.descriptor)
			if tt.wantEmpty {
				require.Empty(t, result)
			} else {
				require.Contains(t, result, tt.contains)
			}
		})
	}
}

func TestDetectDatabaseGemsFromContent(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		wantContainers []string // expected container names
	}{
		{
			name:           "postgres gem with single quotes",
			content:        "gem 'pg'\ngem 'rails'",
			wantContainers: []string{"postgres"},
		},
		{
			name:           "postgres gem with double quotes and version",
			content:        `gem "pg", "~> 1.5"`,
			wantContainers: []string{"postgres"},
		},
		{
			name:           "multiple databases",
			content:        "gem 'pg'\ngem 'redis'\ngem 'mongoid'",
			wantContainers: []string{"postgres", "redis", "mongo"},
		},
		{
			name:           "commented out gem is ignored",
			content:        "# gem 'pg'",
			wantContainers: nil,
		},
		{
			name:           "no database gems",
			content:        "gem 'rails'\ngem 'rspec'",
			wantContainers: nil,
		},
		{
			name:           "mongoid and mongo deduplicated",
			content:        "gem 'mongoid'\ngem 'mongo'",
			wantContainers: []string{"mongo"},
		},
		{
			name:           "mysql2 gem",
			content:        "gem 'mysql2'",
			wantContainers: []string{"mysql"},
		},
		{
			name:           "sqlite3 is not detected",
			content:        "gem 'sqlite3'",
			wantContainers: nil,
		},
		{
			name:           "gem in group block",
			content:        "group :production do\n  gem 'pg'\nend",
			wantContainers: []string{"postgres"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectDatabaseGemsFromContent(tt.content)
			var gotContainers []string
			for _, db := range result {
				gotContainers = append(gotContainers, db.containerName)
			}
			assert.Equal(t, tt.wantContainers, gotContainers)
		})
	}
}

func TestConfigNameWithDatabase(t *testing.T) {
	tests := []struct {
		name       string
		descriptor configDescriptor
		want       string
	}{
		{
			name: "root with bundler, rspec, and postgres",
			descriptor: configDescriptor{
				workdir:       "",
				hasBundler:    true,
				testFramework: "rspec",
				databases:     []databaseGem{{containerName: "postgres"}},
			},
			want: "ruby-root-bundler-rspec-postgres-config",
		},
		{
			name: "with mysql",
			descriptor: configDescriptor{
				workdir:       "",
				hasBundler:    true,
				testFramework: "minitest",
				databases:     []databaseGem{{containerName: "mysql"}},
			},
			want: "ruby-root-bundler-minitest-mysql-config",
		},
		{
			name: "multiple databases uses first",
			descriptor: configDescriptor{
				workdir:       "",
				hasBundler:    true,
				testFramework: "rspec",
				databases: []databaseGem{
					{containerName: "postgres"},
					{containerName: "redis"},
				},
			},
			want: "ruby-root-bundler-rspec-postgres-config",
		},
		{
			name: "no databases keeps original name",
			descriptor: configDescriptor{
				workdir:       "",
				hasBundler:    true,
				testFramework: "rspec",
			},
			want: "ruby-root-bundler-rspec-config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, configName(tt.descriptor))
		})
	}
}

func TestExtractEnvVarFromValue(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantEnvName string
		wantDefault string
	}{
		{
			name:        "ENV.fetch with default",
			value:       `<%= ENV.fetch("DB_PASSWORD") { "password" } %>`,
			wantEnvName: "DB_PASSWORD",
			wantDefault: "password",
		},
		{
			name:        "ENV bracket without default",
			value:       `<%= ENV["MY_DB_PASS"] %>`,
			wantEnvName: "MY_DB_PASS",
			wantDefault: "",
		},
		{
			name:        "plain value",
			value:       "postgres",
			wantEnvName: "",
			wantDefault: "postgres",
		},
		{
			name:        "empty value",
			value:       "",
			wantEnvName: "",
			wantDefault: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractEnvVarFromValue(tt.value)
			assert.Equal(t, tt.wantEnvName, result.name)
			assert.Equal(t, tt.wantDefault, result.defaultValue)
		})
	}
}

func TestParseDatabaseYMLContent(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		wantHostEnv     string
		wantHostDefault string
		wantUsernameEnv string
		wantPasswordEnv string
	}{
		{
			name: "test section with ENV.fetch fields",
			content: `default: &default
  adapter: postgresql

test:
  <<: *default
  database: myapp_test
  username: <%= ENV.fetch("DB_USERNAME") { "postgres" } %>
  password: <%= ENV.fetch("DB_PASSWORD") { "password" } %>
  host: <%= ENV.fetch("DB_HOST") { "localhost" } %>

production:
  <<: *default
  database: myapp_prod`,
			wantHostEnv:     "DB_HOST",
			wantHostDefault: "localhost",
			wantUsernameEnv: "DB_USERNAME",
			wantPasswordEnv: "DB_PASSWORD",
		},
		{
			name: "fields inherited via YAML anchor merge",
			content: `default: &default
  host: <%= ENV.fetch("DB_HOST") { "localhost" } %>
  username: postgres
  password: secret

test:
  <<: *default
  database: myapp_test`,
			wantHostEnv:     "DB_HOST",
			wantHostDefault: "localhost",
			wantUsernameEnv: "",
			wantPasswordEnv: "",
		},
		{
			name: "plain values without ENV references",
			content: `test:
  host: myhost
  username: myuser
  password: mypass`,
			wantHostEnv:     "",
			wantHostDefault: "myhost",
			wantUsernameEnv: "",
			wantPasswordEnv: "",
		},
		{
			name: "no test section falls back to default",
			content: `default: &default
  host: <%= ENV.fetch("DB_HOST") { "localhost" } %>
  username: postgres

development:
  <<: *default
  database: myapp_dev`,
			wantHostEnv:     "DB_HOST",
			wantHostDefault: "localhost",
			wantUsernameEnv: "",
			wantPasswordEnv: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDatabaseYMLContent(tt.content)
			assert.Equal(t, tt.wantHostEnv, result.hostEnvVar.name)
			assert.Equal(t, tt.wantHostDefault, result.hostEnvVar.defaultValue)
			assert.Equal(t, tt.wantUsernameEnv, result.usernameEnvVar.name)
			assert.Equal(t, tt.wantPasswordEnv, result.passwordEnvVar.name)
		})
	}
}

func TestGenerateConfigWithServices(t *testing.T) {
	descriptor := configDescriptor{
		hasBundler:    true,
		testFramework: "rspec",
		databases: []databaseGem{
			{
				gemName:         "pg",
				containerName:   "postgres",
				image:           "postgres:17",
				ports:           []string{"5432:5432"},
				containerEnvKey: "POSTGRES_PASSWORD",
				healthCheck:     `--health-cmd "pg_isready"`,
				isRelationalDB:  true,
			},
		},
		dbYMLInfo: databaseYMLInfo{
			hostEnvVar:     databaseEnvVar{name: "DB_HOST", defaultValue: "localhost"},
			usernameEnvVar: databaseEnvVar{name: "DB_USERNAME", defaultValue: "postgres"},
			passwordEnvVar: databaseEnvVar{name: "DB_PASSWORD", defaultValue: "password"},
		},
	}

	config, err := generateConfigBasedOn(descriptor, models.SSHKeyActivationConditional)
	require.NoError(t, err)

	// Verify containers block
	assert.True(t, strings.Contains(config, "containers:"), "should have containers block")
	assert.True(t, strings.Contains(config, "type: service"), "should have type: service")
	assert.True(t, strings.Contains(config, "image: postgres:17"), "should have postgres image")
	assert.True(t, strings.Contains(config, "POSTGRES_PASSWORD"), "should have POSTGRES_PASSWORD env")

	// Verify service_containers on steps
	assert.True(t, strings.Contains(config, "service_containers:"), "should have service_containers on steps")
	assert.True(t, strings.Contains(config, "- postgres"), "should reference postgres service")

	// Verify database setup step
	assert.True(t, strings.Contains(config, "Database setup"), "should have database setup step")
	assert.True(t, strings.Contains(config, "db:create db:schema:load"), "should have db:create db:schema:load")

	// Verify app-level env vars
	assert.True(t, strings.Contains(config, "DB_HOST: postgres"), "should set DB_HOST to container name")
	assert.True(t, strings.Contains(config, "DB_PASSWORD: password"), "should set DB_PASSWORD default")
}
