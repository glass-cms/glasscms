package auth

import (
	"time"

	"github.com/glass-cms/glasscms/internal/auth"
	"github.com/glass-cms/glasscms/internal/auth/repository"
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/pkg/log"
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
		Short: "Initialize a new token",
		RunE:  cmd.Execute,
	}

	flagset := cmd.Command.Flags()

	flagset.StringVar(
		&cmd.databaseConfig.DSN,
		database.ArgDSN,
		"",
		"The data source name (DSN) for the database",
	)
	_ = viper.BindPFlag(database.ArgDSN, flagset.Lookup(database.ArgDSN))

	flagset.StringVar(
		&cmd.databaseConfig.Driver,
		database.ArgDriver,
		"",
		"The name of the database driver",
	)
	_ = viper.BindPFlag(database.ArgDriver, flagset.Lookup(database.ArgDriver))

	flagset.IntVar(
		&cmd.databaseConfig.MaxConnections,
		database.ArgMaxConnections,
		database.MaxConnectionsDefault,
		"The maximum number of connections that can be opened to the database",
	)
	_ = viper.BindPFlag(database.ArgMaxConnections, flagset.Lookup(database.ArgMaxConnections))

	flagset.IntVar(
		&cmd.databaseConfig.MaxIdleConnections,
		database.ArgMaxIdleConnections,
		database.MaxIdleConnectionsDefault,
		"The maximum number of idle connections that can be maintained",
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
	logger.Warn("Please save this token in a secure location. It will be used to authenticate your requests to the API.")
	return nil
}
