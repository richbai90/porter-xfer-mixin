package xfer

import (
	"encoding/json"
	"os/exec"

	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
)

type KindType string

const (
	Directory KindType = "directory"
	URL       KindType = "url"
	Repo      KindType = "repo"
	Volume    KindType = "volume"
	Archive   KindType = "archive"
	File      KindType = "file"
)

var kindValueToTypeMap = map[string]KindType{
	"directory": Directory,
	"url":       URL,
	"repo":      Repo,
	"volume":    Volume,
	"archive":   Archive,
	"file":      File,
}

var kindTypeToValueMap = map[KindType]string{
	Directory: "directory",
	URL:       "url",
	Repo:      "repo",
	Volume:    "volume",
	Archive:   "archive",
	File:      "file",
}

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

type MixinSource struct {
	Kind  KindType `yaml:"kind,omitempty"`
	Value string   `yaml:"value,omitempty"`
	As    string   `yaml:"as,omitempty"`
}

// MixinConfig represents configuration that can be set on the FileTransfer mixin in porter.yaml
// mixins:
// - FileTransfer:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
	// The path on the container to copy the files to if Files is specified
	Destination string `yaml:"destination,omitempty"`
	// Source file can be a file path, volume name, or url
	Sources []MixinSource `yaml:"sources,omitempty"`
	// Owner ID to use during install
	Chown string `yaml:"chown,omitempty"`
}

func (k KindType) String() string {
	return kindTypeToValueMap[k]
}

func (k KindType) MarshalJSON() ([]byte, error) {
	return []byte(kindTypeToValueMap[k]), nil
}

func (k *KindType) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	*k = kindValueToTypeMap[str]
	return nil
}

func (k KindType) MarshalYAML() ([]byte, error) {
	return []byte(kindTypeToValueMap[k]), nil
}

func (k *KindType) UnmarshalYAML(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	*k = kindValueToTypeMap[str]
	return nil
}

func PopulateInput(m *Mixin, b *BuildInput) error {
	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, b)
		if b.Config.ClientVersion == "" {
			b.Config.ClientVersion = getDockerVersion()
		}
		return b, err
	})
	if m.HandleErr(&err) {
		return err
	}

	return nil
}

func getDockerVersion() string {
	cmd := exec.Command("docker", "version", "--format", "{{.Client.Version}}")
	ver, err := cmd.Output()

	if err != nil {
		// If there was some problem getting the version -- maybe the format changed? -- choose a sensible default
		ver = []byte("20.10.7")
	}

	return string(ver)

}

var Input BuildInput
