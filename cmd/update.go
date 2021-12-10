package cmd

import (
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates cloud resources",
}

func init() {
	RootCmd.AddCommand(UpdateCmd)
}
