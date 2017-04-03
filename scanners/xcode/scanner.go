package xcode

import "github.com/bitrise-core/bitrise-init/utility"

// FilterProjectFiles ...
func FilterProjectFiles(fileList []string, projectTypes ...ProjectType) ([]string, error) {
	filters := []utility.FilterFunc{
		utility.AllowXcodeProjExtFilter,
	}

	for _, projectType := range projectTypes {
		switch projectType {
		case ProjectTypeIOS:
			filters = append(filters, utility.AllowIphoneosSDKFilter)
		case ProjectTypeMacOS:
			filters = append(filters, utility.AllowMacosxSDKFilter)
		}
	}

	return utility.FilterPaths(fileList, filters...)
}

// FilterWorkspaceFiles ...
func FilterWorkspaceFiles(fileList []string, projectTypes ...ProjectType) ([]string, error) {
	filters := []utility.FilterFunc{
		utility.AllowXCWorkspaceExtFilter,
	}

	for _, projectType := range projectTypes {
		switch projectType {
		case ProjectTypeIOS:
			filters = append(filters, utility.AllowIphoneosSDKFilter)
		case ProjectTypeMacOS:
			filters = append(filters, utility.AllowMacosxSDKFilter)
		}
	}

	return utility.FilterPaths(fileList, filters...)
}

// FilterRelevantProjectFiles ...
func FilterRelevantProjectFiles(fileList []string, projectTypes ...ProjectType) ([]string, error) {
	filters := []utility.FilterFunc{
		utility.AllowXcodeProjExtFilter,
		utility.AllowIsDirectoryFilter,
		utility.ForbidEmbeddedWorkspaceRegexpFilter,
		utility.ForbidGitDirComponentFilter,
		utility.ForbidPodsDirComponentFilter,
		utility.ForbidCarthageDirComponentFilter,
		utility.ForbidFramworkComponentWithExtensionFilter,
	}

	for _, projectType := range projectTypes {
		switch projectType {
		case ProjectTypeIOS:
			filters = append(filters, utility.AllowIphoneosSDKFilter)
		case ProjectTypeMacOS:
			filters = append(filters, utility.AllowMacosxSDKFilter)
		}
	}

	return utility.FilterPaths(fileList, filters...)
}

// FilterRelevantWorkspaceFiles ...
func FilterRelevantWorkspaceFiles(fileList []string, projectTypes ...ProjectType) ([]string, error) {
	filters := []utility.FilterFunc{
		utility.AllowXCWorkspaceExtFilter,
		utility.AllowIsDirectoryFilter,
		utility.ForbidEmbeddedWorkspaceRegexpFilter,
		utility.ForbidGitDirComponentFilter,
		utility.ForbidPodsDirComponentFilter,
		utility.ForbidCarthageDirComponentFilter,
		utility.ForbidFramworkComponentWithExtensionFilter,
	}

	for _, projectType := range projectTypes {
		switch projectType {
		case ProjectTypeIOS:
			filters = append(filters, utility.AllowIphoneosSDKFilter)
		case ProjectTypeMacOS:
			filters = append(filters, utility.AllowMacosxSDKFilter)
		}
	}

	return utility.FilterPaths(fileList, filters...)
}

// FilterRelevantPodfiles ...
func FilterRelevantPodfiles(fileList []string) ([]string, error) {
	return utility.FilterPaths(fileList,
		utility.AllowPodfileBaseFilter,
		utility.ForbidGitDirComponentFilter,
		utility.ForbidPodsDirComponentFilter,
		utility.ForbidCarthageDirComponentFilter,
		utility.ForbidFramworkComponentWithExtensionFilter)
}

// FilterRelevantCartFile ...
func FilterRelevantCartFile(fileList []string) ([]string, error) {
	return utility.FilterPaths(fileList,
		utility.AllowCartfileBaseFilter,
		utility.ForbidGitDirComponentFilter,
		utility.ForbidPodsDirComponentFilter,
		utility.ForbidCarthageDirComponentFilter,
		utility.ForbidFramworkComponentWithExtensionFilter)
}
