package cmd

import "github.com/spf13/cobra"

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "managing all the jobs defined in pipeline.yml",
	RunE:  jobRun,
}

func jobRun(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(jobCmd)
	// required fields
	//	resourceCmd.PersistentFlags().StringVarP(&list, "list", "r", "", "list all resources")
}
