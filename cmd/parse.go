package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ArgDestination          = "destination"
	ArgDestinationShorthand = "d"
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

	parsePFlagSet := c.Command.Flags()

	parsePFlagSet.StringP(ArgDestination, ArgDestinationShorthand, ".", "Destination to write the parsed files to")
	viper.BindPFlag(ArgDestination, parsePFlagSet.Lookup(ArgDestination))

	return c
}

func (c *ParseCommand) Execute(_ *cobra.Command, _ []string) error {
	// Get the destination.
	destination := viper.GetString(ArgDestination)
	fmt.Println("Destination:", destination)

	return nil
}
