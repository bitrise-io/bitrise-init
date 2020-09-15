package analytics

import (
	"github.com/bitrise-io/go-utils/log"
)

// LogError sends analytics log using log.RErrorf by setting the stepID and data/build_slug.
func LogError(tag string, data map[string]interface{}, format string, v ...interface{}) {
	log.RErrorf("bitrise-init", tag, data, format, v...)
}
