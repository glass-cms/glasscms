package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/user"

	"github.com/MakeNowJust/heredoc"
	"github.com/glass-cms/glasscms/internal/auth"
	authRepository "github.com/glass-cms/glasscms/internal/auth/repository"
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/internal/item"
	itemRepository "github.com/glass-cms/glasscms/internal/item/repository"
	"github.com/glass-cms/glasscms/internal/server"
	internalMiddleware "github.com/glass-cms/glasscms/internal/server/middleware"
	ctx "github.com/glass-cms/glasscms/pkg/context"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/glass-cms/glasscms/pkg/mediatype"
	"github.com/glass-cms/glasscms/pkg/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartCommand struct {
	Command *cobra.Command

	databaseConfig database.Config
}

func NewStartCommand() *StartCommand {
	sc := &StartCommand{
		databaseConfig: database.Config{},
	}

	sc.Command = &cobra.Command{
		Use:   "start",
		Short: "Start the GlassCMS API server",
		Long: heredoc.Doc(`
			Start the GlassCMS API server with the specified configuration.

			This command initializes and starts the CMS server with database connectivity
			and all required services. It sets up the HTTP server with appropriate middleware
			for authentication, content negotiation, and request tracking.

			The server will continue running until it receives a termination signal.
		`),
		RunE: sc.Execute,
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
	logger, err := log.NewLogger()
	if err != nil {
		return err
	}

	logger.Debug("connecting to database",
		slog.String("driver", c.databaseConfig.Driver),
		slog.String("dsn", c.databaseConfig.DSN),
	)

	db, err := database.NewConnection(c.databaseConfig)
	if err != nil {
		return err
	}

	errHandler, err := database.NewErrorHandler(c.databaseConfig)
	if err != nil {
		return err
	}

	itemRepo := itemRepository.NewRepository(db, errHandler)
	itemService := item.NewService(db, itemRepo)

	authRepo := authRepository.NewRepository(db, errHandler)
	authService := auth.NewAuth(db, authRepo, logger)

	server, err := server.New(logger, itemService, []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.ContentType(mediatype.ApplicationJSON),
		middleware.Accept(mediatype.ApplicationJSON),
		internalMiddleware.AuthMiddleware(authService),
	})
	if err != nil {
		return err
	}

	if err = createServerRootFolder(); err != nil {
		return err
	}

	_ = ctx.SigtermCacellationContext(cmd.Context(), func() {
		slog.Info("shutting down server")
		server.Shutdown()
	})

	logger.Info("starting server")
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
