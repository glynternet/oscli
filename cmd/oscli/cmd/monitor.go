package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/sander/go-osc/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyNoInput = "no-input"
	keyDecodeBlob = "decode-blobs"
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

		if !viper.GetBool(keyNoInput) {
			fmt.Println(`Press "q" then enter to exit`)
			go startQuitterReader(bufio.NewReader(os.Stdin))
		}

		print := getPrinter(viper.GetBool(keyDecodeBlob))

		fmt.Println("Listening on", addr)
		srv := &osc.Server{}

		for {
			packet, err := srv.ReceivePacket(conn)
			if err != nil {
				fmt.Println("Server error: " + err.Error())
				// TODO: add a flag to exit on error instead of loop?
				continue
			}

			if packet != nil {
				switch packet.(type) {
				default:
					fmt.Println("Unknown packet type!")

				case *osc.Message:
					fmt.Printf("-- OSC Message: ")
					print(packet.(*osc.Message))

				case *osc.Bundle:
					fmt.Println("-- OSC Bundle:")
					bundle := packet.(*osc.Bundle)
					for i, message := range bundle.Messages {
						fmt.Printf("  -- OSC Message #%d: ", i+1)
						print(message)
					}
				}
			}
		}

	},
}

func getPrinter(decodeBlobs bool) func(*osc.Message) {
	if decodeBlobs {
		return decodedBlobsPrint
	}
	return rawPrint
}

func rawPrint(msg *osc.Message) {
	fmt.Println(msg)
}

func decodedBlobsPrint(msg *osc.Message) {
	fmt.Println(msg)
	for i, a := range msg.Arguments {
		if bs, ok := a.([]byte); ok {
			fmt.Printf("element[%d]: %s\n", i, string(bs))
		}
	}
}

type byteReader interface {
	ReadByte() (byte, error)
}

func startQuitterReader(r byteReader) {
	for {
		c, err := r.ReadByte()
		if err != nil {
			fmt.Println(errors.Wrap(err, "reading bytes"))
			os.Exit(1)
		}

		if c == 'q' {
			os.Exit(0)
		}
	}
}

func init() {
	rootCmd.AddCommand(cmdMonitor)
	cmdMonitor.Flags().Bool(keyNoInput, false, "turn on no-input mode for when no terminal is available")
	cmdMonitor.Flags().Bool(keyDecodeBlob, false, "decode blob values into strings")
	err := viper.BindPFlags(cmdMonitor.Flags())
	if err != nil {
		log.Fatal(err)
	}
}
