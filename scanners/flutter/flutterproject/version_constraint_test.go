package flutterproject

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVersionConstraint(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		source         VersionConstraintSource
		wantVersion    string
		wantConstraint string
		wantSource     string
		wantErr        string
	}{
		{
			name:        "Exact semver version",
			version:     "1.2.3",
			source:      "test_source",
			wantVersion: "1.2.3",
			wantSource:  "test_source",
		},
		{
			name:           "Version constraint - Caret syntax",
			version:        "^1.2.3",
			source:         "test_source",
			wantConstraint: "^1.2.3",
			wantSource:     "test_source",
		},
		{
			name:           "Version constraint - Traditional syntax",
			version:        ">=1.2.3 <2.0.0",
			source:         "test_source",
			wantConstraint: ">=1.2.3 <2.0.0",
			wantSource:     "test_source",
		},
		{
			name:    "Empty version",
			version: "",
			wantErr: "invalid version (): not a semantic version (Invalid Semantic Version) nor a version constraint (improper constraint: )",
		},
		{
			name:    "Invalid version",
			version: "asdf",
			wantErr: "invalid version (asdf): not a semantic version (Invalid Semantic Version) nor a version constraint (improper constraint: asdf)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersionConstraint(tt.version, tt.source)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)

				if tt.wantVersion != "" {
					require.NotNil(t, got.Version)
					require.Equal(t, tt.wantVersion, got.Version.String())
				} else {
					require.Nil(t, got.Version)
				}

				if tt.wantConstraint != "" {
					require.NotNil(t, got.Constraint)
					require.Equal(t, tt.wantConstraint, got.Constraint.String())
				} else {
					require.Nil(t, got.Constraint)
				}
				require.Equal(t, tt.wantSource, string(got.Source))
			}
		})
	}
}
