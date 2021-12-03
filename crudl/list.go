package crudl

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/fatih/color"
	"github.com/jaymccon/cloudctl/data"
	awsProvider "github.com/jaymccon/cloudctl/providers/aws"
	"github.com/rodaine/table"
	"strings"
)

func ListResource(typeName string) {
	fmt.Printf("Listing resources for %q\n", typeName)
	resources, err := awsProvider.ListResource(typeName)
	if err != nil {
		fmt.Printf("ERROR: %q\n", err.Error())
		return
	}
	headers := append([]interface{}{"Identifier"}, getTableHeaders(*resources)...)
	tbl := table.New(headers...)
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	createRows(tbl, *resources, headers)
	tbl.Print()
}

func getTableHeaders(resources []types.ResourceDescription) []interface{} {
	var headers []interface{}
	for _, r := range resources {
		var props map[string]interface{}
		err := json.Unmarshal([]byte(*r.Properties), &props)
		if err != nil {
			fmt.Printf("failed to unmarshal json properties: %q %q\n", err.Error(), *r.Properties)
		}
		for k, _ := range props {
			lk := strings.ToLower(k)
			// TODO: fields to be included in summary view should be configurable with a declarative file, maybe a schema enhancement
			if lk == "description" || lk == "name" || lk == "subnetid" || lk == "vpcid" {
				var hString []string
				for _, v := range headers {
					hString = append(hString, v.(string))
				}
				if !data.Contains(hString, k) {
					headers = append(headers, k)
				}
			}
		}
	}
	return headers
}

func createRows(tbl table.Table, resources []types.ResourceDescription, headers []interface{}) {
	for _, r := range resources {
		var props map[string]interface{}
		err := json.Unmarshal([]byte(*r.Properties), &props)
		if err != nil {
			fmt.Printf("failed to unmarshal json properties: %q %q\n", err.Error(), *r.Properties)
		}
		row := []interface{}{*r.Identifier}
		for _, header := range headers {
			switch props[header.(string)].(type) {
			case nil:
				if header.(string) != "Identifier" {
					row = append(row, "")
				}
			case string:
				if header.(string) != "Identifier" {
					row = append(row, props[header.(string)].(string))
				}
			default:
				fmt.Printf("unsupported type for output row value %q\n", header)
			}
		}
		tbl.AddRow(row...)
	}
}
