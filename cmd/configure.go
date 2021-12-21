package cmd

import (
	"fmt"
	"github.com/jaymccon/cloudctl/data"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configures cloud resource providers",
}

func init() {
	RootCmd.AddCommand(ConfigureCmd)
}

func Edit(initialContent string, fileExt string) {
	// TODO: work out filename, populate initial file, parse edited file and check if it changed
	tmpFile, err := ioutil.TempFile(os.TempDir(), "cloudctl-*."+fileExt)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	_, err = tmpFile.Write([]byte(initialContent))
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	err = tmpFile.Close()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	editor := exec.Command("vim", tmpFile.Name())
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	err = editor.Start()
	if err != nil {
		log.Printf("2")
		log.Fatal(err)
	}
	err = editor.Wait()
	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
	} else {
		log.Printf("Successfully edited.")
	}
}

func PocEdit() {
	schema, err := data.GetSchema("AWS::EC2::NetworkInterface")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	yamlFile := yamlDoc{}
	for name, props := range schema.Properties {
		fmt.Println(props)
		yamlFile = append(yamlFile, newProp(name, props, *schema, nil, 0))
	}
	fmt.Println(yamlFile)
}

func isPrimitive(t string) bool {
	if t == "array" || t == "object" {
		return false
	}
	return true
}

type property struct {
	Name         string
	Type         string
	Depth        int
	Required     bool
	WriteOnly    bool
	ReadOnly     bool
	CreateOnly   bool
	Interface    interface{}
	Parent       *property
	Children     *yamlDoc
	ItemProperty *property
}

type yamlDoc []property

func newProp(name string, iface interface{}, schema data.CfnSchema, parent *property, depth int) property {
	if val, ok := iface.(map[string]interface{})["$ref"]; ok {
		def := strings.Replace(val.(string), "#/definitions/", "", 1)
		iface = schema.Definitions[def]
	}
	if iface == nil {
		fmt.Println("iface is nil")
	}
	propMap := iface.(map[string]interface{})
	keys := make([]string, 0, len(propMap))
	for k, _ := range propMap {
		keys = append(keys, k)
	}
	prop := property{
		Name:         name,
		Type:         propMap["type"].(string),
		Depth:        0,
		Required:     false,
		WriteOnly:    false,
		ReadOnly:     false,
		CreateOnly:   false,
		Interface:    iface,
		Parent:       parent,
		Children:     nil,
		ItemProperty: nil,
	}
	if parent != nil {
		prop.Depth = parent.Depth
	}
	if data.Contains(schema.Required, name) {
		prop.Required = true
	}
	if data.Contains(schema.WriteOnlyProperties, name) {
		prop.WriteOnly = true
	}
	if data.Contains(schema.CreateOnlyProperties, name) {
		prop.CreateOnly = true
	}
	if data.Contains(schema.ReadOnlyProperties, name) {
		prop.ReadOnly = true
	}
	if prop.Type == "object" {
		var children yamlDoc
		for n, i := range iface.(map[string]interface{})["properties"].(map[string]interface{}) {
			switch v := i.(type) {
			case map[string]interface{}:
				children = append(children, newProp(n, i, schema, &prop, depth+1))
			default:
				fmt.Printf("i is %q", v)
			}
		}
		prop.Children = &children
	} else if prop.Type == "array" {
		items := newProp("", iface.(map[string]interface{})["items"], schema, &prop, depth+1)
		prop.ItemProperty = &items
	}
	return prop
}
