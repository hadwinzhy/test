package utils

import "time"

func CurrentDate(captureAt time.Time) time.Time {
	return time.Date(captureAt.Year(), captureAt.Month(), captureAt.Day(), 0, 0, 0, 0, captureAt.Location())
}
