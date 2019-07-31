package uds

import (
	"errors"
	"fmt"
)

var units = []string{"bytes", "KB", "MB", "GB", "TB"}

type sizeIndex int

const invalid sizeIndex = -1

func calcSize(numOfBytes int64) (size float64, idx sizeIndex) {
	if numOfBytes <= 0 {
		idx = invalid
		return
	}

	step := 1024.
	size = float64(numOfBytes)
	for i := range units {
		idx = sizeIndex(i)
		if i == len(units)-1 || size/step < 1 {
			break
		}
		size /= step
	}
	return
}

func format(numOfBytes int64) (string, error) {
	size, idx := calcSize(numOfBytes)
	if idx == invalid {
		return "", errors.New("size should be positive")
	}

	return fmt.Sprintf("%.1f %s", size, units[idx]), nil
}
