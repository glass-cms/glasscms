package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "glcms <command>",
	Short: "glcms is a CMS for mangaging content based on markdown files",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewParseCommand().Command)
	rootCmd.AddCommand(NewDocsCommand().Command)
}
