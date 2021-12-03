package data

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
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
			fmt.Println(k)
			return true
		}
	}
	return false
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
