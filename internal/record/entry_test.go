package record

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntry_JSONLoops(t *testing.T) {
	bundle := &osc.Bundle{
		// it would be great to not use time.Local in the tests, which is environment dependent, but the
		// underlying implementation used time.Unix which eventually makes a call to time.Local and it's not
		// worth try to make a workaround to this
		Timetag:  *osc.NewTimetag(time.Date(1999, 1, 1, 1, 1, 1, 1, time.Local)),
		Messages: []*osc.Message{osc.NewMessage("/hiyer", int64(1), "woop")},
	}

	for _, tc := range []struct {
		name string
		Entry
	}{
		{
			name: "zero-value",
		},
		{
			name: "non-zero duration",
			Entry: Entry{
				Elapsed: 1,
			},
		},
		{
			name: "non-zero packet - message",
			Entry: Entry{
				Packet: &osc.Message{
					Address:   "/hiyer",
					Arguments: []interface{}{"wooop"},
				},
			},
		},
		{
			name: "non-zero packet - bundle",
			Entry: Entry{
				Packet: bundle,
			},
		}, {
			name: "non-zero all fields",
			Entry: Entry{
				Elapsed: 1,
				Packet:  bundle,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bs, looped := json.Marshal(tc.Entry)
			require.NoError(t, looped)
			assert.NotNil(t, bs)

			var out Entry
			looped = json.Unmarshal(bs, &out)
			require.NoError(t, looped)
			assert.Equal(t, tc.Entry, out)
		})
	}
}
