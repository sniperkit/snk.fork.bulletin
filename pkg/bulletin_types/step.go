package bulletin_types

import (
	"errors"
	"fmt"

	berror "github.com/maplain/bulletin/pkg/error"
	yaml "gopkg.in/yaml.v2"
)

type Steps struct {
	Steps []Step          `yaml:"steps"`
	cache map[string]Step `yaml:",omitempty"`
}

func (s *Steps) String() string {
	b, err := yaml.Marshal(*s)
	berror.CheckError(err)
	return string(b[:])
}

func (s *Steps) Populate(r FuncRef) (Step, error) {
	if s.cache == nil {
		s.cache = make(map[string]Step)
		for _, st := range s.Steps {
			s.cache[st.Name] = st
		}
	}
	v, ok := s.cache[r.Name]
	if !ok {
		return Step{}, errors.New("no referenced Step definition")
	}
	return v.Populate(r)
}

type Step struct {
	FuncDef `yaml:",inline"`
	Step    interface{} `yaml:"step"`
}

func (s *Step) String() string {
	b, err := yaml.Marshal(*s)
	berror.CheckError(err)
	return string(b[:])
}

func (s *Step) Populate(r FuncRef) (Step, error) {
	var err error
	s.Step, err = s.Replace(r, s.Step)
	if err != nil {
		return *s, err
	}
	return *s, nil
}

func (s *Step) GetStep() interface{} {
	return s.Step
}

func GetStepsFromString(data string) Steps {
	r := Steps{}
	err := yaml.Unmarshal([]byte(data), &r)
	berror.CheckError(err)
	return r
}

func getStepFromString(data string) Step {
	r := Step{}
	err := yaml.Unmarshal([]byte(data), &r)
	berror.CheckError(err)
	return r
}

func GetStep(s string) (Step, error) {
	t, err := GetType(s)
	if err != nil {
		return Step{}, err
	}
	switch t {
	case StepType:
		return getStepFromString(s), nil
	default:
		return Step{}, errors.New(fmt.Sprintf("not a Step"))
	}
}
