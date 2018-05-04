package math_test

import (
	"math"
	"testing"

	gohmath "github.com/GlynOwenHanmer/math"
)

func TestFloatRangeMapper_Map(t *testing.T) {
	tests := []struct {
		gohmath.FloatRangeMapper
		io map[float64]float64
	}{
		{
			FloatRangeMapper: gohmath.FloatRangeMapper{
				Source: gohmath.FloatRange{1, 1},
				Target: gohmath.FloatRange{2, 2},
			},
			io: map[float64]float64{
				0.5: 2,
				1:   2,
				1.5: 2,
			},
		},
		{
			FloatRangeMapper: gohmath.FloatRangeMapper{
				Source: gohmath.FloatRange{1, 1},
				Target: gohmath.FloatRange{2, 5},
			},
			io: map[float64]float64{
				0.5: math.Inf(-1),
				1:   3.5,
				1.5: math.Inf(1),
			},
		},
		{
			FloatRangeMapper: gohmath.FloatRangeMapper{
				Source: gohmath.FloatRange{1, 1},
				Target: gohmath.FloatRange{5, 2},
			},
			io: map[float64]float64{
				0.5: math.Inf(1),
				1:   3.5,
				1.5: math.Inf(-1),
			},
		},
		{
			FloatRangeMapper: gohmath.FloatRangeMapper{
				Source: gohmath.FloatRange{1, 2},
				Target: gohmath.FloatRange{2, 4},
			},
			io: map[float64]float64{
				0.5: 1,
				1:   2,
				1.5: 3,
				2:   4,
				2.5: 5,
			},
		},
		{
			FloatRangeMapper: gohmath.FloatRangeMapper{
				Source: gohmath.FloatRange{2, 4},
				Target: gohmath.FloatRange{2, -2},
			},
			io: map[float64]float64{
				1: 4,
				2: 2,
				3: 0,
				4: -2,
				5: -4,
			},
		},
	}
	for _, test := range tests {
		for value, expected := range test.io {
			actual := test.FloatRangeMapper.Map(value)
			if actual != expected {
				t.Errorf("Expected %f, got %f", expected, actual)
			}
		}
	}
}
