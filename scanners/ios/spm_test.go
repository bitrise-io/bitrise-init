package ios

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasSPMDependencies(t *testing.T) {
	tests := []struct {
		name     string
		fileList []string
		want     bool
	}{
		{
			name:     "Empty project",
			fileList: []string{},
			want:     false,
		},
		{
			name: "Project without SPM files",
			fileList: []string{
				"/Users/vagrant/git/README.md",
				"/Users/vagrant/git/App/AppDelegate.swift",
			},
			want: false,
		},
		{
			name: "Pure SPM project",
			fileList: []string{
				"/Users/vagrant/git/README.md",
				"/Users/vagrant/git/Package.swift",
				"/Users/vagrant/git/MyLib/MyLib.swift",
			},
			want: true,
		},
		{
			name: "Nested pure SPM project",
			fileList: []string{
				"/Users/vagrant/git/README.md",
				"/Users/vagrant/git/ios/Package.swift",
				"/Users/vagrant/git/ios/MyLib/MyLib.swift",
			},
			want: true,
		},
		{
			name: "Xcode project with SPM dependencies",
			fileList: []string{
				"/Users/vagrant/git/README.md",
				"/Users/vagrant/git/App/AppDelegate.swift",
				"/Users/vagrant/git/project.xcodeproj/project.xcworkspace/xcshareddata/swiftpm/Package.resolved",
			},
			want: true,
		},
		{
			name: "Xcode project without SPM dependencies",
			fileList: []string{
				"/Users/vagrant/git/README.md",
				"/Users/vagrant/git/App/AppDelegate.swift",
				"/Users/vagrant/git/BitriseTest.xcodeproj/project.xcworkspace/xcshareddata/IDEWorkspaceChecks.plist",
			},
			want: false,
		},
		{
			name: "Swift package descriptor in a vendored dependency folder",
			fileList: []string{
				"/Users/vagrant/git/README.md",
				"/Users/vagrant/git/App/AppDelegate.swift",
				"/Users/vagrant/git/Carthage/Checkouts/Lib/Package.swift",
			},
			want: false,
		},
		{
			name: "Lockfile in a vendored dependency folder",
			fileList: []string{
				"/Users/vagrant/git/README.md",
				"/Users/vagrant/git/App/AppDelegate.swift",
				"/Users/vagrant/git/Carthage/Checkouts/Lib/Package.resolved",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HasSPMDependencies(tt.fileList)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "HasSPMDependencies(%v)", tt.fileList)
		})
	}
}
