package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/glynternet/oscli/internal/record"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
				ctx, cancel := contextWithSignalCancels(context.Background(),
					syscall.SIGINT, syscall.SIGTERM)
				defer cancel()
				addr := fmt.Sprintf("%s:%d", listenHost, listenPort)
				logger.Printf("Recording on address %s", addr)
				r, err := record.Record(ctx, logger, addr)
				if err != nil {
					return errors.Wrap(err, "recording packets")
				}
				logger.Println("Finished recording")
				if err := writeToFile(logger, r, output); err != nil {
					return err
				}
				logger.Printf("Written to %s", output)
				return nil
			},
		}
	)
	parent.AddCommand(recordCmd)
	flagListenHost(recordCmd, &listenHost)
	flagListenPort(recordCmd, &listenPort)
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

func writeToFile(logger *log.Logger, r record.Recording, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return errors.Wrapf(err, "creating file at %s", output)
	}
	logger.Printf("Writing to file at %s", output)
	_, err = r.WriteTo(file)
	err = errors.Wrap(err, "writing recording to writer")
	cErr := errors.Wrap(file.Close(), "closing file")
	if err == nil {
		return cErr
	}
	if cErr != nil {
		logger.Println(cErr)
	}
	return err
}
