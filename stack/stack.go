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
var StackOptionsMap = map[PLATFORM][]STACK{
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

// STACK defines a Stack with it's name.
type STACK string

const (
	LINUX_DOCKER_ANDROID_LTS   STACK = "linux-docker-android-lts"
	LINUX_DOCKER_ANDROID       STACK = "linux-docker-android"
	OSX_VS4MAC_BETA            STACK = "osx-vs4mac-beta"
	OSX_VS4MAC_PREVIOUS_STABLE STACK = "osx-vs4mac-previous-stable"
	OSX_VS4MAC_STABLE          STACK = "osx-vs4mac-stable"
	OSX_XCODE_11_2_X           STACK = "osx-xcode-11.2.x"
	OSX_XCODE_11_1_X           STACK = "osx-xcode-11.1.x"
	OSX_XCODE_11_0_X           STACK = "osx-xcode-11.0.x"
	OSX_XCODE_10_2_X           STACK = "osx-xcode-10.2.x"
	OSX_XCODE_10_1_X           STACK = "osx-xcode-10.1.x"
	OSX_XCODE_10_0_X           STACK = "osx-xcode-10.0.x"
	OSX_XCODE_9_4_X            STACK = "osx-xcode-9.4.x"
	OSX_XCODE_8_3_X            STACK = "osx-xcode-8.3.x"
	OSX_XCODE_EDGE             STACK = "osx-xcode-edge"
)

// Stacks is the array of the all available stacks.
var Stacks = []STACK{
	LINUX_DOCKER_ANDROID_LTS, LINUX_DOCKER_ANDROID, OSX_VS4MAC_BETA, OSX_VS4MAC_PREVIOUS_STABLE, OSX_XCODE_11_2_X,
	OSX_XCODE_11_1_X, OSX_XCODE_11_0_X, OSX_XCODE_10_2_X, OSX_XCODE_10_1_X, OSX_XCODE_10_0_X, OSX_XCODE_9_4_X,
	OSX_XCODE_8_3_X, OSX_XCODE_EDGE,
}

// ParseStack gets the given stack from a string.
func ParseStack(s string) (STACK, error) {
	for _, stack := range Stacks {
		if string(stack) == s {
			return stack, nil
		}
	}
	return "", fmt.Errorf("could not find stack %s", s)
}

// StringValue returns the string value of the given STACK.
func (s STACK) StringValue() string {
	return string(s)
}

// PLATFORM defines a platform with it's name.
type PLATFORM string

const (
	XAMARIN      PLATFORM = "xamarin"
	CORDOVA      PLATFORM = "cordova"
	REACT_NATIVE PLATFORM = "react-native"
	IONIC        PLATFORM = "ionic"
	FLUTTER      PLATFORM = "flutter"
	ANDROID      PLATFORM = "android"
	MACOS        PLATFORM = "macos"
	IOS          PLATFORM = "ios"
)

// platforms is the array of all available platforms.
var Platforms = []PLATFORM{
	XAMARIN, CORDOVA, REACT_NATIVE, IONIC, FLUTTER, ANDROID, MACOS, IOS,
}

// ParsePlatform gets the given platform from a string.
func ParsePlatform(s string) (PLATFORM, error) {
	for _, platform := range Platforms {
		if string(platform) == s {
			return platform, nil
		}
	}
	return "", fmt.Errorf("could not find platorm %s", s)
}

// StringValue returns the string value of the given PLATFORM.
func (p PLATFORM) StringValue() string {
	return string(p)
}
