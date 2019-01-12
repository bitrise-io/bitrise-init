package utility

// Delimiters used for go template execution in bitrise.yml option substitutions
const (
	TemplateDelimiterLeft  = "[["
	TemplateDelimiterRight = "]]"
)

// TemplateWithKey returns a go template string to get a key's value from a map
func TemplateWithKey(key string) string {
	return TemplateDelimiterLeft + "." + key + TemplateDelimiterRight
}
