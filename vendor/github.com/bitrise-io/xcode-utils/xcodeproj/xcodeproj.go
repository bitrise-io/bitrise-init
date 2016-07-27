package xcodeproj

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Extensions
const (
	// XCWorkspaceExt ...
	XCWorkspaceExt = ".xcworkspace"
	// XCodeProjExt ...
	XCodeProjExt = ".xcodeproj"
	// XCSchemeExt ...
	XCSchemeExt = ".xcscheme"
)

// Path Components
const (
	XCSharedData = "xcshareddata"
	XCSchemes    = "xcschemes"
)

// ------------------------------
// Schemes

// IsSharedSchemeFilePath ...
func IsSharedSchemeFilePath(pth string) bool {
	regexpPattern := filepath.Join(".*[/]?xcshareddata", "xcschemes", ".+[.]xcscheme")
	regexp := regexp.MustCompile(regexpPattern)
	return (regexp.FindString(pth) != "")
}

// FilterSharedSchemeFilePaths ...
func FilterSharedSchemeFilePaths(paths []string) []string {
	filteredPaths := []string{}
	for _, pth := range paths {
		if IsSharedSchemeFilePath(pth) {
			filteredPaths = append(filteredPaths, pth)
		}
	}
	return filteredPaths
}

func sharedSchemeFilePaths(projectOrWorkspacePth string) ([]string, error) {
	paths, err := filesInDir(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}
	return FilterSharedSchemeFilePaths(paths), nil
}

// ProjectSharedSchemeFilePaths ...
func ProjectSharedSchemeFilePaths(projectPth string) ([]string, error) {
	return sharedSchemeFilePaths(projectPth)
}

// WorkspaceSharedSchemeFilePaths ...
func WorkspaceSharedSchemeFilePaths(workspacePth string) ([]string, error) {
	workspaceSchemeFilePaths, err := sharedSchemeFilePaths(workspacePth)
	if err != nil {
		return []string{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSchemeFilePaths, err := sharedSchemeFilePaths(project)
		if err != nil {
			return []string{}, err
		}
		workspaceSchemeFilePaths = append(workspaceSchemeFilePaths, projectSchemeFilePaths...)
	}
	return workspaceSchemeFilePaths, nil
}

// SchemeNameFromPath ...
func SchemeNameFromPath(schemePth string) string {
	basename := filepath.Base(schemePth)
	ext := filepath.Ext(schemePth)
	if ext != XCSchemeExt {
		return ""
	}
	return strings.TrimSuffix(basename, ext)
}

func sharedSchemes(projectOrWorkspacePth string) ([]string, error) {
	schemePaths, err := sharedSchemeFilePaths(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}

	schemes := []string{}
	for _, schemePth := range schemePaths {
		schemes = append(schemes, SchemeNameFromPath(schemePth))
	}
	return schemes, nil
}

// ProjectSharedSchemes ...
func ProjectSharedSchemes(projectPth string) ([]string, error) {
	return sharedSchemes(projectPth)
}

// WorkspaceSharedSchemes ...
func WorkspaceSharedSchemes(workspacePth string) ([]string, error) {
	workspaceSchemes, err := sharedSchemes(workspacePth)
	if err != nil {
		return []string{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSchemes, err := sharedSchemes(project)
		if err != nil {
			return []string{}, err
		}
		workspaceSchemes = append(workspaceSchemes, projectSchemes...)
	}
	return workspaceSchemes, nil
}

// IsUserSchemeFilePath ...
func IsUserSchemeFilePath(pth string) bool {
	regexpPattern := filepath.Join(".*[/]?xcuserdata", ".*[.]xcuserdatad", "xcschemes", ".+[.]xcscheme")
	regexp := regexp.MustCompile(regexpPattern)
	return (regexp.FindString(pth) != "")
}

// FilterUserSchemeFilePaths ...
func FilterUserSchemeFilePaths(paths []string) []string {
	filteredPaths := []string{}
	for _, pth := range paths {
		if IsUserSchemeFilePath(pth) {
			filteredPaths = append(filteredPaths, pth)
		}
	}
	return filteredPaths
}

// UserSchemeFilePaths ...
func UserSchemeFilePaths(projectOrWorkspacePth string) ([]string, error) {
	paths, err := filesInDir(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}
	return FilterUserSchemeFilePaths(paths), nil
}

func userSchemes(projectOrWorkspacePth string) ([]string, error) {
	schemePaths, err := UserSchemeFilePaths(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}

	schemes := []string{}
	for _, schemePth := range schemePaths {
		schemes = append(schemes, SchemeNameFromPath(schemePth))
	}
	return schemes, nil
}

// ProjectUserSchemes ...
func ProjectUserSchemes(projectPth string) ([]string, error) {
	return userSchemes(projectPth)
}

// WorkspaceUserSchemes ...
func WorkspaceUserSchemes(workspacePth string) ([]string, error) {
	workspaceSchemes, err := userSchemes(workspacePth)
	if err != nil {
		return []string{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSchemes, err := userSchemes(project)
		if err != nil {
			return []string{}, err
		}
		workspaceSchemes = append(workspaceSchemes, projectSchemes...)
	}
	return workspaceSchemes, nil
}

// SchemeFileContentContainsXCTestBuildAction ...
func SchemeFileContentContainsXCTestBuildAction(schemeFileContent string) (bool, error) {
	regexpPattern := `BuildableName = ".+.xctest"`
	regexp := regexp.MustCompile(regexpPattern)

	scanner := bufio.NewScanner(strings.NewReader(schemeFileContent))
	for scanner.Scan() {
		line := scanner.Text()
		if regexp.FindString(line) != "" {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// SchemeFileContainsXCTestBuildAction ...
func SchemeFileContainsXCTestBuildAction(schemeFilePth string) (bool, error) {
	content, err := fileutil.ReadStringFromFile(schemeFilePth)
	if err != nil {
		return false, err
	}

	return SchemeFileContentContainsXCTestBuildAction(content)
}

// ReCreateProjectUserSchemes ...
func ReCreateProjectUserSchemes(projectPth string) error {
	rubyScriptContent := `require 'xcodeproj'
require 'json'

project_path = ENV['project_path']

begin
  raise 'empty path' if project_path.empty?

  project = Xcodeproj::Project.open(project_path)
  project.recreate_user_schemes
  project.save
rescue => ex
  puts(ex.inspect.to_s)
  puts('--- Stack trace: ---')
  puts(ex.backtrace.to_s)
  exit(1)
end
`

	tmpDir, err := pathutil.NormalizedOSTempDirPath("bitrise")
	if err != nil {
		return err
	}

	rubyScriptPth := path.Join(tmpDir, "recreate_user_schemes.rb")
	if err := fileutil.WriteStringToFile(rubyScriptPth, rubyScriptContent); err != nil {
		return err
	}

	projectBase := filepath.Base(projectPth)
	envs := append(os.Environ(), "project_path="+projectBase, "LC_ALL=en_US.UTF-8")
	projectDir := filepath.Dir(projectPth)

	out, err := cmdex.NewCommand("ruby", rubyScriptPth).SetDir(projectDir).SetEnvs(envs).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) && out != "" {
			return errors.New(out)
		}
		return err
	}

	return nil
}

// ReCreateWorkspaceUserSchemes ...
func ReCreateWorkspaceUserSchemes(workspace string) error {
	projects, err := WorkspaceProjectReferences(workspace)
	if err != nil {
		return err
	}

	for _, project := range projects {
		if err := ReCreateProjectUserSchemes(project); err != nil {
			return err
		}
	}

	return nil
}

// ------------------------------

// ------------------------------
// Project

// IsXCodeProj ...
func IsXCodeProj(pth string) bool {
	return strings.HasSuffix(pth, XCodeProjExt)
}

// ------------------------------
// Workspace

// IsXCWorkspace ...
func IsXCWorkspace(pth string) bool {
	return strings.HasSuffix(pth, XCWorkspaceExt)
}

// WorkspaceProjectReferences ...
func WorkspaceProjectReferences(workspace string) ([]string, error) {
	projects := []string{}

	workspaceDir := filepath.Dir(workspace)

	xcworkspacedataPth := path.Join(workspace, "contents.xcworkspacedata")
	if exist, err := pathutil.IsPathExists(xcworkspacedataPth); err != nil {
		return []string{}, err
	} else if !exist {
		return []string{}, fmt.Errorf("contents.xcworkspacedata does not exist at: %s", xcworkspacedataPth)
	}

	xcworkspacedataStr, err := fileutil.ReadStringFromFile(xcworkspacedataPth)
	if err != nil {
		return []string{}, err
	}

	xcworkspacedataLines := strings.Split(xcworkspacedataStr, "\n")
	fileRefStart := false
	regexp := regexp.MustCompile(`location = "(.+):(.+).xcodeproj"`)

	for _, line := range xcworkspacedataLines {
		if strings.Contains(line, "<FileRef") {
			fileRefStart = true
			continue
		}

		if fileRefStart {
			fileRefStart = false
			matches := regexp.FindStringSubmatch(line)
			if len(matches) == 3 {
				projectName := matches[2]
				project := filepath.Join(workspaceDir, projectName+".xcodeproj")
				projects = append(projects, project)
			}
		}
	}

	return projects, nil
}

// ------------------------------

// ------------------------------
// Utility

func filesInDir(dir string) ([]string, error) {
	files := []string{}
	if err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	}); err != nil {
		return []string{}, err
	}
	return files, nil
}
