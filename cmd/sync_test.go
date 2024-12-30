package cmd_test

import (
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/stretchr/testify/require"
)

func Test_SyncCommandURLValidation(t *testing.T) {
	command := cmd.NewSyncCommand()
	command.SetArgs([]string{
		"filesystem",
		"../docs/commands",
		"--server", "localhost:8080",
	})

	err := command.Command.Execute()
	require.Error(t, err)
}
