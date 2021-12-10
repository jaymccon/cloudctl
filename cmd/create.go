package cmd

import (
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates cloud resources",
}

func init() {
	RootCmd.AddCommand(CreateCmd)
}
