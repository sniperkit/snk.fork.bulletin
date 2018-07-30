package job

import (
	"errors"

	berror "github.com/maplain/bulletin/pkg/error"
	"github.com/maplain/bulletin/pkg/types"
	yaml "gopkg.in/yaml.v2"
)

type Type int

const (
	GetStepType Type = iota
	PutStepType
	TaskStepType
	AggregateStepType
	DoStepType
	TryStepType
	UnrecognizedType

	TypeNotSupportedError types.InternalError = "specified type is not supported"
)

type Step struct {
	StepHooks     `yaml:",inline"`
	StepModifiers `yaml:",inline"`
}

type StepHooks struct {
	OnSuccess interface{} `yaml:"on_success,omitempty"`
	OnFailure interface{} `yaml:"on_failure,omitempty"`
	OnAbort   interface{} `yaml:"on_abort,omitempty"`
	Ensure    interface{} `yaml:"ensure,omitempty"`
}

type StepModifiers struct {
	Tags     []string `yaml:"tags,omitempty"`
	Timeout  string   `yaml:"timeout,omitempty"`
	Attempts string   `yaml:"attempts,omitempty"`
}

type GetStep struct {
	Step `yaml:",inline"`
	Get  string `yaml:"get"`
	// optional fields
	Resource string      `yaml:"resource,omitempty"`
	Version  string      `yaml:"version,omitempty"`
	Passed   []string    `yaml:"passed,omitempty"`
	Params   interface{} `yaml:"params,omitempty"`
	Trigger  bool        `yaml:"trigger,omitempty"`
}

func (s *GetStep) String() string {
	b, err := yaml.Marshal(*s)
	berror.CheckError(err)
	return string(b[:])
}

type PutStep struct {
	Step `yaml:",inline"`
	Put  string `yaml:"put"`
	// optional fields
	Resource  string      `yaml:"resource,omitempty"`
	Params    interface{} `yaml:"params,omitempty"`
	GetParams interface{} `yaml:"get_params,omitempty"`
}

type TaskStep struct {
	Step `yaml:",inline"`
	Task string `yaml:"task"`
	// optional fields
	Config        interface{}                 `yaml:"config,omitempty"`
	File          string                      `yaml:"file,omitempty"`
	Privileged    bool                        `yaml:"privileged,omitempty"`
	Params        map[interface{}]interface{} `yaml:"params,omitempty"`
	Images        string                      `yaml:"images,omitempty"`
	InputMapping  map[interface{}]interface{} `yaml:"input_mapping,omitempty"`
	OutputMapping map[interface{}]interface{} `yaml:"output_mapping,omitempty"`
}

type AggregateStep struct {
	Step      `yaml:",inline"`
	Aggregate []interface{} `yaml:"aggregate"`
}

func (s *AggregateStep) String() string {
	b, err := yaml.Marshal(*s)
	berror.CheckError(err)
	return string(b[:])
}

type DoStep struct {
	Step `yaml:",inline"`
	Do   []interface{} `yaml:"do"`
}

type TryStep struct {
	Step `yaml:",inline"`
	Try  []interface{} `yaml:"try"`
}

func GetType(s string) (Type, error) {
	res := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(s), &res)
	if err != nil {
		return UnrecognizedType, err
	}
	if _, ok := res["get"]; ok {
		return GetStepType, nil
	} else if _, ok := res["put"]; ok {
		return PutStepType, nil
	} else if _, ok := res["task"]; ok {
		return TaskStepType, nil
	} else if _, ok := res["aggregate"]; ok {
		return AggregateStepType, nil
	} else if _, ok := res["do"]; ok {
		return DoStepType, nil
	} else if _, ok := res["try"]; ok {
		return TryStepType, nil
	} else {
		return UnrecognizedType, nil
	}
	return UnrecognizedType, nil
}

func GetTaskStep(s interface{}) (TaskStep, error) {
	d, err := yaml.Marshal(&s)
	if err != nil {
		return TaskStep{}, err
	}
	t, err := GetType(string(d))
	if err != nil {
		return TaskStep{}, err
	}
	switch t {
	case TaskStepType:
		return getTaskStepFromString(string(d)), nil
	default:
		return TaskStep{}, errors.New("not a task Step")
	}
}

func getTaskStepFromString(data string) TaskStep {
	j := TaskStep{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}

func GetPutStep(s interface{}) (PutStep, error) {
	d, err := yaml.Marshal(&s)
	if err != nil {
		return PutStep{}, err
	}
	t, err := GetType(string(d))
	if err != nil {
		return PutStep{}, err
	}
	switch t {
	case PutStepType:
		return getPutStepFromString(string(d)), nil
	default:
		return PutStep{}, errors.New("not a put Step")
	}
}

func getPutStepFromString(data string) PutStep {
	j := PutStep{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}

func GetGetStep(s interface{}) (GetStep, error) {
	d, err := yaml.Marshal(&s)
	if err != nil {
		return GetStep{}, err
	}
	t, err := GetType(string(d))
	if err != nil {
		return GetStep{}, err
	}
	switch t {
	case GetStepType:
		return getGetStepFromString(string(d)), nil
	default:
		return GetStep{}, errors.New("not a get Step")
	}
}

func getGetStepFromString(data string) GetStep {
	j := GetStep{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}

func GetAggregateStep(s interface{}) (AggregateStep, error) {
	d, err := yaml.Marshal(&s)
	if err != nil {
		return AggregateStep{}, err
	}
	t, err := GetType(string(d))
	if err != nil {
		return AggregateStep{}, err
	}
	switch t {
	case AggregateStepType:
		return getAggregateStepFromString(string(d)), nil
	default:
		return AggregateStep{}, errors.New("not a aggregate Step")
	}
}

func getAggregateStepFromString(data string) AggregateStep {
	j := AggregateStep{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}

func GetDoStep(s interface{}) (DoStep, error) {
	d, err := yaml.Marshal(&s)
	if err != nil {
		return DoStep{}, err
	}
	t, err := GetType(string(d))
	if err != nil {
		return DoStep{}, err
	}
	switch t {
	case DoStepType:
		return getDoStepFromString(string(d)), nil
	default:
		return DoStep{}, errors.New("not a do Step")
	}
}

func getDoStepFromString(data string) DoStep {
	j := DoStep{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}

func GetTryStep(s interface{}) (TryStep, error) {
	d, err := yaml.Marshal(&s)
	if err != nil {
		return TryStep{}, err
	}
	t, err := GetType(string(d))
	if err != nil {
		return TryStep{}, err
	}
	switch t {
	case TryStepType:
		return getTryStepFromString(string(d)), nil
	default:
		return TryStep{}, errors.New("not a try Step")
	}
}

func getTryStepFromString(data string) TryStep {
	j := TryStep{}
	err := yaml.Unmarshal([]byte(data), &j)
	berror.CheckError(err)
	return j
}

func GetStepName(i interface{}) (string, error) {
	s, err := yaml.Marshal(&i)
	if err != nil {
		return "", err
	}
	t, err := GetType(string(s))
	if err != nil {
		return "", err
	}
	switch t {
	case PutStepType:
		tv, err := GetPutStep(i)
		berror.CheckError(err)
		return tv.Put, nil
	case GetStepType:
		tv, err := GetGetStep(i)
		berror.CheckError(err)
		return tv.Get, nil
	case TaskStepType:
		tv, err := GetTaskStep(i)
		berror.CheckError(err)
		return tv.Task, nil
	default:
		return "", TypeNotSupportedError
	}
}
