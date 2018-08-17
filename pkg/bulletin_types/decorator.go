/*
Sniperkit-Bot
- Status: analyzed
*/

package bulletin_types

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	template "github.com/maplain/yamltemplate"
	yaml "gopkg.in/yaml.v2"

	berror "github.com/sniperkit/snk.fork.bulletin/pkg/error"
	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
	"github.com/sniperkit/snk.fork.bulletin/pkg/job"
	"github.com/sniperkit/snk.fork.bulletin/pkg/types"
)

const (
	decoratorsDir  = "decorators"
	decoratorsFile = "decorators.yml"
)

type StepDecoratorDefs struct {
	Decorators []StepDecoratorDef `yaml:"decorators"`
}

func (d *StepDecoratorDefs) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

func GetStepDecoratorDefsFromString(data string) StepDecoratorDefs {
	r := StepDecoratorDefs{}
	err := yaml.Unmarshal([]byte(data), &r)
	berror.CheckError(err)
	return r
}

type StepDecoratorDef struct {
	template.TemplateRef `yaml:",inline"`
	Decorate             []string `yaml:"decorate"`
}

func (d *StepDecoratorDef) GetJobTask(s string) (string, string) {
	res := strings.Split(s, "/")
	if len(res) == 2 {
		return res[0], res[1]
	} else if len(res) == 1 {
		return res[0], ""
	} else {
		log.Fatalf("error: invalid decorate target %s", s)
	}
	return res[0], res[1]
}

func (d *StepDecoratorDef) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

type Decorators struct {
	Decorators []Decorator `yaml:"decorators"`
	cache      map[string]Decorator
}

func (d *Decorators) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

func (d *Decorators) Populate(r template.TemplateRef) (Decorator, error) {
	if d.cache == nil {
		d.cache = make(map[string]Decorator)
		for _, dd := range d.Decorators {
			d.cache[dd.Name] = dd
		}
	}
	v, ok := d.cache[r.Name]
	if !ok {
		return Decorator{}, errors.New("no referenced Decorator definition")
	}
	return v.Populate(r)
}

type Decorator struct {
	template.TemplateDef `yaml:",inline"`
	Before               []interface{} `yaml:"before,omitempty"`
	After                []interface{} `yaml:"after,omitempty"`
	// step hook
	job.StepHooks `yaml:",inline"`
	// step modifier
	job.StepModifiers `yaml:",inline"`
}

func (o *Decorator) Populate(r template.TemplateRef) (Decorator, error) {
	var err error
	for i, _ := range o.Before {
		o.Before[i], err = o.Replace(r, o.Before[i])
		if err != nil {
			return *o, err
		}
	}
	for i, _ := range o.After {
		o.After[i], err = o.Replace(r, o.After[i])
		if err != nil {
			return *o, err
		}
	}
	o.OnSuccess, err = o.Replace(r, o.OnSuccess)
	if err != nil {
		return *o, err
	}
	o.OnFailure, err = o.Replace(r, o.OnFailure)
	if err != nil {
		return *o, err
	}
	o.OnAbort, err = o.Replace(r, o.OnAbort)
	if err != nil {
		return *o, err
	}
	o.Ensure, err = o.Replace(r, o.Ensure)
	if err != nil {
		return *o, err
	}
	o.Timeout, err = o.ReplaceOnString(r, o.Timeout)
	if err != nil {
		return *o, err
	}
	o.Attempts, err = o.ReplaceOnString(r, o.Attempts)
	if err != nil {
		return *o, err
	}
	for i, _ := range o.Tags {
		o.Tags[i], err = o.ReplaceOnString(r, o.Tags[i])
		if err != nil {
			return *o, err
		}
	}
	return *o, nil
}

func (d *Decorator) Decorate(s interface{}) (interface{}, error) {
	ss, err := yaml.Marshal(&s)
	if err != nil {
		return s, err
	}
	t, err := job.GetType(string(ss))
	if err != nil {
		return s, err
	}
	switch t {
	//	case job.GetStepType:
	//		fallthrough
	case job.PutStepType:
		task, err := job.GetPutStep(s)
		if err != nil {
			return s, err
		}
		if d.OnSuccess != nil {
			task.OnSuccess = d.OnSuccess
		}
		if d.OnFailure != nil {
			task.OnFailure = d.OnFailure
		}
		if d.OnAbort != nil {
			task.OnAbort = d.OnAbort
		}
		if d.Ensure != nil {
			task.Ensure = d.Ensure
		}
		return task, nil
	case job.TaskStepType:
		//	case job.AggregateSteptype:
		//		fallthrough
		//	case job.DoStepType:
		//		fallthrough
		//	case job.TryStepType:
		task, err := job.GetTaskStep(s)
		if err != nil {
			return s, err
		}
		if d.OnSuccess != nil {
			task.OnSuccess = d.OnSuccess
		}
		if d.OnFailure != nil {
			task.OnFailure = d.OnFailure
		}
		if d.OnAbort != nil {
			task.OnAbort = d.OnAbort
		}
		if d.Ensure != nil {
			task.Ensure = d.Ensure
		}
		return task, nil
	default:
		return s, errors.New(fmt.Sprintf("unsupported decorator type: %+v", t))
	}
}

func Decorate(s interface{}, descs ...Decorator) []interface{} {
	var res []interface{}
	l := len(descs)
	// added as interface{}
	for i := 0; i < l; i++ {
		res = append(res, descs[i].Before...)
	}
	for _, d := range descs {
		ds, err := d.Decorate(s)
		berror.CheckError(err)
		res = append(res, ds)
	}
	// added as interface{}
	for i := l - 1; i >= 0; i-- {
		res = append(res, descs[i].After...)
	}
	// if there is no decorator, add step
	if l == 0 {
		res = append(res, s)
	}
	return res
}

func (d *Decorator) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

//TODO
func (d Decorator) Equal(i interface{}) bool {
	switch v := i.(type) {
	case Decorator:
		return false
	case *Decorator:
		return false
	default:
		fmt.Printf("unsupported type %s\n", v)
	}
	return true
}

func GetDecoratorsFromString(data string) Decorators {
	r := Decorators{}

	err := yaml.Unmarshal([]byte(data), &r)
	berror.CheckError(err)
	return r
}

func GetLocalDecorators(target string) Decorators {
	targetDir := filepath.Join(target, decoratorsDir)
	ioutils.CreateDirIfNotExist(targetDir)
	targetFile := filepath.Join(targetDir, decoratorsFile)
	ioutils.CreateFileIfNotExist(targetFile)
	content := ioutils.ReadFile(targetFile)
	decs := GetDecoratorsFromString(content)
	resSet := types.NewSet()
	for _, d := range decs.Decorators {
		resSet.Add(d)
	}
	res := Decorators{}
	for _, d := range resSet.Get() {
		switch v := d.(type) {
		case Decorator:
			res.Decorators = append(res.Decorators, v)
		default:
			fmt.Printf("unsupported type %s\n", v)
		}
	}
	return res
}
