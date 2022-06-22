package main

import (
	"reflect"
	"testing"

	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

func Test_buildBuildCommand(t *testing.T) {
	type args struct {
		m *xfer.Mixin
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		{
			name: "debug test",
			args: args{
				m: xfer.NewTestMixin(t).Mixin,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := tt.args.m.FileSystem.Fs.Open("/testdata/build-input-with-client-version.yaml")
			if err != nil {
				t.FailNow()
			}
			tt.args.m.Context.In = file
			cmd := buildBuildCommand(tt.args.m)
			cmd.SetArgs([]string{
				"build",
			})
			cmd.Execute()
			if got := buildBuildCommand(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildBuildCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
