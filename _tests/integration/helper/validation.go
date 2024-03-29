package helper

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ValidateConfigExpectation(t *testing.T, ID, expected, actual string, versions ...interface{}) {
	if !assert.ObjectsAreEqual(expected, actual) {
		s, err := replaceVersions(actual, versions...)
		require.NoError(t, err)
		fmt.Println("---------------------")
		fmt.Println("Actual config format:")
		fmt.Println("---------------------")
		fmt.Println(s)
		fmt.Println("---------------------")

		tmpDir, err := pathutil.NormalizedOSTempDirPath("__diffs__")
		require.NoError(t, err)

		expPth := filepath.Join(tmpDir, ID+"-expected.yml")
		actPth := filepath.Join(tmpDir, ID+"-actual.yml")
		require.NoError(t, fileutil.WriteStringToFile(expPth, expected))
		require.NoError(t, fileutil.WriteStringToFile(actPth, actual))
		fmt.Println("Expected: ", expPth)
		fmt.Println("Actual: ", actPth)

		_, err = exec.LookPath("opendiff")
		if err == nil {
			require.NoError(t, exec.Command("opendiff", expPth, actPth).Start())
			t.FailNow()
			return
		}
		log.Warnf("opendiff not installed, unable to open config diff")
		t.FailNow()
	}
}

func replaceVersions(str string, versions ...interface{}) (string, error) {
	for _, f := range versions {
		if format, ok := f.(string); ok {
			beforeCount := strings.Count(str, format)
			if beforeCount < 1 {
				return "", fmt.Errorf("format's original value not found, str: %s versions: %+v", str, versions)
			}
			str = strings.Replace(str, format, "%s", 1)

			afterCount := strings.Count(str, format)
			if beforeCount-1 != afterCount {
				return "", fmt.Errorf("failed to extract all versions")
			}
		}
	}
	return str, nil
}
