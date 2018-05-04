package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hypebeast/go-osc/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdMonitor = &cobra.Command{
	Use:   "monitor",
	Short: "monitor incoming OSC messages",
	Run: func(cmd *cobra.Command, args []string) {
		addr := fmt.Sprintf("%s:%d", viper.GetString(keyListenHost), viper.GetInt(keyListenPort))
		conn, err := net.ListenPacket("udp", addr)
		if err != nil {
			fmt.Println(errors.Wrap(err, "creating listener"))
			os.Exit(1)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				log.Print(errors.Wrap(err, "closing listen connection"))
			}
		}()

		fmt.Println(`Press "q" then enter to exit`)

		go func() {
			fmt.Println("Listening on", addr)

			for {
				packet, err := (&osc.Server{}).ReceivePacket(conn)
				if err != nil {
					fmt.Println("Server error: " + err.Error())
					os.Exit(1)
				}

				if packet != nil {
					switch packet.(type) {
					default:
						fmt.Println("Unknown packet type!")

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
				}
			}
		}()

		reader := bufio.NewReader(os.Stdin)

		for {
			c, err := reader.ReadByte()
			if err != nil {
				fmt.Println(errors.Wrap(err, "reading bytes"))
				os.Exit(1)
			}

			if c == 'q' {
				os.Exit(0)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cmdMonitor)
	err := viper.BindPFlags(cmdMonitor.Flags())
	if err != nil {
		log.Fatal(err)
	}
}
