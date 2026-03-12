package ruby

import (
	"path/filepath"
	"strings"

	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
)

type testFramework struct {
	name           string
	detectionFiles []string
}

var testFrameworks = []testFramework{
	{"rspec", []string{"spec/spec_helper.rb", ".rspec"}},
	{"minitest", []string{"test/test_helper.rb"}},
}

func checkBundler(searchDir string) bool {
	log.TPrintf("Checking for Bundler")
	hasGemfileLock := utility.FileExists(filepath.Join(searchDir, "Gemfile.lock"))

	if !hasGemfileLock {
		log.TPrintf("- Gemfile.lock - not found")
		return false
	}

	log.TPrintf("- Gemfile.lock - found")
	log.TPrintf("Bundler: detected")
	return true
}

func checkRakefile(searchDir string) bool {
	log.TPrintf("Checking for Rakefile")
	hasRakefile := utility.FileExists(filepath.Join(searchDir, "Rakefile"))

	if !hasRakefile {
		log.TPrintf("- Rakefile - not found")
		return false
	}

	log.TPrintf("- Rakefile - found")
	return true
}

func checkRubyVersion(searchDir string) bool {
	log.TPrintf("Checking for Ruby version file")

	versionFiles := []string{".ruby-version", ".tool-versions"}
	for _, versionFile := range versionFiles {
		if utility.FileExists(filepath.Join(searchDir, versionFile)) {
			log.TPrintf("- %s - found", versionFile)
			return true
		}
	}

	log.TPrintf("- Ruby version file - not found")
	return false
}

func detectTestFramework(searchDir string) string {
	log.TPrintf("Checking test framework")

	for _, fw := range testFrameworks {
		for _, detectionFile := range fw.detectionFiles {
			if utility.FileExists(filepath.Join(searchDir, detectionFile)) {
				log.TPrintf("- %s - found (%s)", fw.name, detectionFile)
				return fw.name
			}
		}
	}

	log.TPrintf("- test framework - not detected")
	return ""
}

func detectRails(searchDir string) bool {
	gemfilePath := filepath.Join(searchDir, "Gemfile")
	content, err := fileutil.ReadStringFromFile(gemfilePath)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(content, "\n") {
		match := gemDeclPattern.FindStringSubmatch(line)
		if len(match) >= 2 && match[1] == "rails" {
			return true
		}
	}
	return false
}
