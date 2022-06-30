package main

import (
	"reflect"
	"testing"

	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

func Test_BuildBuildCommand(t *testing.T) {
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
			file, err := tt.args.m.FileSystem.Fs.Open("/testdata/build-input.yaml")
			if err != nil {
				t.FailNow()
			}
			tt.args.m.Context.In = file
			cmd := BuildBuildCommand(tt.args.m)
			cmd.SetArgs([]string{
				"build",
			})
			cmd.Execute()
			if got := BuildBuildCommand(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildBuildCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
