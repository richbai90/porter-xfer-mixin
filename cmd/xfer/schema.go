package main

import (
	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

func BuildSchemaCommand(m *xfer.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Print the json schema for the mixin",
		Run: func(cmd *cobra.Command, args []string) {
			m.PrintSchema()
		},
	}
	return cmd
}
