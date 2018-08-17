/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/PKS/pks-networking/pkg/printer"

	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
	"github.com/sniperkit/snk.fork.bulletin/pkg/resource"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bulletin",
	Short: "A binary to compose Concourse pipeline",
	Long:  `A binary to compose Concourse pipelineu using referenced resources`,
	Run:   rootRun,
}

var (
	registry string
	pipeline string
	log      *printer.Printer
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// print to stderr by default
	log = printer.New(os.Stderr)
	// required fields
	//TODO: combine these two options
	// rootCmd.PersistentFlags().StringVarP(&registry, "registry", "r", "", "folder that include files for all pipeline yaml files")
	rootCmd.PersistentFlags().StringVarP(&pipeline, "pipeline", "p", "", "a pipeline yaml file you want to parse")
	// rootCmd.MarkPersistentFlagRequired("")
	// optional fields
	//	rootCmd.PersistentFlags().BoolVarP(&readOnly, "read-only", "r", true, "Read only mode")
	//	rootCmd.PersistentFlags().BoolVar(&pks, "pks", false, "removes all pks created resources as well")
	//	rootCmd.Flags().StringVar(&floatingIPPoolID, "floating-ip-pool-id", "", "UUID of the floating IP pool configured for the cluster")
	//	rootCmd.Flags().StringVar(&ipBlockID, "ip-block-id", "", "UUID of the IP block configured for the cluster")
}

func rootRun(cmd *cobra.Command, args []string) {
	datas := ioutils.ReadFileDefaultStdin(pipeline)
	resources := resource.GetResourcesFromString(datas)
	for _, r := range resources.Resources {
		fmt.Printf("%+v\n", r.String())
	}
	resourceTypes := resource.GetResourceTypesFromString(datas)
	for _, r := range resourceTypes.ResourceTypes {
		fmt.Printf("%+v\n", r.String())
	}
}
