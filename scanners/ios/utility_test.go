package ios

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
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

func gitClone(t *testing.T, dir, uri string) {
	fmt.Printf("cloning into: %s\n", dir)
	g, err := git.New(dir)
	require.NoError(t, err)
	require.NoError(t, g.Clone(uri).Run())
}

func TestParseProjects(t *testing.T) {
	t.Run("ios-no-shared-schemes", func(t *testing.T) {
		sampleAppDir := t.TempDir()
		sampleAppURL := "https://github.com/bitrise-samples/ios-no-shared-schemes.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		want := DetectResult{
			Warnings: nil,
			Projects: []Project{{
				IsWorkspace:     false,
				IsPodWorkspace:  false,
				RelPath:         "BitriseXcode7Sample.xcodeproj",
				CarthageCommand: "",
				Warnings: []string{
					`No shared schemes found for project: BitriseXcode7Sample.xcodeproj.
Automatically generated schemes may differ from the ones in your project.
Make sure to <a href="http://devcenter.bitrise.io/ios/frequent-ios-issues/#xcode-scheme-not-found">share your schemes</a> for the expected behaviour.`,
				},
				Schemes: []Scheme{{
					Name:       "BitriseXcode7Sample",
					Missing:    true,
					HasXCTests: true,
					HasAppClip: false,
					Icons:      nil,
				}},
			}},
		}

		// While not ideal, the expectation is that the searchDir is the current directory, due to using relative paths.
		// Enforcing this to allow unit test to pass.
		undoChDir, err := pathutil.RevokableChangeDir(sampleAppDir)
		if err != nil {
			t.Fatalf("%s", err)
		}
		defer func() {
			if err := undoChDir(); err != nil {
				log.TWarnf("failed to restore working dir: %s", err)
			}
		}()

		got, err := ParseProjects(XcodeProjectTypeIOS, sampleAppDir, false, true)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}
