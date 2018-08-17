/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var decoratorCmd = &cobra.Command{
	Use:   "decorator",
	Short: "operations on decorators",
}

func init() {
	rootCmd.AddCommand(decoratorCmd)
}
