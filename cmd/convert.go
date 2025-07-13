package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/slug"
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
	opts   WriteItemsOption
}

// WriteItemsOption represents the options for writing items.
type WriteItemsOption struct {
	Output     string // Output specifies the output directory or file path.
	Format     string // Format specifies the format of the output (e.g., JSON, XML).
	Pretty     bool   // Pretty specifies whether to format the output in a human-readable way.
	SingleFile bool   // SingleFile specifies whether to write all items into a single file.
}

func NewConvertCommand() *ConvertCommand {
	c := &ConvertCommand{
		logger: slog.New(
			// TODO: Make handler type configurable.
			tint.NewHandler(os.Stdout, &tint.Options{
				// TODO: Make configurable.
				Level: slog.LevelDebug,
			}),
		),
		opts: WriteItemsOption{},
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

// TODO: This is now coupled to the file system sourcer, but ideally
// once there are multiple sourcer implementation which one to use becomes a flag.

func (c *ConvertCommand) Execute(_ *cobra.Command, args []string) error {
	sourcePath := args[0]
	if err := fs.IsValidFileSystemSource(sourcePath); err != nil {
		return err
	}

	fileSystemSourcer, err := fs.NewSourcer(sourcePath, nil)
	if err != nil {
		return err
	}

	// Iterate over the source files and parse them.
	var items []*api.Item
	for {
		var src sourcer.Source
		src, err = fileSystemSourcer.Next()
		if errors.Is(err, fs.ErrDone) {
			break
		}

		if err != nil {
			return err
		}

		var i *api.Item
		i, err = parser.Parse(src)
		if err != nil {
			c.logger.Warn(fmt.Sprintf("Failed to parse %s: %s", src.Name(), err))
			continue
		}

		items = append(items, i)
	}

	return WriteItems(items, c.opts)
}

func WriteItems(items []*api.Item, opts WriteItemsOption) error {
	marshalFuncs := map[string]marshalerFunc{
		FormatJSON: json.Marshal,
		FormatYAML: yaml.Marshal,
	}

	prettyFuncs := map[string]prettyFunc{
		FormatJSON: pretty.Pretty,
	}

	marshalFunc, ok := marshalFuncs[opts.Format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidFormat, opts.Format)
	}
	prettyFunc := prettyFuncs[opts.Format]

	// Write all items to a single file.
	if opts.SingleFile {
		return writeItems(items, path.Join(opts.Output, "items."+opts.Format), marshalFunc, prettyFunc)
	}

	// Write each item to a separate file.
	for _, i := range items {
		fn := slug.Slug(i.Name)
		path := path.Join(opts.Output, fn+"."+opts.Format)
		if err := writeItems([]*api.Item{i}, path, marshalFunc, prettyFunc); err != nil {
			return err
		}
	}

	return nil
}

type marshalerFunc func(v any) ([]byte, error)
type prettyFunc func(b []byte) []byte

func writeItems(i []*api.Item, path string, marshal marshalerFunc, pretty prettyFunc) error {
	var b []byte
	var err error

	// Marshal the item(s) to the specified format.
	// If there is only one item, marshal it directly, otherwise marshal the slice.
	if len(i) == 1 {
		b, err = marshal(i[0])
	} else {
		b, err = marshal(i)
	}

	if err != nil {
		return err
	}

	// Pretty print the output if requested.
	if pretty != nil {
		b = pretty(b)
	}

	return os.WriteFile(path, b, 0600)
}
