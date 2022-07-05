package xfer

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var dockerfileLines StringBuilder

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {
	// Create new Builder.
	var build BuildInput
	err := PopulateInput(m, &build)
	if m.HandleErr(&err, "Unable to populate input") {
		return err
	}
	
	source := build.Config.Source
	basedir := "/home/${BUNDLE_USER}"

	archive := fmt.Sprintf("%s.tar.gz", m.PackageID)

	dockerfileLines.WriteStrings(
		`ENV PKGID /home/${BUNDLE_USER}/`,
		archive,
		`
# make sure that the home directory for nonroot exists
RUN mkdir -p /home/${BUNDLE_USER} && \
	chown -R ${BUNDLE_USER}:${BUNDLE_GID} `,
		basedir,
		"\n",
	)

	// If the provided source was a volume then everything was done in pre-build so just copy the tar file we know will be there
	if source.Kind == Volume {
		dockerfileLines.WriteStrings(
			`COPY --chown=${BUNDLE_USER}:${BUNDLE_GID} `,
			archive,
			` /home/${BUNDLE_USER}
`,
		)
	}

	if source.Kind == Directory {
		if _, err := os.Stat(source.Value); m.HandleErr(&err, "Could not stat directory %s ", source.Value) {
			return err
		}
		files, err := m.Context.FileSystem.ReadDir(source.Value)
		if m.HandleErr(&err, "Could not read directory %s ", source.Value) {
			return err
		}
		dir := fmt.Sprintf("%s/%s", basedir, path.Base(source.Value))
		dockerfileLines.WriteStrings(`RUN mkdir -p `, dir, "\n")

		for _, f := range files {
			dockerfileLines.WriteStrings(
				`COPY --chown=${BUNDLE_USER}:0 `,
				path.Join(source.Value, f.Name()),
				` /home/${BUNDLE_USER}/`,
				path.Base(source.Value),
				"\n",
			)
		}

		dockerfileLines.WriteStrings("WORKDIR ", dir,
			`
RUN tar -xvf `,
			dir,
			`/`,
			archive,
		)

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
