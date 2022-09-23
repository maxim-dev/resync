package fs_struct

import (
	"github.com/spf13/afero"
	"os"
)

// Важно! Вначале см. файл fs.go

// В этом файле приведен пример того, почему абстракция на структурах неудобна
// При в качестве реализации методов абстрактной файловой системы (далее "ФС") используется In Memory ФС из пакета "afero"
// Но сигнатуры методов Open и Create отличаются. Поэтому напрямую можно присвоить только поле Stat
// В двух других функциях приходится использовать приведение типов к *os.File

/*
func NewMemFS() *AbsFS {
	var (
		fs  = afero.NewMemMapFs()
		AFS = &afero.Afero{Fs: fs}
	)

	return &AbsFS{
		Stat: AFS.Stat,
		Open: AFS.Open,
		Create: AFS.Create,
	}
}
*/

func NewMyMemFS() *AbsFS {
	var (
		fs  = afero.NewMemMapFs()
		AFS = &afero.Afero{Fs: fs}
	)

	return &AbsFS{
		Stat: AFS.Stat,
		Open: func(s string) (*os.File, error) {
			file, err := AFS.Open(s)
			return file.(*os.File), err
		},
		Create: func(s string) (*os.File, error) {
			file, err := AFS.Create(s)
			return file.(*os.File), err
		},
	}
}
