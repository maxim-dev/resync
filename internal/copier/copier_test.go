package copier

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"resync/pkg/fs"
	"resync/pkg/interfaces"
	"testing"
)

func TestCopy(t *testing.T) {

	memFS := fs.NewMemFS()

	mc := gomock.NewController(t)
	defer mc.Finish()

	mock := interfaces.NewMockStatOpenCreater(mc)
	mockW := interfaces.NewMockWorkable(mc)

	t.Run("file copied", func(t *testing.T) {

		req := require.New(t)
		src := "test.txt"
		content := []byte("hello, world")
		expected := int64(len(content))
		destDir := "./dir"
		dest := filepath.Join(destDir, src)

		errWrite := afero.WriteFile(memFS.AFS, src, content, 0644)
		errWriteDir := memFS.AFS.Mkdir(destDir, 0755)

		statSrc, statSrcErr := memFS.Stat(src)
		statDest, statDestErr := memFS.Stat(dest)
		open, openErr := memFS.Open(src)
		create, createErr := memFS.Create(dest)

		mock.EXPECT().Stat(src).Return(statSrc, statSrcErr).Times(1)
		mock.EXPECT().Stat(dest).Return(statDest, statDestErr).Times(1)
		mock.EXPECT().Open(src).Return(open, openErr).Times(1)
		mock.EXPECT().Create(dest).Return(create, createErr).Times(1)

		mockW.EXPECT().IsInProgress(src).Return(false).Times(1)
		mockW.EXPECT().WorkingOn(src).Return().Times(1)
		mockW.EXPECT().WorkingDone(src).Return().Times(1)

		actual, err := Copy(mock, mockW, src, destDir)

		req.Equal(expected, actual)
		req.NoError(err)
		req.NoError(errWrite)
		req.NoError(errWriteDir)
	})

	t.Run("wrong source", func(t *testing.T) {
		var expected int64

		req := require.New(t)
		src := "./src"
		destDir := "test_copy.txt"

		errWrite := memFS.AFS.Mkdir(src, 0755)

		stat, statErr := memFS.Stat(src)

		mock.EXPECT().Stat(src).Return(stat, statErr).Times(1)
		mockW.EXPECT().IsInProgress(src).Return(false).Times(1)

		actual, err := Copy(mock, mockW, src, destDir)

		req.NoError(errWrite)
		req.Equal(expected, actual)
		req.ErrorContains(err, "is not a regular file")
	})

	t.Run("stat error", func(t *testing.T) {
		var expected int64

		req := require.New(t)
		src := "test.txt"
		content := []byte("hello, world")
		destDir := "test_copy.txt"

		errWrite := afero.WriteFile(memFS.AFS, src, content, 0644)

		stat, _ := memFS.Stat(src)
		statErr := &os.PathError{Err: errors.New("custom error")}

		mock.EXPECT().Stat(src).Return(stat, statErr).Times(1)
		mockW.EXPECT().IsInProgress(src).Return(false).Times(1)

		actual, err := Copy(mock, mockW, src, destDir)

		req.NoError(errWrite)
		req.Equal(expected, actual)
		req.ErrorContains(err, "custom error")
	})

	t.Run("open error", func(t *testing.T) {
		var expected int64

		req := require.New(t)
		src := "test.txt"
		content := []byte("hello, world")
		destDir := "./open-error"
		dest := filepath.Join(destDir, src)

		errWrite := afero.WriteFile(memFS.AFS, src, content, 0644)
		errWriteDir := memFS.AFS.Mkdir(destDir, 0755)

		statSrc, statSrcErr := memFS.Stat(src)
		statDest, statDestErr := memFS.Stat(dest)
		open, openErr := memFS.Open("non-existing-file.txt")

		mock.EXPECT().Stat(src).Return(statSrc, statSrcErr).Times(1)
		mock.EXPECT().Stat(dest).Return(statDest, statDestErr).Times(1)
		mock.EXPECT().Open(src).Return(open, openErr).Times(1)

		mockW.EXPECT().IsInProgress(src).Return(false).Times(1)
		mockW.EXPECT().WorkingOn(src).Return().Times(1)
		mockW.EXPECT().WorkingDone(src).Return().Times(1)

		actual, err := Copy(mock, mockW, src, destDir)

		req.NoError(errWrite)
		req.NoError(errWriteDir)
		req.Equal(expected, actual)
		req.ErrorContains(err, "file does not exist")
	})

	t.Run("create error", func(t *testing.T) {
		var expected int64

		req := require.New(t)
		src := "test.txt"
		content := []byte("hello, world")
		destDir := "./create-error"
		dest := filepath.Join(destDir, src)

		errWrite := afero.WriteFile(memFS.AFS, src, content, 0644)
		errWriteDestDir := memFS.AFS.Mkdir(destDir, 0755)

		statSrc, statSrcErr := memFS.Stat(src)
		statDest, statDestErr := memFS.Stat(dest)
		open, openErr := memFS.Open(src)

		create, _ := memFS.Create(dest)
		createErr := &os.PathError{Err: errors.New("custom error")}

		mock.EXPECT().Stat(src).Return(statSrc, statSrcErr).Times(1)
		mock.EXPECT().Open(src).Return(open, openErr).Times(1)
		mock.EXPECT().Stat(dest).Return(statDest, statDestErr).Times(1)
		mock.EXPECT().Create(dest).Return(create, createErr).Times(1)

		mockW.EXPECT().IsInProgress(src).Return(false).Times(1)
		mockW.EXPECT().WorkingOn(src).Return().Times(1)
		mockW.EXPECT().WorkingDone(src).Return().Times(1)

		actual, err := Copy(mock, mockW, src, destDir)

		req.Equal(expected, actual)
		req.NoError(errWrite)
		req.NoError(errWriteDestDir)
		req.ErrorContains(err, "custom error")
	})

	t.Run("dest-is-dir", func(t *testing.T) {
		var expected int64

		req := require.New(t)
		src := "same.txt"
		content := []byte("hello, world")
		destDir := "./dest-is-dir"
		dest := filepath.Join(destDir, src)

		errWrite := afero.WriteFile(memFS.AFS, src, content, 0644)
		errWriteDestDir := memFS.AFS.Mkdir(dest, 0755) // создаем директорию с именем исходного файла

		statSrc, statSrcErr := memFS.Stat(src)
		statDest, statDestErr := memFS.Stat(dest)

		mock.EXPECT().Stat(src).Return(statSrc, statSrcErr).Times(1)
		mock.EXPECT().Stat(dest).Return(statDest, statDestErr).Times(1)

		mockW.EXPECT().IsInProgress(src).Return(false).Times(1)
		mockW.EXPECT().WorkingOn(src).Return().Times(1)
		mockW.EXPECT().WorkingDone(src).Return().Times(1)

		actual, err := Copy(mock, mockW, src, destDir)

		req.NoError(errWrite)
		req.NoError(errWriteDestDir)
		req.Equal(expected, actual)
		req.NoError(err)
	})

}
