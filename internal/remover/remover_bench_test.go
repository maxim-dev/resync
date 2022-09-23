package remover

import (
	"github.com/spf13/afero"
	"path/filepath"
	"resync/internal/dispenser"
	"resync/pkg/fs"
	"testing"
)

var (
	rootDir  = "src"
	destDir  = "dest"
	file     = "benchmark.test"
	destFile = filepath.Join(destDir, file)
	memFS    = fs.NewMemFS()
	wl       = dispenser.NewWorkingList()
)

func init() {
	content := []byte("hello, world")

	memFS.AFS.Mkdir(rootDir, 0755)
	memFS.AFS.Mkdir(destDir, 0755)
	afero.WriteFile(memFS.AFS, destFile, content, 0644)

}

func BenchmarkRemove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Remove(memFS, wl, file, rootDir, destDir)
	}
}
