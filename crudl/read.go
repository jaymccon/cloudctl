package crudl

import (
	"fmt"
	awsProvider "github.com/jaymccon/cloudctl/providers/aws"
)

func ReadResource(typeName string, id string, noPrompts bool, async bool) {
	cc, err := awsProvider.NewCcClient()
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
	}
	awsProvider.AsyncCcReadResource(*cc, typeName, id, async)
}
