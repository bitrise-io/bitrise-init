package scanner

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/bitrise-init/errormapper"

	"github.com/bitrise-io/bitrise-init/models"
)

var GradlewNotFoundRecommendation = errormapper.NewDetailedErrorRecommendation(errormapper.DetailedError{
	Title:       "We couldn’t find your Gradle Wrapper. Please make sure there is a gradlew file in your project’s root directory.",
	Description: `The Gradle Wrapper ensures that the right Gradle version is installed and used for the build. You can find out more about <a target="_blank" href="https://docs.gradle.org/current/userguide/gradle_wrapper.html">the Gradle Wrapper in the Gradle docs</a>.`,
})

var GenericRecommendation = errormapper.NewDetailedErrorRecommendation(newGenericDetail("unexpected end of JSON input"))

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
			args: args{tag: detectPlatformFailedTag, errs: []string{"No Gradle Wrapper (gradlew) found."}},
			want: scannerOutput{errorsWithRecommendation: []models.ErrorWithRecommendations{{Error: "No Gradle Wrapper (gradlew) found.", Recommendations: GradlewNotFoundRecommendation}}},
		},
		{
			name: "Not mapped error",
			args: args{tag: configsFailedTag, errs: []string{"unexpected end of JSON input"}},
			want: scannerOutput{errorsWithRecommendation: []models.ErrorWithRecommendations{{Error: "unexpected end of JSON input", Recommendations: GenericRecommendation}}},
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
