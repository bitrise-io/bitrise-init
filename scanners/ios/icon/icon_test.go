package ios

import "testing"

func Test_getAssetCatalogPath(t *testing.T) {
	tests := []struct {
		name string

		projectPath string
		scheme      string

		want    string
		wantErr bool
	}{
		{
			name:        "Normal",
			projectPath: "/Users/lpusok/Develop/iostest/framework-extension-test/framework-extension-test.xcodeproj",
			scheme:      "framework-extension-test",
			want:        "AppIcon",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAssetCatalogPath(tt.projectPath, tt.scheme)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAssetCatalogPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAssetCatalogPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
