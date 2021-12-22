package crudl

import (
	"fmt"
	awsProvider "github.com/jaymccon/cloudctl/providers/aws"
)

func CreateResource(typeName string, properties string, noPrompts bool, async bool) {
	if noPrompts {
		fmt.Println("--no-prompt flag set, skipping confirmation.")
	}
	cc, err := awsProvider.NewCcClient()
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
	}
	awsProvider.AsyncCcCreateResource(*cc, typeName, properties, async)
}
