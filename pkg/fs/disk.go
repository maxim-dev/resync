package fs

import (
	"io"
	"os"
)

type DiskFS struct{}

func NewDiskFS() *DiskFS {
	return &DiskFS{}
}

func (r *DiskFS) Stat(s string) (os.FileInfo, error) {
	return os.Stat(s)
}

func (r *DiskFS) Open(s string) (io.ReadCloser, error) {
	return os.Open(s)
}

func (r *DiskFS) Create(s string) (io.WriteCloser, error) {
	return os.Create(s)
}

func (r *DiskFS) RemoveAll(s string) error {
	return os.RemoveAll(s)
}

func (r *DiskFS) Append(s string) (io.WriteCloser, error) {
	return os.OpenFile(s, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
