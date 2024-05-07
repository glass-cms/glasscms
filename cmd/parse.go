package cmd

import (
	"github.com/spf13/cobra"
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

	return c
}

func (c *ParseCommand) Execute(_ *cobra.Command, _ []string) error {
	return nil
}
