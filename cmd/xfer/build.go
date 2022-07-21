package main

import (
	"os"
	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

func BuildBuildCommand(m *xfer.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Generate Dockerfile lines for the bundle invocation image",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return m.PreBuild()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, set := os.LookupEnv("debugger"); set {
				// if the debugger env is set then In is a pipe and that pipe must be closed
				// This is strictly for use during debugging sessions
				defer m.Context.In.(*os.File).Close()
			}
			defer cmd.PostRun(cmd, args)
			return m.Build()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// remove the tar file now that it's part of the bundle
			m.PrintDebug(`Cleaning up
Remove %s`, m.PackageID + ".tar.gz")
			err := m.FileSystem.Remove(m.PackageID + ".tar.gz")
			if m.HandleErr(&err, "Warning! Cleanup Failed") {
				m.PrintDebug(err.Error())
			}
		},
	}
	return cmd
}
