/*
Sniperkit-Bot
- Status: analyzed
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sniperkit/snk.fork.bulletin/pkg/group"
	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
)

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "list groups",
	RunE:  groupListRun,
}

var (
	groupName      string
	listGroupNames bool
)

func groupListRun(cmd *cobra.Command, args []string) error {
	datas := ioutils.ReadFileDefaultStdin(pipeline)
	groups := group.GetGroupsFromString(datas)
	if listGroupNames {
		for _, g := range groups.Groups {
			fmt.Printf("%s\n", g.Name)
		}
	} else {
		for _, g := range groups.Groups {
			if groupName == "" {
				fmt.Printf("%+v\n", g.String())
				continue
			}
			if groupName != "" && groupName == g.Name {
				fmt.Printf("%+v\n", g.String())
				continue
			}
		}
	}
	return nil
}

func init() {
	groupCmd.AddCommand(groupListCmd)
	// required fields
	groupCmd.PersistentFlags().StringVarP(&groupName, "group-name", "", "", "list all groups based on provided name")
	groupCmd.PersistentFlags().BoolVarP(&listGroupNames, "names", "n", false, "list all groups names")
}
