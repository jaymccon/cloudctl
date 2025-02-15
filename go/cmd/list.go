package cmd

import (
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists cloud resources",
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
