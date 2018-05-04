package osc

import (
	"strconv"
	"strings"
)

const EmptyValueError = ParseError("empty value")

func Parse(val string) (interface{}, error) {
	if val == "" {
		return nil, EmptyValueError
	}

	if strings.Contains(val, ".") {
		f, err := strconv.ParseFloat(val, 32)
		if err == nil {
			return float32(f), nil
		}
	}

	i, err := strconv.ParseInt(val, 10, 32)
	if err == nil {
		return int32(i), nil
	}

	return val, nil
}

type ParseError string

func (p ParseError) Error() string {
	return string(p)
}
