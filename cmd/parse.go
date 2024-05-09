package cmd

import (
	"fmt"

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
	viper.BindPFlag(ArgOutput, parsePFlagSet.Lookup(ArgOutput))

	return c
}

func (c *ParseCommand) Execute(_ *cobra.Command, _ []string) error {
	// Get the destination.
	destination := viper.GetString(ArgOutput)
	fmt.Println("Destination:", destination)

	return nil
}
