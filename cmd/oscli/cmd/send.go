package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/glynternet/oscli/internal"
	"github.com/glynternet/oscli/pkg/osc"
	"github.com/pkg/errors"
	osc2 "github.com/sander/go-osc/osc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdSend = &cobra.Command{
	Use:   "send",
	Short: "send a single OSC message",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("expects at least 2 arguments, address and message parts. Received %d", len(args))
		}
		msgAddr, err := osc.CleanAddress(args[0])
		if err != nil {
			return errors.Wrap(err, "parsing OSC message address")
		}

		host, err := internal.GetRemoteHost(
			viper.GetBool(keyLocal),
			remoteHost,
		)
		if err != nil {
			return errors.Wrap(err, "getting remote host")
		}

		_, err = net.LookupHost(host)
		if err != nil {
			return errors.Wrapf(err, "looking up host %s on network", host)
		}
		port := int(remotePort)
		client := osc2.NewClient(
			host,
			port,
		)

		msg := osc2.NewMessage(msgAddr)
		parse := getParser(asBlob)
		for _, val := range args[1:] {
			app, err := parse(val)
			if err != nil {
				return errors.Wrap(err, "parsing message argument")
			}
			msg.Append(app)
		}
		err = client.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "sending msg:%v using client:%v", *msg, *client)
		}
		addr := fmt.Sprintf("%s:%d", host, port)
		fmt.Printf("sending to %s: %v\n", addr, msg)
		return nil
	},
}

func init() {

	rootCmd.AddCommand(cmdSend)
	err := viper.BindPFlags(cmdSend.Flags())
	if err != nil {
		log.Fatal(err)
	}
}
