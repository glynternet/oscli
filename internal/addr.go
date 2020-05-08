package internal

import (
	"bufio"
	"fmt"
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

type stringReader interface {
	ReadString(byte) (string, error)
}

func readRemoteAddress(reader stringReader) (string, error) {
	if _, err := fmt.Println("Please enter the remote address."); err != nil {
		return "", errors.Wrap(err, "printing remote address prompt")
	}
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
