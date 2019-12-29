package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	icmd "github.com/glynternet/oscli/internal/cmd"
	"github.com/glynternet/oscli/internal/osc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyNoInput    = "no-input"
	keyDecodeBlob = "decode-blobs"
)

// Monitor adds a generate command to the parent command
func Monitor(logger *log.Logger, _ io.Writer, parent *cobra.Command) error {
	var (
		listenHost  string
		listenPort  uint
		noInput     bool
		decodeBlobs bool

		cmd = &cobra.Command{
			Use:   "monitor",
			Short: "monitor incoming OSC messages",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				if !noInput {
					fmt.Println(`Press "q" then enter to exit`)
					go startQuitterReader(bufio.NewReader(os.Stdin))
				}

				return errors.Wrap(
					osc.ReceivePackets(context.Background(), logger,
						fmt.Sprintf("%s:%d", listenHost, listenPort),
						osc.Print(decodeBlobs),
						printError),
					"receiving packets")
			},
		}
	)

	parent.AddCommand(cmd)
	icmd.FlagListenHost(cmd, &listenHost)
	icmd.FlagListenPort(cmd, &listenPort)
	cmd.Flags().BoolVar(&noInput, keyNoInput, false, "turn on no-input mode for when no terminal is available")
	cmd.Flags().BoolVar(&decodeBlobs, keyDecodeBlob, false, "decode blob values into strings")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}

type byteReader interface {
	ReadByte() (byte, error)
}

func startQuitterReader(r byteReader) {
	for {
		c, err := r.ReadByte()
		if err != nil {
			fmt.Println(errors.Wrap(err, "reading bytes"))
			os.Exit(1)
		}

		if c == 'q' {
			os.Exit(0)
		}
	}
}
