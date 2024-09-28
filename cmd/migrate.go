package cmd

import (
	"log/slog"
	"os"

	"github.com/glass-cms/glasscms/internal/database"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type MigrateCommand struct {
	Command *cobra.Command
	logger  *slog.Logger

	databaseConfig database.Config
}

func NewMigrateCommand() *MigrateCommand {
	mc := &MigrateCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level: slog.LevelDebug,
			}),
		),
	}

	mc.Command = &cobra.Command{
		Use:    "migrate",
		Short:  "Migrate the database schema",
		Hidden: true,
		RunE:   mc.Execute,
		Args:   cobra.NoArgs,
	}

	flagset := mc.Command.Flags()

	flagset.StringVar(
		&mc.databaseConfig.Driver,
		database.ArgDriver,
		"",
		"The name of the database driver",
	)
	_ = viper.BindPFlag(database.ArgDriver, flagset.Lookup(database.ArgDriver))

	flagset.StringVar(
		&mc.databaseConfig.DSN,
		database.ArgDSN,
		"",
		"The data source name (DSN) for the database",
	)
	_ = viper.BindPFlag(database.ArgDSN, flagset.Lookup(database.ArgDSN))

	return mc
}

func (mc *MigrateCommand) Execute(_ *cobra.Command, _ []string) error {
	mc.logger.Info("Migrating the database schema")

	db, err := database.NewConnection(mc.databaseConfig)
	if err != nil {
		mc.logger.Error("Failed to create a new database connection")
		return err
	}

	return database.MigrateDatabase(db, mc.databaseConfig)
}
