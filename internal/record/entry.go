package record

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/pkg/errors"
)

// Entry is a single packet and the elapsed time that it was recorded during a recording
type Entry struct {
	time.Duration
	osc.Packet
}

// MarshalJSON  marshals a given Entry into JSON format
func (e Entry) MarshalJSON() ([]byte, error) {
	var packet *string
	if e.Packet != nil {
		pBin, err := e.Packet.MarshalBinary()
		if err != nil {
			return nil, errors.Wrap(err, "marshalling packet into binary")
		}
		encoded := base64.StdEncoding.EncodeToString(pBin)
		packet = &encoded
	}
	return json.Marshal(entryJSONAlias{
		Duration:  e.Duration,
		B64Packet: packet,
	})
}

// UnmarshalJSON  unmarshals a given Entry from JSON format
func (e *Entry) UnmarshalJSON(data []byte) error {
	var alias entryJSONAlias
	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}
	if alias.B64Packet == nil {
		e.Duration = alias.Duration
		return nil
	}
	decoded, err := base64.StdEncoding.DecodeString(*alias.B64Packet)
	if err != nil {
		return errors.Wrap(err, "decoding packet base64")
	}

	packet, err := osc.ParsePacket(string(decoded))
	if err != nil {
		return errors.Wrap(err, "parsing decoded string as osc packet")
	}
	*e = Entry{
		Duration: alias.Duration,
		Packet:   packet,
	}
	return nil
}

type entryJSONAlias struct {
	time.Duration `json:"duration"`
	B64Packet     *string `json:"packet"`
}

// Entries is a group of Entries
type Entries []Entry

// Len returns the length of an Entries
func (es Entries) Len() int {
	return len(es)
}

// Less returns whether Entry at index i is less than the Entry at index j
func (es Entries) Less(i, j int) bool {
	return es[i].Duration < es[j].Duration
}

// Swap swaps the entries at the given indices, i an j.
func (es Entries) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es *Entries) add(e Entry) {
	*es = append(*es, e)
}

// Call the function, fn, on each Entry
func (es Entries) ForEach(fn func(int, Entry)) {
	for i, e := range es {
		fn(i, e)
	}
}

// Play plays the Entries in realtime, calling the playEntry function on each one at it's original elapsed time
func (es Entries) Play(playEntry func(int, osc.Packet)) {
	player{sleepTime: 5}.play(es, playEntry)
}
