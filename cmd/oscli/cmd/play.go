package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	osc2 "github.com/glynternet/go-osc/osc"
	"github.com/glynternet/oscli/internal/record"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Play adds a play command to the parent command
func Play(logger *log.Logger, _ io.Writer, parent *cobra.Command) error {
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
				client, host, err := initRemoteClient(localMode, remoteHost, int(remotePort))
				if err != nil {
					return errors.Wrap(err, "getting remote host")
				}

				r, err := readFromFile(logger, oscFile)
				if err != nil {
					return errors.Wrapf(err, "reading recording from file:%s", oscFile)
				}

				logger.Printf("Replaying OSC messages")
				addr := fmt.Sprintf("%s:%d", host, remotePort)
				logger.Printf("Sending to  address %s", addr)
				r.Entries.Play(func(_ int, p osc2.Packet) {
					if sErr := client.Send(p); sErr != nil {
						logger.Printf("Error sending message to client: %+v", sErr)
					}
				})

				logger.Println("Finished playing")
				return nil
			},
		}
	)

	parent.AddCommand(cmd)
	flagRemoteHost(cmd, &remoteHost)
	flagRemotePort(cmd, &remotePort)
	flagLocalMode(cmd, &localMode)
	cmd.Flags().StringVar(&oscFile, "osc-file", defaultRecordFile, "recorded osc file")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}

func readFromFile(logger *log.Logger, oscFile string) (record.Recording, error) {
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
		logger.Println(cErr)
	}
	return recording, err
}
