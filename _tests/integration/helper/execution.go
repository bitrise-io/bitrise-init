package helper

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-init/output"
	"github.com/bitrise-io/bitrise-init/scanner"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name             string
	RepoURL          string
	Branch           string
	ExpectedResult   string
	ExpectedVersions []interface{}
}

func Execute(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log("Executing :", testCase.Name)

			sampleAppDir := filepath.Join(t.TempDir(), testCase.Name)

			if testCase.Branch != "" {
				GitCloneBranch(t, sampleAppDir, testCase.RepoURL, testCase.Branch)
			} else {
				GitClone(t, sampleAppDir, testCase.RepoURL)
			}

			_, err := scanner.GenerateAndWriteResults(sampleAppDir, sampleAppDir, output.YAMLFormat)
			require.NoError(t, err)

			scanResultPth := filepath.Join(sampleAppDir, "result.yml")

			result, err := fileutil.ReadStringFromFile(scanResultPth)
			require.NoError(t, err)

			ValidateConfigExpectation(t, testCase.Name, strings.TrimSpace(testCase.ExpectedResult), strings.TrimSpace(result), testCase.ExpectedVersions)
		})

	}
}
