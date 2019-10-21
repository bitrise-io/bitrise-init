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
	XAMARIN: {
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE},
	CORDOVA: {
		OSX_XCODE_10_2_X,
		LINUX_DOCKER_ANDROID_LTS,
		LINUX_DOCKER_ANDROID,
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE,
		OSX_XCODE_11_2_X,
		OSX_XCODE_11_1_X,
		OSX_XCODE_11_0_X,
		OSX_XCODE_10_1_X,
		OSX_XCODE_10_0_X,
		OSX_XCODE_9_4_X,
		OSX_XCODE_8_3_X,
		OSX_XCODE_EDGE},
	REACT_NATIVE: {
		OSX_XCODE_10_2_X,
		LINUX_DOCKER_ANDROID_LTS,
		LINUX_DOCKER_ANDROID,
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE,
		OSX_XCODE_11_2_X,
		OSX_XCODE_11_1_X,
		OSX_XCODE_11_0_X,
		OSX_XCODE_10_1_X,
		OSX_XCODE_10_0_X,
		OSX_XCODE_9_4_X,
		OSX_XCODE_8_3_X,
		OSX_XCODE_EDGE},
	IONIC: {
		OSX_XCODE_10_2_X,
		LINUX_DOCKER_ANDROID_LTS,
		LINUX_DOCKER_ANDROID,
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE,
		OSX_XCODE_11_2_X,
		OSX_XCODE_11_1_X,
		OSX_XCODE_11_0_X,
		OSX_XCODE_10_1_X,
		OSX_XCODE_10_0_X,
		OSX_XCODE_9_4_X,
		OSX_XCODE_8_3_X,
		OSX_XCODE_EDGE,
	},
	FLUTTER: {
		OSX_XCODE_10_2_X,
		LINUX_DOCKER_ANDROID_LTS,
		LINUX_DOCKER_ANDROID,
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE,
		OSX_XCODE_11_2_X,
		OSX_XCODE_11_1_X,
		OSX_XCODE_11_0_X,
		OSX_XCODE_10_1_X,
		OSX_XCODE_10_0_X,
		OSX_XCODE_9_4_X,
		OSX_XCODE_8_3_X,
		OSX_XCODE_EDGE,
	},
	ANDROID: {
		LINUX_DOCKER_ANDROID,
		LINUX_DOCKER_ANDROID_LTS,
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE,
		OSX_XCODE_11_2_X,
		OSX_XCODE_11_1_X,
		OSX_XCODE_11_0_X,
		OSX_XCODE_10_2_X,
		OSX_XCODE_10_1_X,
		OSX_XCODE_10_0_X,
		OSX_XCODE_9_4_X,
		OSX_XCODE_8_3_X,
		OSX_XCODE_EDGE,
	},
	MACOS: {
		OSX_XCODE_10_2_X,
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE,
		OSX_XCODE_11_2_X,
		OSX_XCODE_11_1_X,
		OSX_XCODE_11_0_X,
		OSX_XCODE_10_1_X,
		OSX_XCODE_10_0_X,
		OSX_XCODE_9_4_X,
		OSX_XCODE_8_3_X,
		OSX_XCODE_EDGE},
	IOS: {
		OSX_XCODE_10_2_X,
		OSX_VS4MAC_STABLE,
		OSX_VS4MAC_BETA,
		OSX_VS4MAC_PREVIOUS_STABLE,
		OSX_XCODE_11_2_X,
		OSX_XCODE_11_1_X,
		OSX_XCODE_11_0_X,
		OSX_XCODE_10_1_X,
		OSX_XCODE_10_0_X,
		OSX_XCODE_9_4_X,
		OSX_XCODE_8_3_X,
		OSX_XCODE_EDGE},
}

// Stack defines a Stack with it's name.
type Stack string

const (
	LINUX_DOCKER_ANDROID_LTS   Stack = "linux-docker-android-lts"
	LINUX_DOCKER_ANDROID       Stack = "linux-docker-android"
	OSX_VS4MAC_BETA            Stack = "osx-vs4mac-beta"
	OSX_VS4MAC_PREVIOUS_STABLE Stack = "osx-vs4mac-previous-stable"
	OSX_VS4MAC_STABLE          Stack = "osx-vs4mac-stable"
	OSX_XCODE_11_2_X           Stack = "osx-xcode-11.2.x"
	OSX_XCODE_11_1_X           Stack = "osx-xcode-11.1.x"
	OSX_XCODE_11_0_X           Stack = "osx-xcode-11.0.x"
	OSX_XCODE_10_2_X           Stack = "osx-xcode-10.2.x"
	OSX_XCODE_10_1_X           Stack = "osx-xcode-10.1.x"
	OSX_XCODE_10_0_X           Stack = "osx-xcode-10.0.x"
	OSX_XCODE_9_4_X            Stack = "osx-xcode-9.4.x"
	OSX_XCODE_8_3_X            Stack = "osx-xcode-8.3.x"
	OSX_XCODE_EDGE             Stack = "osx-xcode-edge"
)

// Stacks is the array of the all available stacks.
var Stacks = []Stack{
	LINUX_DOCKER_ANDROID_LTS, LINUX_DOCKER_ANDROID, OSX_VS4MAC_BETA, OSX_VS4MAC_PREVIOUS_STABLE, OSX_XCODE_11_2_X,
	OSX_XCODE_11_1_X, OSX_XCODE_11_0_X, OSX_XCODE_10_2_X, OSX_XCODE_10_1_X, OSX_XCODE_10_0_X, OSX_XCODE_9_4_X,
	OSX_XCODE_8_3_X, OSX_XCODE_EDGE,
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
	XAMARIN      Platform = "xamarin"
	CORDOVA      Platform = "cordova"
	REACT_NATIVE Platform = "react-native"
	IONIC        Platform = "ionic"
	FLUTTER      Platform = "flutter"
	ANDROID      Platform = "android"
	MACOS        Platform = "macos"
	IOS          Platform = "ios"
)

// Platforms is the array of all available platforms.
var Platforms = []Platform{
	XAMARIN, CORDOVA, REACT_NATIVE, IONIC, FLUTTER, ANDROID, MACOS, IOS,
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
