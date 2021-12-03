package commands

import (
	"github.com/jaymccon/cloudctl/cmd"
	"github.com/jaymccon/cloudctl/data"

	"github.com/spf13/cobra"
)

var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade cloudctl resource definitions",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		err := data.UpdateCache()
		if err != nil {
			cmd.PrintErrf("ERROR: %q\n", err.Error())
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(UpgradeCmd)
}
