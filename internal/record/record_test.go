package record

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"

	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/durationstats"
	"github.com/glynternet/oscli/internal/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/glynternet/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: setup benchmarking suite somehow
//   this should probably be something that also runs over quite a long period of time. Say ten minutes, at least.
//   Given that this could be for long performance things
// TODO: benchmarking could output some structured blob of data which can be analysed in python
// TODO: profile this
// TODO: setup profiling, too. It will be important to profile this when we are benchmarking it

func TestRecord(t *testing.T) {
	port := 9000
	addr := "localhost:" + strconv.Itoa(port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c := osc2.NewClient("localhost", port)
	logger := log.NewLogger()

	countCh := make(chan int)
	// WARNING: could be a race here where the first message is send before recording starts
	go generateMessages(ctx, logger, c, wave.Frequency(1000), countCh)

	recorded, err := Record(ctx, logger, addr)
	require.NoError(t, err)
	assert.Equal(t, <-countCh, recorded.entryCount())

	var buf bytes.Buffer
	_, err = recorded.WriteTo(&buf)
	require.NoError(t, err)

	var decoded Recording
	_, err = decoded.ReadFrom(&buf)
	require.NoError(t, err)
	require.Equal(t, recorded.entryCount(), decoded.entryCount())

	replayed := replayPackets(logger, decoded.Data.Entries)
	require.Equal(t, decoded.entryCount(), len(replayed))

	diffs := getDiffs(recorded.Data.Entries, replayed)

	stats := diffs.Stats(2)
	assertQuantileTime(t, stats, 0.99, 10*time.Millisecond)
	assertQuantileTime(t, stats, 0.95, 5*time.Millisecond)
	assertQuantileTime(t, stats, 0.9, time.Millisecond)
}

func assertQuantileTime(t *testing.T, stats durationstats.Stats, quantile float64, d time.Duration) bool {
	duration, ok := stats.Quantiles[quantile]
	if !ok {
		t.Errorf("quantile %f does not exist", quantile)
		return false
	}
	return assert.Truef(t, duration < d, "%f quantile should be under %v but is %v", quantile, d, duration)
}

func replayPackets(logger log.Logger, recorded Entries) []time.Duration {
	var replayed []time.Duration
	copyStart := time.Now()
	recorded.Play(func(_ int, _ osc2.Packet) {
		replayed = append(replayed, time.Since(copyStart))
	})
	_ = logger.Log(log.Message("finished replaying"))
	return replayed
}

func generateMessages(ctx context.Context, logger log.Logger, client *osc2.Client, frequency wave.Frequency, countCh chan<- int) {
	genFn := func() *osc2.Message {
		return &osc2.Message{
			Address:   "/whoop",
			Arguments: []interface{}{int64(1), int64(2), int64(3)},
		}
	}

	_ = logger.Log(log.Message("starting generator"))
	msgCh := osc.Generate(ctx, genFn, frequency.Period())
	var count int
	for msg := range msgCh {
		count++
		if err := client.Send(msg); err != nil {
			_ = logger.Log(log.Message("Error sending packet"),
				log.Error(err),
				log.KV{K: "packet", V: msg})
			return
		}
	}
	countCh <- count
}
