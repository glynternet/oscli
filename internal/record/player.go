package record

import (
	"time"

	"github.com/glynternet/go-osc/osc"
)

type player struct {
	sleepTime time.Duration
}

func (p player) play(es Entries, playEntry func(int, osc.Packet)) {
	start := time.Now()
	es.ForEach(func(i int, e Entry) {
		for time.Since(start) < e.Elapsed {
			time.Sleep(p.sleepTime)
		}
		// Should this spawn a goroutine for each? That could get dangerous.
		playEntry(i, e.Packet)
	})
}
