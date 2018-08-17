/*
Sniperkit-Bot
- Status: analyzed
*/

package bulletin_types

import (
	"errors"
	"fmt"
	"path/filepath"

	template "github.com/maplain/yamltemplate"
	yaml "gopkg.in/yaml.v2"

	berror "github.com/sniperkit/snk.fork.bulletin/pkg/error"
	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
	"github.com/sniperkit/snk.fork.bulletin/pkg/types"
)

const (
	stepsDir  = "steps"
	stepsFile = "steps.yml"
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

func (s *Steps) Populate(r template.TemplateRef) (Step, error) {
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
	template.TemplateDef `yaml:",inline"`
	Step                 interface{} `yaml:"step"`
}

func (s *Step) String() string {
	b, err := yaml.Marshal(*s)
	berror.CheckError(err)
	return string(b[:])
}

func (s *Step) Populate(r template.TemplateRef) (Step, error) {
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

func (s Step) Equal(i interface{}) bool {
	switch v := i.(type) {
	case Step:
		return false
	case *Step:
		return false
	default:
		fmt.Printf("unsupported type %s\n", v)
	}
	return true
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

func GetLocalSteps(target string) Steps {
	targetDir := filepath.Join(target, stepsDir)
	ioutils.CreateDirIfNotExist(targetDir)
	targetFile := filepath.Join(targetDir, stepsFile)
	ioutils.CreateFileIfNotExist(targetFile)
	content := ioutils.ReadFile(targetFile)
	steps := GetStepsFromString(content)
	resSet := types.NewSet()
	for _, d := range steps.Steps {
		resSet.Add(d)
	}
	res := Steps{}
	for _, d := range resSet.Get() {
		switch v := d.(type) {
		case Step:
			res.Steps = append(res.Steps, v)
		default:
			fmt.Printf("unsupported type %s\n", v)
		}
	}
	return res
}
