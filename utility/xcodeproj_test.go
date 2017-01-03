package utility

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllowXcodeProjExtFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := FilterPaths(paths, AllowXcodeProjExtFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestAllowXCWorkspaceExtFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
	}
	actualFiltered, err := FilterPaths(paths, AllowXCWorkspaceExtFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidEmbeddedWorkspaceRegexpFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := FilterPaths(paths, ForbidEmbeddedWorkspaceRegexpFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidGitDirComponentFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := FilterPaths(paths, ForbidGitDirComponentFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidPodsDirComponentFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := FilterPaths(paths, ForbidPodsDirComponentFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidCarthageDirComponentFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := FilterPaths(paths, ForbidCarthageDirComponentFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestForbidFramworkComponentWithExtensionFilter(t *testing.T) {
	paths := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	expectedFiltered := []string{
		"/Users/bitrise/sample-apps-ios-cocoapods/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/.git/SampleAppWithCocoapods.xcodeproj/project.xcworkspace",
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := FilterPaths(paths, ForbidFramworkComponentWithExtensionFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}
