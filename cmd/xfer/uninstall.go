package main

import (
	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

func BuildUninstallCommand(m *xfer.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Execute the uninstall functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute()
		},
	}
	return cmd
}
