package kmp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_detectIncludedProjects(t *testing.T) {
	tests := []struct {
		name                     string
		settingGradleFileContent string
		want                     []string
		wantErr                  bool
	}{
		{
			name: "home test",
			settingGradleFileContent: `
include ':app1', ':app2'
include ':app3'
include(':app4', ':app5')
include(':app6')

include ":app7", ":app8"
include ":app9"
include(":app10", ":app11")
include(":app12")

include "app13"
include 'app14'

include(":backend:app15")
include(":backend:api:app16")
`,
			want: []string{":app1", ":app2", ":app3", ":app4", ":app5", ":app6", ":app7",
				":app8", ":app9", ":app10", ":app11", ":app12", ":app13", ":app14",
				":backend:app15", ":backend:api:app16"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectProjectIncludesInContent(tt.settingGradleFileContent)
			require.Equal(t, tt.want, got)
		})
	}
}
