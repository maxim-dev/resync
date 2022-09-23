package helper

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func IsFilesEqual(src, dest string) bool {
	destHash, err := getHashFor(dest)

	if err != nil {
		return false
	}

	srcHash, err := getHashFor(src)

	if err != nil {
		return false
	}

	// так как хеш возвращает слайс байтов, используем специальную функцию для сравнения слайсов
	i := bytes.Compare(srcHash, destHash)

	if i == 0 {
		return true
	} else {
		return false
	}
}

func getHashFor(path string) ([]byte, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Читаем файл чанками, для этого создаем буффер
	buffer := make([]byte, 30*1024)
	h := md5.New()

	for {
		n, err := f.Read(buffer)

		if n > 0 {
			_, err := h.Write(buffer[:n])
			if err != nil {
				return nil, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	sum := h.Sum(nil)

	return sum, nil
}

func BytesToReadable(b int64) string {
	const base = 1024
	if b < base {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(base), 0

	for n := b / base; n >= base; n /= base {
		div *= base
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
