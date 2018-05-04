package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hypebeast/go-osc/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyForwardHost = "forward-host"
	keyForwardPort = "forward-port"
)

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "listen and relay osc messages",
	Run: func(cmd *cobra.Command, args []string) {
		listenHost := viper.GetString(keyListenHost)
		listenPort := viper.GetInt(keyListenPort)
		forwardHost := viper.GetString(keyForwardHost)
		forwardPort := viper.GetInt(keyForwardPort)

		if listenHost == forwardHost && listenPort == forwardPort {
			log.Fatal(errors.New("cannot forward to listen address: forward loop"))
		}

		listenAddr := fmt.Sprintf("%s:%d", listenHost, listenPort)
		receiveChan, err := receivePackets(listenAddr)
		if err != nil {
			log.Fatal(errors.Wrap(err, "creating packet receiver"))
		}

		c := osc.NewClient(forwardHost, forwardPort)

		for {
			select {
			case p := <-receiveChan:
				if err := c.Send(p); err != nil {
					log.Print(errors.Wrap(err, "forwarding to client"))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(relayCmd)
	relayCmd.Flags().String(keyForwardHost, "", "forwarding host address")
	relayCmd.Flags().Uint(keyForwardPort, 9000, "forwarding port number")
	err := viper.BindPFlags(relayCmd.Flags())
	if err != nil {
		log.Fatal(err)
	}
}

func receivePackets(addr string) (<-chan osc.Packet, error) {
	ch := make(chan osc.Packet)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "creating listener")
	}
	go func() {
		defer func() {
			err := conn.Close()
			if err != nil {
				log.Print(errors.Wrap(err, "closing listen connection"))
			}
		}()
		fmt.Println("Listening on", addr)

		for {
			packet, err := (&osc.Server{}).ReceivePacket(conn)
			if err != nil {
				fmt.Println("Receiving packet: " + err.Error())
				os.Exit(1)
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
