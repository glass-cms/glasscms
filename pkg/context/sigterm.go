package context

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SigtermCacellationContext returns a new context that is canceled when an
// interrupt signal (SIGTERM or SIGINT) is received.
// Additionaly invokes the provided onCancel function before canceling the context.
func SigtermCacellationContext(ctx context.Context, onCancel func()) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-interruptCh
		onCancel()
		cancel()
	}()
	return ctx
}
