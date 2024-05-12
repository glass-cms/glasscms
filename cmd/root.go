package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
	Use:   "glasscms <command>",
	Short: "glasscms is a CMS for mangaging content based on markdown files",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initializeConfig(cmd)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	//
	// Add commands to the root command.
	//
	rootCmd.AddCommand(NewConvertCommand().Command)
	rootCmd.AddCommand(NewDocsCommand().Command)

	//
	// Register flags.
	//

	pflags := rootCmd.PersistentFlags()

	pflags.BoolP(ArgsVerbose, ArgsVerboseShorthand, false, "Enable verbose output")
	_ = viper.BindPFlag(ArgsVerbose, pflags.Lookup(ArgsVerbose))
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")

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
