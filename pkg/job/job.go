package job

import (
	berror "github.com/maplain/bulletin/pkg/error"
	yaml "gopkg.in/yaml.v2"
)

type Jobs struct {
	Jobs []Job `yaml:"jobs"`
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

type Job struct {
	Plan      []interface{} `yaml:"plan"`
	JobBase   `yaml:",inline"`
	StepHooks `yaml:",inline"`
}

func (j *Job) String() string {
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
