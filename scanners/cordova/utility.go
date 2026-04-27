package cordova

import (
	"encoding/xml"

	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

const configXMLBasePath = "config.xml"

// WidgetModel ...
type WidgetModel struct {
	XMLNSCDV string `xml:"xmlns cdv,attr"`
}

func parseConfigXMLContent(content string) (WidgetModel, error) {
	widget := WidgetModel{}
	if err := xml.Unmarshal([]byte(content), &widget); err != nil {
		return WidgetModel{}, err
	}
	return widget, nil
}

// ParseConfigXML ...
func ParseConfigXML(pth string) (WidgetModel, error) {
	content, err := fileutil.ReadStringFromFile(pth)
	if err != nil {
		return WidgetModel{}, err
	}
	return parseConfigXMLContent(content)
}

// FilterRootConfigXMLFile ...
func FilterRootConfigXMLFile(fileList []string) (string, error) {
	filters := append(utility.CommonExcludeFilters(), pathutil.BaseFilter(configXMLBasePath, true))
	configXMLs, err := pathutil.FilterPaths(fileList, filters...)
	if err != nil {
		return "", err
	}

	if len(configXMLs) == 0 {
		return "", nil
	}

	return configXMLs[0], nil
}
