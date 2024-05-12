package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/parser"
	"github.com/glass-cms/glasscms/sourcer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
)

const (
	ArgOutput            = "output"
	ArgOutputShorthand   = "o"
	ArgFilename          = "filename"
	ArgFilenameShorthand = "f"
)

var (
	ErrArgumentInvalid = errors.New("argument is invalid")
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

	flagset := c.Command.Flags()

	flagset.StringP(ArgOutput, ArgOutputShorthand, ".", "Output destination")
	_ = viper.BindPFlag(ArgOutput, flagset.Lookup(ArgOutput))

	flagset.StringP(ArgFilename, ArgFilenameShorthand, "output.json", "Output filename")
	_ = viper.BindPFlag(ArgFilename, flagset.Lookup(ArgFilename))

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
	outputDir := viper.GetString(ArgOutput)
	if err = createOutputDir(outputDir); err != nil {
		return err
	}

	filename, err := filename()
	if err != nil {
		return err
	}

	path := path.Join(outputDir, filename)
	return writeItems(items, path)
}

func writeItems(items []*item.Item, filepath string) error {
	// TODO: Make configurable if all items should be written to a single file or multiple files.
	// TODO: Make content type configurable.
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return err
	}

	if viper.GetBool(ArgsVerbose) {
		j := pretty.Pretty(itemsJSON)
		fmt.Println(string(pretty.Color(j, nil)))
	}

	return os.WriteFile(filepath, itemsJSON, 0600)
}

func createOutputDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}

	return nil
}

func filename() (string, error) {
	arg := viper.GetString(ArgFilename)
	if path.Ext(arg) != "" {
		return arg, fmt.Errorf("%w: %s", ErrArgumentInvalid, "filename must not contain an extension")
	}

	return arg, nil
}
