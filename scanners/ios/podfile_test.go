package ios

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcodeproj"
	"github.com/bitrise-io/go-xcode/xcodeproject/xcscheme"
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

		actualPaths, err := pathutil.FilterPaths(absPaths, AllowPodfileBaseFilter)
		require.NoError(t, err)
		require.Equal(t, expectedPaths, actualPaths)
	}

	t.Log("rel path")
	{
		relPaths := []string{
			".",
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

		actualPaths, err := pathutil.FilterPaths(relPaths, AllowPodfileBaseFilter)
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
project 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		expectedTargetDefinition := map[string]string{
			"Pods": "MyXcodeProject.xcodeproj",
		}

		actualTargetDefinition, err := podparser.getTargetDefinitionProjectMap("")
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}

	t.Log("xcodeproj NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		expectedTargetDefinition := map[string]string{}
		actualTargetDefinition, err := podparser.getTargetDefinitionProjectMap("")
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}

	t.Log("cocoapods 1.8.4")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `source 'https://github.com/CocoaPods/Specs.git'
platform :ios, '8.0'

# pod 'Functional.m', '~> 1.0'

# Add Kiwi as a dependency for the Test target
target :SampleAppWithCocoapodsTests do
  pod 'Kiwi'
end

# post_install do |installer_representation|
#   installer_representation.project.targets.each do |target|
#     target.build_configurations.each do |config|
#       config.build_settings['ONLY_ACTIVE_ARCH'] = 'NO'
#     end
#   end
# end`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		expectedTargetDefinition := map[string]string{}
		actualTargetDefinition, err := podparser.getTargetDefinitionProjectMap("1.8.4")
		require.NoError(t, err)
		require.Equal(t, expectedTargetDefinition, actualTargetDefinition)
	}
}

func TestGetUserDefinedProjectAbsPath(t *testing.T) {
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
project 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		expectedProject := filepath.Join(tmpDir, "MyXcodeProject.xcodeproj")
		actualProject, err := podparser.getUserDefinedProjectAbsPath("")
		require.NoError(t, err)
		require.Equal(t, expectedProject, actualProject)
	}

	t.Log("xcodeproj NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "xcodeproj_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		expectedProject := ""
		actualProject, err := podparser.getUserDefinedProjectAbsPath("")
		require.NoError(t, err)
		require.Equal(t, expectedProject, actualProject)
	}
}

func TestGetUserDefinedWorkspaceAbsPath(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__utility_test__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	t.Log("workspace defined")
	{
		tmpDir = filepath.Join(tmpDir, "workspace_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		expectedWorkspace := filepath.Join(tmpDir, "MyWorkspace.xcworkspace")
		actualWorkspace, err := podparser.getUserDefinedWorkspaceAbsPath("")
		require.NoError(t, err)
		require.Equal(t, expectedWorkspace, actualWorkspace)
	}

	t.Log("workspace NOT defined")
	{
		tmpDir = filepath.Join(tmpDir, "workspace_not_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		expectedWorkspace := ""
		actualWorkspace, err := podparser.getUserDefinedWorkspaceAbsPath("")
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
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir")
	{
		tmpDir = filepath.Join(tmpDir, "one_project")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{projectPth})
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
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{project1Pth, project2Pth})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("0 project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "no_project_project_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'MyXcodeProject'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{})
		require.Error(t, err)
		require.Equal(t, 0, len(workspaceProjectMap))

		require.NoError(t, os.RemoveAll(tmpDir))
	}

	t.Log("1 project in Podfile's dir + project defined in Podfile")
	{
		tmpDir = filepath.Join(tmpDir, "one_project_project_defined")
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'project'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{projectPth})
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
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'project1'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{project1Pth, project2Pth})
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
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		project := ""
		projectPth := filepath.Join(tmpDir, "project.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(projectPth, project))

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{projectPth})
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
		require.NoError(t, os.MkdirAll(tmpDir, 0777))

		podfile := `platform :ios, '9.0'
project 'project1'
workspace 'MyWorkspace'
pod 'Alamofire', '~> 3.4'
`
		podfilePth := filepath.Join(tmpDir, "Podfile")
		require.NoError(t, fileutil.WriteStringToFile(podfilePth, podfile))
		podparser := podfileParser{podfilePth: podfilePth}

		project1 := ""
		project1Pth := filepath.Join(tmpDir, "project1.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project1Pth, project1))

		project2 := ""
		project2Pth := filepath.Join(tmpDir, "project2.xcodeproj")
		require.NoError(t, fileutil.WriteStringToFile(project2Pth, project2))

		workspaceProjectMap, err := podparser.GetWorkspaceProjectMap([]string{project1Pth, project2Pth})
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

type mockContainer struct {
	containerPath string
	innerProjects []string
}

func (m mockContainer) path() string {
	return m.containerPath
}

func (m mockContainer) isWorkspace() bool {
	panic("not implemented")
}

func (m mockContainer) projectPaths() ([]string, error) {
	return m.innerProjects, nil
}

func (m mockContainer) projects() ([]xcodeproj.XcodeProj, []string, error) {
	panic("not implemented")
}

func (m mockContainer) schemes() (map[string][]xcscheme.Scheme, error) {
	panic("not implemented")
}

func TestMergePodWorkspaceProjectMap(t *testing.T) {
	t.Log("workspace is in the repository")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}
		detected := containers{
			workspaces: []container{
				mockContainer{
					containerPath: "MyWorkspace.xcworkspace",
					innerProjects: []string{"MyXcodeProject.xcodeproj"},
				},
			},
		}

		wantDetected := containers{
			workspaces: []container{
				mockContainer{
					containerPath: "MyWorkspace.xcworkspace",
					innerProjects: []string{"MyXcodeProject.xcodeproj"},
				},
			},
			podWorkspacePaths: []string{"MyWorkspace.xcworkspace"},
		}

		gotDetected, err := mergePodWorkspaceProjectMap(podWorkspaceMap, detected)
		require.NoError(t, err)
		require.Equal(t, wantDetected, gotDetected)
	}

	t.Log("workspace is in the repository, but project not attached - ERROR")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		detected := containers{
			workspaces: []container{
				mockContainer{
					containerPath: "MyWorkspace.xcworkspace",
				},
			},
		}

		gotDetected, err := mergePodWorkspaceProjectMap(podWorkspaceMap, detected)
		require.Error(t, err)
		require.Equal(t, containers{}, gotDetected)
	}

	t.Log("workspace is in the repository, but project is also standalone - ERROR")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		detected := containers{
			standaloneProjects: []container{
				mockContainer{
					containerPath: "MyXcodeProject.xcodeproj",
					innerProjects: []string{"MyXcodeProject.xcodeproj"},
				},
			},
			workspaces: []container{
				mockContainer{
					containerPath: "MyWorkspace.xcworkspace",
				},
			},
		}

		gotDetected, err := mergePodWorkspaceProjectMap(podWorkspaceMap, detected)
		require.Error(t, err)
		require.Equal(t, containers{}, gotDetected)
	}

	t.Log("workspace is gitignored")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}
		detected := containers{
			standaloneProjects: []container{
				mockContainer{
					containerPath: "MyXcodeProject.xcodeproj",
					innerProjects: []string{"MyXcodeProject.xcodeproj"},
				},
			},
		}

		want := containers{
			standaloneProjects: []container{},
			workspaces: []container{
				podWorkspace{
					workspacePath: "MyWorkspace.xcworkspace",
					workspaceProjects: []container{
						mockContainer{
							containerPath: "MyXcodeProject.xcodeproj",
							innerProjects: []string{"MyXcodeProject.xcodeproj"},
						},
					},
				},
			},
			podWorkspacePaths: []string{"MyWorkspace.xcworkspace"},
		}

		got, err := mergePodWorkspaceProjectMap(podWorkspaceMap, detected)
		require.NoError(t, err)
		require.Equal(t, want, got)
	}

	t.Log("workspace is gitignored, but standalone project missing - ERROR")
	{
		podWorkspaceMap := map[string]string{
			"MyWorkspace.xcworkspace": "MyXcodeProject.xcodeproj",
		}

		detected := containers{}

		got, err := mergePodWorkspaceProjectMap(podWorkspaceMap, detected)
		require.Error(t, err)
		require.Equal(t, containers{}, got)
	}
}
