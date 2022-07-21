//go:generate packr2
package xfer

import (
	"context"
	"io"
	"os"

	"get.porter.sh/porter/pkg/portercontext"

	"github.com/spf13/cobra"
)

const defaultClientVersion string = "v0.0.0"

type ExpandedContext struct {
	*portercontext.Context
	Cmd *cobra.Command
	Ctx context.Context
	IO  io.ReadWriteCloser
}

type Mixin struct {
	ExpandedContext
	ClientVersion string
	PackageID     string
	Volumes       []string
	Files         []string
	Directories   map[string][]string
	URLs          []URLDetails
	//add whatever other context/state is needed here
}

// New azure mixin client, initialized with useful defaults.
func New() (*Mixin, error) {
	return &Mixin{
		ExpandedContext: ExpandedContext{Context: portercontext.New()},
		ClientVersion:   defaultClientVersion,
	}, nil
}

func(e *ExpandedContext) Getwd() string {
	if _, exists := os.LookupEnv("debugger"); exists {
		return os.Getenv("HostWd")
	} else {
		return e.FileSystem.Getwd()
	}
}
