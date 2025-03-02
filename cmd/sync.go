package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/internal/sync"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/cobra"
)

const (
	ArgSourceType     = "source-type"
	ArgLiveMode       = "live"
	ArgServerURL      = "server"
	ArgToken          = "token"
	ArgHiddenProperty = "hidden-property"
	ArgHiddenValue    = "hidden-value"
)

type SyncCommand struct {
	*cobra.Command

	opts SyncCommandOptions
}

type SyncCommandOptions struct {
	LiveMode       bool
	ServerURL      string
	Token          string
	HiddenProperty string
	HiddenValue    bool
}

// NewSyncCommand returns a new sync command.
func NewSyncCommand() *SyncCommand {
	syncCommand := &SyncCommand{
		opts: SyncCommandOptions{},
	}

	syncCommand.Command = &cobra.Command{
		Use:   "sync [source-type] [source-path]",
		Short: "Synchronize content items from a source to the GlassCMS server",
		Long: heredoc.Doc(`
			Synchronize content items from a source to the GlassCMS API server.

			The sync command allows you to import and update content items from external 
			sources into your GlassCMS instance. It compares the items in the source with 
			those on the server and performs the necessary create, update, or delete operations 
			to keep them in sync.

			Sources are external content repositories that contain structured content items.
			Each source has a specific format and organization, which GlassCMS can interpret
			and import into its content management system.

			Supported source types:
			- filesystem: Read items from a directory on the local filesystem. Items should be
			  organized in a directory structure with JSON or YAML files representing content items.
			  Each file should contain metadata and content according to the GlassCMS schema.

			When run in preview mode (default), the command will show what changes would be made
			without actually applying them. Use the --live flag to apply the changes.
		`),
		Example: heredoc.Doc(`
			# Preview synchronization from a filesystem source
			glasscms sync filesystem /path/to/items

			# Perform live synchronization with server authentication
			glasscms sync filesystem /path/to/items --live --token "your-auth-token"

			# Synchronize to a specific server
			glasscms sync filesystem /path/to/items --server "https://cms.example.com" --token "your-auth-token"

			# Specify a front matter property to determine if an item is hidden
			glasscms sync filesystem /path/to/items --hidden-property "draft" --hidden-value true
		`),
		RunE: syncCommand.RunE,
		Args: cobra.ExactArgs(2),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if _, err := url.Parse(syncCommand.opts.ServerURL); err != nil {
				return fmt.Errorf("invalid server URL: %w", err)
			}
			return nil
		},
	}

	flagset := syncCommand.Command.Flags()

	flagset.BoolVar(&syncCommand.opts.LiveMode, ArgLiveMode, false,
		"When live mode is enabled, items are synchronized to the server, otherwise changes are only previewed")

	flagset.StringVar(&syncCommand.opts.ServerURL, ArgServerURL, "http://localhost:8080",
		"The URL of the server to synchronize items to")

	flagset.StringVar(&syncCommand.opts.Token, ArgToken, "",
		"Bearer token for server authentication")

	flagset.StringVar(&syncCommand.opts.HiddenProperty, ArgHiddenProperty, "",
		"Front matter property name to determine if an item is hidden (e.g., 'draft', 'hidden', 'private')")

	flagset.BoolVar(&syncCommand.opts.HiddenValue, ArgHiddenValue, true,
		`Value of the hidden property that indicates an item is hidden 
		(true = truthy values are hidden, false = falsy values are hidden)`)

	return syncCommand
}

func (c *SyncCommand) RunE(cmd *cobra.Command, args []string) error {
	logger, err := log.NewLogger()
	if err != nil {
		return err
	}

	sr, err := c.initSourcer(args)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	bearerAuth, err := securityprovider.NewSecurityProviderBearerToken(c.opts.Token)
	if err != nil {
		return err
	}

	// Create client with token authentication
	cl, err := api.NewClientWithResponses(c.opts.ServerURL,
		api.WithHTTPClient(httpClient),
		api.WithRequestEditorFn(bearerAuth.Intercept),
	)
	if err != nil {
		return err
	}

	// Create a parser config with the hidden property settings
	parserConfig := parser.Config{
		HiddenProperty: c.opts.HiddenProperty,
		HiddenValue:    c.opts.HiddenValue,
	}

	return sync.NewSyncer(sr, cl, logger, parserConfig).Sync(cmd.Context(), c.opts.LiveMode)
}

// initSourcer initializes a sourcer based on the provided arguments.
// The first argument specifies the source type, and subsequent arguments are source-specific parameters.
// Returns an error if the source type is unrecognized or missing.
func (c *SyncCommand) initSourcer(args []string) (sourcer.Sourcer, error) {
	sourceTypeArg := args[0]
	if sourceTypeArg == "" {
		return nil, errors.New("source type is required")
	}

	sourceType, ok := sourcer.SourceTypeValue[sourceTypeArg]
	if !ok {
		return nil, errors.New("unrecognized source type")
	}

	switch sourceType {
	case sourcer.SourceTypeUnspecified:
		return nil, errors.New("source type is required")
	case sourcer.SourceTypeFilesystem:
		return fs.NewSourcer(args[1])
	}

	return nil, errors.New("unrecognized source type")
}
