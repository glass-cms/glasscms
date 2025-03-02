package auth

import (
	"time"

	"github.com/glass-cms/glasscms/internal/auth"
	"github.com/glass-cms/glasscms/internal/auth/repository"
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type InitCommand struct {
	Command        *cobra.Command
	databaseConfig database.Config
}

func NewInitCommand() *InitCommand {
	cmd := &InitCommand{}

	cmd.Command = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new authentication token",
		Long: dedent.Dedent(`
			Initialize a new authentication token for API access.

			This command creates a new authentication token that can be used to authenticate 
			requests to the GlassCMS API. By default, the token is valid for 24 hours.

			The token is displayed only once upon creation and should be stored securely.
			It cannot be retrieved later, so make sure to save it in a secure location.
		`),
		Example: dedent.Dedent(`
			# Create a new token with default settings
			glasscms auth init

			# Create a new token with a specific database driver and DSN
			glasscms auth init --driver postgres --dsn "postgres://user:password@localhost:5432/glasscms"
		`),
		RunE: cmd.Execute,
	}

	flagset := cmd.Command.Flags()

	flagset.StringVar(
		&cmd.databaseConfig.DSN,
		database.ArgDSN,
		"",
		"The data source name (DSN) for the database connection",
	)
	_ = viper.BindPFlag(database.ArgDSN, flagset.Lookup(database.ArgDSN))

	flagset.StringVar(
		&cmd.databaseConfig.Driver,
		database.ArgDriver,
		"",
		"The database driver to use (e.g., postgres, mysql, sqlite)",
	)
	_ = viper.BindPFlag(database.ArgDriver, flagset.Lookup(database.ArgDriver))

	flagset.IntVar(
		&cmd.databaseConfig.MaxConnections,
		database.ArgMaxConnections,
		database.MaxConnectionsDefault,
		"The maximum number of open connections to the database",
	)
	_ = viper.BindPFlag(database.ArgMaxConnections, flagset.Lookup(database.ArgMaxConnections))

	flagset.IntVar(
		&cmd.databaseConfig.MaxIdleConnections,
		database.ArgMaxIdleConnections,
		database.MaxIdleConnectionsDefault,
		"The maximum number of idle connections in the connection pool",
	)
	_ = viper.BindPFlag(database.ArgMaxIdleConnections, flagset.Lookup(database.ArgMaxIdleConnections))

	return cmd
}

func (c *InitCommand) Execute(cmd *cobra.Command, _ []string) error {
	logger, err := log.NewLogger()
	if err != nil {
		return err
	}

	db, err := database.NewConnection(c.databaseConfig)
	if err != nil {
		return err
	}

	errHandler, err := database.NewErrorHandler(c.databaseConfig)
	if err != nil {
		return err
	}

	authRepo := repository.NewRepository(db, errHandler)
	authService := auth.NewAuth(db, authRepo, logger)

	_, token, err := authService.CreateToken(cmd.Context(), time.Now().Add(24*time.Hour))
	if err != nil {
		return err
	}

	logger.Info("Token created", "token", token)
	//nolint:lll
	logger.Warn("Please save this token in a secure location. It will be used to authenticate your requests to the API and not be shown again.")
	return nil
}
