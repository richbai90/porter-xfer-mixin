
ENV PKG /home/${BUNDLE_USER}/{{ .PackageID }}.tar.gz

# Make sure that the home directory for nonroot exists
RUN mkdir -p /home/${BUNDLE_USER} && \
	chown -R ${BUNDLE_USER}:${BUNDLE_GID} /home/${BUNDLE_USER}

{{if (eq .SrcKind .Volume)}}

COPY --chown=${BUNDLE_USER}:${BUNDLE_GID} {{ .PackageID }}.tar.gz /home/${BUNDLE_USER}/{{ .PackageID }}.tar.gz

{{end}}

{{if (eq .SrcKind .Directory)}}

RUN mkdir -p /home/${BUNDLE_USER}/{{.SrcVal}}

{{range $_, $file := .Files}}
COPY --chown=${BUNDLE_USER}:${BUNDLE_GID} {{$file.Path}} /home/${BUNDLE_USER}/{{.SrcVal}}/{{$file.Name}}
{{end}}

WORKDIR /home/${BUNDLE_USER}/{{.SrcVal}}
RUN tar -zvpf /home/${BUNDLE_USER}/{{.PackageID}}.tar.gz .

{{end}}

{{if (eq .SrcKind .URL)}}
RUN apt update && \ 
	apt install -y curl && \ 
	curl -sL {{.SrcVal}} | tar -zvpf /home/${BUNDLE_USER}/{{.PackageID}}.tar.gz -
{{end}}


{{if (eq .SrcKind .File)}}
COPY --chown=${BUNDLE_USER}:${BUNDLE_GID} {{.SrcVal}} {{call .Base .SrcVal}}
tar -zvf /home/${BUNDLE_USER}/{{.PackageID}}.tar.gz {{call .Base .SrcVal}}
{{end}}

# TODO: Work out the logic for archives and repos

