package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: fmt.Sprintf("generate a bash completion script for %s", appName),
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
