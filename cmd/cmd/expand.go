package cmd

import (
	"fmt"

	"github.com/maplain/bulletin/pkg/bulletin_types"
	"github.com/maplain/bulletin/pkg/ioutils"
	"github.com/spf13/cobra"
)

var expandCmd = &cobra.Command{
	Use:   "expand",
	Short: "expand bulletin pipeline yaml file to concourse configuration yaml file",
	RunE:  expandRun,
}

var expandTarget = "."

func expandRun(cmd *cobra.Command, args []string) error {
	datas := ioutils.ReadFileDefaultStdin(pipeline)
	// update resource types
	jobs := bulletin_types.GetJobsFromString(datas)
	savedDecs := bulletin_types.GetLocalDecorators(expandTarget)
	savedSteps := bulletin_types.GetLocalSteps(expandTarget)

	globalDecs := bulletin_types.GetStepDecoratorDefsFromString(datas)
	for _, gdec := range globalDecs.Decorators {
		for _, jt := range gdec.Decorate {
			job, task := gdec.GetJobTask(jt)
			jobs.AddDecorator(job, task, gdec.TemplateRef)
		}
	}

	cjobs := jobs.Convert(savedDecs, savedSteps)
	fmt.Printf("%s\n", cjobs.String())
	return nil
}

func init() {
	rootCmd.AddCommand(expandCmd)
	expandCmd.PersistentFlags().StringVarP(&expandTarget, "target", "t", "", "a folder to persist pipeline components definitions")
}
