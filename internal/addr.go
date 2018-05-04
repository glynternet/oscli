package internal

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const localhost = "localhost"

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
