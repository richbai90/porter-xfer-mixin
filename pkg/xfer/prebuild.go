package xfer

import (
	"os"
	"path"

	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func Client() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv)
}

// Backup the files so they're ready to be extracted later
func (m *Mixin) PreBuild() error {
	dcli, err := Client()
	setupDebugInput(m)
	if m.HandleErr(&err) {
		return err
	}
	var input BuildInput

	if err := PopulateInput(m, &input); err != nil {
		return err
	}
	cfg := input.Config
	m.PrintDebug("Inspecting Provided Volume %s", cfg.Volume)
	v, err := dcli.VolumeInspect(m.Ctx, cfg.Volume)

	if m.HandleErr(&err, "Problem inspecting volume %s: ", cfg.Volume) {
		return err
	}

	id := uuid.Generate().String()
	// store the id in a private config option so that build can use it later
	m.PackageID = id

	m.PrintDebug("Package ID: %s", id)
	m.PrintDebug("Expected Output File: ", path.Join(m.Getwd(), id+".tar.gz"))

	cmd := []string{`/bin/sh`, `-c`, `cd /src && tar -czvf /dest/` + id + `.tar.gz .`}
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
				Type:   "bind",
				Source: m.Getwd(),
				Target: "/dest",
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

	return nil
}

func setupDebugInput(m *Mixin) {
	r, w, _ := os.Pipe()
	input := "config:\n  volume: test\nactions:\n  install:\n    - xfer:\n        description: File Transfer\n        destination: /Users/Rich/restore\n  uninstall:\n    - xfer:\n        description: Obligatory uninstall step\n  upgrade: []\n"
	if _, dbg := os.LookupEnv("debugger"); dbg {
		w.WriteString(input)
		m.Context.In = r
	}

	defer w.Close()
}
