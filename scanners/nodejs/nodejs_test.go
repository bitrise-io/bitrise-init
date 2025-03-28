package nodejs

import (
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/stretchr/testify/require"
)

func TestDetectPlatform(t *testing.T) {
	scanner := NewScanner()
	searchDir := "/Users/zoltan.szabo/Repos/demo-projects/nestjs-01-cats-app"

	detected, err := scanner.DetectPlatform(searchDir)

	opts, _, _, optsErr := scanner.Options()
	confs, confsErr := scanner.Configs(models.SSHKeyActivationNone)

	defOpts := scanner.DefaultOptions()
	defConfs, defConfsErr := scanner.DefaultConfigs()

	require.NoError(t, err)
	require.Equal(t, true, detected)

	require.NoError(t, optsErr)
	require.NotEmpty(t, opts)
	require.NotEmpty(t, confs)
	require.NoError(t, confsErr)

	require.NotEmpty(t, defOpts)
	require.NotEmpty(t, defConfs)
	require.NoError(t, defConfsErr)

}
