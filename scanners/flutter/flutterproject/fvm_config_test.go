package flutterproject

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseFVMFlutterVersion(t *testing.T) {
	tests := []struct {
		name            string
		fvmConfigReader io.Reader
		wantFlutterSDK  string
		wantErr         string
	}{
		{
			name:            "Real fvm_config.json",
			fvmConfigReader: strings.NewReader(realFVMConfigJSON),
			wantFlutterSDK:  "3.7.12",
		},
		{
			name:            "Empty fvm_config.json",
			fvmConfigReader: strings.NewReader(""),
			wantErr:         "EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlutterSDK, err := parseFVMFlutterVersion(tt.fvmConfigReader)
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

const realFVMConfigJSON = `{
  "flutterSdkVersion": "3.7.12",
  "flavors": {}
}`
