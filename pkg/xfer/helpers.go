package xfer

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	context "get.porter.sh/porter/pkg/portercontext"
	"github.com/carolynvs/aferox"
	"github.com/pkg/errors"
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
			ExpandedContext: ExpandedContext{
				Context: c.Context,
			},
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

func (m *Mixin) HandleErr(err *error, args ...interface{}) bool {
	if len(args) == 0 {
		args = append(args, "")
	}
	e := errors.Wrap(*err, fmt.Sprintf(args[0].(string), args...))
	if m.Context.Debug && e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		*err = e
		return true
	}

	return false
}

func (m *Mixin) PrintDebug(format string, a ...interface{}) {
	if m.Debug {
		format = fmt.Sprintf("=== DEBUG ===\n%s\n", format)
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

// Return true if any of the items in the list cause the callback function f to return true when given as a parameter
func Some(list []interface{}, f func(interface{}, int) bool) bool {
	for i, iface := range list {
		if f(iface, i) {
			return true
		}
	}

	return false
}

// Return the first item in the list that causes the callback function f to return true when given as a parameter
func Find(list []interface{}, f func(interface{}, int) bool) (interface{}, int) {
	for i, iface := range list {
		if f(iface, i) {
			return iface, i
		}
	}

	return nil, -1
}

func (m *Mixin) HandleErrs(errs []error, args ...interface{}) error {
	// If we have no errors in the slice no need to do anything else
	if len(errs) == 0 {
		return nil
	}
	// Create a slice of interfaces that will hold pointers to our errors
	ifaces := make([]interface{}, len(errs))
	for i, e := range errs {
		ifaces[i] = &e
	}
	
	// Search if any of our errors are not nill
	if err, _ := Find(ifaces, func(err interface{}, _ int) bool {
		// We know that our err interface is a *error so convert without checking
		e := err.(*error)
		// Returns true if e is not nil
		return m.HandleErr(e, args...)
	}); err != nil {
		// Return the found error pointer
		e := err.(*error)
		return *e
	}
	return nil
}
