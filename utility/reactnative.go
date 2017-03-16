package utility

import (
	"path/filepath"

	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	reactNativeAndroidProjectFile = "index.android.js"
	reactNativeiOSProjectFile     = "index.ios.js"
	reactNativeAndroidProjectDir  = "android"
	reactNativeiOSProjectDir      = "ios"
	reactNativeTestsDir           = "__tests__"
	reactNativeNodeModulesDir     = "node_modules"
	reactNativeNpmPackageFile     = "package.json"
)

// AllowReactAndroidProjectBaseFilter ...
var AllowReactAndroidProjectBaseFilter = BaseFilter(reactNativeAndroidProjectFile, true)

// AllowReactiOSProjectBaseFilter ...
var AllowReactiOSProjectBaseFilter = BaseFilter(reactNativeiOSProjectFile, true)

// AllowReactNpmPackageBaseFilter ...
var AllowReactNpmPackageBaseFilter = BaseFilter(reactNativeNpmPackageFile, true)

// ForbidReactTestsDir ...
var ForbidReactTestsDir = ComponentFilter(reactNativeTestsDir, false)

// ForbidReactNodeModulesDir ...
var ForbidReactNodeModulesDir = ComponentFilter(reactNativeNodeModulesDir, false)

// HasReactAndroidProjectFileInDirectoryOf ...
func HasReactAndroidProjectFileInDirectoryOf(pth string) bool {
	dir := filepath.Dir(pth)
	reactNativeAndroidProjectFilePth := filepath.Join(dir, reactNativeAndroidProjectFile)
	exists, err := pathutil.IsPathExists(reactNativeAndroidProjectFilePth)
	if err != nil {
		return false
	}
	return exists
}

// HasReactiOSProjectFileInDirectoryOf ...
func HasReactiOSProjectFileInDirectoryOf(pth string) bool {
	dir := filepath.Dir(pth)
	reactNativeiOSProjectFilePth := filepath.Join(dir, reactNativeiOSProjectFile)
	exists, err := pathutil.IsPathExists(reactNativeiOSProjectFilePth)
	if err != nil {
		return false
	}
	return exists
}

// GetReactNativeAndroidProjectDirInDirectoryOf ...
func GetReactNativeAndroidProjectDirInDirectoryOf(pth string) string {
	dir := filepath.Dir(pth)
	reactNativeAndroidProjectDirPth := filepath.Join(dir, reactNativeAndroidProjectDir)
	if exists, err := pathutil.IsDirExists(reactNativeAndroidProjectDirPth); err != nil || !exists {
		return ""
	}
	return reactNativeAndroidProjectDirPth
}

// GetReactNativeiOSProjectDirInDirectoryOf ...
func GetReactNativeiOSProjectDirInDirectoryOf(pth string) string {
	dir := filepath.Dir(pth)
	reactNativeiOSProjectDirPth := filepath.Join(dir, reactNativeiOSProjectDir)
	if exists, err := pathutil.IsDirExists(reactNativeiOSProjectDirPth); err != nil || !exists {
		return ""
	}
	return reactNativeiOSProjectDirPth
}
