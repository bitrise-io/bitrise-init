package flutterproject

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const asdfConfigRelPath = ".tool-versions"

type ASDFVersionReader struct{}

func (r ASDFVersionReader) ReadSDKVersions(projectRootDir string) (*VersionConstraint, *VersionConstraint, error) {
	asdfConfigPth := filepath.Join(projectRootDir, asdfConfigRelPath)
	f, err := os.Open(asdfConfigPth)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, err
	}

	if f == nil {
		return nil, nil, nil
	}

	versionStr, err := parseASDFFlutterVersion(f)
	if err != nil {
		return nil, nil, err
	}
	if versionStr == "" {
		return nil, nil, nil
	}

	flutterSDKVersion, err := NewVersionConstraint(versionStr, ASDFConfigVersionSource)
	if err != nil {
		return nil, nil, err
	}

	return flutterSDKVersion, nil, nil
}

// flutter 3.7.12
func parseASDFFlutterVersion(asdfConfigReader io.Reader) (string, error) {
	scanner := bufio.NewScanner(asdfConfigReader)
	scanner.Split(bufio.ScanLines)
	versionStr := ""
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "flutter ") {
			versionStr = strings.TrimPrefix(line, "flutter ")
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return versionStr, nil
}
