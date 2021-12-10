package cmd

import (
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes cloud resources",
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
}
