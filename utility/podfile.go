package utility

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"encoding/json"

	"github.com/bitrise-io/go-utils/fileutil"
)

/*
If one project exists in the Podfile's directory, workspace name will be the project's name.
If more then one project exists in the Podfile's directory root 'xcodeproj/project' property have to be defined in the Podfile.

Root 'xcodeproj/project' property will be mapped to the default cocoapods target (Pods).

If workspace property defined in the Podfile, it will override the workspace name.
*/

const getWorkspacePathGemfileContent = `source 'https://rubygems.org'
gem 'cocoapods-core'
`

// returns WORKSPACE_NAME.xcworkspace if user defined a workspace name
// returns empty struct {}, if no user defined workspace name exists in Podfile
const getWorkspacePathRubyScriptContent = `require 'cocoapods-core'
podfile_path = ENV['PODFILE_PATH']
podfile = Pod::Podfile.from_file(podfile_path)
puts podfile.workspace_path
`

const getTargetDefinitionProjectMappingGemfileContent = `source 'https://rubygems.org'
gem 'cocoapods-core'
gem 'json'
`

// returns target - project map, if xcodeproj defined in the Podfile
// return empty string if no xcodeproj defined in the Podfile
const getTargetDefinitionProjectMappingRubyScriptContent = `require 'cocoapods-core'
require 'json'

podfile_path = ENV['PODFILE_PATH']
podfile = Pod::Podfile.from_file(podfile_path)
targets = podfile.target_definitions
return '' unless targets

target_project_map = {}
targets.each do |name, target_definition|
  next unless target_definition.user_project_path
  target_project_map[name] = target_definition.user_project_path
end

puts target_project_map.to_json
`

func getTargetDefinitionProjectMap(podfilePth string) (map[string]string, error) {
	envs := []string{fmt.Sprintf("PODFILE_PATH=%s", podfilePth)}
	podfileDir := filepath.Dir(podfilePth)

	out, err := runRubyScriptForOutput(getTargetDefinitionProjectMappingRubyScriptContent, getTargetDefinitionProjectMappingGemfileContent, podfileDir, envs)
	if err != nil {
		return map[string]string{}, err
	}

	var targetDefinitionProjectMap map[string]string
	if err := json.Unmarshal([]byte(out), &targetDefinitionProjectMap); err != nil {
		return map[string]string{}, err
	}

	return targetDefinitionProjectMap, nil
}

func getUserDefinedProjectRelavtivePath(podfilePth string) (string, error) {
	targetProjectMap, err := getTargetDefinitionProjectMap(podfilePth)
	if err != nil {
		return "", err
	}

	for target, project := range targetProjectMap {
		if target == "Pods" {
			return project, nil
		}
	}
	return "", nil
}

func getUserDefinedWorkspaceRelativePath(podfilePth string) (string, error) {
	envs := []string{fmt.Sprintf("PODFILE_PATH=%s", podfilePth)}
	podfileDir := filepath.Dir(podfilePth)

	workspaceBase, err := runRubyScriptForOutput(getWorkspacePathRubyScriptContent, getWorkspacePathGemfileContent, podfileDir, envs)
	if err != nil {
		return "", err
	}

	return workspaceBase, nil
}

func filterXcodeProjectsInDirectory(projects []string, podDir string) []string {
	filtered := []string{}
	for _, project := range projects {
		dir := filepath.Dir(project)
		if dir == podDir {
			filtered = append(filtered, project)
		}
	}
	return filtered
}

// GetWorkspaceProjectMap ...
func GetWorkspaceProjectMap(podfilePth string, projects []string) (map[string]string, error) {
	// fix podfile quotation
	podfileContent, err := fileutil.ReadStringFromFile(podfilePth)
	if err != nil {
		return map[string]string{}, err
	}

	podfileContent = strings.Replace(podfileContent, `‘`, `'`, -1)
	podfileContent = strings.Replace(podfileContent, `’`, `'`, -1)
	podfileContent = strings.Replace(podfileContent, `“`, `"`, -1)
	podfileContent = strings.Replace(podfileContent, `”`, `"`, -1)

	if err := fileutil.WriteStringToFile(podfilePth, podfileContent); err != nil {
		return map[string]string{}, err
	}
	// ----

	podfileDir := filepath.Dir(podfilePth)

	projectRelPth, err := getUserDefinedProjectRelavtivePath(podfilePth)
	if err != nil {
		return map[string]string{}, err
	}

	if projectRelPth == "" {
		projects := filterXcodeProjectsInDirectory(projects, podfileDir)

		if len(projects) == 0 {
			return map[string]string{}, errors.New("failed to determin workspace - project mapping: no explicit project specified and no project found in the Podfile's directory")
		} else if len(projects) > 1 {
			return map[string]string{}, errors.New("failed to determin workspace - project mapping: no explicit project specified and more than one project found in the Podfile's directory")
		}

		projectRelPth = filepath.Base(projects[0])
	}
	projectPth := filepath.Join(podfileDir, projectRelPth)

	workspaceRelPth, err := getUserDefinedWorkspaceRelativePath(podfilePth)
	if err != nil {
		return map[string]string{}, err
	}

	if workspaceRelPth == "" {
		projectName := filepath.Base(strings.TrimSuffix(projectPth, ".xcodeproj"))
		workspaceRelPth = projectName + ".xcworkspace"
	}
	workspacePth := filepath.Join(podfileDir, workspaceRelPth)

	return map[string]string{
		workspacePth: projectPth,
	}, nil
}
