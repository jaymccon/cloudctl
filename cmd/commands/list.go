package commands

import (
	"github.com/jaymccon/cloudctl/cmd"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists cloud resources",
	Long:  `TODO`,
}

func init() {
	cmd.RootCmd.AddCommand(ListCmd)
}
