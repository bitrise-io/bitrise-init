package ios

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/stretchr/testify/require"
)

func TestLookupIconByTarget(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "BitriseStudio.xcodeproj")
	writeFile(t, filepath.Join(projectPath, "project.pbxproj"), testMacOSPbxprojContent)

	appIconSetPath := filepath.Join(dir, "BitriseStudio", "Assets.xcassets", "AppIcon.appiconset")
	writeFile(t, filepath.Join(appIconSetPath, "Contents.json"), `{
  "images": [
    {"size": "16x16", "filename": "icon_16.png"},
    {"size": "512x512", "filename": "icon_512.png"}
  ],
  "info": {"version": 1}
}`)
	writeFile(t, filepath.Join(appIconSetPath, "icon_16.png"), "")
	writeFile(t, filepath.Join(appIconSetPath, "icon_512.png"), "")

	projectFactory, _ := newProjectFactories()
	project, err := projectFactory.Open(projectPath)
	require.NoError(t, err)
	target, ok := project.TargetByName("BitriseStudio")
	require.True(t, ok)

	icons, err := lookupIconByTarget(project, target, dir)
	require.NoError(t, err)

	want := models.Icons{{
		Filename: fmt.Sprintf("%x.png", sha256.Sum256([]byte("BitriseStudio/Assets.xcassets/AppIcon.appiconset/icon_512.png"))),
		Path:     filepath.Join(appIconSetPath, "icon_512.png"),
	}}
	require.Equal(t, want, icons)
}

func writeFile(t *testing.T, pth, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(pth), os.ModePerm))
	require.NoError(t, os.WriteFile(pth, []byte(content), 0o644))
}
