package cmd

import (
	"context"
	"fmt"
	"io"

	icmd "github.com/glynternet/oscli/internal/cmd"
	iosc "github.com/glynternet/oscli/internal/osc"
	"github.com/glynternet/oscli/models"
	"github.com/glynternet/oscli/pkg/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyWaveFrequency = "wave-freq"
)

// Generate adds a generate command to the parent command
func Generate(logger log.Logger, w io.Writer, parent *cobra.Command) error {
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

				client, host, err := icmd.ResolveRemoteClient(localMode, remoteHost, int(remotePort))
				if err != nil {
					return errors.Wrap(err, "initialising host")
				}

				if err := icmd.VerifyFlagMessageFrequency(msgFreq); err != nil {
					return errors.Wrap(err, "verifying message frequency")
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

				remoteKV := log.KV{K: "remote", V: fmt.Sprintf("%s:%d", host, remotePort)}
				if err := logger.Log(
					log.Message("Generating and sending messages"),
					remoteKV); err != nil {
					return errors.Wrap(err, "printing log message")
				}
				// TODO: the third argument to this could be a ticker or something?
				msgCh := iosc.Generate(context.TODO(), genFn, wave.Frequency(msgFreq).Period())
				for msg := range msgCh {
					err := client.Send(msg)
					if err != nil {
						_ = logger.Log(
							log.Message("Error sending message to remote"),
							log.Error(err),
							log.KV{K: "oscMessage", V: msg},
							remoteKV)
						continue
					}
					_ = logger.Log(
						log.Message("Message sent to remote"),
						log.Error(err),
						log.KV{K: "oscMessage", V: msg},
						remoteKV,
					)
				}
				return nil
			},
		}
	)

	parent.AddCommand(cmd)
	icmd.FlagLocalMode(cmd, &localMode)
	icmd.FlagRemoteHost(cmd, &remoteHost)
	icmd.FlagRemotePort(cmd, &remotePort)
	icmd.FlagMessageFrequency(cmd, &msgFreq)
	cmd.Flags().Float64VarP(&waveFreq, keyWaveFrequency, "f", 1, "frequency of generated signal")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
