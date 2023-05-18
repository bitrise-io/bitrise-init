package flutterproject

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const fvmConfigRelPath = ".fvm/fvm_config.json"

type FVMVersionReader struct{}

func (r FVMVersionReader) ReadSDKVersions(projectRootDir string) (*VersionConstraint, *VersionConstraint, error) {
	fvmConfigPth := filepath.Join(projectRootDir, fvmConfigRelPath)
	f, err := os.Open(fvmConfigPth)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, err
	}

	if f == nil {
		return nil, nil, nil
	}

	versionStr, err := parseFVMFlutterVersion(f)
	if err != nil {
		return nil, nil, err
	}

	flutterSDKVersion, err := NewVersionConstraint(versionStr, FVMConfigVersionSource)
	if err != nil {
		return nil, nil, err
	}

	return flutterSDKVersion, nil, nil
}

func parseFVMFlutterVersion(fvmConfigReader io.Reader) (string, error) {
	type fvmConfig struct {
		FlutterSdkVersion string `json:"flutterSdkVersion"`
	}

	var config fvmConfig
	d := json.NewDecoder(fvmConfigReader)
	if err := d.Decode(&config); err != nil {
		return "", err
	}

	return config.FlutterSdkVersion, nil
}
