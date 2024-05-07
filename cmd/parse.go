package cmd

import (
	"fmt"

	"github.com/glass-cms/glasscms/sourcer"
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

func (c *ParseCommand) Execute(cmd *cobra.Command, args []string) error {
	path := args[0]

	source, err := sourcer.NewFileSystemSourcer(path)
	fmt.Println(source.Size())

	return err
}
