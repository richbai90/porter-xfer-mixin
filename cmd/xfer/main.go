package main

import (
	"fmt"
	"io"
	"os"

	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

var commands = make(map[string]*cobra.Command)

func main() {
	cmd, err := BuildRootCommand(os.Stdin)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}
}

func BuildRootCommand(in io.Reader) (*cobra.Command, error) {
	m, err := xfer.New()
	if err != nil {
		return nil, err
	}
	m.In = in
	cmd := &cobra.Command{
		Use:  "xfer",
		Long: "A skeleton mixin to use for building other mixins for porter 👩🏽‍✈️",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Enable swapping out stdout/stderr for testing
			m.Out = cmd.OutOrStdout()
			m.Err = cmd.OutOrStderr()
			m.Cmd = cmd
			m.Ctx = cmd.Context()
		},
		SilenceUsage: true,
	}

	cmd.PersistentFlags().BoolVar(&m.Debug, "debug", false, "Enable debug logging")
	cmd.AddCommand(BuildVersionCommand(m))
	cmd.AddCommand(BuildSchemaCommand(m))
	cmd.AddCommand(BuildBuildCommand(m))
	cmd.AddCommand(BuildInstallCommand(m))
	cmd.AddCommand(BuildInvokeCommand(m))
	cmd.AddCommand(BuildUpgradeCommand(m))
	cmd.AddCommand(BuildUninstallCommand(m))

	return cmd, nil
}