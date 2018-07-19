package resource

import (
	"log"

	yaml "gopkg.in/yaml.v2"
)

type ResourceTypes struct {
	ResourceTypes []ResourceType `yaml:"resource_types"`
}

type ResourceType struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	// optional fields
	Tags       []string    `yaml:"tags,omitempty"`
	Privileged bool        `yaml:"privileged,omitempty"`
	Params     interface{} `yaml:"params,omitempty"`
	Source     interface{} `yaml:"source,omitempty"`
}

type DockerResourceTypeSource struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag,omitempty"`
}

func (r *ResourceType) String() string {
	b, err := yaml.Marshal(*r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(b[:])
}

func GetResourceTypes(data string) ResourceTypes {
	r := ResourceTypes{}
	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return r
}
