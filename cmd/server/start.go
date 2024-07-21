package server

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"

	"github.com/glass-cms/glasscms/ctx"
	"github.com/glass-cms/glasscms/database"
	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/server"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartCommand struct {
	Command *cobra.Command
	logger  *slog.Logger

	databaseConfig *database.Config
}

func NewStartCommand() *StartCommand {
	sc := &StartCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level: slog.LevelDebug,
			}),
		),
		databaseConfig: &database.Config{},
	}

	sc.Command = &cobra.Command{
		Use:   "start",
		Short: "Start the CMS server",
		RunE:  sc.Execute,
	}

	flagset := sc.Command.Flags()

	flagset.StringVar(
		&sc.databaseConfig.DSN,
		database.ArgDSN,
		"",
		"The data source name (DSN) for the database",
	)
	_ = viper.BindPFlag(database.ArgDSN, flagset.Lookup(database.ArgDSN))

	flagset.StringVar(
		&sc.databaseConfig.Driver,
		database.ArgDriver,
		"",
		"The name of the database driver",
	)
	_ = viper.BindPFlag(database.ArgDriver, flagset.Lookup(database.ArgDriver))

	flagset.IntVar(
		&sc.databaseConfig.MaxConnections,
		database.ArgMaxConnections,
		database.MaxConnectionsDefault,
		"The maximum number of connections that can be opened to the database",
	)
	_ = viper.BindPFlag(database.ArgMaxConnections, flagset.Lookup(database.ArgMaxConnections))

	flagset.IntVar(
		&sc.databaseConfig.MaxIdleConnections,
		database.ArgMaxIdleConnections,
		database.MaxIdleConnectionsDefault,
		"The maximum number of idle connections that can be maintained",
	)
	_ = viper.BindPFlag(database.ArgMaxIdleConnections, flagset.Lookup(database.ArgMaxIdleConnections))

	return sc
}

func (c *StartCommand) Execute(cmd *cobra.Command, _ []string) error {
	db, err := database.NewConnection(*c.databaseConfig)
	if err != nil {
		return err
	}

	server, err := server.New(c.logger, item.NewRepository(db))
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
