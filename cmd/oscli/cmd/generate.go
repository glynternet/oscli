package cmd

import (
	"fmt"
	"log"

	iosc "github.com/glynternet/oscli/internal/osc"
	"github.com/glynternet/oscli/models"
	"github.com/glynternet/oscli/pkg/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/pkg/errors"
	hosc "github.com/sander/go-osc/osc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyMsgFrequency  = "msg-freq"
	keyWaveFrequency = "wave-freq"
)

var cmdOSCGen = &cobra.Command{
	Use:   "generate [ADDRESS] [MESSAGE]...",
	Short: "generate a stream of osc messages",
	Long: `generate a stream of osc messages

Generate an osc signal with values ranging from 0 to 1 as a sin wave.
The messages will be sent to the given address.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		msgAddr, err := osc.CleanAddress(args[0])
		if err != nil {
			return errors.Wrap(err, "parsing OSC message address")
		}

		host, err := initRemoteHost()
		if err != nil {
			return errors.Wrap(err, "initialising host")
		}

		client := hosc.NewClient(
			host,
			int(remotePort),
		)

		msgFreq := viper.GetFloat64(keyMsgFrequency)
		if msgFreq <= 0 {
			log.Fatal(fmt.Errorf("%s must be positive, received %f", keyMsgFrequency, msgFreq))
		}

		var staticArgs []interface{}
		if len(args) > 0 {
			for _, arg := range args[1:] {
				a, err := iosc.Parse(arg)
				if err != nil {
					return errors.Wrapf(err, "parsing arg '%s' as value", arg)
				}
				staticArgs = append(staticArgs, a)
			}
		}

		genFn := models.NowSinNormalised(msgAddr, staticArgs, viper.GetFloat64(keyWaveFrequency))

		// TODO: the second argument to this could be a ticker or something?
		msgCh := iosc.Generate(genFn, wave.Frequency(msgFreq).Period())
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
	rootCmd.AddCommand(cmdOSCGen)
	cmdOSCGen.Flags().Float64P(keyMsgFrequency, "m", 25, "frequency to send messages at")
	cmdOSCGen.Flags().Float64P(keyWaveFrequency, "f", 1, "frequency of generated signal")
	err := viper.BindPFlags(cmdOSCGen.Flags())
	if err != nil {
		log.Fatal(err)
	}
}
