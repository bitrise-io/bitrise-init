package utility

// TemplateWithKey returns a go template string to get a key's value from a map
func TemplateWithKey(key string) string {
	return "{{." + key + "}}"
}
