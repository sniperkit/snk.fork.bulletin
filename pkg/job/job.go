package job

import (
	berror "github.com/maplain/bulletin/pkg/error"
	yaml "gopkg.in/yaml.v2"
)

var data = `
jobs:
- name: run-build-master
  serial: true
  plan:
    - aggregate:
      - get: git-pks-networking-master
        trigger: true

    - task: run-build
      file: git-pks-networking-master/ci/tasks/build.yml
      input_mapping:
        git-pks-networking: git-pks-networking-master

  
  on_failure:
    put: notify
    params:
      channel: pks-ci-bots
      attachments:
      - color: danger
        text: $BUILD_PIPELINE_NAME build failed. See results at <https://((ci_url))/teams/$BUILD_TEAM_NAME/pipelines/$BUILD_PIPELINE_NAME/jobs/$BUILD_JOB_NAME/builds/$BUILD_NAME>

- name: run-build-merge-request
  serial: true
  plan:
  - aggregate:
    - get: git-pks-networking-merge-request
      version: every
      trigger: true

  - put: git-pks-networking-merge-request
    params:
      repository: git-pks-networking-merge-request
      status: pending

  - task: run-build
    file: git-pks-networking-merge-request/ci/tasks/build.yml
    input_mapping:
      git-pks-networking: git-pks-networking-merge-request

    on_success:
      put: git-pks-networking-merge-request
      params:
        repository: git-pks-networking-merge-request
        status: success

    on_failure:
      put: git-pks-networking-merge-request
      params:
        repository: git-pks-networking-merge-request
        status: failed
`

type Jobs struct {
	Jobs []Job `yaml:"jobs"`
}

type Job struct {
	Name string        `yaml:"name"`
	Plan []interface{} `yaml:"plan"`
	// optional fields
	Serial               bool     `yaml:"serial,omitempty"`
	BuildLogsToRetain    int      `yaml:"build_logs_to_retain,omitempty"`
	SerialGroups         []string `yaml:"serial_groups,omitempty"`
	MaxInFlight          int      `yaml:"max_in_flight,omitempty"`
	Public               bool     `yaml:"public,omitempty"`
	DisableManualTrigger bool     `yaml:"disable_manual_trigger,omitempty"`
	Interruptible        bool     `yaml:"interruptible,omitempty"`
	StepHooks            `yaml:",inline"`
}

type StepHooks struct {
	OnSuccess interface{} `yaml:"on_success,omitempty"`
	OnFailure interface{} `yaml:"on_failure,omitempty"`
	OnAbort   interface{} `yaml:"on_abort,omitempty"`
	Ensure    interface{} `yaml:"ensure,omitempty"`
}

type Step struct {
	StepHooks `yaml:",inline"`
	Tags      []string `yaml:"tags,omitempty"`
	Timeout   string   `yaml:"timeout,omitempty"`
	Attempts  int      `yaml:"attempts,omitempty"`
}

type GetStep struct {
	Get string `yaml:"get"`
	// optional fields
	Resource string      `yaml:"resource,omitempty"`
	Version  string      `yaml:"version,omitempty"`
	Passed   []string    `yaml:"passed,omitempty"`
	Params   interface{} `yaml:"params,omitempty"`
	Trigger  bool        `yaml:"trigger,omitempty"`
}

type PubStep struct {
	Put string `yaml:"put"`
	// optional fields
	Resource  string      `yaml:"resource,omitempty"`
	Params    interface{} `yaml:"params,omitempty"`
	GetParams interface{} `yaml:"get_params,omitempty"`
}

type TaskStep struct {
	Task string `yaml:"task"`
	// optional fields
	Config        interface{}                 `yaml:"config,omitempty"`
	File          string                      `yaml:"config,omitempty"`
	Privileged    bool                        `yaml:"privileged,omitempty"`
	Params        map[interface{}]interface{} `yaml:"params,omitempty"`
	Images        string                      `yaml:"images,omitempty"`
	InputMapping  map[interface{}]interface{} `yaml:"input_mapping,omitempty"`
	OutputMapping map[interface{}]interface{} `yaml:"output_mapping,omitempty"`
}

type AggregateStep struct {
	Aggregate []Step `yaml:"aggregate"`
}

type DoStep struct {
	Do []Step `yaml:"do"`
}

type TryStep struct {
	Try []Step `yaml:"try"`
}

func (j *Job) String() string {
	b, err := yaml.Marshal(*j)
	berror.CheckError(err)
	return string(b[:])
}

func GetJobs(data string) Jobs {
	j := Jobs{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}
