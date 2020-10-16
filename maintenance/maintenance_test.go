package maintenance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

func stacks() []string {
	return []string{
		"linux-docker-android-lts",
		"linux-docker-android",
		"osx-vs4mac-beta",
		"osx-vs4mac-previous-stable",
		"osx-vs4mac-stable",
		"osx-xcode-10.1.x",
		"osx-xcode-10.2.x",
		"osx-xcode-10.3.x",
		"osx-xcode-11.0.x",
		"osx-xcode-11.1.x",
		"osx-xcode-11.2.x",
		"osx-xcode-11.3.x",
		"osx-xcode-11.4.x",
		"osx-xcode-11.5.x",
		"osx-xcode-11.6.x",
		"osx-xcode-11.7.x",
		"osx-xcode-12.0.x",
		"osx-xcode-12.1.x",
		"osx-xcode-12.2.x",
		"osx-xcode-9.4.x",
		"osx-xcode-edge",
	}
}

type report struct {
	Name string `json:"name"`
}

type systemReports []report

func (reports systemReports) Stacks() (s []string) {
	for _, report := range reports {
		s = append(s, strings.TrimSuffix(report.Name, ".log"))
	}
	return
}

func TestStackChange(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/bitrise-io/bitrise.io/contents/system_reports", nil)
	if err != nil {
		t.Fatalf("setup: faile to create request: %s", err)
	}

	token := os.Getenv("GIT_BOT_USER_ACCESS_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GIT_BOT_USER_ACCESS_TOKEN")))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error getting current stack list from GitHub: %s", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("Error closing response body")
		}
	}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading stack info from GitHub response: %s", err)
	}

	var reports systemReports
	if err := json.Unmarshal(bytes, &reports); err != nil {
		t.Fatalf("Error unmarshalling stack data from string (%s): %s", bytes, err)
	}

	if expected := reports.Stacks(); !reflect.DeepEqual(expected, stacks()) {
		t.Fatalf("Stack list changed, current: %v, expecting: %v", stacks(), expected)
	}
}
