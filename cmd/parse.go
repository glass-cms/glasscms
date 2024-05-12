package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/parser"
	"github.com/glass-cms/glasscms/sourcer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ArgOutput          = "output"
	ArgOutputShorthand = "o"

	ArgFormat          = "format"
	ArgFormatShorthand = "f"

	FormatJSON = "json"
	FormatYAML = "yaml"
)

var (
	ErrArgumentInvalid = errors.New("argument is invalid")
	ErrInvalidFormat   = errors.New("invalid format")
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			dir := viper.GetString(ArgOutput)

			// Create the output directory if it doesn't exist.
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
			}

			format := viper.GetString(ArgFormat)
			if format != FormatJSON && format != FormatYAML {
				return fmt.Errorf("%w: %s", ErrArgumentInvalid, format)
			}

			return nil
		},
	}

	flagset := c.Command.Flags()

	flagset.StringP(ArgOutput, ArgOutputShorthand, ".", "Output directory")
	_ = viper.BindPFlag(ArgOutput, flagset.Lookup(ArgOutput))

	flagset.StringP(ArgFormat, ArgFormatShorthand, "json", "Output format (json, yaml)")
	_ = viper.BindPFlag(ArgFormat, flagset.Lookup(ArgFormat))

	return c
}

func (c *ParseCommand) Execute(_ *cobra.Command, args []string) error {
	sourcePath := args[0]
	if err := sourcer.IsValidFileSystemSource(sourcePath); err != nil {
		return err
	}

	dir := viper.GetString(ArgOutput)
	format := viper.GetString(ArgFormat)

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

	return writeItems(items, dir, format)
}

func writeItems(items []*item.Item, dir string, format string) error {
	for _, i := range items {
		fn := strings.ReplaceAll(i.Name, "/", "_")

		switch format {
		case FormatJSON:
			if err := writeItemJson(i, path.Join(dir, fn+".json")); err != nil {
				return err
			}
		case FormatYAML:
			if err := writeItemYaml(i, path.Join(dir, fn+".yaml")); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %s", ErrInvalidFormat, format)
		}
	}

	return nil
}

func writeItemJson(i *item.Item, path string) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0600)
}

func writeItemYaml(i *item.Item, path string) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0600)
}
