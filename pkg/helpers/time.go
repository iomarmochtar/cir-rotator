package helpers

import (
	"strconv"
	"time"
)

const (
	unitMilliSecond = 1000
)

func ConvertTimeStrToUnix(timeStr string) (time.Time, error) {
	parsed, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(parsed/unitMilliSecond, 0), nil
}

func ConvertTimeStrToReadAble(timeStr string) (string, error) {
	uTime, err := ConvertTimeStrToUnix(timeStr)
	if err != nil {
		return "", err
	}
	return uTime.Format(time.RFC3339), nil
}
