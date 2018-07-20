package cmd

import "github.com/spf13/cobra"

var resourceTypeCmd = &cobra.Command{
	Use:   "resource-types",
	Short: "rt",
}

func init() {
	rootCmd.AddCommand(resourceTypeCmd)
}
