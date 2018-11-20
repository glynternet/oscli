package models

import (
	"math"
	"time"

	gmath "github.com/glynternet/math"
	osc3 "github.com/glynternet/oscli/internal/osc"
	"github.com/sander/go-osc/osc"
)

// NowSinNormalised returns a MsgGenFunc that, when called, will return a
// message formed with the address and static arguments, appended with the
// value of sin at the time, with a phase calculated as if the epoch was a
// phase of zero, with a frequency freq.
func NowSinNormalised(msgAddr string, staticArgs []interface{}, freq float64) osc3.MsgGenFunc {
	// TODO: revise this documentation
	floatFn := sinNowNormalised(freq)
	return func() *osc.Message {
		mapped := float32(floatFn())
		return osc.NewMessage(msgAddr, append(staticArgs, interface{}(mapped))...)
	}
}

type float64GenFunc func() float64

func sinNowNormalised(freq float64) float64GenFunc {
	mapper := gmath.FloatRangeMapper{
		Source: gmath.FloatRange{From: -1, To: 1},
		Target: gmath.FloatRange{From: 0, To: 1},
	}
	return mappedFloat64GenFunc(mapper, sinNow(freq))
}

func sinNow(waveFreq float64) float64GenFunc {
	return func() float64 {
		unixSecs := float64(time.Now().UnixNano()) * float64(1e-9)
		phase := waveFreq * math.Pi * 2 * unixSecs
		return math.Sin(phase)
	}
}

func mappedFloat64GenFunc(m mapper, fFn float64GenFunc) float64GenFunc {
	return func() float64 {
		return m.Map(fFn())
	}
}

type mapper interface {
	Map(float64) float64
}
