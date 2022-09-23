package interfaces

//go:generate mockgen -source=./fs.go -destination=./fs_mocks.go -package=interfaces

import (
	"io"
	"os"
)

type Stater interface {
	Stat(string) (os.FileInfo, error)
}

type Opener interface {
	Open(string) (io.ReadCloser, error)
}

type Creater interface {
	Create(string) (io.WriteCloser, error)
}

type Remover interface {
	RemoveAll(string) error
}

type Appender interface {
	Append(string) (io.WriteCloser, error)
}

type StatRemover interface {
	Stater
	Remover
}

type StatOpenCreater interface {
	Stater
	Opener
	Creater
}

type Fs interface {
	Stater
	Opener
	Creater
	Remover
	Appender
}
