package internal

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const localhost = "localhost"

// GetRemoteHost will return remote host address based on whether localMode
// is set to true or whether it has been set in the environment. If neither
// are set, GetRemoteHost will request a remote address from the terminal
// user.
// GetRemoteHost will return an error if both localMode and a host in the
// environment are set.
func GetRemoteHost(localMode bool, envHost string) (string, error) {
	if localMode {
		if envHost != "" {
			return "", errors.New("remote envHost cannot be set in local mode")
		}
		return localhost, nil
	}
	if envHost != "" {
		return envHost, nil
	}
	var err error
	envHost, err = readRemoteAddress(bufio.NewReader(os.Stdin))
	return envHost, errors.Wrap(err, "reading remote address")
}

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

type stringReader interface {
	ReadString(byte) (string, error)
}

func readRemoteAddress(reader stringReader) (string, error) {
	log.Println("Please enter the remote address.")
	addr, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "reading address from reader string")
	}
	addr = strings.TrimSpace(addr)
	if addr == "" {
		addr = localhost
	}
	return addr, err
}
