package cmd

import (
	"github.com/spf13/cobra"
)

var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configures cloud resource providers",
}

func init() {
	RootCmd.AddCommand(ConfigureCmd)
}
