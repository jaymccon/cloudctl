package commands

import (
	"fmt"
	"github.com/jaymccon/cloudctl/crudl"
	"github.com/jaymccon/cloudctl/data"
	"github.com/jaymccon/cloudctl/providers"
	"github.com/spf13/cobra"
	"strings"
)

func Dummy() {
	return
}

var ProviderCreateCmds = map[string]*cobra.Command{}
var ProviderReadCmds = map[string]*cobra.Command{}
var ProviderUpdateCmds = map[string]*cobra.Command{}
var ProviderDeleteCmds = map[string]*cobra.Command{}
var ProviderListCmds = map[string]*cobra.Command{}

var ServiceCreateCmds = map[string]map[string]*cobra.Command{}
var ServiceReadCmds = map[string]map[string]*cobra.Command{}
var ServiceUpdateCmds = map[string]map[string]*cobra.Command{}
var ServiceDeleteCmds = map[string]map[string]*cobra.Command{}
var ServiceListCmds = map[string]map[string]*cobra.Command{}

var ResourceCreateCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceReadCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceUpdateCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceDeleteCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceListCmds = map[string]map[string]map[string]*cobra.Command{}

func init() {
	schemas, err := data.GetSchemas()
	if err != nil {
		fmt.Printf("ERROR: %q\n", err)
		return
	}
	for provider, services := range *schemas {
		provider = strings.ToLower(provider)
		use, short, long := providers.GetCmdDetails(provider)
		ProviderCreateCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		ProviderReadCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		ProviderUpdateCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		ProviderDeleteCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		ProviderListCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}

		CreateCmd.AddCommand(ProviderCreateCmds[provider])
		ReadCmd.AddCommand(ProviderReadCmds[provider])
		UpdateCmd.AddCommand(ProviderUpdateCmds[provider])
		DeleteCmd.AddCommand(ProviderDeleteCmds[provider])
		ListCmd.AddCommand(ProviderListCmds[provider])

		ServiceCreateCmds[provider] = map[string]*cobra.Command{}
		ServiceReadCmds[provider] = map[string]*cobra.Command{}
		ServiceUpdateCmds[provider] = map[string]*cobra.Command{}
		ServiceDeleteCmds[provider] = map[string]*cobra.Command{}
		ServiceListCmds[provider] = map[string]*cobra.Command{}

		ResourceCreateCmds[provider] = map[string]map[string]*cobra.Command{}
		ResourceReadCmds[provider] = map[string]map[string]*cobra.Command{}
		ResourceUpdateCmds[provider] = map[string]map[string]*cobra.Command{}
		ResourceDeleteCmds[provider] = map[string]map[string]*cobra.Command{}
		ResourceListCmds[provider] = map[string]map[string]*cobra.Command{}

		for service, resources := range services {
			service = strings.ToLower(service)
			use, short, long := providers.GetCmdDetails(service)
			ServiceCreateCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			ServiceReadCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			ServiceUpdateCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			ServiceDeleteCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			ServiceListCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}

			ProviderCreateCmds[provider].AddCommand(ServiceCreateCmds[provider][service])
			ProviderReadCmds[provider].AddCommand(ServiceReadCmds[provider][service])
			ProviderUpdateCmds[provider].AddCommand(ServiceUpdateCmds[provider][service])
			ProviderDeleteCmds[provider].AddCommand(ServiceDeleteCmds[provider][service])
			ProviderListCmds[provider].AddCommand(ServiceListCmds[provider][service])

			ResourceCreateCmds[provider][service] = map[string]*cobra.Command{}
			ResourceReadCmds[provider][service] = map[string]*cobra.Command{}
			ResourceUpdateCmds[provider][service] = map[string]*cobra.Command{}
			ResourceDeleteCmds[provider][service] = map[string]*cobra.Command{}
			ResourceListCmds[provider][service] = map[string]*cobra.Command{}

			for resource, schema := range resources {
				resource = strings.ToLower(resource)
				short = schema.Description
				long = schema.DocumentationUrl
				ResourceCreateCmds[provider][service][resource] = &cobra.Command{
					Use:   resource,
					Short: short,
					Long:  long,
					Run: func(cmd *cobra.Command, args []string) {
						fmt.Println("TODO: implementation")
					},
				}
				ResourceReadCmds[provider][service][resource] = &cobra.Command{
					Use:   resource,
					Short: short,
					Long:  long,
					Run: func(cmd *cobra.Command, args []string) {
						fmt.Println("TODO: implementation")
					},
				}
				ResourceUpdateCmds[provider][service][resource] = &cobra.Command{
					Use:   resource,
					Short: short,
					Long:  long,
					Run: func(cmd *cobra.Command, args []string) {
						fmt.Println("TODO: implementation")
					},
				}
				ResourceDeleteCmds[provider][service][resource] = &cobra.Command{
					Use:   resource,
					Short: short,
					Long:  long,
					Run: func(cmd *cobra.Command, args []string) {
						fmt.Println("TODO: implementation")
					},
				}
				ResourceListCmds[provider][service][resource] = &cobra.Command{
					Use: resource,
					Annotations: map[string]string{
						"typeName": schema.TypeName,
						"resource": resource,
						"service":  service,
						"provider": provider,
					},
					Short: short,
					Long:  long,
					Run: func(cmd *cobra.Command, args []string) {
						crudl.ListResource(cmd.Annotations["typeName"])
					},
				}

				ServiceCreateCmds[provider][service].AddCommand(ResourceCreateCmds[provider][service][resource])
				ServiceReadCmds[provider][service].AddCommand(ResourceReadCmds[provider][service][resource])
				ServiceUpdateCmds[provider][service].AddCommand(ResourceUpdateCmds[provider][service][resource])
				ServiceDeleteCmds[provider][service].AddCommand(ResourceDeleteCmds[provider][service][resource])
				ServiceListCmds[provider][service].AddCommand(ResourceListCmds[provider][service][resource])
			}
		}
	}
}
