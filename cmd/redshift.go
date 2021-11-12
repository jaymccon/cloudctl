//TODO: this file should be auto-generated
package cmd

import (
	"github.com/spf13/cobra"
)

// redshiftCmd represents the redshift command
var redshiftCmd = &cobra.Command{
	Use:   "redshift",
	Short: "Analyze all of your data with the fastest and most widely used cloud data warehouse",
	Long: `Amazon Redshift uses SQL to analyze structured and semi-structured data across data warehouses, operational 
databases, and data lakes, using AWS-designed hardware and ML to deliver the best price performance at any scale.`,
}

func init() {
	createCmd.AddCommand(redshiftCmd)
}
