package main

import (
	"reflect"
	"testing"

	"github.com/richbai90/xfer/pkg/xfer"
	"github.com/spf13/cobra"
)

func Test_buildInstallCommand(t *testing.T) {
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
			file, err := tt.args.m.FileSystem.Fs.Open("/testdata/step-input.yaml")
			if err != nil {
				t.FailNow()
			}
			tt.args.m.Context.In = file
			cmd := buildInstallCommand(tt.args.m)
			cmd.Execute()
			if got := buildInstallCommand(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildInstallCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
