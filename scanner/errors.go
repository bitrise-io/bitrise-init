package scanner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bitrise-io/bitrise-init/scanners"
	"github.com/bitrise-io/bitrise-init/step"
)

// Detail ...
type Detail struct {
	Title       string
	Description string
}

func newDetailedErrorRecommendation(detail Detail) step.Recommendation {
	return step.Recommendation{
		"DetailedError": map[string]string{
			"Title":       detail.Title,
			"Description": detail.Description,
		},
	}
}

// DetailBuilder ...
type DetailBuilder = func(...string) Detail

// PatternErrorMatcher ...
type PatternErrorMatcher struct {
	defaultHandler DetailBuilder
	handlers       map[string]DetailBuilder
}

func newPatternErrorMatcher(defaultHandler DetailBuilder, handlers map[string]DetailBuilder) *PatternErrorMatcher {
	m := PatternErrorMatcher{
		handlers:       handlers,
		defaultHandler: defaultHandler,
	}

	return &m
}

// Run ...
func (m *PatternErrorMatcher) Run(msg string) step.Recommendation {
	for pattern, handler := range m.handlers {
		re := regexp.MustCompile(pattern)
		if re.MatchString(msg) {
			matches := re.FindStringSubmatch((msg))

			if len(matches) > 1 {
				matches = matches[1:]
			}

			if matches != nil {
				detail := handler(matches...)
				return newDetailedErrorRecommendation(detail)
			}
		}
	}

	detail := m.defaultHandler(msg)
	return newDetailedErrorRecommendation(detail)
}

func mapRecommendation(tag, err string) step.Recommendation {
	var matcher *PatternErrorMatcher
	switch tag {
	case noPlatformDetectedTag:
		matcher = newNoPlatformDetectedMatcher()
	case detectPlatformFailedTag:
		matcher = newDetectPlatformFailedMatcher()
	case optionsFailedTag:
		matcher = newOptionsFailedMatcher()
	}

	if matcher != nil {
		return matcher.Run(err)
	}
	return nil
}

// noPlatformDetectedTag
func newNoPlatformDetectedMatcher() *PatternErrorMatcher {
	return newPatternErrorMatcher(
		newNoPlatformDetectedGenericDetail,
		nil,
	)
}

func newNoPlatformDetectedGenericDetail(params ...string) Detail {
	return Detail{
		Title:       "We couldn’t recognize your platform.",
		Description: fmt.Sprintf("Our auto-configurator supports %s projects. If you’re adding something else, skip this step and configure your Workflow manually.", strings.Join(availableScanners(), ", ")),
	}
}

func availableScanners() (scannerNames []string) {
	for _, scanner := range scanners.ProjectScanners {
		scannerNames = append(scannerNames, scanner.Name())
	}
	for _, scanner := range scanners.AutomationToolScanners {
		scannerNames = append(scannerNames, scanner.Name())
	}
	return
}

// detectPlatformFailedTag
func newDetectPlatformFailedMatcher() *PatternErrorMatcher {
	return newPatternErrorMatcher(
		newDetectPlatformFailedGenericDetail,
		nil,
	)
}

func newDetectPlatformFailedGenericDetail(params ...string) Detail {
	err := params[0]
	return Detail{
		Title:       "We couldn’t parse your project files.",
		Description: fmt.Sprintf("Our auto-configurator returned the following error:\n%s", err),
	}
}

// optionsFailedTag
func newOptionsFailedMatcher() *PatternErrorMatcher {
	return newPatternErrorMatcher(
		newOptionsFailedGenericDetail,
		map[string]DetailBuilder{
			`No Gradle Wrapper \(gradlew\) found\.`:                                                                                 newGradlewNotFoundDetail,
			`app\.json file \((.+)\) missing or empty (.+) entry\nThe app\.json file needs to contain:`:                             newAppJSONIssueDetail,
			`app\.json file \((.+)\) missing or empty (.+) entry\nIf the project uses Expo Kit the app.json file needs to contain:`: newExpoAppJSONIssueDetail,
		},
	)
}

var newOptionsFailedGenericDetail = newDetectPlatformFailedGenericDetail

func newGradlewNotFoundDetail(params ...string) Detail {
	return Detail{
		Title:       "We couldn’t find your Gradle Wrapper. Please make sure there is a gradlew file in your project’s root directory.",
		Description: `The Gradle Wrapper ensures that the right Gradle version is installed and used for the build. You can find out more about <a target="_blank" href="https://docs.gradle.org/current/userguide/gradle_wrapper.html">the Gradle Wrapper in the Gradle docs</a>.`,
	}
}

func newAppJSONIssueDetail(params ...string) Detail {
	appJSONPath := params[0]
	entryName := params[1]
	return Detail{
		Title: fmt.Sprintf("Your app.json file (%s) doesn’t have a %s field.", appJSONPath, entryName),
		Description: `The app.json file needs to contain the following entries:
- name
- displayName`,
	}
}

func newExpoAppJSONIssueDetail(params ...string) Detail {
	appJSONPath := params[0]
	entryName := params[1]
	return Detail{
		Title: fmt.Sprintf("Your app.json file (%s) doesn’t have a %s field.", appJSONPath, entryName),
		Description: `If your project uses Expo Kit, the app.json file needs to contain the following entries:
- expo/name
- expo/ios/bundleIdentifier
- expo/android/package`,
	}
}
