package expo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func configName(hasAndroidProject, hasIosProject, hasNPMTest bool) string {
	name := "react-native-expo"
	if hasAndroidProject {
		name += "-android"
	}
	if hasIosProject {
		name += "-ios"
	}
	if hasNPMTest {
		name += "-test"
	}
	name += "-config"
	return name
}

func defaultConfigName() string {
	return "default-react-native-expo-config"
}

// FindDependency ...
func FindDependency(filePath string, dependency string) (bool, error) {
	rawJSON, err := readFile(filePath)
	if err != nil {
		return false, err
	}

	value, err := getValueByUnmarshalToInterface(rawJSON["dependencies"], dependency)
	return len(value) > 0, err
}

func readFile(pth string) (map[string]*json.RawMessage, error) {
	fmt.Printf("Reading file - %s", pth)
	raw, err := ioutil.ReadFile(pth)
	if err != nil {
		return nil, err
	}

	fmt.Printf("File content: %s", string(raw))

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(raw, &objmap)

	return objmap, err
}

func getValueByUnmarshalToInterface(foo *json.RawMessage, key string) (string, error) {
	var tmp map[string]interface{}
	if err := json.Unmarshal(*foo, &tmp); err != nil {
		return "", err
	}

	value, ok := tmp[key].(string)
	if !ok {
		return "", nil
	}
	return value, nil
}
