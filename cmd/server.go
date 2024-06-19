package cmd

import (
	"log/slog"
	"os"

	"github.com/glass-cms/glasscms/ctx"
	"github.com/glass-cms/glasscms/server"
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

func (c *ServerCommand) Execute(cmd *cobra.Command, _ []string) error {
	server, err := server.New(c.logger)
	if err != nil {
		return err
	}

	_ = ctx.SigtermCacellationContext(cmd.Context(), func() {
		c.logger.Info("shutting down server")
		server.Shutdown()
	})

	c.logger.Info("starting server")
	return server.ListenAndServer()
}
