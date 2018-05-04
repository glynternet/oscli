package main

import (
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cobra.OnInitialize(initConfig)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

const (
	appName = "osc"

	keyListenHost   = "listen-host"
	usageListenHost = "host address to listen on"

	keyListenPort   = "listen-port"
	usageListenPort = "port to listen on"

	keyRemoteHost   = "remote-host"
	usageRemoteHost = "address to send any messages to"

	keyRemotePort   = "remote-port"
	usageRemotePort = "port of the remote host to send any messages to"

	keyLocal   = "local"
	usageLocal = "send messages to localhost"
)

var rootCmd = &cobra.Command{
	Use: appName,
}

func init() {
	rootCmd.PersistentFlags().BoolP(keyLocal, "l", false, usageLocal)
	rootCmd.PersistentFlags().String(keyListenHost, "", usageListenHost)
	rootCmd.PersistentFlags().Uint(keyListenPort, 9000, usageListenPort)
	rootCmd.PersistentFlags().StringP(keyRemoteHost, "r", "", usageRemoteHost)
	rootCmd.PersistentFlags().Uint(keyRemotePort, 9000, usageRemotePort)
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatal(errors.Wrap(err, "binding PFlags"))
	}
}

// initConfig sets AutomaticEnv in viper to true.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
