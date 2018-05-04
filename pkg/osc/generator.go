package osc

import (
	"time"

	"github.com/hypebeast/go-osc/osc"
)

func Generate(fn MsgGenFunc, msgPeriod time.Duration) <-chan *osc.Message {
	ch := make(chan *osc.Message)
	go func() {
		// what's the best way to kill this goroutine?
		for {
			ch <- fn()
			time.Sleep(msgPeriod)
		}
	}()
	return ch
}

type MsgGenFunc func() *osc.Message
