package record_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	bs, err := ioutil.ReadFile("./api/v0.1.0.osc")
	require.NoError(t, err)

	inBuf := bytes.NewBuffer(bs)
	var r record.Recording
	_, err = r.ReadFrom(inBuf)
	require.NoError(t, err)

	expected := record.Recording{
		Data: record.RecordingData{
			Entries: record.Entries{
				{Elapsed: 9783662815, Packet: osc.NewMessage("/zoop", "woop", int32(1))},
				{Elapsed: 21553480721, Packet: osc.NewMessage("/zoop", "woop", float32(0.5))},
				{Elapsed: 26755783749, Packet: osc.NewMessage("/zoop", "woop", "soup")},
			},
		}}

	assertEqualRecording(t, expected, r)

	var outBuf bytes.Buffer
	_, err = r.WriteTo(&outBuf)
	require.NoError(t, err)
	assert.Equal(t, bs, outBuf.Bytes())
}

func assertEqualRecording(t *testing.T, expected, actual record.Recording) {
	assert.Equal(t, len(expected.Data.Entries), len(actual.Data.Entries), "expected equal length entries")
	for i := range expected.Data.Entries {
		assert.Equal(t, expected.Data.Entries[i].Elapsed, actual.Data.Entries[i].Elapsed)
		// We don't support anything other than messages just yet
		require.IsType(t, &osc.Message{}, expected.Data.Entries[i].Packet)
		require.IsType(t, &osc.Message{}, actual.Data.Entries[i].Packet)
		expectedArgs := expected.Data.Entries[i].Packet.(*osc.Message).Arguments
		actualArgs := actual.Data.Entries[i].Packet.(*osc.Message).Arguments
		for j := range expectedArgs {
			assert.Equal(t, expectedArgs[j], actualArgs[j])
		}
	}
}
