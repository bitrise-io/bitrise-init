package integration

import (
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
)

func TestKotlinMultiplatform(t *testing.T) {
	var testCases = []helper.TestCase{
		{
			"Taskman",
			"https://github.com/godrei/Taskman.git",
			"main",
			sampleAppsAndroid22ResultYML,
			sampleAppsAndroid22Versions,
		},
	}

	helper.Execute(t, testCases)
}
