package math

// FloatRangeMapper is used to proportionally map values from one FloatRange onto another FloatRange.
type FloatRangeMapper struct {
	Source, Target FloatRange
}

// Map calculates the result for if a given value were to be proportionally from the Source to the Target FloatRange of the mapper.
func (fmr FloatRangeMapper) Map(value float64) float64 {
	return mapUsingRanges(fmr.Source, fmr.Target, value)
}

func mapUsingRanges(source, target FloatRange, value float64) float64 {
	tr := target.Range()
	if tr == 0 {
		return target.From
	}
	sr := source.Range()
	if sr == 0 && value == source.From {
		return (target.From + target.To) / 2
	}
	return (value-source.From)*tr/source.Range() + target.From
}
