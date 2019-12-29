package file

import (
	"io"
	"log"

	"github.com/glynternet/oscli/internal/record"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultCombinedFile = "./combined.osc"

func Combine(logger *log.Logger, w io.Writer, parent *cobra.Command) error {
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
					logger.Printf("Recording read from %s\n", f)
					rs = append(rs, r)
				}

				combined, err := record.Combine(rs...)
				if err != nil {
					return errors.Wrap(err, "combining recordings")
				}
				logger.Print("Recordings combined.")

				if err := writeToFile(logger, combined, output); err != nil {
					return errors.Wrapf(err, "writing recording to file:%s", output)
				}
				logger.Printf("Combined file written to %s", output)
				return nil
			},
		}
	)

	parent.AddCommand(cmd)
	cmd.Flags().StringVar(&output, "output-file", defaultCombinedFile, "file to write combined osc content to")
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
