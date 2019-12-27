package record

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/pkg/errors"
)

type Entry struct {
	time.Duration
	osc.Packet
}

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

type Entries []Entry

func (es Entries) Len() int {
	return len(es)
}

func (es Entries) Less(i, j int) bool {
	return es[i].Duration < es[j].Duration
}

func (es Entries) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es *Entries) add(e Entry) {
	*es = append(*es, e)
}

func (es Entries) Count() int {
	return len(es)
}

func (es Entries) forEach(fn func(int, Entry)) {
	for i, e := range es {
		fn(i, e)
	}
}

func (es Entries) Play(playEntry func(int, osc.Packet)) {
	player{sleepTime: 5}.play(es, playEntry)
}
