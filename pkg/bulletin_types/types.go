package bulletin_types

import (
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Type int

const (
	StepType Type = iota
	TaskDecoratorType
	UnrecognizedType

	typeKey string = "type"
)

func GetType(s string) (Type, error) {
	res := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(s), &res)
	if err != nil {
		return UnrecognizedType, err
	}
	v, ok := res[typeKey]
	if !ok {
		return UnrecognizedType, errors.New(fmt.Sprintf("no type field"))
	}
	vs, ok := v.(string)
	if !ok {
		return UnrecognizedType, errors.New(fmt.Sprintf("type value is not string"))
	}
	switch vs {
	case "step":
		return StepType, nil
	case "task-decorator":
		return TaskDecoratorType, nil
	default:
		return UnrecognizedType, errors.New(fmt.Sprintf("Unrecognized type:%s", vs))
	}
}

func GetTemplatedString(s string) string {
	return "((" + s + "))"
}

func Replace(s interface{}, old, new string) (interface{}, error) {
	ss, err := yaml.Marshal(&s)
	if err != nil {
		return s, err
	}
	news := strings.Replace(string(ss), old, new, -1)
	var res interface{}
	err = yaml.Unmarshal([]byte(news), &res)
	if err != nil {
		return s, err
	}
	return res, nil
}

type FuncRef struct {
	Name   string            `yaml:"name"`
	Inputs map[string]string `yaml:"inputs,omitempty"`
}

type FuncDef struct {
	Name   string   `yaml:"name"`
	Type   string   `yaml:"type"`
	Inputs []string `yaml:"inputs,omitempty"`
}

func (f *FuncDef) Replace(fr FuncRef, i interface{}) (interface{}, error) {
	var err error
	var res = i
	for _, in := range f.Inputs {
		v, ok := fr.Inputs[in]
		if !ok {
			return i, errors.New(fmt.Sprintf("missing key %s in inputs", in))
		}
		res, err = Replace(i, GetTemplatedString(in), v)
		if err != nil {
			return i, err
		}
	}
	return res, nil
}

func (f *FuncDef) ReplaceString(fr FuncRef, i string) (string, error) {
	var res = i
	for _, in := range f.Inputs {
		v, ok := fr.Inputs[in]
		if !ok {
			return i, errors.New(fmt.Sprintf("missing key %s in inputs", in))
		}
		res = strings.Replace(i, GetTemplatedString(in), v, -1)
	}
	return res, nil
}
