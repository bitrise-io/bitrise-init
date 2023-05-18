package flutterproject

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseASDFFlutterVersion(t *testing.T) {
	tests := []struct {
		name             string
		asdfConfigReader io.Reader
		wantFlutterSDK   string
		wantErr          string
	}{
		{
			name:             "Real .tool-versions",
			asdfConfigReader: strings.NewReader(realToolVersions),
			wantFlutterSDK:   "3.7.12",
		},
		{
			name:             "Empty .tool-versions",
			asdfConfigReader: strings.NewReader(""),
			wantErr:          "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlutterSDK, err := parseASDFFlutterVersion(tt.asdfConfigReader)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Empty(t, gotFlutterSDK)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantFlutterSDK, gotFlutterSDK)
			}
		})
	}
}

const realToolVersions = `flutter 3.7.12
`
