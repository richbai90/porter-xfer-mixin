package xfer

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"get.porter.sh/porter/pkg/context"
	"github.com/carolynvs/aferox"
	"github.com/spf13/afero"
)

type TestMixin struct {
	*Mixin
	TestContext *MyTestContext
}

type MyTestContext struct {
	*context.TestContext
}

// NewTestMixin initializes a mixin test client, with the output buffered, and an in-memory file system.
func NewTestMixin(t *testing.T) *TestMixin {
	c := NewContext(context.NewTestContext(t))
	m := &TestMixin{
		Mixin: &Mixin{
			Context: c.Context,
		},
		TestContext: c,
	}

	_, b, _, ok := runtime.Caller(0)
	basepath := filepath.Dir(b)
	if ok {
		c.CopyDirectoryToTestFs(path.Join(basepath, "testdata"), true)
	}

	return m
}

func NewContext(cxt *context.TestContext) *MyTestContext {
	return &MyTestContext{
		TestContext: cxt,
	}
}

func (c *MyTestContext) CopyDirectoryToTestFs(srcDir string, includeBaseDir bool) error {
	var stripPrefix string
	if includeBaseDir {
		stripPrefix = filepath.Dir(srcDir)
	} else {
		stripPrefix = srcDir
	}

	osFs := aferox.NewAferox(srcDir, afero.NewOsFs())

	return osFs.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Translate the path from the src to the final destination
		dest := filepath.Join("/", strings.TrimPrefix(path, stripPrefix))
		if dest == "" {
			return nil
		}

		if info.IsDir() {
			return c.FileSystem.MkdirAll(dest, info.Mode())
		}

		file, err := osFs.Open(path)

		if err != nil {
			return err
		}

		return CopyFile(file, c.FileSystem.Fs, dest)
	})
}

func CopyFile(file afero.File, Fs afero.Fs, filepath string) error {
	outFile, _ := Fs.Create(filepath)

	buf := make([]byte, 1024)

	for {
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			outFile.Close()
			file.Close()
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := outFile.Write(buf[:n]); err != nil {
			outFile.Close()
			file.Close()
			return err
		}
	}

	return nil
}
