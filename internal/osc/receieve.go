package osc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sander/go-osc/osc"
)

// ReceivePackets will receive all packets at the given address.
// Successfully received packets will be handled by the PacketHandler.
// Errors whilst receiving will be handled by the given ErrorHandler.
func ReceivePackets(ctx context.Context, logger *log.Logger, addr string,
	handlePacket PacketHandler,
	handleReceiveError ErrorHandler,
) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return errors.Wrap(err, "creating listener")
	}
	srv := &osc.Server{ReadTimeout: time.Second}
	fmt.Println("Listening on", addr)

	for {
		select {
		case <-ctx.Done():
			err = conn.Close()
			if err != nil {
				return errors.Wrapf(err, "closing connection")
			}
			logger.Println("Listen connection closed")
			return nil
		default:
			packet, err := srv.ReceivePacket(conn)
			if err != nil {
				handleError(handleReceiveError, err)
				continue
			}

			if packet != nil {
				handlePacket(packet)
			}
		}
	}
}

func handleError(handler ErrorHandler, err error) {
	if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
		return
	}
	handler(err)
}

// PacketHandler handles a packet
type PacketHandler func(osc.Packet)

// ErrorHandler handles an error
type ErrorHandler func(error)
