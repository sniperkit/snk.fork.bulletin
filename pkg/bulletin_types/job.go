package bulletin_types

import (
	berror "github.com/maplain/bulletin/pkg/error"
	"github.com/maplain/bulletin/pkg/job"
	yaml "gopkg.in/yaml.v2"
)

var stepss = `
steps:
- name: run-build
  type: step
  inputs:
  - target
  step:
    task: run-build
    file: ((target))/ci/tasks/build.yml
    input_mapping:
      git-pks-networking: ((target))
`

var decoraterss = `
decorators:
- name: merge-request
  type: task-decorator
  inputs:
  - target
  before:
  - put: ((target))
    params:
      repository: ((target))
      status: pending
  on_success:
    put: ((target))
    params:
      repository: ((target))
      status: success
  on_failure:
    put: ((target))
    params:
      repository: ((target))
      status: failed
`
var data = `
jobs:
- name: run-build-merge-request
  serial: true
  plan:
  - name: run-build
    inputs:
      target: git-pks-networking-merge-request
    decorators:
    - name: merge-request
      inputs:
        target: git-pks-networking-merge-request
`

type Jobs struct {
	Jobs []JobRef `yaml:"jobs"`
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
	Plan        []StepRef `yaml:"plan"`
	job.JobBase `yaml:",inline"`
	//Decorators []FuncRef `yaml:"decorators,omitempty"`
}

type StepRef struct {
	FuncRef    `yaml:",inline"`
	Decorators []FuncRef `yaml:"decorators,omitempty"`
}

func (j *StepRef) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

func (s *StepRef) DeRef(decs Decorators, ss Steps) ([]interface{}, error) {
	var res []interface{}
	step, err := ss.Populate(s.FuncRef)
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

func ConvertJobs(jobs Jobs, decs Decorators, ss Steps) job.Jobs {
	res := job.Jobs{}
	for _, j := range jobs.Jobs {
		res.Jobs = append(res.Jobs, ConvertJob(j, decs, ss))
	}
	return res
}

func ConvertJob(j JobRef, decs Decorators, ss Steps) job.Job {
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
		//st := ss.Populate(s.StepRef).Step()
		st, err := sref.DeRef(decs, ss)
		berror.CheckError(err)
		res.Plan = append(res.Plan, st)
	}
	return res
}
