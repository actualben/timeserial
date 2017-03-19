package main

import (
	"fmt"
	"time"
)

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

// TimeSerial returns the timeserial string for the given time.Time
func TimeSerial(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d%02d", t.Year(), t.Month(), t.Day(),
		min(((t.Hour()*3600)+(t.Minute()*60)+t.Second())/864, 99))
}
