package record

import (
	"bytes"
	"encoding/json"
)

func jsonUnmarshalStrictFields(data []byte, value interface{}) error {
	d := json.NewDecoder(bytes.NewBuffer(data))
	d.DisallowUnknownFields()
	return d.Decode(&value)
}
