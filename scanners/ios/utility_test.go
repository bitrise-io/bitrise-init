package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const workspaceSettingsWithAutocreateSchemesDisabledContent = `<?xml version="1.0" encoding="UTF-8"?>
 <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
 <plist version="1.0">
 <dict>
 	<key>IDEWorkspaceSharedSettings_AutocreateContextsIfNeeded</key>
 	<false/>
 </dict>
 </plist>
 `

func TestNewConfigDescriptor(t *testing.T) {
	descriptor := NewConfigDescriptor(false, "", false, false, true, false, "development", true)
	require.Equal(t, false, descriptor.HasPodfile)
	require.Equal(t, false, descriptor.HasTest)
	require.Equal(t, false, descriptor.HasAppClip)
	require.Equal(t, true, descriptor.HasSPMDependencies)
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
			descriptor:         NewConfigDescriptor(false, "", false, false, false, false, "development", false),
			expectedConfigName: "ios-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "", false, false, false, false, "development", false),
			expectedConfigName: "ios-pod-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", false, false, true, false, "development", false),
			expectedConfigName: "ios-spm-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "bootsrap", false, false, false, false, "development", false),
			expectedConfigName: "ios-carthage-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", true, false, false, false, "development", false),
			expectedConfigName: "ios-test-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", false, false, false, false, "development", true),
			expectedConfigName: "ios-missing-shared-schemes-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "bootstrap", false, false, false, false, "development", false),
			expectedConfigName: "ios-pod-carthage-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "bootstrap", true, false, false, false, "development", false),
			expectedConfigName: "ios-pod-carthage-test-config",
		},
		{
			descriptor:         NewConfigDescriptor(true, "bootstrap", true, false, false, false, "development", true),
			expectedConfigName: "ios-pod-carthage-test-missing-shared-schemes-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", false, true, false, false, "development", false),
			expectedConfigName: "ios-app-clip-development-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", false, true, false, false, "ad-hoc", false),
			expectedConfigName: "ios-app-clip-ad-hoc-config",
		},
		{
			descriptor:         NewConfigDescriptor(false, "", true, true, false, false, "development", false),
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
	t.Run("ios-no-shared-schemes-files-and-autocreate-schemes-disabled", func(t *testing.T) {
		sampleAppDir := t.TempDir()
		sampleAppURL := "https://github.com/bitrise-samples/ios-no-shared-schemes.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		xcodeProjectPath := filepath.Join(sampleAppDir, "BitriseXcode7Sample.xcodeproj")
		projectEmbeddedWorksaceSettingsPth := filepath.Join(xcodeProjectPath, "project.xcworkspace/xcshareddata/WorkspaceSettings.xcsettings")
		require.NoError(t, os.MkdirAll(filepath.Dir(projectEmbeddedWorksaceSettingsPth), os.ModePerm))
		require.NoError(t, fileutil.WriteStringToFile(projectEmbeddedWorksaceSettingsPth, workspaceSettingsWithAutocreateSchemesDisabledContent))

		got, err := ParseProjects(XcodeProjectTypeIOS, sampleAppDir, false, true)
		require.EqualError(t, err, fmt.Sprintf("failed to read Schemes: failed to list Schemes in Project (%s/BitriseXcode7Sample.xcodeproj): no schemes found and the Xcode project's 'Autocreate schemes' option is disabled", sampleAppDir))
		require.Equal(t, DetectResult{}, got)
	})

	t.Run("ios-no-shared-schemes-files", func(t *testing.T) {
		sampleAppDir := t.TempDir()
		sampleAppURL := "https://github.com/bitrise-samples/ios-no-shared-schemes.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		want := DetectResult{
			Warnings: nil,
			Projects: []Project{{
				RelPath:     "BitriseXcode7Sample.xcodeproj",
				IsWorkspace: false,
				Schemes: []Scheme{{
					Name:       "BitriseXcode7Sample",
					Missing:    false,
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
