package android

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

const appManifest = `
<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="sample.results.test.multiple.bitrise.com.multipletestresultssample">
    <application
        android:allowBackup="true"
        android:icon="@mipmap/ic_launcher"
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
`

func TestLookupIcons(t *testing.T) {
	projectDir, err := ioutil.TempDir("", "android-dummy-project")
	if err != nil {
		t.Errorf("setup: failed to create temp dir")
	}
	defer func() {
		if err := os.RemoveAll(projectDir); err != nil {
			t.Logf("Failed to clean up after test, error: %s", err)
		}
	}()

	app1Dir := filepath.Join(projectDir, "app", "src", "main")
	app2Dir := filepath.Join(projectDir, "another_app", "src", "main")
	app1ResDir := filepath.Join(app1Dir, "res", "mipmap-xxxhdpi")
	app2ResDir := filepath.Join(app2Dir, "res", "mipmap-xxxhdpi")

	if err := os.MkdirAll(app1ResDir, 0755); err != nil {
		t.Errorf("setup: failed top create dir %s", app1ResDir)
	}
	if err := os.MkdirAll(app2ResDir, 0755); err != nil {
		t.Errorf("setup: failed top create dir %s", app2ResDir)
	}

	if err := ioutil.WriteFile(filepath.Join(app1Dir, "AndroidManifest.xml"), []byte(appManifest), 0755); err != nil {
		t.Error("setup: failed to create file")
	}
	if err := ioutil.WriteFile(filepath.Join(app2Dir, "AndroidManifest.xml"), []byte(appManifest), 0755); err != nil {
		t.Error("setup: failed to create file")
	}

	if err := ioutil.WriteFile(filepath.Join(app1ResDir, "ic_launcher.png"), []byte{}, 0755); err != nil {
		t.Errorf("setup: failed to create file")
	}
	if err := ioutil.WriteFile(filepath.Join(app2ResDir, "ic_launcher.png"), []byte{}, 0755); err != nil {
		t.Errorf("setup: failed to create file")
	}

	tests := []struct {
		name       string
		projectDir string
		basepath   string
		want       []string
		wantErr    bool
	}{
		{
			name:       "android sample app",
			projectDir: projectDir,
			basepath:   projectDir,
			want: []string{
				filepath.Join(projectDir, "another_app", "src", "main", "res", "mipmap-xxxhdpi", "ic_launcher.png"),
				filepath.Join(projectDir, "app", "src", "main", "res", "mipmap-xxxhdpi", "ic_launcher.png"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}
