package ruby

import (
	"testing"

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
