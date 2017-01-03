package utility

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestAllowPodfileBaseFilter(t *testing.T) {
	t.Log("abs path")
	{
		absPaths := []string{
			"/Users/bitrise/Test.txt",
			"/Users/bitrise/.git/Podfile",
			"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		expectedPaths := []string{
			"/Users/bitrise/.git/Podfile",
			"/Users/bitrise/sample-apps-ios-cocoapods/Pods/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"/Users/bitrise/ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		actualPaths, err := FilterPaths(absPaths, AllowPodfileBaseFilter)
		require.NoError(t, err)
		require.Equal(t, expectedPaths, actualPaths)
	}

	t.Log("rel path")
	{
		relPaths := []string{
			"Test.txt",
			".git/Podfile",
			"sample-apps-ios-cocoapods/Pods/Podfile",
			"ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		expectedPaths := []string{
			".git/Podfile",
			"sample-apps-ios-cocoapods/Pods/Podfile",
			"ios-no-shared-schemes/Carthage/Checkouts/Result/Podfile",
			"ios-no-shared-schemes/test.framework/Checkouts/Result/Podfile",
		}

		actualPaths, err := FilterPaths(relPaths, AllowPodfileBaseFilter)
		require.NoError(t, err)
		require.Equal(t, expectedPaths, actualPaths)
	}
}

func TestGetTargetDefinitionProjectMap(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("xcodeproj defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
xcodeproj 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedTargetDefinition := map[string]string{
			"Pods": "MyXcodeProject",
		}
		actualTargetDefinition, err := getTargetDefinitionProjectMap(podfilePth)
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}

	t.Log("xcodeproj NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedTargetDefinition := map[string]string{}
		actualTargetDefinition, err := getTargetDefinitionProjectMap(podfilePth)
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}
}

func TestGetUserDefinedProjectRelavtivePath(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("xcodeproj defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_defined")

		podfile := `platform :ios, '9.0'
xcodeproj 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedProject := "MyXcodeProject"
		actualProject, err := getUserDefinedProjectRelavtivePath(podfilePth)
		require.NoError(t, err)
		require.Equal(t, expectedProject, actualProject)
	}

	t.Log("xcodeproj NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedProject := ""
		actualProject, err := getUserDefinedProjectRelavtivePath(podfilePth)
		require.NoError(t, err)
		require.Equal(t, expectedProject, actualProject)
	}
}

func TestGetUserDefinedWorkspaceRelativePath(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("workspace defined")
	{
		tmpDir = filepath.Join(tmpDir, "workspace_defined")

		podfile := `platform :ios, '9.0'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedWorkspace := "MyWorkspace.xcworkspace"
		actualWorkspace, err := getUserDefinedWorkspaceRelativePath(podfilePth)
		require.NoError(t, err)
		require.Equal(t, expectedWorkspace, actualWorkspace)
	}

	t.Log("workspace NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "workspace_not_defined")

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		expectedWorkspace := ""
		actualWorkspace, err := getUserDefinedWorkspaceRelativePath(podfilePth)
		require.NoError(t, err)
		require.Equal(t, expectedWorkspace, actualWorkspace)
	}
}

func TestGetWorkspaceProjectMap(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("0 project in Podfile's dir")
	{
		tmpDir = filepath.Join(tmpDir, "no_project")
		require.NoError(t, err)

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir")
	{
		tmpDir = filepath.Join(tmpDir, "one_project")
		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{projectPth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "project", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("Multiple project in Podfile's dir")
	{
		tmpDir = filepath.Join(tmpDir, "multiple_project")
		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{project1Pth, project2Pth})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("0 project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "no_project_project_defined")
		require.NoError(t, err)

		podfile := `platform :ios, '9.0'
xcodeproj 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "one_project_project_defined")
		podfile := `platform :ios, '9.0'
xcodeproj 'project'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{projectPth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "project", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("Multiple project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "multiple_project")
		podfile := `platform :ios, '9.0'
xcodeproj 'project1'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{project1Pth, project2Pth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "project1", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project1", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir + workspace defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "one_project")
		podfile := `platform :ios, '9.0'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{projectPth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "MyWorkspace", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("Multiple project in Podfile's dir + workspace defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "multiple_project")
		podfile := `platform :ios, '9.0'
xcodeproj 'project1'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := GetWorkspaceProjectMap(podfilePth, []string{project1Pth, project2Pth})
		require.NoError(t, err)
		require.Equal(t, 1, len(workspaceProjectMap))

		for workspace, project := range workspaceProjectMap {
			workspaceBasename := filepath.Base(workspace)
			workspaceName := strings.TrimSuffix(workspaceBasename, ".xcworkspace")

			projectBasename := filepath.Base(project)
			projectName := strings.TrimSuffix(projectBasename, ".xcodeproj")

			require.Equal(t, "MyWorkspace", workspaceName, fmt.Sprintf("%v", workspaceProjectMap))
			require.Equal(t, "project1", projectName, fmt.Sprintf("%v", workspaceProjectMap))
		}

		require.NoError(t, os.RemoveAll(tmpDir))
	}
}
