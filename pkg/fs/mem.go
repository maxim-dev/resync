package fs

import (
	"github.com/spf13/afero"
	"io"
	"os"
)

type MemFS struct {
	AFS *afero.Afero
}

func NewMemFS() *MemFS {
	var fs = afero.NewMemMapFs()
	return &MemFS{&afero.Afero{Fs: fs}}
}

func (r *MemFS) Stat(s string) (os.FileInfo, error) {
	return r.AFS.Stat(s)
}

func (r *MemFS) Open(s string) (io.ReadCloser, error) {
	return r.AFS.Open(s)
}

func (r *MemFS) Create(s string) (io.WriteCloser, error) {
	return r.AFS.Create(s)
}

func (r *MemFS) RemoveAll(s string) error {
	return r.AFS.RemoveAll(s)
}

func (r *MemFS) Append(s string) (io.WriteCloser, error) {
	return r.AFS.OpenFile(s, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
