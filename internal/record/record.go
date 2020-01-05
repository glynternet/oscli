package record

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/osc"
	"github.com/pkg/errors"
)

var version = "v0.1.0"

type Recording struct {
	Data RecordingData
}

type RecordingData struct {
	Entries Entries `json:"entries"`
}

func (rd Recording) MarshalJSON() ([]byte, error) {
	return json.Marshal(serialisedRecording{
		Schema: version,
		Data:   rd.Data,
	})
}

func (rd *Recording) UnmarshalJSON(bs []byte) error {
	var sr serialisedRecording
	if err := json.Unmarshal(bs, &sr); err != nil {
		return errors.Wrap(err, "decoding into intermediate recording type")
	}
	if sr.Schema != version {
		return UnsupportedSchemaError(sr.Schema)
	}
	*rd = Recording{Data: sr.Data}
	return nil
}

type serialisedRecording struct {
	Schema string        `json:"schema"`
	Data   RecordingData `json:"data"`
}

func (rd Recording) entryCount() int {
	return rd.Data.Entries.Len()
}

func (rd Recording) WriteTo(w io.Writer) (int64, error) {
	return 0, json.NewEncoder(w).Encode(&serialisedRecording{
		Schema: version,
		Data:   rd.Data,
	})
}

func (rd *Recording) ReadFrom(r io.Reader) (int64, error) {
	return 0, json.NewDecoder(r).Decode(&rd)
}

func Record(ctx context.Context, logger *log.Logger, addr string) (Recording, error) {
	return recorder{
		start:   time.Now(),
		address: addr,
	}.record(ctx, logger)
}

type recorder struct {
	start   time.Time
	address string
}

func (r recorder) record(ctx context.Context, logger *log.Logger) (Recording, error) {
	var recorded Entries
	err := osc.ReceivePackets(ctx, logger, r.address, func(packet osc2.Packet) {
		since := time.Since(r.start)
		recorded.add(Entry{
			Duration: since,
			Packet:   packet},
		)
	}, func(err error) {
		logger.Printf("Error while receiving packet from connection: %v", err)
	})
	return Recording{
		Data: RecordingData{Entries: recorded},
	}, errors.Wrap(err, "receiving packets")
}

// UnsupportedSchemaError is the error type returned when a Recording with an unsupported schema is encountered
type UnsupportedSchemaError string

func (us UnsupportedSchemaError) Error() string {
	return fmt.Sprintf("invalid schema: %s", string(us))
}
