package osc

import (
	"time"

	"github.com/hypebeast/go-osc/osc"
)

// Generate will call the MsgGenFunc periodically after every period given by
// msgPeriod. The resultant message will be sent to the channel that is
// returned by the function.
func Generate(fn MsgGenFunc, msgPeriod time.Duration) <-chan *osc.Message {
	ch := make(chan *osc.Message)
	go func() {
		// TODO: what's the best way to kill this goroutine?
		for {
			ch <- fn()
			time.Sleep(msgPeriod)
		}
	}()
	return ch
}

// MsgGenFunc is a function that returns an osc.Message when called, for
// passing to the Generate function.
type MsgGenFunc func() *osc.Message
