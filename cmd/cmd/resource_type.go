package cmd

import "github.com/spf13/cobra"

var resourceTypeCmd = &cobra.Command{
	Use:   "resource-types",
	Short: "rt",
	RunE:  resourceTypeRun,
}

func resourceTypeRun(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(resourceTypeCmd)
	// required fields
	//resourceTypeCmd.PersistentFlags().StringVarP(&list, "list", "r", "", "list all resources")
}
