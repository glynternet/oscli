package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	hosc "github.com/glynternet/go-osc/osc"
	iosc "github.com/glynternet/oscli/internal/osc"
	"github.com/glynternet/oscli/models"
	"github.com/glynternet/oscli/pkg/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyMsgFrequency  = "msg-freq"
	keyWaveFrequency = "wave-freq"
)

// Generate adds a generate command to the parent command
func Generate(logger *log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		localMode  bool
		remoteHost string
		remotePort uint
		msgFreq    float64
		waveFreq   float64

		cmd = &cobra.Command{
			Use:   "generate ADDRESS [ ARGS... ]",
			Short: "generate a stream of osc messages",
			Long: `generate a stream of osc messages

Generate an osc signal with values ranging from 0 to 1 as a sin wave.
The messages will be sent to the given address.`,
			Args: cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				msgAddr, err := osc.CleanAddress(args[0])
				if err != nil {
					return errors.Wrap(err, "parsing OSC message address")
				}

				host, err := initRemoteHost(localMode, remoteHost)
				if err != nil {
					return errors.Wrap(err, "initialising host")
				}

				client := hosc.NewClient(
					host,
					int(remotePort),
				)

				if msgFreq <= 0 {
					return fmt.Errorf("%s must be positive, received %f", keyMsgFrequency, msgFreq)
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

				genFn := models.NowSinNormalised(msgAddr, staticArgs, waveFreq)

				logger.Printf("Generating and sending to %s:%d", host, remotePort)
				// TODO: the third argument to this could be a ticker or something?
				msgCh := iosc.Generate(context.TODO(), genFn, wave.Frequency(msgFreq).Period())
				for msg := range msgCh {
					err := client.Send(msg)
					if err != nil {
						logger.Print(errors.Wrap(err, "sending message to client"))
						continue
					}
					logger.Printf("Message (%+v) sent to client at %s:%d", msg, client.IP(), client.Port())
				}
				return nil
			},
		}
	)

	parent.AddCommand(cmd)
	flagLocalMode(cmd, &localMode)
	flagRemoteHost(cmd, &remoteHost)
	flagRemotePort(cmd, &remotePort)
	flagMessageFrequency(cmd, &msgFreq)
	cmd.Flags().Float64VarP(&waveFreq, keyWaveFrequency, "f", 1, "frequency of generated signal")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
