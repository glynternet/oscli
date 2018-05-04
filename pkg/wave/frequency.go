package wave

import "time"

type Frequency float64

func (f Frequency) Period() time.Duration {
	return time.Duration(1.0 / float64(f) * float64(time.Second))
}
