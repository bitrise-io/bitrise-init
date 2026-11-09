package ios

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testXCProjNestedAssetCatalogContent = `{
  "configurations": ["Release"],
  "files": [
    {
      "kind": "group",
      "path": "Branding",
      "children": [
        { "kind": "group", "path": "Icons", "children": [
          { "path": "Icons.xcassets", "target-membership": ["App/resources"] },
        ] },
      ],
    },
  ],
  "targets": [
    {
      "name": "App",
      "id": "A1",
      "product-type": "application",
      "build-phases": ["resources"],
      "build-settings": { "ASSETCATALOG_COMPILER_APPICON_NAME": "AppIcon" },
    },
  ],
}
`

func TestLookupIconByTarget(t *testing.T) {
	tests := []struct {
		name               string
		projectFileName    string
		projectFileContent string
		targetName         string
		appIconSetRelPath  string
	}{
		{
			name:               "project.pbxproj",
			projectFileName:    "project.pbxproj",
			projectFileContent: testMacOSPbxprojContent,
			targetName:         "BitriseStudio",
			appIconSetRelPath:  filepath.Join("BitriseStudio", "Assets.xcassets", "AppIcon.appiconset"),
		},
		{
			name:               "project.xcproj with the asset catalog in nested groups",
			projectFileName:    "project.xcproj",
			projectFileContent: testXCProjNestedAssetCatalogContent,
			targetName:         "App",
			appIconSetRelPath:  filepath.Join("Branding", "Icons", "Icons.xcassets", "AppIcon.appiconset"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			projectPath := filepath.Join(dir, "App.xcodeproj")
			writeTestFile(t, filepath.Join(projectPath, tt.projectFileName), tt.projectFileContent)

			appIconSetPath := filepath.Join(dir, tt.appIconSetRelPath)
			writeTestFile(t, filepath.Join(appIconSetPath, "Contents.json"), `{
  "images": [
    {"size": "16x16", "filename": "icon_16.png"},
    {"size": "512x512", "filename": "icon_512.png"}
  ],
  "info": {"version": 1}
}`)
			writeTestFile(t, filepath.Join(appIconSetPath, "icon_16.png"), "")
			writeTestFile(t, filepath.Join(appIconSetPath, "icon_512.png"), "")

			projectFactory, _ := newProjectFactories()
			project, err := projectFactory.Open(projectPath)
			require.NoError(t, err)
			target, ok := project.TargetByName(tt.targetName)
			require.True(t, ok)

			icons, err := lookupIconByTarget(project, target, dir)
			require.NoError(t, err)

			require.Len(t, icons, 1)
			require.Equal(t, filepath.Join(appIconSetPath, "icon_512.png"), icons[0].Path)
		})
	}
}

func writeTestFile(t *testing.T, pth, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(pth), os.ModePerm))
	require.NoError(t, os.WriteFile(pth, []byte(content), 0o644))
}
