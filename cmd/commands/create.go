package commands

import (
	"github.com/jaymccon/cloudctl/cmd"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates cloud resources",
	Long:  `TODO`,
}

func init() {
	cmd.RootCmd.AddCommand(CreateCmd)
}
