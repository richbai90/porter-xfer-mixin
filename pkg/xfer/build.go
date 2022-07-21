package xfer

import (
	"path"
	"text/template"

	"github.com/richbai90/xfer/pkg/xfer/resources"
)

type TplData struct {
	PackageID   string
	Volume      bool
	Directory   bool
	SrcVal      string
	URL         bool
	File        bool
	Base        func(string) string
	Volumes     []string
	Directories map[string][]string
	URLs        []URLDetails
	Files       []string
}

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build() error {
	// Create new Builder

	tplData := TplData{
		PackageID:   m.PackageID,
		Volume:      len(m.Volumes) > 0,
		Directory:   len(m.Directories) > 0,
		URL:         len(m.URLs) > 0,
		File:        len(m.Files) > 0,
		Base:        path.Base,
		Volumes:     m.Volumes,
		Directories: m.Directories,
		URLs:        m.URLs,
		Files:       m.Files,
	}

	tpl, err := template.ParseFS(resources.FS, "Dockerfile.tpl")

	if m.HandleErr(&err, "Unable to parse Dockerfile Template") {
		return err
	}

	err = tpl.Execute(m.Out, tplData)
	return err
}
