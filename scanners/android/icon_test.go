package android

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func manifestWithIcon(iconName string) string {
	return fmt.Sprintf(`
<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="sample.results.test.multiple.bitrise.com.multipletestresultssample">
    <application
        android:allowBackup="true"
        android:icon="%s"
        android:label="@string/app_name"
        android:roundIcon="@mipmap/ic_launcher_round"
        android:supportsRtl="true"
        android:theme="@style/AppTheme">
        <activity android:name=".MainActivity">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>
</manifest>
`, iconName)
}

func TestLookupIconsMultipleApps(t *testing.T) {
	projectDir:= t.TempDir()
	defer func() {
		if err := os.RemoveAll(projectDir); err != nil {
			t.Logf("Failed to clean up after test, error: %s", err)
		}
	}()

	type dummyAppParams struct {
		projectDir, appName, appManifest, iconFileName string
	}
	createDummyApp := func(params dummyAppParams) {
		appDir := filepath.Join(params.projectDir, params.appName, "src", "main")
		appResDir := filepath.Join(appDir, "res", "mipmap-xxxhdpi")

		if err := os.MkdirAll(appResDir, 0755); err != nil {
			t.Errorf("setup: failed top create dir %s", appResDir)
		}
		if err := os.WriteFile(filepath.Join(appDir, "AndroidManifest.xml"), []byte(params.appManifest), 0755); err != nil {
			t.Error("setup: failed to create file")
		}
		if err := os.WriteFile(filepath.Join(appResDir, params.iconFileName), []byte{}, 0755); err != nil {
			t.Errorf("setup: failed to create file")
		}
	}

	tests := []struct {
		name           string
		dummyAppParams []dummyAppParams
		projectDir     string
		basepath       string
		want           []string
		wantErr        bool
	}{
		{
			name: "multiple android apps",
			dummyAppParams: []dummyAppParams{
				{
					projectDir:   projectDir,
					appName:      "app",
					appManifest:  manifestWithIcon("@mipmap/custom_icon"),
					iconFileName: "ic_launcher.png",
				},
				{
					projectDir:   projectDir,
					appName:      "another_app",
					appManifest:  manifestWithIcon("@mipmap/ic_launcher"),
					iconFileName: "custom_icon.png",
				},
			},
			projectDir: projectDir,
			basepath:   projectDir,
			want: []string{
				filepath.Join(projectDir, "another_app", "src", "main", "res", "mipmap-xxxhdpi", "custom_icon.png"),
				filepath.Join(projectDir, "app", "src", "main", "res", "mipmap-xxxhdpi", "ic_launcher.png"),
			},
		},
		{
			name: "unknown icon format in manifest",
			dummyAppParams: []dummyAppParams{
				{
					projectDir:   projectDir,
					appName:      "app",
					appManifest:  manifestWithIcon("${appIcon}"),
					iconFileName: "ic_launcher.png",
				},
			},
			projectDir: projectDir,
			basepath:   projectDir,
			want: []string{
				filepath.Join(projectDir, "app", "src", "main", "res", "mipmap-xxxhdpi", "ic_launcher.png"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, app := range tt.dummyAppParams {
				createDummyApp(app)
			}

			got, err := lookupIcons(tt.projectDir, tt.basepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupPossibleMatches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LookupPossibleMatches() = %v, want %v", got, tt.want)
			}

			for _, app := range tt.dummyAppParams {
				if err := os.RemoveAll(filepath.Join(tt.projectDir, app.appName)); err != nil {
					t.Logf("Failed to clean up after test, error: %s", err)
				}
			}
		})
	}
}
