package cmd

import (
	"errors"

	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/spf13/cobra"
)

const (
	ArgSourceType = "source-type"
	ArgLiveMode   = "live"
)

type SyncCommand struct {
	*cobra.Command

	opts SyncCommandOptions
}

type SyncCommandOptions struct {
	SourceType string
	LiveMode   bool
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
		Args:  cobra.ExactArgs(1),
	}

	flagset := cmd.Command.Flags()

	flagset.StringVar(&cmd.opts.SourceType, ArgSourceType, "",
		"The source type to synchronize items from")
	flagset.BoolVar(&cmd.opts.LiveMode, ArgLiveMode, false,
		"When live mode is enabled, items are synchronized to the server, otherwise changes are only previewed")

	return cmd
}

func (c *SyncCommand) RunE(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	logger, err := log.NewLogger()
	if err != nil {
		return err
	}

	if c.opts.SourceType == "" {
		return errors.New("source type is required")
	}

	// Print the source type and live mode.
	logger.InfoContext(ctx, "Synchronizing items", "sourceType", c.opts.SourceType, "liveMode", c.opts.LiveMode)
	return nil
}
