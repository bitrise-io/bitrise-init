package stack

import (
	"github.com/bitrise-io/bitrise-init/stack/platform"
	"github.com/bitrise-io/bitrise-init/stack/tool"
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

// stack infos
var (
	tools = map[string][]string{
		tool.Xcode: {
			"osx-xcode-11.2.x",
			"osx-xcode-11.1.x",
			"osx-xcode-11.0.x",
			"osx-xcode-10.2.x",
			"osx-xcode-10.1.x",
			"osx-xcode-10.0.x",
			"osx-xcode-9.4.x",
			"osx-xcode-8.3.x",
			"osx-xcode-edge",
		},
		tool.Vs4Mac: {
			"osx-vs4mac-stable",
			"osx-vs4mac-beta",
			"osx-vs4mac-previous-stable",
		},
		tool.AndroidSDK: {
			"linux-docker-android",
			"linux-docker-android-lts",
		},
	}

	All = merge(tool.AndroidSDK, tool.Xcode, tool.Vs4Mac)

	StackOptionsMap = map[string][]string{
		platform.Xamarin: merge(
			tool.Vs4Mac,
		),
		platform.Cordova: merge(
			tool.Xcode,
			tool.Vs4Mac,
			tool.AndroidSDK,
		),
		platform.ReactNative: merge(
			tool.Xcode,
			tool.Vs4Mac,
			tool.AndroidSDK,
		),
		platform.Ionic: merge(
			tool.Xcode,
			tool.Vs4Mac,
			tool.AndroidSDK,
		),
		platform.Flutter: merge(
			tool.Xcode,
			tool.Vs4Mac,
			tool.AndroidSDK,
		),
		platform.Android: merge(
			tool.AndroidSDK,
			tool.Xcode,
			tool.Vs4Mac,
		),
		platform.MacOS: merge(
			tool.Xcode,
			tool.Vs4Mac,
		),
		platform.IOS: merge(
			tool.Xcode,
			tool.Vs4Mac,
		),
	}
)

func merge(tt ...string) (o []string) {
	for _, t := range tt {
		o = append(o, tools[t]...)
	}
	return
}
