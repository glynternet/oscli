package osc

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/sander/go-osc/osc"
)

func ReceivePackets(logger *log.Logger, addr string,
	handlePacket PacketHandler,
	handleReceiveError ErrorHandler,
) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return errors.Wrap(err, "creating listener")
	}
	srv := &osc.Server{}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Print(errors.Wrap(err, "closing listen connection"))
		}
	}()
	fmt.Println("Listening on", addr)
	for {
		packet, err := srv.ReceivePacket(conn)
		if err != nil {
			handleReceiveError(err)
			continue
		}

		if packet != nil {
			handlePacket(packet)
		}
	}
}

type PacketHandler func(osc.Packet)

type ErrorHandler func(error)
