package xfer

import (
	"errors"
	"io/fs"
	"os"
	"path"

	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/richbai90/xfer/pkg/xfer/testdata"
)

func Client() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv)
}

type URLDetails struct {
	URL string
	As string
}

// Archive the files if they come from a volume.
// Must be done pre-build as the volume cannot easily be mounted during build
func (m *Mixin) PreBuild() error {
	id := uuid.Generate().String()
	m.PrintDebug("Pre Build started")
	m.PrintDebug("ID: %s", m.PackageID)
	// store the id in a config option so that build can use it during build
	m.PackageID = id
	setupDebugInput(m)

	if err := PopulateInput(m, &Input); err != nil {
		return err
	}
	var errs []error
	for _, source := range Input.Config.Sources {
		switch source.Kind {
		case Volume:
			errs = append(errs, m.handleVolumeSource(source.Value))
			break
		case URL:
			errs = append(errs, m.handleURLSource(URLDetails{URL: source.Value, As: source.As}))
			break
		case File:
			errs = append(errs, m.handleFileSource(source.Value))
			break
		case Directory:
			errs = append(errs, m.handleDirSource(source.Value))
			break
		default:
			errs = append(errs, errors.New("Not Implemented"))
		}
	}

	err := m.HandleErrs(errs, "Problem processings sources")
	return err
}

func setupDebugInput(m *Mixin) {
	r, w, _ := os.Pipe()
	if _, dbg := os.LookupEnv("debugger"); dbg {
		m.IO = w
		w.WriteString(testdata.BuildInput)
		m.Context.In = r
	} else {
		m.IO = os.Stdin
	}

	defer w.Close()
}

func (m *Mixin) handleURLSource(dets URLDetails) error {
	m.URLs = append(m.URLs, dets)
	return nil
}

func (m *Mixin) handleFileSource(file string) error {
	m.Files = append(m.Files, file)
	return nil
}

func (m *Mixin) handleDirSource(dir string) error {
	if m.Directories == nil {
		m.Directories = make(map[string][]string)
	}
	err := m.FileSystem.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		m.Directories[dir] = append(m.Directories[dir], path)
		return nil
	})

	return err
}

func (m *Mixin) handleVolumeSource(volume string) error {
	dcli, err := Client()
	if err != nil {
		return err
	}
	volid := uuid.Generate().String()
	m.PrintDebug("Inspecting Provided Volume %s", volume)
	v, err := dcli.VolumeInspect(m.Ctx, volume)

	if m.HandleErr(&err, "Problem inspecting volume %s: ", volume) {
		return err
	}

	m.PrintDebug("Package ID: %s", volid)
	m.PrintDebug("Expected Output File: ", path.Join(m.Getwd(), volid+".tar.gz"))

	cmd := []string{`/bin/sh`, `-c`, `cd /src && tar -czvf /dest/` + volid + `.tar.gz .`}

	m.PrintDebug("Pulling debian image")
	reader, err := dcli.ImagePull(m.Ctx, "debian:latest", types.ImagePullOptions{})
	defer reader.Close()
	if m.HandleErr(&err, "Problem pulling docker image debian:latest") {
		return err;
	}

	backupConfig := container.Config{
		AttachStderr: true,
		Cmd:          cmd,
		Image:        "debian:latest",
	}

	backupHostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     "volume",
				Source:   v.Name,
				Target:   "/src",
				ReadOnly: true,
			},
			{
				Type: "bind",
				Source: m.Getwd(),
				Target: "/dest",
				ReadOnly: false,
			},
		}}
	m.PrintDebug("Starting xfer container")
	dkr, err := dcli.ContainerCreate(m.Ctx, &backupConfig, &backupHostConfig, nil, nil, "xfer")
	if m.HandleErr(&err) {
		dcli.ContainerKill(m.Ctx, dkr.ID, "SIGKILL")
		return err
	}

	err = dcli.ContainerStart(m.Ctx, dkr.ID, types.ContainerStartOptions{})
	chok, cherr := dcli.ContainerWait(m.Ctx, dkr.ID, container.WaitConditionNotRunning)

	select {
	case e := <-cherr:
		if m.HandleErr(&e) {
			// container isn't stopped, force the removal
			dcli.ContainerRemove(m.Ctx, dkr.ID, types.ContainerRemoveOptions{
				Force: true,
			})
			return e
		}
		break
	case <-chok:
		dcli.ContainerRemove(m.Ctx, dkr.ID, types.ContainerRemoveOptions{})
	}

	if m.HandleErr(&err) {
		return err
	}

	m.Volumes = append(m.Volumes, volid+".tar.gz")
	return nil
}
