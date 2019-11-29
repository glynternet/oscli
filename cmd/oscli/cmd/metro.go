package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	osc3 "github.com/glynternet/oscli/internal/osc"
	osc2 "github.com/glynternet/oscli/pkg/osc"
	"github.com/glynternet/oscli/pkg/wave"
	"github.com/pkg/errors"
	"github.com/sander/go-osc/osc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Metro adds a generate command to the parent command
func Metro(logger *log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		remoteHost string
		remotePort uint
		msgFreq    float64
		asBlob     bool
		localMode  bool

		cmd = &cobra.Command{
			Use:   "metro [ADDRESS] [MESSAGE]...",
			Short: "generate a ticker of the same OSC message",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				msgAddr, err := osc2.CleanAddress(args[0])
				if err != nil {
					return errors.Wrap(err, "parsing OSC message address")
				}

				host, err := initRemoteHost(localMode, remoteHost)
				if err != nil {
					return errors.Wrap(err, "initialising host")
				}

				client := osc.NewClient(
					host,
					int(remotePort),
				)

				if msgFreq <= 0 {
					return fmt.Errorf("%s must be positive, received %f", keyMsgFrequency, msgFreq)
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
				msgCh := osc3.Generate(context.TODO(), genFn, wave.Frequency(msgFreq).Period())
				for {
					select {
					case msg := <-msgCh:
						err := client.Send(msg)
						if err != nil {
							logger.Print(errors.Wrap(err, "sending message to client"))
							continue
						}
						logger.Printf("Message (%+v) sent to client at %s:%d", msg, client.IP(), client.Port())
					}
				}
			},
		}
	)

	parent.AddCommand(cmd)
	flagLocalMode(cmd, &localMode)
	flagRemoteHost(cmd, &remoteHost)
	flagRemotePort(cmd, &remotePort)
	flagMessageFrequency(cmd, &msgFreq)
	flagAsBlob(cmd, &asBlob)
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
