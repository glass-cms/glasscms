package server

import (
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/spf13/cobra"
)

type Config struct {
	Database *database.Config `mapstructure:"database"`
}

type Command struct {
	Command *cobra.Command
}

func NewCommand() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:    "server",
			Hidden: true, // Server command is a group command not meant to be run directly.
		},
	}

	cmd.Command.AddCommand(NewStartCommand().Command)
	return cmd
}
