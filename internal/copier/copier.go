package copier

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"resync/internal/helper"
	"resync/internal/witness"
	"resync/pkg/interfaces"
)

// Copy копируем файл src в директорию destDir
func Copy(o interfaces.StatOpenCreater, w interfaces.Workable, src, destDir string) (written int64, err error) {
	var (
		source      io.ReadCloser
		destination io.WriteCloser
		infoSrc     fs.FileInfo
	)

	if w.IsInProgress(src) {
		return 0, nil
	}

	infoSrc, err = o.Stat(src)

	// Проверяем источник. Все, что не относится к обычным файлам (например, directories,
	// symlinks, devices, etc.) не сможем скопировать и вернем ошибку
	if err != nil {
		return 0, err
	} else if mode := infoSrc.Mode(); !mode.IsRegular() {
		return 0, fmt.Errorf("the source file %q is not a regular file (%q)", infoSrc.Name(), mode.String())
	}

	w.WorkingOn(src)
	defer w.WorkingDone(src)

	// Строим путь до файла, который будет скопирован.
	// Используем встроенную функцию, чтобы учесть оба варианта, со слешем на конце и без
	dest := filepath.Join(destDir, infoSrc.Name())

	// Проверяем финальный файл. Если это директория или это один и тот же файл, ничего не копируем
	if infoDest, err := o.Stat(dest); err == nil && (infoDest.IsDir() || os.SameFile(infoSrc, infoDest) || helper.IsFilesEqual(src, dest)) {
		return 0, nil
	}

	if source, err = o.Open(src); err != nil {
		return 0, err
	}

	defer func() {
		scopeErr := source.Close()
		if err == nil {
			err = scopeErr
		}
	}()

	if destination, err = o.Create(dest); err != nil {
		return 0, err
	}

	defer func() {
		scopeErr := destination.Close()
		if err == nil {
			err = scopeErr
		}
	}()

	written, err = io.Copy(destination, source)

	if err != nil {
		witness.Println("Copy:", src, "->", dest, "failed with error:", err)
	} else {
		witness.Println("Copy:", src, "->", dest, "written: ", helper.BytesToReadable(written))
		os.Chmod(dest, infoSrc.Mode())
	}

	return written, err
}

// CloneDir рекурсивно копируем src директорию в destDir.
func CloneDir(o interfaces.StatOpenCreater, w interfaces.Workable, src, destDir string) error {
	witness.Println("CloneDir:", src, "->", destDir)

	items, err := ioutil.ReadDir(src)

	if err != nil {
		return err
	}

	for _, item := range items {
		srcPath := filepath.Join(src, item.Name())
		destPath := filepath.Join(destDir, item.Name())

		fileInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		// Поддиректорию создаем и продолжаем в рекурсии. Обычный файл просто копируем
		if fileInfo.IsDir() {
			if err := helper.MakeDir(destPath); err != nil {
				return err
			}
			CloneDir(o, w, srcPath, destPath)
		} else {
			Copy(o, w, srcPath, destDir)
		}

	}

	return nil
}

func WatchFor(o interfaces.StatOpenCreater, w interfaces.Workable, src, dest string) error {
	items, err := ioutil.ReadDir(src)

	if err != nil {
		return err
	}

	for _, item := range items {
		srcPath := filepath.Join(src, item.Name())
		destPath := filepath.Join(dest, item.Name())

		// если item это директория, то
		// -- целевой директории нет – клонируем целиком
		// -- такая уже есть – заходим в рекурсию
		// если item это файл, копируем его через Copy()

		if item.IsDir() {
			if !helper.IsDirExist(destPath) {
				helper.MakeDir(destPath)
				CloneDir(o, w, srcPath, destPath)
				continue
			}
			WatchFor(o, w, srcPath, destPath)
		} else {
			Copy(o, w, srcPath, dest)
		}
	}

	return nil
}
