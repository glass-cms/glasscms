package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ArgOutput          = "output"
	ArgOutputShorthand = "o"
)

type ParseCommand struct {
	*cobra.Command
}

func NewParseCommand() *ParseCommand {
	c := &ParseCommand{}
	c.Command = &cobra.Command{
		Use:   "parse <source>",
		Short: "Parses source files",
		Long:  "Parses source files and pumps them to the desired destination",
		RunE:  c.Execute,
		Args:  cobra.ExactArgs(1),
	}

	parsePFlagSet := c.Command.PersistentFlags()

	parsePFlagSet.StringP(ArgOutput, ArgOutputShorthand, ".", "Output destination")
	_ = viper.BindPFlag(ArgOutput, parsePFlagSet.Lookup(ArgOutput))

	return c
}

func (c *ParseCommand) Execute(_ *cobra.Command, args []string) error {
	_ = args[0]

	return nil
}
