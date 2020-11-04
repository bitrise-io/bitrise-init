package scanner

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
)

var NoKnowPlatformDetectedRecommendation = newDetailedErrorRecommendation(Detail{
	Title:       "We couldn’t recognize your platform.",
	Description: "Our auto-configurator supports react-native, flutter, ionic, cordova, ios, macos, android, xamarin, fastlane projects. If you’re adding something else, skip this step and configure your Workflow manually.",
})

func Test_scannerOutput_AddErrors(t *testing.T) {
	type args struct {
		tag  string
		errs []string
	}
	tests := []struct {
		name string
		args args
		want scannerOutput
	}{
		{
			name: "Mapped error",
			args: args{tag: noPlatformDetectedTag, errs: []string{"No known platform detected"}},
			want: scannerOutput{errorsWithRecommendation: []models.ErrorWithRecommendations{{Error: "No known platform detected", Recommendations: NoKnowPlatformDetectedRecommendation}}},
		},
		{
			name: "Not mapped error",
			args: args{tag: configsFailedTag, errs: []string{"unexpected end of JSON input"}},
			want: scannerOutput{errors: models.Errors{"unexpected end of JSON input"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := scannerOutput{}
			if o.AddErrors(tt.args.tag, tt.args.errs...); !reflect.DeepEqual(o, tt.want) {
				t.Errorf("mapRecommendation() = %v, want %v", o, tt.want)
			}
		})
	}
}
