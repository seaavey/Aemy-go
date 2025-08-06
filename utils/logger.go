package utils

import (
	"fmt"
	"time"
);

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

func timestamp() string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func Info(msg string) {
	fmt.Printf("%s[%s] [INFO] %s%s\n", Green, timestamp(), msg, Reset)
}

func Error(msg string) {
	fmt.Printf("%s[%s] [ERROR] %s%s\n", Red, timestamp(), msg, Reset)
}

func Warn(msg string) {
	fmt.Printf("%s[%s] [WARN] %s%s\n", Yellow, timestamp(), msg, Reset)
}

func Debug(msg string) {
	fmt.Printf("%s[%s] [DEBUG] %s%s\n", Cyan, timestamp(), msg, Reset)
}
