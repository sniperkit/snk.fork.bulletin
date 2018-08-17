/*
Sniperkit-Bot
- Status: analyzed
*/

package resource

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	berror "github.com/sniperkit/snk.fork.bulletin/pkg/error"
	"github.com/sniperkit/snk.fork.bulletin/pkg/ioutils"
	"github.com/sniperkit/snk.fork.bulletin/pkg/types"
)

const (
	GCSResourceType               = "gcs-resource"
	GithubReleaseResourceType     = "github-release"
	BoshIOStemcellResourceType    = "bosh-io-stemcell"
	GitResourceType               = "git"
	MergeRequestResourceType      = "merge-request"
	SlackNotificationResourceType = "slack-notification"
	PoolResourceType              = "pool"
	SemverResourceType            = "semver"

	SemverResourceDriverGit = "git"
	SemverResourceDriverS3  = "s3"
	SemverResourceDriverGCS = "gcs"
	//	SemverResourceDriverSwift = "swift"

	resourcesDir  = "resources"
	resourcesFile = "resources.yml"
)

type Resources struct {
	Resources []Resource `yaml:"resources"`
}

func (r *Resources) String() string {
	b, err := yaml.Marshal(*r)
	berror.CheckError(err)
	return string(b[:])
}

func (r *Resources) Map() map[string]Resource {
	res := make(map[string]Resource)
	for _, r := range r.Resources {
		res[r.Name] = r
	}
	return res
}

func NewResourcesFromMap(i map[string]Resource) Resources {
	var rs []Resource
	for _, v := range i {
		rs = append(rs, v)
	}
	return Resources{rs}
}

func (r *Resources) UpdateWith(n Resources) Resources {
	rm := r.Map()
	nm := n.Map()
	var res []Resource
	for n, rmr := range rm {
		if nr, ok := nm[n]; ok {
			res = append(res, rmr.UpdateWith(nr))
		} else {
			res = append(res, rmr)
		}
	}
	return Resources{res}
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
	berror.CheckError(err)
	return string(b[:])
}

// simple replace
//TODO: based on fields
func (r *Resource) UpdateWith(n Resource) Resource {
	res := *r
	if !r.Equal(n) {
		if n.Type != "" {
			if r.Type != n.Type {
				log.Fatalf("Can not update resource:\n%+v\nwith resource:\n%+v\n", r.String(), n.String())
			}
		}
		if n.CheckEvery != "" {
			res.CheckEvery = n.CheckEvery
		}
		if n.WebhookToken != "" {
			res.WebhookToken = n.WebhookToken
		}
		if len(n.Tags) != 0 {
			res.Tags = n.Tags
		}
		switch r.Type {
		case GCSResourceType:
			if !GCSResourceEqual(r.Source, n.Source) {
				res.Source = UpdateGCSResource(r.Source, n.Source)
			}
		case GithubReleaseResourceType:
			if !GithubReleaseResourceEqual(r.Source, n.Source) {
				res.Source = UpdateGithubReleaseResource(r.Source, n.Source)
			}
		case BoshIOStemcellResourceType:
			if !BoshIOStemcellResourceEqual(r.Source, n.Source) {
				res.Source = UpdateBoshIOStemcellResource(r.Source, n.Source)
			}
		case GitResourceType:
			if !GitResourceEqual(r.Source, n.Source) {
				res.Source = UpdateGitResource(r.Source, n.Source)
			}
		case MergeRequestResourceType:
			if !MergeRequestResourceEqual(r.Source, n.Source) {
				res.Source = UpdateMergeRequestResource(r.Source, n.Source)
			}
		case SlackNotificationResourceType:
			if !SlackNotificationResourceEqual(r.Source, n.Source) {
				res.Source = UpdateSlackNotificationResource(r.Source, n.Source)
			}
		case PoolResourceType:
			if !PoolResourceEqual(r.Source, n.Source) {
				res.Source = UpdatePoolResource(r.Source, n.Source)
			}
		case SemverResourceType:
			if !SemverResourceEqual(r.Source, n.Source) {
				res.Source = UpdateSemverResource(r.Source, n.Source)
			}
		default:
			fmt.Printf("Unrecognized resource type: %s", r.Type)
		}
	}
	return res
}

func (r *Resource) Equal(n Resource) bool {
	if r.Name != n.Name {
		return false
	}
	if r.Type != n.Type {
		return false
	}
	if !types.StringSliceEqual(r.Tags, n.Tags) {
		return false
	}
	if r.CheckEvery != n.CheckEvery {
		return false
	}
	if r.WebhookToken != n.WebhookToken {
		return false
	}
	switch r.Type {
	case GCSResourceType:
		if !GCSResourceEqual(r.Source, n.Source) {
			return false
		}
	case GithubReleaseResourceType:
		if !GithubReleaseResourceEqual(r.Source, n.Source) {
			return false
		}
	case BoshIOStemcellResourceType:
		if !BoshIOStemcellResourceEqual(r.Source, n.Source) {
			return false
		}
	case GitResourceType:
		if !GitResourceEqual(r.Source, n.Source) {
			return false
		}
	case MergeRequestResourceType:
		if !MergeRequestResourceEqual(r.Source, n.Source) {
			return false
		}
	case SlackNotificationResourceType:
		if !SlackNotificationResourceEqual(r.Source, n.Source) {
			return false
		}
	case PoolResourceType:
		if !PoolResourceEqual(r.Source, n.Source) {
			return false
		}
	case SemverResourceType:
		if !SemverResourceEqual(r.Source, n.Source) {
			return false
		}
	default:
		fmt.Printf("Unrecognized resource type: %s", r.Type)
	}
	return true
}

func SaveResourcesLocally(target string, res Resources) error {
	err := ioutil.WriteFile(filepath.Join(target, resourcesDir, resourcesFile), []byte(res.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalResources(target string) ResourceSet {
	targetDir := filepath.Join(target, resourcesDir)
	ioutils.CreateDirIfNotExist(targetDir)
	targetFile := filepath.Join(targetDir, resourcesFile)
	ioutils.CreateFileIfNotExist(targetFile)
	content := ioutils.ReadFile(targetFile)
	resources := GetResourcesFromString(content)
	res := ResourceSet{}
	for _, rt := range resources.Resources {
		res.Add(rt)
	}
	return res
}

func GetResourcesFromString(data string) Resources {
	r := Resources{}

	err := yaml.Unmarshal([]byte(data), &r)
	berror.CheckError(err)
	return r
}

func GetResourcesFromFile(filename string) Resources {
	return GetResourcesFromString(ioutils.ReadFile(filename))
}

type ResourceSet struct {
	rt []Resource
}

func (rs *ResourceSet) Get() []Resource {
	return rs.rt
}

func (rs *ResourceSet) Add(t Resource) {
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

func GCSResourceEqual(a, b interface{}) bool {
	rsource, err := GetGCSResource(a)
	if err != nil {
		log.Fatalf("Unrecognized GCSResource: %+v\n", a)
	}
	nsource, err := GetGCSResource(b)
	if err != nil {
		log.Fatalf("Unrecognized GCSResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func GithubReleaseResourceEqual(a, b interface{}) bool {
	rsource, err := GetGithubReleaseResource(a)
	if err != nil {
		log.Fatalf("Unrecognized GCSResource: %+v\n", a)
	}
	nsource, err := GetGithubReleaseResource(b)
	if err != nil {
		log.Fatalf("Unrecognized GCSResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func BoshIOStemcellResourceEqual(a, b interface{}) bool {
	rsource, err := GetBoshIOStemcellResource(a)
	if err != nil {
		log.Fatalf("Unrecognized BoshIOStemcellResource: %+v\n", a)
	}
	nsource, err := GetBoshIOStemcellResource(b)
	if err != nil {
		log.Fatalf("Unrecognized BoshIOStemcellResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func GitResourceEqual(a, b interface{}) bool {
	rsource, err := GetGitResource(a)
	if err != nil {
		log.Fatalf("Unrecognized GitResource: %+v\n", a)
	}
	nsource, err := GetGitResource(b)
	if err != nil {
		log.Fatalf("Unrecognized GitResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func MergeRequestResourceEqual(a, b interface{}) bool {
	rsource, err := GetMergeRequestResource(a)
	if err != nil {
		log.Fatalf("Unrecognized MergeRequestResource: %+v\n", a)
	}
	nsource, err := GetMergeRequestResource(b)
	if err != nil {
		log.Fatalf("Unrecognized MergeRequestResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func SlackNotificationResourceEqual(a, b interface{}) bool {
	rsource, err := GetSlackNotificationResource(a)
	if err != nil {
		log.Fatalf("Unrecognized SlackNotificationResource: %+v\n", a)
	}
	nsource, err := GetSlackNotificationResource(b)
	if err != nil {
		log.Fatalf("Unrecognized SlackNotificationResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func PoolResourceEqual(a, b interface{}) bool {
	rsource, err := GetPoolResource(a)
	if err != nil {
		log.Fatalf("Unrecognized PoolResource: %+v\n", a)
	}
	nsource, err := GetPoolResource(b)
	if err != nil {
		log.Fatalf("Unrecognized PoolResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func SemverResourceEqual(a, b interface{}) bool {
	rsourceBase, err := GetSemverResourceBase(a)
	if err != nil {
		log.Fatalf("Unrecognized SemverResource: %+v\n", a)
	}
	nsourceBase, err := GetSemverResourceBase(b)
	if err != nil {
		log.Fatalf("Unrecognized SemverResource: %+v\n", b)
	}
	if rsourceBase != nsourceBase {
		return false
	}
	//TODO: add support for other drivers
	switch rsourceBase.Driver {
	//	case SemverResourceDriverGit:
	//		if !SemverGitResourceEqual(a, b) {
	//			return false
	//		}
	//	case SemverResourceDriverS3:
	//		if !SemverS3ResourceEqual(a, b) {
	//			return false
	//		}
	case SemverResourceDriverGCS:
		if !SemverGCSResourceEqual(a, b) {
			return false
		}
	default:
		log.Fatalf("Unsupported semver resource driver %s\n", rsourceBase.Driver)
	}
	return true
}

func SemverGCSResourceEqual(a, b interface{}) bool {
	rsource, err := GetSemverGCSResource(a)
	if err != nil {
		log.Fatalf("Unrecognized SemverGCSResource: %+v\n", a)
	}
	nsource, err := GetSemverGCSResource(b)
	if err != nil {
		log.Fatalf("Unrecognized SemverGCSResource: %+v\n", b)
	}
	if !rsource.Equal(nsource) {
		return false
	}
	return true
}

func GetGCSResource(i interface{}) (GCSResource, error) {
	res := GCSResource{}
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

type GCSResource struct {
	Bucket  string `yaml:"bucket`
	JsonKey string `yaml:"json_key"`
	// optional field
	Regexp        string `yaml:"regexp,omitempty"`
	VersionedFile string `yaml:"versioned_file,omitempty"`
}

func (g *GCSResource) Equal(n GCSResource) bool {
	return *g == n
}

func GetGithubReleaseResource(i interface{}) (GithubReleaseResource, error) {
	res := GithubReleaseResource{}
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

type GithubReleaseResource struct {
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repository"`
	// optional field
	AccessToken      string `yaml:"access_token,omitempty"`
	GithubApiUrl     string `yaml:"github_api_url,omitempty"`
	GithubUploadsUrl string `yaml:"github_uploads_url,omitempty"`
	Insecure         bool   `yaml:"insecure,omitempty"`
	Release          bool   `yaml:"release,omitempty"`
	PreRelease       bool   `yaml:"pre_release,omitempty"`
	Drafts           bool   `yaml:"drafts,omitempty"`
	TagFilter        string `yaml:"tag_filter,omitempty"`
}

func (g *GithubReleaseResource) Equal(n GithubReleaseResource) bool {
	return *g == n
}

func GetBoshIOStemcellResource(i interface{}) (BoshIOStemcellResource, error) {
	res := BoshIOStemcellResource{}
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

type BoshIOStemcellResource struct {
	Repository string `yaml:"repository"`
}

func (b *BoshIOStemcellResource) Equal(n BoshIOStemcellResource) bool {
	return *b == n
}

func GetGitResource(i interface{}) (GitResource, error) {
	res := GitResource{}
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

type GitResource struct {
	Uri string `yaml:"uri"`
	// optional field
	Branch                   string                 `yaml:"branch,omitempty"`
	PrivateKey               string                 `yaml:"private_key,omitempty"`
	Username                 string                 `yaml:"username,omitempty"`
	Password                 string                 `yaml:"password,omitempty"`
	Paths                    []string               `yaml:"paths,omitempty"`
	IgnorePaths              []string               `yaml:"ignore_paths,omitempty"`
	SkipSslVerification      bool                   `yaml:"skip_ssl_verification,omitempty"`
	TagFilter                string                 `yaml:"tag_filter,omitempty"`
	GitConfig                GitResourceGitConfigs  `yaml:"git_config,omitempty"`
	DisableCiSkip            bool                   `yaml:"disable_ci_skip,omitempty"`
	CommitVerificationKeys   []string               `yaml:"commit_verification_keys,omitempty"`
	CommitVerificationKeyIds []string               `yaml:"commit_verification_key_ids,omitempty"`
	GpgKeyserver             string                 `yaml:"gpg_keyserver,omitempty"`
	GitCryptKey              string                 `yaml:"git_crypt_key,omitempty"`
	HttpsTunnel              GitResourceHttpsTunnel `yaml:"https_tunnel,omitempty"`
}

type GitResourceHttpsTunnel struct {
	ProxyHost     string `yaml:"proxy_host"`
	ProxyPort     string `yaml:"proxy_port"`
	ProxyUser     string `yaml:"proxy_user,omitempty"`
	ProxyPassword string `yaml:"proxy_password,omitempty"`
}

type GitResourceGitConfigs struct {
	configs []map[string]string
}

func (g *GitResourceGitConfigs) Len() int {
	return len(g.configs)
}

func (g *GitResourceGitConfigs) Equal(n GitResourceGitConfigs) bool {
	cg := make(map[string]string)
	cn := make(map[string]string)
	for _, gv := range g.configs {
		cg[gv["name"]] = cg[gv["value"]]
	}
	for _, nv := range n.configs {
		cn[nv["name"]] = cn[nv["value"]]
	}
	return types.StringMapEqual(cg, cn)
}

func (g *GitResource) Equal(n GitResource) bool {
	var c types.Comparator
	if !c.Strings(g.Uri, n.Uri).
		Strings(g.Branch, n.Branch).
		Strings(g.PrivateKey, n.PrivateKey).
		Strings(g.Username, n.Username).
		Strings(g.Password, n.Password).
		Strings(g.TagFilter, n.TagFilter).
		Strings(g.GpgKeyserver, n.GpgKeyserver).
		Strings(g.GitCryptKey, n.GitCryptKey).
		StringSlice(g.Paths, n.Paths).
		StringSlice(g.IgnorePaths, n.IgnorePaths).
		StringSlice(g.CommitVerificationKeys, n.CommitVerificationKeys).
		StringSlice(g.CommitVerificationKeyIds, n.CommitVerificationKeyIds).
		Bool(g.SkipSslVerification, n.SkipSslVerification).
		Bool(g.DisableCiSkip, n.DisableCiSkip).
		Equal() {
		return false
	}
	if !g.GitConfig.Equal(n.GitConfig) {
		return false
	}
	if g.HttpsTunnel != n.HttpsTunnel {
		return false
	}
	return true
}

func GetMergeRequestResource(i interface{}) (MergeRequestResource, error) {
	res := MergeRequestResource{}
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

type MergeRequestResource struct {
	Uri          string `yaml:"uri"`
	PrivateToken string `yaml:"private_token"`
	// optional field
	PrivateKey          string `yaml:"private_key,omitempty"`
	Username            string `yaml:"username,omitempty`
	Password            string `yaml:"password,omitempty`
	NoSSL               bool   `yaml:"no_ssl,omitempty"`
	SkipSslVerification bool   `yaml:"skip_ssl_verification,omitempty"`
}

func (m *MergeRequestResource) Equal(n MergeRequestResource) bool {
	return *m == n
}

func GetSlackNotificationResource(i interface{}) (SlackNotificationResource, error) {
	res := SlackNotificationResource{}
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

type SlackNotificationResource struct {
	URL string `yaml:"url"`
}

func (s *SlackNotificationResource) Equal(n SlackNotificationResource) bool {
	return *s == n
}

func GetPoolResource(i interface{}) (PoolResource, error) {
	res := PoolResource{}
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

type PoolResource struct {
	Uri    string `yaml:"uri"`
	Branch string `yaml:"branch"`
	Pool   string `yaml:"pool"`
	// optional field
	PrivateKey string `yaml:"private_key,omitempty"`
	Username   string `yaml:"username,omitempty`
	Password   string `yaml:"password,omitempty`
	RetryDelay string `yaml:"retry_delay,omitempty"`
}

func (s *PoolResource) Equal(n PoolResource) bool {
	return *s == n
}

func GetSemverResourceBase(i interface{}) (SemverResourceBase, error) {
	res := SemverResourceBase{}
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

type SemverResourceBase struct {
	InitialVersion string `yaml:"initial_version,omitempty"`
	Driver         string `yaml:"driver,omitempty"`
}

func (s *SemverResourceBase) Equal(n SemverResourceBase) bool {
	return *s == n
}

func GetSemverGCSResource(i interface{}) (SemverGCSResource, error) {
	res := SemverGCSResource{}
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

type SemverGCSResource struct {
	SemverResourceBase `yaml:",inline"`
	Bucket             string `yaml:"Bucket,omitempty"`
	Key                string `yaml:"Key,omitempty"`
	JsonKey            string `yaml:"JsonKey,omitempty"`
}

func (s *SemverGCSResource) Equal(n SemverGCSResource) bool {
	return *s == n
}

func UpdateGCSResource(a, b interface{}) interface{} {
	rsource, err := GetGCSResource(a)
	if err != nil {
		log.Fatalf("Unrecognized MergeRequestResource: %+v\n", a)
	}
	nsource, err := GetGCSResource(b)
	if err != nil {
		log.Fatalf("Unrecognized MergeRequestResource: %+v\n", b)
	}
	if nsource.Bucket != "" {
		rsource.Bucket = nsource.Bucket
	}
	if nsource.JsonKey != "" {
		rsource.JsonKey = nsource.JsonKey
	}
	if nsource.Regexp != "" {
		rsource.Regexp = nsource.Regexp
	}
	if nsource.VersionedFile != "" {
		rsource.VersionedFile = nsource.VersionedFile
	}
	return rsource
}

func UpdateGithubReleaseResource(a, b interface{}) interface{} {
	rsource, err := GetGithubReleaseResource(a)
	if err != nil {
		log.Fatalf("Unrecognized GithubReleaseResource: %+v\n", a)
	}
	nsource, err := GetGithubReleaseResource(b)
	if err != nil {
		log.Fatalf("Unrecognized GithubReleaseResource: %+v\n", b)
	}
	if nsource.Owner != "" {
		rsource.Owner = nsource.Owner
	}
	if nsource.Repository != "" {
		rsource.Repository = nsource.Repository
	}
	if nsource.AccessToken != "" {
		rsource.AccessToken = nsource.AccessToken
	}
	if nsource.GithubApiUrl != "" {
		rsource.GithubApiUrl = nsource.GithubApiUrl
	}
	if nsource.TagFilter != "" {
		rsource.TagFilter = nsource.TagFilter
	}

	// reset all boolean values
	rsource.Insecure = nsource.Insecure
	rsource.Release = nsource.Release
	rsource.PreRelease = nsource.PreRelease
	rsource.Drafts = nsource.Drafts
	return rsource
}

func UpdateBoshIOStemcellResource(a, b interface{}) interface{} {
	rsource, err := GetBoshIOStemcellResource(a)
	if err != nil {
		log.Fatalf("Unrecognized BoshIOStemcellResource: %+v\n", a)
	}
	nsource, err := GetBoshIOStemcellResource(b)
	if err != nil {
		log.Fatalf("Unrecognized BoshIOStemcellResource: %+v\n", b)
	}
	if nsource.Repository != "" {
		rsource.Repository = nsource.Repository
	}
	return rsource
}

func UpdateGitResource(a, b interface{}) interface{} {
	rsource, err := GetGitResource(a)
	if err != nil {
		log.Fatalf("Unrecognized GitResource: %+v\n", a)
	}
	nsource, err := GetGitResource(b)
	if err != nil {
		log.Fatalf("Unrecognized GitResource: %+v\n", b)
	}
	if nsource.Uri != "" {
		rsource.Uri = nsource.Uri
	}
	if nsource.Branch != "" {
		rsource.Branch = nsource.Branch
	}
	if nsource.PrivateKey != "" {
		rsource.PrivateKey = nsource.PrivateKey
	}
	if nsource.Username != "" {
		rsource.Username = nsource.Username
	}
	if nsource.Password != "" {
		rsource.Password = nsource.Password
	}
	if nsource.TagFilter != "" {
		rsource.TagFilter = nsource.TagFilter
	}
	if nsource.GpgKeyserver != "" {
		rsource.GpgKeyserver = nsource.GpgKeyserver
	}
	if nsource.GitCryptKey != "" {
		rsource.GitCryptKey = nsource.GitCryptKey
	}
	rsource.SkipSslVerification = nsource.SkipSslVerification
	rsource.DisableCiSkip = nsource.DisableCiSkip
	if len(nsource.Paths) != 0 {
		rsource.Paths = nsource.Paths
	}
	if len(nsource.IgnorePaths) != 0 {
		rsource.IgnorePaths = nsource.IgnorePaths
	}
	if len(nsource.CommitVerificationKeys) != 0 {
		rsource.CommitVerificationKeys = nsource.CommitVerificationKeys
	}
	if len(nsource.CommitVerificationKeyIds) != 0 {
		rsource.CommitVerificationKeyIds = nsource.CommitVerificationKeyIds
	}
	// no partial update on git config
	if nsource.GitConfig.Len() != 0 {
		rsource.GitConfig = nsource.GitConfig
	}
	if nsource.HttpsTunnel.ProxyHost != "" {
		rsource.HttpsTunnel.ProxyHost = nsource.HttpsTunnel.ProxyHost
	}
	if nsource.HttpsTunnel.ProxyPort != "" {
		rsource.HttpsTunnel.ProxyPort = nsource.HttpsTunnel.ProxyPort
	}
	if nsource.HttpsTunnel.ProxyUser != "" {
		rsource.HttpsTunnel.ProxyUser = nsource.HttpsTunnel.ProxyUser
	}
	if nsource.HttpsTunnel.ProxyPassword != "" {
		rsource.HttpsTunnel.ProxyPassword = nsource.HttpsTunnel.ProxyPassword
	}
	return rsource
}

func UpdateMergeRequestResource(a, b interface{}) interface{} {
	rsource, err := GetMergeRequestResource(a)
	if err != nil {
		log.Fatalf("Unrecognized MergeRequestResource: %+v\n", a)
	}
	nsource, err := GetMergeRequestResource(b)
	if err != nil {
		log.Fatalf("Unrecognized MergeRequestResource: %+v\n", b)
	}
	if nsource.Uri != "" {
		rsource.Uri = nsource.Uri
	}
	if nsource.PrivateToken != "" {
		rsource.PrivateToken = nsource.PrivateToken
	}
	if nsource.PrivateKey != "" {
		rsource.PrivateKey = nsource.PrivateKey
	}
	if nsource.PrivateToken != "" {
		rsource.PrivateToken = nsource.PrivateToken
	}
	if nsource.Username != "" {
		rsource.Username = nsource.Username
	}
	if nsource.Password != "" {
		rsource.Password = nsource.Password
	}
	rsource.NoSSL = nsource.NoSSL
	rsource.SkipSslVerification = nsource.SkipSslVerification
	return rsource
}

func UpdateSlackNotificationResource(a, b interface{}) interface{} {
	rsource, err := GetSlackNotificationResource(a)
	if err != nil {
		log.Fatalf("Unrecognized SlackNotificationResource: %+v\n", a)
	}
	nsource, err := GetSlackNotificationResource(b)
	if err != nil {
		log.Fatalf("Unrecognized SlackNotificationResource: %+v\n", b)
	}
	if nsource.URL != "" {
		rsource.URL = nsource.URL
	}
	return rsource
}

func UpdatePoolResource(a, b interface{}) interface{} {
	rsource, err := GetPoolResource(a)
	if err != nil {
		log.Fatalf("Unrecognized PoolResource: %+v\n", a)
	}
	nsource, err := GetPoolResource(b)
	if err != nil {
		log.Fatalf("Unrecognized PoolResource: %+v\n", b)
	}
	if nsource.Uri != "" {
		rsource.Uri = nsource.Uri
	}
	if nsource.Branch != "" {
		rsource.Branch = nsource.Branch
	}
	if nsource.Pool != "" {
		rsource.Pool = nsource.Pool
	}
	if nsource.PrivateKey != "" {
		rsource.PrivateKey = nsource.PrivateKey
	}
	if nsource.Username != "" {
		rsource.Username = nsource.Username
	}
	if nsource.Password != "" {
		rsource.Password = nsource.Password
	}
	if nsource.RetryDelay != "" {
		rsource.RetryDelay = nsource.RetryDelay
	}
	return rsource
}

func UpdateSemverResource(a, b interface{}) interface{} {
	rsource, err := GetSemverResourceBase(a)
	if err != nil {
		log.Fatalf("Unrecognized SemverResource: %+v\n", a)
	}
	switch rsource.Driver {
	case SemverResourceDriverGCS:
		return UpdateSemverGCSResource(a, b)
	default:
		log.Fatalf("unsupported semver resource driver:%s", rsource.Driver)
	}
	return rsource
}

func UpdateSemverGCSResource(a, b interface{}) interface{} {
	rsource, err := GetSemverGCSResource(a)
	if err != nil {
		log.Fatalf("Unrecognized SemverGCSResource: %+v\n", a)
	}
	nsource, err := GetSemverGCSResource(b)
	if err != nil {
		log.Fatalf("Unrecognized SemverGCSResource: %+v\n", b)
	}
	if nsource.InitialVersion != "" {
		rsource.InitialVersion = nsource.InitialVersion
	}
	if nsource.Driver != "" {
		rsource.Driver = nsource.Driver
	}
	if nsource.Bucket != "" {
		rsource.Bucket = nsource.Bucket
	}
	if nsource.Key != "" {
		rsource.Key = nsource.Key
	}
	if nsource.JsonKey != "" {
		rsource.JsonKey = nsource.JsonKey
	}
	return rsource
}
