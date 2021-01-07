package android

import "fmt"

type containsLocalPropertiesError struct {
	Path string
}

func (e *containsLocalPropertiesError) Error() string {
	return fmt.Sprintf("the local.properties file should NOT be checked into Version Control Systems, as it contains information specific to your local configuration, the location of the file is: %s", e.Path)
}

func newContainsLocalPropertiesError(path string) *containsLocalPropertiesError {
	return &containsLocalPropertiesError{Path: path}
}

func isContainsLocalPropertiesError(err error) bool {
	_, ok := err.(*containsLocalPropertiesError)
	return ok
}
