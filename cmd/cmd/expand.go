/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sniperkit/snk.fork.bulletin/pkg/bulletin_types"
	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
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
	deps := bulletin_types.GetDepsFromString(datas)
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
	for _, d := range deps.Deps {
		d.AddResource(cjobs)
	}
	fmt.Printf("%s\n", cjobs.String())
	return nil
}

func init() {
	rootCmd.AddCommand(expandCmd)
	expandCmd.PersistentFlags().StringVarP(&expandTarget, "target", "t", "", "a folder to persist pipeline components definitions")
}
