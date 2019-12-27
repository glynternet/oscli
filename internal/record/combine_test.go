package record_test

import (
	"testing"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const expectedVersion = "v0.1.0"

func TestCombine_schemaValidation(t *testing.T) {
	for _, tc := range []struct {
		name          string
		recordings    []record.Recording
		invalidSchema string
	}{{
		name: "zero-value",
	}, {
		name:       "single valid schema",
		recordings: []record.Recording{{Schema: expectedVersion}},
	}, {
		name:          "single invalid schema",
		recordings:    []record.Recording{{Schema: "babbedyboopedy"}},
		invalidSchema: "babbedyboopedy",
	}, {
		name: "valid followed by invalid schema",
		recordings: []record.Recording{
			{Schema: expectedVersion},
			{Schema: expectedVersion},
			{Schema: "babbedyboopedy"},
		},
		invalidSchema: "babbedyboopedy",
	}, {
		name: "return first invalid schema",
		recordings: []record.Recording{
			{Schema: expectedVersion},
			{Schema: expectedVersion},
			{Schema: "babbedyboopedy"},
			{Schema: "foopedywoopedy"},
		},
		invalidSchema: "babbedyboopedy",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := record.Combine(tc.recordings...)
			if tc.invalidSchema == "" {
				assert.Equal(t, expectedVersion, out.Schema)
				return
			}
			assert.Empty(t, out)
			assert.Equal(t, record.UnsupportedSchemaError(tc.invalidSchema), err)
		})
	}
}

func TestCombine_output(t *testing.T) {
	for _, tc := range []struct {
		name       string
		recordings []record.Recording
		out        record.Recording
	}{{
		name: "zero-values",
		out:  record.Recording{Schema: expectedVersion},
	}, {
		name:       "single recording returns self but sorted",
		recordings: []record.Recording{recording(100, 90)},
		out:        recording(90, 100),
	}, {
		name: "multiple recordings with non-overlapping durations are concatenated",
		recordings: []record.Recording{
			recording(100, 90),
			recording(110, 120)},
		out: recording(90, 100, 110, 120),
	}, {
		name: "multiple recordings with overlapping durations are combined and sorted",
		recordings: []record.Recording{
			recording(90, 120),
			recording(110, 100),
			recording(80, 200, 110)},
		out: recording(80, 90, 100, 110, 110, 120, 200),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := record.Combine(tc.recordings...)
			require.NoError(t, err)
			assert.Equal(t, tc.out, out)
		})
	}
}

func recording(entryTimes ...time.Duration) record.Recording {
	var es record.Entries
	for _, entryTime := range entryTimes {
		es = append(es, record.Entry{
			Duration: entryTime,
			Packet:   osc.NewMessage("/" + entryTime.String()),
		})
	}
	return record.Recording{
		Schema:  expectedVersion,
		Entries: es,
	}
}
