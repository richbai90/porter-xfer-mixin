package xfer

import (
	"os/exec"

	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the FileTransfer mixin in porter.yaml
// mixins:
// - FileTransfer:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
	// The path on the container to copy the files to if Files is specified
	Destination string `yaml:"destination,omitempty"`
	// The files to be copied into the container at build time
	Files []string `yaml:"files,omitempty"`
	// The name of a predefined volume to use instead of a file list
	Volume string `yaml:"volume,omitempty"`
	// The ID to pass to the chown flag of docker copy
	Chown int `yaml:"chown,omitempty"`
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