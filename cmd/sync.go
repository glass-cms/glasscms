package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/internal/sync"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/lithammer/dedent"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/cobra"
)

const (
	ArgSourceType = "source-type"
	ArgLiveMode   = "live"
	ArgServerURL  = "server"
	ArgToken      = "token"
)

type SyncCommand struct {
	*cobra.Command

	opts SyncCommandOptions
}

type SyncCommandOptions struct {
	LiveMode  bool
	ServerURL string
	Token     string
}

// NewSyncCommand returns a new sync command.
func NewSyncCommand() *SyncCommand {
	syncCommand := &SyncCommand{
		opts: SyncCommandOptions{},
	}

	syncCommand.Command = &cobra.Command{
		Use:   "sync [source-type] [source-path]",
		Short: "Synchronize items from a source to the server",
		Long: dedent.Dedent(`
			Synchronises content items from a source to the server.

			Source types:
			- filesystem: Read items from a directory on the local filesystem.

			Example:
			glasscms sync filesystem /path/to/items
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

	return sync.NewSyncer(sr, cl, logger).Sync(cmd.Context(), c.opts.LiveMode)
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
