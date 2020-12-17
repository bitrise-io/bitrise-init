package ios

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigDescriptor(t *testing.T) {
	descriptor := NewConfigDescriptor(false, "", false, false, "development", true)
	require.Equal(t, false, descriptor.HasPodfile)
	require.Equal(t, false, descriptor.HasTest)
	require.Equal(t, false, descriptor.HasAppClip)
	require.Equal(t, "development", descriptor.ExportMethod)
	require.Equal(t, true, descriptor.MissingSharedSchemes)
	require.Equal(t, "", descriptor.CarthageCommand)
}

func TestConfigName(t *testing.T) {
	type testCase struct {
		descriptor         ConfigDescriptor
		expectedConfigName string
	}

	testCases := []testCase{
		{
			descriptor:         NewConfigDescriptor(false, "", false, false, "development", false),
			expectedConfigName: "ios-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "", false, false, "development", false),
			expectedConfigName: "ios-pod-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "bootsrap", false, false, "development", false),
			expectedConfigName: "ios-carthage-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", true, false, "development", false),
			expectedConfigName: "ios-test-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", false, false, "development", true),
			expectedConfigName: "ios-missing-shared-schemes-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "bootstrap", false, false, "development", false),
			expectedConfigName: "ios-pod-carthage-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "bootstrap", true, false, "development", false),
			expectedConfigName: "ios-pod-carthage-test-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "bootstrap", true, false, "development", true),
			expectedConfigName: "ios-pod-carthage-test-missing-shared-schemes-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", false, true, "development", false),
			expectedConfigName: "ios-app-clip-development-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", false, true, "ad-hoc", false),
			expectedConfigName: "ios-app-clip-ad-hoc-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", true, true, "development", false),
			expectedConfigName: "ios-test-app-clip-development-config",
		},
	}

	for _, testcase := range testCases {
		assert.Equal(t, testcase.expectedConfigName, testcase.descriptor.ConfigName(XcodeProjectTypeIOS))
	}
}
