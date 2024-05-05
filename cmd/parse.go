package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ParseCommand struct {
	*cobra.Command
}

func NewParseCommand() *ParseCommand {
	c := &ParseCommand{}
	c.Command = &cobra.Command{
		Use:   "parse",
		Short: "Parses source files",
		Long:  "Parses source files and pumps them to the desired destination",
		RunE:  c.Execute,
	}

	return c
}

func (c *ParseCommand) Execute(cmd *cobra.Command, args []string) error {
	fmt.Println("Parsing source files")
	return nil
}
