package xfer

import (
	"path"
	"text/template"

	"github.com/richbai90/xfer/pkg/xfer/resources"
)

type TplData struct {
	PackageID string
	Volume bool
	Directory bool
	SrcVal string
	URL bool
	File bool
	Base func(string) string
}

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {
	// Create new Builder
	var build BuildInput
	err := PopulateInput(m, &build)
	if m.HandleErr(&err, "Unable to populate input") {
		return err
	}

	tplData := TplData{
		PackageID: m.PackageID,
		Volume: build.Config.Source.Kind == Volume,
		Directory: build.Config.Source.Kind == Directory,
		SrcVal: build.Config.Source.Value,
		URL: build.Config.Source.Kind == URL,
		File: build.Config.Source.Kind == File,
		Base: path.Base,
	}

	tpl, err := template.ParseFS(resources.FS, "Dockerfile.tpl")

	if m.HandleErr(&err, "Unable to parse Dockerfile Template") {
		return err;
	}

	err = tpl.Execute(m.Out, tplData)
	return err
}


