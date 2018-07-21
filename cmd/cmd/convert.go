package cmd

import (
	"github.com/maplain/bulletin/pkg/ioutils"
	"github.com/maplain/bulletin/pkg/resource"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert concourse pipeline yaml file to bulletin configuration yaml file",
	RunE:  convertRun,
}

var target = "."

const (
	resourcesDir  = "resources"
	resourcesFile = "resources.yml"
)

func convertRun(cmd *cobra.Command, args []string) error {
	if pipeline == "" {
		return nil
	}
	p := ioutils.ReadFile(pipeline)

	// update resource types
	rt := resource.GetResourceTypesFromString(p)
	savedRT := resource.GetLocalResourceTypes(target)
	for _, r := range rt.ResourceTypes {
		savedRT.Add(r)
	}
	err := resource.SaveResourceTypesLocally(target, resource.ResourceTypes{savedRT.Get()})
	if err != nil {
		log.Warn("failed to save resource types locally")
	}

	// update resources
	rs := resource.GetResourcesFromString(p)
	savedRs := resource.GetLocalResources(target)
	for _, r := range rs.Resources {
		savedRs.Add(r)
	}
	err = resource.SaveResourcesLocally(target, resource.Resources{savedRs.Get()})
	if err != nil {
		log.Warn("failed to save resources locally")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "a folder to persist pipeline components definitions")
}
