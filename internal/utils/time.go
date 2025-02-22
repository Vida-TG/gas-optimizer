package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	h := d / time.Hour
	d = d % time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh%02dm", h, m)
} 