package util

import (
	"time"
)

func GetNowSecondTimestamp() int64 {
	return time.Now().Unix()
}

func GetNowMillisecondTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetNowNanosecondTimestamp() int64 {
	return time.Now().UnixNano()
}

func GetNowFormatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func StrToFormatDate(input string) string {
	parsedTime, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return ""
	}

	return parsedTime.Format("2006-01-02")
}
