package cmd

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/sander/go-osc/osc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyForwardHost = "forward-host"
	keyForwardPort = "forward-port"
)

// Relay adds a generate command to the parent command
func Relay(logger *log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		listenHost  string
		listenPort  uint
		forwardHost string
		forwardPort uint

		relayCmd = &cobra.Command{
			Use:   "relay",
			Short: "listen and relay osc messages",
			RunE: func(cmd *cobra.Command, args []string) error {

				if listenHost == forwardHost && listenPort == forwardPort {
					return errors.New("cannot forward to listen address: forward loop")
				}

				listenAddr := fmt.Sprintf("%s:%d", listenHost, listenPort)
				receiveChan, err := receivePackets(logger, listenAddr)
				if err != nil {
					return errors.Wrap(err, "creating packet receiver")
				}

				c := osc.NewClient(forwardHost, int(forwardPort))
				logger.Printf("forwarding to %s:%d", forwardHost, forwardPort)
				for {
					select {
					case p := <-receiveChan:
						if err := c.Send(p); err != nil {
							logger.Print(errors.Wrap(err, "forwarding to client"))
						}
					}
				}
			},
		}
	)
	parent.AddCommand(relayCmd)
	flagListenHost(relayCmd, &listenHost)
	flagListenPort(relayCmd, &listenPort)
	relayCmd.Flags().StringVar(&forwardHost, keyForwardHost, "", "forwarding host address")
	relayCmd.Flags().UintVar(&forwardPort, keyForwardPort, 9000, "forwarding port number")
	return errors.Wrap(viper.BindPFlags(relayCmd.Flags()), "binding pflags")
}

func receivePackets(logger *log.Logger, addr string) (<-chan osc.Packet, error) {
	ch := make(chan osc.Packet)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "creating listener")
	}
	go func() {
		defer func() {
			close(ch)
			logger.Println("closed channel")
			if err := conn.Close(); err != nil {
				logger.Print(errors.Wrap(err, "closing listen connection"))
			}
		}()
		fmt.Println("Listening on", addr)

		for {
			packet, err := (&osc.Server{}).ReceivePacket(conn)
			if err != nil {
				fmt.Println("Receiving packet: " + err.Error())
				return
			}

			if packet != nil {
				switch packet.(type) {
				default:
					fmt.Println("Unknown packet type!")
					continue
				case *osc.Message:
					fmt.Printf("-- OSC Message: ")
					osc.PrintMessage(packet.(*osc.Message))
				case *osc.Bundle:
					fmt.Println("-- OSC Bundle:")
					bundle := packet.(*osc.Bundle)
					for i, message := range bundle.Messages {
						fmt.Printf("  -- OSC Message #%d: ", i+1)
						osc.PrintMessage(message)
					}
				}
				ch <- packet
			}
		}
	}()
	return ch, nil
}
