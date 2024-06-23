package server

import (
	"github.com/spf13/cobra"
)

type Command struct {
	Command *cobra.Command
}

func NewCommand() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use: "server",
		},
	}

	cmd.Command.AddCommand(NewStartCommand().Command)
	return cmd
}
