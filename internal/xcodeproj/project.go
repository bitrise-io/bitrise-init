package xcodeproj

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// ProjectModel ...
type ProjectModel struct {
	Pth           string
	Name          string
	SDKs          []string
	SharedSchemes []SchemeModel
	Targets       []TargetModel
}

// NewProject ...
func NewProject(xcodeprojPth string) (ProjectModel, error) {
	project := ProjectModel{
		Pth:  xcodeprojPth,
		Name: strings.TrimSuffix(filepath.Base(xcodeprojPth), filepath.Ext(xcodeprojPth)),
	}

	// SDK
	pbxprojPth := filepath.Join(xcodeprojPth, "project.pbxproj")

	if exist, err := pathutil.IsPathExists(pbxprojPth); err != nil {
		return ProjectModel{}, err
	} else if !exist {
		return ProjectModel{}, fmt.Errorf("Project descriptor not found at: %s", pbxprojPth)
	}

	sdks, err := getBuildConfigSDKs(pbxprojPth)
	if err != nil {
		return ProjectModel{}, err
	}

	project.SDKs = sdks

	// Shared Schemes
	schemes, err := projectSharedSchemes(xcodeprojPth)
	if err != nil {
		return ProjectModel{}, err
	}

	project.SharedSchemes = schemes

	// Targets
	targets, err := projectTargets(xcodeprojPth)
	if err != nil {
		return ProjectModel{}, err
	}

	project.Targets = targets

	return project, nil
}

func getBuildConfigSDKs(pbxprojPth string) ([]string, error) {
	content, err := fileutil.ReadStringFromFile(pbxprojPth)
	if err != nil {
		return []string{}, err
	}

	return getBuildConfigSDKsFromContent(content)
}

func getBuildConfigSDKsFromContent(pbxprojContent string) ([]string, error) {
	sdkMap := map[string]bool{}

	beginXCBuildConfigurationSection := `/* Begin XCBuildConfiguration section */`
	endXCBuildConfigurationSection := `/* End XCBuildConfiguration section */`
	isXCBuildConfigurationSection := false

	// SDKROOT = macosx;
	pattern := `SDKROOT = (?P<sdk>.*);`
	regexp := regexp.MustCompile(pattern)

	scanner := bufio.NewScanner(strings.NewReader(pbxprojContent))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == endXCBuildConfigurationSection {
			break
		}

		if strings.TrimSpace(line) == beginXCBuildConfigurationSection {
			isXCBuildConfigurationSection = true
			continue
		}

		if !isXCBuildConfigurationSection {
			continue
		}

		if match := regexp.FindStringSubmatch(line); len(match) == 2 {
			sdk := match[1]
			sdkMap[sdk] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	sdks := []string{}
	for sdk := range sdkMap {
		sdks = append(sdks, sdk)
	}

	return sdks, nil
}
