package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/glynternet/oscli/internal"
	"github.com/glynternet/oscli/models"
	osc2 "github.com/glynternet/oscli/pkg/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/Pocketbrain/go-logger"
	"github.com/hypebeast/go-osc/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyMsgFrequency  = "msg-freq"
	keyWaveFrequency = "wave-freq"
)

var cmdOSCGen = &cobra.Command{
	Use:   "gen",
	Short: "generate a stream of osc messages",
	Long: `generate a stream of osc messages

Generate an osc signal with values ranging from 0 to 1 as a sin wave.
The messages will be sent to the given address.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("arguments required to form OSC message")
		}
		msgAddr, err := osc2.CleanAddress(args[0])
		if err != nil {
			return errors.Wrap(err, "parsing OSC message address")
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
				staticArgs = append(staticArgs, arg)
			}
		}

		genFn := models.NowSinNormalised(msgAddr, staticArgs, viper.GetFloat64(keyWaveFrequency))

		// TODO: the second argument to this could be a ticker or something?
		msgCh := osc2.Generate(genFn, wave.Frequency(msgFreq).Period())
		for {
			select {
			case msg := <-msgCh:
				err := client.Send(msg)
				if err != nil {
					log.Print(errors.Wrap(err, "sending message to client"))
				}
				log.Printf("Message (%+v) sent to client at %s:%d", msg, client.IP(), client.Port())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cmdOSCGen)
	cmdOSCGen.Flags().Float64P(keyMsgFrequency, "m", 25, "frequency to send messages at")
	cmdOSCGen.Flags().Float64P(keyWaveFrequency, "f", 1, "frequency of generated signal")
	err := viper.BindPFlags(cmdOSCGen.Flags())
	if err != nil {
		plog.Fatal(err)
	}
}
