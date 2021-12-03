package commands

import (
	"github.com/jaymccon/cloudctl/cmd"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes cloud resources",
	Long:  `TODO`,
}

func init() {
	cmd.RootCmd.AddCommand(DeleteCmd)
}
