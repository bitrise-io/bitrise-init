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
			name: "multiple databases includes all container names",
			descriptor: configDescriptor{
				workdir:       "",
				hasBundler:    true,
				testFramework: "rspec",
				databases: []databaseGem{
					{containerName: "postgres"},
					{containerName: "redis"},
				},
			},
			want: "ruby-root-bundler-rspec-postgres-redis-config",
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

func TestGenerateConfigWithMySQLSystemDeps(t *testing.T) {
	descriptor := configDescriptor{
		hasBundler:    true,
		testFramework: "rspec",
		databases: []databaseGem{
			{
				gemName:         "mysql2",
				containerName:   "mysql",
				image:           "mysql:8",
				ports:           []string{"3306:3306"},
				containerEnvKey: "MYSQL_ROOT_PASSWORD",
				healthCheck:     `--health-cmd "mysqladmin ping"`,
				isRelationalDB:  true,
				aptPackages:     []string{"libmariadb-dev"},
				hostValue:       "127.0.0.1",
			},
		},
		dbYMLInfo: databaseYMLInfo{
			hostEnvVar:     databaseEnvVar{name: "DB_HOST", defaultValue: "localhost"},
			usernameEnvVar: databaseEnvVar{name: "DB_USERNAME", defaultValue: "root"},
			passwordEnvVar: databaseEnvVar{name: "DB_PASSWORD", defaultValue: "password"},
		},
	}

	config, err := generateConfigBasedOn(descriptor, models.SSHKeyActivationConditional)
	require.NoError(t, err)

	assert.True(t, strings.Contains(config, "Install system dependencies"), "should have system deps step")
	assert.True(t, strings.Contains(config, "apt-get install -y libmariadb-dev"), "should install libmariadb-dev")
	assert.True(t, strings.Contains(config, "DB_HOST: 127.0.0.1"), "MySQL must use 127.0.0.1, not localhost (localhost uses Unix socket)")
	// System deps step must come before Install dependencies
	sysDepsIdx := strings.Index(config, "Install system dependencies")
	installDepsIdx := strings.Index(config, "Install dependencies")
	assert.True(t, sysDepsIdx < installDepsIdx, "system deps step should come before Install dependencies")
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
	assert.True(t, strings.Contains(config, "DB_HOST: localhost"), "should set DB_HOST to localhost (scripts run on host, not in Docker)")
	assert.True(t, strings.Contains(config, "DB_PASSWORD: password"), "should set DB_PASSWORD default")
}
