package cmd

import "github.com/spf13/cobra"

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "r",
	RunE:  resourceRun,
}

func resourceRun(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(resourceCmd)
	// required fields
	//	resourceCmd.PersistentFlags().StringVarP(&list, "list", "r", "", "list all resources")
}
