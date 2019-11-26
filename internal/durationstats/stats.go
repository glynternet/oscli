package durationstats

import (
	"math"
	"time"

	"github.com/influxdata/tdigest"
)

type Stats struct {
	Mean      time.Duration
	Quantiles map[float64]time.Duration
}

type Durations []time.Duration

func (ds Durations) Stats(quantilePrecision int) Stats {
	var d tdigest.TDigest
	ds.forEach(func(duration time.Duration) {
		d.Add(float64(duration), 1)
	})
	max := maxxer()
	ds.forEach(max.accept)
	min := minner()
	ds.forEach(min.accept)
	sum := summer()
	ds.forEach(sum.accept)
	qs := make(map[float64]time.Duration, 101)

	// not sure if factor is actually the right name for this
	factor := math.Pow10(quantilePrecision)
	start := 0.9 * factor
	for i := start; i <= factor; i++ {
		quantile := i / factor
		qs[quantile] = time.Duration(d.Quantile(quantile))
	}
	return Stats{
		Mean:      time.Duration(float64(sum.get()) / float64(len(ds))),
		Quantiles: qs,
	}
}

func (ds Durations) forEach(fn func(time.Duration)) {
	for _, d := range ds {
		fn(d)
	}
}

type reducer struct {
	current time.Duration
	result  func(current, new time.Duration) time.Duration
}

func (m *reducer) accept(d time.Duration) {
	m.current = m.result(m.current, d)
}

func (m *reducer) get() time.Duration {
	return m.current
}

func maxxer() reducer {
	return reducer{
		current: -1 << 63,
		result: func(current, new time.Duration) time.Duration {
			if new >= current {
				return new
			}
			return current
		},
	}
}

func minner() reducer {
	return reducer{
		current: 1<<63 - 1,
		result: func(current, new time.Duration) time.Duration {
			if new <= current {
				return new
			}
			return current
		},
	}
}

func summer() reducer {
	return reducer{
		result: func(current, new time.Duration) time.Duration {
			return current + new
		},
	}
}
