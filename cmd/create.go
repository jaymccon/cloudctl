package cmd

import (
	"fmt"
	"github.com/jaymccon/cloudctl/crudl"
	"github.com/jaymccon/cloudctl/data"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates cloud resources",
}

func CreateEdit(typeName string) {
	schema, err := data.GetSchema(typeName)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	yamlFile := data.YamlDoc{}
	for name, props := range schema.Properties {
		yamlFile = append(yamlFile, data.NewProp(name, props, *schema, nil, 0))
	}
	yamlFile.Sort()
	outp := string(yamlFile.Marshal())
	jsonDoc, err := data.Edit(outp, "yml")
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}
	// TODO: propagate noPrompts and async
	crudl.CreateResource(typeName, string(jsonDoc), true, false)
}

func init() {
	RootCmd.AddCommand(CreateCmd)
}
