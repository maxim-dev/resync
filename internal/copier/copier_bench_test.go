package copier

import (
	"github.com/spf13/afero"
	"resync/internal/dispenser"
	"resync/pkg/fs"
	"testing"
)

var (
	srcBench                            = "benchmark.test"
	destDirBench                        = "./bench"
	memFS        *fs.MemFS              = fs.NewMemFS()
	wl           *dispenser.WorkingList = dispenser.NewWorkingList()
)

func init() {
	content := []byte("hello, world")

	afero.WriteFile(memFS.AFS, srcBench, content, 0644)
	memFS.AFS.Mkdir(destDirBench, 0755)

}

func BenchmarkCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Copy(memFS, wl, srcBench, destDirBench)
	}
}
