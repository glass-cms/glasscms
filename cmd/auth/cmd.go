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
			Use: "auth",
		},
	}

	cmd.Command.AddCommand(NewInitCommand().Command)
	return cmd
}
