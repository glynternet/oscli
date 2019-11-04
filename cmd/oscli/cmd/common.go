package cmd

import (
	"net"

	"github.com/glynternet/oscli/internal"
	"github.com/glynternet/oscli/internal/osc"
	"github.com/pkg/errors"
)

func initRemoteHost(localMode bool, remoteHost string) (string, error) {
	host, err := internal.GetRemoteHost(
		localMode,
		remoteHost,
	)
	if err != nil {
		return "", errors.Wrap(err, "getting remote host")
	}

	return host, errors.Wrap(verifyHost(host), "verifying host")
}

// verifyHost checks that the given string can be resolved through the current
// DNS/networking state
func verifyHost(host string) error {
	_, err := net.LookupHost(host)
	return errors.Wrapf(err, "looking up host %s on network", host)
}

func getParser(asBlobs bool) func(string) (interface{}, error) {
	if asBlobs {
		return osc.BlobParse
	}
	return osc.Parse
}
