package server

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"

	"github.com/glass-cms/glasscms/ctx"
	"github.com/glass-cms/glasscms/server"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type StartCommand struct {
	Command *cobra.Command
	logger  *slog.Logger
}

func NewStartCommand() *StartCommand {
	sc := &StartCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level: slog.LevelDebug,
			}),
		),
	}

	sc.Command = &cobra.Command{
		Use:   "start",
		Short: "Start the CMS server",
		RunE:  sc.Execute,
	}

	return sc
}

func (c *StartCommand) Execute(cmd *cobra.Command, _ []string) error {
	server, err := server.New(c.logger)
	if err != nil {
		return err
	}

	if err = createServerRootFolder(); err != nil {
		return err
	}

	_ = ctx.SigtermCacellationContext(cmd.Context(), func() {
		c.logger.Info("shutting down server")
		server.Shutdown()
	})

	c.logger.Info("starting server")
	return server.ListenAndServer()
}

func createServerRootFolder() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/.glasscms", usr.HomeDir)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}

	return nil
}
