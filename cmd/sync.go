package cmd

import (
	"errors"
	"net/http"
	"time"

	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/internal/sync"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/spf13/cobra"
)

const (
	ArgSourceType = "source-type"
	ArgLiveMode   = "live"
	ArgServerURL  = "server"
)

type SyncCommand struct {
	*cobra.Command

	opts SyncCommandOptions
}

type SyncCommandOptions struct {
	LiveMode  bool
	ServerURL string
}

// NewSyncCommand returns a new sync command.
func NewSyncCommand() *SyncCommand {
	cmd := &SyncCommand{
		opts: SyncCommandOptions{},
	}

	cmd.Command = &cobra.Command{
		Use:   "sync",
		Short: "Synchronize items from a source to the server",
		RunE:  cmd.RunE,
		Args:  cobra.ExactArgs(2),
	}

	flagset := cmd.Command.Flags()

	flagset.BoolVar(&cmd.opts.LiveMode, ArgLiveMode, false,
		"When live mode is enabled, items are synchronized to the server, otherwise changes are only previewed")

	flagset.StringVar(&cmd.opts.ServerURL, ArgServerURL, "http://localhost:8080",
		"The URL of the server to synchronize items to")

	return cmd
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

	// TODO: The client should also have logging middleware and retry logic.

	cl, err := api.NewClientWithResponses(c.opts.ServerURL, api.WithHTTPClient(httpClient))
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