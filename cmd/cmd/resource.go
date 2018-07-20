package cmd

import "github.com/spf13/cobra"

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "r",
}

func init() {
	rootCmd.AddCommand(resourceCmd)
}
