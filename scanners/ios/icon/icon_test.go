package icon

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xcode-project/xcodeproj"
)

func Test_getAppIconSetName(t *testing.T) {
	const projectPath = "/Users/lpusok/Develop/keybase-client/shared/react-native/ios/Keybase.xcodeproj"
	const schemeName = "Keybase"

	project, err := xcodeproj.Open(projectPath)
	if err != nil {
		t.Errorf("setup: failed to open project %s, error: %s", projectPath, err)
	}

	log.Printf("name: %s", project.Name)

	scheme, found := project.Scheme(schemeName)
	if !found {
		t.Errorf("setup: scheme (%s) not found in project", schemeName)
	}

	mainTarget, err := mainTargetOfScheme(project, scheme.Name)
	log.Printf("main target: %s", mainTarget.Name)

	tests := []struct {
		name string

		project xcodeproj.XcodeProj
		target  xcodeproj.Target

		want    string
		wantErr bool
	}{
		{
			name:    "Normal",
			project: project,
			target:  mainTarget,
			want:    "AppIcon",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAppIconSetName(tt.project, tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAppIconSetName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAppIconSetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIcon(t *testing.T) {
	type args struct {
		projectPath string
		schemeName  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				projectPath: "/Users/lpusok/Develop/keybase-client/shared/react-native/ios/Keybase.xcodeproj",
				schemeName:  "Keybase",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIcon(tt.args.projectPath, tt.args.schemeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIcon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getIcon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileOpen(t *testing.T) {
	const appIconPath = "/Users/lpusok/Develop/keybase-client/shared/react-native/ios/Keybase/Images.xcassets/AppIcon.appiconset"
	const metaDataFileName = "Contents.json"
	_, err := os.Open(filepath.Join(appIconPath, metaDataFileName))
	if err != nil {
		t.Errorf("file not found")
	}
}

func Test_parseAppIconMetadata(t *testing.T) {
	tests := []struct {
		name    string
		input   io.Reader
		want    []appIcon
		wantErr bool
	}{
		{
			name:  "Minimal",
			input: bytes.NewReader(jsonRaw),
			want: []appIcon{
				{
					Size:     20,
					Filename: "Icon-App-20x20@2x.png",
				},
				{
					Size:     1024,
					Filename: "ItunesArtwork@2x.png",
				},
			},
			wantErr: false,
		},
		{
			name:  "Full",
			input: bytes.NewReader(jsonRawMissingSize),
			want: []appIcon{
				{
					Size:     20,
					Filename: "Icon-App-20x20@2x.png",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAppIconMetadata(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAppIconMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAppIconMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

var jsonRaw = []byte(`
{
	"images" : [
		{
		"size" : "20x20",
		"idiom" : "iphone",
		"filename" : "Icon-App-20x20@2x.png",
		"scale" : "2x"
		},
		{
		"size" : "1024x1024",
		"idiom" : "ios-marketing",
		"filename" : "ItunesArtwork@2x.png",
		"scale" : "1x"
		}
	],
	"info" : {
		"version" : 1,
		"author" : "xcode"
	}
	}
`)

var jsonRawMissingSize = []byte(`
{
	"images" : [
		{
		"size" : "20x20",
		"idiom" : "iphone",
		"filename" : "Icon-App-20x20@2x.png",
		"scale" : "2x"
		},
		{
			"size" : "1024x1024",
			"idiom" : "ios-marketing",
			"scale" : "1x"
		}
	],
	"info" : {
		"version" : 1,
		"author" : "xcode"
	}
	}
`)
