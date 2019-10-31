package maintenance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-add-new-project/config"
	"github.com/bitrise-io/go-utils/sliceutil"
)

type ResponseBody []DirectoryEntry
type DirectoryEntry struct {
	Name string `json:"name"`
}

var stacks = []string{
	"linux-docker-android-lts",
	"linux-docker-android",
	"osx-vs4mac-beta",
	"osx-vs4mac-previous-stable",
	"osx-vs4mac-stable",
	"osx-xcode-10.0.x",
	"osx-xcode-10.1.x",
	"osx-xcode-10.2.x",
	"osx-xcode-10.3.x",
	"osx-xcode-11.0.x",
	"osx-xcode-11.1.x",
	"osx-xcode-11.2.x",
	"osx-xcode-8.3.x",
	"osx-xcode-9.4.x",
	"osx-xcode-edge",
}

func TestStackChange(t *testing.T) {
	resp, err := http.Get("https://api.github.com/repos/bitrise-io/bitrise.io/contents/system_reports")
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

	var rb ResponseBody
	if err := json.Unmarshal(bytes, &rb); err != nil {
		t.Fatalf("Error unmarshalling stack data from string (%s): %s", bytes, err)
	}

	if len(config.Stacks) != len(rb) {
		t.Fatalf("Stack list changed")
	}

	for _, de := range rb {
		trimmed := strings.TrimSuffix(de.Name, ".log")
		if !sliceutil.IsStringInSlice(trimmed, stacks) {
			t.Fatalf("Stack list changed")
		}
	}

}
