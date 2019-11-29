package osc

import (
	"context"
	"time"

	"github.com/sander/go-osc/osc"
)

// Generate will call the MsgGenFunc periodically after every period given by
// msgPeriod. The resultant message will be sent to the channel that is
// returned by the function.
func Generate(ctx context.Context, fn MsgGenFunc, msgPeriod time.Duration) <-chan *osc.Message {
	ch := make(chan *osc.Message)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-time.After(msgPeriod):
				ch <- fn()
			}
		}
	}()
	return ch
}

// MsgGenFunc is a function that returns an osc.Message when called, for
// passing to the Generate function.
type MsgGenFunc func() *osc.Message
