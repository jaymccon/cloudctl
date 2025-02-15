package crudl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	"github.com/jaymccon/cloudctl/data"
	"github.com/rodaine/table"
	"log"
	"os"
	"strings"
)

func confirm(s string) bool {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/n]: ", s)
	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	// Empty input (i.e. "\n")
	if len(res) < 2 {
		return false
	}
	return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
}

func GetTableHeaders(resources []types.ResourceDescription) []interface{} {
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
		row := GetRow(r, headers)
		tbl.AddRow(row...)
	}
}

func GetRow(r types.ResourceDescription, headers []interface{}) []interface{} {
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
	return row
}
