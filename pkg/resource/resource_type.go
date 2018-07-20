package resource

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/maplain/bulletin/pkg/ioutils"
	"github.com/maplain/bulletin/pkg/types"
	yaml "gopkg.in/yaml.v2"
)

const (
	DockerImageResourceType = "docker-image"
	resourceTypesDir        = "resource_types"
	resourceTypesFile       = "resource_types.yml"
)

type ResourceTypes struct {
	ResourceTypes []ResourceType `yaml:"resource_types"`
}

func (r *ResourceTypes) String() string {
	b, err := yaml.Marshal(*r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(b[:])
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

func (r *ResourceType) String() string {
	b, err := yaml.Marshal(*r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(b[:])
}

func (r *ResourceType) Equal(n ResourceType) bool {
	if r.Name != n.Name {
		return false
	}
	if r.Type != n.Type {
		return false
	}
	if r.Privileged != n.Privileged {
		return false
	}
	rparams, err := types.GetStringMap(r.Params)
	if err != nil {
		log.Fatalf("Unrecognized params: %+v\n", r.Params)
	}
	nparams, err := types.GetStringMap(n.Params)
	if err != nil {
		log.Fatalf("Unrecognized params: %+v\n", r.Params)
	}
	if !types.StringMapEqual(rparams, nparams) {
		return false
	}

	switch r.Type {
	case DockerImageResourceType:
		rsource, err := GetDockerResourceTypeSource(r.Source)
		if err != nil {
			log.Fatalf("Unrecognized source: %+v\n", r.Source)
		}
		nsource, err := GetDockerResourceTypeSource(n.Source)
		if err != nil {
			log.Fatalf("Unrecognized source: %+v\n", n.Source)
		}
		if !rsource.Equal(nsource) {
			return false
		}
	}
	return true
}

func GetResourceTypes(data string) ResourceTypes {
	r := ResourceTypes{}
	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return r
}

func SaveResourceTypesLocally(target string, res ResourceTypes) error {
	err := ioutil.WriteFile(filepath.Join(target, resourceTypesDir, resourceTypesFile), []byte(res.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalResourceTypes(target string) ResourceTypeSet {
	targetDir := filepath.Join(target, resourceTypesDir)
	ioutils.CreateDirIfNotExist(targetDir)
	targetFile := filepath.Join(targetDir, resourceTypesFile)
	ioutils.CreateFileIfNotExist(targetFile)
	content := ioutils.ReadFile(targetFile)
	resourceTypes := GetResourceTypes(content)
	res := ResourceTypeSet{}
	for _, rt := range resourceTypes.ResourceTypes {
		res.Add(rt)
	}
	return res
}

type ResourceTypeSet struct {
	rt []ResourceType
}

func (rs *ResourceTypeSet) Get() []ResourceType {
	return rs.rt
}

func (rs *ResourceTypeSet) Add(t ResourceType) {
	found := false
	for _, r := range rs.rt {
		if r.Equal(t) {
			found = true
			break
		}
	}
	if !found {
		rs.rt = append(rs.rt, t)
	}
}

type DockerResourceTypeSource struct {
	Repository             string                         `yaml:"repository"`
	Tag                    string                         `yaml:"tag,omitempty"`
	Username               string                         `yaml:"username,omitempty"`
	Password               string                         `yaml:"password,omitempty"`
	AwsAccessKeyId         string                         `yaml:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey     string                         `yaml:"aws_secret_access_key,omitempty"`
	AwsSessionToken        string                         `yaml:"aws_session_token,omitempty"`
	InsecureRegistries     []string                       `yaml:"insecure_registries,omitempty"`
	RegistryMirror         string                         `yaml:"registry_mirror,omitempty"`
	CaCerts                []DockerResourceTypeCACert     `yaml:"ca_certs,omitempty"`
	ClientCerts            []DockerResourceTypeClientCert `yaml:"client_certs,omitempty"`
	MaxConcurrentDownloads int                            `yaml:"max_concurrent_downloads,omitempty"`
	MaxConcurrentUploads   int                            `yaml:"max_concurrent_uploads,omitempty"`
}

func (d *DockerResourceTypeSource) Equal(n DockerResourceTypeSource) bool {
	if d.Repository != n.Repository {
		return false
	}
	if d.Tag != n.Tag {
		return false
	}
	if d.Username != n.Username {
		return false
	}
	if d.Password != n.Password {
		return false
	}
	if d.AwsAccessKeyId != n.AwsAccessKeyId {
		return false
	}
	if d.AwsSecretAccessKey != n.AwsSecretAccessKey {
		return false
	}
	if d.AwsSessionToken != n.AwsSessionToken {
		return false
	}
	if d.RegistryMirror != n.RegistryMirror {
		return false
	}
	if d.MaxConcurrentDownloads != n.MaxConcurrentDownloads {
		return false
	}
	if d.MaxConcurrentUploads != n.MaxConcurrentUploads {
		return false
	}
	if len(d.InsecureRegistries) != len(n.InsecureRegistries) {
		return false
	}
	if len(d.CaCerts) != len(n.CaCerts) {
		return false
	}
	if len(d.ClientCerts) != len(n.ClientCerts) {
		return false
	}
	sort.Strings(d.InsecureRegistries)
	sort.Strings(n.InsecureRegistries)
	for i, v := range d.InsecureRegistries {
		if v != n.InsecureRegistries[i] {
			return false
		}
	}
	equal := true
	for _, dcacert := range d.CaCerts {
		iequal := false
		for _, ncacert := range n.CaCerts {
			if dcacert.Equal(ncacert) {
				iequal = true
			}
		}
		if !iequal {
			equal = false
			break
		}
	}
	if !equal {
		return false
	}
	for _, dclientcert := range d.ClientCerts {
		iequal := false
		for _, nclientcert := range n.ClientCerts {
			if dclientcert.Equal(nclientcert) {
				iequal = true
			}
		}
		if !iequal {
			equal = false
			break
		}
	}
	if !equal {
		return false
	}
	return true
}

type DockerResourceTypeCACert struct {
	Domain string `yaml:"domain"`
	Cert   string `yaml:"cert"`
}

func (d *DockerResourceTypeCACert) Equal(n DockerResourceTypeCACert) bool {
	if d.Domain != n.Domain {
		return false
	}
	if d.Cert != n.Cert {
		return false
	}
	return true
}

type DockerResourceTypeClientCert struct {
	DockerResourceTypeCACert `yaml:",inline"`
	Key                      string `yaml:"key"`
}

func (d *DockerResourceTypeClientCert) Equal(n DockerResourceTypeClientCert) bool {
	if !d.DockerResourceTypeCACert.Equal(n.DockerResourceTypeCACert) {
		return false
	}
	if d.Key != n.Key {
		return false
	}
	return true
}

func GetDockerResourceTypeSource(i interface{}) (DockerResourceTypeSource, error) {
	res := DockerResourceTypeSource{}
	d, err := yaml.Marshal(i)
	if err != nil {
		return res, err
	}
	err = yaml.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
