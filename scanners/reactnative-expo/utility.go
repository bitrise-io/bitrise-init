package expo

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
