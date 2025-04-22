package android

// TODO: restore tests
//import (
//	"fmt"
//	"path/filepath"
//	"testing"
//
//	"github.com/bitrise-io/bitrise-init/models"
//	"github.com/bitrise-io/go-utils/command/git"
//	"github.com/stretchr/testify/require"
//)
//
//func gitClone(t *testing.T, dir, uri string) {
//	fmt.Printf("cloning into: %s\n", dir)
//	g, err := git.New(dir)
//	require.NoError(t, err)
//	require.NoError(t, g.Clone(uri).Run())
//}
//
//func Test_detect(t *testing.T) {
//	t.Run("Sample app", func(t *testing.T) {
//		sampleAppDir := t.TempDir()
//		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-android-sdk22.git"
//		gitClone(t, sampleAppDir, sampleAppURL)
//
//		want := []Project{{
//			RelPath: ".",
//			Icons: models.Icons{{
//				Filename: "81af22c35b03b30a1931a6283349eae094463aa69c52af3afe804b40dbe6dc12.png",
//				Path:     filepath.Join(sampleAppDir, "app", "src", "main", "res", "mipmap-xxxhdpi", "ic_launcher.png"),
//			}},
//		}}
//
//		got, err := detect(sampleAppDir)
//		require.NoError(t, err)
//		require.Equal(t, want, got)
//	})
//}
//
//func Test_moduleFromBuildScriptPath(t *testing.T) {
//	tests := []struct {
//		name           string
//		projectRootDir string
//		buildScriptPth string
//		want           string
//	}{
//		{
//			name:           "empty",
//			projectRootDir: "./",
//			buildScriptPth: "./build.gradle.kts",
//			want:           "",
//		},
//		{
//			name:           "module",
//			projectRootDir: "./",
//			buildScriptPth: "./androidApp/build.gradle.kts",
//			want:           ":androidApp",
//		},
//		{
//			name:           "submodule",
//			projectRootDir: "./",
//			buildScriptPth: "./backend/datastore/build.gradle.kts",
//			want:           ":backend:datastore",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			require.Equal(t, tt.want, moduleFromBuildScriptPath(tt.projectRootDir, tt.buildScriptPth))
//		})
//	}
//}
