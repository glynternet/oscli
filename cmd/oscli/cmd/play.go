package cmd

import (
	"fmt"
	"io"
	"os"

	osc2 "github.com/glynternet/go-osc/osc"
	icmd "github.com/glynternet/oscli/internal/cmd"
	"github.com/glynternet/oscli/internal/record"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Play adds a play command to the parent command
func Play(logger log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		localMode  bool
		remoteHost string
		remotePort uint
		oscFile    string

		cmd = &cobra.Command{
			Use:   "play",
			Short: "play a recorded osc file",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				client, host, err := icmd.ResolveRemoteClient(localMode, remoteHost, int(remotePort))
				if err != nil {
					return errors.Wrap(err, "getting remote host")
				}

				r, err := readFromFile(logger, oscFile)
				if err != nil {
					return errors.Wrapf(err, "reading recording from file:%s", oscFile)
				}

				if err := logger.Log(log.Message("Replaying OSC messages")); err != nil {
					return errors.Wrap(err, "writing log message")
				}

				addr := fmt.Sprintf("%s:%d", host, remotePort)
				if err := logger.Log(log.Message("Sending OSC messages"),
					log.KV{K: "address", V: addr}); err != nil {
					return errors.Wrap(err, "writing log message")
				}
				r.Entries.Play(func(_ int, p osc2.Packet) {
					if sErr := client.Send(p); sErr != nil {
						_ = logger.Log(
							log.Message("Error sending message to client"),
							log.Error(err))
					}
				})

				return errors.Wrap(logger.Log(log.Message("Finished playing")), "writing log message")
			},
		}
	)

	parent.AddCommand(cmd)
	icmd.FlagRemoteHost(cmd, &remoteHost)
	icmd.FlagRemotePort(cmd, &remotePort)
	icmd.FlagLocalMode(cmd, &localMode)
	cmd.Flags().StringVar(&oscFile, "osc-file", defaultRecordFile, "recorded osc file")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}

func readFromFile(logger log.Logger, oscFile string) (record.Recording, error) {
	f, err := os.OpenFile(oscFile, os.O_RDONLY, 0400)
	if err != nil {
		return record.Recording{}, errors.Wrap(err, "opening file")
	}

	var recording record.Recording
	_, err = recording.ReadFrom(f)
	err = errors.Wrap(err, "reading recording from file")
	cErr := errors.Wrap(f.Close(), "closing file")
	if err == nil {
		return recording, cErr
	}
	if cErr != nil {
		_ = logger.Log(log.Message("Error closing file"),
			log.Error(cErr))
	}
	return recording, err
}
