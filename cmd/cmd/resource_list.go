/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
	"github.com/sniperkit/snk.fork.bulletin/pkg/resource"
)

var resourceListCmd = &cobra.Command{
	Use:   "list",
	Short: "l",
	RunE:  resourceListRun,
}

var (
	resourceName      string
	resourceType      string
	listResourceNames bool
	listResourceTypes bool
)

func resourceListRun(cmd *cobra.Command, args []string) error {

	datas := ioutils.ReadFileDefaultStdin(pipeline)
	resources := resource.GetResourcesFromString(datas)
	printTypes := make(map[string]interface{})
	if listResourceNames {
		for _, r := range resources.Resources {
			fmt.Printf("%s\n", r.Name)
		}
	} else if listResourceTypes {
		for _, r := range resources.Resources {
			printTypes[r.Type] = nil
		}
	} else {
		for _, r := range resources.Resources {
			if resourceType == "" && resourceName == "" {
				fmt.Printf("%+v\n", r.String())
				continue
			}
			if resourceType != "" && resourceType == r.Type {
				if resourceName != "" && resourceName == r.Name {
					fmt.Printf("%+v\n", r.String())
					continue
				}
				if resourceName == "" {
					fmt.Printf("%+v\n", r.String())
					continue
				}
			}
			if resourceName != "" && resourceName == r.Name {
				if resourceType == "" {
					fmt.Printf("%+v\n", r.String())
					continue
				}
			}
		}
	}
	for k, _ := range printTypes {
		fmt.Printf("%s\n", k)
	}
	return nil
}

func init() {
	resourceCmd.AddCommand(resourceListCmd)
	// required fields
	resourceCmd.PersistentFlags().StringVarP(&resourceName, "resource-name", "", "", "list all resources based on provided name")
	resourceCmd.PersistentFlags().StringVarP(&resourceType, "resource-type", "", "", "list all resources based on provided type")
	resourceCmd.PersistentFlags().BoolVarP(&listResourceNames, "name", "n", false, "list all resources names")
	resourceCmd.PersistentFlags().BoolVarP(&listResourceTypes, "type", "t", false, "list all resources types")
}
