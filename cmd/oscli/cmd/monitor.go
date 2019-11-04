package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/sander/go-osc/osc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyNoInput    = "no-input"
	keyDecodeBlob = "decode-blobs"
)

func Monitor(logger *log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		listenHost  string
		listenPort  uint
		noInput     bool
		decodeBlobs bool

		cmd = &cobra.Command{
			Use:   "monitor",
			Short: "monitor incoming OSC messages",
			RunE: func(cmd *cobra.Command, args []string) error {
				addr := fmt.Sprintf("%s:%d", listenHost, listenPort)
				conn, err := net.ListenPacket("udp", addr)
				if err != nil {
					return errors.Wrap(err, "creating listener")
				}
				defer func() {
					err := conn.Close()
					if err != nil {
						logger.Print(errors.Wrap(err, "closing listen connection"))
					}
				}()

				if !noInput {
					fmt.Println(`Press "q" then enter to exit`)
					go startQuitterReader(bufio.NewReader(os.Stdin))
				}

				printMsg := getPrinter(decodeBlobs)

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
							printMsg(packet.(*osc.Message))

						case *osc.Bundle:
							fmt.Println("-- OSC Bundle:")
							bundle := packet.(*osc.Bundle)
							for i, message := range bundle.Messages {
								fmt.Printf("  -- OSC Message #%d: ", i+1)
								printMsg(message)
							}
						}
					}
				}

			},
		}
	)

	parent.AddCommand(cmd)
	flagListenHost(cmd, &listenHost)
	flagListenPort(cmd, &listenPort)
	cmd.Flags().BoolVar(&noInput, keyNoInput, false, "turn on no-input mode for when no terminal is available")
	cmd.Flags().BoolVar(&decodeBlobs, keyDecodeBlob, false, "decode blob values into strings")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
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
