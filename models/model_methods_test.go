package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOption(t *testing.T) {
	actual := NewOption("Project (or Workspace) path", "BITRISE_PROJECT_PATH")
	expected := &OptionModel{
		Title:          "Project (or Workspace) path",
		EnvKey:         "BITRISE_PROJECT_PATH",
		ChildOptionMap: map[string]*OptionModel{},
	}

	require.Equal(t, expected, actual)
}

func TestGetValues(t *testing.T) {
	option := OptionModel{
		ChildOptionMap: map[string]*OptionModel{},
	}
	option.ChildOptionMap["assembleAndroidTest"] = &OptionModel{}
	option.ChildOptionMap["assembleDebug"] = &OptionModel{}
	option.ChildOptionMap["assembleRelease"] = &OptionModel{}

	values := option.GetValues()

	expectedMap := map[string]bool{
		"assembleAndroidTest": false,
		"assembleDebug":       false,
		"assembleRelease":     false,
	}

	for _, value := range values {
		delete(expectedMap, value)
	}

	require.Equal(t, 0, len(expectedMap))
}

func TestLastOptions(t *testing.T) {
	// 1. level
	opt0 := NewOption("OPT0", "OPT0'_KEY")

	// 2. level
	opt01 := NewOption("OPT01", "OPT01'_KEY") // has no child
	opt0.AddOption("value1", opt01)

	opt02 := NewOption("OPT02", "OPT02'_KEY")
	opt0.AddOption("value2", opt02)

	// 3. level
	opt021 := NewOption("OPT021", "OPT021'_KEY")
	opt02.AddOption("value1", opt021)

	// 4. level
	opt0211 := NewOption("OPT0211", "OPT0211'_KEY") // has no child
	opt021.AddOption("value1", opt0211)

	opt0212 := NewOption("OPT0212", "OPT0212'_KEY")
	opt021.AddOption("value2", opt0212)

	// 5. level
	opt02121 := NewOption("OPT02121", "OPT02121'_KEY") // has no child
	opt0212.AddOption("value1", opt02121)

	lastOptions := opt0.LastOptions()
	require.Equal(t, true, len(lastOptions) == 3, fmt.Sprintf("%d", len(lastOptions)))

	optionsMap := map[string]bool{}
	for _, opt := range lastOptions {
		optionsMap[opt.Title] = true
	}

	require.Equal(t, true, optionsMap["OPT01"])
	require.Equal(t, true, optionsMap["OPT0211"])
	require.Equal(t, true, optionsMap["OPT02121"])
}
