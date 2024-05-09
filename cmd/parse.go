package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/parser"
	"github.com/glass-cms/glasscms/sourcer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
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
	sourcePath := args[0]
	if err := sourcer.IsValidFileSystemSource(sourcePath); err != nil {
		return err
	}

	fileSystemSourcer, err := sourcer.NewFileSystemSourcer(sourcePath)
	if err != nil {
		return err
	}

	// Iterate over the source files and parse them.
	var items []*item.Item
	for {
		var src sourcer.Source
		src, err = fileSystemSourcer.Next()
		if errors.Is(err, sourcer.ErrDone) {
			break
		}

		if err != nil {
			return err
		}

		var i *item.Item
		i, err = parser.Parse(src)
		if err != nil {
			return err
		}

		items = append(items, i)
	}

	// Write the parsed items to the output destination.
	output := viper.GetString(ArgOutput)
	return writeItems(items, output)
}

func writeItems(items []*item.Item, output string) error {
	// TODO: Make configurable if all items should be written to a single file or multiple files.
	// TODO: Make content type configurable.
	// TODO: Make the filename configurable.
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return err
	}

	// Print the JSON to the console if verbose mode is enabled.
	if viper.GetBool(ArgsVerbose) {
		j := pretty.Pretty(itemsJSON)
		fmt.Println(string(pretty.Color(j, nil)))
	}

	// If the output directory does not exist, create it.
	if _, err = os.Stat(output); os.IsNotExist(err) {
		if err = os.MkdirAll(output, 0755); err != nil {
			return err
		}
	}

	// Combine the output path with the filename.
	fn := output + "/output.json"
	return os.WriteFile(fn, itemsJSON, 0600)
}
