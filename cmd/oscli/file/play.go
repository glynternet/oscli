package file

import (
	"fmt"
	"io"
	"log"

	"github.com/glynternet/go-osc/osc"
	icmd "github.com/glynternet/oscli/internal/cmd"
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

		cmd = &cobra.Command{
			Use:   "play",
			Short: "play a recorded osc file",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, host, err := icmd.ResolveRemoteClient(localMode, remoteHost, int(remotePort))
				if err != nil {
					return errors.Wrap(err, "getting remote host")
				}

				oscFile := defaultRecordFile
				if len(args) == 1 {
					oscFile = args[0]
				}
				r, err := readFromFile(logger, oscFile)
				if err != nil {
					return errors.Wrapf(err, "reading recording from file:%s", oscFile)
				}
				logger.Printf("Messages read from %s\n", oscFile)

				logger.Printf("Replaying OSC messages")
				addr := fmt.Sprintf("%s:%d", host, remotePort)
				logger.Printf("Sending to  address %s", addr)
				r.Data.Entries.Play(func(_ int, p osc.Packet) {
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
	icmd.FlagRemoteHost(cmd, &remoteHost)
	icmd.FlagRemotePort(cmd, &remotePort)
	icmd.FlagLocalMode(cmd, &localMode)
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
