package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePyproject(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    pyprojectInfo
	}{
		{
			name: "poetry name only",
			content: `[tool.poetry]
name = "my-app"
version = "0.1.0"
`,
			want: pyprojectInfo{poetryName: "my-app"},
		},
		{
			name: "package-mode false",
			content: `[tool.poetry]
name = "my-app"
package-mode = false
`,
			want: pyprojectInfo{poetryName: "my-app", poetryPackageModeFalse: true},
		},
		{
			name: "explicit packages declaration",
			content: `[tool.poetry]
name = "my-lib"
packages = [{ include = "my_lib" }]
`,
			want: pyprojectInfo{poetryName: "my-lib", poetryHasPackagesField: true},
		},
		{
			name: "PEP 621 project name only",
			content: `[project]
name = "modern-lib"
`,
			want: pyprojectInfo{projectName: "modern-lib"},
		},
		{
			name: "both project and tool.poetry",
			content: `[project]
name = "modern-lib"

[tool.poetry]
name = "legacy-name"
`,
			want: pyprojectInfo{projectName: "modern-lib", poetryName: "legacy-name"},
		},
		{
			name: "single quotes",
			content: `[tool.poetry]
name = 'my-app'
`,
			want: pyprojectInfo{poetryName: "my-app"},
		},
		{
			name: "fields outside tool.poetry are ignored",
			content: `[tool.poetry.dependencies]
name = "this-is-a-dep-not-the-project"
`,
			want: pyprojectInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePyproject(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}
