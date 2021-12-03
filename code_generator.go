package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	basePath            = "schemas/"
	serviceTemplate     = "templates/service.go.tmpl"
	serviceInitTemplate = "templates/service_init.go.tmpl"
)

type CfnType struct {
	Full      string
	Path      string
	Provider  string
	Service   string
	Resource  string
	Schema    map[string]interface{}
	Updatable bool
	Command   string
}

func NewCfnType(name string, path string, schema map[string]interface{}, updatable bool) CfnType {
	parts := strings.Split(name, "::")
	return CfnType{
		Full:      name,
		Path:      path,
		Provider:  parts[0],
		Service:   parts[1],
		Resource:  parts[2],
		Schema:    schema,
		Updatable: updatable,
		Command:   strings.ToLower(parts[1]),
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
func sliceOfInterfacesToStrings(interf []interface{}) []string {
	var slStr []string
	for _, s := range interf {
		str := s.(string)
		str = strings.Replace(str, "/properties/", "", 1)
		slStr = append(slStr, str)
	}
	return slStr
}

func isUpdatable(schema map[string]interface{}) bool {
	var properties []string
	for k, _ := range schema["properties"].(map[string]interface{}) {
		properties = append(properties, k)
	}
	if schema["createOnlyProperties"] == nil {
		schema["createOnlyProperties"] = []interface{}{}
	}
	if schema["readOnlyProperties"] == nil {
		schema["readOnlyProperties"] = []interface{}{}
	}
	notUpdatable := append(
		sliceOfInterfacesToStrings(schema["createOnlyProperties"].([]interface{})),
		sliceOfInterfacesToStrings(schema["readOnlyProperties"].([]interface{}))...,
	)
	for _, k := range properties {
		if !contains(notUpdatable, k) {
			fmt.Println(k)
			return true
		}
	}
	return false
}

func readJson(filePath string) (map[string]interface{}, error) {
	var jsonInterface map[string]interface{}
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byteValue, &jsonInterface)
	if err != nil {
		return nil, err
	}
	return jsonInterface, nil
}

func writeJson(filePath string, jsonInterface map[string]interface{}) error {
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	jsonString, err := json.Marshal(jsonInterface)
	if err != nil {
		return err
	}
	_, err = f.Write(jsonString)
	if err != nil {
		return err
	}
	return nil
}

func generate() {
	resources := map[string]CfnType{}
	var services []string
	var updatable []string

	// gather schemas
	schemaFiles := getSchemaFiles(basePath)
	for _, schemaPath := range schemaFiles {
		schema, err := readJson(schemaPath)
		if err != nil {
			fmt.Printf("ERROR: failed to open %q %q\n", schemaPath, err.Error())
			continue
		}
		svc := NewCfnType(schema["typeName"].(string), schemaPath, schema, isUpdatable(schema))
		if !contains(services, svc.Command) {
			services = append(services, svc.Command)
		}
		if svc.Updatable && !contains(updatable, svc.Command) {
			updatable = append(updatable, svc.Command)
		}
		resources[svc.Full] = svc
	}

	// render service commands
	for name, svc := range resources {
		err := renderService(svc, updatable)
		if err != nil {
			fmt.Printf("ERROR: failed to save %q %q\n", name, err.Error())
			continue
		}
	}

	// render init function for services
	err := renderServiceInit(resources, updatable)
	if err != nil {
		fmt.Printf("ERROR: failed to save services init file %q\n", err.Error())
		panic(err)
	}
}

func renderService(svc CfnType, updatable []string) error {
	upd := ""
	if contains(updatable, svc.Command) {
		upd = "true"
	}
	serviceTemplate, err := template.ParseFiles(serviceTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create("cmd/services/" + svc.Command + ".go")
	if err != nil {
		return err
	}
	err = serviceTemplate.Execute(f, map[string]string{
		"Command":      svc.Command,
		"ShortDesc":    "TODO: short",
		"LongDesc":     "TODO: long",
		"CamelCommand": svc.Service,
		"Updateable":   upd,
	})
	if err != nil {
		return err
	}
	return nil
}

func renderServiceInit(resources map[string]CfnType, updatable []string) error {
	serviceList := map[string]bool{}
	// inefficient loop, should re-work data structures so we don't need this
	for _, svc := range resources {
		serviceList[svc.Service] = contains(updatable, svc.Command)
	}
	serviceInitTemplate, err := template.ParseFiles(serviceInitTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create("cmd/services/init.go")
	if err != nil {
		return err
	}
	err = serviceInitTemplate.Execute(f, serviceList)
	if err != nil {
		return err
	}
	return nil
}

func getSchemaFiles(basePath string) []string {
	var files []string
	basePath, err := filepath.Abs(basePath)
	if err != nil {
		panic(err)
	}
	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".json" {
			fmt.Printf("WARNING: found non-schema entity in the schema's path %q\n", path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
