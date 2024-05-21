package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/parser"
	"github.com/glass-cms/glasscms/sourcer"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
)

const (
	ArgOutput          = "output"
	ArgOutputShorthand = "o"

	ArgFormat          = "format"
	ArgFormatShorthand = "f"

	ArgPretty          = "pretty"
	ArgPrettyShorthand = "p"

	ArgSingleFile          = "single-file"
	ArgSingleFileShorthand = "s"

	FormatJSON = "json"
	FormatYAML = "yaml"

	TitleProperty = "title"
)

var (
	ErrArgumentInvalid = errors.New("argument is invalid")
	ErrInvalidFormat   = errors.New("invalid format")
)

type ConvertCommand struct {
	*cobra.Command

	logger *slog.Logger
	opts   ConvertCommandOptions
}

type ConvertCommandOptions struct {
	Output     string
	Format     string
	Pretty     bool
	SingleFile bool
}

func NewConvertCommand() *ConvertCommand {
	c := &ConvertCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.TimeOnly,
			}),
		),
		opts: ConvertCommandOptions{},
	}

	c.Command = &cobra.Command{
		Use:   "convert <source>",
		Short: "Convert source files",
		Long:  "Convert source files to a structured format at the specified output directory.",
		RunE:  c.Execute,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			dir := viper.GetString(ArgOutput)

			// Create the output directory if it doesn't exist.
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err = os.MkdirAll(dir, 0755); err != nil {
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

	flagset.StringVarP(&c.opts.Output, ArgOutput, ArgOutputShorthand, ".", "Output directory")
	_ = viper.BindPFlag(ArgOutput, flagset.Lookup(ArgOutput))

	flagset.StringVarP(&c.opts.Format, ArgFormat, ArgFormatShorthand, "json", "Output format (json, yaml)")
	_ = viper.BindPFlag(ArgFormat, flagset.Lookup(ArgFormat))

	flagset.BoolVarP(&c.opts.Pretty, ArgPretty, ArgPrettyShorthand, false, "Pretty print output")
	_ = viper.BindPFlag("pretty", flagset.Lookup("pretty"))

	flagset.BoolVarP(&c.opts.SingleFile, ArgSingleFile, ArgSingleFileShorthand, false, "Write all items to a single file")
	_ = viper.BindPFlag(ArgSingleFile, flagset.Lookup(ArgSingleFile))

	return c
}

func (c *ConvertCommand) Execute(_ *cobra.Command, args []string) error {
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
			c.logger.Warn(fmt.Sprintf("Failed to parse %s: %s", src.Name(), err))
			continue
		}

		items = append(items, i)
	}

	return writeItems(items, c.opts)
}

func writeItems(items []*item.Item, opts ConvertCommandOptions) error {
	for _, i := range items {
		fn := i.Name
		if title := i.Title(); title != nil {
			fn = *title
		}

		switch opts.Format {
		case FormatJSON:
			if err := writeItemJSON(i, path.Join(opts.Output, fn+".json"), opts.Pretty); err != nil {
				return err
			}
		case FormatYAML:
			if err := writeItemYAML(i, path.Join(opts.Output, fn+".yaml")); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %s", ErrInvalidFormat, opts.Format)
		}
	}

	return nil
}

func writeItemJSON(i *item.Item, path string, format bool) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	if format {
		b = pretty.Pretty(b)
	}

	return os.WriteFile(path, b, 0600)
}

func writeItemsJSON(items []*item.Item, path string, format bool) error {
	b, err := json.Marshal(items)
	if err != nil {
		return err
	}

	if format {
		b = pretty.Pretty(b)
	}

	return os.WriteFile(path, b, 0600)
}

func writeItemYAML(i *item.Item, path string) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0600)
}

func writeItemsYAML(items []*item.Item, path string) error {
	b, err := yaml.Marshal(items)
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0600)
}
