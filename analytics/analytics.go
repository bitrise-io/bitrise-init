package analytics

import (
	"os"

	"github.com/bitrise-io/go-utils/log"
)

// Log sends analytics log using log.RInfof by setting the stepID and data/build_slug.
func Log(tag string, format string, v ...interface{}) {
	data := map[string]interface{}{
		"build_slug": os.Getenv("BITRISE_BUILD_SLUG"),
	}
	log.RInfof("bitrise-init", tag, data, format, v)
}
