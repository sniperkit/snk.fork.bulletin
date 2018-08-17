/*
Sniperkit-Bot
- Status: analyzed
*/

package bulletin_types

import (
	"errors"
	"fmt"

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
