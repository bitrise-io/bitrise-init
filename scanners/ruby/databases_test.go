package ruby

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			wantContainers: []string{"postgres", "redis", "mongodb"},
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
			wantContainers: []string{"mongodb"},
		},
		{
			name:           "mysql2 gem",
			content:        "gem 'mysql2'",
			wantContainers: []string{"mysql"},
		},
		{
			name:           "sqlite3 detected without container",
			content:        "gem 'sqlite3'",
			wantContainers: []string{""},
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
	pgGem := databaseGem{gemName: "pg", adapterName: "postgresql", isRelationalDB: true}
	mysqlGem := databaseGem{gemName: "mysql2", adapterName: "mysql2", isRelationalDB: true}

	tests := []struct {
		name      string
		content   string
		databases []databaseGem
		want      databaseYMLInfo
	}{
		{
			name: "postgresql adapter matches pg gem",
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
			databases: []databaseGem{pgGem},
			want: databaseYMLInfo{
				adapter:        "postgresql",
				hostEnvVar:     databaseEnvVar{name: "DB_HOST", defaultValue: "localhost"},
				usernameEnvVar: databaseEnvVar{name: "DB_USERNAME", defaultValue: "postgres"},
				passwordEnvVar: databaseEnvVar{name: "DB_PASSWORD", defaultValue: "password"},
			},
		},
		{
			name: "mysql2 adapter matches mysql2 gem",
			content: `test:
  adapter: mysql2
  host: myhost
  username: myuser
  password: mypass`,
			databases: []databaseGem{mysqlGem},
			want: databaseYMLInfo{
				adapter:        "mysql2",
				hostEnvVar:     databaseEnvVar{defaultValue: "myhost"},
				usernameEnvVar: databaseEnvVar{defaultValue: "myuser"},
				passwordEnvVar: databaseEnvVar{defaultValue: "mypass"},
			},
		},
		{
			name: "no test section falls back to default",
			content: `default: &default
  host: <%= ENV.fetch("DB_HOST") { "localhost" } %>
  username: postgres

development:
  <<: *default
  database: myapp_dev`,
			databases: []databaseGem{pgGem},
			want:      databaseYMLInfo{},
		},
		{
			name: "adapter mismatch returns empty result",
			content: `test:
  adapter: mysql2
  host: myhost
  username: myuser
  password: mypass`,
			databases: []databaseGem{pgGem},
			want:      databaseYMLInfo{},
		},
		{
			name: "no adapter in database.yml returns empty result",
			content: `default: &default
  host: <%= ENV.fetch("DB_HOST") { "localhost" } %>
  username: postgres
  password: secret

test:
  <<: *default
  database: myapp_test`,
			databases: []databaseGem{pgGem},
			want:      databaseYMLInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDatabaseYMLContent(tt.content, tt.databases)
			assert.Equal(t, tt.want, result)
		})
	}
}
