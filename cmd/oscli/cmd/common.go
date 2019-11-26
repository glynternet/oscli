package cmd

import (
	"github.com/glynternet/oscli/internal/osc"
)

const defaultRecordFile = "./recording.osc"

func getParser(asBlobs bool) func(string) (interface{}, error) {
	if asBlobs {
		return osc.BlobParse
	}
	return osc.Parse
}
