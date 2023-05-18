package flutterproject

import (
	"encoding/json"
	"io"
	"path/filepath"
)

const fvmConfigRelPath = ".fvm/fvm_config.json"

type FVMVersionReader struct {
	fileOpener FileOpener
}

func NewFVMVersionReader(fileOpener FileOpener) FVMVersionReader {
	return FVMVersionReader{
		fileOpener: fileOpener,
	}
}

func (r FVMVersionReader) ReadSDKVersions(projectRootDir string) (*VersionConstraint, *VersionConstraint, error) {
	fvmConfigPth := filepath.Join(projectRootDir, fvmConfigRelPath)
	f, err := r.fileOpener.OpenFile(fvmConfigPth)
	if err != nil {
		return nil, nil, err
	}

	if f == nil {
		return nil, nil, nil
	}

	versionStr, err := parseFVMFlutterVersion(f)
	if err != nil {
		return nil, nil, err
	}
	if versionStr == "" {
		return nil, nil, nil
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
