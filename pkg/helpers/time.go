package helpers

import (
	"strconv"
	"time"
)

const (
	UnitMilliSecond = 1000
)

func ConvertTimeStrToUnix(timeStr string, unit int64) (time.Time, error) {
	parsed, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(parsed/unit, 0), nil
}

func ConvertTimeStrToReadAble(timeStr string, unit int64) (string, error) {
	uTime, err := ConvertTimeStrToUnix(timeStr, unit)
	if err != nil {
		return "", err
	}
	return uTime.Format(time.RFC3339), nil
}
