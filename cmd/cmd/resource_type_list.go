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

var resourceTypeListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all resource types",
	RunE:  resourceTypeListRun,
}

var (
	resourceTypeType     string
	resourceTypeName     string
	listResourceTypeType bool
	listResourceTypeName bool
)

func resourceTypeListRun(cmd *cobra.Command, args []string) error {
	datas := ioutils.ReadFileDefaultStdin(pipeline)
	resourceTypes := resource.GetResourceTypesFromString(datas)
	printTypes := make(map[string]interface{})
	if listResourceTypeName {
		for _, r := range resourceTypes.ResourceTypes {
			fmt.Printf("%s\n", r.Name)
		}
	} else if listResourceTypeType {
		for _, r := range resourceTypes.ResourceTypes {
			printTypes[r.Type] = nil
		}
	} else {
		for _, r := range resourceTypes.ResourceTypes {
			if resourceTypeType == "" && resourceTypeName == "" {
				fmt.Printf("%+v\n", r.String())
				continue
			}
			if resourceTypeType != "" && resourceTypeType == r.Type {
				if resourceTypeName != "" && resourceTypeName == r.Name {
					fmt.Printf("%+v\n", r.String())
					continue
				}
				if resourceTypeName == "" {
					fmt.Printf("%+v\n", r.String())
					continue
				}
			}
			if resourceTypeName != "" && resourceTypeName == r.Name {
				if resourceTypeType == "" {
					fmt.Printf("%+v\n", r.String())
					continue
				}
			}
		}
	}
	for t, _ := range printTypes {
		fmt.Printf("%s\n", t)
	}
	return nil
}

func init() {
	resourceTypeCmd.AddCommand(resourceTypeListCmd)
	// required fields
	resourceTypeCmd.PersistentFlags().StringVarP(&resourceTypeType, "rt-type", "", "", "list all resource types based on provided type")
	resourceTypeCmd.PersistentFlags().StringVarP(&resourceTypeName, "rt-name", "", "", "list all resource types based on provided name")
	resourceTypeCmd.PersistentFlags().BoolVarP(&listResourceTypeType, "type", "t", false, "list all resource types types")
	resourceTypeCmd.PersistentFlags().BoolVarP(&listResourceTypeName, "name", "n", false, "list all resource types names")
}
