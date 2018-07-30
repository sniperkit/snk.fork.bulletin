package job

import (
	"log"

	berror "github.com/maplain/bulletin/pkg/error"
	"github.com/maplain/bulletin/pkg/types"
	yaml "gopkg.in/yaml.v2"
)

const (
	JobNotExistError  types.InternalError = "could not find specified job"
	NoTypedStepsError types.InternalError = "no steps defined for specified type"
	NoStepError       types.InternalError = "no such step"
)

type Jobs struct {
	Jobs  []Job `yaml:"jobs"`
	cache map[string]int
}

func (j *Jobs) buildCache(force bool) {
	if j.cache == nil || force {
		j.cache = make(map[string]int)
		for i, job := range j.Jobs {
			j.cache[job.Name] = i
		}
	}
}

func (j *Jobs) GetJob(name string) (Job, error) {
	j.buildCache(false)
	i, ok := j.cache[name]
	if !ok {
		return Job{}, JobNotExistError
	}
	return j.Jobs[i], nil
}

func (j *Jobs) UpdateJob(newj Job) error {
	j.buildCache(false)
	i, ok := j.cache[newj.Name]
	if !ok {
		return JobNotExistError
	}
	j.Jobs[i] = newj
	return nil
}

func (j *Jobs) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

type JobBase struct {
	Name string `yaml:"name"`
	// optional fields
	Serial               bool     `yaml:"serial,omitempty"`
	BuildLogsToRetain    int      `yaml:"build_logs_to_retain,omitempty"`
	SerialGroups         []string `yaml:"serial_groups,omitempty"`
	MaxInFlight          int      `yaml:"max_in_flight,omitempty"`
	Public               bool     `yaml:"public,omitempty"`
	DisableManualTrigger bool     `yaml:"disable_manual_trigger,omitempty"`
	Interruptible        bool     `yaml:"interruptible,omitempty"`
}

func (j *JobBase) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

type Job struct {
	Plan      []interface{} `yaml:"plan"`
	JobBase   `yaml:",inline"`
	StepHooks `yaml:",inline"`
	stepCache map[Type]map[string]interface{}
	typeCache map[Type][]interface{}
}

func (j *Job) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

func (j *Job) buildCache() {
	if j.typeCache == nil {
		j.typeCache = make(map[Type][]interface{})
	}
	if j.stepCache == nil {
		j.stepCache = make(map[Type]map[string]interface{})
		for _, p := range j.Plan {
			s, err := yaml.Marshal(&p)
			berror.CheckError(err)
			t, _ := GetType(string(s))
			if j.stepCache[t] == nil {
				j.stepCache[t] = make(map[string]interface{})
			}
			switch t {
			case PutStepType:
				fallthrough
			case GetStepType:
				fallthrough
			case TaskStepType:
				name, err := GetStepName(p)
				berror.CheckError(err)
				j.stepCache[t][name] = &p
				j.typeCache[t] = append(j.typeCache[t], &p)
			case DoStepType:
				tv, err := GetDoStep(p)
				berror.CheckError(err)
				j.typeCache[t] = append(j.typeCache[t], &tv)
				for _, as := range tv.Do {
					n, err := GetStepName(as)
					berror.CheckError(err)
					j.stepCache[t][n] = &as
				}
			case TryStepType:
				tv, err := GetTryStep(p)
				berror.CheckError(err)
				j.typeCache[t] = append(j.typeCache[t], &tv)
				for _, as := range tv.Try {
					n, err := GetStepName(as)
					berror.CheckError(err)
					j.stepCache[t][n] = &as
				}
			case AggregateStepType:
				tv, err := GetAggregateStep(p)
				berror.CheckError(err)
				j.typeCache[t] = append(j.typeCache[t], &tv)
				for _, as := range tv.Aggregate {
					n, err := GetStepName(as)
					berror.CheckError(err)
					j.stepCache[t][n] = &as
				}
			default:
				log.Fatalf("unsupported type %+v", t)
			}
		}
	}
}

func (j *Job) AddStepByTypeAndName(t Type, name string, i interface{}) {
	j.buildCache()
	if j.stepCache[t] == nil {
		j.stepCache[t] = make(map[string]interface{})
	}
	_, exist := j.stepCache[t][name]
	if !exist {
		j.typeCache[t] = append(j.typeCache[t], i)
	}
	j.stepCache[t][name] = &i

}

func (j *Job) AddStepByType(t Type, i interface{}) {
	j.buildCache()
	j.typeCache[t] = append(j.typeCache[t], i)
}

func (j *Job) GetStepsByType(t Type) []interface{} {
	j.buildCache()
	return j.typeCache[t]
}

func (j *Job) GetStepByTypeAndName(t Type, name string) interface{} {
	j.buildCache()
	steps, ok := j.stepCache[t]
	if !ok {
		return nil
	}
	res, ok := steps[name]
	if !ok {
		return nil
	}
	return res
}

func GetJobsFromString(data string) Jobs {
	j := Jobs{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}
