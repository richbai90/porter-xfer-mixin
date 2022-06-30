package xfer

import (
	"fmt"
	"strings"
)

var dockerfileLines StringBuilder

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {
	// Create new Builder.
	dockerfileLines.WriteStrings(
		`ENV PKGID /home/${BUNDLE_USER}/`,
		m.PackageID,
		`.tar.gz
# make sure that the home directory for nonroot exists
RUN mkdir -p /home/${BUNDLE_USER} && \
	chown -R ${BUNDLE_USER}:${BUNDLE_GID} /home/${BUNDLE_USER}
`,
	)

	if m.PackageID != "" {
		dockerfileLines.WriteStrings(
			`COPY --chown=${BUNDLE_USER}:${BUNDLE_USER} `,
			m.PackageID,
			`.tar.gz /home/${BUNDLE_USER}`,
		)
	}

	// TODO: Instead of a copy command in the prebuild phase create a tar file and then copy that to the appropriate folder
	// if cfg.Files != nil {
	// 	for _, src := range cfg.Files {
	// 		dockerfileLines.WriteStrings(
	// 			`RUN mkdir /backup
	// 			COPY --chown=`,
	// 			strconv.Itoa(cfg.Chown),
	// 			`:`,
	// 			strconv.Itoa(cfg.Chown),
	// 			` `,
	// 			src,
	// 			` `,
	// 			`/backup`,
	// 		)
	// 	}
	// }

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


