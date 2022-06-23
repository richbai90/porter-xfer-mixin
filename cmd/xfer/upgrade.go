package main

import (
	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *xfer.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute()
		},
	}
	return cmd
}
