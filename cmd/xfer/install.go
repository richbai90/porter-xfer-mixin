package main

import (
	"os"
	"strings"

	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

var (
	commandFile string
)

func BuildInstallCommand(m *xfer.Mixin) *cobra.Command {
	r, w, _ := os.Pipe()
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Execute the install functionality of this mixin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Handle debug pipe
			if _, dbg := os.LookupEnv("debugger"); dbg {
				w.WriteString("install:\n  - xfer:\n      description: File Transfer\n      destination: /Users/Rich/restore\n")
				m.Context.In = r
			} else if dbg, _ := cmd.Flags().GetBool("debug"); dbg {
				m.PrintDebug("%s", strings.Join(args, " "))
			}
			defer w.Close()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			defer r.Close()
			return m.Execute()
		},
	}
	return cmd
}
