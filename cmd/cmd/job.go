/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "managing all the jobs defined in pipeline.yml",
}

func init() {
	rootCmd.AddCommand(jobCmd)
}
