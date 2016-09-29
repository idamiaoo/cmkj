package util

import (
	"time"
)

const (
	Layout = "2006-01-02 15:04:05.999"
)

func NewDate(layout string) string {
	return time.Now().Format(layout)
}
