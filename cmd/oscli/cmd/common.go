package cmd

import (
	"github.com/glynternet/oscli/internal/osc"
)

func getParser(asBlobs bool) func(string) (interface{}, error) {
	if asBlobs {
		return osc.BlobParse
	}
	return osc.Parse
}
