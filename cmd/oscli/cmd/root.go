package cmd

import (
	"log"
	"net"
	"strings"

	"github.com/glynternet/oscli/internal"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute() error {
	cobra.OnInitialize(initConfig)
	return rootCmd.Execute()
}

const (
	appName = "oscli"

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

	rootCmd.PersistentFlags().Float64P(keyMsgFrequency, "m", 25, "frequency to send messages at")

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

func initRemoteHost() (string, error) {
	host, err := internal.GetRemoteHost(
		viper.GetBool(keyLocal),
		viper.GetString(keyRemoteHost),
	)
	if err != nil {
		return "", errors.Wrap(err, "getting remote host")
	}

	return host, errors.Wrap(verifyHost(host), "verifying host")
}

// verifyHost checks that the given string can be resolved through the current
// DNS/networking state
func verifyHost(host string) error {
	_, err := net.LookupHost(host)
	return errors.Wrapf(err, "looking up %s host %s", keyRemoteHost, host)
}
