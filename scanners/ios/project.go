package ios

import (
	"fmt"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/v2/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/v2/xcodeproject/xcscheme"
	"github.com/bitrise-io/go-xcode/v2/xcodeproject/xcworkspace"
)

type container interface {
	isWorkspace() bool
	// path returns container absolute path
	path() string
	// schemes returns schemes mapped to the project or workspace path
	schemes() (map[string][]xcscheme.Scheme, error)
	projects() ([]*xcodeproj.XcodeProj, []string, error)
	projectPaths() ([]string, error)
}

// projectContainer ...
type projectContainer struct {
	project *xcodeproj.XcodeProj
}

func newProject(path string, projectFactory xcodeproj.Factory) (projectContainer, error) {
	if !xcodeproj.IsXcodeProj(path) {
		return projectContainer{}, fmt.Errorf("path (%s) is not a Project", path)
	}

	project, err := projectFactory.Open(path)
	if err != nil {
		return projectContainer{}, fmt.Errorf("failed to open Project (%s): %w", path, err)
	}

	return projectContainer{
		project: project,
	}, nil
}

func (p projectContainer) isWorkspace() bool {
	return false
}

func (p projectContainer) path() string {
	return p.project.Path
}

func (p projectContainer) schemes() (map[string][]xcscheme.Scheme, error) {
	projectSchemes, err := p.project.Schemes()
	if err != nil {
		return nil, fmt.Errorf("failed to list Schemes in Project (%s): %w", p.project.Path, err)
	}

	containerToSchemes := make(map[string][]xcscheme.Scheme)
	containerToSchemes[p.project.Path] = projectSchemes

	return containerToSchemes, nil
}

func (p projectContainer) projects() ([]*xcodeproj.XcodeProj, []string, error) {
	return []*xcodeproj.XcodeProj{p.project}, []string{}, nil
}

func (p projectContainer) projectPaths() ([]string, error) {
	return []string{p.project.Path}, nil
}

// workspaceContainer ...
type workspaceContainer struct {
	workspace      *xcworkspace.Workspace
	projectFactory xcodeproj.Factory
}

func newWorkspace(path string, workspaceFactory xcworkspace.Factory, projectFactory xcodeproj.Factory) (workspaceContainer, error) {
	if !xcworkspace.IsWorkspace(path) {
		return workspaceContainer{}, fmt.Errorf("path (%s) is not a Workspace", path)
	}

	workspace, err := workspaceFactory.Open(path)
	if err != nil {
		return workspaceContainer{}, fmt.Errorf("failed to open Workspace (%s): %w", path, err)
	}

	return workspaceContainer{
		workspace:      workspace,
		projectFactory: projectFactory,
	}, nil
}

func (w workspaceContainer) isWorkspace() bool {
	return true
}

func (w workspaceContainer) path() string {
	return w.workspace.Path
}

func (w workspaceContainer) schemes() (map[string][]xcscheme.Scheme, error) {
	containerToSchemes, err := w.workspace.Schemes()
	if err != nil {
		return nil, fmt.Errorf("failed to list Schemes in Workspace (%s): %w", w.workspace.Path, err)
	}

	return containerToSchemes, nil
}

func (w workspaceContainer) projects() ([]*xcodeproj.XcodeProj, []string, error) {
	projPaths, err := w.workspace.ProjectFileLocations()
	if err != nil {
		return nil, nil, err
	}

	var projects []*xcodeproj.XcodeProj
	var missingProjects []string
	for _, projPath := range projPaths {
		if exist, err := pathutil.IsPathExists(projPath); err != nil {
			return nil, nil, fmt.Errorf("failed to list Projects in the Workspace (%s), can not check if path (%s) exists: %w", w.workspace.Path, projPath, err)
		} else if !exist {
			missingProjects = append(missingProjects, projPath)
			continue
		}

		project, err := w.projectFactory.Open(projPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open Project (%s) in the Workspace (%s): %w", projPath, w.workspace.Path, err)
		}

		projects = append(projects, project)
	}

	return projects, missingProjects, nil
}

func (w workspaceContainer) projectPaths() ([]string, error) {
	return w.workspace.ProjectFileLocations()
}
