package file

import (
	"io"

	"github.com/glynternet/oscli/internal/record"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultCombinedFile = "./combined.osc"

func Combine(logger log.Logger, w io.Writer, parent *cobra.Command) error {
	var (
		output string

		cmd = &cobra.Command{
			Use:   "combine OSC_FILES...",
			Short: "combine multiple a osc files into a single file",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, oscFiles []string) error {
				var rs []record.Recording
				for _, f := range oscFiles {
					r, err := readFromFile(logger, f)
					if err != nil {
						return errors.Wrapf(err, "reading recording from file:%s", f)
					}

					if err := logger.Log(log.Message("Recording read from file"),
						log.KV{K: "path", V: f}); err != nil {
						return errors.Wrap(err, "writing log message")
					}

					rs = append(rs, r)
				}

				combined, err := record.Combine(rs...)
				if err != nil {
					return errors.Wrap(err, "combining recordings")
				}

				if err := logger.Log(log.Message("Recordings combined.")); err != nil {
					return errors.Wrap(err, "writing log message")
				}

				if err := writeToFile(logger, combined, output); err != nil {
					return errors.Wrapf(err, "writing recording to file:%s", output)
				}

				return errors.Wrap(logger.Log(log.Message("Combined file written to file"),
					log.KV{K: "path", V: output}), "writing log message")
			},
		}
	)

	parent.AddCommand(cmd)
	cmd.Flags().StringVar(&output, "output-file", defaultCombinedFile, "file to write combined osc content to")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
