package main

import (
	"io"
	"log"
	"os"

	"github.com/glynternet/oscli/cmd/oscli/cmd"
	pkgcmd "github.com/glynternet/pkg/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func buildCmdTree(logger *log.Logger, out io.Writer, rootCmd *cobra.Command) {
	rootCmd.AddCommand(pkgcmd.NewBashCompletion(rootCmd, os.Stdout))
	for _, addCmd := range []func(*log.Logger, io.Writer, *cobra.Command) error{
		cmd.Generate,
		cmd.Metro,
		cmd.Monitor,
		cmd.Relay,
		cmd.Send,
	} {
		err := addCmd(logger, out, rootCmd)
		if err != nil {
			log.Fatal(errors.Wrap(err, "adding subcommand"))
		}
	}
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatal(errors.Wrap(err, "binding PFlags"))
	}
}
