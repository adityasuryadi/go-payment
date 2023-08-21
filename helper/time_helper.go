package helper

import "time"

func GetNowStringFormat() string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return now
}
