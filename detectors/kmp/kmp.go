package kmp

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/bitrise-init/detectors/direntry"
	"github.com/bitrise-io/bitrise-init/detectors/gradle"
	"github.com/bitrise-io/go-utils/log"
)

type Project struct {
	GradleProject gradle.Project

	XcodeprojFile     *direntry.DirEntry
	AndroidAppDir     *direntry.DirEntry
	AndroidWearAppDir *direntry.DirEntry
}

func ScanProject(gradleProject gradle.Project) (*Project, error) {
	log.TInfof("Searching for Kotlin Multiplatform dependencies...")
	kotlinMultiplatformDetected, err := gradleProject.DetectAnyDependencies([]string{
		"org.jetbrains.kotlin.multiplatform",
		`kotlin("multiplatform")`,
	})
	if err != nil {
		return nil, err
	}

	log.TDonef("Kotlin Multiplatform dependencies found: %v", kotlinMultiplatformDetected)
	if !kotlinMultiplatformDetected {
		return nil, nil
	}

	log.TInfof("Scanning Kotlin Multiplatform targets...")
	xcodeProjectFile := gradleProject.RootDirEntry.FindFirstFileEntryByExtension(".xcodeproj")
	if xcodeProjectFile != nil {
		log.TPrintf("iOS App target: %s", xcodeProjectFile.RelPath)
	}

	androidApplicationPluginID, err := gradleProject.GetDependencyID(`com.android.application`)
	if err != nil {
		return nil, fmt.Errorf("failed to get Android application plugin ID: %w", err)
	}

	androidAppDependencies := []string{
		`"com.android.application"`,
	}
	if androidApplicationPluginID != "" {
		// alias(libs.plugins.androidApplication)
		// TODO?: make it more reliable by using a regex/glob matching (alias.*<androidApplicationPluginID>)
		androidAppDependencies = append(androidAppDependencies, fmt.Sprintf("alias(libs.plugins.%s)", androidApplicationPluginID))
	}

	androidProjects, err := gradleProject.FindSubProjectsWithAnyDependencies(androidAppDependencies)
	if err != nil {
		return nil, err
	}

	// Wear projects Manifest files contains this: <uses-feature android:name="android.hardware.type.watch" />
	var androidAppDirs []direntry.DirEntry
	var androidWearAppDirs []direntry.DirEntry
	if len(androidProjects) > 0 {
		for _, androidProject := range androidProjects {
			androidProjectDir := androidProject.BuildScriptFileEntry.Parent()
			manifestFiles := androidProjectDir.FindAllEntriesByName("AndroidManifest.xml", false)
			isWearApp := false
			if len(manifestFiles) > 0 {
				for _, manifestFile := range manifestFiles {
					manifestContent, err := os.ReadFile(manifestFile.AbsPath)
					if err != nil {
						return nil, fmt.Errorf("failed to read AndroidManifest.xml file: %w", err)
					}
					if strings.Contains(string(manifestContent), "android.hardware.type.watch") {
						isWearApp = true
						break
					}
				}
			}

			if isWearApp {
				androidWearAppDirs = append(androidWearAppDirs, *androidProjectDir)
			} else {
				androidAppDirs = append(androidAppDirs, *androidProjectDir)
			}
		}
	}

	var androidAppDir *direntry.DirEntry
	if len(androidAppDirs) > 0 {
		androidAppDir = &androidAppDirs[0]
		if len(androidAppDirs) > 1 {
			log.TWarnf("%d Android targets found in the Gradle project, using the first one: %s", len(androidAppDirs), androidAppDir.RelPath)
		} else {
			log.TPrintf("Android App target: %s", androidAppDir.RelPath)
		}
	}

	var androidWearAppDir *direntry.DirEntry
	if len(androidWearAppDirs) > 0 {
		androidWearAppDir = &androidWearAppDirs[0]
		if len(androidWearAppDirs) > 1 {
			log.TWarnf("%d Android Wear targets found in the Gradle project, using the first one: %s", len(androidWearAppDirs), androidWearAppDir.RelPath)
		} else {
			log.TPrintf("Android Wear target: %s", androidWearAppDir.RelPath)
		}
	}

	return &Project{
		GradleProject:     gradleProject,
		XcodeprojFile:     xcodeProjectFile,
		AndroidAppDir:     androidAppDir,
		AndroidWearAppDir: androidWearAppDir,
	}, nil
}
