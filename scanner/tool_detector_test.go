package scanner

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_kotlinMultiplatformDetector(t *testing.T) {
	tests := []struct {
		name     string
		rootPath string
		want     DetectionResult
		wantErr  bool
	}{
		{
			name:     "Empty project",
			rootPath: t.TempDir(),
			want:     DetectionResult{Detected: false, ProjectTree: ""},
		},
		{
			name:     "Android project",
			rootPath: createAndroidProjectFiles(t, t.TempDir()),
			want: DetectionResult{Detected: false, ProjectTree: `app/
· build.gradle
build.gradle
settings.gradle
`},
		},
		{
			name:     "Kotlin Multiplatform project",
			rootPath: createKotlinMultiplatformFiles(t, t.TempDir()),
			want: DetectionResult{Detected: true, ProjectTree: `app/
· build.gradle
build.gradle
settings.gradle
shared/
· build.gradle.kts
`},
		},
		{
			name:     "Nested Kotlin Multiplatform project",
			rootPath: createNestedKotlinMultiplatformFiles(t),
			want: DetectionResult{Detected: true, ProjectTree: `my-project/
· mobile/
· · android/
· · · app/
· · · · build.gradle
· · · build.gradle
· · · settings.gradle
· · · shared/
· · · · build.gradle.kts
`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := kotlinMultiplatformDetector{}
			got, err := d.DetectToolIn(tt.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectToolIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectToolIn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toolDetector(t *testing.T) {
	tests := []struct {
		name         string
		toolDetector toolDetector
		rootPath     string
		want         DetectionResult
		wantErr      bool
	}{
		{
			name: "Primary file or optional files in empty project",
			toolDetector: toolDetector{
				toolName:      "Bazel",
				primaryFile:   "WORKSPACE",
				optionalFiles: []string{"BUILD.bazel", ".bazelrc"},
			},
			rootPath: t.TempDir(),
			want:     DetectionResult{Detected: false, ProjectTree: ""},
		},
		{
			name: "Primary file in Tuist project",
			toolDetector: toolDetector{
				toolName:      "Tuist",
				primaryFile:   "Project.swift",
				optionalFiles: nil,
			},
			rootPath: createTuistFiles(t),
			want:     DetectionResult{Detected: true, ProjectTree: "Project.swift\n"},
		},
		{
			name: "Optional files in Bazel project",
			toolDetector: toolDetector{
				toolName:      "Bazel",
				primaryFile:   "WORKSPACE",
				optionalFiles: []string{"WORKSPACE.bazel"},
			},
			rootPath: createBazelFiles(t),
			want:     DetectionResult{Detected: true, ProjectTree: "WORKSPACE.bazel\n"},
		},
		{
			name: "Nested project files",
			toolDetector: toolDetector{
				toolName:      "Kotlin Gradle script",
				primaryFile:   "build.gradle.kts",
				optionalFiles: nil,
			},
			rootPath: createNestedKotlinMultiplatformFiles(t),
			want: DetectionResult{Detected: true, ProjectTree: `my-project/
· mobile/
· · android/
· · · app/
· · · · build.gradle
· · · build.gradle
· · · settings.gradle
· · · shared/
· · · · build.gradle.kts
`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.toolDetector.DetectToolIn(tt.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectToolIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectToolIn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createAndroidProjectFiles(t *testing.T, rootPath string) string {
	projectPath := path.Join(rootPath, "android")
	err := os.Mkdir(projectPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	rootBuildGradle := `
// Top-level build file where you can add configuration options common to all sub-projects/modules.
buildscript {
    ext.kotlin_version = "1.5.0"
    repositories {
        google()
        jcenter()
    }
    dependencies {
        classpath 'com.android.tools.build:gradle:4.2.1'
        classpath "org.jetbrains.kotlin:kotlin-gradle-plugin:$kotlin_version"

        // NOTE: Do not place your application dependencies here; they belong
        // in the individual module build.gradle files
    }
}

allprojects {
    repositories {
        google()
        jcenter()
    }
}

task clean(type: Delete) {
    delete rootProject.buildDir
}`
	err = os.WriteFile(path.Join(projectPath, "build.gradle"), []byte(rootBuildGradle), 0777)
	if err != nil {
		t.Fatal(err)
	}

	settingsGradle := `
include ':app'
rootProject.name = "Bitrise Sample"`
	err = os.WriteFile(path.Join(projectPath, "settings.gradle"), []byte(settingsGradle), 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(path.Join(projectPath, "app"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	appGradle := `
plugins {
    id 'com.android.application'
    id 'kotlin-android'
}

android {
    compileSdkVersion 30

    defaultConfig {
        applicationId "io.bitrise.sample.android"
        minSdkVersion 21
        targetSdkVersion 30
        versionCode 1
        versionName "1.0"

        testInstrumentationRunner "androidx.test.runner.AndroidJUnitRunner"
    }

    buildTypes {
        release {
            minifyEnabled true
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
    compileOptions {
        sourceCompatibility JavaVersion.VERSION_1_8
        targetCompatibility JavaVersion.VERSION_1_8
    }
    kotlinOptions {
        jvmTarget = '1.8'
    }
}
dependencies {
    implementation "org.jetbrains.kotlin:kotlin-stdlib:$kotlin_version"
    implementation 'androidx.core:core-ktx:1.3.2'
    implementation 'androidx.appcompat:appcompat:1.2.0'
}`
	err = os.WriteFile(path.Join(projectPath, "app/build.gradle"), []byte(appGradle), 0777)
	if err != nil {
		t.Fatal(err)
	}

	return projectPath
}

func createKotlinMultiplatformFiles(t *testing.T, rootPath string) string {
	projectPath := createAndroidProjectFiles(t, rootPath)

	err := os.Mkdir(path.Join(projectPath, "shared"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	sharedModuleGradle := `
plugins {
    kotlin("multiplatform")
    id("com.android.library")
    kotlin("plugin.serialization")
}

kotlin {
}`
	err = os.WriteFile(path.Join(projectPath, "shared/build.gradle.kts"), []byte(sharedModuleGradle), 0777)
	if err != nil {
		t.Fatal(err)
	}

	return projectPath
}

func createNestedKotlinMultiplatformFiles(t *testing.T) string {
	rootPath := path.Join(t.TempDir(), "nested")
	projectPath := path.Join(rootPath, "my-project/mobile")
	err := os.MkdirAll(projectPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	createKotlinMultiplatformFiles(t, projectPath)

	return rootPath
}

func createTuistFiles(t *testing.T) string {
	projectPath := path.Join(t.TempDir(), "tuist")
	err := os.Mkdir(projectPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(path.Join(projectPath, "node_modules"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(path.Join(projectPath, "Pods"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	projectSwift := `
let Project = Project(name: "MyProject")
`
	assert.NoError(t, os.WriteFile(path.Join(projectPath, "Project.swift"), []byte(projectSwift), 0777))

	return projectPath
}

func createBazelFiles(t *testing.T) string {
	projectPath := path.Join(t.TempDir(), "bazel")
	err := os.Mkdir(projectPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(path.Join(projectPath, "node_modules"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(path.Join(projectPath, "Pods"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, os.WriteFile(path.Join(projectPath, "WORKSPACE.bazel"), []byte(""), 0777))

	return projectPath
}
