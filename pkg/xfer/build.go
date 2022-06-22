package xfer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	Xfer struct {
		// The path on the container to copy the files to if Files is specified
		Destination string `yaml:"destination,omitempty"`
		// The files to be copied into the container at build time
		Files []string `yaml:"files,omitempty"`
		// The name of a predefined volume to use instead of a file list
		Volume string `yaml:"volume,omitempty"`
		// The ID to pass to the chown flag of docker copy
		Chown int `yaml:"chown,omitempty"`
	} `yaml:"xfer,omitempty"`
}

var dockerfileLines StringBuilder

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {
	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	cfg := input.Config.Xfer

	if (cfg.Files != nil && cfg.Destination == "") || (cfg.Destination != "" && cfg.Files == nil) {
		return errors.New("Invalid Mixin Configuration: Destination and Files must be supplied together or not at all")
	}

	if cfg.Chown > 0 {
		dockerfileLines.WriteStrings(
			`RUN useradd --uid `,
			strconv.Itoa(cfg.Chown),
			` --gid `,
			strconv.Itoa(cfg.Chown),
			` porter`,
		)
	}

	if cfg.Chown < 0 {
		return errors.New("Invalid Mixin Configuration: chown must be an integer value > 0 or left unspecified")
	}

	if cfg.Volume != "" {
		dockerfileLines.WriteStrings(
			`RUN apt-get update && \
				 apt-get install jq && \
				 docker volume inspect `,
			cfg.Volume,
			// this is kind of a hack. Backup will either be a file or a folder, but not both.
			// the check is done in the xfer-bin program
			` | jq '.[0]' > /backup && \
				 cd `,
			`$( \
					docker inspect $(docker ps -a --filter volume=`,
			cfg.Volume,
			` \
						| awk '(NR>1) {print $1}' \
						| sed 's/ //g') \
						| jq '.[0].Mounts[] | select(.Name == "`,
			cfg.Volume,
			`").Destination' \
			&& tar -czvf /backup.tar.gz .
			`,
		)

		if _, err := fmt.Fprint(m.Out, dockerfileLines); err != nil {
			return err
		}
	}

	if cfg.Files != nil {
		for _, src := range cfg.Files {
			dockerfileLines.WriteStrings(
				`RUN mkdir /backup
				COPY --chown=`,
				strconv.Itoa(cfg.Chown),
				`:`,
				strconv.Itoa(cfg.Chown),
				` `,
				src,
				` `,
				`/backup`,
			)
		}
	}

	fmt.Fprint(m.Out, dockerfileLines.String())
	return nil
}

type StringBuilder struct {
	strings.Builder
}

func (b *StringBuilder) WriteStrings(strings ...string) (int, error) {
	bytes := 0
	for _, str := range strings {
		if n, err := b.WriteString(str); err != nil {
			return bytes, err
		} else {
			bytes += n
		}
	}

	return bytes, nil
}
