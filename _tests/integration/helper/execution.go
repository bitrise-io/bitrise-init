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
	Name              string
	RepoURL           string
	Branch            string
	RelativeSearchDir string
	ExpectedResult    string
	ExpectedVersions  []interface{}
}

var globalHelper *testHelper

type testHelper struct {
	repoCache map[string]string
}

func newTestHelper() *testHelper {
	return &testHelper{
		repoCache: make(map[string]string),
	}
}

func Execute(t *testing.T, testCases []TestCase) {
	if globalHelper == nil {
		globalHelper = newTestHelper()
	}

	cloneDir := t.TempDir()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log("Executing :", testCase.Name)

			cacheKey := testCase.RepoURL
			if testCase.Branch != "" {
				cacheKey = cacheKey + "@" + testCase.Branch
			}

			var sampleAppDir string
			if _, ok := globalHelper.repoCache[cacheKey]; !ok {
				sampleAppDir = filepath.Join(cloneDir, testCase.Name)

				if testCase.Branch != "" {
					GitCloneBranch(t, sampleAppDir, testCase.RepoURL, testCase.Branch)
				} else {
					GitClone(t, sampleAppDir, testCase.RepoURL)
				}

				globalHelper.repoCache[cacheKey] = sampleAppDir
			} else {
				sampleAppDir = globalHelper.repoCache[cacheKey]
			}

			resultDir := t.TempDir()
			searchDir := filepath.Join(sampleAppDir, testCase.RelativeSearchDir)

			_, err := scanner.GenerateAndWriteResults(searchDir, resultDir, output.YAMLFormat)
			require.NoError(t, err)

			scanResultPth := filepath.Join(resultDir, "result.yml")

			result, err := fileutil.ReadStringFromFile(scanResultPth)
			require.NoError(t, err)

			ValidateConfigExpectation(t, testCase.Name, strings.TrimSpace(testCase.ExpectedResult), strings.TrimSpace(result), testCase.ExpectedVersions)
		})

	}
}
