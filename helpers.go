package rrkurrentdb

import "time"

func deadlineMsToDuration(ms *int64) *time.Duration {
	if ms == nil {
		return nil
	}
	d := time.Duration(*ms) * time.Millisecond
	return &d
}
