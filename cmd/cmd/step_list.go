/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sniperkit/snk.fork.bulletin/pkg/bulletin_types"
	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
)

var stepListCmd = &cobra.Command{
	Use:   "list",
	Short: "list steps",
	RunE:  stepListRun,
}

var (
	stepInputs    string
	stepName      string
	listStepNames bool
)

func stepListRun(cmd *cobra.Command, args []string) error {
	datas := ioutils.ReadFile(stepInputs)
	steps := bulletin_types.GetStepsFromString(datas)
	if listStepNames {
		for _, j := range steps.Steps {
			fmt.Printf("%s\n", j.Name)
		}
	} else {
		for _, j := range steps.Steps {
			if stepName == "" {
				fmt.Printf("%+v\n", j.String())
				continue
			}
			if stepName != "" && stepName == j.Name {
				fmt.Printf("%+v\n", j.String())
				continue
			}
		}
	}
	return nil
}

func init() {
	stepCmd.AddCommand(stepListCmd)
	// required fields
	stepCmd.PersistentFlags().StringVarP(&stepInputs, "steps", "s", "", "path to a file that includes definitions of steps")
	stepCmd.PersistentFlags().StringVarP(&stepName, "step-name", "", "", "list all steps based on provided name")
	stepCmd.PersistentFlags().BoolVarP(&listStepNames, "names", "n", false, "list all step names")
}
