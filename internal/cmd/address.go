package cmd

import (
	"net"

	"github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal"
	"github.com/pkg/errors"
)

func ResolveRemoteHost(localMode bool, remoteHost string) (string, error) {
	host, err := internal.GetRemoteHost(
		localMode,
		remoteHost,
	)
	if err != nil {
		return "", errors.Wrap(err, "getting remote host")
	}

	return host, errors.Wrap(VerifyHost(host), "verifying host")
}

func ResolveRemoteClient(localMode bool, remoteHost string, port int) (*osc.Client, string, error) {
	host, err := ResolveRemoteHost(localMode, remoteHost)
	if err != nil {
		return nil, "", errors.Wrap(err, "initialising remote host")
	}
	return osc.NewClient(host, port), host, nil
}

// verifyHost checks that the given string can be resolved through the current
// DNS/networking state
func VerifyHost(host string) error {
	_, err := net.LookupHost(host)
	return errors.Wrapf(err, "looking up host %s on network", host)
}
