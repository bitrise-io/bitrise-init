package flutterproject

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitrise-io/bitrise-init/scanners/flutter/flutterproject/mocks"
)

func TestProject_FlutterAndDartSDKVersions(t *testing.T) {
	fileOpener := new(mocks.FileOpener)
	fileOpener.On("OpenFile", ".fvm/fvm_config.json").Return(strings.NewReader(realFVMConfigJSON), nil)
	fileOpener.On("OpenFile", ".tool-versions").Return(strings.NewReader(realToolVersions), nil)
	fileOpener.On("OpenFile", "pubspec.lock").Return(strings.NewReader(realPubspecLock), nil)
	fileOpener.On("OpenFile", "pubspec.yaml").Return(strings.NewReader(realPubspecYaml), nil)

	proj := New("", fileOpener)
	sdkVersions, err := proj.FlutterAndDartSDKVersions()
	require.NoError(t, err)

	b, err := json.MarshalIndent(sdkVersions, "", "\t")
	require.NoError(t, err)

	require.Equal(t, string(b), `{
	"FlutterSDKVersions": [
		{
			"Version": "3.7.12",
			"Constraint": null,
			"Source": "fvm_config_json"
		},
		{
			"Version": "3.7.12",
			"Constraint": null,
			"Source": "tool_versions"
		},
		{
			"Version": null,
			"Constraint": "^3.7.12",
			"Source": "pubspec_yaml"
		}
	],
	"DartSDKVersions": [
		{
			"Version": null,
			"Constraint": "\u003e=2.19.6 \u003c3.0.0",
			"Source": "pubspec_lock"
		},
		{
			"Version": null,
			"Constraint": "\u003e=2.19.6 \u003c3.0.0",
			"Source": "pubspec_yaml"
		}
	]
}`)
}
