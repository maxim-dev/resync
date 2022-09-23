package helper

/*
Этот хелпер можно вынести в pkg/fs и доработать интерфейсы
Тогда в функциях копирования будет возможность вызывать их из зависимостей (первый параметр с интерфейсом)

Либо, переделать это в отдельную структуру, у которой будет свойство path и два метода IsDirExist и MakeDir.
Тогда эти методы не будут принимать параметры
*/

import (
	"fmt"
	"os"
)

func IsDirExist(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func MakeDir(path string) error {
	if IsDirExist(path) {
		return nil
	}

	err := os.MkdirAll(path, 0755)

	if err != nil {
		return fmt.Errorf("ошибка при создании директории %q: %s", path, err)
	}

	return nil
}
