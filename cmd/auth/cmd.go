package auth

import (
	"github.com/spf13/cobra"
)

type Command struct {
	Command *cobra.Command
}

func NewAuthCommand() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:    "auth",
			Hidden: true, // Auth command is a group command not meant to be run directly.
		},
	}

	cmd.Command.AddCommand(NewInitCommand().Command)
	return cmd
}
