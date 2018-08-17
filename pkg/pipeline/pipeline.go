/*
Sniperkit-Bot
- Status: analyzed
*/

package pipeline

import (
	yaml "gopkg.in/yaml.v2"

	berror "github.com/sniperkit/snk.fork.bulletin/pkg/error"
	"github.com/sniperkit/snk.fork.bulletin/pkg/group"
	"github.com/sniperkit/snk.fork.bulletin/pkg/job"
	"github.com/sniperkit/snk.fork.bulletin/pkg/resource"
)

type Pipeline struct {
	resource.Resources     `yaml:",inline"`
	resource.ResourceTypes `yaml:",inline"`
	group.Groups           `yaml:",inline"`
	job.Jobs               `yaml:",inline"`
}

func (p *Pipeline) String() string {
	b, err := yaml.Marshal(*p)
	berror.CheckError(err)
	return string(b[:])
}

func GetPipelineFromString(data string) Pipeline {
	g := Pipeline{}
	err := yaml.Unmarshal([]byte(data), &g)
	berror.CheckError(err)
	return g
}

func (p *Pipeline) UpdateWith(n Pipeline) {
	p.Resources = p.Resources.UpdateWith(n.Resources)
	//	p.ResourceTypes.UpdateWith(n.ResourceTypes)
	//	p.Groups.UpdateWith(n.Groups)
	//	p.Jobs.UpdateWith(n.Jobs)
}
