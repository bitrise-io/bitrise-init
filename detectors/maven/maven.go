package maven

import (
	"fmt"

	"github.com/bitrise-io/bitrise-init/detectors/direntry"
)

type Project struct {
	RootDirEntry                direntry.DirEntry
	ProjectObjectModelFileEntry direntry.DirEntry
	MavenWrapperFileEntry       direntry.DirEntry
}

func ScanProject(searchDir string) (*Project, error) {
	rootEntry, err := direntry.WalkDir(searchDir, 6)
	if err != nil {
		return nil, err
	}
	return detectMavenProjectRoot(*rootEntry)
}

func detectMavenProjectRoot(searchDir direntry.DirEntry) (*Project, error) {
	projectObjectModelEntry := searchDir.FindFirstEntryByName("pom.xml", false)
	if projectObjectModelEntry == nil {
		return nil, nil
	}

	projectRootDirEntry := projectObjectModelEntry.Parent()
	if projectRootDirEntry == nil {
		return nil, fmt.Errorf("unable to detect project root")
	}

	mavenWrapperEntry := projectRootDirEntry.FindImmediateChildByName("mvnw", false)
	if mavenWrapperEntry == nil {
		return nil, nil
	}

	return &Project{
		RootDirEntry:                *projectRootDirEntry,
		ProjectObjectModelFileEntry: *projectObjectModelEntry,
		MavenWrapperFileEntry:       *mavenWrapperEntry,
	}, nil
}
