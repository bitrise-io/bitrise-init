package helper

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/stretchr/testify/require"
	"testing"
)

func GitClone(t *testing.T, dir, uri string) {
	fmt.Printf("cloning into: %s\n", dir)
	g, err := git.New(dir)
	require.NoError(t, err)
	require.NoError(t, g.Clone(uri).Run())
}

func GitCloneBranch(t *testing.T, dir, uri, branch string) {
	fmt.Printf("cloning into: %s\n", dir)
	g, err := git.New(dir)
	require.NoError(t, err)
	require.NoError(t, g.CloneTagOrBranch(uri, branch).Run())
}
