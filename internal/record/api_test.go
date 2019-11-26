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
		Schema: "v0.1.0",
		Entries: record.Entries{
			{Duration: 9783662815, Packet: osc.NewMessage("/zoop", "woop", int32(1))},
			{Duration: 21553480721, Packet: osc.NewMessage("/zoop", "woop", float32(0.5))},
			{Duration: 26755783749, Packet: osc.NewMessage("/zoop", "woop", "soup")},
		},
	}

	assertEqualRecording(t, expected, r)

	var outBuf bytes.Buffer
	_, err = r.WriteTo(&outBuf)
	require.NoError(t, err)
	assert.Equal(t, bs, outBuf.Bytes())
}

func assertEqualRecording(t *testing.T, expected, actual record.Recording) {
	assert.Equal(t, expected.Schema, actual.Schema, "expected equal recording schemas")
	assert.Equal(t, len(expected.Entries), len(actual.Entries), "expected equal length arguments")
	for i := range expected.Entries {
		assert.Equal(t, expected.Entries[i].Duration, actual.Entries[i].Duration)
		// We don't support anything other than messages just yet
		require.IsType(t, &osc.Message{}, expected.Entries[i].Packet)
		require.IsType(t, &osc.Message{}, actual.Entries[i].Packet)
		expectedArgs := expected.Entries[i].Packet.(*osc.Message).Arguments
		actualArgs := actual.Entries[i].Packet.(*osc.Message).Arguments
		for j := range expectedArgs {
			assert.Equal(t, expectedArgs[j], actualArgs[j])
		}
	}
}
