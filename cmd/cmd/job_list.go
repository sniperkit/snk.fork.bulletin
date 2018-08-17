/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
	"github.com/sniperkit/snk.fork.bulletin/pkg/job"
)

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "list jobs",
	RunE:  jobListRun,
}

var (
	jobName      string
	listJobNames bool
)

func jobListRun(cmd *cobra.Command, args []string) error {
	datas := ioutils.ReadFileDefaultStdin(pipeline)
	jobs := job.GetJobsFromString(datas)
	if listJobNames {
		for _, j := range jobs.Jobs {
			fmt.Printf("%s\n", j.Name)
		}
	} else {
		for _, j := range jobs.Jobs {
			if jobName == "" {
				fmt.Printf("%+v\n", j.String())
				continue
			}
			if jobName != "" && jobName == j.Name {
				fmt.Printf("%+v\n", j.String())
				continue
			}
		}
	}
	return nil
}

func init() {
	jobCmd.AddCommand(jobListCmd)
	// required fields
	jobCmd.PersistentFlags().StringVarP(&jobName, "job-name", "", "", "list all jobs based on provided name")
	jobCmd.PersistentFlags().BoolVarP(&listJobNames, "names", "n", false, "list all jobs names")
}
