package config

import (
	"os"
	"strconv"
)

func ProgressThreshold() float64 {
	const defaultProgressThreshold = 5.0

	if v := os.Getenv("PROGRESS_THRESHOLD"); v != "" {
		if t, err := strconv.ParseFloat(v, 64); err == nil && t >= 1 && t <= 100 {
			return t
		}
	}
	return defaultProgressThreshold
}
