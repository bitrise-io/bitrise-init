package icon

import "testing"

func Test_getAppIconSetName(t *testing.T) {
	tests := []struct {
		name string

		projectPath string
		scheme      string

		want    string
		wantErr bool
	}{
		{
			name:        "Normal",
			projectPath: "/Users/lpusok/Develop/keybase-client/shared/react-native/ios/Keybase.xcodeproj",
			scheme:      "Keybase",
			want:        "AppIcon",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAppIconSetName(tt.projectPath, tt.scheme)
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
