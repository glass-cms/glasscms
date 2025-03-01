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

// TODO: Find a way to omit this command from the help command and docs, whilst still allowing its children to be shown.

func NewCommand() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use: "server",
		},
	}

	cmd.Command.AddCommand(NewStartCommand().Command)
	return cmd
}
