package helper

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/stretchr/testify/require"
	"testing"
)

func GitClone(t *testing.T, dir, uri string) {
	fmt.Printf("cloning into: %s\n", dir)
	g, err := git.New(dir)
	require.NoError(t, err)

	command := g.Clone(uri)
	addShallowOption(command)

	require.NoError(t, command.Run())
}

func GitCloneBranch(t *testing.T, dir, uri, branch string) {
	fmt.Printf("cloning into: %s\n", dir)
	g, err := git.New(dir)
	require.NoError(t, err)

	command := g.CloneTagOrBranch(uri, branch)
	addShallowOption(command)

	require.NoError(t, command.Run())
}

func addShallowOption(command *command.Model) {
	args := command.GetCmd().Args
	// The command arguments will always start "git clone" as the first two words,
	// so it is safe to just use the 2nd index.
	firstPart := args[:2]
	secondPart := args[2:]

	command.GetCmd().Args = append(firstPart, append([]string{"--depth=1"}, secondPart...)...)
}
