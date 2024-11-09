package cmd_test

import (
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/stretchr/testify/assert"
)

func TestSyncCommand(t *testing.T) {
	t.Parallel()

	t.Run("should return an error when the source type is not provided", func(t *testing.T) {
		t.Parallel()

		command := cmd.NewSyncCommand()
		command.SetArgs([]string{"test", "--live"})

		err := command.RunE(command.Command, []string{"--live"})
		assert.EqualError(t, err, "source type is required")
	})
}
