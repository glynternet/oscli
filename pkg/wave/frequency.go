package wave

import (
	"math"
	"time"
)

// Frequency in Hz
type Frequency float64

// Period returns the Duration in which it takes for 1 cycle of a certain
// frequency to occur.
// For very small Frequency values, behaviour is not considered to be stable.
// The only small Frequency value that is considered is a Frequency of 0Hz,
// where the Duration represented by math.MaxInt64 will be returned.
func (f Frequency) Period() time.Duration {
	if f == 0 {
		return math.MaxInt64
	}
	return time.Duration(1.0 / float64(f) * float64(time.Second))
}
