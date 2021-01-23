package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/glynternet/go-osc/osc"
	icmd "github.com/glynternet/oscli/internal/cmd"
	iosc "github.com/glynternet/oscli/internal/osc"
	osc2 "github.com/glynternet/oscli/pkg/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Metro adds a generate command to the parent command
func Metro(logger log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		remoteHost string
		remotePort uint
		msgFreq    float64
		asBlob     bool
		localMode  bool

		cmd = &cobra.Command{
			Use:   "metro ADDRESS [ ARGS... ]",
			Short: "generate a ticker of the same OSC message",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				msgAddr, err := osc2.CleanAddress(args[0])
				if err != nil {
					return errors.Wrap(err, "parsing OSC message address")
				}

				client, _, err := icmd.ResolveRemoteClient(localMode, remoteHost, int(remotePort))
				if err != nil {
					return errors.Wrap(err, "initialising host")
				}

				if err := icmd.VerifyFlagMessageFrequency(msgFreq); err != nil {
					return errors.Wrap(err, "verifying message frequency")
				}

				parse := getParser(asBlob)
				var staticArgs []interface{}
				if len(args) > 0 {
					for _, arg := range args[1:] {
						a, err := parse(arg)
						if err != nil {
							return errors.Wrapf(err, "parsing arg:%q", arg)
						}
						staticArgs = append(staticArgs, a)
					}
				}

				genFn := func() *osc.Message {
					return osc.NewMessage(msgAddr, staticArgs...)
				}

				// TODO: the third argument to this could be a ticker or something?
				msgCh := iosc.Generate(context.TODO(), genFn, wave.Frequency(msgFreq).Period())
				for msg := range msgCh {
					err := client.Send(msg)
					if err != nil {
						_ = logger.Log(
							log.Message("Error sending message to client"),
							log.ErrorMessage(err),
							log.KV{K: "oscMessage", V: msg},
							log.KV{K: "clientAddress", V: fmt.Sprintf("%s:%d", client.IP(), client.Port())})
						continue
					}
					_ = logger.Log(
						log.Message("Message sent to client"),
						log.KV{K: "oscMessage", V: msg},
						log.KV{K: "clientAddress", V: fmt.Sprintf("%s:%d", client.IP(), client.Port())})
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
	icmd.FlagAsBlob(cmd, &asBlob)
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
