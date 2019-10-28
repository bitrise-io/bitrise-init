package stack

import (
	"github.com/bitrise-io/bitrise-init/stack/platform"
	"github.com/bitrise-io/bitrise-init/stack/tool"
)

// stack infos
var (
	tools = map[string][]string{
		tool.Xcode: []string{
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
		tool.Vs4Mac: []string{
			"osx-vs4mac-stable",
			"osx-vs4mac-beta",
			"osx-vs4mac-previous-stable",
		},
		tool.AndroidSDK: []string{
			"linux-docker-android",
			"linux-docker-android-lts",
		},
	}

	All = merge(tool.AndroidSDK, tool.Xcode, tool.Vs4Mac)

	Platforms = map[string][]string{
		platform.Xamarin: merge(
			tool.Vs4Mac,
		),
		platform.Cordova: merge(
			tool.Vs4Mac,
			tool.Xcode,
			tool.AndroidSDK,
		),
		platform.ReactNative: merge(
			tool.Vs4Mac,
			tool.Xcode,
			tool.AndroidSDK,
		),
		platform.Ionic: merge(
			tool.Vs4Mac,
			tool.Xcode,
			tool.AndroidSDK,
		),
		platform.Flutter: merge(
			tool.Vs4Mac,
			tool.Xcode,
			tool.AndroidSDK,
		),
		platform.Android: merge(
			tool.AndroidSDK,
			tool.Vs4Mac,
			tool.Xcode,
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

func merge(ts ...string) (o []string) {
	for _, tool := range ts {
		o = append(o, tools[tool]...)
	}
	return
}
