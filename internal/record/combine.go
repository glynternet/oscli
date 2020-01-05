package record

import (
	"sort"
)

// Combine combines a group of Entries, sorting them by their given timestamp in the process
func Combine(ess ...Entries) Entries {
	if len(ess) == 0 {
		return Entries{}
	}
	var combined Entries
	for _, es := range ess {
		combined = append(combined, es...)
	}
	sort.Sort(combined)
	return combined
}
