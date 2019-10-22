package stack

import (
	"fmt"
)

// defaultStacks contains the default stacks for given project types.
//
// Deprecated: please use StackOptionsMap instead.
var defaultStacks = map[string]string{
	"xamarin":      "osx-vs4mac-stable",
	"cordova":      "osx-vs4mac-stable",
	"react-native": "osx-vs4mac-stable",
	"ionic":        "osx-vs4mac-stable",
	"flutter":      "osx-vs4mac-stable",
	"android":      "linux-docker-android",
	"macos":        "osx-xcode-10.0.x",
	"ios":          "osx-xcode-10.0.x",
}

// optionsStacks is the array of all available stacks.
//
// Deprecated: please use the Stacks member defined in this class.
var optionsStacks = []string{
	"linux-docker-android-lts",
	"linux-docker-android",
	"osx-vs4mac-beta",
	"osx-vs4mac-previous-stable",
	"osx-vs4mac-stable",
	"osx-xamarin-stable",
	"osx-xcode-11.0.x",
	"osx-xcode-11.1.x",
	"osx-xcode-11.2.x",
	"osx-xcode-10.0.x",
	"osx-xcode-10.1.x",
	"osx-xcode-10.2.x",
	"osx-xcode-8.3.x",
	"osx-xcode-9.4.x",
	"osx-xcode-edge",
}

// StackOptionsMap holds the available stack for the given platforms, where the first item is the default.
//
// This should be in sync with: https://github.com/bitrise-io/bitrise-website/blob/master/config/available_stacks.yml
var StackOptionsMap = map[Platform][]Stack{
	Xamarin: {
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable},
	Cordova: {
		OsxXcode102X,
		LinuxDockerAndroidLts,
		LinuxDockerAndroid,
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable,
		OsxXcode112X,
		OsxXcode111X,
		OsxXcode110X,
		OsxXcode101X,
		OsxXcode100X,
		OsxXcode94X,
		OsxXcode83X,
		OsxXcodeEdge},
	ReactNative: {
		OsxXcode102X,
		LinuxDockerAndroidLts,
		LinuxDockerAndroid,
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable,
		OsxXcode112X,
		OsxXcode111X,
		OsxXcode110X,
		OsxXcode101X,
		OsxXcode100X,
		OsxXcode94X,
		OsxXcode83X,
		OsxXcodeEdge},
	Ionic: {
		OsxXcode102X,
		LinuxDockerAndroidLts,
		LinuxDockerAndroid,
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable,
		OsxXcode112X,
		OsxXcode111X,
		OsxXcode110X,
		OsxXcode101X,
		OsxXcode100X,
		OsxXcode94X,
		OsxXcode83X,
		OsxXcodeEdge,
	},
	Flutter: {
		OsxXcode102X,
		LinuxDockerAndroidLts,
		LinuxDockerAndroid,
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable,
		OsxXcode112X,
		OsxXcode111X,
		OsxXcode110X,
		OsxXcode101X,
		OsxXcode100X,
		OsxXcode94X,
		OsxXcode83X,
		OsxXcodeEdge,
	},
	Android: {
		LinuxDockerAndroid,
		LinuxDockerAndroidLts,
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable,
		OsxXcode112X,
		OsxXcode111X,
		OsxXcode110X,
		OsxXcode102X,
		OsxXcode101X,
		OsxXcode100X,
		OsxXcode94X,
		OsxXcode83X,
		OsxXcodeEdge,
	},
	MacOS: {
		OsxXcode102X,
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable,
		OsxXcode112X,
		OsxXcode111X,
		OsxXcode110X,
		OsxXcode101X,
		OsxXcode100X,
		OsxXcode94X,
		OsxXcode83X,
		OsxXcodeEdge},
	Ios: {
		OsxXcode102X,
		OsxVs4MACStable,
		OsxVs4MACBeta,
		OsxVs4MACPreviousStable,
		OsxXcode112X,
		OsxXcode111X,
		OsxXcode110X,
		OsxXcode101X,
		OsxXcode100X,
		OsxXcode94X,
		OsxXcode83X,
		OsxXcodeEdge},
}

// Stack defines a Stack with it's name.
type Stack string

const (
	// LinuxDockerAndroidLts the stack for linux-docker-android-lts.
	LinuxDockerAndroidLts Stack = "linux-docker-android-lts"
	// LinuxDockerAndroid the stack for linux-docker-android.
	LinuxDockerAndroid Stack = "linux-docker-android"
	// OsxVs4MACBeta the stack for osx-vs4mac-beta.
	OsxVs4MACBeta Stack = "osx-vs4mac-beta"
	// OsxVs4MACPreviousStable the stack for osx-vs4mac-previous-stable.
	OsxVs4MACPreviousStable Stack = "osx-vs4mac-previous-stable"
	// OsxVs4MACStable the stack for osx-vs4mac-stable.
	OsxVs4MACStable Stack = "osx-vs4mac-stable"
	// OsxXcode112X the stack for osx-xcode-11.2.x.
	OsxXcode112X Stack = "osx-xcode-11.2.x"
	// OsxXcode111X the stack for osx-xcode-11.1.x.
	OsxXcode111X Stack = "osx-xcode-11.1.x"
	// OsxXcode110X the stack for osx-xcode-11.0.x.
	OsxXcode110X Stack = "osx-xcode-11.0.x"
	// OsxXcode102X the stack for osx-xcode-10.2.x.
	OsxXcode102X Stack = "osx-xcode-10.2.x"
	// OsxXcode101X the stack for osx-xcode-10.1.x.
	OsxXcode101X Stack = "osx-xcode-10.1.x"
	// OsxXcode100X the stack for osx-xcode-10.0.x.
	OsxXcode100X Stack = "osx-xcode-10.0.x"
	// OsxXcode94X the stack for osx-xcode-9.4.x
	OsxXcode94X Stack = "osx-xcode-9.4.x"
	// OsxXcode83X the stack for osx-xcode-8.3.x.
	OsxXcode83X Stack = "osx-xcode-8.3.x"
	// OsxXcodeEdge the stack for osx-xcode-edge.
	OsxXcodeEdge Stack = "osx-xcode-edge"
)

// Stacks is the array of the all available stacks.
var Stacks = []Stack{
	LinuxDockerAndroidLts, LinuxDockerAndroid, OsxVs4MACBeta, OsxVs4MACPreviousStable, OsxXcode112X,
	OsxXcode111X, OsxXcode110X, OsxXcode102X, OsxXcode101X, OsxXcode100X, OsxXcode94X,
	OsxXcode83X, OsxXcodeEdge,
}

// ParseStack gets the given stack from a string.
func ParseStack(s string) (Stack, error) {
	for _, stack := range Stacks {
		if string(stack) == s {
			return stack, nil
		}
	}
	return "", fmt.Errorf("could not find stack %s", s)
}

// StringValue returns the string value of the given Stack.
func (s Stack) StringValue() string {
	return string(s)
}

// Platform defines a platform with it's name.
type Platform string

const (
	// Xamarin the constant identifier for Xamarin platform.
	Xamarin Platform = "xamarin"
	// Cordova the constant identifier for Cordova platform.
	Cordova Platform = "cordova"
	// ReactNative the constant identifier for React Native platform.
	ReactNative Platform = "react-native"
	// Ionic the constant identifier for Ionic platform.
	Ionic Platform = "ionic"
	// Flutter the constant identifier for Flutter platform.
	Flutter Platform = "flutter"
	// Android the constant identifier for Android platform.
	Android Platform = "android"
	// MacOS the constant identifier for MacOS platform.
	MacOS Platform = "macos"
	// Ios the constant identifier for iOS platform.
	Ios Platform = "ios"
)

// Platforms is the array of all available platforms.
var Platforms = []Platform{
	Xamarin, Cordova, ReactNative, Ionic, Flutter, Android, MacOS, Ios,
}

// ParsePlatform gets the given platform from a string.
func ParsePlatform(s string) (Platform, error) {
	for _, platform := range Platforms {
		if string(platform) == s {
			return platform, nil
		}
	}
	return "", fmt.Errorf("could not find platorm %s", s)
}

// StringValue returns the string value of the given Platform.
func (p Platform) StringValue() string {
	return string(p)
}
