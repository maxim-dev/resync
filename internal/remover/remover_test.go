package remover

import (
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"resync/pkg/fs"
	"resync/pkg/interfaces"
	"testing"
)

func TestRemove(t *testing.T) {

	mc := gomock.NewController(t)
	defer mc.Finish()

	mock := interfaces.NewMockWorkable(mc)

	t.Run("remove a file", func(t *testing.T) {
		memFS := fs.NewMemFS()

		req := require.New(t)

		rootDir := "src"
		destDir := "dest"
		file := "test.txt"
		content := []byte("hello, world")
		destFile := filepath.Join(destDir, file)

		errWriteSrcDir := memFS.AFS.Mkdir(rootDir, 0755)
		errWriteDestDir := memFS.AFS.Mkdir(destDir, 0755)
		errWriteFile := afero.WriteFile(memFS.AFS, destFile, content, 0644)

		mock.EXPECT().IsInProgress(destFile).Return(false).Times(1)
		mock.EXPECT().WorkingOn(destFile).Return().Times(1)
		mock.EXPECT().WorkingDone(destFile).Return().Times(1)

		isDeleted, err := Remove(memFS, mock, file, rootDir, destDir)

		// Проверяем, что файл удален
		actual, existsErr := memFS.AFS.Exists(destFile)

		req.NoError(errWriteSrcDir)
		req.NoError(errWriteDestDir)
		req.NoError(errWriteFile)

		req.Equal(false, actual)
		req.Equal(true, isDeleted)
		req.NoError(existsErr)
		req.NoError(err)
	})

	t.Run("remove a dir", func(t *testing.T) {
		memFS := fs.NewMemFS()

		req := require.New(t)

		rootDir := "src"
		destDir := "dest"
		destSubDir := "sub"
		destSubDirPath := filepath.Join(destDir, destSubDir)

		errWriteSrcDir := memFS.AFS.Mkdir(rootDir, 0755)
		errWriteDestSubDir := memFS.AFS.MkdirAll(destSubDirPath, 0755)

		mock.EXPECT().IsInProgress(destSubDirPath).Return(false).Times(1)
		mock.EXPECT().WorkingOn(destSubDirPath).Return().Times(1)
		mock.EXPECT().WorkingDone(destSubDirPath).Return().Times(1)

		isDeleted, err := Remove(memFS, mock, destSubDir, rootDir, destDir)

		// Проверяем, что сущность удалена
		actual, existsErr := memFS.AFS.Exists(destSubDirPath)

		req.NoError(errWriteSrcDir)
		req.NoError(errWriteDestSubDir)

		req.Equal(false, actual)
		req.Equal(true, isDeleted)
		req.NoError(existsErr)
		req.NoError(err)

	})

	t.Run("leave a file", func(t *testing.T) {
		memFS := fs.NewMemFS()

		req := require.New(t)

		rootDir := "src"
		destDir := "dest"
		file := "test.txt"
		content := []byte("hello, world")
		rootFile := filepath.Join(rootDir, file)
		destFile := filepath.Join(destDir, file)

		errWriteSrcDir := memFS.AFS.Mkdir(rootDir, 0755)
		errWriteDestDir := memFS.AFS.Mkdir(destDir, 0755)
		errWriteRootFile := afero.WriteFile(memFS.AFS, rootFile, content, 0644)
		errWriteDestFile := afero.WriteFile(memFS.AFS, destFile, content, 0644)

		mock.EXPECT().IsInProgress(destFile).Return(false).Times(1)
		mock.EXPECT().WorkingOn(destFile).Return().Times(1)
		mock.EXPECT().WorkingDone(destFile).Return().Times(1)

		isDeleted, err := Remove(memFS, mock, file, rootDir, destDir)

		// Проверяем, что файл остался
		actual, existsErr := memFS.AFS.Exists(destFile)

		req.NoError(errWriteSrcDir)
		req.NoError(errWriteDestDir)
		req.NoError(errWriteRootFile)
		req.NoError(errWriteDestFile)

		req.Equal(true, actual)
		req.Equal(false, isDeleted)
		req.NoError(existsErr)
		req.NoError(err)

	})

	t.Run("leave a dir", func(t *testing.T) {
		memFS := fs.NewMemFS()

		req := require.New(t)

		rootDir := "src"
		destDir := "dest"
		subDir := "sub"
		destRootDirPath := filepath.Join(rootDir, subDir)
		destSubDirPath := filepath.Join(destDir, subDir)

		errWriteSrcDir := memFS.AFS.MkdirAll(destRootDirPath, 0755)
		errWriteDestDir := memFS.AFS.MkdirAll(destSubDirPath, 0755)

		mock.EXPECT().IsInProgress(destSubDirPath).Return(false).Times(1)
		mock.EXPECT().WorkingOn(destSubDirPath).Return().Times(1)
		mock.EXPECT().WorkingDone(destSubDirPath).Return().Times(1)

		isDeleted, err := Remove(memFS, mock, subDir, rootDir, destDir)

		// Проверяем, что сущность осталась
		actual, existsErr := memFS.AFS.Exists(destSubDirPath)

		req.NoError(errWriteSrcDir)
		req.NoError(errWriteDestDir)

		req.Equal(true, actual)
		req.Equal(false, isDeleted)
		req.NoError(existsErr)
		req.NoError(err)

	})
}
