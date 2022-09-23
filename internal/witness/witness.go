package witness

import (
	"io"
	"log"
	"resync/pkg/interfaces"
)

// Имя файла, куда всё логируется можем задать константой, потому что в требованиях зафиксировано определенное имя
const logFileName = "log.txt"

type Witness struct {
	Output io.WriteCloser // Поле необходимо в структуре для того, чтобы можно было сделать отложенный вызов Close()
}

func NewWitness(appender interfaces.Appender) *Witness {
	writer, err := appender.Append(logFileName)

	if err == nil {
		log.Printf("The logger Output is now set to the %q file.", logFileName)
		log.SetOutput(writer)
	} else {
		log.Printf("Can't set logger Output to the %q file (error: %s). will log to stdOut", logFileName, err)
	}

	return &Witness{Output: writer}
}

// Println обертка вокруг стандартного GoLang логгера. Ресивер этому методу не нужен,
// чтобы была возможность вызывать как обычный Go log без привязки к структуре
func Println(v ...interface{}) {
	log.Println(v...)
}
