package ios

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/pathfilters"
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
	actualFiltered, err := pathutil.FilterPaths(paths, pathfilters.AllowXcodeProjExtFilter)
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
	actualFiltered, err := pathutil.FilterPaths(paths, pathfilters.AllowXCWorkspaceExtFilter)
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
		"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Pods.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Result.xcodeproj",
		"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Result.xcodeproj",
	}
	actualFiltered, err := pathutil.FilterPaths(paths, pathfilters.ForbidEmbeddedWorkspaceRegexpFilter)
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
	actualFiltered, err := pathutil.FilterPaths(paths, pathfilters.ForbidGitDirComponentFilter)
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
	actualFiltered, err := pathutil.FilterPaths(paths, pathfilters.ForbidPodsDirComponentFilter)
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
	actualFiltered, err := pathutil.FilterPaths(paths, pathfilters.ForbidCarthageDirComponentFilter)
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
	actualFiltered, err := pathutil.FilterPaths(paths, pathfilters.ForbidFramworkComponentWithExtensionFilter)
	require.NoError(t, err)
	require.Equal(t, expectedFiltered, actualFiltered)
}

func TestSDKFilter(t *testing.T) {
	dir := t.TempDir()

	iosProject := filepath.Join(dir, "iphoneos.xcodeproj")
	writeTestFile(t, filepath.Join(iosProject, "project.pbxproj"), testIOSPbxprojContent)

	macOSProject := filepath.Join(dir, "macosx.xcodeproj")
	writeTestFile(t, filepath.Join(macOSProject, "project.pbxproj"), testMacOSPbxprojContent)

	iosWorkspace := filepath.Join(dir, "iphoneos.xcworkspace")
	writeTestFile(t, filepath.Join(iosWorkspace, "contents.xcworkspacedata"), `<?xml version="1.0" encoding="UTF-8"?>
<Workspace version = "1.0">
   <FileRef location = "group:missing.xcodeproj"></FileRef>
   <FileRef location = "group:iphoneos.xcodeproj"></FileRef>
</Workspace>
`)

	autoProject := filepath.Join(dir, "auto.xcodeproj")
	writeTestFile(t, filepath.Join(autoProject, "project.pbxproj"),
		sdkFilterTestPbxproj("", `SDKROOT = auto; SUPPORTED_PLATFORMS = "macosx";`))

	macDestinationProject := filepath.Join(dir, "mac-destination.xcodeproj")
	writeTestFile(t, filepath.Join(macDestinationProject, "project.pbxproj"),
		sdkFilterTestPbxproj("SDKROOT = iphoneos;", `SUPPORTED_PLATFORMS = "iphoneos iphonesimulator macosx";`))

	tests := []struct {
		name      string
		path      string
		wantIOS   bool
		wantMacOS bool
	}{
		{name: "iOS project", path: iosProject, wantIOS: true},
		{name: "macOS project", path: macOSProject, wantMacOS: true},
		{name: "workspace with an iOS project and a missing project", path: iosWorkspace, wantIOS: true},
		{name: "target with SDKROOT auto", path: autoProject, wantMacOS: true},
		{name: "iOS target listing macosx in SUPPORTED_PLATFORMS", path: macDestinationProject, wantIOS: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectFactory, workspaceFactory := newProjectFactories()

			gotIOS, err := newSDKFilter(XcodeProjectTypeIOS, projectFactory, workspaceFactory)(tt.path)
			require.NoError(t, err)
			require.Equal(t, tt.wantIOS, gotIOS)

			gotMacOS, err := newSDKFilter(XcodeProjectTypeMacOS, projectFactory, workspaceFactory)(tt.path)
			require.NoError(t, err)
			require.Equal(t, tt.wantMacOS, gotMacOS)
		})
	}
}

func TestSDKFilter_OpenError(t *testing.T) {
	brokenProject := filepath.Join(t.TempDir(), "broken.xcodeproj")
	writeTestFile(t, filepath.Join(brokenProject, "project.pbxproj"), "not a project")
	projectFactory, workspaceFactory := newProjectFactories()

	_, err := newSDKFilter(XcodeProjectTypeIOS, projectFactory, workspaceFactory)(brokenProject)
	require.Error(t, err)
}

func sdkFilterTestPbxproj(projectBuildSettings, targetBuildSettings string) string {
	return fmt.Sprintf(`// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 77;
	objects = {
		A10000000000000000000001 = {isa = PBXProject; buildConfigurationList = A10000000000000000000002; targets = (A10000000000000000000004, ); };
		A10000000000000000000002 = {isa = XCConfigurationList; buildConfigurations = (A10000000000000000000003, ); defaultConfigurationName = Release; };
		A10000000000000000000003 = {isa = XCBuildConfiguration; name = Release; buildSettings = {%s}; };
		A10000000000000000000004 = {isa = PBXNativeTarget; name = App; buildConfigurationList = A10000000000000000000005; productType = "com.apple.product-type.application"; };
		A10000000000000000000005 = {isa = XCConfigurationList; buildConfigurations = (A10000000000000000000006, ); defaultConfigurationName = Release; };
		A10000000000000000000006 = {isa = XCBuildConfiguration; name = Release; buildSettings = {%s}; };
	};
	rootObject = A10000000000000000000001;
}
`, projectBuildSettings, targetBuildSettings)
}
