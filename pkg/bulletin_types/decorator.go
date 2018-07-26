package bulletin_types

import (
	"errors"
	"fmt"

	berror "github.com/maplain/bulletin/pkg/error"
	"github.com/maplain/bulletin/pkg/job"
	yaml "gopkg.in/yaml.v2"
)

type Decorators struct {
	Decorators []Decorator `yaml:"decorators"`
	cache      map[string]Decorator
}

func (d *Decorators) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

func (d *Decorators) Populate(r FuncRef) (Decorator, error) {
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
	FuncDef `yaml:",inline"`
	Before  []interface{} `yaml:"before,omitempty"`
	After   []interface{} `yaml:"after,omitempty"`
	// step hook
	OnSuccess interface{} `yaml:"on_success,omitempty"`
	OnFailure interface{} `yaml:"on_failure,omitempty"`
	OnAbort   interface{} `yaml:"on_abort,omitempty"`
	Ensure    interface{} `yaml:"ensure,omitempty"`
	// step modifier
	Tags     []string `yaml:"tags,omitempty"`
	Timeout  string   `yaml:"timeout,omitempty"`
	Attempts string   `yaml:"attempts,omitempty"`
}

func (o *Decorator) Populate(r FuncRef) (Decorator, error) {
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
	o.Timeout, err = o.ReplaceString(r, o.Timeout)
	if err != nil {
		return *o, err
	}
	o.Attempts, err = o.ReplaceString(r, o.Attempts)
	if err != nil {
		return *o, err
	}
	for i, _ := range o.Tags {
		o.Tags[i], err = o.ReplaceString(r, o.Tags[i])
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
	case job.TaskStepType:
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
		return task, nil
	default:
		return s, errors.New("unsupported decorator type")
	}
}

func Decorate(s interface{}, descs ...Decorator) []interface{} {
	var res []interface{}
	l := len(descs)
	for i := 0; i < l; i++ {
		res = append(res, descs[i].Before...)
	}
	for _, d := range descs {
		ds, err := d.Decorate(s)
		if err != nil {
			fmt.Printf("err %s\n", err.Error())
		}
		res = append(res, ds)
	}
	for i := l - 1; i >= 0; i-- {
		res = append(res, descs[i].After...)
	}
	return res
}

func (d *Decorator) String() string {
	b, err := yaml.Marshal(*d)
	berror.CheckError(err)
	return string(b[:])
}

func GetDecoratorsFromString(data string) Decorators {
	r := Decorators{}

	err := yaml.Unmarshal([]byte(data), &r)
	berror.CheckError(err)
	return r
}
