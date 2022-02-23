package ios

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/go-utils/command/git"
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
				RelPath:     "BitriseXcode7Sample.xcodeproj",
				IsWorkspace: false,
				Warnings: []string{
					`No shared schemes found for project: BitriseXcode7Sample.xcodeproj.
Automatically generated schemes may differ from the ones in your project.
Make sure to <a href="https://support.bitrise.io/hc/en-us/articles/4405779956625">share your schemes</a> for the expected behaviour.`,
				},
				Schemes: []Scheme{{
					Name:       "BitriseXcode7Sample",
					Missing:    true,
					HasXCTests: true,
					Icons:      nil,
				}},
			}},
		}

		got, err := ParseProjects(XcodeProjectTypeIOS, sampleAppDir, false, true)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("ios-cocoapods-at-root", func(t *testing.T) {
		sampleAppDir := t.TempDir()
		sampleAppURL := "https://github.com/bitrise-io/ios-cocoapods-at-root"
		gitClone(t, sampleAppDir, sampleAppURL)

		want := DetectResult{
			Projects: []Project{{
				RelPath:        "iOSMinimalCocoaPodsSample.xcworkspace",
				IsWorkspace:    true,
				IsPodWorkspace: true,
				Schemes: []Scheme{{
					Name:       "iOSMinimalCocoaPodsSample",
					HasXCTests: true,
				}},
			}},
		}

		got, err := ParseProjects(XcodeProjectTypeIOS, sampleAppDir, false, true)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("sample-apps-ios-watchkit", func(t *testing.T) {
		sampleAppDir := t.TempDir()
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-ios-watchkit.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		want := DetectResult{
			Projects: []Project{{
				RelPath:     "watch-test.xcodeproj",
				IsWorkspace: false,
				Schemes: []Scheme{
					{
						Name: "Complication - watch-test WatchKit App",
					},
					{
						Name: "Glance - watch-test WatchKit App",
					},
					{
						Name: "Notification - watch-test WatchKit App",
					},
					{
						Name: "watch-test WatchKit App",
					},
					{
						Name:       "watch-test",
						HasXCTests: true,
					},
				},
			}},
		}

		got, err := ParseProjects(XcodeProjectTypeIOS, sampleAppDir, false, true)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("sample-apps-carthage", func(t *testing.T) {
		sampleAppDir := t.TempDir()
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-carthage.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		want := DetectResult{
			Projects: []Project{{
				RelPath:         "sample-apps-carthage.xcodeproj",
				IsWorkspace:     false,
				CarthageCommand: "bootstrap",
				Schemes: []Scheme{{
					Name:       "sample-apps-carthage",
					HasXCTests: true,
				}},
			}},
		}

		got, err := ParseProjects(XcodeProjectTypeIOS, sampleAppDir, false, true)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("sample-apps-appclip", func(t *testing.T) {
		sampleAppDir := t.TempDir()
		sampleAppURL := "https://github.com/bitrise-io/sample-apps-ios-with-appclip.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		want := DetectResult{
			Projects: []Project{{
				IsWorkspace: true,
				RelPath:     "Sample.xcworkspace",
				Schemes: []Scheme{{
					Name:       "SampleAppClipApp",
					HasAppClip: true,
				}},
			}},
		}

		got, err := ParseProjects(XcodeProjectTypeIOS, sampleAppDir, false, true)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}
