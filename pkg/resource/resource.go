package resource

import (
	"github.com/maplain/bulletin/pkg/error"
	yaml "gopkg.in/yaml.v2"
)

type Resources struct {
	Resources []Resource `yaml:"resources"`
}

type Resource struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	// optional fields
	Tags         []string    `yaml:"tags,omitempty"`
	CheckEvery   string      `yaml:"check_every,omitempty"`
	Source       interface{} `yaml:"source,omitempty"`
	WebhookToken string      `yaml:"webhook_token,omitempty"`
}

func (r *Resource) String() string {
	b, err := yaml.Marshal(*r)
	error.CheckError(err)
	return string(b[:])
}

func GetResources(data string) Resources {
	r := Resources{}

	err := yaml.Unmarshal([]byte(data), &r)
	error.CheckError(err)
	return r
}
