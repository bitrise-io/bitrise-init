package kmp

import (
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type DirEntry struct {
	Path    string
	Name    string
	IsDir   bool
	Entries []DirEntry
}

var ignoreDirs = []string{".git", ".github", ".gradle", ".idea", "build", ".kotlin", ".fleet", "CordovaLib", "node_modules"}

func ListEntries(rootDir string, depth uint) (*DirEntry, error) {
	if depth == 0 {
		return nil, nil
	}

	root := DirEntry{
		Path:    rootDir,
		Name:    filepath.Base(rootDir),
		IsDir:   true,
		Entries: nil,
	}

	if err := recursiveListEntries(&root, 0, depth); err != nil {
		return nil, err
	}

	return &root, nil
}

func recursiveListEntries(root *DirEntry, currentDepth, maxDepth uint) error {
	if currentDepth >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(root.Path)
	if err != nil {
		// TODO: log error
		return nil
	}

	root.Entries = make([]DirEntry, 0, len(entries))
	for _, entry := range entries {
		if slices.Contains(ignoreDirs, entry.Name()) {
			continue
		}

		dirEntry := DirEntry{
			Path:    filepath.Join(root.Path, entry.Name()),
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Entries: nil,
		}

		if dirEntry.IsDir {
			if err := recursiveListEntries(&dirEntry, currentDepth+1, maxDepth); err != nil {
				return err
			}
		}

		root.Entries = append(root.Entries, dirEntry)
	}

	return nil
}

func listDirEntries(root string, depth uint) ([]DirEntry, error) {
	var entries []DirEntry
	dirsToRead := []string{root}
	for i := 0; i < int(depth); i++ {
		var nextDirsToRead []string
		for _, dir := range dirsToRead {
			dirEntries, err := os.ReadDir(dir)
			if err != nil {
				// log.Warnf("Failed to read dir: %s", dir)
				continue
			}

			for _, entry := range dirEntries {
				if entry.IsDir() {
					if slices.Contains(ignoreDirs, entry.Name()) {
						continue
					}

					nextDirsToRead = append(nextDirsToRead, filepath.Join(dir, entry.Name()))
				}
				entries = append(entries, DirEntry{Path: path.Join(dir, entry.Name()), Name: entry.Name(), IsDir: entry.IsDir()})
			}
		}
		if len(nextDirsToRead) == 0 {
			break
		}
		dirsToRead = nextDirsToRead
	}

	slices.SortFunc(entries, func(a, b DirEntry) int {
		componentsA := strings.Split(a.Path, string(os.PathSeparator))
		componentsB := strings.Split(b.Path, string(os.PathSeparator))
		if len(componentsA) < len(componentsB) {
			return -1
		} else if len(componentsA) > len(componentsB) {
			return 1
		}
		for i := 0; i < len(componentsA); i++ {
			if componentsA[i] < componentsB[i] {
				return -1
			} else if componentsA[i] > componentsB[i] {
				return 1
			}
		}
		return 0
	})
	return entries, nil
}
