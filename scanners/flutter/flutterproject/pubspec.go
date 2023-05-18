package flutterproject

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const pubspecRelPath = "pubspec.yaml"

type PubspecVersionReader struct{}

func (r PubspecVersionReader) ReadSDKVersions(projectRootDir string) (*VersionConstraint, *VersionConstraint, error) {
	pubspecPth := filepath.Join(projectRootDir, pubspecRelPath)
	f, err := os.Open(pubspecPth)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, err
	}

	if f == nil {
		return nil, nil, nil
	}

	flutterVersionStr, dartVersionStr, err := parsePubspecSDKVersions(f)
	if err != nil {
		return nil, nil, err
	}

	var flutterVersion *VersionConstraint
	if flutterVersionStr != "" {
		flutterVersion, err = NewVersionConstraint(flutterVersionStr, PubspecVersionSource)
		if err != nil {
			return nil, nil, err
		}
	}

	var dartVersion *VersionConstraint
	if dartVersionStr != "" {
		dartVersion, err = NewVersionConstraint(dartVersionStr, PubspecVersionSource)
		if err != nil {
			return nil, nil, err
		}
	}

	return flutterVersion, dartVersion, nil
}

func parsePubspecSDKVersions(pubspecReader io.Reader) (string, string, error) {
	type pubspec struct {
		Environment struct {
			Dart    string `yaml:"sdk"`
			Flutter string `yaml:"flutter"`
		} `yaml:"environment"`
	}

	var config pubspec
	d := yaml.NewDecoder(pubspecReader)
	if err := d.Decode(&config); err != nil {
		return "", "", err
	}

	return config.Environment.Flutter, config.Environment.Dart, nil
}
