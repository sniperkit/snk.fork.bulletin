/*
Sniperkit-Bot
- Status: analyzed
*/

package bulletin_types

import (
	"fmt"
	"strconv"

	yaml "gopkg.in/yaml.v2"

	berror "github.com/sniperkit/snk.fork.bulletin/pkg/error"
	"github.com/sniperkit/snk.fork.bulletin/pkg/job"
)

type Deps struct {
	Deps []Dep `yaml:"deps"`
}

func (ds *Deps) SetDefault() Deps {
	for i, d := range ds.Deps {
		ds.Deps[i] = d.SetDefault()
	}
	return *ds
}

func (ds *Deps) String() string {
	b, err := yaml.Marshal(*ds)
	berror.CheckError(err)
	return string(b[:])
}

type Dep struct {
	Name       string         `yaml:"name"`
	RequiredBy []Requirements `yaml:"required_by"`
}

func (dep *Dep) AddResource(jobs job.Jobs) {
	for _, r := range dep.RequiredBy {
		r.AddResource(dep.Name, jobs)
	}
}

type Requirements []DepJobRef

func (req *Requirements) AddResource(name string, jobs job.Jobs) error {
	for i, ref := range *req {
		oldj, err := jobs.GetJob(ref.Name)
		if err != nil {
			return err
		}
		oldGetStepI := oldj.GetStepByTypeAndName(job.GetStepType, name)
		if oldGetStepI != nil {
			//TODO: overwrite other fields?
			if i >= 1 {
				switch v := oldGetStepI.(type) {
				case *job.GetStep:
					v.Passed = append(v.Passed, (*req)[i-1].Name)
				}
			}
			continue
		}
		getStep := job.GetStep{
			Get:     name,
			Version: ref.Version,
			Params:  ref.Params,
			Trigger: ref.Trigger,
		}
		if i >= 1 {
			getStep.Passed = append(getStep.Passed, (*req)[i-1].Name)
		}
		if ref.aggregatableB {
			aggregateSteps := oldj.GetStepsByType(job.AggregateStepType)
			if len(aggregateSteps) == 0 {
				aggregateStep := &job.AggregateStep{
					Aggregate: []interface{}{
						getStep,
					},
				}
				oldj.Plan = append([]interface{}{aggregateStep}, oldj.Plan...)
				oldj.AddStepByType(job.AggregateStepType, aggregateStep)
			} else {
				aggregateStepI := aggregateSteps[0]
				switch v := aggregateStepI.(type) {
				case *job.AggregateStep:
					v.Aggregate = append(v.Aggregate, getStep)
				default:
					fmt.Printf("!!!\nerror: unsupported type %T\n", v)
				}
			}
		} else {
			oldj.Plan = append(oldj.Plan, getStep)
		}
		oldj.AddStepByTypeAndName(job.GetStepType, name, getStep)
		err = jobs.UpdateJob(oldj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dep) SetDefault() Dep {
	for i, a1 := range d.RequiredBy {
		for j, a2 := range a1 {
			a1[j] = a2.SetDefault()
		}
		d.RequiredBy[i] = a1
	}
	return *d
}

func (d *Dep) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

type DepJobRef struct {
	Name               string                 `yaml:"name,omitempty"`
	Version            string                 `yaml:"version,omitempty"`
	Params             map[string]interface{} `yaml:"params,omitempty"`
	Trigger            bool                   `yaml:"trigger,omitempty"`
	AggregatableString string                 `yaml:"aggregatable,omitempty"`
	aggregatableB      bool
}

func (d *DepJobRef) SetDefault() DepJobRef {
	res := d
	// no need to do this, concourse will handle it
	//if res.Version == "" {
	//	res.Version = "latest"
	//}
	if res.AggregatableString == "" {
		res.aggregatableB = true
	} else {
		b, err := strconv.ParseBool(res.AggregatableString)
		berror.CheckError(err)
		res.aggregatableB = b
	}
	return *res
}

func (d *DepJobRef) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

func GetDepsFromString(data string) Deps {
	d := Deps{}
	err := yaml.Unmarshal([]byte(data), &d)
	berror.CheckError(err)
	return d.SetDefault()
}
