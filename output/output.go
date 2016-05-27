package output

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/bitrise-io/go-utils/fileutil"
)

// Format ...
type Format uint8

const (
	// RawFormat ...
	RawFormat Format = iota
	// JSONFormat ...
	JSONFormat
	// YAMLFormat ...
	YAMLFormat
)

// ParseFormat ...
func ParseFormat(format string) (Format, error) {
	switch strings.ToLower(format) {
	case "raw":
		return RawFormat, nil
	case "json":
		return JSONFormat, nil
	case "yaml":
		return YAMLFormat, nil
	}

	var f Format
	return f, fmt.Errorf("not a valid format: %s", format)
}

// String ...
func (format Format) String() string {
	switch format {
	case RawFormat:
		return "raw"
	case JSONFormat:
		return "json"
	case YAMLFormat:
		return "yaml"
	}

	return "unknown"
}

// Print ...
func Print(a interface{}, format Format, pth string) error {
	str := ""
	ext := ""

	switch format {
	case RawFormat:
		str = fmt.Sprint(a)
		ext = ".txt"
	case JSONFormat:
		bytes, err := json.MarshalIndent(a, "", "\t")
		if err != nil {
			return err
		}
		str = string(bytes)
		ext = ".json"
	case YAMLFormat:
		bytes, err := yaml.Marshal(a)
		if err != nil {
			return err
		}
		str = string(bytes)
		ext = ".yml"
	default:
		return fmt.Errorf("not a valid format: %s", format)
	}

	if pth != "" {
		fileExt := filepath.Ext(pth)
		if fileExt != "" {
			pth = strings.TrimSuffix(pth, fileExt)
		}
		pth = pth + ext

		return fileutil.WriteStringToFile(pth, str)
	}

	fmt.Println(str)
	return nil
}
