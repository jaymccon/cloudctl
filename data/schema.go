package data

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	schemaDir = "~/.cloudctl/schemas/"
)

type CfnSchema struct {
	TypeName              string                 `json:"typeName"`
	Description           string                 `json:"description"`
	SourceUrl             string                 `json:"sourceUrl"`
	DocumentationUrl      string                 `json:"documentationUrl"`
	Definitions           map[string]interface{} `json:"definitions"`
	Properties            map[string]interface{} `json:"properties"`
	Required              []string               `json:"required"`
	ReadOnlyProperties    []string               `json:"readOnlyProperties"`
	WriteOnlyProperties   []string               `json:"writeOnlyProperties"`
	CreateOnlyProperties  []string               `json:"createOnlyProperties"`
	DeprecatedProperties  []string               `json:"deprecatedProperties"`
	PrimaryIdentifier     []string               `json:"primaryIdentifier"`
	AdditionalIdentifiers [][]string             `json:"additionalIdentifiers"`
	Handlers              CfnSchemaHandlers      `json:"handlers"`
	TypeConfiguration     map[string]interface{} `json:"typeConfiguration"`
}

type CfnSchemaHandlersPermissions struct {
	Permissions []string `json:"permissions"`
}

type CfnSchemaHandlers struct {
	Create CfnSchemaHandlersPermissions `json:"create"`
	Read   CfnSchemaHandlersPermissions `json:"read"`
	Update CfnSchemaHandlersPermissions `json:"update"`
	Delete CfnSchemaHandlersPermissions `json:"delete"`
	List   CfnSchemaHandlersPermissions `json:"list"`
}

func (s CfnSchema) ToJsonString() (*string, error) {
	jsonB, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	jsonStr := string(jsonB)
	return &jsonStr, nil
}

func ParseSchemas() (*map[string]CfnSchema, error) {
	schemas := map[string]CfnSchema{}
	schemasPath, err := absPath(schemaDir)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(*schemasPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		jsonFile, err := os.Open(path)
		if err != nil {
			fmt.Printf("ERROR: cannot open %q: %q\n", path, err.Error())
			return nil
		}
		defer jsonFile.Close()
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Printf("ERROR: cannot read %q: %q\n", path, err.Error())
			return nil
		}
		var schema CfnSchema
		err = json.Unmarshal(byteValue, &schema)
		if err != nil {
			fmt.Printf("ERROR: cannot parse %q: %q\n", path, err.Error())
			return nil
		}
		schemas[schema.TypeName] = schema
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &schemas, nil
}

func (s CfnSchema) IsUpdatable() bool {
	var properties []string
	for k, _ := range s.Properties {
		properties = append(properties, k)
	}
	co := []string{}
	ro := []string{}
	if s.CreateOnlyProperties != nil {
		co = s.CreateOnlyProperties
	}
	if s.ReadOnlyProperties != nil {
		ro = s.ReadOnlyProperties
	}
	notUpdatable := append(co, ro...)
	for _, k := range properties {
		if !Contains(notUpdatable, k) {
			return true
		}
	}
	return false
}

func (s CfnSchema) IsConfigurable() bool {
	if s.TypeConfiguration == nil {
		return false
	}
	if s.TypeConfiguration["properties"] == nil {
		return false
	}
	if len(s.TypeConfiguration["properties"].(map[string]interface{})) == 0 {
		return false
	}
	return true
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		v := strings.Replace(v, "/properties/", "", 1)
		if v == str {
			return true
		}
	}

	return false
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
	Default      interface{}
	Description  *string
	Parent       *property
	Children     *YamlDoc
	ItemProperty *property
}

type YamlDoc []property

func (y YamlDoc) Sort() {
	sort.Slice(y, func(i, j int) bool {
		return y[i].Name < y[j].Name
	})
}

func (y YamlDoc) Marshal() []byte {
	req := ""
	optDef := ""
	optNoDef := ""
	for _, prop := range y {
		indent := prop.Depth * 2
		line := ""
		line = strings.Repeat(" ", indent)
		if isOptNoDef(prop.Parent, prop) {
			line = "# " + line
		}
		current := &req
		if prop.ReadOnly {
			continue
		}
		if prop.Name == "" {
			line = line + "- "
		} else {
			line = line + prop.Name + ": "
		}
		if isPrimitive(prop.Type) {
			if prop.Default != nil {
				line = line + prop.Default.(string)
			} else {
				line = line + prop.Type
			}
		}
		if prop.Description != nil {
			line = line + "  # " + *prop.Description
		}
		if !prop.Required && prop.Default == nil {
			current = &optNoDef
		} else if !prop.Required && prop.Default != nil {
			current = &optDef
		}
		if prop.Parent != nil {
			if prop.Parent.Type == "array" {
				if isPrimitive(prop.Type) {
					line = line + "\n"
				}
			} else {
				line = line + "\n"
			}
		} else {
			line = line + "\n"
		}
		*current = *current + line
		if prop.ItemProperty != nil {
			if isPrimitive(prop.ItemProperty.Type) {
				indent := prop.ItemProperty.Depth * 2
				line := strings.Repeat(" ", indent)
				if current == &optNoDef {
					line = "# " + line
				}
				line = line + "- " + prop.ItemProperty.Type
				*current = *current + line + "\n"
			} else {
				yd := string(YamlDoc{*prop.ItemProperty}.Marshal())
				if strings.HasSuffix(*current, "- ") {
					yd = strings.TrimLeft(yd, "# ")
				}
				*current = *current + yd
			}
		} else if prop.Children != nil {
			yd := string(prop.Children.Marshal())
			if strings.HasSuffix(*current, "- ") {
				yd = strings.TrimLeft(yd, "# ")
			}
			*current = *current + yd
		}
	}
	return []byte(req + optDef + optNoDef)
}

func isOptNoDef(parent *property, prop property) bool {
	if !prop.Required && prop.Default == nil {
		return true
	} else if parent == nil {
		return false
	} else if !parent.Required && parent.Default == nil {
		return true
	}
	return false
}

func NewProp(name string, iface interface{}, schema CfnSchema, parent *property, depth int) property {
	if val, ok := iface.(map[string]interface{})["$ref"]; ok {
		def := strings.Replace(val.(string), "#/definitions/", "", 1)
		iface = schema.Definitions[def]
	}
	propMap := iface.(map[string]interface{})
	keys := make([]string, 0, len(propMap))
	for k := range propMap {
		keys = append(keys, k)
	}
	prop := property{
		Name:         name,
		Type:         propMap["type"].(string),
		Depth:        depth,
		Required:     false,
		WriteOnly:    false,
		ReadOnly:     false,
		CreateOnly:   false,
		Interface:    iface,
		Default:      propMap["default"],
		Description:  nil,
		Parent:       parent,
		Children:     nil,
		ItemProperty: nil,
	}
	if Contains(schema.Required, name) {
		prop.Required = true
	}
	if Contains(schema.WriteOnlyProperties, name) {
		prop.WriteOnly = true
	}
	if Contains(schema.CreateOnlyProperties, name) {
		prop.CreateOnly = true
	}
	if Contains(schema.ReadOnlyProperties, name) {
		prop.ReadOnly = true
	}
	if propMap["description"] != nil {
		desc := propMap["description"].(string)
		prop.Description = &desc
	}
	if prop.Type == "object" {
		var children YamlDoc
		for n, i := range iface.(map[string]interface{})["properties"].(map[string]interface{}) {
			switch v := i.(type) {
			case map[string]interface{}:
				children = append(children, NewProp(n, i, schema, &prop, depth+1))
			default:
				fmt.Printf("i is %q", v)
			}
		}
		children.Sort()
		prop.Children = &children
	} else if prop.Type == "array" {
		items := NewProp("", iface.(map[string]interface{})["items"], schema, &prop, depth+1)
		prop.ItemProperty = &items
	}
	return prop
}

func Edit(initialContent string, fileExt string) ([]byte, error) {
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
		fmt.Printf("ERROR: %s\n", err.Error())
		return nil, err
	}
	err = editor.Wait()
	if err != nil {
		fmt.Printf("Error while editing. Error: %s\n", err.Error())
	}
	readFile, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	desiredState := map[string]interface{}{}
	err = yaml.Unmarshal(readFile, desiredState)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	jsonBytes, err := json.Marshal(desiredState)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return jsonBytes, nil
}
