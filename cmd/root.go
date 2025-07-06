package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/glass-cms/glasscms/cmd/auth"
	"github.com/glass-cms/glasscms/cmd/server"
	"github.com/glass-cms/glasscms/internal/version"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ArgsVerbose          = "verbose"
	ArgsVerboseShorthand = "v"

	defaultConfigFilename = "config"
	envPrefix             = "GLASS"
)

var rootCmd = &cobra.Command{
	Use:          "glasscms <command>",
	Short:        "glasscms is a headless CMS powered by markdown",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initializeConfig(cmd)
	},
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		// If --version flag is used, show version and exit
		if viper.GetBool("version") {
			info := version.Get()
			fmt.Println(info.String())
			return nil
		}
		// Otherwise show help
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewConvertCommand().Command)
	rootCmd.AddCommand(NewDocsCommand().Command)
	rootCmd.AddCommand(server.NewCommand().Command)
	rootCmd.AddCommand(NewMigrateCommand().Command)
	rootCmd.AddCommand(NewSyncCommand().Command)
	rootCmd.AddCommand(auth.NewAuthCommand().Command)
	rootCmd.AddCommand(NewVersionCommand().Command)

	// Register flags.
	pflags := rootCmd.PersistentFlags()

	pflags.BoolP(ArgsVerbose, ArgsVerboseShorthand, false, "Enable verbose output")
	_ = viper.BindPFlag(ArgsVerbose, pflags.Lookup(ArgsVerbose))

	pflags.Bool("version", false, "Show version information")
	_ = viper.BindPFlag("version", pflags.Lookup("version"))

	pflags.String(log.ArgLevel, "INFO", "Log level")
	_ = viper.BindPFlag(log.ArgLevel, pflags.Lookup(log.ArgLevel))

	pflags.String(log.ArgFormat, "TEXT", "Log format")
	_ = viper.BindPFlag(log.ArgFormat, pflags.Lookup(log.ArgFormat))
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")
	// Add glasscms config directory to the search path.

	var cfgNotFoundError viper.ConfigFileNotFoundError
	if err := v.ReadInConfig(); err != nil {
		if !errors.As(err, &cfgNotFoundError) {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable).
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
