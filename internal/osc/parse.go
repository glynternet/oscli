package osc

import (
	"strconv"
	"strings"
)

// EmptyValueError is returned if an empty string is passed to Parse
const EmptyValueError = ParseError("empty value")

// Parse attempts to parse a given string and return a value that is a
// supported osc argument that represents the string value.
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

// BlobParse will convert any string argument to a []byte so that it will be sent
// as a blob argument
func BlobParse(val string) (interface{}, error) {
	if val == "" {
		return nil, EmptyValueError
	}
	return []byte(val), nil
}

// ParseError is an error type that can be returned when there is an issue with
// parsing a string value in order to turn it into an osc.Message.
type ParseError string

// Error returns a description of the error that occurred whilst parsing a
// string value.
func (p ParseError) Error() string {
	return string(p)
}
