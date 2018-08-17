/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "managing all the groups defined in pipeline.yml",
}

func init() {
	rootCmd.AddCommand(groupCmd)
}
