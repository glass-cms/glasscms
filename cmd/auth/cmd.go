package auth

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type Command struct {
	Command *cobra.Command
}

func NewAuthCommand() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use: "auth",
			Long: heredoc.Doc(`
				Authentication commands for the GlassCMS API.

				This command provides subcommands for managing authentication tokens for the GlassCMS API.
			`),
		},
	}

	cmd.Command.AddCommand(NewInitCommand().Command)
	return cmd
}
