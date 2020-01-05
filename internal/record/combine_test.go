package record_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/record"
	"github.com/stretchr/testify/assert"
)

func TestCombine_output(t *testing.T) {
	for _, tc := range []struct {
		name    string
		entries []record.Entries
		out     record.Entries
	}{{
		name: "zero-values",
		out:  record.Entries{},
	}, {
		name:    "single entries returns self but sorted",
		entries: []record.Entries{entries(100, 90)},
		out:     entries(90, 100),
	}, {
		name: "multiple entries with non-overlapping durations are concatenated",
		entries: []record.Entries{
			entries(100, 90),
			entries(110, 120)},
		out: entries(90, 100, 110, 120),
	}, {
		name: "multiple entries with overlapping durations are combined and sorted",
		entries: []record.Entries{
			entries(90, 120),
			entries(110, 100),
			entries(80, 200, 110)},
		out: entries(80, 90, 100, 110, 110, 120, 200),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			out := record.Combine(tc.entries...)
			assert.Equal(t, tc.out, out)
		})
	}
}

func entries(times ...time.Duration) record.Entries {
	var es record.Entries
	for _, t := range times {
		es = append(es, record.Entry{
			Elapsed: t,
			Packet:  osc.NewMessage("/" + t.String()),
		})
	}
	return es
}
