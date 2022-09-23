package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/eiannone/keyboard"
	"resync/internal/concierge"
	"resync/internal/copier"
	"resync/internal/dispenser"
	"resync/internal/remover"
	"resync/internal/witness"
	"resync/pkg/fs"
	"resync/pkg/interfaces"
	"time"
)

var (
	src, dest string
	fileSys   interfaces.Fs
	logger    *witness.Witness
	copierD   *dispenser.WorkingList
	removerD  *dispenser.WorkingList
)

func init() {
	flag.Parse()

	src, dest = flag.Arg(0), flag.Arg(1)

	fileSys = fs.NewDiskFS()

	// Вызываем NewWitness() для того, чтобы задать Output у log пакета
	// Сама переменная используется только для того, чтобы в отложенном вызове закрыть Writer после
	// завершения работы основной функции
	logger = witness.NewWitness(fileSys)

	copierD = dispenser.NewWorkingList()
	removerD = dispenser.NewWorkingList()
}

func main() {
	from, to, err := concierge.Obtain(fileSys, src, dest)

	if err != nil {
		panic(err)
	}

	fmt.Println(concierge.Status(from, to))
	fmt.Println("Press ESC to quit")

	keysEvents, err := keyboard.GetKeys(10)
	defer func() {
		_ = keyboard.Close()
	}()

	defer logger.Output.Close()

	watchForTicker := time.NewTicker(2 * time.Second)
	watchForCtx, cancelWatchFor := context.WithCancel(context.Background())

mainLoop:
	for {
		select {
		case <-watchForTicker.C:
			go remover.WatchFor(fileSys, removerD, to, from)
			go copier.WatchFor(fileSys, copierD, from, to)
		case <-watchForCtx.Done():
			watchForTicker.Stop()
			break mainLoop
		case event := <-keysEvents:
			if event.Key == keyboard.KeyEsc {
				break mainLoop
			}
		}
	}

	fmt.Println("\nend")

	cancelWatchFor()

	if err != nil {
		panic(err)
	}
}
