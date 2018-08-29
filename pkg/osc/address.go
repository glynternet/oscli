package osc

import (
	"github.com/pkg/errors"
	"strings"
)

// CleanAddress will ensure that the given address is in a format that is
// suitable for OSC messaging, appending a / before the address if there isn't
// one already.
func CleanAddress(addr string) (string, error) {
	msgAddr := strings.TrimSpace(addr)
	if len(msgAddr) == 0 {
		return "", errors.New("address must be non-zero length")
	}
	if strings.ContainsAny(addr, "\n\t\r ") {
		return "", errors.New("cannot contain whitespace")
	}
	if msgAddr[0] != '/' {
		msgAddr = "/" + msgAddr
	}
	return msgAddr, nil
}
