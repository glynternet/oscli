package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const keyMsgFreq = "msg-freq"

func FlagLocalMode(cmd *cobra.Command, localMode *bool) {
	cmd.Flags().BoolVar(localMode, "local", false, "send messages to localhost")
}

func FlagListenHost(cmd *cobra.Command, listenHost *string) {
	cmd.Flags().StringVar(listenHost, "listen-host", "", "host address to listen on")
}

func FlagListenPort(cmd *cobra.Command, listenPort *uint) {
	cmd.Flags().UintVar(listenPort, "listen-port", 9000, "port to listen on")
}

func FlagRemoteHost(cmd *cobra.Command, remoteHost *string) {
	cmd.Flags().StringVarP(remoteHost, "remote-host", "r", "", "address to send any messages to")
}

func FlagRemotePort(cmd *cobra.Command, remotePort *uint) {
	cmd.Flags().UintVar(remotePort, "remote-port", 9000, "remote host post")
}

func FlagAsBlob(cmd *cobra.Command, asBlob *bool) {
	cmd.Flags().BoolVar(asBlob, "as-blob", false, "send all arguments as blobs")
}

func FlagMessageFrequency(cmd *cobra.Command, msgFreq *float64) {
	cmd.Flags().Float64VarP(msgFreq, keyMsgFreq, "m", 25, "frequency to send messages at")
}

func VerifyFlagMessageFrequency(msgFreq float64) error {
	if msgFreq <= 0 {
		return fmt.Errorf("%s must be positive, received %f", keyMsgFreq, msgFreq)
	}
	return nil
}
