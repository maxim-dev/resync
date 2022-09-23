package concierge

import (
	"fmt"
	"os"
	"resync/internal/witness"
	"resync/pkg/interfaces"
)

func greet() {
	fmt.Println("ReSync is a program to make two directories synchronized")
	fmt.Println("Usage: resync SRC DEST")
	os.Exit(1)
}

func Status(src, dest string) string {
	status := fmt.Sprintf("Sync between %q and %q is running\n", src, dest)
	witness.Println(status)
	return status
}

func Obtain(checker interfaces.Stater, src, dest string) (string, string, error) {

	if src == "" || dest == "" {
		greet()
	}

	if err := checkPath(checker, src); err != nil {
		return "", "", err
	}

	if err := checkPath(checker, dest); err != nil {
		return "", "", err
	}

	//Проверяем, что src и dest это разные директории. Если это одна и та же директория – вернуть ошибку
	isTheSame, err := checkTheSame(checker, src, dest)

	if isTheSame {
		return "", "", &TheSamePath{}
	} else if err != nil {
		return "", "", err
	}

	return src, dest, nil
}

func checkTheSame(checker interfaces.Stater, src, dest string) (bool, error) {
	if src == dest {
		return true, nil
	}
	infoSrc, err := checker.Stat(src)

	if err != nil {
		return false, err
	}

	infoDest, err := checker.Stat(dest)

	if err != nil {
		return false, err
	}

	return os.SameFile(infoSrc, infoDest), nil
}

func checkPath(checker interfaces.Stater, path string) error {
	if path == "" {
		return &EmptyPath{}
	}

	info, err := checker.Stat(path)

	if err != nil {
		return fmt.Errorf("failed to open directory, error: %w", err)
	}

	if !info.IsDir() {
		return &IsNotADir{info.Name()}
	}

	return nil
}
