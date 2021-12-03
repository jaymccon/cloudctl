package commands

import (
	"github.com/jaymccon/cloudctl/cmd"

	"github.com/spf13/cobra"
)

var ReadCmd = &cobra.Command{
	Use:   "read",
	Short: "reads cloud resources",
	Long:  `TODO`,
}

func init() {
	cmd.RootCmd.AddCommand(ReadCmd)
}
