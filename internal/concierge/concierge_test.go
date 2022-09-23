package concierge

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/fs"
	fsPkg "resync/pkg/fs"
	"resync/pkg/interfaces"
	"testing"
	"testing/fstest"
)

func TestStatus(t *testing.T) {
	type args struct {
		src string
		dst string
	}

	req := require.New(t)

	tests := []struct {
		name string
		args args
	}{
		{name: "simple", args: args{
			src: "/tmp",
			dst: "/src",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := `Sync between "` + tt.args.src + `" and "` + tt.args.dst + `" is running` + "\n"
			actual := Status(tt.args.src, tt.args.dst)
			req.Equal(expected, actual)
		})
	}
}

func Test_checkPath(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockChecker := interfaces.NewMockStater(mockController)

	t.Run("no error", func(t *testing.T) {
		req := require.New(t)
		arg := "tmp"

		testFS := fstest.MapFS{arg: {Mode: fs.ModeDir}}
		st, stErr := testFS.Stat(arg)

		mockChecker.EXPECT().Stat(arg).Return(st, stErr).Times(1)

		err := checkPath(mockChecker, arg)

		req.NoError(err)
	})

	t.Run("not a dir error", func(t *testing.T) {
		req := require.New(t)
		arg := "test.txt"

		testFS := fstest.MapFS{arg: {
			Data: []byte("hello, world"),
		}}

		st, stErr := testFS.Stat(arg)

		mockChecker.EXPECT().Stat(arg).Return(st, stErr).Times(1)

		err := checkPath(mockChecker, arg)

		req.ErrorIs(err, &IsNotADir{arg})
	})

	t.Run("no dir error", func(t *testing.T) {
		req := require.New(t)
		arg := "tmp"

		testFS := fstest.MapFS{arg: {Mode: fs.ModeDir}}
		st, stErr := testFS.Stat("no_dir")

		mockChecker.EXPECT().Stat(arg).Return(st, stErr).Times(1)

		err := checkPath(mockChecker, arg)

		req.ErrorContains(err, "file does not exist")
	})

	t.Run("empty path", func(t *testing.T) {
		req := require.New(t)

		err := checkPath(mockChecker, "")

		req.ErrorIs(err, &EmptyPath{})
	})
}

func Test_checkTheSame(t *testing.T) {
	t.Run("not the same", func(t *testing.T) {
		memFS := fsPkg.NewMemFS()
		req := require.New(t)

		srcDir := "src"
		destDir := "dest"

		errWriteSrcDir := memFS.AFS.Mkdir(srcDir, 0755)
		errWriteDestDir := memFS.AFS.Mkdir(destDir, 0755)

		actual, checkTheSameErr := checkTheSame(memFS, srcDir, destDir)

		req.NoError(errWriteSrcDir)
		req.NoError(errWriteDestDir)
		req.NoError(checkTheSameErr)

		req.Equal(false, actual)
	})

	t.Run("the same", func(t *testing.T) {
		memFS := fsPkg.NewMemFS()
		req := require.New(t)

		srcDir := "src"
		destDir := "src"

		errWriteSrcDir := memFS.AFS.Mkdir(srcDir, 0755)

		actual, checkTheSameErr := checkTheSame(memFS, srcDir, destDir)

		req.NoError(errWriteSrcDir)
		req.NoError(checkTheSameErr)

		req.Equal(true, actual)
	})
}
