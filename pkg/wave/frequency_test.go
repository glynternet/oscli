package wave_test

import (
	"math"
	"testing"
	"time"

	"github.com/glynternet/oscli/pkg/wave"
	"github.com/stretchr/testify/assert"
)

func TestFrequency_Period(t *testing.T) {
	for _, test := range []struct {
		name string
		f    wave.Frequency
		p    time.Duration
	}{
		{
			name: "zero-values",
			p:    math.MaxInt64,
		},
		{
			name: "1=1",
			f:    1,
			p:    time.Second,
		},
		{
			name: "negative",
			f:    -0.25,
			p:    time.Second * -4,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			p := test.f.Period()
			assert.Equal(t, test.p, p)
		})
	}
}
