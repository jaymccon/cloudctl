package crudl

import (
	"fmt"
	"github.com/fatih/color"
	awsProvider "github.com/jaymccon/cloudctl/providers/aws"
	"github.com/rodaine/table"
)

func ListResource(typeName string) {
	resources, err := awsProvider.ListResource(typeName)
	if err != nil {
		fmt.Printf("ERROR: %q\n", err.Error())
		return
	}
	headers := append([]interface{}{"Identifier"}, GetTableHeaders(*resources)...)
	tbl := table.New(headers...)
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	createRows(tbl, *resources, headers)
	tbl.Print()
}
