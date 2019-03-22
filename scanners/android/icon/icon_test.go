package icon

import (
	"reflect"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
)

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
			got, err := fetchIcon(tt.args.appPath)
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

func Test_getAllIcons(t *testing.T) {
	tests := []struct {
		name       string
		projectDir string
		want       models.Icons
		wantErr    bool
	}{
		// {
		// 	name:       "normal",
		// 	projectDir: "/Users/lpusok/Develop/AndroidStudioProjects/AppIconTest",
		// 	want:       nil,
		// 	wantErr:    false,
		// },
		{
			name:       "normal",
			projectDir: "/Users/lpusok/Develop/android_github/iosched-master",
			want:       nil,
			wantErr:    false,
		},
		// 	FastHub-development
		// 	Telecine-master
		// 	 iosched-master
		// 	NewPipe-dev
		// 	Telegram-master
		// 	 k-9-master
		// SeeWeather-master
		//  WordPress-Android-develop plaid-master
		// Signal-Android-master
		// android-oss-master
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAllIcons(tt.projectDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllIcons() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAllIcons() = %v, want %v", got, tt.want)
			}
		})
	}
}
