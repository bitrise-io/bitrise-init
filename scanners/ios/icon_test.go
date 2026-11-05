package ios

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupIconByTarget(t *testing.T) {
	dir := t.TempDir()
	projectPath := filepath.Join(dir, "BitriseStudio.xcodeproj")
	writeIconTestFile(t, filepath.Join(projectPath, "project.pbxproj"), testMacOSPbxprojContent)

	appIconSetPath := filepath.Join(dir, "BitriseStudio", "Assets.xcassets", "AppIcon.appiconset")
	writeIconTestFile(t, filepath.Join(appIconSetPath, "Contents.json"), `{
  "images": [
    {"size": "16x16", "filename": "icon_16.png"},
    {"size": "512x512", "filename": "icon_512.png"}
  ],
  "info": {"version": 1}
}`)
	writeIconTestFile(t, filepath.Join(appIconSetPath, "icon_16.png"), "")
	writeIconTestFile(t, filepath.Join(appIconSetPath, "icon_512.png"), "")

	projectFactory, _ := newProjectFactories()
	project, err := projectFactory.Open(projectPath)
	require.NoError(t, err)
	target, ok := project.TargetByName("BitriseStudio")
	require.True(t, ok)

	icons, err := lookupIconByTarget(project, target, dir)
	require.NoError(t, err)

	require.Len(t, icons, 1)
	require.Equal(t, filepath.Join(appIconSetPath, "icon_512.png"), icons[0].Path)
}

func writeIconTestFile(t *testing.T, pth, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(pth), os.ModePerm))
	require.NoError(t, os.WriteFile(pth, []byte(content), 0o644))
}
