package cmd

import (
	"fmt"
	"io"
	"log"

	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/pkg/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Send adds a generate command to the parent command
func Send(_ *log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		localMode  bool
		remoteHost string
		remotePort uint
		asBlob     bool

		cmdSend = &cobra.Command{
			Use:   "send ADDRESS [ ARGS... ]",
			Short: "send a single OSC message",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				msgAddr, err := osc.CleanAddress(args[0])
				if err != nil {
					return errors.Wrap(err, "parsing OSC message address")
				}

				host, err := initRemoteHost(localMode, remoteHost)
				if err != nil {
					return errors.Wrap(err, "getting remote host")
				}

				port := int(remotePort)
				client := osc2.NewClient(host, port)
				msg := osc2.NewMessage(msgAddr)
				parse := getParser(asBlob)
				if len(args) > 1 {
					for _, val := range args[1:] {
						app, err := parse(val)
						if err != nil {
							return errors.Wrap(err, "parsing message argument")
						}
						msg.Append(app)
					}
				}

				if err := client.Send(msg); err != nil {
					return errors.Wrapf(err, "sending msg:%v using client:%v", *msg, *client)
				}
				addr := fmt.Sprintf("%s:%d", host, port)
				fmt.Printf("sending to %s: %v\n", addr, msg)
				return nil
			},
		}
	)

	parent.AddCommand(cmdSend)
	flagRemoteHost(cmdSend, &remoteHost)
	flagRemotePort(cmdSend, &remotePort)
	flagLocalMode(cmdSend, &localMode)
	flagAsBlob(cmdSend, &asBlob)
	return errors.Wrap(viper.BindPFlags(cmdSend.Flags()), "binding pflags")
}
