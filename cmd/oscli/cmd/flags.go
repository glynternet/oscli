package cmd

import "github.com/spf13/cobra"

const (
	keyListenHost   = "listen-host"
	usageListenHost = "host address to listen on"

	keyListenPort   = "listen-port"
	usageListenPort = "port to listen on"
)

func flagLocalMode(cmd *cobra.Command, localMode *bool) {
	cmd.Flags().BoolVar(localMode, "local", false, "send messages to localhost")
}

func flagListenHost(cmd *cobra.Command, listenHost *string) {
	cmd.Flags().StringVar(listenHost, keyListenHost, "", usageListenHost)
}

func flagListenPort(cmd *cobra.Command, listenPort *uint) {
	cmd.Flags().UintVar(listenPort, keyListenPort, 9000, usageListenPort)
}

func flagRemoteHost(cmd *cobra.Command, remoteHost *string) {
	cmd.Flags().StringVarP(remoteHost, "remote-host", "r", "", "address to send any messages to")
}

func flagRemotePort(cmd *cobra.Command, remotePort *uint) {
	cmd.Flags().UintVar(remotePort, "remote-port", 9000, "remote host post")
}

func flagAsBlob(cmd *cobra.Command, asBlob *bool) {
	cmd.Flags().BoolVar(asBlob, "as-blob", false, "send all arguments as blobs")
}

func flagMessageFrequency(cmd *cobra.Command, msgFreq *float64) {
	cmd.Flags().Float64VarP(msgFreq, keyMsgFrequency, "m", 25, "frequency to send messages at")
}
