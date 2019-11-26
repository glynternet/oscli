package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/glynternet/oscli/internal/cmd"
	"github.com/glynternet/oscli/internal/record"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Record adds a record command to the parent command
func Record(logger log.Logger, w io.Writer, parent *cobra.Command) error {
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
				ctx, cancel := contextWithSignalCancels(context.Background(),
					syscall.SIGINT, syscall.SIGTERM)
				defer cancel()
				addr := fmt.Sprintf("%s:%d", listenHost, listenPort)
				if err := logger.Log(log.Message("Recording"),
					log.KV{K: "address", V: addr}); err != nil {
					return errors.Wrap(err, "writing log message")
				}

				r, err := record.Record(ctx, logger, addr)
				if err != nil {
					return errors.Wrap(err, "recording packets")
				}
				if err := logger.Log(log.Message("Finished recording"),
					log.KV{K: "address", V: addr}); err != nil {
					return errors.Wrap(err, "writing log message")
				}
				if err := writeToFile(logger, r, output); err != nil {
					return err
				}
				if err := logger.Log(log.Message("Recording file written"),
					log.KV{K: "output", V: output}); err != nil {
					return errors.Wrap(err, "writing log message")
				}
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

func contextWithSignalCancels(ctx context.Context, ss ...os.Signal) (context.Context, context.CancelFunc) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, ss...)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-sigs
		// does it matter that cancel will be called twice?
		cancel()
	}()
	return ctx, cancel
}

func writeToFile(logger log.Logger, r record.Recording, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return errors.Wrapf(err, "creating file at %s", output)
	}
	if err := logger.Log(log.Message("Writing to file"),
		log.KV{K: "output", V: output}); err != nil {
		return errors.Wrap(err, "writing log message")
	}
	_, err = r.WriteTo(file)
	err = errors.Wrap(err, "writing recording to writer")
	cErr := errors.Wrap(file.Close(), "closing file")
	if err == nil {
		return cErr
	}
	if cErr != nil {
		_ = logger.Log(log.Message("Error closing file"),
			log.Error(cErr))
	}
	return err
}
