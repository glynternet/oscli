package cmd

import (
	"net"

	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal"
	"github.com/glynternet/oscli/internal/osc"
	"github.com/pkg/errors"
)

const defaultRecordFile = "./recording.osc"

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

func initRemoteClient(localMode bool, remoteHost string, port int) (*osc2.Client, string, error) {
	host, err := initRemoteHost(localMode, remoteHost)
	if err != nil {
		return nil, "", errors.Wrap(err, "initialising remote host")
	}
	return osc2.NewClient(host, port), host, nil
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
