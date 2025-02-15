package crudl

import (
	"fmt"
	awsProvider "github.com/jaymccon/cloudctl/providers/aws"
)

func DeleteResources(typeName string, ids []string, noPrompts bool, async bool) {
	if !noPrompts {
		if !confirm(fmt.Sprintf("Are you sure you want to delete %s resources with identifiers %s", typeName, ids)) {
			fmt.Println("Exiting without deleting anything.")
			return
		}
		if noPrompts {
			fmt.Println("--no-prompt flag set, skipping confirmation.")
		}
	}
	cc, err := awsProvider.NewCcClient()
	if err != nil {
		fmt.Printf("ERROR: %q", err.Error())
	}
	awsProvider.AsyncCcDeleteResource(*cc, typeName, ids, async)
}
