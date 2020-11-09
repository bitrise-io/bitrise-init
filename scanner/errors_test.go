package scanner

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/bitrise-init/step"
)

func Test_mapRecommendation(t *testing.T) {
	type args struct {
		tag string
		err string
	}
	tests := []struct {
		name string
		args args
		want step.Recommendation
	}{
		{
			name: "noPlatformDetected generic error",
			args: args{tag: noPlatformDetectedTag, err: "No known platform detected"},
			want: newDetailedErrorRecommendation(Detail{Title: "We couldn’t recognize your platform.", Description: "Our auto-configurator supports react-native, flutter, ionic, cordova, ios, macos, android, xamarin, fastlane projects. If you’re adding something else, skip this step and configure your Workflow manually."}),
		},
		{
			name: "detectPlatformFailed generic error",
			args: args{tag: detectPlatformFailedTag, err: "No file found at path: Bitrise.xcodeproj/project.pbxproj"},
			want: newDetailedErrorRecommendation(Detail{Title: "We couldn’t parse your project files.", Description: "Our auto-configurator returned the following error:\nNo file found at path: Bitrise.xcodeproj/project.pbxproj"}),
		},
		{
			name: "optionsFailed generic error",
			args: args{tag: optionsFailedTag, err: "No file found at path: ios/App/App/package.json"},
			want: newDetailedErrorRecommendation(Detail{Title: "We couldn’t parse your project files.", Description: "Our auto-configurator returned the following error:\nNo file found at path: ios/App/App/package.json"}),
		},
		{
			name: "optionsFailed gradlew error",
			args: args{tag: optionsFailedTag, err: `<b>No Gradle Wrapper (gradlew) found.</b>
Using a Gradle Wrapper (gradlew) is required, as the wrapper is what makes sure that the right Gradle version is installed and used for the build. More info/guide: <a>https://docs.gradle.org/current/userguide/gradle_wrapper.html</a>`},
			want: newDetailedErrorRecommendation(Detail{Title: "We couldn’t find your Gradle Wrapper. Please make sure there is a gradlew file in your project’s root directory.", Description: `The Gradle Wrapper ensures that the right Gradle version is installed and used for the build. You can find out more about <a target="_blank" href="https://docs.gradle.org/current/userguide/gradle_wrapper.html">the Gradle Wrapper in the Gradle docs</a>.`}),
		},
		{
			name: "optionsFailed app.json error",
			args: args{tag: optionsFailedTag, err: `app.json file (bitrise/app.json) missing or empty name entry
The app.json file needs to contain:
- name
- displayName
entries.`},
			want: newDetailedErrorRecommendation(Detail{Title: "Your app.json file (bitrise/app.json) doesn’t have a name field.", Description: `The app.json file needs to contain the following entries:
- name
- displayName`}),
		},
		{
			name: "optionsFailed Expo app.json error",
			args: args{tag: optionsFailedTag, err: `app.json file (app.json) missing or empty expo/ios/bundleIdentifier entry
If the project uses Expo Kit the app.json file needs to contain:
- expo/name
- expo/ios/bundleIdentifier
- expo/android/package
- entries.`},
			want: newDetailedErrorRecommendation(Detail{Title: "Your app.json file (app.json) doesn’t have a expo/ios/bundleIdentifier field.", Description: `If your project uses Expo Kit, the app.json file needs to contain the following entries:
- expo/name
- expo/ios/bundleIdentifier
- expo/android/package`}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapRecommendation(tt.args.tag, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapRecommendation() = %v, want %v", got, tt.want)
			}
		})
	}
}
