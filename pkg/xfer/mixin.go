//go:generate packr2
package xfer

import (
	"context"

	"get.porter.sh/porter/pkg/portercontext"

	"github.com/spf13/cobra"
)

const defaultClientVersion string = "v0.0.0"

type ExpandedContext struct {
	*portercontext.Context
	Cmd *cobra.Command
	Ctx context.Context
}

type Mixin struct {
	ExpandedContext
	ClientVersion string
	PackageID     string
	//add whatever other context/state is needed here
}

// New azure mixin client, initialized with useful defaults.
func New() (*Mixin, error) {
	return &Mixin{
		ExpandedContext: ExpandedContext{Context: portercontext.New()},
		ClientVersion:   defaultClientVersion,
	}, nil

}
