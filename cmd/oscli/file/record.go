package file

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/glynternet/oscli/internal"
	"github.com/glynternet/oscli/internal/cmd"
	"github.com/glynternet/oscli/internal/record"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultRecordFile = "./recording.osc"

// Record adds a record command to the parent command
func Record(logger *log.Logger, w io.Writer, parent *cobra.Command) error {
	var (
		listenHost string
		listenPort uint
		output     string

		recordCmd = &cobra.Command{
			Use:   "record",
			Short: "record received osc messages",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				if output = strings.TrimSpace(output); output == "" {
					return errors.New("no output provided")
				}
				if _, err := os.Stat(output); !os.IsNotExist(err) {
					return fmt.Errorf("file already exists at %s", output)
				}
				wc, err := fileCreatingWriteCloser(logger, output)
				if err != nil {
					return errors.Wrapf(err, "creating new WriterCloser for output:%s", output)
				}
				ctx, cancel := internal.ContextWithSignalCancels(context.Background(),
					syscall.SIGINT, syscall.SIGTERM)
				defer cancel()
				addr := fmt.Sprintf("%s:%d", listenHost, listenPort)
				logger.Printf("Recording on address %s", addr)
				r, err := record.Record(ctx, logger, addr)
				if err != nil {
					return errors.Wrap(err, "recording packets")
				}
				logger.Println("Finished recording")
				if err := catchFirstLogOthers(logger, writeToWriteCloser(r, wc)); err != nil {
					return errors.Wrap(err, "writing recording")
				}
				logger.Print("Finished writing")
				return nil
			},
		}
	)
	parent.AddCommand(recordCmd)
	cmd.FlagListenHost(recordCmd, &listenHost)
	cmd.FlagListenPort(recordCmd, &listenPort)
	recordCmd.Flags().StringVar(&output, "osc-file", defaultRecordFile, "file to record osc stream to")
	return errors.Wrap(viper.BindPFlags(recordCmd.Flags()), "binding pflags")
}
