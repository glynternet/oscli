package osc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanAddress(t *testing.T) {
	for _, test := range []struct {
		name    string
		address string
		cleaned string
		err     bool
	}{
		{
			name: "zero-values",
			err:  true,
		},
		{
			name:    "simple valid",
			address: "/synth1",
			cleaned: "/synth1",
		},
		{
			name:    "with whitespace 0",
			address: "/synth1 woop",
			err:     true,
		},
		{
			name: "with whitespace 1",
			address: `/synth1
woop`,
			err: true,
		},
		{
			name:    "with no slashes",
			address: "synth2",
			cleaned: "/synth2",
		},
		{
			name:    "with multiple slashes",
			address: "//synth2",
			cleaned: "//synth2",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cleaned, err := CleanAddress(test.address)
			assert.Equal(t, test.cleaned, cleaned)
			if test.err {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
