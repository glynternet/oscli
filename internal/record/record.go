package record

import (
	"context"
	"encoding/json"
	"io"
	"time"

	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/osc"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
)

var version = "v0.1.0"

type Recording struct {
	Schema  string  `json:"schema"`
	Entries Entries `json:"entries"`
}

func (rd Recording) entryCount() int {
	return rd.Entries.Len()
}
func (rd Recording) WriteTo(w io.Writer) (int64, error) {
	return 0, json.NewEncoder(w).Encode(&rd)
}

func (rd *Recording) ReadFrom(r io.Reader) (int64, error) {
	return 0, json.NewDecoder(r).Decode(rd)
}

func Record(ctx context.Context, logger log.Logger, addr string) (Recording, error) {
	return recorder{
		start:   time.Now(),
		address: addr,
	}.record(ctx, logger)
}

type recorder struct {
	start   time.Time
	address string
}

func (r recorder) record(ctx context.Context, logger log.Logger) (Recording, error) {
	var recorded Entries
	err := osc.ReceivePackets(ctx, logger, r.address, func(packet osc2.Packet) {
		since := time.Since(r.start)
		recorded.add(Entry{
			Duration: since,
			Packet:   packet},
		)
	}, func(err error) {
		_ = logger.Log(
			log.Message("Error while receiving packet from connection"),
			log.Error(err))
	})
	return Recording{
		Schema:  version,
		Entries: recorded,
	}, errors.Wrap(err, "receiving packets")
}
