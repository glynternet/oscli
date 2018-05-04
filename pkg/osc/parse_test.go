package osc_test

import (
	"fmt"
	"testing"

	"github.com/glynternet/oscli/pkg/osc"
	"github.com/magiconair/properties/assert"
)

func TestParse(t *testing.T) {
	for _, test := range []struct {
		name     string
		val      string
		expected interface{}
		error
	}{
		{
			name:  "zero-values",
			error: osc.EmptyValueError,
		},
		{
			name:     "string value",
			val:      "hey man wassup",
			expected: string("hey man wassup"),
		},
		{
			name:     "float32 value with awkward decimal",
			val:      "6.75",
			expected: float32(6.75),
		},
		{
			name:     "float32 value with rounded decimal",
			val:      "13.0",
			expected: float32(13.0),
		},
		{
			name:     "float32 value with decimal but nothing after",
			val:      "100203.",
			expected: float32(100203),
		},
		{
			name:     "int32",
			val:      "56",
			expected: int32(56),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			val, err := osc.Parse(test.val)
			assert.Equal(t, test.error, err)
			assert.Equal(t, test.expected, val, fmt.Sprintf("type expected: %T, type returned: %T", test.expected, val))
		})
	}
}
