package cmd

import (
	"github.com/spf13/cobra"
)

var ReadCmd = &cobra.Command{
	Use:   "read",
	Short: "reads cloud resources",
}

func init() {
	RootCmd.AddCommand(ReadCmd)
}
