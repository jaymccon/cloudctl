package cmd

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates cloud resources",
	Run: func(cmd *cobra.Command, args []string) {
		yamlIface := map[string]interface{}{
			"ReqProp2": "integer", "ReqProp1": "string", // required props
			"OptProp2": "myDefaultValue", "OptProp1": true, "OptProp3": 1, // optional_with_default
			"OptNoDefProp2": map[string]interface{}{"nested": "string"}, "OptNoDefProp1": "boolean", "OptNoDefProp3": "integer", // optional without defaults
		}
		comments := yaml.CommentMap{
			"$.ReqProp1":      yaml.HeadComment("Required properties"),
			"$.OptProp1":      yaml.HeadComment("Optional Properties with defaults"),
			"$.OptProp3":      FootComment("Optional Properties without defaults"),
			"$.OptNoDefProp1": EntireLineComment(),
			"$.OptNoDefProp2": EntireLineComment(),
			"$.OptNoDefProp3": EntireLineComment(),
		}
		yamlString, err := yaml.MarshalWithOptions(yamlIface, yaml.WithComment(comments))
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			return
		}
		// This is what we would want the marshalled output to look like:
		yamlString = []byte(
			`# Required properties
ReqProp1: string
ReqProp2: integer
# Optional Properties with defaults
OptProp1: true
OptProp2: "myDefaultValue"
OptProp3: 1
# Optional Properties without defaults
# OptNoDefProp1: boolean
# OptNoDefProp2:
#   nested: string
# OptNoDefProp3: integer
`)
		Edit(string(yamlString), "yml")
	},
}

func FootComment(_ string) *yaml.Comment {
	return &yaml.Comment{
		Texts:    []string{"TODO: implement FootComment"},
		Position: 0,
	}
}

func EntireLineComment() *yaml.Comment {
	return &yaml.Comment{
		Texts:    []string{"TODO: implement EntireLineComment"},
		Position: 0,
	}
}

func init() {
	RootCmd.AddCommand(CreateCmd)
}
