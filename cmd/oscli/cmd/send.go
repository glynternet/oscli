package cmd

import (
	"fmt"
	osc3 "github.com/glynternet/oscli/internal/osc"
	"log"
	"net"

	"github.com/glynternet/oscli/internal"
	"github.com/glynternet/oscli/pkg/osc"
	osc2 "github.com/sander/go-osc/osc"
	"github.com/pkg/errors"
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
		msg := osc2.NewMessage(msgAddr)

		for _, val := range args[1:] {
			app, err := osc3.Parse(val)
			if err != nil {
				return errors.Wrap(err, "parsing message argument")
			}
			msg.Append(app)
		}

		host, err := internal.GetRemoteHost(
			viper.GetBool(keyLocal),
			viper.GetString(keyRemoteHost),
		)
		if err != nil {
			log.Fatal(errors.Wrap(err, "getting remote host"))
		}

		_, err = net.LookupHost(host)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "looking up %s host %s", keyRemoteHost, host))
		}
		client := osc2.NewClient(
			host,
			viper.GetInt(keyRemotePort),
		)
		fmt.Printf("sending: %v\n", msg)
		return client.Send(msg)
	},
}

func init() {
	rootCmd.AddCommand(cmdSend)
	err := viper.BindPFlags(cmdSend.Flags())
	if err != nil {
		log.Fatal(err)
	}
}
