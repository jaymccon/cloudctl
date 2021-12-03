package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string
var namespace string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cloudctl",
	Short: "A better way to control your cloud",
	Long: `cloudctl provides a consistent and intuitive interface to AWS and other cloudy services that are available 
in the AWS CloudFormation Registry.

cloudctl is built on top of the AWS Cloud Control API and aims to keep you in the terminal without needing to 
constantly be flipping between documentation and google searches to know how to manage your cloud infrastructure. It 
follows the verb-noun adjective pattern popularised by modern cli's like kubectl`,
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// global flags
	flags := RootCmd.PersistentFlags()
	flags.StringVarP(
		&cfgFile,
		"config",
		"c",
		"",
		"config file (default \"$HOME/.cloudctl.yaml\")",
	)
	flags.StringVarP(
		&namespace,
		"namespace",
		"n",
		"aws",
		"namespace",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cloudctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cloudctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, err := fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			panic(err)
		}
	}
}
