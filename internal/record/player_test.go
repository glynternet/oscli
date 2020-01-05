package record

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/durationstats"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/stretchr/testify/require"
)

func TestPlayer(t *testing.T) {
	// TODO(glynternet): actually implement automated benchmarking
	//totalLength := 180 * time.Second
	//for totalLength < time.Minute {
	//runBench(t, totalLength)
	//totalLength = totalLength * 2
	//}
}

// tuple of sleep duration and the durationstats.Stats for that given run
type result struct {
	sleep time.Duration
	durationstats.Stats
}

func runBench(t *testing.T, totalLength time.Duration) {
	sleep := 10 * time.Nanosecond
	quantilePrecision := 4
	var results []result
	for sleep < 100*time.Nanosecond {
		t.Run(fmt.Sprintf("totalLength:%s sleep:%s", totalLength, sleep), func(t *testing.T) {
			stats := testPlayerSleep(generateEntries(totalLength, wave.Frequency(500).Period()), sleep, quantilePrecision)
			t.Logf("%+v", stats)
			results = append(results, result{sleep: sleep, Stats: stats})
		})
		inc := time.Duration(float64(sleep) * 0.01)
		// if sleep was 1, inc evaluates to 0 and would eternally loop
		if inc == 0 {
			inc = 1
		}
		sleep += inc
	}
	filename := fmt.Sprintf("./length_%d.csv", totalLength)
	writeResultsToFile(t, filename, results, quantilePrecision)
	t.Logf("file written to %s", filename)
}

func writeResultsToFile(t *testing.T, filename string, results []result, quantilePrecision int) {
	f, err := os.Create(filename)
	require.NoError(t, err)
	w := csv.NewWriter(f)
	defer func() {
		w.Flush()
		require.NoError(t, f.Close(), "closing results file")
	}()
	for _, r := range results {
		sleepStr := strconv.Itoa(int(r.sleep))
		for q, duration := range r.Quantiles {
			require.NoError(t,
				w.Write([]string{sleepStr,
					strconv.FormatFloat(q, 'f', quantilePrecision, 32),
					strconv.Itoa(int(duration))}))
		}
	}
}

// testPlayerSleep plays the given entries, records the time at which they are replayed, returning
// the Stats produced from the difference between the original play times and the resultant play times.
func testPlayerSleep(es Entries, sleep time.Duration, quantilePrecision int) durationstats.Stats {
	var played []time.Duration
	start := time.Now()
	player{sleepTime: sleep}.play(es, func(i int, packet osc.Packet) {
		played = append(played, time.Since(start))
	})
	return getDiffs(es, played).Stats(quantilePrecision)
}

// generateEntries creates a set of periodic entries that total no longer than the totalTime
func generateEntries(totalTime, period time.Duration) Entries {
	var es Entries
	var current time.Duration
	for current <= totalTime {
		es.add(Entry{Duration: current})
		current += period
	}
	return es
}

// getDiffs calculates the durations that are the difference between the given entries and
// the set of replayed times that were recorded when the entries were replayed
func getDiffs(es Entries, replayed []time.Duration) durationstats.Durations {
	var diffs durationstats.Durations
	es.ForEach(func(i int, e Entry) {
		diffs = append(diffs, replayed[i]-es[i].Duration)
	})
	return diffs
}
