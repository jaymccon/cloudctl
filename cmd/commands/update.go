package commands

import (
	"github.com/jaymccon/cloudctl/cmd"

	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates cloud resources",
	Long:  `TODO`,
}

func init() {
	cmd.RootCmd.AddCommand(UpdateCmd)
}
