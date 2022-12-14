package fs_struct

import (
	"os"
)

// Пример абстракции файловой системы на структурах
// Здесь определена структура абстрактной файловой системы, которая содержит поля типа "функция"
// без конкретной реализации
// Предполагается, что на этапе инициализации структуры соответствующие поля будут заданы
//
// В целом это менее удобный подход, чем абстракция на интерфейсах, потому что сигнатуры функций могут не совпадать
// и придется выполнять приведение типов
// см. disk.go и mem.go

type AbsFS struct {
	Stat   func(string) (os.FileInfo, error)
	Open   func(string) (*os.File, error)
	Create func(string) (*os.File, error)
}
