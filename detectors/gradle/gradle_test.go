package gradle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_detectProjectIncludesInContent(t *testing.T) {
	tests := []struct {
		name                     string
		settingGradleFileContent string
		want                     []string
	}{
		{
			name:                     "empty",
			settingGradleFileContent: "",
			want:                     nil,
		},
		{
			name:                     "gradle dsl",
			settingGradleFileContent: "include ':app'",
			want:                     []string{":app"},
		},
		{
			name:                     "kotlin dsl",
			settingGradleFileContent: `include(":androidApp")`,
			want:                     []string{":androidApp"},
		},
		{
			name:                     "multiple components",
			settingGradleFileContent: `include(":backend:datastore")`,
			want:                     []string{":backend:datastore"},
		},
		{
			name: "multiple includes",
			settingGradleFileContent: `include(":androidApp")
//include(":androidBenchmark")
//include(":automotiveApp")
include(":common:car")`,
			want: []string{":androidApp", ":common:car"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectProjectIncludesInContent(tt.settingGradleFileContent)
			require.Equal(t, tt.want, got)
		})
	}
}
