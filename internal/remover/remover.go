package remover

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"resync/internal/witness"
	"resync/pkg/interfaces"
)

func Remove(o interfaces.StatRemover, w interfaces.Workable, suspect, root, dest string) (isDeleted bool, err error) {

	suspectPath := filepath.Join(dest, suspect)

	if w.IsInProgress(suspectPath) {
		return false, nil
	}

	w.WorkingOn(suspectPath)
	defer w.WorkingDone(suspectPath)

	if exist, err := exists(o, filepath.Join(root, suspect)); err != nil {
		return false, err
	} else if exist {
		return false, nil
	} else if err := o.RemoveAll(suspectPath); err != nil {
		return false, err
	}

	witness.Println("has been removed:", suspectPath)

	return true, nil
}

// exists проверяет наличие файла или директории
// замечание: если функция понадобится за пределами пакета, то ее необходимо вынести в интерфейс файловой системы
// и реализовать методы в каждой из структур интерфейса ФС
func exists(o interfaces.Stater, path string) (bool, error) {
	_, err := o.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WatchFor(o interfaces.StatRemover, w interfaces.Workable, path, root string) error {
	items, err := ioutil.ReadDir(path)

	if err != nil {
		return err
	}

	for _, item := range items {
		isDeleted, err := Remove(o, w, item.Name(), root, path)

		if err != nil {
			return err
		}

		if item.IsDir() && !isDeleted {
			WatchFor(o, w, filepath.Join(path, item.Name()), filepath.Join(root, item.Name()))
		}
	}

	return nil
}
