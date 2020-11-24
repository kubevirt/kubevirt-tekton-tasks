package utils

import "time"

func WithTimeout(originalTimeout time.Duration) (runWithTimeout func(callback func(timeout time.Duration, finished bool))) {
	start := time.Now()

	return func(callback func(timeout time.Duration, finished bool)) {
		// special case - run indefinitely
		if originalTimeout <= 0 {
			callback(originalTimeout, false)
			return
		}
		elapsedDuration := time.Since(start)
		newTimeout := originalTimeout - elapsedDuration

		callback(newTimeout, newTimeout <= 0)
	}
}
