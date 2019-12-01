package osc

import (
	"context"
	"testing"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	expectedMsg := &osc.Message{}

	for _, tc := range []struct {
		name    string
		count   int
		period  time.Duration
		timeout time.Duration
	}{
		{
			name:    "cancel before first message",
			period:  110 * time.Millisecond,
			timeout: 100 * time.Millisecond,
		},
		{
			name:    "cancel after one message",
			period:  70 * time.Millisecond,
			timeout: 100 * time.Millisecond,
			count:   1,
		},
		{
			name:    "cancel after five messages",
			period:  18 * time.Millisecond,
			timeout: 100 * time.Millisecond,
			count:   5,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var fn MsgGenFunc = func() *osc.Message {
				return expectedMsg
			}
			ctx, cancelFunc := context.WithTimeout(context.Background(), tc.timeout)
			defer cancelFunc()
			start := time.Now()
			ch := Generate(ctx, fn, tc.period)

			var msgs []*osc.Message
			for msg := range ch {
				msgs = append(msgs, msg)
			}
			stop := time.Now()
			for _, msg := range msgs {
				assert.Equal(t, expectedMsg, msg)
			}
			assert.Len(t, msgs, tc.count)
			assert.WithinDuration(t, start.Add(tc.timeout), stop, tc.timeout/10) // fairly lax 10% for CI test runners
		})
	}
}
