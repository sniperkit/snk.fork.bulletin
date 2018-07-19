package cmd

import (
	"fmt"
	"io/ioutil"

	berror "github.com/maplain/bulletin/pkg/error"
	"github.com/maplain/bulletin/pkg/job"
	"github.com/spf13/cobra"
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
	if pipeline != "" {
		dat, err := ioutil.ReadFile(pipeline)
		berror.CheckError(err)
		jobs := job.GetJobs(string(dat))
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
	}
	return nil
}

func init() {
	jobCmd.AddCommand(jobListCmd)
	// required fields
	jobCmd.PersistentFlags().StringVarP(&jobName, "job-name", "", "", "list all jobs based on provided name")
	jobCmd.PersistentFlags().BoolVarP(&listJobNames, "names", "n", false, "list all jobs names")
}
