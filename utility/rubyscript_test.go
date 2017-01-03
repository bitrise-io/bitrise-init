package utility

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunRubyScriptForOutput(t *testing.T) {
	gemfileContent := `source 'https://rubygems.org'
gem 'json'
`

	rubyScriptContent := `require 'json'

hash = {
    'test_key': 'test_value'
}
puts hash.to_json
`

	expectedOut := ""
	actualOut, err := runRubyScriptForOutput(rubyScriptContent, gemfileContent, "", []string{})
	require.NoError(t, err)
	require.Equal(t, expectedOut, actualOut)
}
