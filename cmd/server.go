package cmd

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type ServerCommand struct {
	Command *cobra.Command

	logger *slog.Logger
}

// NewServerCommand creates a new cobra.Command for
// starting the CMS server.
func NewServerCommand() *ServerCommand {
	sc := &ServerCommand{
		logger: slog.New(
			// TODO: Make handler type configurable.
			tint.NewHandler(os.Stdout, &tint.Options{
				// TODO: Make configurable.
				Level: slog.LevelDebug,
			}),
		),
	}

	sc.Command = &cobra.Command{
		Use:   "server",
		Short: "Start the cms server",
		RunE:  sc.Execute,
	}

	return sc
}

func (c *ServerCommand) Execute(_ *cobra.Command, _ []string) error {
	c.logger.Info("Starting server")
	return nil
}
