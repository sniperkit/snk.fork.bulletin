/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
	ppl "github.com/sniperkit/snk.fork.bulletin/pkg/pipeline"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update provided pipeline based on given files. Note: following one will overwrite pre-existing one",
	RunE:  updateRun,
}

var (
	destination string
)

func updateRun(cmd *cobra.Command, args []string) error {
	datas := ioutils.ReadFileDefaultStdin(pipeline)
	// update resource types
	pp := ppl.GetPipelineFromString(datas)
	for _, f := range args {
		tp := ppl.GetPipelineFromString(ioutils.ReadFile(f))
		pp.UpdateWith(tp)
	}
	if destination != "" {
		err := ioutil.WriteFile(destination, []byte(pp.String()), 0644)
		if err != nil {
			return err
		}
	} else {
		io.WriteString(os.Stdout, pp.String())
	}

	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().StringVarP(&destination, "destination", "d", "", "destination file to write updated pipeline yaml")
}
