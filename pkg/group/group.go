/*
Sniperkit-Bot
- Status: analyzed
*/

package group

import (
	yaml "gopkg.in/yaml.v2"

	berror "github.com/sniperkit/snk.fork.bulletin/pkg/error"
)

type Groups struct {
	Groups []Group `yaml:"groups"`
}

type Group struct {
	Name string `yaml:"name"`
	// optional fields
	Jobs      []string `yaml:"jobs"`
	Resources []string `yaml:"resources"`
}

func (g *Group) String() string {
	b, err := yaml.Marshal(*g)
	berror.CheckError(err)
	return string(b[:])
}

func GetGroupsFromString(data string) Groups {
	g := Groups{}
	err := yaml.Unmarshal([]byte(data), &g)
	berror.CheckError(err)
	return g
}
