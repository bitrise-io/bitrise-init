package icon

import "testing"

func TestFetchIcon(t *testing.T) {
	type args struct {
		appPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal case",
			args: args{
				appPath: "/Users/lpusok/Develop/AndroidStudioProjects/AppIconTest/app/src/main",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchIcon(tt.args.appPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchIcon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchIcon() = %v, want %v", got, tt.want)
			}
		})
	}
}
