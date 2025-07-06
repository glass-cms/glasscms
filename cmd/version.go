package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/glass-cms/glasscms/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ArgVersionFormat          = "format"
	ArgVersionFormatShorthand = "f"
)

type VersionCommand struct {
	Command *cobra.Command
}

func NewVersionCommand() *VersionCommand {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information including version, commit hash, build date, Go version, and platform",
		RunE:  runVersion,
	}

	flags := cmd.Flags()
	flags.StringP(ArgVersionFormat, ArgVersionFormatShorthand, "text", "Output format (text, json)")
	_ = viper.BindPFlag(ArgVersionFormat, flags.Lookup(ArgVersionFormat))

	return &VersionCommand{
		Command: cmd,
	}
}

func runVersion(_ *cobra.Command, _ []string) error {
	format := viper.GetString(ArgVersionFormat)
	info := version.Get()

	switch format {
	case "json":
		output, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal version info: %w", err)
		}
		fmt.Println(string(output))
	case "text":
		fmt.Println(info.String())
	default:
		return fmt.Errorf("unsupported format: %s (supported: text, json)", format)
	}

	return nil
}
