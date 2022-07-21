
ENV PKG /home/${BUNDLE_USER}/{{.PackageID}}.tar.gz

# Make sure that the home directory for nonroot exists
RUN mkdir -p /home/${BUNDLE_USER}/files && \
	chown -R ${BUNDLE_USER}:${BUNDLE_GID} /home/${BUNDLE_USER}
VOLUME /home/${BUNDLE_USER}/files
{{if (.Volume)}}
{{range $_, $volume := .Volumes}}
COPY --chown=${BUNDLE_USER}:${BUNDLE_GID} {{$volume}} /home/${BUNDLE_USER}/files/{{$volume}}
RUN cd /home/${BUNDLE_USER}/files && tar -zxvf {{$volume}} && rm {{$volume}}
{{end}}
{{end}}

{{if (.Directory)}}
{{range $directory, $files := .Directories}}
RUN mkdir -p /home/${BUNDLE_USER}/files/{{$directory}}
{{range $_, $file := $files}}
COPY --chown=${BUNDLE_USER}:${BUNDLE_GID} {{$file}} /home/${BUNDLE_USER}/files/{{$file}}
{{end}}
{{end}}
{{end}}

{{if (.URL)}}
RUN apt-get update && \
	apt-get install -y curl
{{range $_, $url := .URLs}}
	RUN curl -L --globoff "{{$url.URL}}" --output "/home/${BUNDLE_USER}/files/{{$url.As}}"
{{end}}
{{end}}


{{if (.File)}}
{{range $_, $file := .Files}}
COPY --chown=${BUNDLE_USER}:${BUNDLE_GID} {{$file}} /home/${BUNDLE_USER}/files/{{call $.Base $file}}
{{end}}
{{end}}

WORKDIR /home/${BUNDLE_USER}/files
RUN tar -zcvpf ${PKG} .

WORKDIR ${BUNDLE_DIR}

# TODO: Work out the logic for archives and repos

