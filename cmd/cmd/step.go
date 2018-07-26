package cmd

import "github.com/spf13/cobra"

var stepCmd = &cobra.Command{
	Use:   "step",
	Short: "managing all the steps defined in specified library",
}

func init() {
	rootCmd.AddCommand(stepCmd)
}
