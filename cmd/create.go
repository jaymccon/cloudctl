package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates cloud resources",
	Long:  `TODO`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
