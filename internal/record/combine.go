package record

import (
	"fmt"
	"sort"
)

// Combine combines a group of Recordings, sorting all of the entries by their given timestamp in the process
func Combine(rs ...Recording) (Recording, error) {
	if len(rs) == 0 {
		return Recording{Schema: version}, nil
	}
	for _, r := range rs {
		if r.Schema != version {
			return Recording{}, UnsupportedSchemaError(r.Schema)
		}
	}
	var combined Entries
	for _, r := range rs {
		combined = append(combined, r.Entries...)
	}
	sort.Sort(combined)
	return Recording{
		Schema:  version,
		Entries: combined,
	}, nil
}

// UnsupportedSchemaError is the error type returned when a Recording wit an unsupported schema is encountered
type UnsupportedSchemaError string

func (us UnsupportedSchemaError) Error() string {
	return fmt.Sprintf("invalid schema: %s", string(us))
}
