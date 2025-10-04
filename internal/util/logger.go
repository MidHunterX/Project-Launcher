package util

import (
	"fmt"
	"time"
)

const (
	Red   = "\033[1;31m"
	Green = "\033[1;32m"
	Cyan  = "\033[1;36m"
	Reset = "\033[0;0m"
)

func Log(message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s%s%s] %s\n", Green, timestamp, Reset, message)
}
