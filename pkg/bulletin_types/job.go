/*
Sniperkit-Bot
- Status: analyzed
*/

package bulletin_types

import (
	template "github.com/maplain/yamltemplate"
	yaml "gopkg.in/yaml.v2"

	berror "github.com/sniperkit/snk.fork.bulletin/pkg/error"
	"github.com/sniperkit/snk.fork.bulletin/pkg/job"
)

type Jobs struct {
	Jobs  []JobRef `yaml:"jobs"`
	cache map[string]int
}

func (j *Jobs) AddDecorator(job, task string, d template.TemplateRef) {
	if j.cache == nil {
		j.cache = make(map[string]int)
		for i, jr := range j.Jobs {
			j.cache[jr.Name] = i
		}
	}
	i, ok := j.cache[job]
	if !ok {
		return
	}
	j.Jobs[i].AddDecorator(task, d)
}

func (j *Jobs) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

func GetJobsFromString(data string) Jobs {
	j := Jobs{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}

type JobRef struct {
	Plan          []StepRef `yaml:"plan"`
	job.JobBase   `yaml:",inline"`
	job.StepHooks `yaml:",inline"`
	Decorators    []template.TemplateRef `yaml:"decorators,omitempty"`
	cache         map[string]int
}

func (j *JobRef) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

func (j *JobRef) buildCache() {
	if j.cache == nil {
		j.cache = make(map[string]int)
		for i, s := range j.Plan {
			j.cache[s.Name] = i
		}
	}

}

func (j *JobRef) AddDecorator(task string, d template.TemplateRef) {
	j.buildCache()
	// this is a job decorator
	if task == "" {
		j.Decorators = append(j.Decorators, d)
	}
	i, ok := j.cache[task]
	if !ok {
		return
	}
	j.Plan[i].Decorators = append(j.Plan[i].Decorators, d)
}

type StepRef struct {
	template.TemplateRef `yaml:",inline"`
	Decorators           []template.TemplateRef `yaml:"decorators,omitempty"`
}

func (j *StepRef) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

func (s *StepRef) DeRef(decs Decorators, ss Steps) ([]interface{}, error) {
	var res []interface{}
	step, err := ss.Populate(s.TemplateRef)
	if err != nil {
		return res, err
	}
	i := step.GetStep()
	var ds []Decorator
	for _, d := range s.Decorators {
		dec, err := decs.Populate(d)
		if err != nil {
			return res, err
		}
		ds = append(ds, dec)
	}
	return Decorate(i, ds...), nil
}

func (jobs *Jobs) Convert(decs Decorators, ss Steps) job.Jobs {
	res := job.Jobs{}
	for _, j := range jobs.Jobs {
		res.Jobs = append(res.Jobs, j.Convert(decs, ss))
	}
	return res
}

func (j *JobRef) Convert(decs Decorators, ss Steps) job.Job {
	res := job.Job{}
	// copy job base
	res.Name = j.Name
	res.Serial = j.Serial
	res.BuildLogsToRetain = j.BuildLogsToRetain
	res.SerialGroups = j.SerialGroups
	res.MaxInFlight = j.MaxInFlight
	res.Public = j.Public
	res.DisableManualTrigger = j.DisableManualTrigger
	res.Interruptible = j.Interruptible

	// dereference step refs
	for _, sref := range j.Plan {
		// get real step
		st, err := sref.DeRef(decs, ss)
		berror.CheckError(err)
		// aggregate step is the first step
		for _, step := range st {
			b, err := yaml.Marshal(step)
			berror.CheckError(err)
			t, err := job.GetType(string(b[:]))
			berror.CheckError(err)
			switch t {
			case job.AggregateStepType:
				steps, err := job.GetAggregateStep(step)
				berror.CheckError(err)
				res.Plan = append([]interface{}{&steps}, res.Plan...)
			default:
				res.Plan = append(res.Plan, step)
			}
		}
	}

	// dereference job decorators
	for _, dref := range j.Decorators {
		d, err := decs.Populate(dref)
		berror.CheckError(err)
		if d.OnSuccess != nil {
			res.OnSuccess = d.OnSuccess
		}
		if d.OnFailure != nil {
			res.OnFailure = d.OnFailure
		}
		if d.OnAbort != nil {
			res.OnAbort = d.OnAbort
		}
		if d.Ensure != nil {
			res.Ensure = d.Ensure
		}
	}

	return res
}
