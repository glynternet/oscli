package cmd

import (
	"io"
	"log"

	"github.com/glynternet/oscli/cmd/oscli/file"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func File(logger *log.Logger, w io.Writer, parent *cobra.Command) error {
	var cmd = &cobra.Command{
		Use:   "file",
		Short: "record, play and combine files",
	}

	for _, subcommand := range []func(*log.Logger, io.Writer, *cobra.Command) error{
		file.Combine,
		file.Play,
		file.Record,
	} {
		if err := subcommand(logger, w, cmd); err != nil {
			return errors.Wrap(err, "adding subcommand to file command")
		}
	}
	parent.AddCommand(cmd)
	return errors.Wrap(viper.BindPFlags(cmd.Flags()), "binding pflags")
}
