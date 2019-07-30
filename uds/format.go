package uds

import (
	"errors"
	"fmt"
)

func format(numOfBytes int64) (string, error) {
	if numOfBytes <= 0 {
		return "", errors.New("size should be positive")
	}

	step := 1024.

	size := float64(numOfBytes)
	units := []string{"bytes", "KB", "MB", "GB", "TB"}

	var unit string
	for i, u := range units {
		unit = u
		if i == len(units)-1 || size/step < 1 {
			break
		}
		size /= step
	}

	return fmt.Sprintf("%.1f %s", size, unit), nil
}
