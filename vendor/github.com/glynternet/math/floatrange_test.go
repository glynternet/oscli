package math_test

import (
	"testing"

	"github.com/GlynOwenHanmer/math"
)

func TestFloatrange_Range(t *testing.T) {
	tests := []struct{ from, to, expected float64 }{
		{from: 0, to: 1, expected: 1},
		{from: -1, to: -2, expected: -1},
	}
	for _, test := range tests {
		fr := math.FloatRange{From: test.from, To: test.to}
		actual := fr.Range()
		if actual != test.expected {
			t.Errorf("Expected Range %f, got %f", test.expected, actual)
		}
	}
}

func TestFloatRange_AbsRange(t *testing.T) {
	tests := []struct{ from, to, expected float64 }{
		{from: 0, to: 1, expected: 1},
		{from: -1, to: -2, expected: 1},
	}
	for _, test := range tests {
		fr := math.FloatRange{From: test.from, To: test.to}
		actual := fr.AbsRange()
		if actual != test.expected {
			t.Errorf("Expected AbsRange %f, got %f", test.expected, actual)
		}
	}
}

func TestFloatRange_Cap(t *testing.T) {
	tests := []struct{ from, to, value, expected float64 }{
		{from: 1, to: 5, value: 0, expected: 1},
		{from: 1, to: 5, value: 1, expected: 1},
		{from: 1, to: 5, value: 3, expected: 3},
		{from: 1, to: 5, value: 5, expected: 5},
		{from: 1, to: 5, value: 6, expected: 5},
		{from: 5, to: 1, value: 0, expected: 1},
		{from: 5, to: 1, value: 1, expected: 1},
		{from: 5, to: 1, value: 3, expected: 3},
		{from: 5, to: 1, value: 5, expected: 5},
		{from: 5, to: 1, value: 6, expected: 5},
		{from: -4, to: -4, value: 6, expected: -4},
		{from: -4, to: -4, value: -4, expected: -4},
		{from: -4, to: -4, value: -6, expected: -4},
	}
	for _, test := range tests {
		fr := math.FloatRange{From: test.from, To: test.to}
		actual := fr.Cap(test.value)
		if actual != test.expected {
			t.Errorf("Expected Cap %f, got %f.", test.expected, actual)
			t.Logf("From: %f, To: %f, Value: %f", test.from, test.to, test.value)
		}
	}
}

func TestFloatRange_Normalise(t *testing.T) {
	tests := []struct{ from, to, value, expected float64 }{
		{from: 2, to: 6, value: 1, expected: -0.25},
		{from: 2, to: 6, value: 2, expected: 0},
		{from: 2, to: 6, value: 4, expected: 0.5},
		{from: 2, to: 6, value: 6, expected: 1},
		{from: 2, to: 6, value: 7, expected: 1.25},
	}
	for _, test := range tests {
		r := math.FloatRange{From: test.from, To: test.to}
		normalised := r.Normalise(test.value)
		if normalised != test.expected {
			t.Errorf("Expected %f, but got %f", test.expected, normalised)
		}
	}
}

func TestFloatRange_ScaleFromNormalised(t *testing.T) {
	tests := []struct{ from, to, value, expected float64 }{
		{from: 3, to: 9, value: -1, expected: -3},
		{from: 3, to: 9, value: 0, expected: 3},
		{from: 3, to: 9, value: 0.5, expected: 6},
		{from: 3, to: 9, value: 1, expected: 9},
		{from: 3, to: 9, value: 1.5, expected: 12},
	}
	for _, test := range tests {
		r := math.FloatRange{From: test.from, To: test.to}
		normalised := r.ScaleFromNormalised(test.value)
		if normalised != test.expected {
			t.Errorf("Expected %f, but got %f", test.expected, normalised)
		}
	}
}
