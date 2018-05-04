package math

import "math"

// FloatRange represents a range of values between two limits, From and To.
type FloatRange struct {
	From, To float64
}

// Range returns the difference between the limits of a FloatRange.
// If the FloatRange as a From that is less than the To, Range will return a negative value.
func (fr FloatRange) Range() float64 {
	return fr.To - fr.From
}

// AbsRange returns the absolute difference between the limits of a FloatRange.
func (fr FloatRange) AbsRange() float64 {
	return math.Abs(fr.Range())
}

// Cap returns the given value but limited to within the range of the FloatRange limits.
func (fr FloatRange) Cap(value float64) float64 {
	if max := math.Max(fr.From, fr.To); value >= max {
		return max
	}
	if min := math.Min(fr.From, fr.To); value <= min {
		return min
	}
	return value
}

// Normalise returns the result representing where the given value would sit if the FloatRange were to be represented by a range from 0 to 1.
// For example, with a FloatRange From 5 to 9, Normalise of 7 would return a result of 0.5.
func (fr FloatRange) Normalise(value float64) float64 {
	return mapUsingRanges(fr, FloatRange{From: 0, To: 1}, value)
}

// ScaleFromNormalised returns the number that would be represented by the given value if the FloatRange were to be represented by a range of 0 to 1.
// For example, with a FloatRange from -1 to -3, ScaleFromNormalised of 0.25 would return -1.5.
func (fr FloatRange) ScaleFromNormalised(value float64) float64 {
	return mapUsingRanges(FloatRange{From: 0, To: 1}, fr, value)
}
