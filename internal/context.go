package internal

import (
	"context"
	"os"
	"os/signal"
)

func ContextWithSignalCancels(ctx context.Context, ss ...os.Signal) (context.Context, context.CancelFunc) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, ss...)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-sigs
		// does it matter that cancel will be called twice?
		cancel()
	}()
	return ctx, cancel
}
