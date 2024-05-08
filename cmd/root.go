package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "glass <command>",
	Short: "glass is a CMS for mangaging content based on markdown files",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewParseCommand().Command)
	rootCmd.AddCommand(NewDocsCommand().Command)

	cobra.OnInitialize(intializeConfig)
}

func intializeConfig() {
	viper := viper.New()

	viper.SetEnvPrefix("GLASS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
