package cmd

import (
	"fmt"

	"github.com/maplain/bulletin/pkg/bulletin_types"
	"github.com/maplain/bulletin/pkg/ioutils"
	"github.com/spf13/cobra"
)

var decoratorListCmd = &cobra.Command{
	Use:   "list",
	Short: "l",
	RunE:  decoratorListRun,
}

var (
	decoratorInputs    string
	decoratorName      string
	listDecoratorNames bool
)

func decoratorListRun(cmd *cobra.Command, args []string) error {

	datas := ioutils.ReadFile(decoratorInputs)
	decorators := bulletin_types.GetDecoratorsFromString(datas)
	if listDecoratorNames {
		for _, r := range decorators.Decorators {
			fmt.Printf("%s\n", r.Name)
		}
	} else {
		for _, r := range decorators.Decorators {
			if decoratorName == "" {
				fmt.Printf("%s\n", r.String())
			} else if decoratorName == r.Name {
				fmt.Printf("%s\n", r.String())
			}
		}
	}
	return nil
}

func init() {
	decoratorCmd.AddCommand(decoratorListCmd)
	// required fields
	decoratorCmd.PersistentFlags().StringVarP(&decoratorInputs, "decorators", "d", "", "path to a file that includes definitions of decorators")
	decoratorCmd.PersistentFlags().StringVarP(&decoratorName, "decorator-name", "", "", "list decorator definition based on provided name")
	decoratorCmd.PersistentFlags().BoolVarP(&listDecoratorNames, "name", "n", false, "list all decorator names")
}
