package cmd

import (
	"fmt"
	"log"

	"github.com/glynternet/oscli/internal"
	osc2 "github.com/glynternet/oscli/pkg/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/hypebeast/go-osc/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdMetro = &cobra.Command{
	Use:   "metro [ADDRESS] [MESSAGE]...",
	Short: "generate a ticker of the same OSC message",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		msgAddr, err := internal.CleanAddress(args[0])
		if err != nil {
			return errors.Wrap(err, "parsing OSC message address")
		}

		host, err := initRemoteHost()
		if err != nil {
			return errors.Wrap(err, "initialising host")
		}

		client := osc.NewClient(
			host,
			viper.GetInt(keyRemotePort),
		)

		msgFreq := viper.GetFloat64(keyMsgFrequency)
		if msgFreq <= 0 {
			log.Fatal(fmt.Errorf("%s must be positive, received %f", keyMsgFrequency, msgFreq))
		}

		var staticArgs []interface{}
		if len(args) > 0 {
			for _, arg := range args[1:] {
				a, err := osc2.Parse(arg)
				if err != nil {
					return errors.Wrapf(err, "parsing arg '%s' as value", arg)
				}
				staticArgs = append(staticArgs, a)
			}
		}

		genFn := func() *osc.Message {
			return osc.NewMessage(msgAddr, staticArgs...)
		}

		// TODO: the second argument to this could be a ticker or something?
		msgCh := osc2.Generate(genFn, wave.Frequency(msgFreq).Period())
		for {
			select {
			case msg := <-msgCh:
				err := client.Send(msg)
				if err != nil {
					log.Print(errors.Wrap(err, "sending message to client"))
					continue
				}
				log.Printf("Message (%+v) sent to client at %s:%d", msg, client.IP(), client.Port())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cmdMetro)
	err := viper.BindPFlags(cmdMetro.Flags())
	if err != nil {
		log.Fatal(err)
	}
}
