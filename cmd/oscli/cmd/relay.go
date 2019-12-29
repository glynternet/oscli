package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/glynternet/go-osc/osc"
	icmd "github.com/glynternet/oscli/internal/cmd"
	osc2 "github.com/glynternet/oscli/internal/osc"
	"github.com/pkg/errors"
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
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				if listenHost == forwardHost && listenPort == forwardPort {
					return errors.New("cannot forward to listen address: forward loop")
				}

				c := osc.NewClient(forwardHost, int(forwardPort))
				printPacket := osc2.Print(false)
				handle := func(p osc.Packet) {
					printPacket(p)
					if err := c.Send(p); err != nil {
						logger.Print(errors.Wrap(err, "forwarding to client"))
					}
				}

				listenAddr := fmt.Sprintf("%s:%d", listenHost, listenPort)
				logger.Printf("forwarding to %s:%d", forwardHost, forwardPort)
				return errors.Wrap(
					osc2.ReceivePackets(context.Background(), logger, listenAddr, handle, printError),
					"receiving packets")
			},
		}
	)
	parent.AddCommand(relayCmd)
	icmd.FlagListenHost(relayCmd, &listenHost)
	icmd.FlagListenPort(relayCmd, &listenPort)
	relayCmd.Flags().StringVar(&forwardHost, keyForwardHost, "", "forwarding host address")
	relayCmd.Flags().UintVar(&forwardPort, keyForwardPort, 9000, "forwarding port number")
	return errors.Wrap(viper.BindPFlags(relayCmd.Flags()), "binding pflags")
}
func printError(err error) {
	fmt.Println("Receiving packet: " + err.Error())
}
