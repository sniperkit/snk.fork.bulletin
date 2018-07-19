package cmd

import (
	"fmt"
	"io/ioutil"

	berror "github.com/maplain/bulletin/pkg/error"
	"github.com/maplain/bulletin/pkg/group"
	"github.com/spf13/cobra"
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
	if pipeline != "" {
		dat, err := ioutil.ReadFile(pipeline)
		berror.CheckError(err)
		groups := group.GetGroups(string(dat))
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
	}
	return nil
}

func init() {
	groupCmd.AddCommand(groupListCmd)
	// required fields
	groupCmd.PersistentFlags().StringVarP(&groupName, "group-name", "", "", "list all groups based on provided name")
	groupCmd.PersistentFlags().BoolVarP(&listGroupNames, "names", "n", false, "list all groups names")
}
