package cmd

import (
	"fmt"
	"github.com/jaymccon/cloudctl/crudl"
	"github.com/jaymccon/cloudctl/data"
	"github.com/jaymccon/cloudctl/providers"
	awsProvider "github.com/jaymccon/cloudctl/providers/aws"
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
var ProviderConfigureCmds = map[string]*cobra.Command{}

var ServiceCreateCmds = map[string]map[string]*cobra.Command{}
var ServiceReadCmds = map[string]map[string]*cobra.Command{}
var ServiceUpdateCmds = map[string]map[string]*cobra.Command{}
var ServiceDeleteCmds = map[string]map[string]*cobra.Command{}
var ServiceListCmds = map[string]map[string]*cobra.Command{}
var ServiceConfigureCmds = map[string]map[string]*cobra.Command{}

var ResourceCreateCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceReadCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceUpdateCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceDeleteCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceListCmds = map[string]map[string]map[string]*cobra.Command{}
var ResourceConfigureCmds = map[string]map[string]map[string]*cobra.Command{}

func init() {
	schemas, err := data.GetSchemas()
	if err != nil {
		fmt.Printf("ERROR: %q\n", err)
		return
	}
	for provider, services := range *schemas {
		updatable := data.IsUpdatable(services)
		configurable := data.IsConfigurable(services)
		provider = strings.ToLower(provider)
		use, short, long := providers.GetCmdDetails(provider)
		ProviderCreateCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		ProviderReadCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		if updatable {
			ProviderUpdateCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		}
		ProviderDeleteCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		ProviderListCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		if configurable {
			ProviderConfigureCmds[provider] = &cobra.Command{Use: use, Short: short, Long: long}
		}

		CreateCmd.AddCommand(ProviderCreateCmds[provider])
		ReadCmd.AddCommand(ProviderReadCmds[provider])
		if updatable {
			UpdateCmd.AddCommand(ProviderUpdateCmds[provider])
		}
		DeleteCmd.AddCommand(ProviderDeleteCmds[provider])
		ListCmd.AddCommand(ProviderListCmds[provider])
		if configurable {
			ConfigureCmd.AddCommand(ProviderConfigureCmds[provider])
		}

		ServiceCreateCmds[provider] = map[string]*cobra.Command{}
		ServiceReadCmds[provider] = map[string]*cobra.Command{}
		if updatable {
			ServiceUpdateCmds[provider] = map[string]*cobra.Command{}
		}
		ServiceDeleteCmds[provider] = map[string]*cobra.Command{}
		ServiceListCmds[provider] = map[string]*cobra.Command{}
		if configurable {
			ServiceConfigureCmds[provider] = map[string]*cobra.Command{}
		}

		ResourceCreateCmds[provider] = map[string]map[string]*cobra.Command{}
		ResourceReadCmds[provider] = map[string]map[string]*cobra.Command{}
		if updatable {
			ResourceUpdateCmds[provider] = map[string]map[string]*cobra.Command{}
		}
		ResourceDeleteCmds[provider] = map[string]map[string]*cobra.Command{}
		ResourceListCmds[provider] = map[string]map[string]*cobra.Command{}
		if configurable {
			ResourceConfigureCmds[provider] = map[string]map[string]*cobra.Command{}
		}

		for service, resources := range services {
			service = strings.ToLower(service)
			updatable := data.IsUpdatable(services)
			configurable := data.IsConfigurable(services)
			use, short, long := providers.GetCmdDetails(service)
			ServiceCreateCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			ServiceReadCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			if updatable {
				ServiceUpdateCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			}
			ServiceDeleteCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			ServiceListCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			if configurable {
				ServiceConfigureCmds[provider][service] = &cobra.Command{Use: use, Short: short, Long: long}
			}

			ProviderCreateCmds[provider].AddCommand(ServiceCreateCmds[provider][service])
			ProviderReadCmds[provider].AddCommand(ServiceReadCmds[provider][service])
			if updatable {
				ProviderUpdateCmds[provider].AddCommand(ServiceUpdateCmds[provider][service])
			}
			ProviderDeleteCmds[provider].AddCommand(ServiceDeleteCmds[provider][service])
			ProviderListCmds[provider].AddCommand(ServiceListCmds[provider][service])
			if configurable {
				ProviderConfigureCmds[provider].AddCommand(ServiceConfigureCmds[provider][service])
			}

			ResourceCreateCmds[provider][service] = map[string]*cobra.Command{}
			ResourceReadCmds[provider][service] = map[string]*cobra.Command{}
			if updatable {
				ResourceUpdateCmds[provider][service] = map[string]*cobra.Command{}
			}
			ResourceDeleteCmds[provider][service] = map[string]*cobra.Command{}
			ResourceListCmds[provider][service] = map[string]*cobra.Command{}
			if configurable {
				ResourceConfigureCmds[provider][service] = map[string]*cobra.Command{}
			}

			for resource, schema := range resources {
				resource = strings.ToLower(resource)
				updatable := data.IsUpdatable(services)
				configurable := data.IsConfigurable(services)
				short = schema.Description
				long = schema.DocumentationUrl
				ResourceCreateCmds[provider][service][resource] = &cobra.Command{
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
						CreateEdit(cmd.Annotations["typeName"])
					},
				}
				ResourceReadCmds[provider][service][resource] = &cobra.Command{
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
						if len(args) != 1 {
							fmt.Println("read command requires an identifier to be supplied as a single argument")
						}
						crudl.ReadResource(cmd.Annotations["typeName"], args[0], noPrompts, async)
					},
					ValidArgsFunction: completeId,
				}
				if updatable {
					ResourceUpdateCmds[provider][service][resource] = &cobra.Command{
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
							fmt.Println("TODO: implementation")
						},
					}
				}
				ResourceDeleteCmds[provider][service][resource] = &cobra.Command{
					Use: resource,
					Annotations: map[string]string{
						"typeName": schema.TypeName,
						"resource": resource,
						"service":  service,
						"provider": provider,
					},
					Short: short,
					Long:  long,
					Args:  cobra.MinimumNArgs(1),
					Run: func(cmd *cobra.Command, args []string) {
						crudl.DeleteResources(cmd.Annotations["typeName"], args, noPrompts, async)
					},
					ValidArgsFunction: completeId,
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
				if configurable {
					ResourceConfigureCmds[provider][service][resource] = &cobra.Command{
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
							fmt.Println("TODO: implementation")
						},
					}
				}

				ServiceCreateCmds[provider][service].AddCommand(ResourceCreateCmds[provider][service][resource])
				ServiceReadCmds[provider][service].AddCommand(ResourceReadCmds[provider][service][resource])
				if updatable {
					ServiceUpdateCmds[provider][service].AddCommand(ResourceUpdateCmds[provider][service][resource])
				}
				ServiceDeleteCmds[provider][service].AddCommand(ResourceDeleteCmds[provider][service][resource])
				ServiceListCmds[provider][service].AddCommand(ResourceListCmds[provider][service][resource])
				if configurable {
					ServiceConfigureCmds[provider][service].AddCommand(ResourceConfigureCmds[provider][service][resource])
				}
			}
		}
	}
}

func completeId(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// this was in the example from cobra docs, but it's not clear whether it's needed
	//if len(args) != 0 {
	//	return nil, cobra.ShellCompDirectiveNoFileComp
	//}
	fmt.Println(cmd.Parent().Parent().Parent().Name())
	resources, err := awsProvider.ListResource(cmd.Annotations["typeName"])
	if err != nil {
		fmt.Printf("ERROR: %q\n", err.Error())
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	headers := append([]interface{}{"Identifier"}, crudl.GetTableHeaders(*resources)...)
	var completeList []string
	for _, r := range *resources {
		row := crudl.GetRow(r, headers)
		completeStr := row[0].(string) + "\t"
		for _, i := range row[1:] {
			completeStr = completeStr + i.(string) + " "
		}
		completeList = append(completeList, completeStr)
	}
	return completeList, cobra.ShellCompDirectiveNoFileComp
}
