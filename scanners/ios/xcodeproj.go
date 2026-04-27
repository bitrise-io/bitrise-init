package ios

import (
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-io/go-xcode/pathfilters"
)

// XcodeProjectType ...
type XcodeProjectType string

const (
	// XcodeProjectTypeIOS ...
	XcodeProjectTypeIOS XcodeProjectType = "ios"
	// XcodeProjectTypeMacOS ...
	XcodeProjectTypeMacOS XcodeProjectType = "macos"
)

func findInList(path string, containers []container) (container, bool) {
	for _, container := range containers {
		if container.path() == path {
			return container, true
		}
	}
	return nil, false
}

func removeProjectFromList(projectPth string, projects []container) []container {
	newProjects := []container{}
	for _, project := range projects {
		if project.path() != projectPth {
			newProjects = append(newProjects, project)
		}
	}
	return newProjects
}

func createStandaloneProjectsAndWorkspaces(projectFiles, workspaceFiles []string) (containers, error) {
	var (
		workspaces         []container
		standaloneProjects []container
	)

	for _, workspaceFile := range workspaceFiles {
		workspace, err := newWorkspace(workspaceFile)
		if err != nil {
			return containers{}, err
		}

		workspaces = append(workspaces, workspace)
	}

	for _, projectFile := range projectFiles {
		workspaceContains := false
		for _, workspace := range workspaces {
			workspaceProjectFiles, err := workspace.projectPaths()
			if err != nil {
				return containers{}, err
			}

			if found := sliceutil.IsStringInSlice(projectFile, workspaceProjectFiles); found {
				workspaceContains = true
				break
			}
		}

		if !workspaceContains {
			project, err := newProject(projectFile)
			if err != nil {
				return containers{}, err
			}

			standaloneProjects = append(standaloneProjects, project)
		}
	}

	return containers{
		standaloneProjects: standaloneProjects,
		workspaces:         workspaces,
	}, nil
}

// FilterRelevantProjectFiles ...
func FilterRelevantProjectFiles(fileList []string, projectTypes ...XcodeProjectType) ([]string, error) {
	filters := []pathutil.FilterFunc{
		pathfilters.AllowXcodeProjExtFilter,
		pathfilters.AllowIsDirectoryFilter,
		pathfilters.ForbidEmbeddedWorkspaceRegexpFilter,
	}
	filters = append(filters, utility.CommonExcludeFilters()...)

	for _, projectType := range projectTypes {
		switch projectType {
		case XcodeProjectTypeIOS:
			filters = append(filters, pathfilters.AllowIphoneosSDKFilter)
		case XcodeProjectTypeMacOS:
			filters = append(filters, pathfilters.AllowMacosxSDKFilter)
		}
	}

	return pathutil.FilterPaths(fileList, filters...)
}

// FilterRelevantWorkspaceFiles ...
func FilterRelevantWorkspaceFiles(fileList []string, projectTypes ...XcodeProjectType) ([]string, error) {
	filters := []pathutil.FilterFunc{
		pathfilters.AllowXCWorkspaceExtFilter,
		pathfilters.AllowIsDirectoryFilter,
		pathfilters.AllowWorkspaceWithContentsFile,
		pathfilters.ForbidEmbeddedWorkspaceRegexpFilter,
	}
	filters = append(filters, utility.CommonExcludeFilters()...)

	for _, projectType := range projectTypes {
		switch projectType {
		case XcodeProjectTypeIOS:
			filters = append(filters, pathfilters.AllowIphoneosSDKFilter)
		case XcodeProjectTypeMacOS:
			filters = append(filters, pathfilters.AllowMacosxSDKFilter)
		}
	}

	return pathutil.FilterPaths(fileList, filters...)
}

// FilterRelevantPodfiles ...
func FilterRelevantPodfiles(fileList []string) ([]string, error) {
	filters := append(utility.CommonExcludeFilters(), AllowPodfileBaseFilter)
	return pathutil.FilterPaths(fileList, filters...)
}

// FilterRelevantCartFile ...
func FilterRelevantCartFile(fileList []string) ([]string, error) {
	filters := append(utility.CommonExcludeFilters(), AllowCartfileBaseFilter)
	return pathutil.FilterPaths(fileList, filters...)
}
