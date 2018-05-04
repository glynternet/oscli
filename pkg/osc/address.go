package osc

import (
	"errors"
)

func CleanAddress(addr string) (string, error) {
	msgAddr := addr
	if len(msgAddr) == 0 {
		return "", errors.New("address must be non-zero length")
	}
	if msgAddr[0] != '/' {
		msgAddr = "/" + msgAddr
	}
	return msgAddr, nil
}
