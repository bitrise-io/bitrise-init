package android

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_modulePathFromBuildScriptPath(t *testing.T) {
	tests := []struct {
		name           string
		projectRootDir string
		buildScriptPth string
		want           string
	}{
		{
			name:           "empty",
			projectRootDir: "./",
			buildScriptPth: "./build.gradle.kts",
			want:           "",
		},
		{
			name:           "module",
			projectRootDir: "./",
			buildScriptPth: "./androidApp/build.gradle.kts",
			want:           "androidApp",
		},
		{
			name:           "project in subdir",
			projectRootDir: ".src",
			buildScriptPth: ".src/androidApp/build.gradle.kts",
			want:           "androidApp",
		},
		{
			name:           "submodule",
			projectRootDir: "./",
			buildScriptPth: "./backend/datastore/build.gradle.kts",
			want:           "backend/datastore",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, modulePathFromBuildScriptPath(tt.projectRootDir, tt.buildScriptPth))
		})
	}
}
