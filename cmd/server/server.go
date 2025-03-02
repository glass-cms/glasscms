package server

import (
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/lithammer/dedent"
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
			Use:   "server",
			Short: "Server management commands",
			Long: dedent.Dedent(`
				Server management commands for the GlassCMS API.
				This command provides subcommands for managing the GlassCMS API server,
				including starting the server and other server-related operations.
			`),
		},
	}
	cmd.Command.AddCommand(NewStartCommand().Command)
	return cmd
}
